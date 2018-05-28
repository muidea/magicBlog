package core

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"muidea.com/magicCenter/application/common"
	common_result "muidea.com/magicCommon/common"
	"muidea.com/magicCommon/foundation/net"
	"muidea.com/magicCommon/model"
	engine "muidea.com/magicEngine"
)

type route struct {
	pattern string
	method  string
	handler interface{}
}

func (s *route) Pattern() string {
	return s.pattern
}

func (s *route) Method() string {
	return s.method
}

func (s *route) Handler() interface{} {
	return s.handler
}

func newRoute(pattern, method string, handler interface{}) engine.Route {
	return &route{pattern: pattern, method: method, handler: handler}
}

// NewBlog 新建Blog
func NewBlog(centerServer, name, endpointID, authToken string) (Blog, bool) {
	blog := Blog{centerAgent: NewCenterAgent()}

	agent := NewCenterAgent()
	if !agent.Start(centerServer, endpointID, authToken) {
		return blog, false
	}
	blogCatalog, ok := agent.FetchCatalog(name)
	if !ok {
		ok = agent.CreateCatalog(name, "MagicBlog auto create catalog.")
		if !ok {
			log.Print("create blog root catalog failed.")
			return blog, false
		}
	}
	blogCatalog, ok = agent.FetchCatalog(name)
	if !ok {
		log.Print("fetch blog root ctalog failed.")
		return blog, false
	}

	blogContent := agent.QuerySummary(blogCatalog.ID)

	blog.centerAgent = agent
	blog.blogInfo = blogCatalog
	blog.blogContent = blogContent

	return blog, true
}

// Blog Blog对象
type Blog struct {
	centerAgent Agent
	blogInfo    model.CatalogDetailView
	blogContent []model.SummaryView
}

// Startup 启动
func (s *Blog) Startup(router engine.Router) {
	mainRoute := newRoute("/", "GET", s.mainPage)
	router.AddRoute(mainRoute)

	catalogListRoute := newRoute("/catalog/", "GET", s.catalogListPage)
	router.AddRoute(catalogListRoute)

	catalogRoute := newRoute("/catalog/:id", "GET", s.catalogPage)
	router.AddRoute(catalogRoute)

	contentRoute := newRoute("/content/:id", "GET", s.contentPage)
	router.AddRoute(contentRoute)

	aboutRoute := newRoute("/about/", "GET", s.aboutPage)
	router.AddRoute(aboutRoute)

	contactRoute := newRoute("/contact/", "GET", s.contactPage)
	router.AddRoute(contactRoute)

	noFoundRoute := newRoute("/404.html", "GET", s.noFoundPage)
	router.AddRoute(noFoundRoute)

	statusRoute := newRoute("/maintain/status", "GET", s.statusAction)
	router.AddRoute(statusRoute)

	loginRoute := newRoute("/maintain/login", "POST", s.loginAction)
	router.AddRoute(loginRoute)

	logoutRoute := newRoute("/maintain/logout", "POST", s.logoutAction)
	router.AddRoute(logoutRoute)
}

// Teardown 销毁
func (s *Blog) Teardown() {
	if s.centerAgent != nil {
		s.centerAgent.Stop()
	}
}

func (s *Blog) getIndexView() (model.SummaryView, bool) {
	for _, v := range s.blogContent {
		if v.Name == "Index" && v.Type == model.CATALOG {
			return v, true
		}
	}

	return model.SummaryView{}, false
}

func (s *Blog) getCatalogView() (model.SummaryView, bool) {
	for _, v := range s.blogContent {
		if v.Name == "Catalog" && v.Type == model.CATALOG {
			return v, true
		}
	}

	return model.SummaryView{}, false
}

func (s *Blog) getAboutView() (model.SummaryView, bool) {
	for _, v := range s.blogContent {
		if v.Name == "About" && v.Type == model.ARTICLE {
			return v, true
		}
	}

	return model.SummaryView{}, false
}

func (s *Blog) getContactView() (model.SummaryView, bool) {
	for _, v := range s.blogContent {
		if v.Name == "Contact" && v.Type == model.ARTICLE {
			return v, true
		}
	}

	return model.SummaryView{}, false
}

func (s *Blog) get404View() (model.SummaryView, bool) {
	for _, v := range s.blogContent {
		if v.Name == "404" && v.Type == model.ARTICLE {
			return v, true
		}
	}

	return model.SummaryView{}, false
}

func (s *Blog) mainPage(res http.ResponseWriter, req *http.Request) {
	log.Print("mainPage")

	summary := []model.SummaryView{}
	indexView, ok := s.getIndexView()
	if ok {
		summary = s.centerAgent.QuerySummary(indexView.ID)
		block, err := json.Marshal(summary)
		if err == nil {
			res.Write(block)
			return
		}

		log.Print("mainPage, json.Marshal, failed, err:" + err.Error())
	}

	http.Redirect(res, req, "/default/index.html", http.StatusMovedPermanently)
}

