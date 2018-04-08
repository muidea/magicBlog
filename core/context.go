package core

import (
	"log"
	"net/http"
	"reflect"
)

// Context request base context
type Context interface {
	SetData(key string, value interface{})
	GetData(key string) (interface{}, bool)
}

// RequestContext represents a request context. Services can be mapped on the request level from this interface.
type RequestContext interface {
	Context
	// Next is an optional function that Middleware Handlers can call to yield the until after
	// the other Handlers have been executed. This works really well for any operations that must
	// happen after an http request
	Next()
	// Written returns whether or not the response for this context has been written.
	Written() bool

	Run()
}

// ValidateMiddleWareHandler 校验MiddleWareHandler
func ValidateMiddleWareHandler(handler interface{}) {
	log.Print("ValidateMiddleWareHandler")
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Ptr {
		panic("middleware handler must be a callable interface")
	}

	handlerMethod, ok := handlerType.MethodByName("Handle")
	if !ok {
		panic("middleware handler isn\\'t have Handle func")
	}

	methodType := handlerMethod.Type
	paramNum := methodType.NumIn()
	if paramNum != 4 {
		panic("middleware handler invalid handle func param number")
	}

	// param0 := methodType.In(0).String()
	param1 := methodType.In(1)
	if param1.Kind() != reflect.Interface {
		panic("middleware handler invalid handle func param0 type")
	}
	if param1.Name() != "RequestContext" {
		panic("middleware handler invalid handle func param0 type")
	}
	param2 := methodType.In(2)
	if param2.Kind() != reflect.Interface {
		panic("middleware handler invalid handle func param1 type")
	}
	if param2.String() != "http.ResponseWriter" {
		panic("middleware handler invalid handle func param1 type")
	}

	param3 := methodType.In(3)
	if param3.Kind() != reflect.Ptr {
		panic("middleware handler invalid handle func param2 type")
	}
	if param3.String() != "*http.Request" {
		panic("middleware handler invalid handle func param2 type")
	}
}

// InvokeMiddleWareHandler 执行MiddleWareHandle
func InvokeMiddleWareHandler(handler interface{}, ctx RequestContext, res http.ResponseWriter, req *http.Request) {
	log.Print("InvokeMiddleWareHandler")
	params := make([]reflect.Value, 4)
	params[0] = reflect.ValueOf(handler)
	params[1] = reflect.ValueOf(ctx)
	params[2] = reflect.ValueOf(res)
	params[3] = reflect.ValueOf(req)

	handlerType := reflect.TypeOf(handler)
	// 已经验证通过，所以这里就不用继续判断
	//if handlerType.Kind() != reflect.Ptr {
	//	panic("middleware handler must be a callable interface")
	//}

	handlerMethod, ok := handlerType.MethodByName("Handle")
	if !ok {
		panic("middleware handler isn\\'t have Handle func")
	}

	fv := handlerMethod.Func
	fv.Call(params)
}

// ValidateRouteHandler 校验RouteHandler
func ValidateRouteHandler(handler interface{}) {
	log.Print("ValidateRouteHandler")
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		panic("route handler must be a callable func")
	}

	paramNum := handlerType.NumIn()
	if paramNum == 3 {
		param0 := handlerType.In(0)
		if param0.Kind() != reflect.Interface {
			panic("route handler invalid handle func param0 type")
		}
		if param0.Name() != "RequestContext" {
			panic("route handler invalid handle func param0 type")
		}
		param1 := handlerType.In(1)
		if param1.Kind() != reflect.Interface {
			panic("route handler invalid handle func param1 type")
		}
		if param1.String() != "http.ResponseWriter" {
			panic("route handler invalid handle func param1 type")
		}

		param2 := handlerType.In(2)
		if param2.Kind() != reflect.Ptr {
			panic("route handler invalid handle func param2 type")
		}
		if param2.String() != "*http.Request" {
			panic("route handler invalid handle func param2 type")
		}
	} else if paramNum == 2 {
		param0 := handlerType.In(0)
		if param0.Kind() != reflect.Interface {
			panic("route handler invalid handle func param0 type")
		}
		if param0.String() != "http.ResponseWriter" {
			panic("route handler invalid handle func param0 type")
		}

		param1 := handlerType.In(1)
		if param1.Kind() != reflect.Ptr {
			panic("route handler invalid handle func param0 type")
		}
		if param1.String() != "*http.Request" {
			panic("route handler invalid handle func param0 type")
		}
	} else {
		panic("illegal callable func")
	}
}

