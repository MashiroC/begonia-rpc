package entity

type Request struct {
	UUID    string `json:"1"`
	Service string `json:"2"`
	Fun     string `json:"3"`
	Data    Param  `json:"4"`
}

type RespForm struct {
	Uuid string `json:"1"`
	Data Param  `json:"2"`
}

type SignForm struct {
	Sign []SignEntity `json:"1"`
}

type SignEntity struct {
	Name   string   `json:"1"`
	Fun    []string `json:"2"`
	IsMore bool     `json:"3"`
}

type ErrForm struct {
	Uuid    string `json:"1"`
	ErrCode string `json:"2"`
	ErrMsg  string `json:"3"`
}
