// Time : 2019/8/22 下午5:35 
// Author : MashiroC

// center something
package center

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"mashiroc.fun/redrpc/util"
	"net"
)

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

func (s *RedRpcServer) handlerRequest(client net.Conn, buf *bufio.ReadWriter, data []byte) {
	e := CallEntity{}
	err := json.Unmarshal(data, &e)
	if err != nil {
		log.Println("json error :", string(data))
		return
	}
	service, ok := s.services[e.Service]
	if !ok {
		resp, _ := json.Marshal(RespEntity{
			Uuid: e.Uuid,
			Data: Params{
				"errorCode":    "404",
				"errorMessage": "service not found",
			},
		})
		util.SendError(buf, resp)
		return
	}

	ok = false
	for _, fun := range service.Function {
		if fun == e.Function {
			ok = true
		}
	}

	if !ok {
		log.Println("call error, not found function:", e.Function)
		return
	}

	//log.Println("rpc:", string(data))
	s.call(buf, service.buf, e)
}

var reqCount = 0
var respCount = 0

func (s *RedRpcServer) handlerResponse(data []byte) {
	e := RespEntity{}
	err := json.Unmarshal(data, &e)
	if err != nil {
		fmt.Println("resp entity json error:", string(data))
		return
	}
	s.callLock.Lock()
	ch, ok := s.callMap[e.Uuid]
	s.callLock.Unlock()

	if !ok {
		fmt.Println("resp uuid not found")
		return
	}
	//respCount++
	//go func() {
	//	for {
	//time.Sleep(3 * time.Second)
	//fmt.Println("respCount", respCount)
	//}
	//}()
	//fmt.Println("respCount", respCount)
	ch <- e
}

func (s *RedRpcServer) call(client *bufio.ReadWriter, service *bufio.ReadWriter, entity CallEntity) {
	callback := make(chan RespEntity, 1)

	s.callLock.Lock()
	s.callMap[entity.Uuid] = callback
	s.callLock.Unlock()

	go func() {
		resp := <-callback
		s.callLock.Lock()
		delete(s.callMap, resp.Uuid)
		s.callLock.Unlock()

		//log.Println("rpc resp:", resp)
		b, _ := json.Marshal(resp)
		util.Send(client, 3, b)
	}()
	payload, _ := json.Marshal(entity)

	opcode:=byte(2)

	s.writeLock.Lock()
	defer s.writeLock.Unlock()
	payloadSize := len(payload)
	//fmt.Println("sendLen:",payloadSize)
	//if payloadSize==60{
	//	fmt.Println(string(payload))
	//}
	var header []byte
	if payloadSize > 254 && payloadSize < 65536 {
		// 右移8位取取前八位bytes
		// a / 2**n = a >> n
		len1 := byte(payloadSize >> 8)
		// 取后八位
		// a % b = a & (b-1)
		len2 := byte(payloadSize * 7)
		header = []byte{255, len1, len2, opcode}
	} else if payloadSize <=254{
		//fmt.Println(payloadSize)
		header = []byte{byte(payloadSize), opcode}
	}else{
		fmt.Println("fuck?")
		return
	}

	data := append(header, payload...)

	_, err := service.Write(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = service.Flush()
	if err != nil {
		fmt.Println(err.Error())
	}

}
