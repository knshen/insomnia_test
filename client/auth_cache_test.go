package client

import (
	"code.sk.org/insomnia_test/consts"
	"code.sk.org/insomnia_test/dal/kv"
	"code.sk.org/insomnia_test/mocks/mockclient"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"time"

	"testing"
)

type AuthCacheTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	ctx      context.Context

	mockAuthClient *mockclient.MockIAuthClient
}

func TestAuthCacheTestSuite(t *testing.T) {
	suite.Run(t, new(AuthCacheTestSuite))
}

func (s *AuthCacheTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.mockCtrl = gomock.NewController(s.T())
	kv.InitRedisClient()

	s.mockAuthClient = mockclient.NewMockIAuthClient(s.mockCtrl)
}

func (s *AuthCacheTestSuite) TearDownSuite() {
	s.mockCtrl.Finish()
}

func (s *AuthCacheTestSuite) TestGetUserRoleFromCache() {
	Convey("test get role from cache", s.T(), func() {
		authClient := NewAuthCacheClient(s.mockAuthClient, time.Second*10, time.Hour)
		token := "mock_token"
		var orgID int64 = 111
		s.mockAuthClient.EXPECT().GetUserRole(s.ctx, token, orgID).Return(consts.AdminRole, nil).Times(1)

		// get by rpc
		role, err := authClient.GetUserRole(s.ctx, token, orgID)
		So(role, ShouldEqual, consts.AdminRole)
		So(err, ShouldBeNil)

		// get from cache
		time.Sleep(time.Second) // wait till async goroutine finish
		role, err = authClient.GetUserRole(s.ctx, token, orgID)
		So(role, ShouldEqual, consts.AdminRole)
		So(err, ShouldBeNil)
	})
}

func (s *AuthCacheTestSuite) TestUseExpiredCache() {
	Convey("test use expire cache", s.T(), func() {
		authClient := NewAuthCacheClient(s.mockAuthClient, time.Second*1, time.Hour)
		token := "mock_token"
		var orgID int64 = 222

		s.mockAuthClient.EXPECT().GetUserRole(s.ctx, token, orgID).Return(consts.AdminRole, nil).Times(1)
		role, err := authClient.GetUserRole(s.ctx, token, orgID)
		So(role, ShouldEqual, consts.AdminRole)
		So(err, ShouldBeNil)

		// get from cache, use expired value
		time.Sleep(time.Second * 3) // cache expire but still use
		s.mockAuthClient.EXPECT().GetUserRole(s.ctx, token, orgID).Return(consts.NonAdminRole, nil).AnyTimes()
		role, err = authClient.GetUserRole(s.ctx, token, orgID)
		So(role, ShouldEqual, consts.AdminRole)
		So(err, ShouldBeNil)

		// this time cache has been updated
		time.Sleep(time.Second)
		role, err = authClient.GetUserRole(s.ctx, token, orgID)
		So(role, ShouldEqual, consts.NonAdminRole)
		So(err, ShouldBeNil)
	})
}

func (s *AuthCacheTestSuite) TestWhenAuthServiceDown() {
	Convey("test when auth service down", s.T(), func() {
		authClient := NewAuthCacheClient(s.mockAuthClient, time.Second*1, time.Hour)
		token := "mock_token"
		var orgID int64 = 333
		s.mockAuthClient.EXPECT().GetUserRole(s.ctx, token, orgID).Return(consts.AdminRole, nil).Times(1)

		// get by rpc
		role, err := authClient.GetUserRole(s.ctx, token, orgID)
		So(role, ShouldEqual, consts.AdminRole)
		So(err, ShouldBeNil)

		// now auth service is down, should always use expired value
		time.Sleep(time.Second * 2)
		s.mockAuthClient.EXPECT().GetUserRole(s.ctx, token, orgID).Return(consts.Unknown, errors.New("network err")).AnyTimes()
		role, err = authClient.GetUserRole(s.ctx, token, orgID)
		So(role, ShouldEqual, consts.AdminRole)
		So(err, ShouldBeNil)

		role, err = authClient.GetUserRole(s.ctx, token, orgID)
		So(role, ShouldEqual, consts.AdminRole)
		So(err, ShouldBeNil)
	})
}
