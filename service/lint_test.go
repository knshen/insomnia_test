package service

import (
	"code.sk.org/insomnia_test/client"
	"code.sk.org/insomnia_test/consts"
	"code.sk.org/insomnia_test/container"
	"code.sk.org/insomnia_test/dal/kv"
	"code.sk.org/insomnia_test/mocks/mockclient"
	"code.sk.org/insomnia_test/model"
	"code.sk.org/insomnia_test/service/checkers"
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LintServiceTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	ctx      context.Context

	mockAuthClient  *mockclient.MockIAuthClient
	mockProjectAuth *mockclient.MockIProjectClient
	idGen           *client.DefaultIDGenerateClient
}

func TestLintServiceTestSuite(t *testing.T) {
	suite.Run(t, new(LintServiceTestSuite))
}

func (s *LintServiceTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.mockCtrl = gomock.NewController(s.T())
	kv.InitRedisClient()

	s.idGen = &client.DefaultIDGenerateClient{}
}

func (s *LintServiceTestSuite) TearDownSuite() {
	s.mockCtrl.Finish()
}

func (s *LintServiceTestSuite) TestCreateLintRule() {
	adminUserToken := "admin_user"
	commonUserToken := "non_admin_user"
	s.mockAuthClient = mockclient.NewMockIAuthClient(s.mockCtrl)
	s.mockProjectAuth = mockclient.NewMockIProjectClient(s.mockCtrl)

	Convey("test create and get lint rule", s.T(), func() {
		createSvc := NewLintService(container.LintDal, checkers.BuildLintCheckers(
			container.LintIDsChecker,
			checkers.NewProjectExistChecker(s.mockProjectAuth),
			container.LintContentChecker,
			checkers.NewOrgAdminChecker(s.mockAuthClient)), container.IDGen)

		querySvc := NewLintService(container.LintDal, checkers.BuildLintCheckers(
			container.LintIDsChecker,
			checkers.NewProjectExistChecker(s.mockProjectAuth),
			checkers.NewProjectPermChecker(s.mockAuthClient)), container.IDGen)

		projectData := &model.Project{
			ProjectID:    1001,
			OrgID:        100,
			ProjectTitle: "mock_title",
		}
		s.mockAuthClient.EXPECT().GetUserRole(s.ctx, adminUserToken, gomock.Any()).Return(consts.AdminRole, nil).AnyTimes()
		s.mockAuthClient.EXPECT().GetUserRole(s.ctx, commonUserToken, gomock.Any()).Return(consts.NonAdminRole, nil).AnyTimes()
		s.mockAuthClient.EXPECT().HasProjectPermission(s.ctx, gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()
		s.mockProjectAuth.EXPECT().GetByProjectID(s.ctx, gomock.Any()).Return(projectData, nil).AnyTimes()
		// test create success
		err := createSvc.CreateLintRule(s.ctx, &model.LintRuleRequest{
			LintRule: &model.LintRule{
				RuleID:        0,
				ProjectID:     1001,
				RawRuleString: "name: xxxx",
			},
			UserToken: adminUserToken,
		})
		So(err, ShouldBeNil)

		// query
		rule, err := querySvc.GetLintRule(s.ctx, &model.LintRuleRequest{
			LintRule: &model.LintRule{
				ProjectID: 1001,
			},
			UserToken: adminUserToken,
		})
		So(err, ShouldBeNil)
		So(rule, ShouldNotBeNil)
		So(rule.RuleID, ShouldBeGreaterThan, 0)
		So(rule.ProjectID, ShouldEqual, int64(1001))
		So(rule.OrgID, ShouldEqual, int64(100))
		So(rule.RawRuleString, ShouldEqual, "name: xxxx")

		// test can only create one rule for each project
		err = createSvc.CreateLintRule(s.ctx, &model.LintRuleRequest{
			LintRule: &model.LintRule{
				RuleID:        0,
				ProjectID:     1001,
				RawRuleString: "name: yyyy",
			},
			UserToken: adminUserToken,
		})
		So(err, ShouldEqual, consts.LintingRuleAlreadyExist)

		// test non-admin user create rule
		err = createSvc.CreateLintRule(s.ctx, &model.LintRuleRequest{
			LintRule: &model.LintRule{
				RuleID:        0,
				ProjectID:     1002,
				RawRuleString: "name: xxxx",
			},
			UserToken: commonUserToken,
		})
		So(err, ShouldEqual, consts.NoOrgAdmin)
	})
}

func (s *LintServiceTestSuite) TestUpdateLintRule() {
	adminUserToken := "admin_user"
	commonUserToken := "non_admin_user"
	s.mockAuthClient = mockclient.NewMockIAuthClient(s.mockCtrl)
	s.mockProjectAuth = mockclient.NewMockIProjectClient(s.mockCtrl)

	Convey("test update and get lint rule", s.T(), func() {
		updateSvc := NewLintService(container.LintDal, checkers.BuildLintCheckers(
			container.LintIDsChecker,
			checkers.NewProjectExistChecker(s.mockProjectAuth),
			container.LintContentChecker,
			checkers.NewOrgAdminChecker(s.mockAuthClient)), container.IDGen)

		querySvc := NewLintService(container.LintDal, checkers.BuildLintCheckers(
			container.LintIDsChecker,
			checkers.NewProjectExistChecker(s.mockProjectAuth),
			checkers.NewProjectPermChecker(s.mockAuthClient)), container.IDGen)

		projectData := &model.Project{
			ProjectID:    1001,
			OrgID:        222,
			ProjectTitle: "mock_title",
		}
		s.mockAuthClient.EXPECT().GetUserRole(s.ctx, adminUserToken, gomock.Any()).Return(consts.AdminRole, nil).AnyTimes()
		s.mockAuthClient.EXPECT().GetUserRole(s.ctx, commonUserToken, gomock.Any()).Return(consts.NonAdminRole, nil).AnyTimes()
		s.mockAuthClient.EXPECT().HasProjectPermission(s.ctx, gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()
		s.mockProjectAuth.EXPECT().GetByProjectID(s.ctx, gomock.Any()).Return(projectData, nil).AnyTimes()

		// test create success
		err := updateSvc.CreateLintRule(s.ctx, &model.LintRuleRequest{
			LintRule: &model.LintRule{
				ProjectID:     1001,
				RawRuleString: "name: xxxx",
			},
			UserToken: adminUserToken,
		})
		So(err, ShouldBeNil)

		// test update
		err = updateSvc.UpdateLintRule(s.ctx, &model.LintRuleRequest{
			LintRule: &model.LintRule{
				ProjectID:     1001,
				RawRuleString: "name: yyyy",
			},
			UserToken: adminUserToken,
		})
		So(err, ShouldBeNil)

		// query the updated rule
		rule, err := querySvc.GetLintRule(s.ctx, &model.LintRuleRequest{
			LintRule: &model.LintRule{
				OrgID:     222,
				ProjectID: 1001,
			},
			UserToken: adminUserToken,
		})
		So(err, ShouldBeNil)
		So(rule, ShouldNotBeNil)
		So(rule.RuleID, ShouldBeGreaterThan, 0)
		So(rule.ProjectID, ShouldEqual, int64(1001))
		So(rule.OrgID, ShouldEqual, int64(222))
		So(rule.RawRuleString, ShouldEqual, "name: yyyy")

	})
}

func (s *LintServiceTestSuite) TestGetLintRule() {
	Convey("test get lint rule without project permission", s.T(), func() {
		s.mockAuthClient = mockclient.NewMockIAuthClient(s.mockCtrl)
		s.mockProjectAuth = mockclient.NewMockIProjectClient(s.mockCtrl)
		querySvc := NewLintService(container.LintDal, checkers.BuildLintCheckers(
			container.LintIDsChecker,
			checkers.NewProjectExistChecker(s.mockProjectAuth),
			checkers.NewProjectPermChecker(s.mockAuthClient)), container.IDGen)

		projectData := &model.Project{
			ProjectID:    1001,
			OrgID:        333,
			ProjectTitle: "mock_title",
		}
		s.mockAuthClient.EXPECT().HasProjectPermission(s.ctx, gomock.Any(), gomock.Any()).Return(false, nil).AnyTimes()
		s.mockProjectAuth.EXPECT().GetByProjectID(s.ctx, gomock.Any()).Return(projectData, nil).AnyTimes()

		rule, err := querySvc.GetLintRule(s.ctx, &model.LintRuleRequest{
			LintRule: &model.LintRule{
				ProjectID: 1001,
			},
			UserToken: "some_token",
		})
		So(rule, ShouldBeNil)
		So(err, ShouldEqual, consts.NoProjectPerm)
	})

	Convey("test get lint rule project not exist", s.T(), func() {
		s.mockAuthClient = mockclient.NewMockIAuthClient(s.mockCtrl)
		s.mockProjectAuth = mockclient.NewMockIProjectClient(s.mockCtrl)
		querySvc := NewLintService(container.LintDal, checkers.BuildLintCheckers(
			container.LintIDsChecker,
			checkers.NewProjectExistChecker(s.mockProjectAuth),
			checkers.NewProjectPermChecker(s.mockAuthClient)), container.IDGen)

		s.mockAuthClient.EXPECT().HasProjectPermission(s.ctx, gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()
		s.mockProjectAuth.EXPECT().GetByProjectID(s.ctx, gomock.Any()).Return(nil, nil).AnyTimes()

		rule, err := querySvc.GetLintRule(s.ctx, &model.LintRuleRequest{
			LintRule: &model.LintRule{
				ProjectID: 1001,
			},
			UserToken: "some_token",
		})
		So(rule, ShouldBeNil)
		So(err, ShouldEqual, consts.ProjectNotExist)
	})

	Convey("test project does not have lint rule", s.T(), func() {
		s.mockAuthClient = mockclient.NewMockIAuthClient(s.mockCtrl)
		s.mockProjectAuth = mockclient.NewMockIProjectClient(s.mockCtrl)
		querySvc := NewLintService(container.LintDal, checkers.BuildLintCheckers(
			container.LintIDsChecker,
			checkers.NewProjectExistChecker(s.mockProjectAuth),
			checkers.NewProjectPermChecker(s.mockAuthClient)), container.IDGen)

		project := &model.Project{
			ProjectID:    1007,
			OrgID:        100,
			ProjectTitle: "mock_title",
		}
		s.mockAuthClient.EXPECT().HasProjectPermission(s.ctx, gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()
		s.mockProjectAuth.EXPECT().GetByProjectID(s.ctx, gomock.Any()).Return(project, nil).AnyTimes()

		rule, err := querySvc.GetLintRule(s.ctx, &model.LintRuleRequest{
			LintRule: &model.LintRule{
				ProjectID: 1007,
			},
			UserToken: "some_token",
		})
		So(rule, ShouldBeNil)
		So(err, ShouldEqual, consts.NoLintingRuleBind)
	})
}
