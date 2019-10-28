package begoniarpc

// handler_call.go 远程调用的处理中心

import (
	"github.com/MashiroC/begonia-rpc/conn"
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/util/log"
)

// callHandler 处理远程调用请求的实体
type callHandler struct {
	remoteFun *serviceMap // 所有注册的服务
}

// signService 注册服务
func (h *callHandler) signService(s service) (err error) {
	if _, ok := h.remoteFun.Get(s.name); ok {
		log.Warn("servive [%s] has signed", s.name)
		return entity.ServiceSignedErr
	}
	h.remoteFun.Set(s.name, s)
	return
}

// unsignService 注销服务
func (h *callHandler) unsignService(s service) (err error) {
	if _, ok := h.remoteFun.Get(s.name); !ok {
		log.Warn("service [%s] not sign", s.name)
		return entity.ServiceSignedErr
	}

	h.remoteFun.Remove(s.name)
	return
}

// unBindService 解绑和conn有关的服务
func (h *callHandler) unBindService(c conn.Conn) (err error) {
	h.remoteFun.Unbind(c)
	return
}

// call 远程调用
func (h *callHandler) call(req entity.Request,b []byte) (err error) {
	service, ok := h.remoteFun.Get(req.Service)
	if !ok {
		log.Warn("call service [%s] not found", req.Service)
		return entity.ServiceNotFoundErr
	}

	for _, f := range service.fun {
		if f.Name == req.Fun {
			if len(req.Data) != f.Size {
				err = entity.ParamsNumErr
				return
			}
			//b, _ := json.Marshal(req)
			err = service.c.WriteRequest(req)
			return
		}
	}
	log.Warn("call function [%s] not found", req.Fun)
	return entity.FunctionNotFoundErr
}

// newCallHandler 构造函数
func newCallHandler() *callHandler {
	return &callHandler{remoteFun: newServiceMap(5)}
}
