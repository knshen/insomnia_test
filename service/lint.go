package service

import (
	"code.sk.org/insomnia_test/client"
	"code.sk.org/insomnia_test/consts"
	"code.sk.org/insomnia_test/dal/kv"
	"code.sk.org/insomnia_test/model"
	"code.sk.org/insomnia_test/service/checkers"
	"context"
)

// core service to manage lint rules
type LintService struct {
	dal         kv.ILintRuleDal
	lintChecker checkers.ILintChecker
	idGen       client.IIDGeneratClient
}

func NewLintService(dal kv.ILintRuleDal, checker checkers.ILintChecker, idGen client.IIDGeneratClient) *LintService {
	return &LintService{
		dal:         dal,
		lintChecker: checker,
		idGen:       idGen,
	}
}

func (s *LintService) CreateLintRule(ctx context.Context, rule *model.LintRuleRequest) error {
	if err := s.lintChecker.IsValid(ctx, rule); err != nil {
		return err
	}

	ids, err := s.idGen.GenInt64ID(ctx, 1)
	if err != nil {
		return err
	}
	rule.RuleID = ids[0]

	return s.dal.CreateLintRule(ctx, rule.LintRule)
}

func (s *LintService) UpdateLintRule(ctx context.Context, rule *model.LintRuleRequest) error {
	if err := s.lintChecker.IsValid(ctx, rule); err != nil {
		return err
	}

	if rule.LintRule.RuleID == 0 {
		ids, err := s.idGen.GenInt64ID(ctx, 1)
		if err != nil {
			return err
		}
		rule.RuleID = ids[0]
	}

	s.dal.UpdateLintRule(ctx, rule.LintRule)
	return nil
}

func (s *LintService) GetLintRule(ctx context.Context, rule *model.LintRuleRequest) (*model.LintRule, error) {
	if err := s.lintChecker.IsValid(ctx, rule); err != nil {
		return nil, err
	}

	r, err := s.dal.GetLintRule(ctx, rule.LintRule.OrgID, rule.LintRule.ProjectID)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, consts.NoLintingRuleBind
	}
	return r, nil
}
