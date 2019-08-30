// Time : 2019/8/22 下午5:35 
// Author : MashiroC

// center something
package center

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)

type SignEntity struct {
	Service []ServiceEntity `json:"service"`
}

func (s *RedRpcServer) signServer(conn net.Conn, buf *bufio.ReadWriter, data []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()
	e := SignEntity{}
	err := json.Unmarshal(data, &e)
	if err != nil {
		log.Println("json error :", string(data))
		return
	}
	for _, service := range e.Service {
		service.conn = conn
		service.buf = buf
		service.addr = conn.RemoteAddr().String()
		s.services[service.Name] = service
	}
	log.Println("服务注册:", string(data))
}
