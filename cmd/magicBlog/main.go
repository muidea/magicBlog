package main

import (
	"flag"
	"log"

	magicblog "muidea.com/magicBlog/core"
	engine "muidea.com/magicEngine"
)

var serviceAddress = ":8018"
var centerServer = "http://127.0.0.1:8888"
var blogName = "magicBlog"
var account = "magicBlog"
var password = "123"

func main() {
	flag.StringVar(&serviceAddress, "Address", serviceAddress, "blogService listen address")
	flag.StringVar(&centerServer, "CenterSvr", centerServer, "magicCenter server")
	flag.StringVar(&blogName, "BlogName", blogName, "blog name.")
	flag.StringVar(&account, "Account", account, "magicBlog account")
	flag.StringVar(&password, "Password", password, "magicBlog password")
	flag.Parse()

	router := engine.NewRouter()

	blog, ok := magicblog.NewBlog(centerServer, blogName, account, password)
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
