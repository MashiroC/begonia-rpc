# Begonia-RPC

A golang rpc framework for efficient and concise.

## User Guide

### Prerequisites

Golang Version >= 1.11.2

### Installation

```bash
$ go get github.com/MashiroC/begonia
```

### Example

**service center**

```go
package main

import (
	"github.com/MashiroC/begonia-rpc"
)

func main() {
	rpc := begoniarpc.Default()
	rpc.Run("0.0.0.0:1234")
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

func (s *MathService) Sum(c *begonia.Context) {
	a := c.Param["a"].(float64)
	b := c.Param["b"].(float64)
	c.Int(int(a + b))
}


func (h *HelloService) Hello(c *begonia.Context) {
	c.String("World!")
}

func (h *HelloService) Fun1(c *begonia.Context){
	c.JSON(entity.Param{"hello":"world"})
}

func main() {
	cli := begonia.New("localhost:1234")
	cli.Sign("Hello", &HelloService{})
	cli.Sign("Math",&MathService{})
	cli.Wait()
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
	cli := begonia.New("localhost:1234")

	req := begonia.Request{
		Service:  "Math",
		Function: "Sum",
		Param:    entity.Param{"a": 1, "b": 1},
	}

	resp := cli.Call(req)
	i := resp.Int()

	if resp.Error() != nil {
		fmt.Println(resp.Error())
	} else {
		fmt.Println(i)
	}

	cli.CallAsync(req, func(resp *begonia.Response) {
		i := resp.Int()
		if resp.Error() != nil {
			fmt.Println(resp.Error())
		} else {
			fmt.Println(i)
		}
	})


	wait := make(chan bool)
	<-wait
}
```

### Run

```bash
$ go build xxxx
$ ./xxxx
```

### Release History

0.1.0

### Proof-of-concept code