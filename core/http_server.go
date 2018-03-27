package core

import (
	"log"
	"net/http"
)

// HTTPServer HTTPServer
type HTTPServer interface {
	Bind(router Router)
	Run()
}

type httpServer struct {
	listentAddr string
	router      Router
}

// NewHTTPServer 新建HTTPServer
func NewHTTPServer(listentAddr string) HTTPServer {
	return &httpServer{listentAddr: listentAddr}
}

func (s *httpServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Printf("url:%s, method:%s", req.URL.Path, req.Method)
}

func (s *httpServer) Bind(router Router) {
	s.router = router
}

func (s *httpServer) Run() {
	traceInfo("listening on " + s.listentAddr)

	err := http.ListenAndServe(s.listentAddr, s)
	log.Fatalf("run httpserver fatal, err:%s", err.Error())
}
