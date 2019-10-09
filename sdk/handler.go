package begonia

import (
	"encoding/json"
	"fmt"
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/util/log"
)

// handlerRequest 处理远程调用请求帧
// 客户端收到的请求帧只能是被rpc调用的帧
func (cli *Client) handlerRequest(data []byte) {

	// 先检查data的json对不对 json不对直接关了
	req := entity.Request{}
	if err := json.Unmarshal(data, &req); err != nil {
		log.Error("request frame json [%s] error [%s] req addr: [%s]", string(data), err.Error(), cli.conn.Addr())
		cli.closeWith(err)
		return
	}

	// 检查一个帧要调用的service和function是否存在
	if req.Service == "" || req.Fun == "" {
		cli.respError(req.UUID, entity.ServiceNotFoundErr)
		return
	}

	log.Info("received [%s] call %s.%s", cli.conn.Addr(), req.Service, req.Fun)

	resp, err := cli.pc.call(req)
	if err != nil {
		// TODO:Err
		//cli.conn.WriteError()
		fmt.Println("errorerror", err, req.UUID)
		//cli.respError(req.UUID,err)
		return
	}

	_ = cli.conn.WriteResponse(resp.Response())
}

// handlerResponse 处理远程调用响应帧
// 客户端收到的响应帧只能是自己调用的响应
func (cli *Client) handlerResponse(data []byte) {
	// 先检查data的json对不对 json不对直接关了
	form := entity.RespForm{}
	if err := json.Unmarshal(data, &form); err != nil {
		log.Error("resp frame json [%s] error [%s] form addr: [%s]", string(data), err.Error(), cli.conn.Addr())
		cli.closeWith(err)
		return
	}

	// 先找uuid有没有 uuid没有就是没有注册回调
	if form.Uuid == "" {
		// Uuid not found 这个应该直接返回给这条连接
		cli.respError("", entity.CallbackNotSignedErr)
		return
	}

	// 有uuid 去回调
	if err := cli.resp.callback(form.Uuid, form.Data); err != nil {
		cli.respError(form.Uuid, err)
		return
	}
}

// handlerError 处理错误帧
// 这个错误帧指的是收到的error frame 不是异常帧
// 客户端收到的错误帧只能是rpc响应的错误
func (cli *Client) handlerError(data []byte) {
	// 先检查data的json对不对 json不对直接关了
	form := entity.ErrForm{}
	if err := json.Unmarshal(data, &form); err != nil {
		log.Error("error frame json [%s] error [%s] form addr: [%s]", string(data), err.Error(), cli.conn.Addr())
		cli.closeWith(err)
		return
	}

	// 这里和响应的逻辑基本一样 只不过回调传的是error
	if form.Uuid == "" {
		//cli.respError(entity.CallbackNotSignedErr)
		log.Error("fuck Uuid")
		return
	}

	// 回调的error
	cErr := entity.CallError{
		ErrCode:    form.ErrCode,
		ErrMessage: form.ErrMsg,
	}
	if err := cli.resp.callbackErr(form.Uuid, cErr); err != nil {
		cli.respError(form.Uuid, err)
	}
}
