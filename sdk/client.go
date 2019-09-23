package begonia

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	begoniarpc "mashiroc.fun/begonia"
	"mashiroc.fun/begonia/conn"
	"mashiroc.fun/begonia/entity"
)

type Client struct {
	call *CallHandler
	resp *ResponseHandler
	pc   *ProcessCallHandler
	conn conn.Conn
	wait chan bool
}

type Callback = func(*Response)

type CallbackChan = chan entity.Response

// New 创建客户端并监听端口
func New(addr string) *Client {
	cli := &Client{
		call: &CallHandler{},
		resp: &ResponseHandler{cbMap: begoniarpc.NewWaitChan(255)},
		pc:   &ProcessCallHandler{service: newServiceMap(5)},
		wait: make(chan bool, 2),
	}
	c := cli.connection(addr)
	cli.conn = c
	cli.call.conn = c
	go cli.listen()
	return cli
}

// call 同步调用
func (cli *Client) Call(r Request) (res *Response) {
	ch := make(chan *Response, 1)
	cli.CallAsync(r, func(resp *Response) {
		ch <- resp
	})

	res = <-ch
	return
}

// CallAsync 异步调用
func (cli *Client) CallAsync(r Request, cb Callback) {
	uuid := uuid.NewV4().String()
	cbCh, err := cli.resp.signCallback(uuid, r)
	if err != nil {
		//TODO: handler err
		cb(newErrorResponse(err))
	}

	if err := cli.call.call(uuid, r); err != nil {
		//TODO: handler err
		cb(newErrorResponse(err))
	}

	go func(cbCh CallbackChan) {
		resp := <-cbCh
		res := newResponseFromEntity(resp)
		cb(res)
	}(cbCh)

}

// Sign 注册服务
func (cli *Client) Sign(name string, in interface{}) {
	fun := cli.pc.sign(name, in)

	e := entity.SignEntity{
		Name:   name,
		Fun:    fun,
		IsMore: false,
	}
	form := entity.SignForm{Sign: []entity.SignEntity{e}}
	b, _ := json.Marshal(form)

	_ = cli.conn.WriteSign(b)
}

func (cli *Client) Wait() {
	<-cli.wait
}
