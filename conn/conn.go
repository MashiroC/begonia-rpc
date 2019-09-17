// conn 封装过的连接 可以直接读出帧
//
// 帧分成三个部分 分别是len、opcode、payload
//
// len是payload的长度 有baseLen和extendLen两种
// 如果len小于255 占一个byte
// 如果length的长度大于等于255 先1byte的255 然后用2byte来表示长度 最大2byte
//
// opcode占1个byte 具体有哪些看下面的常量
//
// 不管payload是什么类型的数据 这里只处理读写数据
package conn

const (
	// 第一个length的最大值
	BaseLenMax = 255
	// 上述变量的byte
	BaseLenMaxByte = byte(BaseLenMax)

	// length的最大值
	ExtendLenMax = 255 * 255

	// opcode
	// 注册服务
	OpSign = 1
	// 请求
	OpRequest = 2
	// 响应
	OpResponse = 3
	// 错误
	OpError = 4
)

// Conn 封装的tcp连接的接口
// 现在实现的是一个普通的连接 后面或许会上连接池
type Conn interface {
	ReadData() (opcode uint8, data []byte, err error)
	read(len uint) (data []byte, err error)
	readByte() (d byte, err error)
	WriteSign([]byte) error
	WriteRequest([]byte) error
	WriteError([]byte) error
	Close() error
	Addr() string
	WriteResponse(response []byte) error
}

// handlerErr 如果err不为nil直接抛出来
func handlerErr(err error) {
	if err != nil {
		panic(err)
	}
}
