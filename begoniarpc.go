// sdk-rpc 轻量级rpc框架
// 架构分为客户端 服务端 服务中心三部分
// By MashiroC
package begoniarpc

import (
	"encoding/json"
	begoniaConn "github.com/MashiroC/begonia-rpc/conn"
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/util/log"
	"net"
)

// ServerCenter 服务中心的实体
type ServerCenter struct {
	resp *respHandler
	call *callHandler
}

// Default 返回一个默认配置的服务中心
func Default() *ServerCenter {
	return &ServerCenter{
		resp: newRespHandler(),
		call: newCallHandler(),
	}
}

// Run 开始监听
func (s *ServerCenter) Run(addr string) {
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		panic(err)
	}
	log.Info("RPC center start on %s", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("conn failed : %s", err)
			continue
		}
		c := begoniaConn.New(conn)
		go s.work(c)
	}
}

// respError 向某条连接写一个异常
func respError(conn begoniaConn.Conn, uuid string, cErr entity.CallError) {
	log.Error("remote [%s] frame has some error: %s", conn.Addr(), cErr.Error())
	eForm := entity.ErrForm{
		Uuid:    uuid,
		ErrCode: cErr.ErrCode,
		ErrMsg:  cErr.ErrMessage,
	}
	b, _ := json.Marshal(eForm)

	_ = conn.WriteError(b)
}
