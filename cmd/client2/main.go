// Time : 2019/8/22 下午9:30
// Author : MashiroC

// client2 something
package main

import (
	"fmt"
	"mashiroc.fun/begonia/entity"
	begonia "mashiroc.fun/begonia/sdk"
)

func main() {
	cli := begonia.New("localhost:1234")
	req := begonia.Request{
		Service:  "Math",
		Function: "Sum",
		Param:    entity.Param{"a": 1, "b": 1},
	}
	resp := cli.Call(req)
	fmt.Println("call")

	i := resp.Int()
	if resp.Error() != nil {
		fmt.Println(resp.Error())
	} else {
		fmt.Println(i)
	}

	cli.CallAsync(req, func(resp *begonia.Response) {
		fmt.Println("async")

		i := resp.Int()
		if resp.Error() != nil {
			fmt.Println(resp.Error())
		} else {
			fmt.Println(i)
		}
	})
	wait:=make(chan bool)
	<-wait
	//cli.CallAsync(resp.S)
}
