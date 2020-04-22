package core

import (
	"fmt"
	"log"

	batisClient "github.com/muidea/magicBatis/client"
	"github.com/muidea/magicBlog/config"
	"github.com/muidea/magicBlog/core/route"
	"github.com/muidea/magicBlog/model"
	"github.com/muidea/magicCommon/session"

	engine "github.com/muidea/magicEngine"
)

// New 新建Protal
func New() (*Core, error) {
	var err error
	clnt := batisClient.NewClient(config.BatisService(), config.EndpointName())
	defer func() {
		if err != nil {
			clnt.Release()
		}
	}()
	err = model.InitializeModel(clnt)
	if err != nil {
		log.Printf("initializeModel failed, err:%s", err.Error())
		return nil, err
	}

	sessionRegistry := session.CreateRegistry(nil)

	routeRegister := route.NewRoute(sessionRegistry)
	if routeRegister == nil {
		return nil, fmt.Errorf("NewRoute failed")
	}

	core := &Core{}
	core.routeRegister = routeRegister

	return core, nil
}

// Core Protal对象
type Core struct {
	routeRegister *route.Registry
}

// Startup 启动
func (s *Core) Startup(router engine.Router) {
	s.routeRegister.RegisterRoute(router)
}

// Teardown 销毁
func (s *Core) Teardown() {
}
