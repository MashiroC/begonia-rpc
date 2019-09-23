package begonia

import (
	"encoding/json"
	"mashiroc.fun/begonia/conn"
	"mashiroc.fun/begonia/entity"
)

type CallHandler struct {
	conn conn.Conn
}

func (h *CallHandler) call(uuid string, request Request) (err error) {
	req := entity.Request{
		UUID:    uuid,
		Service: request.Service,
		Fun:     request.Function,
		Data:    request.Param,
	}
	b, err := json.Marshal(req)
	err = h.conn.WriteRequest(b)
	return
}


