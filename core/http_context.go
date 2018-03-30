package core

import "net/http"

// Context represents a request context. Services can be mapped on the request level from this interface.
type Context interface {
	// Next is an optional function that Middleware Handlers can call to yield the until after
	// the other Handlers have been executed. This works really well for any operations that must
	// happen after an http request
	Next()
	// Written returns whether or not the response for this context has been written.
	Written() bool

	Run()
}

type context struct {
	filters []MiddleWareHandler
	router  Router
	rw      ResponseWriter
	req     *http.Request
	index   int
}

// NewContext 新建Context
func NewContext(filters []MiddleWareHandler, router Router, res http.ResponseWriter, req *http.Request) Context {
	return &context{filters: filters, router: router, rw: NewResponseWriter(res), req: req, index: 0}
}

func (c *context) Next() {
	c.index++
	c.Run()
}

func (c *context) Written() bool {
	return c.rw.Written()
}

func (c *context) Run() {
	lenSize := len(c.filters)
	for c.index < lenSize {
		handler := c.filters[c.index]
		handler.Handle(c, c.rw, c.req)

		c.index++
		if c.Written() {
			return
		}
	}

	c.router.Handle(c.rw, c.req)
}
