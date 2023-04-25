package checkers

import (
	"code.sk.org/insomnia_test/consts"
	"code.sk.org/insomnia_test/model"
	"context"
	"gopkg.in/yaml.v2"
)

type LintYamlContentChecker struct {
}

// check if lint content is valid
func (c *LintYamlContentChecker) IsValid(ctx context.Context, rule *model.LintRuleRequest) error {
	if rule.RawRuleString == "" {
		return consts.YamlNotValid
	}
	err := yaml.Unmarshal([]byte(rule.RawRuleString), map[string]interface{}{})
	if err != nil {
		return consts.YamlNotValid
	}
	return nil
}