// InvokeRouteHandler 执行RouteHandle
func InvokeRouteHandler(handler interface{}, ctx Context, res http.ResponseWriter, req *http.Request) {
	log.Print("InvokeRouteHandler")
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		panic("route handler must be a callable func")
	}

	var params []reflect.Value

	paramNum := handlerType.NumIn()
	if paramNum == 3 {
		params = make([]reflect.Value, 3)
		params[0] = reflect.ValueOf(ctx)
		params[1] = reflect.ValueOf(res)
		params[2] = reflect.ValueOf(req)
	} else if paramNum == 2 {
		params = make([]reflect.Value, 2)
		params[0] = reflect.ValueOf(res)
		params[1] = reflect.ValueOf(req)
	} else {
		panic("illegal callable func")
	}

	fv := reflect.ValueOf(handler)
	fv.Call(params)
}

type requestContext struct {
	filters []MiddleWareHandler
	rw      ResponseWriter
	req     *http.Request
	index   int

	router  Router
	dataMap map[string]interface{}
}

// NewRequestContext 新建Context
func NewRequestContext(filters []MiddleWareHandler, router Router, res http.ResponseWriter, req *http.Request) RequestContext {
	return &requestContext{filters: filters, router: router, rw: NewResponseWriter(res), req: req, index: 0, dataMap: make(map[string]interface{})}
}

func (c *requestContext) Next() {
	c.index++
	c.Run()
}

func (c *requestContext) Written() bool {
	return c.rw.Written()
}

func (c *requestContext) Run() {
	totalSizxe := len(c.filters)
	for c.index < totalSizxe {
		handler := c.filters[c.index]
		InvokeMiddleWareHandler(handler, c, c.rw, c.req)

		c.index++
		if c.Written() {
			return
		}
	}

	if !c.Written() && c.router != nil {
		c.router.Handle(c, c.rw, c.req)
	} else {
		http.NotFound(c.rw, c.req)
	}
}

func (c *requestContext) SetData(key string, value interface{}) {
	c.dataMap[key] = value
}

func (c *requestContext) GetData(key string) (interface{}, bool) {
	val, ok := c.dataMap[key]
	return val, ok
}

type routeContext struct {
	filters []MiddleWareHandler
	rw      ResponseWriter
	req     *http.Request
	index   int

	route   Route
	context Context
}

// NewRouteContext 新建Context
func NewRouteContext(reqCtx Context, filters []MiddleWareHandler, route Route, res http.ResponseWriter, req *http.Request) RequestContext {
	return &routeContext{filters: filters, route: route, rw: NewResponseWriter(res), req: req, index: 0, context: reqCtx}
}

func (c *routeContext) Next() {
	c.index++
	c.Run()
}

func (c *routeContext) Written() bool {
	return c.rw.Written()
}

func (c *routeContext) Run() {
	totalSizxe := len(c.filters)
	for c.index < totalSizxe {
		handler := c.filters[c.index]
		InvokeMiddleWareHandler(handler, c, c.rw, c.req)

		c.index++
		if c.Written() {
			return
		}
	}

	if !c.Written() {
		InvokeRouteHandler(c.route.Handler(), c.context, c.rw, c.req)
	}
}

func (c *routeContext) SetData(key string, value interface{}) {
	c.context.SetData(key, value)
}

func (c *routeContext) GetData(key string) (interface{}, bool) {
	return c.context.GetData(key)
}
