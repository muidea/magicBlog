package main

import (
	"muidea.com/magicBlog/core"
	"muidea.com/magicBlog/test"
)

func main() {

	router := core.NewRouter()

	test.Append(router)

	svr := core.NewHTTPServer(":8010")
	svr.Bind(router)

	svr.Use(&test.Hello{})

	//svr.Use(&test.Test{})

	svr.Run()
}
