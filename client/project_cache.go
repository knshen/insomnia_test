package client

import (
	"code.sk.org/insomnia_test/dal/kv"
	"code.sk.org/insomnia_test/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type ProjectCacheClient struct {
	realClient IProjectClient

	ttl    time.Duration
	delTTL time.Duration
}

func NewProjectCacheClient(realClient IProjectClient, ttl, delTTL time.Duration) IProjectClient {
	if delTTL.Seconds() < ttl.Seconds() {
		delTTL = ttl
	}
	return &ProjectCacheClient{
		realClient: realClient,
		ttl:        ttl,
		delTTL:     delTTL,
	}
}

func (c *ProjectCacheClient) GetByProjectID(ctx context.Context, projectID int64) (*model.Project, error) {
	key := fmt.Sprintf("project:%v", projectID)

	raw, rErr := kv.GetRedisClient().Get(ctx, key).Result()
	var project *model.Project
	var err error
	needUpdateCache := false

	defer func() {
		if needUpdateCache {
			go func() {
				if project == nil {
					project, err = c.realClient.GetByProjectID(ctx, projectID)
				}
				if err == nil {
					projectVal, _ := json.Marshal(project)
					entry := &model.CacheEntry{
						Key:       key,
						Value:     string(projectVal),
						CreatedAt: time.Now().Unix(),
					}

					kv.GetRedisClient().Set(ctx, key, entry.Dump(), c.delTTL)
				}
			}()
		}
	}()

	if rErr == redis.Nil {
		// cache not exist
		project, err = c.realClient.GetByProjectID(ctx, projectID)
		if err == nil {
			needUpdateCache = true
		}
		return project, err
	}

	if rErr != nil {
		// cache read fail: let us fail fast
		return nil, rErr
	}

	entry := model.NewCacheFromRawString(raw)
	if time.Now().After(time.Unix(entry.CreatedAt, 0).Add(c.ttl)) {
		// check if logically expire, if so return expired one and update cache async
		needUpdateCache = true
	}

	if entry.Value == "null" {
		return nil, nil
	}

	var ret model.Project
	json.Unmarshal([]byte(entry.Value), &ret)
	return &ret, nil
}
