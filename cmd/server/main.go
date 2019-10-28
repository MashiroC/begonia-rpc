// Time : 2019/8/22 下午2:58
// Author : MashiroC

// redrpc something
package main

import (
	"github.com/MashiroC/begonia-rpc"
)

func main() {
	rpc := begoniarpc.Default()
	rpc.Run("0.0.0.0:8080")
}
