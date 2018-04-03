package main

import (
	"fmt"
	"reflect"

	"muidea.com/magicBlog/core"
	"muidea.com/magicBlog/test"
)

func validateHandler(handler interface{}) {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		panic("middleware handler must be a callable func")
	}

	paramNum := handlerType.NumIn()
	param0 := handlerType.In(0).Name()

	fmt.Printf("param Num:%d\n", paramNum)
	fmt.Printf("param name:%s\n", param0)
}

func testFunc(val string) {
	fmt.Print(val)
}

func main() {

	validateHandler(testFunc)

	svr := core.NewHTTPServer(":8010")

	svr.Use(&test.Hello{})

	svr.Use(&test.Test{})

	svr.Run()
}
