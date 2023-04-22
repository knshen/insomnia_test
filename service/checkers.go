package service

import "context"

type ILintChecker interface {
	IsValid(ctx context.Context, orgID, projectID int64, rawRuleString string) error
}

type LintUpdatePermChecker struct {
}

func (c *LintUpdatePermChecker) IsValid(ctx context.Context, orgID, projectID int64, rawRuleString string) error {

}
