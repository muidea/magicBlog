package core

import (
	"log"
	"net/http"

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
func NewBlog(centerServer, name, account, password string) (Blog, bool) {
	blog := Blog{centerAgent: NewCenterAgent()}

	agent := NewCenterAgent()
	if !agent.Start(centerServer, name, account, password) {
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

	blog.centerAgent = agent
	blog.blogCatalogView = blogCatalog

	return blog, true
}

// Blog Blog对象
type Blog struct {
	centerAgent     Agent
	blogCatalogView model.CatalogDetailView
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

	aboutRoute := newRoute("/about", "GET", s.aboutPage)
	router.AddRoute(aboutRoute)
}

// Teardown 销毁
func (s *Blog) Teardown() {
	if s.centerAgent != nil {
		s.centerAgent.Stop()
	}
}

func (s *Blog) mainPage(res http.ResponseWriter, req *http.Request) {
	log.Print("mainPage")
}

func (s *Blog) catalogListPage(res http.ResponseWriter, req *http.Request) {
	log.Print("catalogListPage")
}

func (s *Blog) catalogPage(res http.ResponseWriter, req *http.Request) {
	log.Print("catalogPage")
}

func (s *Blog) contentPage(res http.ResponseWriter, req *http.Request) {
	log.Print("contentPage")
}

func (s *Blog) aboutPage(res http.ResponseWriter, req *http.Request) {
	log.Print("aboutPage")
}
