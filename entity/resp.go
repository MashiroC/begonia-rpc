package entity

import (
	"encoding/json"
)

type Response interface {
	Response() []byte
}

type ErrResponse struct {
	Uuid    string
	ErrCode string
	ErrMsg  string
}

func (r ErrResponse) Response() []byte {
	tmp :=RespForm{
		Uuid: r.Uuid,
		Type: ErrorResponse,
		Data: Param{
			"errCode": r.ErrCode,
			"errMsg":  r.ErrMsg,
		},
	}
	b, _ := json.Marshal(tmp)
	return b
}

type DefaultResponse struct {
	Uuid string `json:"1"`
	Data interface{}  `json:"2"`
}

func (r DefaultResponse) Response() []byte {
	tmp := RespForm{
		Uuid: r.Uuid,
		Type: NormalResponse,
		Data: r.Data,
	}
	b, _ := json.Marshal(tmp)
	return b
}
