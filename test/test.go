package test

import (
	"log"
	"net/http"

	"muidea.com/magicBlog/core"
)

// Test hello middleware
type Test struct {
}

// Handle handle request
func (s *Test) Handle(ctx core.RequestContext, res http.ResponseWriter, req *http.Request) {
	log.Print("Test Handle")
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Write([]byte("test world"))
}
