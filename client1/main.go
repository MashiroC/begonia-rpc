// Time : 2019/8/22 下午3:03 
// Author : MashiroC

// client1 something
package main

import (
	"mashiroc.fun/redrpc/sdk"
	"strconv"
)

type MathService struct {
}

type HelloService struct {
}

func (s *MathService) Sum(params sdk.Params) (res sdk.Params) {
	res = make(sdk.Params,1)
	a := params["a"].(string)
	b := params["b"].(string)
	pa, _ := strconv.ParseInt(a, 10, 0)
	pb, _ := strconv.ParseInt(b, 10, 0)
	res["sum"] = strconv.FormatInt(pa+pb,10)
	return
}

func (s *MathService) Test(params sdk.Params) (res sdk.Params) {
	res["hello"] = "world"
	return
}

func (h *HelloService) Hello(params sdk.Params) (res sdk.Params) {
	res["errorCode"] = "1551"
	res["errorMessage"] = "hello"
	return
}

func main() {
	cli := sdk.Default("127.0.0.1:1234")

	cli.Sign("math", &MathService{})
	cli.Sign("hello", &HelloService{})


	wait := make(chan bool)
	<-wait
}
