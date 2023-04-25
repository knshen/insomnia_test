package model

type Project struct {
	ProjectID    int64  `json:"project_id"`
	OrgID        int64  `json:"org_id"`
	ProjectTitle string `json:"project_title"`
}
