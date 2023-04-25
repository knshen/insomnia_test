package checkers

import (
	"code.sk.org/insomnia_test/model"
	"context"
)

type ILintChecker interface {
	IsValid(ctx context.Context, rule *model.LintRuleRequest) error
}

type CompositeLintChecker struct {
	checkers []ILintChecker
}

func (c *CompositeLintChecker) IsValid(ctx context.Context, rule *model.LintRuleRequest) error {
	for _, ch := range c.checkers {
		if err := ch.IsValid(ctx, rule); err != nil {
			return err
		}
	}
	return nil
}

func (c *CompositeLintChecker) AppendChecker(ch ILintChecker) {
	if c.checkers == nil {
		c.checkers = make([]ILintChecker, 0)
	}
	c.checkers = append(c.checkers, ch)
}
