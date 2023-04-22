package client

import (
	"code.sk.org/insomnia_test/consts"
	"context"
)

type IAuthClient interface {
	GetUserRole(ctx context.Context, token string, orgID int64) (consts.Role, error)
	GetPermitProjectIDs(ctx context.Context, token string, limit, offset int64) ([]int64, error)
	HasProjectPermission(ctx context.Context, token string, orgID, projectID int64) (bool, error)
}

type DefaultAuthClient struct {
}

func (c *DefaultAuthClient) GetUserRole(ctx context.Context, token string, orgID int64) (consts.Role, error) {
	// TODO call rpc method of auth service

	return consts.AdminRole, nil
}

func (c *DefaultAuthClient) GetPermitProjectIDs(ctx context.Context, token string, limit, offset int64) ([]int64, error) {
	// TODO call rpc method of auth service

	return nil, nil
}

func (c *DefaultAuthClient) HasProjectPermission(ctx context.Context, token string, orgID, projectID int64) (bool, error) {
	// TODO call rpc method of auth service

	return true, nil
}
