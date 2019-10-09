package begonia

import (
	"github.com/MashiroC/begonia-rpc/entity"
	"reflect"
)

var (
	defineType = []reflect.Kind{reflect.String, reflect.Int, reflect.Bool, reflect.Array, reflect.Slice, reflect.Float32, reflect.Float64}
)

const (
	emptyString = ""
	emptyInt    = 0
)

type Response struct {
	Uuid string
	Data entity.Param
	Err  error
}

func (r Response) Error() error {
	return r.Err
}

func (r *Response) Int() (res int) {
	return r.respOnce(emptyInt, reflect.Int).(int)

}

func (r *Response) IntOr(or int) (res int) {
	return r.respOnce(or, reflect.Int).(int)

}

func (r *Response) String() (res string) {
	return r.respOnce(emptyString, reflect.String).(string)
}

func (r *Response) StringOr(or string) (res string) {
	return r.respOnce(or, reflect.String).(string)
}

func (r *Response) respOnce(or interface{}, k reflect.Kind) (res interface{}) {
	if r.Err != nil {
		res = or
		return
	}

	defer func() {
		if err := recover(); err != nil {
			r.Err = err.(error)
			// or
			res = or
		}
	}()

	res = r.Data[r.Uuid]

	if res == nil {
		panic(entity.RespEmptyError)
	}
	kind := reflect.TypeOf(res).Kind()
	if kind != k {
		if kind == reflect.Float64 && k == reflect.Int {
			res = int(res.(float64))
			return
		}
		panic(entity.RespTypeError)
	}
	return
}

func (r *Response) Params() entity.Param {
	return r.Data
}

func newResponseFromEntity(resp entity.Response) (res *Response) {
	if r, ok := resp.(entity.DefaultResponse); ok {
		res = &Response{
			Uuid: r.Uuid,
			Err:  nil,
		}
		if _, ok := r.Data["errCode"]; ok {
			res.Err = entity.NewError(r.Data["errCode"].(string), r.Data["errMsg"].(string))
		} else {
			res.Data = r.Data
		}

		return
	}

	if r, ok := resp.(entity.ErrResponse); ok {
		res = &Response{
			Uuid: r.Uuid,
			Data: nil,
			Err:  entity.NewError(r.ErrCode, r.ErrMsg),
		}
		return
	}
	return
}

func newErrorResponse(uuid string, err error) *Response {
	return &Response{
		Uuid: uuid,
		Data: nil,
		Err:  err,
	}
}
