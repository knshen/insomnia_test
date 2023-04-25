package container

import (
	"code.sk.org/insomnia_test/client"
	"code.sk.org/insomnia_test/dal/kv"
	"code.sk.org/insomnia_test/service/checkers"
	"time"
)

// singletons, thread-safe
var (
	// client
	AuthClient         = &client.DefaultAuthClient{}
	AuthCacheClient    = client.NewAuthCacheClient(AuthClient, time.Minute, time.Hour*24)
	ProjectClient      = &client.DefaultProjectClient{}
	ProjectCacheClient = client.NewProjectCacheClient(ProjectClient, time.Minute*10, time.Hour*24)
	IDGen              = &client.DefaultIDGenerateClient{}

	// kv
	LintDal = &kv.RedisLintRuleDal{}

	// checkers
	LintContentChecker  = &checkers.LintYamlContentChecker{}
	LintIDsChecker      = &checkers.LintIDsChecker{}
	OrgAdminChecker     = checkers.NewOrgAdminChecker(AuthCacheClient)
	ProjectPermChecker  = checkers.NewProjectPermChecker(AuthCacheClient)
	ProjectExistChecker = checkers.NewProjectExistChecker(ProjectCacheClient)
)
