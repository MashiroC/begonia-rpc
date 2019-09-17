package begonia

// listen.go 监听和初步解析帧

import (
	"io"
	"mashiroc.fun/begonia/conn"
	"mashiroc.fun/begonia/util/log"
	"net"
)

// work 监听并循环从连接拿到一个帧
func (s *ServerCenter) work(conn conn.Conn) {
	for {
		opcode, data, err := conn.ReadData()
		if err != nil {
			// 如果拿到数据这里有错 直接把连接关了 然后解绑所有和这条连接有关的服务
			_ = conn.Close()
			_ = s.call.unBindService(conn)
			if _, ok := err.(*net.OpError); ok || err == io.EOF {
				log.Info("remote addr [" + conn.Addr() + "]连接断开")
			} else {
				log.Error("remote addr ["+conn.Addr()+"] error:", err)
			}
			break
		}

		// 读到帧之后开一个协程处理这个帧 然后继续读新的
		go s.operator(conn, opcode, data)
	}
}

// operator 根据操作码分发到各个函数 操作吗见conn包
func (s *ServerCenter) operator(c conn.Conn, opcode uint8, data []byte) {
	switch opcode {
	case conn.OpSign: // 注册服务
		s.handlerSign(c, data)
	case conn.OpRequest: // 远程调用request
		s.handlerRequest(c, data)
	case conn.OpResponse: // 远程调用resp
		s.handlerResponse(c, data)
	case conn.OpError: // resp error
		s.handlerError(c, data)
	default:
		log.Warn("unknow opcode:", opcode, string(data))
	}
}
