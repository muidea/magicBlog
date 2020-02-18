package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"

	commonCommon "github.com/muidea/magicCommon/common"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/magicCommon/session"

	"github.com/muidea/magicBlog/config"
	"github.com/muidea/magicBlog/core/handler"

	cmsClient "github.com/muidea/magicCMS/client"
	casModel "github.com/muidea/magicCas/model"
	casRegistry "github.com/muidea/magicCas/toolkit/cas"
	privateRegistry "github.com/muidea/magicCas/toolkit/private"
	engine "github.com/muidea/magicEngine"
)

// Registry 路由信息
type Registry struct {
	commonHandler        handler.CommonHandler
	sessionRegistry      session.Registry
	privateRouteRegistry privateRegistry.RouteRegistry
	casRouteRegistry     casRegistry.RouteRegistry

	userService string
	casService  string
	fileService string
	cmsService  string
	cmsCatalog  int

	bashPath string
}

// NewRoute create route
func NewRoute(
	sessionRegistry session.Registry,
	commonHandler handler.CommonHandler,
) *Registry {
	casService := config.CasService()

	route := &Registry{
		sessionRegistry:      sessionRegistry,
		commonHandler:        commonHandler,
		privateRouteRegistry: privateRegistry.NewRouteRegistry(casService, sessionRegistry),
		casRouteRegistry:     casRegistry.NewRouteRegistry(casService, sessionRegistry),
		casService:           casService,
		userService:          config.UserService(),
		fileService:          config.FileService(),
		cmsService:           config.CMSService(),
		cmsCatalog:           config.CMSCatalog(),
		bashPath:             "static/default",
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

	cmsClient := cmsClient.NewClient(s.casService)
	defer cmsClient.Release()
	cmsClient.BindSession(sessionInfo)

	var err error
	defer func() {
		if err != nil {
			curSession.RemoveOption(commonCommon.AuthAccount)
			curSession.RemoveOption(commonCommon.SessionIdentity)
		}
	}()

	accountPtr, accountSession, accountErr := cmsClient.StatusAccount()
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
	case "/api/v1/account/create/":
		s.recordCreateAccount(res, req)
	case "/api/v1/account/login/":
		s.recordLoginAccount(res, req)
	case "/api/v1/account/logout/":
		s.recordLogoutAccount(res, req)
	case "/api/v1/account/status/":
		s.updateSessionAccount(res, req)
	case "/api/v1/account/change/password/":
		s.recordChangePassword(res, req)
	}
}

