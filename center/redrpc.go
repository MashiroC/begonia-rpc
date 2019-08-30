// Time : 2019/8/22 下午4:32 
// Author : MashiroC

// redrpc something
package center

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type RedRpcServer struct {
	callMap   map[string]chan RespEntity
	listener  net.Listener
	services  map[string]ServiceEntity
	lock      sync.Mutex
	callLock  sync.Mutex
	writeLock sync.Mutex
}

type ServiceEntity struct {
	conn     net.Conn
	buf      *bufio.ReadWriter
	addr     string
	Name     string   `json:"name"`
	Function []string `json:"function"`
}

type Params map[string]interface{}

func Server(addr string) (s *RedRpcServer) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	service := make(map[string]ServiceEntity, 5)
	call := make(map[string]chan RespEntity, 5)
	s = &RedRpcServer{
		listener: listen,
		services: service,
		callMap:  call,
	}
	//go func() {
	//	for {
	//		time.Sleep(5 * time.Second)
	//		fmt.Println(len(s.callMap))
	//	}
	//}()
	go s.Run()
	return
}

func (s *RedRpcServer) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		buf := CreateBuf(conn)
		go s.work(conn, buf)
	}
}

func CreateBuf(conn net.Conn) *bufio.ReadWriter {
	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)
	buf := bufio.NewReadWriter(r, w)
	return buf
}

func (s *RedRpcServer) work(conn net.Conn, buf *bufio.ReadWriter) {

	for {

		opcode, data, err := ReadData(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println(conn.RemoteAddr().String(), "连接断开")
			} else {
				fmt.Println(conn.RemoteAddr().String(), "error:", err.Error())
			}
			s.onError(conn)
			return
		}

		go s.operator(conn, buf, opcode, data)

	}
}

func (s *RedRpcServer) onError(conn net.Conn) {
	keys := []string{}
	for k, v := range s.services {
		if v.addr == conn.RemoteAddr().String() {
			keys = append(keys, k)
		}
	}
	for _, key := range keys {
		delete(s.services, key)
	}
}

func (s *RedRpcServer) operator(conn net.Conn, buf *bufio.ReadWriter, opcode uint8, data []byte) {

	switch opcode {
	case 1: // 注册服务
		s.signServer(conn, buf, data)
	case 2: // 远程调用request
		s.handlerRequest(conn, buf, data)
	case 3: // 远程调用resp
		s.handlerResponse(data)
	case 4: // resp error
		s.handlerError(conn, buf, data)
	}
}

func (s *RedRpcServer) Close() {
	s.listener.Close()
}

func ReadData(conn net.Conn) (opcode uint8, data []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	len1 := make([]byte, 1)
	_, err = conn.Read(len1)
	//if l != 1 {
	//	log.Fatal("fuck len1")
	//}
	//len1, err := buf.ReadByte()
	handlerErr(err)

	// 拿payload长度
	payloadLen := uint(len1[0])

	// len1如果是255 则去看扩展len
	if len1[0] == byte(255) {
		fmt.Println("ffffff")
		overLen := make([]byte, 2)
		_, err = conn.Read(overLen)
		handlerErr(err)
		payloadLen = uint(binary.BigEndian.Uint16(overLen))
	}

	tmp := make([]byte, 1)
	opl, err := conn.Read(tmp)
	opcode = tmp[0]
	if opl != 1 {
		log.Fatal("fuck opcode len")
	}
	//opcode, err = buf.ReadByte()
	handlerErr(err)

	// 拿数据
	data = make([]byte, payloadLen)
	n, err := conn.Read(data)
	//n, err := buf.Read(data)
	handlerErr(err)
	for uint(n) != payloadLen {
		fmt.Println("over")
		//fmt.Println(payloadLen,n,string())
		oversize := make([]byte, payloadLen-uint(n))
		size, err := conn.Read(oversize)
		for i:=n;i<int(payloadLen);i++{
			data[i]=oversize[i-n]
		}
		n += size
		handlerErr(err)
		//log.Fatal(string(data))
	}
	//fmt.Println("????fuck")

	return
}

func handlerErr(err error) {
	if err != nil {
		panic(err)
	}
}
