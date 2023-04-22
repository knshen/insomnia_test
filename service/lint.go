package service

import (
	"code.sk.org/insomnia_test/model"
	"context"
)

// core service to manage lint rules
type LintService struct {
}

func (s *LintService) CreateLintRule(ctx context.Context, orgID, projectID int64, rawLint string) error {
	// TODO
	return nil
}

func (s *LintService) UpdateLintRule(ctx context.Context, orgID, projectID int64, rawLint string) error {
	// TODO
	return nil
}

func (s *LintService) GetLintRule(ctx context.Context, orgID, projectID int64) (string, error) {
	return "", nil
}

func (s *LintService) MGetLintRule(ctx context.Context, param *model.MGetLintRuleParam) (map[int64]string, error) {
	return nil, nil
}
