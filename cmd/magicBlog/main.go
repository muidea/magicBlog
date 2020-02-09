package main

import (
	"flag"
	"log"

	engine "github.com/muidea/magicEngine"
	_ "github.com/muidea/magicBlog/config"
	core "github.com/muidea/magicBlog/core"
)

var listenPort = "8880"

func main() {
	flag.StringVar(&listenPort, "ListenPort", listenPort, "magicBlog listen address")
	flag.Parse()

	router := engine.NewRouter()
	core, err := core.New()

	if err == nil {
		core.Startup(router)

		svr := engine.NewHTTPServer(listenPort)
		svr.Bind(router)

		svr.Run()
	} else {
		log.Printf("start magicBlog failed.")
	}

	core.Teardown()
}
