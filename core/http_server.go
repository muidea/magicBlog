package core

import (
	"log"
	"net/http"
)

// HTTPServer HTTPServer
type HTTPServer interface {
	Use(handler MiddleWareHandler)
	Bind(router Router)
	Run()
}

type httpServer struct {
	listenAddr string
	router     Router
	filter     MiddleWareChains
}

// NewHTTPServer 新建HTTPServer
func NewHTTPServer(listenAddr string) HTTPServer {
	return &httpServer{listenAddr: listenAddr, filter: NewMiddleWareChains()}
}

func (s *httpServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	ctx := NewRequestContext(s.filter.GetHandlers(), s.router, res, req)

	ctx.Run()
}

func (s *httpServer) Use(handler MiddleWareHandler) {
	s.filter.Append(handler)
}

func (s *httpServer) Bind(router Router) {
	s.router = router
}

func (s *httpServer) Run() {
	traceInfo("listening on " + s.listenAddr)

	err := http.ListenAndServe(s.listenAddr, s)
	log.Fatalf("run httpserver fatal, err:%s", err.Error())
}
