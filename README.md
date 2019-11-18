# Begonia-RPC

A golang rpc framework for efficient and concise.

## User Guide

### Prerequisites

Golang Version >= 1.11.2

### Installation

```bash
$ go get github.com/MashiroC/begonia-rpc
```

### Example

**service center**

```go
package main

import (
	"github.com/MashiroC/begonia-sdk"
)

func main() {
	rpc := begoniarpc.Default()
	rpc.Run("0.0.0.0:8080")
}
```

How to provide a service?

```go
package main

import (
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/sdk"
)

type MathService struct {
}

type HelloService struct {
}

func (s *MathService) Sum(a, b int) (res int) {
    return a + b
}


func (h *HelloService) Hello(name string) (res string) {
	return "Hello " + name
}

func main() {
	cli := begonia.Default(":8080")
	cli.Sign("Hello", &HelloService{})
	cli.Sign("Math",&MathService{})
	cli.KeepConnect()
}
```

And How to call remote function?

```go
package main

import (
	"fmt"
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/sdk"
)

func main() {
    
    // get a begonia client
	cli := begonia.Default(":8080")
	
    // get a service
    helloService := cli.Service("Hello")
    
    // get a sync Function
    hello := helloService.FunSync("Hello")
    
    // call it!
    res, err := hello("MashiroC")
    
    fmt.Println(res, err)
    // Hello Mashiroc <nil>
    
   
    // get a async function
    helloAsync := helloService.FunAsync("Hello")

    // call it too!
    helloAsync(func(res interface{}, err error) {
        fmt.Println(res,err)
    	}, "MashiroC")
    
}
```

### Run

> this is rpc center, if you want to provide a service, please use sdk.

```bash
$ go build -o rpccenter ./cmd/server/main.go
$ ./rpccenter
```

or

```bash
$ go run ./cmd/server/main.go
```

### Release History

0.1.0

### Proof-of-concept code