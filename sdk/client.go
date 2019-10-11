package begonia

import (
	"encoding/json"
	begoniarpc "github.com/MashiroC/begonia-rpc"
	"github.com/MashiroC/begonia-rpc/conn"
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/util/log"
	"github.com/satori/go.uuid"
	"time"
)

type Client struct {
	addr      string
	signCache []entity.SignEntity
	call      *CallHandler
	resp      *ResponseHandler
	pc        *ProcessCallHandler
	conn      conn.Conn
	wait      chan bool
}

type Callback = func(*Response)

type CallbackChan = chan entity.Response

// New 创建客户端并监听端口
func New(addr string) *Client {
	cli := &Client{
		addr:      addr,
		signCache: make([]entity.SignEntity, 0),
		call:      &CallHandler{},
		resp:      &ResponseHandler{cbMap: begoniarpc.NewWaitChan(255)},
		pc:        &ProcessCallHandler{service: newServiceMap(5)},
		wait:      make(chan bool, 2),
	}
	cli.connectAndListen()
	return cli
}

func (cli *Client) connectAndListen() {
	c := cli.connection(cli.addr)
	cli.conn = c
	cli.call.conn = c
	go cli.listen()
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
		//TODO: handler Err
		cb(newErrorResponse(uuid, err))
	}

	if err := cli.call.call(uuid, r); err != nil {
		//TODO: handler Err
		cb(newErrorResponse(uuid, err))
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
	cli.signCache = append(cli.signCache, e)
	form := entity.SignForm{Sign: []entity.SignEntity{e}}
	b, _ := json.Marshal(form)

	_ = cli.conn.WriteSign(b)
}

func (cli *Client) Wait() {
	<-cli.wait
}

func (cli *Client) KeepConnect() {

	for {
		<-cli.wait

		time.Sleep(3 * time.Second)
		ok := Must(cli.connectAndListen, cli.wait)
		if ok {
			Must(cli.reSign, cli.wait)
		}

	}

}

// reSign 断开连接后重新注册服务
func (cli *Client) reSign() {
	if len(cli.signCache) != 0 {
		form := entity.SignForm{Sign: cli.signCache}
		b, _ := json.Marshal(form)

		_ = cli.conn.WriteSign(b)
	}
}

type RemoteService struct {
	cli     *Client
	Service string
}

func (s RemoteService) Call(fun string, argv entity.Param) *Response {
	req := Request{
		Service:  s.Service,
		Function: fun,
		Param:    argv,
	}
	return s.cli.Call(req)
}

func (s RemoteService) CallAsync(fun string, argv entity.Param, cb Callback) {
	req := Request{
		Service:  s.Service,
		Function: fun,
		Param:    argv,
	}
	s.cli.CallAsync(req, cb)
}

type RemoteFun struct {
	cli     *Client
	Service string
	Fun     string
}

func (f RemoteFun) Call(argv entity.Param) *Response {
	req := Request{
		Service:  f.Service,
		Function: f.Fun,
		Param:    argv,
	}
	return f.cli.Call(req)
}

func (f RemoteFun) CallAsync(argv entity.Param, cb Callback) {
	req := Request{
		Service:  f.Service,
		Function: f.Fun,
		Param:    argv,
	}
	f.cli.CallAsync(req, cb)
}

func (cli *Client) Service(s string) RemoteService {
	return RemoteService{
		cli:     cli,
		Service: s,
	}
}

func (s RemoteService) Fun(f string) RemoteFun {
	return RemoteFun{
		cli:     s.cli,
		Service: s.Service,
		Fun:     f,
	}
}

func Must(fun func(), ch chan bool) (res bool) {
	defer func() {
		if re := recover(); re != nil {
			log.Warn("recover something : %s", re)
			ch <- true
			res = false
		}
	}()
	fun()
	res = true
	return
}
