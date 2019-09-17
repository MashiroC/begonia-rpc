// Time : 2019/8/22 下午3:03
// Author : MashiroC

// client1 something
package main

import (
	"fmt"
	"mashiroc.fun/begonia/entity"
	"mashiroc.fun/begonia/sdk"
	"strconv"
)

type MathService struct {
}

type HelloService struct {
}

func (s *MathService) Sum(params entity.Param) (res entity.Param) {
	res = make(entity.Param, 1)
	a := params["a"].(string)
	b := params["b"].(string)
	pa, _ := strconv.ParseInt(a, 10, 0)
	pb, _ := strconv.ParseInt(b, 10, 0)
	res["sum"] = strconv.FormatInt(pa+pb, 10)
	return
}

func (s *MathService) Test(params entity.Param) (res entity.Param) {
	res["hello"] = "world"
	return
}

func (h *HelloService) Hello(params entity.Param) (res entity.Param) {
	res = make(entity.Param, 2)
	fmt.Println("fuck")
	//res["errorCode"] = "1551"
	res["hello"] = "hello"
	return
}

func main() {
	cli := sdk.Default("127.0.0.1:1234")

	cli.Sign("hello", &HelloService{})

	//cli.CallAsyn("math", "Sum",param, func(params entity.Param) {
	//	fmt.Println(params["sum"])
	//})
	wait := make(chan bool)
	<-wait
}

type Ttt struct {
}
