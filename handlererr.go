package begonia_rpc

import (
	"encoding/json"
	"mashiroc.fun/begoniarpc/conn"
	"mashiroc.fun/begoniarpc/entity"
	"mashiroc.fun/begoniarpc/util/log"
)

func respError(conn conn.Conn, cErr entity.CallError) {
	log.Error(conn.Addr(), cErr)
	b, _ := json.Marshal(cErr)
	_ = conn.WriteError(b)
}
