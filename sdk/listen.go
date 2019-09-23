package begonia

import (
	"encoding/json"
	"io"
	"mashiroc.fun/begonia/conn"
	"mashiroc.fun/begonia/entity"
	"mashiroc.fun/begonia/util/log"
	"net"
)

// connection 开启连接
func (cli *Client) connection(addr string) conn.Conn {
	c, err := net.Dial("tcp4", addr)
	if err != nil {
		panic(err.Error())
	}

	return conn.New(c)
}

// listen 监听 应该开启连接之后调用
func (cli *Client) listen() {
	for {
		opcode, data, err := cli.conn.ReadData()
		if err != nil {
			// 如果拿到数据这里有错 直接把连接关了
			_ = cli.conn.Close()
			if _, ok := err.(*net.OpError); ok || err == io.EOF {
				log.Info("remote addr [" + cli.conn.Addr() + "]连接断开")
			} else {
				log.Error("remote addr ["+cli.conn.Addr()+"] error:", err)
			}
			break
		}

		// 读到帧之后开一个协程处理这个帧 然后继续读新的
		go cli.operator(cli.conn, opcode, data)
	}
	cli.wait <- true

}

func (cli *Client) operator(c conn.Conn, opcode uint8, data []byte) {
	switch opcode {
	case conn.OpRequest: // 远程调用request 其他cli想要请求这个服务
		cli.handlerRequest(data)
	case conn.OpResponse: // 远程调用resp center返回了其他服务的resp
		cli.handlerResponse(data)
	case conn.OpError: // resp error
		cli.handlerError(data)
	default:
		log.Warn("unknow opcode:", opcode, string(data))
	}
}

func (cli *Client) closeWith(err error) {
	//cli.respError(err)
	_ = cli.conn.Close()
}

func (cli *Client) respError(err error) {
	// 如果这个错误是CallError 那么直接序列化 如果不是 把err的信息写一个新的CallError
	cErr, ok := err.(entity.CallError)
	if !ok {
		cErr = entity.NewError(entity.ErrCodeUnknow, err.Error())
	}
	b, _ := json.Marshal(cErr)
	_ = cli.conn.WriteError(b)
}
