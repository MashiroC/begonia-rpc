package entity

type Request struct {
	UUID    string `json:"1"`
	Service string `json:"2"`
	Fun     string `json:"3"`
	Data    []interface{}  `json:"4"`
}

type RespForm struct {
	Uuid string `json:"1"`
	Type int `json:"2"`
	Data interface{}  `json:"3"`
}

type SignForm struct {
	Sign []SignEntity `json:"1"`
}

type SignEntity struct {
	Name   string   `json:"1"`
	Fun    []FunEntity `json:"2"`
	IsMore bool     `json:"3"`
}

type FunEntity struct {
	Name string `json:"1"`
	Size int `json:"2"`
}

type ErrForm struct {
	Uuid    string `json:"1"`
	ErrCode string `json:"2"`
	ErrMsg  string `json:"3"`
}
