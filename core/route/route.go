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
	cmsCommon "github.com/muidea/magicCMS/common"
	cmsModel "github.com/muidea/magicCMS/model"

	casModel "github.com/muidea/magicCas/model"
	casRoute "github.com/muidea/magicCas/toolkit/route"
	engine "github.com/muidea/magicEngine"
)

// Registry 路由信息
type Registry struct {
	commonHandler    handler.CommonHandler
	sessionRegistry  session.Registry
	casRouteRegistry casRoute.CasRegistry

	casService string
	cmsService string
	cmsCatalog int

	basePath       string
	currentCatalog *cmsModel.CatalogTree
	archiveCatalog *cmsModel.CatalogTree
}

// NewRoute create route
func NewRoute(
	sessionRegistry session.Registry,
	commonHandler handler.CommonHandler,
) *Registry {
	route := &Registry{
		sessionRegistry: sessionRegistry,
		commonHandler:   commonHandler,
		casService:      config.CasService(),
		cmsService:      config.CMSService(),
		cmsCatalog:      config.CMSCatalog(),
		basePath:        "static/default",
	}

	route.casRouteRegistry = casRoute.NewCasRegistry(route)

	return route
}

// Verify verify current session
func (s *Registry) Verify(res http.ResponseWriter, req *http.Request) (err error) {
	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := curSession.GetSessionInfo()

	cmsClient := cmsClient.NewClient(s.cmsService)
	defer cmsClient.Release()
	cmsClient.BindSession(sessionInfo)

	sessionInfo, sessionErr := cmsClient.VerifySession()
	if sessionErr != nil {
		err = sessionErr
		log.Printf("verify current session failed, err:%s", sessionErr.Error())
		return
	}

	curSession.SetSessionInfo(sessionInfo)

	return
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

func (s *Registry) recordPostBlog(res http.ResponseWriter, req *http.Request, title string) {
	curSession := s.sessionRegistry.GetSession(res, req)
	authPtr, authOK := curSession.GetOption(commonCommon.AuthAccount)
	if authOK {
		acountPtr := authPtr.(*casModel.AccountView)
		memo := fmt.Sprintf("%s发布%s", acountPtr.Account, title)
		s.writelog(res, req, memo)
	}
}

// Handle middleware handler
func (s *Registry) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := curSession.GetSessionInfo()
	sessionInfo.Scope = commonCommon.ShareSession

	values := req.URL.Query()
	sessionInfo.Merge(values)
	req.URL.RawQuery = values.Encode()

	ctx.Next()

	switch req.URL.Path {
	case cmsCommon.LoginAccountURL:
		s.recordLoginAccount(res, req)
	case cmsCommon.LogoutAccountURL:
		s.recordLogoutAccount(res, req)
	case cmsCommon.StatusAccountURL:
		s.Verify(res, req)
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

	curSession := s.sessionRegistry.GetSession(res, req)
	result := &loginResult{}
	for {
		sessionInfo := curSession.GetSessionInfo()
		sessionInfo.Scope = commonCommon.ShareSession

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

		cmsClient := cmsClient.NewClient(s.cmsService)
		defer cmsClient.Release()
		cmsClient.BindSession(sessionInfo)

		accountPtr, sessionPtr, err := cmsClient.LoginAccount(param.Account, param.Password)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = err.Error()
			break
		}

		curSession.SetOption(commonCommon.AuthAccount, accountPtr)
		curSession.SetSessionInfo(sessionPtr)

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

	curSession := s.sessionRegistry.GetSession(res, req)
	_, authOk := curSession.GetOption(commonCommon.AuthAccount)

	view := &viewResult{IsAuthOK: authOk}
	fileName := ""
	for {
		filter := &filter{}
		err := filter.Decode(req)
		if err != nil {
			fileName = "404.html"
			break
		}

		switch fileName {
		case "about.html":
			view.Content = s.filterAbout(res, req)
		case "contact.html":
			view.Content = s.filterContact(res, req)
		case "index.html":
			view.Content = s.filterPostList(res, req)
		case "post.html":
		case "edit.html":
			if !authOk {
				http.Redirect(res, req, "/", http.StatusMovedPermanently)
				return
			}
		case "login.html":
			if authOk {
				http.Redirect(res, req, "/", http.StatusMovedPermanently)
				return
			}
		default:
			view.Content = s.filterPostList(res, req)
		}

		break
	}

	fullFilePath := path.Join(s.basePath, fileName)
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

	logoutURL := net.JoinURL(s.cmsService, cmsCommon.LoginAccountURL)
	logoutRoute := engine.CreateProxyRoute("/api/v1/account/logout/", "DELETE", logoutURL, true)
	router.AddRoute(logoutRoute, s)

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

			break
		}
	}

	_, logErr := s.commonHandler.WriteOpLog(account, address, memo)
	if logErr != nil {
		log.Printf("WriteOpLog failed, err:%s", logErr.Error())
	}
}
