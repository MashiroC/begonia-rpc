package begonia_rpc

import (
	"io"
	"mashiroc.fun/begoniarpc/conn"
	"mashiroc.fun/begoniarpc/util/log"
	"net"
)

func (s *Server) work(conn conn.Conn) {
	for {
		opcode, data, err := conn.ReadData()
		if err != nil {
			_ = conn.Close()
			_ = s.call.unBindService(conn)
			if _, ok := err.(*net.OpError); ok || err == io.EOF {
				log.Info("remote addr [" + conn.Addr() + "]连接断开")
			} else {
				log.Error("remote addr ["+conn.Addr()+"] error:", err)
			}
			break
		}

		go s.operator(conn, opcode, data)
	}
}

func (s *Server) operator(conn conn.Conn, opcode uint8, data []byte) {
	switch opcode {
	case 1: // 注册服务
		s.handlerSign(conn, data)
	case 2: // 远程调用request
		s.handlerRequest(conn, data)
	case 3: // 远程调用resp
		s.handlerResponse(conn, data)
	case 4: // resp error
		s.handlerError(conn, data)
	default:
		log.Warn("unknow opcode:", opcode, string(data))
	}
}
