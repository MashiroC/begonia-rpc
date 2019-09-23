package begonia

import (
	"mashiroc.fun/begonia/entity"
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
	uuid string
	data entity.Param
	err  error
}

func (r Response) Error() error {
	return r.err
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
	if r.err != nil {
		res = or
		return
	}

	defer func() {
		if err := recover(); err != nil {
			r.err = err.(error)
			// or
			res = or
		}
	}()

	res = r.data[r.uuid]

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
	return r.data
}

func newResponseFromEntity(resp entity.Response) (res *Response) {
	if r, ok := resp.(entity.DefaultResponse); ok {
		res = &Response{
			uuid: r.Uuid,
			data: r.Data,
			err:  nil,
		}
		return
	}

	if r, ok := resp.(entity.ErrResponse); ok {
		res = &Response{
			uuid: r.Uuid,
			data: nil,
			err:  entity.NewError(r.ErrCode, r.ErrMsg),
		}
		return
	}
	return
}

func newErrorResponse(err error) *Response {
	return &Response{
		uuid: "",
		data: nil,
		err:  err,
	}
}
