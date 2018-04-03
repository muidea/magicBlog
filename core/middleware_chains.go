package core

import (
	"net/http"
	"reflect"
	"sync"
)

// MiddleWareHandler 中间件处理器
type MiddleWareHandler interface {
	Handle(ctx RequestContext, res http.ResponseWriter, req *http.Request)
}

// MiddleWareChains 处理器链
type MiddleWareChains interface {
	Append(handler MiddleWareHandler)

	GetHandlers() []MiddleWareHandler
}

type chainsImpl struct {
	handlers    []MiddleWareHandler
	handlesLock sync.RWMutex
}

func validateHandler(handler interface{}) {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic("middleware handler must be a callable func")
	}

	reflect.TypeOf(handler).NumField()
}

// NewMiddleWareChains 新建MiddleWareChains
func NewMiddleWareChains() MiddleWareChains {
	return &chainsImpl{handlers: []MiddleWareHandler{}}
}

func (s *chainsImpl) GetHandlers() []MiddleWareHandler {
	s.handlesLock.RLock()
	defer s.handlesLock.RUnlock()

	return s.handlers[:]
}

func (s *chainsImpl) Append(handler MiddleWareHandler) {
	s.handlesLock.Lock()
	defer s.handlesLock.Unlock()

	s.handlers = append(s.handlers, handler)
}
