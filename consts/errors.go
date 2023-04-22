package consts

import "errors"

var (
	NoLintUpdatePermError = errors.New("no permission to update linting rule")
	NoLintCreatePerm      = errors.New("no permission to create linting rule")

	OrgIDListEmpty                 = errors.New("org id list is empty")
	ProjectIDListEmpty             = errors.New("project id list is empty")
	OrgIDListProjectIDListNotMatch = errors.New("org id list and project id list is not match")
)
