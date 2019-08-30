// Time : 2019/8/22 下午9:30 
// Author : MashiroC

// client2 something
package main

import (
	"fmt"
	"mashiroc.fun/redrpc/sdk"
)

func main() {
	param := sdk.Params{
		"a": "1",
		"b": "1",
	}

	cli := sdk.Default("127.0.0.1:1234")

	resp := cli.Call("math", "Sum", param)
	fmt.Println(resp["sum"])

	cli.CallAsyn("math", "Sum",param, func(params sdk.Params) {
		fmt.Println(params["sum"])
	})
	wait:=make(chan bool)
	<-wait
}

//func main() {
//	cli := sdk.Default("127.0.0.1:1234")
//	size := 100000
//	work := make(chan string, size)
//	for i := 0; i < size; i++ {
//		go func() {
//			resp := cli.Call("math", "Sum", sdk.Params{
//				"a": "1",
//				"b": "1",
//			})
//
//			work <- resp["sum"].(string)
//		}()
//	}
//
//	for i := 0; i < size; i++ {
//		_= <-work
//	}
//}
