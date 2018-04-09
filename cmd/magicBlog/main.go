package main

import (
	"muidea.com/magicBlog/test"
	engine "muidea.com/magicEngine"
)

func main() {

	router := engine.NewRouter()

	test.Append(router)

	svr := engine.NewHTTPServer(":8010")
	svr.Bind(router)

	svr.Use(&test.Hello{})

	svr.Use(&test.Test{})

	svr.Run()
}
