package model

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type LintRule struct {
	RuleID        int64  `json:"rule_id"`
	OrgID         int64  `json:"org_id"`
	ProjectID     int64  `json:"project_id"`
	RawRuleString string `json:"raw_rule_string"`
}

func (r *LintRule) Dump() string {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(r); err != nil {
		return ""
	}

	return b.String()
}

func NewLintRuleFromRaw(raw string) *LintRule {
	var r LintRule
	json.Unmarshal([]byte(raw), &r)
	return &r
}

type LintRuleRequest struct {
	*LintRule

	UserToken string
}

func GetLintRuleFromReq(req *http.Request) (*LintRuleRequest, error) {
	token := req.Header.Get("user_token")

	var r LintRule
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	return &LintRuleRequest{
		LintRule:  &r,
		UserToken: token,
	}, nil
}
