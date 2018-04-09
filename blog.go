package magicblog

import (
	"net/http"

	engine "muidea.com/magicEngine"
)

type blog struct {
}

func (s *blog) Startup(router engine.Router) {

}

func (s *blog) Teardown() {

}

func (s *blog) MainPage(res http.ResponseWriter, req *http.Request) {

}

func (s *blog) CatalogPage(res http.ResponseWriter, req *http.Request) {

}

func (s *blog) DetailPage(res http.ResponseWriter, req *http.Request) {

}
