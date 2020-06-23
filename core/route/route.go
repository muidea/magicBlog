package route

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"
	"text/template"
	"time"

	commonCommon "github.com/muidea/magicCommon/common"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/session"
	"github.com/muidea/magicCommon/task"

	casToolkit "github.com/muidea/magicCas/toolkit"

	"github.com/muidea/magicBlog/config"
	"github.com/muidea/magicBlog/model"

	cmsCommon "github.com/muidea/magicCMS/common"
	cmsModel "github.com/muidea/magicCMS/model"

	engine "github.com/muidea/magicEngine"
)

// Registry 路由信息
type Registry struct {
	sessionRegistry  session.Registry
	casRouteRegistry casToolkit.CasRegistry

	cmsService string
	cmsCatalog int

	basePath       string
	currentCatalog *cmsModel.CatalogLite
	systemCatalog  *cmsModel.CatalogLite
	archiveCatalog *cmsModel.CatalogLite
	authorCatalog  *cmsModel.CatalogLite

	blogSetting map[string]string

	backgroundRoutine *task.BackgroundRoutine
}

// NewRoute create route
func NewRoute(
	sessionRegistry session.Registry,
) *Registry {

	backgroundRoutine := task.NewBackgroundRoutince()

	route := &Registry{
		sessionRegistry: sessionRegistry,
		cmsService:      config.CMSService(),
		cmsCatalog:      config.CMSCatalog(),
		basePath:        "static/default",
	}

	route.casRouteRegistry = toolkit.NewCasRegistry(route)
	route.backgroundRoutine = backgroundRoutine

	backgroundRoutine.Timer(&archiveBlogTask{registry: route}, 24*time.Hour, 2*time.Hour)

	return route
}

// Verify verify current session
func (s *Registry) Verify(res http.ResponseWriter, req *http.Request) (err error) {
	curSession := s.sessionRegistry.GetSession(res, req)

	cmsClient, cmsErr := s.getCMSClient(curSession)
	if cmsErr != nil {
		err = cmsErr
		log.Printf("getCMSClient failed, err:%s", err.Error())
		return
	}
	defer cmsClient.Release()

	sessionEntity, _, sessionErr := cmsClient.RefreshStatus()
	if sessionErr != nil {
		err = sessionErr
		log.Printf("verify current session failed, err:%s", err.Error())
		return
	}

	curSession.SetOption(commonCommon.AuthAccount, sessionEntity)

	return
}

// Handle middleware handler
func (s *Registry) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := curSession.GetSessionInfo()
	sessionInfo.Scope = commonCommon.ShareSession

	values := req.URL.Query()
	values = sessionInfo.Encode(values)
	req.URL.RawQuery = values.Encode()

	ctx.Next()
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

		cmsClient, cmsErr := s.getCMSClient(curSession)
		if cmsErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数,输入参数为空"
			break
		}
		defer cmsClient.Release()

		accountPtr, _, err := cmsClient.LoginAccount(param.Account, param.Password)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = err.Error()
			break
		}

		curSession.SetOption(commonCommon.AuthAccount, accountPtr)

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

