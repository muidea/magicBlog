package main

import (
	"flag"

	engine "muidea.com/magicEngine"
)

var serviceAddress = ":8010"
var centerServer = "http://127.0.0.1:8080/api/remote"

func main() {
	flag.StringVar(&serviceAddress, "Address", serviceAddress, "blogService listen address")
	flag.StringVar(&centerServer, "CenterSvr", centerServer, "magicCenter server")
	flag.Parse()

	router := engine.NewRouter()

	//test.Append(router)

	svr := engine.NewHTTPServer(serviceAddress)
	svr.Bind(router)

	//svr.Use(&test.Hello{})

	//svr.Use(&test.Test{})

	svr.Run()
}
