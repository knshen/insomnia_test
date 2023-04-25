package checkers

import (
	"code.sk.org/insomnia_test/client"
	"code.sk.org/insomnia_test/consts"
	"code.sk.org/insomnia_test/model"
	"context"
)

type ProjectPermChecker struct {
	authClient client.IAuthClient
}

func NewProjectPermChecker(authClient client.IAuthClient) *ProjectPermChecker {
	return &ProjectPermChecker{
		authClient: authClient,
	}
}

func (c *ProjectPermChecker) IsValid(ctx context.Context, rule *model.LintRuleRequest) error {
	hasPerm, err := c.authClient.HasProjectPermission(ctx, rule.UserToken, rule.ProjectID)
	if err != nil {
		return err
	}

	if !hasPerm {
		return consts.NoProjectPerm
	}
	return nil
}
