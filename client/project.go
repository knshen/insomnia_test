package client

import (
	"code.sk.org/insomnia_test/model"
	"context"
	"time"
)

type IProjectClient interface {
	GetByProjectID(ctx context.Context, projectID int64) (*model.Project, error)
}

type DefaultProjectClient struct {
}

func (c *DefaultProjectClient) GetByProjectID(ctx context.Context, projectID int64) (*model.Project, error) {
	// call rpc of project service
	time.Sleep(time.Millisecond * 200)

	data := map[int64]*model.Project{
		1001: &model.Project{
			ProjectID:    1001,
			OrgID:        100,
			ProjectTitle: "mock_title",
		},
		1002: &model.Project{
			ProjectID:    1002,
			OrgID:        100,
			ProjectTitle: "mock_title",
		},
		1003: &model.Project{
			ProjectID:    1003,
			OrgID:        100,
			ProjectTitle: "mock_title",
		},
	}

	return data[projectID], nil
}
