package begonia

import (
	"github.com/MashiroC/begonia-rpc/entity"
	"reflect"
)

var (
	contextType = reflect.TypeOf(&Context{})
)

type ProcessCallHandler struct {
	service *serviceMap
}

func (pc *ProcessCallHandler) call(request entity.Request) (resp entity.Response, err error) {

	s, ok := pc.service.Get(request.Service)
	if !ok {
		err = entity.ServiceNotFoundErr
		return
	}

	//c := newContext(request)
	res, err := s.do(request.Fun,request.Data)
	if err != nil {
		cErr, ok := err.(entity.CallError)
		if !ok {
			cErr = entity.CallError{
				ErrCode:    "114514",
				ErrMessage: err.Error(),
			}
		}
		resp = entity.ErrResponse{
			Uuid:    request.UUID,
			ErrCode: cErr.ErrCode,
			ErrMsg:  cErr.ErrMessage,
		}
		err = nil
	} else {
		resp = entity.DefaultResponse{
			Uuid: request.UUID,
			Data: res,
		}
	}

	return
}

// sign 注册服务
func (pc *ProcessCallHandler) sign(name string, in interface{}) (res []entity.FunEntity) {
	funs := make([]reflect.Method, 0)

	v := reflect.ValueOf(in)
	t := reflect.TypeOf(in)
	num := t.NumMethod()
	for i := 0; i < num; i++ {
		m := t.Method(i)
		funs = append(funs, m)
		en := entity.FunEntity{
			Name: m.Name,
			Size: m.Type.NumIn(),
		}
		res = append(res, en)
	}
	pc.addService(name, v, funs)
	return
}

func (pc *ProcessCallHandler) addService(name string, in reflect.Value, funs []reflect.Method) {
	rfs := make([]*remoteFun, len(funs))
	for i, fun := range funs {
		in := make([]reflect.Type, fun.Type.NumIn())
		for i := 0; i < len(in); i++ {
			funIn := fun.Type.In(i)
			in[i] = funIn
		}

		rf := &remoteFun{
			name: fun.Name,
			in:   in,
			fun:  fun.Func,
		}
		rfs[i] = rf
	}
	s := &service{
		name: name,
		fun:  rfs,
		in:   in,
	}
	pc.service.Set(name, s)
}
