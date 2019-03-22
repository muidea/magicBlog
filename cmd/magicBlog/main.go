package main

import (
	"flag"
	"log"

	magicblog "github.com/muidea/magicBlog/core"
	engine "github.com/muidea/magicEngine"
)

var bindPort = "8866"
var centerServer = "127.0.0.1:8888"
var endpointName = "magicBlog"
var endpointID = "f0e078a8-6de8-4273-88a4-dccef60ff88f"
var authToken = "yTtWiuuoGifPVfcK5Mf4mdu8mGl78E3y"

func main() {
	flag.StringVar(&bindPort, "ListenPort", bindPort, "magicBlog listen address")
	flag.StringVar(&centerServer, "CenterSvr", centerServer, "magicCenter server")
	flag.StringVar(&endpointName, "EndpointName", endpointName, "magicBlog endpoint name.")
	flag.StringVar(&endpointID, "EndpointID", endpointID, "magicBlog endpoint id")
	flag.StringVar(&authToken, "AuthToken", authToken, "magicBlog authtoken")
	flag.Parse()

	router := engine.NewRouter()

	blog, ok := magicblog.New(centerServer, endpointName, endpointID, authToken)
	if ok {
		blog.Startup(router)

		svr := engine.NewHTTPServer(bindPort)
		svr.Bind(router)

		svr.Run()
	} else {
		log.Printf("new Blog failed.")
	}

	blog.Teardown()
}
