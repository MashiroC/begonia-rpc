package entity

import (
	"fmt"
)

type CallError struct {
	ErrCode    string `json:"errCode"`
	ErrMessage string `json:"errMsg"`
}

var (
	ErrCodeUnknow = "114514"

	ServiceSignedErr = CallError{
		ErrCode:    "114514",
		ErrMessage: "server has signed",
	}

	ServiceNotFoundErr = CallError{
		ErrCode:    "114514",
		ErrMessage: "server not found",
	}

	FunctionNotFoundErr = CallError{
		ErrCode:    "114514",
		ErrMessage: "function not found",
	}

	CallbackNotSignedErr = CallError{
		ErrCode:    "114514",
		ErrMessage: "callback uuid not found",
	}

	ParamsNumErr = CallError{
		ErrCode:    "114514",
		ErrMessage: "params num failed",
	}

	RespEmptyError = CallError{
		ErrCode:    "114514",
		ErrMessage: "response empty",
	}

	RespTypeError = CallError{
		ErrCode:    "114514",
		ErrMessage: "response kind not allow",
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
