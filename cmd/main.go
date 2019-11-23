// Time : 2019/10/28 14:36
// Author : MashiroC

// cmd
package main

import begoniarpc "github.com/MashiroC/begonia-rpc"

// main.go something

func main(){
	center :=begoniarpc.Default()
	center.Run(":12306")
}
