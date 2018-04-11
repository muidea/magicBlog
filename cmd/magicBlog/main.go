package main

import (
	"flag"

	magicblog "muidea.com/magicBlog"
	engine "muidea.com/magicEngine"
)

var serviceAddress = ":8010"
var centerServer = "http://127.0.0.1:8080/api/remote"
var account = "test"
var authToken = "token"

func main() {
	flag.StringVar(&serviceAddress, "Address", serviceAddress, "blogService listen address")
	flag.StringVar(&centerServer, "CenterSvr", centerServer, "magicCenter server")
	flag.StringVar(&account, "Account", account, "magicCenter server")
	flag.StringVar(&authToken, "AuthToken", authToken, "magicCenter server")
	flag.Parse()

	router := engine.NewRouter()

	blog := magicblog.NewBlog(centerServer, account, authToken)

	blog.Startup(router)

	svr := engine.NewHTTPServer(serviceAddress)
	svr.Bind(router)

	svr.Run()

	blog.Teardown()
}
