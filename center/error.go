// Time : 2019/8/23 下午2:18 
// Author : MashiroC

// center something
package center

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func (s *RedRpcServer) handlerError(conn net.Conn, buf *bufio.ReadWriter, data []byte) {
	e := RespEntity{}
	err := json.Unmarshal(data, &e)
	if err != nil {
		log.Println(err.Error())
		return
	}
	errCode, ok1 := e.Data["errorCode"]
	errMsg, ok2 := e.Data["errorMessage"]

	if !ok1 || !ok2 {
		log.Println("not found error param")
		return
	}

	fmt.Println(conn.RemoteAddr().String(), errCode, errMsg)
	ch, ok := s.callMap[e.Uuid]
	if !ok {
		log.Println("error uuid not found:", e.Uuid)
		return
	}

	ch <- e
}
