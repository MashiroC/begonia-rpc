package begoniarpc

// handler.go 处理各种帧的函数

import (
	"encoding/json"
	"fmt"
	"github.com/MashiroC/begonia-rpc/conn"
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/util/log"
)

// handlerSign 处理注册帧
func (s *ServerCenter) handlerSign(conn conn.Conn, data []byte) {
	// 先检查data的json对不对 json不对直接关了
	form := entity.SignForm{}
	err := json.Unmarshal(data, &form)
	if err != nil {
		log.Error("sign frame json [%s] error [%s] form addr: [%s]", string(data), err.Error(), conn.Addr())
		s.closeWith(conn, err)
		return
	}

	// 遍历json里的每一个服务 服务端一次性可以注册多个服务
	for _, si := range form.Sign {

		if si.IsMore {

			// TODO:已经注册服务 服务用了连接池 又开了个连接
			log.Error("service has signed")

		} else {

			// 第一次注册连接
			log.Info("addr [%s] signed service [%s] for function %s", conn.Addr(), si.Name, si.Fun)

			ser := service{
				name: si.Name,
				fun:  si.Fun,
				c:    conn,
			}
			// 注册服务 如果有错误 直接把连接关了 错误返回到注册方
			if err := s.call.signService(ser); err != nil {
				respError(conn, "", err.(entity.CallError))
			}

		}
	}

}

// handlerRequest 处理远程调用请求帧
func (s *ServerCenter) handlerRequest(conn conn.Conn, data []byte) {
	// 先检查data的json对不对 json不对直接关了
	form := entity.Request{}
	if err := json.Unmarshal(data, &form); err != nil {
		log.Error("request frame json [%s] error [%s] form addr: [%s]", string(data), err.Error(), conn.Addr())
		s.closeWith(conn, err)
		return
	}

	// 检查一个帧要调用的service和function是否存在
	if form.Service == "" || form.Fun == "" {
		respError(conn, form.UUID, entity.ServiceNotFoundErr)
		return
	}

	log.Info("remote [%s] call %s.%s", conn.Addr(), form.Service, form.Fun)

	// 这里要先注册回调再发request
	// 否则的话 在高并发场景下 因为协程调度的问题
	// 会先收到response包再注册回调 这样子uuid会找不到

	// 注册一个回调
	if err := s.resp.signCallback(form.UUID, conn); err != nil {
		respError(conn, form.UUID, err.(entity.CallError))
	}

	// Remote Process call
	if err := s.call.call(form,data); err != nil {
		// call error 这个应该直接返回给这条连接
		respError(conn, form.UUID, err.(entity.CallError))
		return
	}

}

// handlerResponse 处理远程调用响应帧
func (s *ServerCenter) handlerResponse(conn conn.Conn, data []byte) {

	// 先检查data的json对不对 json不对直接关了
	form := entity.RespForm{}
	if err := json.Unmarshal(data, &form); err != nil {
		log.Error("resp frame json [%s] error [%s] form addr: [%s]", string(data), err.Error(), conn.Addr())
		s.closeWith(conn, err)
		return
	}

	// 先找uuid有没有 uuid没有就是没有注册回调
	if form.Uuid == "" {
		// uuid not found 这个应该直接返回给这条连接
		respError(conn, form.Uuid, entity.CallbackNotSignedErr)
		return
	}

	// 有uuid 去回调
	if err := s.resp.Callback(form.Uuid, form.Data); err != nil {
		//respError(conn,form.Uuid,err.(entity.CallError))
		fmt.Println("fuckfuck")
		return
	}
}

// handlerError 处理错误帧
// 这个错误帧指的是收到的error frame 不是异常帧
func (s *ServerCenter) handlerError(conn conn.Conn, data []byte) {
	// 先检查data的json对不对 json不对直接关了
	form := entity.ErrForm{}
	if err := json.Unmarshal(data, &form); err != nil {
		fmt.Println(err)
		fmt.Println(string(data))
		//log.Error("error frame json [%s] error [%s] form addr: [%s]", string(data), err.Error(), conn.Addr())
		s.closeWith(conn, err)
		return
	}

	// 这里和响应的逻辑基本一样 只不过回调传的是error
	if form.Uuid == "" {
		//respError(conn, "", entity.CallbackNotSignedErr)
		return
	}

	// 回调的error
	cErr := entity.CallError{
		ErrCode:    form.ErrCode,
		ErrMessage: form.ErrMsg,
	}
	if err := s.resp.CallbackErr(form.Uuid, cErr); err != nil {
		//respError(conn,form.Uuid ,err.(entity.CallError))
	}
}

// closeWith 向这条连接写个错误并关闭这条连接
func (s *ServerCenter) closeWith(conn conn.Conn, err error) {

	// 如果这个错误是CallError 那么直接序列化 如果不是 把err的信息写一个新的CallError
	cErr, ok := err.(entity.CallError)
	if !ok {
		cErr = entity.NewError(entity.ErrCodeUnknow, err.Error())
	}

	b, _ := json.Marshal(cErr)
	_ = conn.WriteError(b)
	if err := s.call.unBindService(conn); err != nil {
		log.Warn(err.Error())
	}
}
