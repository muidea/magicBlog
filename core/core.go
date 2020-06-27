package core

import (
	"log"

	engine "github.com/muidea/magicEngine"

	batisClient "github.com/muidea/magicBatis/client"
	casToolkit "github.com/muidea/magicCas/toolkit"
	"github.com/muidea/magicCommon/session"

	"github.com/muidea/magicBlog/config"
	"github.com/muidea/magicBlog/core/service/blog"
	"github.com/muidea/magicBlog/model"
)

// New 新建Core
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

	blogService := blog.New()

	core := &Core{
		blogService: blogService,
	}

	return core, nil
}

// Core Core对象
type Core struct {
	blogService *blog.Blog
}

// Startup 启动
func (s *Core) Startup(router engine.Router) {
	sessionRegistry := session.CreateRegistry(nil)
	casRouteRegistry := casToolkit.NewCasRegistry(s.blogService, router)

	s.blogService.BindRegistry(sessionRegistry, casRouteRegistry)

	s.blogService.RegisterHandler(router)
}

// Teardown 销毁
func (s *Core) Teardown() {
}
