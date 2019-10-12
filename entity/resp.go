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
	tmp := DefaultResponse{
		Uuid: r.Uuid,
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
	b, _ := json.Marshal(r)
	return b
}
