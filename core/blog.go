package core

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"muidea.com/magicCommon/agent"
	common_def "muidea.com/magicCommon/def"
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

// New 新建Blog
func New(centerServer, name, endpointID, authToken string) (Blog, bool) {
	blog := Blog{}

	agent := agent.New()
	sessionID, ok := agent.Start(centerServer, endpointID, authToken)
	if !ok {
		return blog, false
	}
	blogCatalog, ok := agent.FetchSummary(name, model.CATALOG, authToken, sessionID)
	if !ok {
		_, ok = agent.CreateCatalog(name, "magicBlog auto create catalog.", []model.Catalog{}, authToken, sessionID)
		if !ok {
			log.Print("create blog root catalog failed.")
			return blog, false
		}
	}
	blogCatalog, ok = agent.FetchSummary(name, model.CATALOG, authToken, sessionID)
	if !ok {
		log.Print("fetch blog root ctalog failed.")
		return blog, false
	}

	blogContent := agent.QuerySummaryContent(blogCatalog.ID, model.CATALOG, authToken, sessionID)

	blog.centerAgent = agent
	blog.blogInfo = blogCatalog
	blog.blogContent = blogContent
	blog.endpointID = endpointID
	blog.authToken = authToken
	blog.sessionID = sessionID

	return blog, true
}

// Blog Blog对象
type Blog struct {
	centerAgent agent.Agent
	blogInfo    model.SummaryView
	blogContent []model.SummaryView
	endpointID  string
	authToken   string
	sessionID   string
}

