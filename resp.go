package begoniarpc

// repo.go 响应处理的handler

import (
	"encoding/json"
	"github.com/MashiroC/begonia-rpc/conn"
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/util/log"
)

// respHandler 响应处理的实体
type respHandler struct {
	ch *WaitChan
}

// newRespHandler 构造函数
func newRespHandler() *respHandler {
	return &respHandler{ch: NewWaitChan(100)}
}

// signCallback 注册一个回调
func (h *respHandler) signCallback(uuid string, conn conn.Conn) (err error) {
	ch := make(chan entity.Response, 1)
	h.ch.Set(uuid, func(resp entity.Response) {
		ch <- resp
	})
	go h.waitCallback(ch, conn)
	return
}

// waitCallback 等待回调
// 这里做成这样而不是直接用连接放到map里是为了做超时(大概
func (h *respHandler) waitCallback(ch chan entity.Response, conn conn.Conn) {
	resp := <-ch
	var in interface{}
	_ = json.Unmarshal(resp.Response(), &in)
	if err := conn.WriteResponse(in); err != nil {
		log.Warn("write error ", err)
	}

}

// callback 请求完了之后等待响应 响应到了就回调
func (h *respHandler) Callback(uuid string,model int, data interface{}) (err error) {
	f, ok := h.ch.Get(uuid)
	h.ch.Remove(uuid)

	if !ok {
		return entity.CallError{
			ErrCode:    "114514",
			ErrMessage: "callback uuid not found :" + uuid,
		}
		//return entity.CallbackNotSignedErr
	}

	var resp entity.Response
	if model == entity.NormalResponse{
		resp = entity.DefaultResponse{
			Uuid: uuid,
			Data: data,
		}
	}else if model == entity.ErrorResponse{
		errMap := data.(map[string]interface{})
		resp = entity.ErrResponse{
			Uuid:    uuid,
			ErrCode: errMap["errCode"].(string),
			ErrMsg:  errMap["errMsg"].(string),
		}
	}


	log.Info("call %s resp", uuid)
	f(resp)
	return
}

// callbackErr 服务端发来err 回调把这个err返回给client
func (h *respHandler) CallbackErr(uuid string, cErr entity.CallError) (err error) {
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
