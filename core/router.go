package core

import "net/http"

// 基本HTTP行为定义
const (
	GET     = "GET"
	PUT     = "PUT"
	POST    = "POST"
	DELETE  = "DELETE"
	OPTIONS = "OPTIONS"
)

// Route 路由接口
type Route interface {
	// Action 路由行为GET/PUT/POST/DELETE
	Method() string
	// Pattern 路由规则, 以'/'开始
	Pattern() string
	// Handler 路由处理器
	Handler() interface{}
}

// Router 路由器对象
type Router interface {
	// 增加路由
	AddRoute(rt Route, filters ...MiddleWareHandler)
	// 清除路由
	RemoveRoute(rt Route)
	// 分发一条请求
	Handle(ctx Context, res http.ResponseWriter, req *http.Request)
}

// 路由对象
type route struct {
}

type router struct {
}

func (s *router) AddRoute(rt Route, filters ...MiddleWareHandler) {

}

func (s *router) RemoveRoute(rt Route) {

}

func (s *router) Handle(ctx RequestContext, res ResponseWriter, req *http.Request) {

}
