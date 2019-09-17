package begonia

import (
	"mashiroc.fun/begonia/conn"
	"mashiroc.fun/begonia/entity"
	"mashiroc.fun/begonia/util/log"
)

type respHandler struct {
	ch *WaitChan
}

func newRespHandler() *respHandler {
	return &respHandler{ch: NewWaitChan(100)}
}

func (h *respHandler) signCallBack(uuid string, conn conn.Conn) (err error) {
	ch := make(chan entity.Response, 1)
	h.ch.Set(uuid, func(resp entity.Response) {
		ch <- resp
	})
	go h.waitCallBack(ch, conn)
	return
}

// waitCallBack 等待回调成功
// 这里做成这样而不是直接用连接放到map里是为了做超时
func (h *respHandler) waitCallBack(ch chan entity.Response, conn conn.Conn) {
	resp := <-ch

	if _, ok := resp.(entity.DefaultResponse); ok {
		if err := conn.WriteResponse(resp.Response()); err != nil {
			log.Warn("write error ", err)
		}
	}

	if _, ok := resp.(entity.ErrResponse); ok {
		if err := conn.WriteError(resp.Response()); err != nil {
			log.Warn("write error ", err)
		}
	}

}

func (h *respHandler) CallBack(uuid string, data entity.Param) (err error) {
	f, ok := h.ch.Get(uuid)
	if !ok {
		return entity.CallError{
			ErrCode:    "114514",
			ErrMessage: "callback uuid not found :" + uuid,
		}
		//return entity.CallbackNotSignedErr
	}
	resp := entity.DefaultResponse{
		Uuid: uuid,
		Data: data,
	}
	f(resp)
	h.ch.Remove(uuid)
	return
}

func (h *respHandler) CallBackErr(uuid string, cErr entity.CallError) (err error) {
	f, ok := h.ch.Get(uuid)
	if !ok {
		return entity.CallbackNotSignedErr
	}
	resp := entity.ErrResponse{
		Uuid:    uuid,
		ErrCode: cErr.ErrCode,
		ErrMsg:  cErr.ErrMessage,
	}
	f(resp)
	return
}
