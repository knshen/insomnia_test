package client

import (
	"code.sk.org/insomnia_test/consts"
	"context"
)

type AuthCacheClient struct {
	realClient IAuthClient
}

func (c *AuthCacheClient) GetUserRole(ctx context.Context, token string, orgID int64) (consts.Role, error) {

}

func (c *AuthCacheClient) GetPermitProjectIDs(ctx context.Context, token string, limit, offset int64) ([]int64, error) {
	return c.realClient.GetPermitProjectIDs(ctx, token, limit, offset)
}

func (c *AuthCacheClient) HasProjectPermission(ctx context.Context, token string, orgID, projectID int64) (bool, error) {

}
