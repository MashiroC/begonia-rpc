// Time : 2019/8/23 下午2:38 
// Author : MashiroC

// sdk something
package sdk

import (
	"bufio"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"mashiroc.fun/redrpc/center"
	"mashiroc.fun/redrpc/util"
	"net"
	"reflect"
	"sync"
)

type RedRpcClient struct {
	serviceMap map[string]ServiceEntity
	conn       net.Conn
	buf        *bufio.ReadWriter
	callMap    map[string]chan RespEntity
	lock       sync.Mutex
	callLock   sync.Mutex
}

type SignEntity struct {
	Service []ServiceEntity `json:"service"`
}

type ServiceEntity struct {
	service  interface{}
	Name     string   `json:"name"`
	Function []string `json:"function"`
}

type CallEntity struct {
	Uuid     string `json:"1"`
	Service  string `json:"2"`
	Function string `json:"3"`
	Param    Params `json:"4"`
}

type RespEntity struct {
	Uuid string `json:"1"`
	Data Params `json:"2"`
}

func Default(centerAddr string) (cli *RedRpcClient) {
	cli = &RedRpcClient{
		serviceMap: make(map[string]ServiceEntity, 10),
		callMap:    make(map[string]chan RespEntity, 10),
	}
	cli.server(centerAddr)
	return
}

func (cli *RedRpcClient) Test() {
	fmt.Println(cli.callMap)
}

func (cli *RedRpcClient) Call(service, function string, param Params) (res Params) {
	u1 := uuid.NewV4()
	e := CallEntity{
		Uuid:     u1.String(),
		Service:  service,
		Function: function,
		Param:    param,
	}
	data, _ := json.Marshal(e)
	util.Send(cli.buf, 2, data)

	ch := make(chan RespEntity, 1)
	cli.callLock.Lock()
	cli.callMap[e.Uuid] = ch
	cli.callLock.Unlock()
	resp := <-ch
	res = resp.Data
	return
}

func (cli *RedRpcClient) CallAsyn(service, fun string, param Params, callback func(Params)) {
	u1 := uuid.NewV4()
	e := CallEntity{
		Uuid:     u1.String(),
		Service:  service,
		Function: fun,
		Param:    param,
	}
	data, _ := json.Marshal(e)
	util.Send(cli.buf, 2, data)

	ch := make(chan RespEntity, 1)
	cli.callLock.Lock()
	cli.callMap[e.Uuid] = ch
	cli.callLock.Unlock()
	go func() {
		resp := <-ch
		callback(resp.Data)
	}()
}

func (cli *RedRpcClient) Sign(name string, service interface{}) {

	t := reflect.TypeOf(service)
	server := ServiceEntity{
		service:  service,
		Name:     name,
		Function: []string{},
	}
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if m.Type.NumOut() != 1 || m.Type.Out(0) != reflect.TypeOf(Params{}) {
			log.Fatal("return must Params")
		}
		if m.Type.NumIn() != 2 || m.Type.In(1) != reflect.TypeOf(Params{}) {
			log.Fatal("input must Params")
		}
		// 检查好了

		server.Function = append(server.Function, m.Name)

	}
	cli.lock.Lock()
	cli.serviceMap[name] = server
	cli.lock.Unlock()
	e := SignEntity{Service: []ServiceEntity{server}}

	b, _ := json.Marshal(e)

	util.Send(cli.buf, 1, b)

}

func (cli *RedRpcClient) server(centerAddr string) {
	conn, err := net.Dial("tcp", centerAddr)

	if err != nil {
		log.Fatal(err.Error())
	}

	buf := center.CreateBuf(conn)
	cli.conn = conn
	cli.buf = buf
	go func() {
		lock := sync.Mutex{}
		i := 0
		for {
			lock.Lock()
			opcode, data, err := center.ReadData(conn)
			lock.Unlock()
			i++
			//fmt.Println(i)
			if err != nil {
				break
			}
			go cli.operator(opcode, data)

		}
	}()

}

func (cli *RedRpcClient) operator(opcode uint8, data []byte) {
	switch opcode {
	case 2:
		cli.handlerRequest(data)
	case 3:
		cli.handlerResponse(data)
	case 4:
		cli.handerError(data)
	default:
		fmt.Println(opcode, string(data))
	}
}

func (cli *RedRpcClient) handlerRequest(data []byte) {
	e := CallEntity{}
	err := json.Unmarshal(data, &e)
	if err != nil {
		return
	}

	server, ok := cli.serviceMap[e.Service]
	if !ok {
		r := RespEntity{
			Uuid: e.Uuid,
			Data: Params{
				"errorCode":    "404",
				"errorMessage": "service not found",
			},
		}
		b, _ := json.Marshal(r)
		util.SendError(cli.buf, b)
	}
	for _, f := range server.Function {
		if f == e.Function {
			cli.request(server, e)
			break
		}
	}
}

func (cli *RedRpcClient) handlerResponse(data []byte) {
	resp := RespEntity{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		log.Println("json error:", err.Error(), string(data))
		return
	}
	cli.callLock.Lock()
	ch, ok := cli.callMap[resp.Uuid]
	cli.callLock.Unlock()
	if !ok {
		log.Println("uuid not found")
		return
	}
	cli.callLock.Lock()
	delete(cli.callMap, resp.Uuid)
	cli.callLock.Unlock()

	ch <- resp
}

func (cli *RedRpcClient) request(server ServiceEntity, e CallEntity) {
	v := reflect.ValueOf(server.service)
	attr := []reflect.Value{reflect.ValueOf(e.Param)}
	res := v.MethodByName(e.Function).Call(attr)
	p := res[0].Interface().(Params)
	if res != nil {
		data, _ := json.Marshal(RespEntity{
			Uuid: e.Uuid,
			Data: p,
		})
		util.Send(cli.buf, 3, data)
	}
}

func (cli *RedRpcClient) handerError(bytes []byte) {
	fmt.Println(string(bytes))
}

type Params map[string]interface{}
