package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	commonCommon "github.com/muidea/magicCommon/common"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/magicCommon/session"

	"github.com/muidea/magicBlog/config"
	"github.com/muidea/magicBlog/core/handler"

	casClient "github.com/muidea/magicCas/client"
	casCommon "github.com/muidea/magicCas/common"
	casModel "github.com/muidea/magicCas/model"
	casPrivate "github.com/muidea/magicCas/toolkit/private"
	engine "github.com/muidea/magicEngine"
	userClient "github.com/muidea/magicUser/client"
	userModel "github.com/muidea/magicUser/model"
)

// Registry 路由信息
type Registry struct {
	commonHandler        handler.CommonHandler
	sessionRegistry      session.Registry
	privateRouteRegistry casPrivate.RouteRegistry

	userService string
	casService  string
	fileService string
}

// NewRoute create route
func NewRoute(
	sessionRegistry session.Registry,
	commonHandler handler.CommonHandler,
) *Registry {
	route := &Registry{
		sessionRegistry:      sessionRegistry,
		commonHandler:        commonHandler,
		privateRouteRegistry: casPrivate.NewRouteRegistry(sessionRegistry),
		casService:           config.CasService(),
		userService:          config.UserService(),
		fileService:          config.FileService(),
	}

	return route
}

func (s *Registry) recordCreateAccount(res http.ResponseWriter, req *http.Request) {
	account := &casModel.AccountView{}
	err := net.ParseJSONBody(req, account)
	if err != nil {
		return
	}

	memo := fmt.Sprintf("新建账号%s", account.Account)
	s.writelog(res, req, memo)
}

func (s *Registry) recordLoginAccount(res http.ResponseWriter, req *http.Request) {
	curSession := s.sessionRegistry.GetSession(res, req)
	authPtr, authOK := curSession.GetOption(commonCommon.AuthAccount)
	if authOK {
		acountPtr := authPtr.(*casModel.AccountView)
		memo := fmt.Sprintf("账号%s登录", acountPtr.Account)
		s.writelog(res, req, memo)
	}
}

func (s *Registry) recordLogoutAccount(res http.ResponseWriter, req *http.Request) {
	curSession := s.sessionRegistry.GetSession(res, req)
	authPtr, authOK := curSession.GetOption(commonCommon.AuthAccount)
	if authOK {
		acountPtr := authPtr.(*casModel.AccountView)
		memo := fmt.Sprintf("账号%s登出", acountPtr.Account)
		s.writelog(res, req, memo)
	}
}

func (s *Registry) updateSessionAccount(res http.ResponseWriter, req *http.Request) {
	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := &commonCommon.SessionInfo{}
	sessionInfo.Decode(req)

	casClient := casClient.NewClient(s.casService)
	defer casClient.Release()
	casClient.BindSession(sessionInfo)

	var err error
	defer func() {
		if err != nil {
			curSession.RemoveOption(commonCommon.AuthAccount)
			curSession.RemoveOption(commonCommon.SessionIdentity)
		}
	}()

	accountPtr, accountSession, accountErr := casClient.StatusAccount()
	if accountErr != nil {
		err = accountErr
		log.Printf("get account status failed, err:%s", accountErr.Error())
		return
	}

	curSession.SetOption(commonCommon.AuthAccount, accountPtr)
	curSession.SetOption(commonCommon.SessionIdentity, accountSession)
}

func (s *Registry) recordChangePassword(res http.ResponseWriter, req *http.Request) {
	curSession := s.sessionRegistry.GetSession(res, req)
	authPtr, authOK := curSession.GetOption(commonCommon.AuthAccount)
	if authOK {
		acountPtr := authPtr.(*casModel.AccountView)
		memo := fmt.Sprintf("账号%s修改密码", acountPtr.Account)
		s.writelog(res, req, memo)
	}
}

// Handle middleware handler
func (s *Registry) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := &commonCommon.SessionInfo{}
	sessionInfo.Decode(req)
	sessionInfo.ID = curSession.ID()
	sessionInfo.Scope = commonCommon.ShareSession

	values := req.URL.Query()
	sessionInfo.Merge(values)
	req.URL.RawQuery = values.Encode()

	ctx.Next()

	switch req.URL.Path {
	case "/account/create/":
		s.recordCreateAccount(res, req)
	case "/account/login/":
		s.recordLoginAccount(res, req)
	case "/account/logout/":
		s.recordLogoutAccount(res, req)
	case "/account/status/":
		s.updateSessionAccount(res, req)
	case "/account/change/password/":
		s.recordChangePassword(res, req)
	}
}

