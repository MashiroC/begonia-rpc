// Time : 2019/8/22 下午3:03
// Author : MashiroC

// client1 something
package main

import (
	"encoding/json"
	"fmt"
	conn2 "github.com/MashiroC/begonia-rpc/conn"
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/util/log"
	"net"
)

type MathService struct {
}

type HelloService struct {
}

func (s *MathService) Sum(a,b int) int {
	return a+b
}

func (h *HelloService) Hello() string {
	return "world"
}

//func (h *HelloService) Fun1(c *begonia.Context) {
//	c.JSON(entity.Param{"hello": "world"})
//}
//
//


func main() {
	//cli := begonia.New("localhost:4949")
	//cli.Sign("Hello", &HelloService{})
	//cli.Sign("Math", &MathService{})
	//cli.KeepConnect()
	//i := 1
	//res := testMarshal(1.0)
	//testUnmarshal(res)
	//fmt.Println()
	//
	//testCall("hc", Person{Name:"hc",Age:18}, 4949.01)
	conn, err := net.Dial("tcp4", ":4949")
	if err!=nil{
		log.Fatal(err.Error())
	}
	c:=conn2.New(conn)
	en:=entity.SignForm{
		Sign:[]entity.SignEntity{
			{
				Name:   "Hello",
				Fun:    []entity.FunEntity{{
					Name: "Hello",
					Size: 1,
				}},
				IsMore: false,
			},
		},
	}
	b, _ := json.Marshal(en)
	c.WriteSign(b)

	opcode,data, err := c.ReadData()
	fmt.Println(opcode)
	fmt.Println(string(data))
	resp:=entity.RespForm{
		Uuid: "11111",
		Data: "hello hc",
	}
	rb, err := json.Marshal(resp)
	if err!=nil{
		log.Fatal(err.Error())
	}
	c.WriteResponse(rb)
	wait:=make(chan bool)
	<-wait
}