// Login account login
func (s *Registry) Login(res http.ResponseWriter, req *http.Request) {
	type loginParam struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	type loginResult struct {
		commonDef.Result
		Redirect string `json:"redirect"`
	}

	sessionInfo := &commonCommon.SessionInfo{}
	sessionInfo.Decode(req)

	curSession := s.sessionRegistry.GetSession(res, req)
	result := &loginResult{}
	for {
		param := &loginParam{}
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

		cmsClient := cmsClient.NewClient(s.casService)
		defer cmsClient.Release()
		cmsClient.BindSession(sessionInfo)

		accountPtr, sessionPtr, err := cmsClient.LoginAccount(param.Account, param.Password)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = err.Error()
			break
		}

		curSession.SetOption(commonCommon.AuthAccount, accountPtr)
		curSession.SetOption(commonCommon.SessionIdentity, sessionPtr)

		result.ErrorCode = commonDef.Success
		result.Redirect = "/"
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

// View static view
func (s *Registry) View(res http.ResponseWriter, req *http.Request) {
	type viewResult struct {
		IsAuthOK bool        `json:"isAuthOK"`
		Content  interface{} `json:"content"`
	}

	_, fileName := path.Split(req.URL.EscapedPath())
	if fileName == "" {
		fileName = "index.html"
	}

	view := &viewResult{}
	switch fileName {
	case "contact.html":
	case "about.html":
	case "post.html":
	case "edit.html":
	case "login.html":
	default:
		view.Content = s.filterPostList()
	}

	curSession := s.sessionRegistry.GetSession(res, req)
	_, authOk := curSession.GetOption(commonCommon.AuthAccount)
	if authOk {
		view.IsAuthOK = authOk
	}

	fullFilePath := path.Join(s.bashPath, fileName)
	t, err := template.ParseFiles(fullFilePath)
	if err != nil {
		log.Println(err)
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(res, view)
}

// RegisterRoute 注册路由
func (s *Registry) RegisterRoute(router engine.Router) {
	// blog view routes
	indexURL := "/view/index.html"
	indexRoute := engine.CreateProxyRoute("/", "GET", indexURL, true)
	router.AddRoute(indexRoute, s)

	viewRoute := engine.CreateRoute("/view/**", "GET", s.View)
	router.AddRoute(viewRoute, s)

	// blog api routes
	s.casRouteRegistry.AddHandler("/api/v1/blog/post/", "POST", s.PostBlog)

	// account login,logout,status,changepassword
	//---------------------------------------------------------------------------------------
	loginRoute := engine.CreateRoute("/api/v1/account/login/", "POST", s.Login)
	router.AddRoute(loginRoute, s)

	logoutURL := net.JoinURL(s.casService, "/access/logout/")
	logoutRoute := engine.CreateProxyRoute("/api/v1/account/logout/", "DELETE", logoutURL, true)
	router.AddRoute(logoutRoute, s)

	statusURL := net.JoinURL(s.casService, "/access/status/")
	statusRoute := engine.CreateProxyRoute("/api/v1/account/status/", "GET", statusURL, true)
	router.AddRoute(statusRoute, s)

	createAccountURL := net.JoinURL(s.casService, "/account/create/")
	createAccountRoute := engine.CreateProxyRoute("/api/v1/account/create/", "POST", createAccountURL, true)
	router.AddRoute(createAccountRoute, s)

	changePasswordURL := net.JoinURL(s.casService, "/account/password/")
	changePasswordRoute := engine.CreateProxyRoute("/api/v1/account/change/password/", "PUT", changePasswordURL, true)
	router.AddRoute(changePasswordRoute, s)

	// update account,query account,delete account
	//---------------------------------------------------------------------------------------
	s.privateRouteRegistry.AddHandler("/api/v1/account/query/all/", "GET", casModel.ReadPrivate, s.QueryAllAccount)
	s.privateRouteRegistry.AddHandler("/api/v1/account/query/:id", "GET", casModel.ReadPrivate, s.QueryAccount)
	s.privateRouteRegistry.AddHandler("/api/v1/account/delete/:id", "DELETE", casModel.DeletePrivate, s.DeleteAccount)
	s.privateRouteRegistry.AddHandler("/api/v1/account/update/:id", "PUT", casModel.WritePrivate, s.UpdateAccount)

	// private
	//---------------------------------------------------------------------------------------
	s.privateRouteRegistry.AddHandler("/api/v1/private/query/", "GET", casModel.ReadPrivate, s.QueryPrivateGroup)
	s.privateRouteRegistry.AddHandler("/api/v1/private/save/", "POST", casModel.WritePrivate, s.SavePrivateGroup)
	s.privateRouteRegistry.AddHandler("/api/v1/private/destory/", "GET", casModel.DeletePrivate, s.DestoryPrivateGroup)

	// upload file
	//---------------------------------------------------------------------------------------
	uploadFileURL := net.JoinURL(s.fileService, "/file/upload/")
	uploadFileRoute := engine.CreateProxyRoute("/api/v1/file/upload/", "POST", uploadFileURL, true)
	router.AddRoute(uploadFileRoute, s)

	viewFileURL := net.JoinURL(s.fileService, "/file/download/")
	viewFileRoute := engine.CreateProxyRoute("/api/v1/file/view/", "GET", viewFileURL, true)
	router.AddRoute(viewFileRoute, s)

	// add more route define

	s.privateRouteRegistry.RegisterRoute(router)
	s.casRouteRegistry.RegisterRoute(router)
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