// Logout account logout
func (s *Registry) Logout(res http.ResponseWriter, req *http.Request) {
	type logoutResult struct {
		commonDef.Result
		Redirect string `json:"redirect"`
	}

	curSession := s.sessionRegistry.GetSession(res, req)
	result := &logoutResult{}
	for {
		cmsClient, cmsErr := s.getCMSClient(curSession)
		if cmsErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数,输入参数为空"
			break
		}
		defer cmsClient.Release()

		_, err := cmsClient.LogoutAccount()
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = err.Error()
			break
		}

		curSession.RemoveOption(commonCommon.AuthAccount)

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
		IsAuthOK    bool                    `json:"isAuthOK"`
		ElapsedTime string                  `json:"elapsedTime"`
		CurrentURL  string                  `json:"currentUrl"`
		Setting     *model.Setting          `json:"setting"`
		Catalogs    []*cmsModel.CatalogLite `json:"catalogs"`
		Archives    []*cmsModel.CatalogLite `json:"archives"`
		Content     interface{}             `json:"content"`
	}

	preTime := time.Now()
	curSession := s.sessionRegistry.GetSession(res, req)
	_, authOk := curSession.GetOption(commonCommon.AuthAccount)

	var content interface{}
	var contentErr error
	view := &viewResult{IsAuthOK: authOk}
	fileName := ""
	filter := &filter{}
	for {
		cmsClnt, cmsErr := s.getCMSClient(curSession)
		if cmsErr != nil {
			log.Printf("getCMSClient failed, err:%s", cmsErr.Error())
			fileName = "500.html"
			break
		}
		defer cmsClnt.Release()

		catalogs, archives, articles, commonErr := s.queryBlogCommon(cmsClnt)
		if commonErr != nil {
			log.Printf("queryBlogCommon failed, err:%s", commonErr.Error())
			fileName = "500.html"
			break
		}

		settingPtr, settingErr := s.getBlogSetting(articles)
		if settingErr == nil {
			view.Setting = settingPtr
		}

		view.Catalogs = catalogs
		view.Archives = archives

		err := filter.decode(req)
		if err != nil {
			fileName = "404.html"
			break
		}

		if filter.isArchive() {
			fileName, content, contentErr = s.filterBlogArchive(filter, archives, cmsClnt)
			break
		}
		if filter.isCatalog() {
			fileName, content, contentErr = s.filterBlogCatalog(filter, catalogs, cmsClnt)
			break
		}

		if filter.isAuthor() {
			fileName, content, contentErr = s.filterBlogAuthor(filter, cmsClnt)
			break
		}

		fileName = filter.fileName
		if fileName == "" {
			fileName = "index.html"
		}

		switch fileName {
		case "index.html":
			content, contentErr = s.filterBlogPostList(filter, cmsClnt)
		case "edit.html":
			if !authOk {
				http.Redirect(res, req, "/", http.StatusMovedPermanently)
				return
			}

			action := filter.action
			if action == "update_post" {
				fileName, content, contentErr = s.queryBlogPostEdit(filter, cmsClnt)
			} else if action == "delete_post" {
				contentErr = s.deleteBlogPost(filter, cmsClnt)
				if contentErr == nil {
					http.Redirect(res, req, "/", http.StatusMovedPermanently)
					return
				}
			} else if action == "delete_catalog" {
				s.deleteBlogCatalog(filter, catalogs, cmsClnt)
				http.Redirect(res, req, "/", http.StatusMovedPermanently)
				return
			} else {
				content = map[string]interface{}{"ID": 0, "Title": "", "Content": "", "Catalog": ""}
			}
		case "setting.html":
			if !authOk {
				http.Redirect(res, req, "/", http.StatusMovedPermanently)
				return
			}
			fileName, content, contentErr = s.queryBlogSetting(filter, articles, cmsClnt)
		case "login.html":
			if authOk {
				http.Redirect(res, req, "/", http.StatusMovedPermanently)
				return
			}
			fileName, content, contentErr = s.queryBlogLogin(filter, articles, cmsClnt)
		case "about.html":
			fileName, content, contentErr = s.queryBlogAbout(filter, articles, cmsClnt)
		case "contact.html":
			fileName, content, contentErr = s.queryBlogContact(filter, articles, cmsClnt)
		default:
			fileName, content, contentErr = s.queryBlogPost(filter, cmsClnt)
		}

		break
	}

	if contentErr != nil {
		fileName = "500.html"
	} else {
		view.Content = content
	}

	fullFilePath := path.Join(s.basePath, fileName)
	t, err := template.ParseFiles(fullFilePath)
	if err != nil {
		log.Println(err)
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")

	if filter.pageFilter != nil {
		filter.pageFilter.PageNum++

		qv := url.Values{}
		qv = filter.pageFilter.Encode(qv)

		curURL := url.URL{}
		curURL.Path = req.URL.Path
		curURL.RawQuery = qv.Encode()
		view.CurrentURL = curURL.String()
	}

	elapsedTime := time.Now().Sub(preTime)
	view.ElapsedTime = elapsedTime.String()
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

	// setting api routes
	s.casRouteRegistry.AddHandler("/api/v1/blog/setting/", "POST", s.SettingBlog)

	// comment api routes
	commentRoute := engine.CreateRoute("/api/v1/comment/post/", "POST", s.PostComment)
	router.AddRoute(commentRoute, s)

	// reply comment api routes
	s.casRouteRegistry.AddHandler("/api/v1/comment/reply/", "POST", s.ReplyComment)

	// delete comment api routes
	s.casRouteRegistry.AddHandler("/api/v1/comment/delete/", "POST", s.DeleteComment)

	// account login,logout,status,changepassword
	//---------------------------------------------------------------------------------------
	loginRoute := engine.CreateRoute(cmsCommon.LoginAccountURL, "POST", s.Login)
	router.AddRoute(loginRoute, s)

	logoutRoute := engine.CreateRoute(cmsCommon.LogoutAccountURL, "DELETE", s.Logout)
	router.AddRoute(logoutRoute, s)

	s.casRouteRegistry.RegisterRoute(router)
}
