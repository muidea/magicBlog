package magicblog

import (
	"log"
	"net/http"

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
func NewBlog(centerServer, account, authToken string) Blog {
	return Blog{}
}

// Blog Blog对象
type Blog struct {
}

// Startup 启动
func (s *Blog) Startup(router engine.Router) {
	mainRoute := newRoute("/", "GET", s.mainPage)
	router.AddRoute(mainRoute)

	detailRoute := newRoute("/detail/:id", "GET", s.detailPage)
	router.AddRoute(detailRoute)

	catalogRoute := newRoute("/catalog/:id", "GET", s.catalogPage)
	router.AddRoute(catalogRoute)
}

// Teardown 销毁
func (s *Blog) Teardown() {

}

func (s *Blog) mainPage(res http.ResponseWriter, req *http.Request) {
	log.Print("mainPage")
}

func (s *Blog) catalogPage(res http.ResponseWriter, req *http.Request) {
	log.Print("catalogPage")
}

func (s *Blog) detailPage(res http.ResponseWriter, req *http.Request) {
	log.Print("detailPage")
}
