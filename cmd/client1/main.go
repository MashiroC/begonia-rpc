// Time : 2019/8/22 下午3:03
// Author : MashiroC

// client1 something
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

func (h *HelloService) Fun1(c *begonia.Context) {
	c.JSON(entity.Param{"hello": "world"})
}

func main() {
	cli := begonia.New("localhost:1234")
	cli.Sign("Hello", &HelloService{})
	cli.Sign("Math", &MathService{})
	cli.Wait()
}
