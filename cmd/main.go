// Time : 2019/8/22 下午2:58
// Author : MashiroC

// redrpc something
package main

import (
	"mashiroc.fun/begoniarpc"
)

func main() {
	//app:= demo1.Server("127.0.0.1:1234")
	//defer app.Close()
	rpc := begonia_rpc.Default()
	rpc.Run("localhost:1234")
	wait := make(chan bool)
	<-wait
}
