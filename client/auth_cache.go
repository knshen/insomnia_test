package client

import (
	"code.sk.org/insomnia_test/consts"
	"code.sk.org/insomnia_test/dal/kv"
	"code.sk.org/insomnia_test/model"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type AuthCacheClient struct {
	realClient IAuthClient

	ttl    time.Duration
	delTTL time.Duration
}

func NewAuthCacheClient(realClient IAuthClient, ttl, delTTL time.Duration) IAuthClient {
	if delTTL.Seconds() < ttl.Seconds() {
		delTTL = ttl
	}
	return &AuthCacheClient{
		realClient: realClient,
		ttl:        ttl,
		delTTL:     delTTL,
	}
}

func (c *AuthCacheClient) GetUserRole(ctx context.Context, token string, orgID int64) (consts.Role, error) {
	key := fmt.Sprintf("role:%v:%v", orgID, token)

	raw, rErr := kv.GetRedisClient().Get(ctx, key).Result()
	var loadRole consts.Role
	var err error
	needUpdateCache := false

	defer func() {
		if needUpdateCache {
			go func() {
				if loadRole == consts.Unknown {
					loadRole, err = c.realClient.GetUserRole(ctx, token, orgID)
				}
				if loadRole != consts.Unknown {
					entry := &model.CacheEntry{
						Key:       key,
						Value:     strconv.FormatInt(int64(loadRole), 10),
						CreatedAt: time.Now().Unix(),
					}

					kv.GetRedisClient().Set(ctx, key, entry.Dump(), c.delTTL)
				}
			}()
		}
	}()

	if rErr == redis.Nil {
		// cache not exist
		loadRole, err = c.realClient.GetUserRole(ctx, token, orgID)
		if err == nil {
			needUpdateCache = true
		}
		return loadRole, err
	}

	if rErr != nil {
		// cache read fail: let us fail fast
		return consts.Unknown, rErr
	}

	entry := model.NewCacheFromRawString(raw)
	if time.Now().After(time.Unix(entry.CreatedAt, 0).Add(c.ttl)) {
		// check if logically expire, if so return expired one and update cache async
		needUpdateCache = true
	}

	role, _ := strconv.ParseInt(entry.Value, 10, 64)
	return consts.Role(role), nil
}

func (c *AuthCacheClient) GetPermitProjectIDs(ctx context.Context, token string, limit, offset int64) ([]int64, error) {
	return c.realClient.GetPermitProjectIDs(ctx, token, limit, offset)
}

func (c *AuthCacheClient) HasProjectPermission(ctx context.Context, token string, projectID int64) (bool, error) {
	key := fmt.Sprintf("project_perm:%v:%v", projectID, token)

	raw, rErr := kv.GetRedisClient().Get(ctx, key).Result()
	var hasPerm *bool
	var err error
	needUpdateCache := false

	defer func() {
		if needUpdateCache {
			go func() {
				if hasPerm == nil {
					perm := false
					perm, err = c.realClient.HasProjectPermission(ctx, token, projectID)
					if err == nil {
						hasPerm = &perm
					}
				}
				if hasPerm != nil {
					entry := &model.CacheEntry{
						Key:       key,
						Value:     strconv.FormatBool(*hasPerm),
						CreatedAt: time.Now().Unix(),
					}
					kv.GetRedisClient().Set(ctx, key, entry.Dump(), c.delTTL)
				}
			}()
		}
	}()

	if rErr == redis.Nil {
		// cache not exist
		perm := false
		perm, err = c.realClient.HasProjectPermission(ctx, token, projectID)
		if err == nil {
			hasPerm = &perm
			needUpdateCache = true
		}
		return perm, err
	}

	if rErr != nil {
		// cache read fail: let us fail fast
		return false, rErr
	}

	entry := model.NewCacheFromRawString(raw)
	if time.Now().After(time.Unix(entry.CreatedAt, 0).Add(c.ttl)) {
		// check if logically expire, if so return expired one and update cache async
		needUpdateCache = true
	}

	perm, _ := strconv.ParseBool(entry.Value)
	return perm, nil
}
