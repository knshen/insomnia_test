package client

import (
	"code.sk.org/insomnia_test/consts"
	"code.sk.org/insomnia_test/utils"
	"context"
	"time"
)

type IAuthClient interface {
	GetUserRole(ctx context.Context, token string, orgID int64) (consts.Role, error)
	GetPermitProjectIDs(ctx context.Context, token string, limit, offset int64) ([]int64, error)
	HasProjectPermission(ctx context.Context, token string, projectID int64) (bool, error)
}

type DefaultAuthClient struct {
}

func (c *DefaultAuthClient) GetUserRole(ctx context.Context, token string, orgID int64) (consts.Role, error) {
	// TODO: call rpc method of auth service
	time.Sleep(time.Millisecond * 500)
	if token == "admin_token" {
		return consts.AdminRole, nil
	}
	return consts.NonAdminRole, nil
}

func (c *DefaultAuthClient) GetPermitProjectIDs(ctx context.Context, token string, limit, offset int64) ([]int64, error) {
	// TODO: call rpc method of auth service
	time.Sleep(time.Millisecond * 500)
	if token == "admin_token" {
		return []int64{1001, 1002, 1003}, nil
	}
	return []int64{1001}, nil
}

func (c *DefaultAuthClient) HasProjectPermission(ctx context.Context, token string, projectID int64) (bool, error) {
	var limit int64 = 10
	var offset int64 = 0

	for {
		pids, err := c.GetPermitProjectIDs(ctx, token, limit, offset)
		if err != nil {
			return false, err
		}

		if utils.Int64SliceContainEle(pids, projectID) {
			return true, nil
		}

		if len(pids) < int(limit) {
			break
		}
		offset += limit
	}

	return false, nil
}
