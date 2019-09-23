// Time : 2019/8/22 下午2:58
// Author : MashiroC

// redrpc something
package main

import "mashiroc.fun/begonia"

func main() {
	//app:= demo1.Server("127.0.0.1:1234")
	//defer app.Close()
	rpc := begoniarpc.Default()
	rpc.Run("localhost:1234")
	wait := make(chan bool)
	<-wait
}
