package checkers

import (
	"code.sk.org/insomnia_test/consts"
	"code.sk.org/insomnia_test/model"
	"context"
)

type LintIDsChecker struct {
}

func (c *LintIDsChecker) IsValid(ctx context.Context, rule *model.LintRuleRequest) error {
	if rule.ProjectID <= 0 {
		return consts.ProjectIDEmpty
	}

	return nil
}
