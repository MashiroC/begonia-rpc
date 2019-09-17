package begonia_rpc

import (
	"encoding/json"
	"mashiroc.fun/begoniarpc/conn"
	"mashiroc.fun/begoniarpc/entity"
	"mashiroc.fun/begoniarpc/util/log"
)

type callHandler struct {
	remoteFun *serviceMap
}

type request struct {
	Uuid  string                 `json:"1"`
	Name  string                 `json:"2"`
	Fun   string                 `json:"3"`
	Param map[string]interface{} `json:"4"`
}

type service struct {
	name string
	fun  []string
	c    conn.Conn
}

func (h *callHandler) signService(s service) (err error) {
	if _, ok := h.remoteFun.Get(s.name); ok {
		log.Warn("service has signed", s)
		return entity.ServiceSignedErr
	}
	log.Info(s.c.Addr(), "注册服务:", s.name, "[", s.fun, "]")
	h.remoteFun.Set(s.name, s)
	return
}

func (h *callHandler) unsignService(s service) (err error) {
	if _, ok := h.remoteFun.Get(s.name); !ok {
		log.Warn("service not sign", s)
		return entity.ServiceSignedErr
	}

	h.remoteFun.Remove(s.name)
	return
}

func (h *callHandler) unBindService(c conn.Conn) (err error) {
	h.remoteFun.Unbind(c)
	return
}

func (h *callHandler) call(uuid, name, fun string, param entity.Param) (err error) {
	service, ok := h.remoteFun.Get(name)
	if !ok {
		log.Warn("service not found", err)
		return entity.ServiceNotFoundErr
	}

	for _, f := range service.fun {
		if f == fun {
			req := request{
				Uuid:  uuid,
				Name:  name,
				Fun:   fun,
				Param: param,
			}
			b, _ := json.Marshal(req)
			if err := service.c.WriteRequest(b); err != nil {
				log.Warn(err)
			}
			return
		}
	}
	log.Warn("function not found", err)
	return entity.FunctionNotFoundErr
}

func newCallHandler() *callHandler {
	return &callHandler{remoteFun: NewServiceMap(5)}
}
