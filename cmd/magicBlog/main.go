package main

import (
	"flag"

	magicblog "muidea.com/magicBlog/core"
	engine "muidea.com/magicEngine"
)

var serviceAddress = ":8018"
var centerServer = "http://127.0.0.1:8888/api/remote"
var account = "test"
var password = "token"

func main() {
	flag.StringVar(&serviceAddress, "Address", serviceAddress, "blogService listen address")
	flag.StringVar(&centerServer, "CenterSvr", centerServer, "magicCenter server")
	flag.StringVar(&account, "Account", account, "magicBlog account")
	flag.StringVar(&password, "Password", password, "magicBlog password")
	flag.Parse()

	router := engine.NewRouter()

	blog := magicblog.NewBlog(centerServer, account, password)

	blog.Startup(router)

	svr := engine.NewHTTPServer(serviceAddress)
	svr.Bind(router)

	svr.Run()

	blog.Teardown()
}
