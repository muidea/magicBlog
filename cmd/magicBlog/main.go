package main

import (
	"flag"
	"log"

	magicblog "muidea.com/magicBlog/core"
	engine "muidea.com/magicEngine"
)

var serviceAddress = ":8866"
var centerServer = "http://127.0.0.1:8888"
var blogName = "magicBlog"
var endpointID = "f0e078a8-6de8-4273-88a4-dccef60ff88f"
var authToken = "yTtWiuuoGifPVfcK5Mf4mdu8mGl78E3y"

func main() {
	flag.StringVar(&serviceAddress, "Address", serviceAddress, "blogService listen address")
	flag.StringVar(&centerServer, "CenterSvr", centerServer, "magicCenter server")
	flag.StringVar(&blogName, "BlogName", blogName, "blog name.")
	flag.StringVar(&endpointID, "EndpointID", endpointID, "magicBlog endpoint id")
	flag.StringVar(&authToken, "AuthToken", authToken, "magicBlog authtoken")
	flag.Parse()

	router := engine.NewRouter()

	blog, ok := magicblog.NewBlog(centerServer, blogName, endpointID, authToken)
	if ok {
		blog.Startup(router)

		svr := engine.NewHTTPServer(serviceAddress)
		svr.Bind(router)

		svr.Run()
	} else {
		log.Printf("new Blog failed.")
	}

	blog.Teardown()
}
