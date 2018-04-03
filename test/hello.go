package test

import (
	"log"
	"net/http"

	"muidea.com/magicBlog/core"
)

// Hello hello middleware
type Hello struct {
}

// Handle handle request
func (s *Hello) Handle(ctx core.RequestContext, res http.ResponseWriter, req *http.Request) {
	log.Print("Hello Handle")
}
