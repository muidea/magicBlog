package main

import (
	"fmt"
	"reflect"

	"muidea.com/magicBlog/core"
	"muidea.com/magicBlog/test"
)

func validateHandler(handler interface{}) {
	fmt.Print("validateHandler...........\n")
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		panic("middleware handler must be a callable func")
	}

	paramNum := handlerType.NumIn()
	param0 := handlerType.In(0).Name()
	param1 := handlerType.In(1).Name()
	param2 := handlerType.In(2).Name()

	fmt.Printf("param Num:%d\n", paramNum)
	fmt.Printf("param0 name:%s\n", param0)
	fmt.Printf("param1 name:%s\n", param1)
	fmt.Printf("param2 name:%s\n", param2)

	fv := reflect.ValueOf(handler)
	params := make([]reflect.Value, 3)
	params[0] = reflect.ValueOf("Hello World")
	params[1] = reflect.ValueOf(123)
	params[2] = reflect.ValueOf(&testStruct{})

	fv.Call(params)
}

type testInterface interface {
	Demo()
}

type testStruct struct {
}

func (s *testStruct) Demo() {
	fmt.Printf("testStruct\n")
}

func (s *testStruct) DemoParam(val string, intVal int, testPtr testInterface) {
	fmt.Print(val)
	testPtr.Demo()
}

func testFunc(val string, intVal int, testPtr testInterface) {
	fmt.Print(val)
	testPtr.Demo()
}

func main() {

	//validateHandler(testFunc)

	//temp := testStruct{}
	//validateHandler(temp.DemoParam)

	router := core.NewRouter()

	test.Append(router)

	svr := core.NewHTTPServer(":8010")
	svr.Bind(router)

	svr.Use(&test.Hello{})

	//svr.Use(&test.Test{})

	svr.Run()
}
