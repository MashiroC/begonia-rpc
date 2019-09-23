// Time : 2019/8/22 下午3:03
// Author : MashiroC

// client1 something
package main

import (
	"fmt"
	"mashiroc.fun/begonia/sdk"
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

func (h *HelloService) Test(name string) {
	fmt.Println("reflect ok")
}

func (h *HelloService) Hello(c *begonia.Context) {
	c.Int(1234)
}

func main() {
	cli := begonia.New("localhost:1234")
	cli.Sign("Hello", &HelloService{})
	cli.Sign("Math",&MathService{})
	cli.Wait()
}