func (s *Blog) catalogListPage(res http.ResponseWriter, req *http.Request) {
	log.Print("catalogListPage")

	summary := []model.SummaryView{}
	catalogView, ok := s.getCatalogView()
	if ok {
		summary = s.centerAgent.QuerySummary(catalogView.ID)
		block, err := json.Marshal(summary)
		if err == nil {
			res.Write(block)
			return
		}

		log.Print("catalogListPage, json.Marshal, failed, err:" + err.Error())
	}

	http.Redirect(res, req, "/default/catalog.html", http.StatusMovedPermanently)
}

func (s *Blog) catalogPage(res http.ResponseWriter, req *http.Request) {
	log.Print("catalogPage")

	_, value := net.SplitRESTAPI(req.URL.Path)
	id, err := strconv.Atoi(value)
	if err == nil {
		summary := s.centerAgent.QuerySummary(id)
		block, err := json.Marshal(summary)
		if err == nil {
			res.Write(block)
			return
		}

		log.Print("catalogPage, json.Marshal, failed, err:" + err.Error())
	} else {
		log.Print("catalogPage, strconv.atoi failed, err:" + err.Error())
	}

	http.Redirect(res, req, "/404.html", http.StatusMovedPermanently)
}

func (s *Blog) contentPage(res http.ResponseWriter, req *http.Request) {
	log.Print("contentPage")

	_, value := net.SplitRESTAPI(req.URL.Path)
	id, err := strconv.Atoi(value)
	if err == nil {
		article, ok := s.centerAgent.QueryArticle(id)
		if ok {
			block, err := json.Marshal(article)
			if err == nil {
				res.Write(block)
				return
			}

			log.Print("contentPage, json.Marshal, failed, err:" + err.Error())
		} else {
			log.Printf("contentPage, nofound article content, id:%d", id)
		}
	}

	http.Redirect(res, req, "/404.html", http.StatusMovedPermanently)
}

func (s *Blog) aboutPage(res http.ResponseWriter, req *http.Request) {
	log.Print("aboutPage")

	aboutView, ok := s.getAboutView()
	if ok {
		article, ok := s.centerAgent.QueryArticle(aboutView.ID)
		if ok {
			block, err := json.Marshal(article)
			if err == nil {
				res.Write(block)
				return
			}

			log.Print("aboutPage, json.Marshal, failed, err:" + err.Error())
		} else {
			log.Printf("aboutPage, nofound about content, id:%d", aboutView.ID)
		}
	}

	http.Redirect(res, req, "/default/about.html", http.StatusMovedPermanently)
}

func (s *Blog) contactPage(res http.ResponseWriter, req *http.Request) {
	log.Print("contactPage")

	contactView, ok := s.getContactView()
	if ok {
		article, ok := s.centerAgent.QueryArticle(contactView.ID)

		if ok {
			block, err := json.Marshal(article)
			if err == nil {
				res.Write(block)
				return
			}

			log.Print("contactPage, json.Marshal, failed, err:" + err.Error())
		} else {
			log.Printf("contactPage, nofound contact content, id:%d", contactView.ID)
		}
	}

	http.Redirect(res, req, "/default/contact.html", http.StatusMovedPermanently)
}

func (s *Blog) noFoundPage(res http.ResponseWriter, req *http.Request) {
	log.Print("noFoundPage")

	noFoundView, ok := s.get404View()
	if ok {
		article, ok := s.centerAgent.QueryArticle(noFoundView.ID)

		if ok {
			block, err := json.Marshal(article)
			if err == nil {
				res.Write(block)
				return
			}

			log.Print("noFoundPage, json.Marshal, failed, err:" + err.Error())
		} else {
			log.Printf("noFoundPage, nofound 404 content, id:%d", noFoundView.ID)
		}
	}

	http.Redirect(res, req, "/default/404.html", http.StatusMovedPermanently)
}

func (s *Blog) statusAction(res http.ResponseWriter, req *http.Request) {
	log.Print("statusAction")

	type statusResult struct {
		common_result.Result
		OnlineUser model.AccountOnlineView `json:"onlineUser"`
	}

	result := statusResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		userView, ok := s.centerAgent.StatusAccount(authToken, sessionID)
		if !ok {
			log.Print("statusAccount failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		result.OnlineUser = userView
		result.ErrorCode = common_result.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Blog) loginAction(res http.ResponseWriter, req *http.Request) {
	log.Print("loginAction")

	type loginParam struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	type loginResult struct {
		common_result.Result
		OnlineUser model.AccountOnlineView `json:"onlineUser"`
	}

	param := loginParam{}
	result := loginResult{}
	for {
		err := net.ParsePostJSON(req, param)
		if err != nil {
			log.Printf("ParsePostJSON failed, err:%s", err.Error())
			result.ErrorCode = common_result.Failed
			result.Reason = "非法请求"
			break
		}

		userView, ok := s.centerAgent.LoginAccount(param.Account, param.Password)
		if !ok {
			log.Print("login failed, illegal account or password")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效账号或密码"
			break
		}

		result.OnlineUser = userView
		result.ErrorCode = common_result.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Blog) logoutAction(res http.ResponseWriter, req *http.Request) {
	log.Print("logoutAction")

	type logoutResult struct {
		common_result.Result
	}

	result := logoutResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		ok := s.centerAgent.LogoutAccount(authToken, sessionID)
		if !ok {
			log.Print("logout failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		result.ErrorCode = common_result.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
