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

	casClient "github.com/muidea/magicCas/client"
	casModel "github.com/muidea/magicCas/model"
	casRegistry "github.com/muidea/magicCas/toolkit/cas"
	engine "github.com/muidea/magicEngine"
)

// Registry 路由信息
type Registry struct {
	commonHandler    handler.CommonHandler
	sessionRegistry  session.Registry
	casRouteRegistry casRegistry.RouteRegistry

	casService string
	cmsService string
	cmsCatalog int

	bashPath string
}

// NewRoute create route
func NewRoute(
	sessionRegistry session.Registry,
	commonHandler handler.CommonHandler,
) *Registry {
	casService := config.CasService()

	route := &Registry{
		sessionRegistry:  sessionRegistry,
		commonHandler:    commonHandler,
		casRouteRegistry: casRegistry.NewRouteRegistry(casService, sessionRegistry),
		casService:       casService,
		cmsService:       config.CMSService(),
		cmsCatalog:       config.CMSCatalog(),
		bashPath:         "static/default",
	}

	return route
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
	case "/api/v1/account/login/":
		s.recordLoginAccount(res, req)
	case "/api/v1/account/logout/":
		s.recordLogoutAccount(res, req)
	case "/api/v1/account/status/":
		s.updateSessionAccount(res, req)
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
