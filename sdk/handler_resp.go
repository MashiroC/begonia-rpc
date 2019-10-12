package begonia

import (
	begoniarpc "github.com/MashiroC/begonia-rpc"
	"github.com/MashiroC/begonia-rpc/entity"
)

type ResponseHandler struct {
	cbMap *begoniarpc.WaitChan
}

// return callback chan
func (h *ResponseHandler) signCallback(uuid string, request Request) (CallbackChan, error) {
	ch := make(CallbackChan, 1)
	h.cbMap.Set(uuid, func(resp entity.Response) {
		ch <- resp
	})
	return ch, nil
}

func (h *ResponseHandler) callback(uuid string, params interface{}) (err error) {
	f, ok := h.cbMap.Get(uuid)
	h.cbMap.Remove(uuid)

	if !ok {
		err = entity.CallbackNotSignedErr
		return
	}
	resp := entity.DefaultResponse{
		Uuid: uuid,
		Data: params,
	}
	f(resp)

	return
}

func (h *ResponseHandler) callbackErr(uuid string, cErr entity.CallError) (err error) {
	f, ok := h.cbMap.Get(uuid)
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
