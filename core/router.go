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
	AddRoute(rt Route)
	// 清除路由
	RemoveRoute(rt Route)

	// 分发一条请求
	Dispatch(res http.ResponseWriter, req *http.Request)
}
