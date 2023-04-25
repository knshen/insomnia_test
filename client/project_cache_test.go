package client

import (
	"code.sk.org/insomnia_test/dal/kv"
	"code.sk.org/insomnia_test/mocks/mockclient"
	"code.sk.org/insomnia_test/model"
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ProjectCacheTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	ctx      context.Context

	mockProClient *mockclient.MockIProjectClient
}

func TestProjectCacheTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectCacheTestSuite))
}

func (s *ProjectCacheTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.mockCtrl = gomock.NewController(s.T())
	kv.InitRedisClient()

	s.mockProClient = mockclient.NewMockIProjectClient(s.mockCtrl)
}

func (s *ProjectCacheTestSuite) TearDownSuite() {
	s.mockCtrl.Finish()
}

func (s *ProjectCacheTestSuite) TestGetProjectFromCache() {
	Convey("test get project from cache", s.T(), func() {
		projectClient := NewProjectCacheClient(s.mockProClient, time.Second*10, time.Hour)
		var pid int64 = 111
		project := &model.Project{
			ProjectID:    111,
			OrgID:        11,
			ProjectTitle: "mock_111",
		}
		s.mockProClient.EXPECT().GetByProjectID(s.ctx, pid).Return(project, nil).Times(1)

		// get by rpc
		res, err := projectClient.GetByProjectID(s.ctx, pid)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(res.ProjectTitle, ShouldEqual, "mock_111")
		So(res.OrgID, ShouldEqual, int64(11))

		// get from cache
		time.Sleep(time.Second) // wait till async goroutine finish
		res, err = projectClient.GetByProjectID(s.ctx, pid)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(res.ProjectTitle, ShouldEqual, "mock_111")
		So(res.OrgID, ShouldEqual, int64(11))
	})
}

func (s *ProjectCacheTestSuite) TestProjectNotExist() {
	Convey("test project not exist", s.T(), func() {
		projectClient := NewProjectCacheClient(s.mockProClient, time.Second*10, time.Hour)
		var pid int64 = 222

		s.mockProClient.EXPECT().GetByProjectID(s.ctx, pid).Return(nil, nil).MaxTimes(2)

		// get by rpc, project service return nil
		res, err := projectClient.GetByProjectID(s.ctx, pid)
		So(err, ShouldBeNil)
		So(res, ShouldBeNil)

		// get from cache
		time.Sleep(time.Second) // wait till async goroutine finish
		// now value of cache should be null and should not call rpc
		res, err = projectClient.GetByProjectID(s.ctx, pid)
		So(err, ShouldBeNil)
		So(res, ShouldBeNil)
	})
}
