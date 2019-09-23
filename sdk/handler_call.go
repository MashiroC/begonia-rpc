package begonia

import (
	"encoding/json"
	"github.com/MashiroC/begonia-rpc/conn"
	"github.com/MashiroC/begonia-rpc/entity"
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


