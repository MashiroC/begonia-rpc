package entity

import "encoding/json"

type Response interface {
	Response() []byte
}

type ErrResponse struct {
	Uuid    string
	ErrCode string
	ErrMsg  string
}

func (r ErrResponse) Response() []byte {
	b, _ := json.Marshal(r)
	return b
}

type DefaultResponse struct {
	Uuid string `json:"1"`
	Data Param  `json:"2"`
}

func (r DefaultResponse) Response() []byte {
	b, _ := json.Marshal(r)
	return b
}
