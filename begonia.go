package begonia_rpc

import (
	conn2 "mashiroc.fun/begoniarpc/conn"
	"mashiroc.fun/begoniarpc/util/log"
	"net"
)

type Server struct {
	resp *respHandler
	call *callHandler
}

func Default() *Server {
	return &Server{
		resp: newRespHandler(),
		call: newCallHandler(),
	}
}

func (s *Server) Run(addr string) {
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("conn failed : ", err)
			continue
		}
		c := conn2.NewConn(conn)
		go s.work(c)
	}
}