// Login account login
func (s *Registry) Login(res http.ResponseWriter, req *http.Request) {

	sessionInfo := &commonCommon.SessionInfo{}
	sessionInfo.Decode(req)

	curSession := s.sessionRegistry.GetSession(res, req)
	result := &casCommon.LoginResult{}
	for {
		param := &casCommon.LoginParam{}
		err := net.ParseJSONBody(req, param)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		if param.Account == "" || param.Password == "" {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数,输入参数为空"
			break
		}

		casClient := casClient.NewClient(s.casService)
		defer casClient.Release()
		casClient.BindSession(sessionInfo)

		accountPtr, sessionPtr, err := casClient.LoginAccount(param.Account, param.Password)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = err.Error()
			break
		}

		curSession.SetOption(commonCommon.AuthAccount, accountPtr)
		curSession.SetOption(commonCommon.SessionIdentity, sessionPtr)

		result.Account = accountPtr
		result.SessionInfo = sessionPtr
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

// RegisterRoute 注册路由
func (s *Registry) RegisterRoute(router engine.Router) {
	indexURL := "/static/default/"
	indexRoute := engine.CreateProxyRoute("/", "GET", indexURL, true)
	router.AddRoute(indexRoute, s)

	// account login,logout,status,changepassword
	//---------------------------------------------------------------------------------------
	loginRoute := engine.CreateRoute("/account/login/", "POST", s.Login)
	router.AddRoute(loginRoute, s)

	logoutURL := net.JoinURL(s.casService, "/access/logout/")
	logoutRoute := engine.CreateProxyRoute("/account/logout/", "DELETE", logoutURL, true)
	router.AddRoute(logoutRoute, s)

	statusURL := net.JoinURL(s.casService, "/access/status/")
	statusRoute := engine.CreateProxyRoute("/account/status/", "GET", statusURL, true)
	router.AddRoute(statusRoute, s)

	createAccountURL := net.JoinURL(s.casService, "/account/create/")
	createAccountRoute := engine.CreateProxyRoute("/account/create/", "POST", createAccountURL, true)
	router.AddRoute(createAccountRoute, s)

	changePasswordURL := net.JoinURL(s.casService, "/account/password/")
	changePasswordRoute := engine.CreateProxyRoute("/account/change/password/", "PUT", changePasswordURL, true)
	router.AddRoute(changePasswordRoute, s)

	// update account,query account,delete account
	//---------------------------------------------------------------------------------------
	s.privateRouteRegistry.AddHandler("/account/query/all/", "GET", casModel.ReadPrivate, s.QueryAllAccount)
	s.privateRouteRegistry.AddHandler("/account/query/:id", "GET", casModel.ReadPrivate, s.QueryAccount)
	s.privateRouteRegistry.AddHandler("/account/delete/:id", "DELETE", casModel.DeletePrivate, s.DeleteAccount)
	s.privateRouteRegistry.AddHandler("/account/update/:id", "PUT", casModel.WritePrivate, s.UpdateAccount)

	// private
	//---------------------------------------------------------------------------------------
	s.privateRouteRegistry.AddHandler("/private/query/", "GET", casModel.ReadPrivate, s.QueryPrivateGroup)
	s.privateRouteRegistry.AddHandler("/private/save/", "POST", casModel.WritePrivate, s.SavePrivateGroup)
	s.privateRouteRegistry.AddHandler("/private/destory/", "GET", casModel.DeletePrivate, s.DestoryPrivateGroup)

	// upload file
	//---------------------------------------------------------------------------------------
	uploadFileURL := net.JoinURL(s.fileService, "/file/upload/")
	uploadFileRoute := engine.CreateProxyRoute("/file/upload/", "POST", uploadFileURL, true)
	router.AddRoute(uploadFileRoute, s)

	viewFileURL := net.JoinURL(s.fileService, "/file/download/")
	viewFileRoute := engine.CreateProxyRoute("/file/view/", "GET", viewFileURL, true)
	router.AddRoute(viewFileRoute, s)

	// add more route define

	s.privateRouteRegistry.RegisterRoute(router)
}

func (s *Registry) writelog(res http.ResponseWriter, req *http.Request, memo string) {
	address := net.GetHTTPRemoteAddress(req)
	account := ""
	curSession := s.sessionRegistry.GetSession(res, req)
	authVal, ok := curSession.GetOption(commonCommon.AuthAccount)
	if ok {
		for {
			accountView, accountOK := authVal.(*casModel.AccountView)
			if accountOK {
				account = accountView.Account
				break
			}

			endpointView, endpointOK := authVal.(*casModel.EndpointView)
			if endpointOK {
				account = endpointView.Endpoint
				break
			}

			break
		}
	}

	_, logErr := s.commonHandler.WriteOpLog(account, address, memo)
	if logErr != nil {
		log.Printf("WriteOpLog failed, err:%s", logErr.Error())
	}
}

func (s *Registry) queryUser(sessionInfo *commonCommon.SessionInfo, id int) (ret *userModel.User, err error) {
	userClient := userClient.NewClient(s.userService)
	defer userClient.Release()
	userClient.BindSession(sessionInfo)

	userPtr, userErr := userClient.QueryUserByID(id)
	if userErr != nil {
		err = userErr
		log.Printf("query user failed, err:%s", userErr.Error())
		return
	}

	ret = userPtr
	return
}
