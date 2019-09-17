// Time : 2019/8/23 下午2:38
// Author : MashiroC

// sdk something
package sdk

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	begoniaConn "mashiroc.fun/begonia/conn"
	"mashiroc.fun/begonia/entity"
	"mashiroc.fun/begonia/util/log"
	"net"
	"reflect"
	"sync"
)

type RedRpcClient struct {
	serviceMap map[string]ServiceEntity
	conn       begoniaConn.Conn
	callMap    map[string]chan RespEntity
	lock       sync.Mutex
	callLock   sync.Mutex
}

type SignEntity struct {
	Service []ServiceEntity `json:"1"`
}

type ServiceEntity struct {
	service  interface{}
	Name     string   `json:"1"`
	Function []string `json:"2"`
}

type CallEntity struct {
	UUID     string       `json:"1"`
	Service  string       `json:"2"`
	Function string       `json:"3"`
	Param    entity.Param `json:"4"`
}

type RespEntity struct {
	UUID string       `json:"1"`
	Data entity.Param `json:"2"`
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

func (cli *RedRpcClient) Call(service, function string, param entity.Param) (res entity.Param) {
	u1 := uuid.NewV4()
	e := CallEntity{
		UUID:     u1.String(),
		Service:  service,
		Function: function,
		Param:    param,
	}
	data, _ := json.Marshal(e)
	_ = cli.conn.WriteRequest(data)
	ch := make(chan RespEntity, 1)
	cli.callLock.Lock()
	cli.callMap[e.UUID] = ch
	cli.callLock.Unlock()
	resp := <-ch
	res = resp.Data
	return
}

func (cli *RedRpcClient) CallAsyn(service, fun string, param entity.Param, callback func(entity.Param)) {
	u1 := uuid.NewV4()
	e := CallEntity{
		UUID:     u1.String(),
		Service:  service,
		Function: fun,
		Param:    param,
	}
	data, _ := json.Marshal(e)

	cli.conn.WriteRequest(data)

	ch := make(chan RespEntity, 1)
	cli.callLock.Lock()
	cli.callMap[e.UUID] = ch
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
		if m.Type.NumOut() != 1 || m.Type.Out(0) != reflect.TypeOf(entity.Param{}) {
			log.Fatal("return must entity.Param")
		}
		if m.Type.NumIn() != 2 || m.Type.In(1) != reflect.TypeOf(entity.Param{}) {
			log.Fatal("input must entity.Param")
		}
		// 检查好了

		server.Function = append(server.Function, m.Name)

	}
	cli.lock.Lock()
	cli.serviceMap[name] = server
	cli.lock.Unlock()
	e := SignEntity{Service: []ServiceEntity{server}}

	b, _ := json.Marshal(e)

	cli.conn.WriteSign(b)

}

func (cli *RedRpcClient) server(centerAddr string) {
	conn, err := net.Dial("tcp", centerAddr)

	if err != nil {
		log.Fatal(err.Error())
	}

	c := begoniaConn.NewConn(conn)
	cli.conn = c
	go func() {
		for {
			opcode, data, err := c.ReadData()
			if err != nil {
				log.Error("readData error: %s", err.Error())
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
	//fmt.Println(cli.serviceMap)
	if !ok {
		r := RespEntity{
			UUID: e.UUID,
			Data: entity.Param{
				"errorCode":    "404",
				"errorMessage": "server not found",
			},
		}
		b, _ := json.Marshal(r)
		cli.conn.WriteError(b)
	}

	for _, f := range server.Function {
		if f == e.Function {
			cli.request(server, e)
			break
		}
	}
}

func (cli *RedRpcClient) handlerResponse(data []byte) {
	//fmt.Println("handler resp", string(data))
	resp := RespEntity{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		log.Error("json error: %s for data: %s", err.Error(), string(data))
		return
	}
	//fmt.Println(cli.callMap)
	//fmt.Println(resp.UUID)
	cli.callLock.Lock()
	ch, ok := cli.callMap[resp.UUID]
	cli.callLock.Unlock()
	if !ok {
		log.Error("uuid [%s] not found", resp.UUID)
		return
	}
	cli.callLock.Lock()
	delete(cli.callMap, resp.UUID)
	cli.callLock.Unlock()

	ch <- resp
}

func (cli *RedRpcClient) request(server ServiceEntity, e CallEntity) {
	v := reflect.ValueOf(server.service)
	if e.Param == nil {
		//fmt.Println("fuck")
		e.Param = make(entity.Param, 5)
	}
	attr := []reflect.Value{reflect.ValueOf(e.Param)}
	res := v.MethodByName(e.Function).Call(attr)
	p := res[0].Interface().(entity.Param)
	if res != nil {
		data, _ := json.Marshal(RespEntity{
			UUID: e.UUID,
			Data: p,
		})
		//fmt.Println(string(data))
		cli.conn.WriteResponse(data)
	}
}

func (cli *RedRpcClient) handerError(bytes []byte) {
	fmt.Println("handler error", string(bytes))
}
