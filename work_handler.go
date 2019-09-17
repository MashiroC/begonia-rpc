package begonia_rpc

import (
	"encoding/json"
	"fmt"
	"mashiroc.fun/begoniarpc/conn"
	"mashiroc.fun/begoniarpc/entity"
	"mashiroc.fun/begoniarpc/util/log"
)

//.handlerSign(conn,data)

func (s *Server) handlerSign(conn conn.Conn, data []byte) {
	form := entity.SignForm{}
	err := json.Unmarshal(data, &form)
	if err != nil {
		log.Warn("addr:", conn.Addr(), "data json error:", string(data))
		s.closeWith(conn, err)
		return
	}

	for _, si := range form.Sign {
		if si.IsMore {
			// TODO:已经注册服务 又开了个连接
			log.Error("isMore")
		} else {
			// 第一次注册连接
			ser := service{
				name: si.Name,
				fun:  si.Fun,
				c:    conn,
			}
			if err := s.call.signService(ser); err != nil {
				respError(conn, err.(entity.CallError))
			}
		}
	}

}

func (s *Server) handlerRequest(conn conn.Conn, data []byte) {
	form := entity.RequestForm{}
	if err := json.Unmarshal(data, &form); err != nil {
		log.Warn("addr:", conn.Addr(), "data json error:", string(data))
		s.closeWith(conn, err)
		return
	}

	if form.Service == "" || form.Fun == "" {
		log.Warn("call param error")
		return
	}

	log.Info(conn.Addr(), "call", form.Uuid, form.Service, form.Fun)

	// 注册一个回调
	if err := s.resp.signCallBack(form.Uuid, conn); err != nil {
		respError(conn, err.(entity.CallError))
	}

	if err := s.call.call(form.Uuid, form.Service, form.Fun, form.Data); err != nil {
		// call error 这个应该直接返回给这条连接
		respError(conn, err.(entity.CallError))
		return
	}

}

func (s *Server) handlerResponse(conn conn.Conn, data []byte) {
	form := entity.RespForm{}
	if err := json.Unmarshal(data, &form); err != nil {
		log.Warn("resp addr:", conn.Addr(), "data json error:", string(data))
		s.closeWith(conn, err)
		return
	}

	if form.Uuid == "" {
		// TODO:errcode
		err := entity.CallError{
			ErrCode:    "114514",
			ErrMessage: "uuid not found",
		}
		// uuid not found 这个应该直接返回给这条连接
		respError(conn, err)
		return
	}

	if err := s.resp.CallBack(form.Uuid, form.Data); err != nil {
		s.resp.ch.lock.Lock()
		fmt.Println(s.resp.ch.data)
		s.resp.ch.lock.Unlock()

		respError(conn, err.(entity.CallError))
		return
	}
}

func (s *Server) handlerError(conn conn.Conn, data []byte) {
	form := entity.ErrForm{}
	if err := json.Unmarshal(data, &form); err != nil {
		log.Warn("error addr:", conn.Addr(), "data json error:", string(data))
		s.closeWith(conn, err)
		return
	}

	if form.Uuid == "" {
		// TODO:errcode
		err := entity.CallError{
			ErrCode:    "114514",
			ErrMessage: "uuid not found",
		}
		respError(conn, err)
		log.Warn(err.Error())
		return
	}

	err := entity.CallError{
		ErrCode:    form.ErrCode,
		ErrMessage: form.ErrMsg,
	}
	if err := s.resp.CallBackErr(form.Uuid, err); err != nil {
		respError(conn, err.(entity.CallError))
	}
}

func (s *Server) closeWith(conn conn.Conn, err error) {
	cErr := entity.CallError{
		ErrCode:    "114514",
		ErrMessage: err.Error(),
	}
	b, _ := json.Marshal(cErr)
	_ = conn.WriteError(b)
	if err := s.call.unBindService(conn); err != nil {
		log.Warn(err.Error())
	}
}
