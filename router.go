package main

import (
	"code.sk.org/insomnia_test/consts"
	"code.sk.org/insomnia_test/container"
	"code.sk.org/insomnia_test/model"
	"code.sk.org/insomnia_test/mws"
	"code.sk.org/insomnia_test/service"
	"code.sk.org/insomnia_test/service/checkers"
	"context"
	"net/http"
	"strings"
	"time"
)

func Ping(w http.ResponseWriter, req *http.Request) {
	RequestHandleWrapper(w, req, func(ctx context.Context, req *http.Request, writer http.ResponseWriter) (interface{}, error) {
		return time.Now().Format("2006-01-02 15:04:05"), nil
	})
}

func HandleLintRule(w http.ResponseWriter, req *http.Request) {
	RequestHandleWrapper(w, req, func(ctx context.Context, req *http.Request, writer http.ResponseWriter) (interface{}, error) {
		rule, err := model.GetLintRuleFromReq(req)
		if err != nil {
			return nil, consts.RequestBodyIllegal
		}

		var svc *service.LintService

		switch strings.ToLower(req.Method) {
		case "get":
			svc = service.NewLintService(container.LintDal, checkers.BuildLintCheckers(
				container.LintIDsChecker,
				container.ProjectExistChecker,
				container.ProjectPermChecker), container.IDGen)
			return svc.GetLintRule(ctx, rule)
		case "put":
			svc = service.NewLintService(container.LintDal, checkers.BuildLintCheckers(
				container.LintIDsChecker,
				container.ProjectExistChecker,
				container.LintContentChecker,
				container.OrgAdminChecker), container.IDGen)
			if err = svc.CreateLintRule(ctx, rule); err != nil {
				return "", err
			}
			return "ok", nil
		case "post":
			svc = service.NewLintService(container.LintDal, checkers.BuildLintCheckers(
				container.LintIDsChecker,
				container.ProjectExistChecker,
				container.LintContentChecker,
				container.OrgAdminChecker), container.IDGen)
			if err = svc.UpdateLintRule(ctx, rule); err != nil {
				return "", err
			}
			return "ok", nil
		default:
			return "", consts.HttpMethodNotAllow
		}
	})
}

// common request handler with mws
func RequestHandleWrapper(w http.ResponseWriter, req *http.Request, bizHandler mws.RequestHandler) {
	chain := mws.Build(mws.AccessLogMiddleware(), mws.UserTokenCheckerMiddleware(), mws.RespHandlerMiddleware())(bizHandler)
	ctx := context.Background()
	chain(ctx, req, w)
}
