package model

type LintRule struct {
	RuleID        int64
	RawRuleString string
}

func (r *LintRule) ToKVModel() map[string]