// Startup 启动
func (s *Blog) Startup(router engine.Router) {
	mainRoute := newRoute("/", "GET", s.mainPage)
	router.AddRoute(mainRoute)

	catalogSummaryRoute := newRoute("/catalog/", "GET", s.catalogSummaryPage)
	router.AddRoute(catalogSummaryRoute)

	catalogSummaryByIDRoute := newRoute("/catalog/:id", "GET", s.catalogSummaryByIDPage)
	router.AddRoute(catalogSummaryByIDRoute)

	contentRoute := newRoute("/content/:id", "GET", s.contentPage)
	router.AddRoute(contentRoute)

	aboutRoute := newRoute("/about", "GET", s.aboutPage)
	router.AddRoute(aboutRoute)

	contactRoute := newRoute("/contact", "GET", s.contactPage)
	router.AddRoute(contactRoute)

	noFoundRoute := newRoute("/404.html", "GET", s.noFoundPage)
	router.AddRoute(noFoundRoute)

	statusRoute := newRoute("/maintain/status", "GET", s.statusAction)
	router.AddRoute(statusRoute)

	loginRoute := newRoute("/maintain/login", "POST", s.loginAction)
	router.AddRoute(loginRoute)

	logoutRoute := newRoute("/maintain/logout", "DELETE", s.logoutAction)
	router.AddRoute(logoutRoute)

	summaryRoute := newRoute("/maintain/summary", "GET", s.summaryAction)
	router.AddRoute(summaryRoute)

	catalogCreateRoute := newRoute("/maintain/catalog", "POST", s.catalogCreateAction)
	router.AddRoute(catalogCreateRoute)

	catalogQueryRoute := newRoute("/maintain/catalog/:id", "GET", s.catalogQueryAction)
	router.AddRoute(catalogQueryRoute)

	articleCreateRoute := newRoute("/maintain/article", "POST", s.articleCreateAction)
	router.AddRoute(articleCreateRoute)

	catalogUpdateRoute := newRoute("/maintain/catalog/:id", "PUT", s.catalogUpdateAction)
	router.AddRoute(catalogUpdateRoute)

	articleUpdateRoute := newRoute("/maintain/article/:id", "PUT", s.articleUpdateAction)
	router.AddRoute(articleUpdateRoute)

	catalogDeleteRoute := newRoute("/maintain/catalog/:id", "DELETE", s.catalogDeleteAction)
	router.AddRoute(catalogDeleteRoute)

	articleDeleteRoute := newRoute("/maintain/article/:id", "DELETE", s.articleDeleteAction)
	router.AddRoute(articleDeleteRoute)
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

	result := common_def.QuerySummaryListResult{}
	indexView, ok := s.getIndexView()
	if ok {
		result.Summary = s.centerAgent.QuerySummaryContent(indexView.ID, model.CATALOG, s.authToken, s.sessionID)
		result.ErrorCode = common_def.Success
	} else {
		result.ErrorCode = common_def.Redirect
		result.Reason = "/default/index.html"
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	log.Print("mainPage, json.Marshal, failed, err:" + err.Error())
}

func (s *Blog) catalogSummaryPage(res http.ResponseWriter, req *http.Request) {
	log.Print("catalogSummaryPage")

	result := common_def.QuerySummaryListResult{}
	catalogView, ok := s.getCatalogView()
	if ok {
		result.Summary = s.centerAgent.QuerySummaryContent(catalogView.ID, model.CATALOG, s.authToken, s.sessionID)
		result.ErrorCode = common_def.Success
	} else {
		result.ErrorCode = common_def.Redirect
		result.Reason = "/default/catalog.html"
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	log.Print("catalogSummaryPage, json.Marshal, failed, err:" + err.Error())
}

func (s *Blog) catalogSummaryByIDPage(res http.ResponseWriter, req *http.Request) {
	log.Print("catalogSummaryByIDPage")

	result := common_def.QuerySummaryListResult{}
	_, value := net.SplitRESTAPI(req.URL.Path)
	id, err := strconv.Atoi(value)
	if err == nil {
		result.Summary = s.centerAgent.QuerySummaryContent(id, model.CATALOG, s.authToken, s.sessionID)
		result.ErrorCode = common_def.Success
	} else {
		result.ErrorCode = common_def.IllegalParam
		result.Reason = "非法参数"
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	log.Print("catalogSummaryByIDPage, json.Marshal, failed, err:" + err.Error())
}

type contentResult struct {
	common_def.Result
	Content model.ArticleDetailView `json:"content"`
}

func (s *Blog) contentPage(res http.ResponseWriter, req *http.Request) {
	log.Print("contentPage")

	result := contentResult{}
	_, value := net.SplitRESTAPI(req.URL.Path)
	id, err := strconv.Atoi(value)
	if err == nil {
		article, ok := s.centerAgent.QueryArticle(id, s.authToken, s.sessionID)
		if ok {
			result.Content = article
			result.ErrorCode = common_def.Success
		} else {
			result.ErrorCode = common_def.NoExist
			result.Reason = "对象不存在"
		}

	} else {
		result.ErrorCode = common_def.IllegalParam
		result.Reason = "非法参数"
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	log.Print("contentPage, json.Marshal, failed, err:" + err.Error())
}

func (s *Blog) aboutPage(res http.ResponseWriter, req *http.Request) {
	log.Print("aboutPage")

	result := contentResult{}
	aboutView, ok := s.getAboutView()
	if ok {
		article, ok := s.centerAgent.QueryArticle(aboutView.ID, s.authToken, s.sessionID)
		if ok {
			result.Content = article
			result.ErrorCode = common_def.Success
		} else {
			result.ErrorCode = common_def.NoExist
			result.Reason = "对象不存在"
		}
	} else {
		result.ErrorCode = common_def.Redirect
		result.Reason = "/default/about.html"
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	log.Print("aboutPage, json.Marshal, failed, err:" + err.Error())
}

func (s *Blog) contactPage(res http.ResponseWriter, req *http.Request) {
	log.Print("contactPage")

	result := contentResult{}
	contactView, ok := s.getContactView()
	if ok {
		article, ok := s.centerAgent.QueryArticle(contactView.ID, s.authToken, s.sessionID)
		if ok {
			result.Content = article
			result.ErrorCode = common_def.Success
		} else {
			result.ErrorCode = common_def.NoExist
			result.Reason = "对象不存在"
		}
	} else {
		result.ErrorCode = common_def.Redirect
		result.Reason = "/default/contact.html"
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	log.Print("contactPage, json.Marshal, failed, err:" + err.Error())
}

func (s *Blog) noFoundPage(res http.ResponseWriter, req *http.Request) {
	log.Print("noFoundPage")

	result := contentResult{}
	noFoundView, ok := s.get404View()
	if ok {
		article, ok := s.centerAgent.QueryArticle(noFoundView.ID, s.authToken, s.sessionID)
		if ok {
			result.Content = article
			result.ErrorCode = common_def.Success
		} else {
			result.ErrorCode = common_def.NoExist
			result.Reason = "对象不存在"
		}
	} else {
		result.ErrorCode = common_def.Redirect
		result.Reason = "/default/404.html"
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}
	log.Print("noFoundPage, json.Marshal, failed, err:" + err.Error())
}
