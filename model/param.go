package model

import (
	"code.sk.org/insomnia_test/consts"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type MGetLintRuleParam struct {
	OrgIDs     []int64 `json:"org_ids"`
	ProjectIDs []int64 `json:"project_ids"`
}

func (p *MGetLintRuleParam) Check() error {
	if len(p.OrgIDs) == 0 {
		return consts.OrgIDListEmpty
	}
	if len(p.ProjectIDs) == 0 {
		return consts.ProjectIDListEmpty
	}

	if len(p.OrgIDs) != len(p.ProjectIDs) {
		return consts.OrgIDListProjectIDListNotMatch
	}

	return nil
}

func NewMGetLintRuleParamFromReq(req *http.Request) *MGetLintRuleParam {
	body, _ := ioutil.ReadAll(req.Body)
	var p MGetLintRuleParam
	json.Unmarshal(body, &p)
	return &p
}
