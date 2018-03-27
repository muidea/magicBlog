package main

import "muidea.com/magicBlog/core"

func main() {

	svr := core.NewHTTPServer(":8010")

	svr.Run()
}
