package conn

// conn_default.go 默认的连接

import (
	"bufio"
	"encoding/binary"
	"errors"
	"net"
	"sync"
)

// DefaultConn 默认的不带连接池的连接
type DefaultConn struct {
	conn net.Conn          // 连接
	buf  *bufio.ReadWriter // 带缓冲的流
	lock sync.Mutex        // 写的排他锁
}

// Addr 拿到连接的远程地址
func (c *DefaultConn) Addr() (res string) {
	return c.conn.RemoteAddr().String()
}

// read 读一定长度的数据 读不到会阻塞
// 设计成这样是因为之前改一个高并发的bug的时候 发现了一个问题
// 在高并发场景下 client发了一个json 这边一次性read不完
// 报序列化错误 导致后面所有的包都乱序
func (c *DefaultConn) read(len uint) (data []byte, err error) {
	data = make([]byte, len)
	n, err := c.buf.Read(data)

	if err != nil {
		return
	}

	// 一次没读够指定的len 继续读
	for n != int(len) {
		overSize := make([]byte, int(len)-n)
		size, err := c.buf.Read(overSize)
		handlerErr(err)
		for i := 0; i < size; i++ {
			data[n+i] = overSize[i]
		}
		n += size
	}

	return
}

// readByte 读一个byte
func (c *DefaultConn) readByte() (data byte, err error) {
	data, err = c.buf.ReadByte()
	return
}

// ReadData 读一个包 返回opcode和json数据
func (c *DefaultConn) ReadData() (opcode uint8, data []byte, err error) {

	// 检测一个有没有panic出来错误 有的话把连接关了
	defer func() {
		if err := recover(); err != nil {
			c.Close()
			return
		}
	}()

	// 读第一个长度 是包的头 不需要等超时 直接等第一个包来就行
	// 除了第一个包 剩下的都要等超时 TODO:这个应该不需要超时
	baseLen, err := c.buf.ReadByte()
	handlerErr(err)

	// 拿payload长度
	payloadLen := uint(baseLen)

	// baseLen如果是255 则去看扩展len
	if baseLen == BaseLenMaxByte {
		// 这里读了两个byte 然后转化成int
		extendLen, err := c.read(2)
		handlerErr(err)
		payloadLen = uint(binary.BigEndian.Uint16(extendLen))
		// 我们不支持超过一定大小的包
		if payloadLen >= ExtendLenMax {
			err = errors.New("payload len oversize")
			return 0, nil, err
		}
	}

	// 拿opcode
	opcode, err = c.readByte()
	handlerErr(err)

	// 拿数据
	data, err = c.read(payloadLen)
	//n, err := buf.Read(data)
	handlerErr(err)

	return
}

// WriteData 写数据进去
func (c *DefaultConn) WriteData(opcode int, data []byte) (err error) {

	// 因为bufio的Writer不是并发安全的 所以这里要加个锁
	c.lock.Lock()
	defer c.lock.Unlock()

	payloadLenSize := 1
	size := make([]byte, 0)
	if len(data) >= BaseLenMax {
		binary.BigEndian.PutUint16(size, uint16(len(data)))
		payloadLenSize = 2
	} else {
		size = append(size, byte(len(data)))
	}
	if len(data) >= ExtendLenMax {
		err = errors.New("payload len oversize")
		return
	}
	tmp := make([]byte, 0, payloadLenSize+1+len(data))
	tmp = append(tmp, size...)
	tmp = append(tmp, byte(opcode))
	tmp = append(tmp, data...)
	_, err = c.buf.Write(tmp)
	if err != nil {
		return
	}
	err = c.buf.Flush()
	return
}

// WriteSign 写一个注册帧
func (c *DefaultConn) WriteSign(data []byte) (err error) {
	return c.WriteData(OpSign, data)
}

// WriteRequest 写一个请求帧
func (c *DefaultConn) WriteRequest(data []byte) (err error) {
	return c.WriteData(OpRequest, data)
}

// WriteResponse 写一个响应帧
func (c *DefaultConn) WriteResponse(data []byte) (err error) {
	return c.WriteData(OpResponse, data)
}

// WriteError 写一个错误帧
func (c *DefaultConn) WriteError(data []byte) (err error) {
	return c.WriteData(OpError, data)
}

// Close 关闭连接
func (c *DefaultConn) Close() error {
	return c.conn.Close()
}

func New(conn net.Conn) (res Conn) {
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	buf := bufio.NewReadWriter(r, w)
	res = &DefaultConn{conn: conn, buf: buf}
	return
}
