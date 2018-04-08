package test

import (
	"log"
	"net/http"

	"muidea.com/magicBlog/core"
)

// Append append router
func Append(router core.Router) {
	router.AddRoute(&getRoute{})
}

type getRoute struct {
}

func (s *getRoute) Method() string {
	return "GET"
}

func (s *getRoute) Pattern() string {
	return "/demo/:id"
}

func (s *getRoute) Handler() interface{} {
	return s.getDemo
}

func (s *getRoute) getDemo(res http.ResponseWriter, req *http.Request) {
	log.Print("getDemo....")
}
