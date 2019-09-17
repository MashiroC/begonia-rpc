package entity

import (
	"fmt"
)

type CallError struct {
	ErrCode    string `json:"errCode"`
	ErrMessage string `json:"errMsg"`
}

var (
	ServiceSignedErr = CallError{
		ErrCode:    "114514",
		ErrMessage: "service has signed",
	}

	ServiceNotFoundErr = CallError{
		ErrCode:    "114514",
		ErrMessage: "service not found",
	}

	FunctionNotFoundErr = CallError{
		ErrCode:    "114514",
		ErrMessage: "function not found",
	}

	CallbackUuidNotFoundErr = CallError{
		ErrCode:    "114514",
		ErrMessage: "callback uuid not found",
	}
)

func (c CallError) Error() string {
	return fmt.Sprintf("errCode:%s errMsg: %s", c.ErrCode, c.ErrMessage)
}

func NewError(code, msg string) CallError {
	return CallError{
		ErrCode:    code,
		ErrMessage: msg,
	}
}
