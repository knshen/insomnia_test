package checkers

import (
	"code.sk.org/insomnia_test/client"
	"code.sk.org/insomnia_test/consts"
	"code.sk.org/insomnia_test/model"
	"context"
)

type OrgAdminChecker struct {
	authClient client.IAuthClient
}

func NewOrgAdminChecker(authClient client.IAuthClient) *OrgAdminChecker {
	return &OrgAdminChecker{
		authClient: authClient,
	}
}

func (c *OrgAdminChecker) IsValid(ctx context.Context, rule *model.LintRuleRequest) error {
	role, err := c.authClient.GetUserRole(ctx, rule.UserToken, rule.OrgID)
	if err != nil {
		return err
	}
	if role != consts.AdminRole {
		return consts.NoOrgAdmin
	}
	return nil
}
