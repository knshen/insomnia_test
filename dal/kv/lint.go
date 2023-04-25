package kv

import (
	"code.sk.org/insomnia_test/consts"
	"code.sk.org/insomnia_test/model"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

const (
	LintRuleKeyForm = "linting:%v:%v"
)

type ILintRuleDal interface {
	CreateLintRule(ctx context.Context, rule *model.LintRule) error
	UpdateLintRule(ctx context.Context, rule *model.LintRule) error
	GetLintRule(ctx context.Context, orgID, projectID int64) (*model.LintRule, error)
	MGetLintRule(ctx context.Context, orgIDs, projectIDs []int64) (map[int64]*model.LintRule, error)
}

type RedisLintRuleDal struct {
}

func (o *RedisLintRuleDal) CreateLintRule(ctx context.Context, rule *model.LintRule) error {
	key := o.buildKey(rule.OrgID, rule.ProjectID)

	// never expire
	ok, err := redisCli.SetNX(ctx, key, rule.Dump(), 0).Result()
	if err != nil {
		return err
	}
	if !ok {
		return consts.LintingRuleAlreadyExist
	}

	return nil
}

func (o *RedisLintRuleDal) UpdateLintRule(ctx context.Context, rule *model.LintRule) error {
	// do upsert operation(create if not exist, otherwise update)
	key := o.buildKey(rule.OrgID, rule.ProjectID)
	if err := o.CreateLintRule(ctx, rule); err != nil {
		if err != consts.LintingRuleAlreadyExist {
			return err
		}

		// update
		existRule, err := o.GetLintRule(ctx, rule.OrgID, rule.ProjectID)
		if err != nil {
			return err
		}
		existRule.RawRuleString = rule.RawRuleString
		_, err = redisCli.Set(ctx, key, existRule.Dump(), 0).Result()
		return err
	}

	// create success
	return nil
}

func (o *RedisLintRuleDal) GetLintRule(ctx context.Context, orgID, projectID int64) (*model.LintRule, error) {
	key := o.buildKey(orgID, projectID)
	raw, err := redisCli.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	return model.NewLintRuleFromRaw(raw), nil
}

func (o *RedisLintRuleDal) MGetLintRule(ctx context.Context, orgIDs, projectIDs []int64) (map[int64]*model.LintRule, error) {
	keys := make([]string, 0, len(orgIDs))
	for i := 0; i < len(orgIDs); i++ {
		keys = append(keys, o.buildKey(orgIDs[i], projectIDs[i]))
	}
	res, err := redisCli.MGet(ctx, keys...).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	ret := make(map[int64]*model.LintRule, len(orgIDs))
	for _, r := range res {
		if r == "" {
			continue
		}

		if v, ok := r.(string); ok {
			rule := model.NewLintRuleFromRaw(v)
			ret[rule.ProjectID] = rule
		}

	}

	return ret, nil
}

func (o *RedisLintRuleDal) buildKey(orgID, projectID int64) string {
	return fmt.Sprintf(LintRuleKeyForm, orgID, projectID)
}
