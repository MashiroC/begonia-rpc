// Time : 2019/8/22 下午2:58 
// Author : MashiroC

// redrpc something
package main

import (
	"mashiroc.fun/redrpc/center"
)

func main() {
	app:= center.Server("127.0.0.1:1234")
	defer app.Close()
	wait:=make(chan bool)
	<-wait
}
