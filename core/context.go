package core

import "net/http"

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
		handler.Handle(c, c.rw, c.req)

		c.index++
		if c.Written() {
			return
		}
	}

	if !c.Written() {
		c.router.Handle(c, c.rw, c.req)
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
func NewRouteContext(filters []MiddleWareHandler, route Route, res ResponseWriter, req *http.Request, reqCtx Context) RequestContext {
	return &routeContext{filters: filters, route: route, rw: res, req: req, index: 0, context: reqCtx}
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
		handler.Handle(c, c.rw, c.req)

		c.index++
		if c.Written() {
			return
		}
	}

	if !c.Written() {

	}
}

func (c *routeContext) SetData(key string, value interface{}) {
	c.context.SetData(key, value)
}

func (c *routeContext) GetData(key string) (interface{}, bool) {
	return c.context.GetData(key)
}
