package consts

import "errors"

var (
	NoOrgAdmin    = errors.New("not org admin")
	NoProjectPerm = errors.New("no project permission")

	ProjectIDEmpty     = errors.New("project id is empty")
	ProjectNotExist    = errors.New("project not exist")
	RequestBodyIllegal = errors.New("request body illegal")

	LintingRuleAlreadyExist = errors.New("linting rule of this project already exist")
	NoLintingRuleBind       = errors.New("there is no linting rule for this project")

	YamlNotValid = errors.New("yaml content not valid")

	BatchSizeTooLarge = errors.New("cannot query more than 10 rules one time")

	HttpMethodNotAllow = errors.New("http method not allow")

	Error2Code = map[error]int32{
		NoOrgAdmin:              403,
		NoProjectPerm:           403,
		LintingRuleAlreadyExist: 403,
		NoLintingRuleBind:       404,
		ProjectIDEmpty:          400,
		YamlNotValid:            400,
		HttpMethodNotAllow:      405,
		ProjectNotExist:         400,
	}
)
