package main

import (
	"github.com/hsedjame/webfast/web"
)

type TestReq struct {
	Name string `validate:"required" json:"name"`
}

func (t TestReq) Validate(fn web.ValidationFn) error {
	return fn(t)
}

func main() {

	//var t TestReq
	//
	//check := web.CheckVal(t)
	//
	//fmt.Println(check)
}
