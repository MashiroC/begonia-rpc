// Time : 2019/8/22 下午9:30
// Author : MashiroC

// client2 something
package main

import (
	"fmt"
	"mashiroc.fun/begoniarpc/entity"
	"mashiroc.fun/begoniarpc/sdk"
)

func main() {
	param := entity.Param{
		"a": "1",
		"b": "1",
	}

	cli := sdk.Default("127.0.0.1:1234")
	size := 100000
	for i := 0; i < size; i++ {
		go func() {
			resp := cli.Call("hello", "Hello", param)
			if resp == nil || resp["hello"] != "hello" {
				fmt.Println(resp)
			}
		}()
	}

	wait := make(chan bool)
	<-wait
}
