package checkers

import (
	"code.sk.org/insomnia_test/client"
	"code.sk.org/insomnia_test/consts"
	"code.sk.org/insomnia_test/model"
	"context"
)

type ProjectExistChecker struct {
	projectClient client.IProjectClient
}

func NewProjectExistChecker(projectClient client.IProjectClient) *ProjectExistChecker {
	return &ProjectExistChecker{
		projectClient: projectClient,
	}
}

func (c *ProjectExistChecker) IsValid(ctx context.Context, rule *model.LintRuleRequest) error {
	p, err := c.projectClient.GetByProjectID(ctx, rule.ProjectID)
	if err != nil {
		return err
	}
	if p == nil {
		return consts.ProjectNotExist
	}
	rule.OrgID = p.OrgID
	return nil
}
