package begonia

import (
	"mashiroc.fun/begonia/entity"
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

	c := newContext(request)
	res := s.do(request.Fun, c)
	resp = entity.DefaultResponse{
		Uuid: request.UUID,
		Data: res,
	}
	return
}

// sign 注册服务
func (pc *ProcessCallHandler) sign(name string, in interface{}) (res []string) {
	funs := make([]reflect.Method, 0)

	v := reflect.ValueOf(in)
	t := reflect.TypeOf(in)
	num := t.NumMethod()
	for i := 0; i < num; i++ {
		m := t.Method(i)
		if checkFun(m) {
			funs = append(funs, m)
			res = append(res, m.Name)
		}
	}
	pc.addService(name, v, funs)
	return
}

func (pc *ProcessCallHandler) addService(name string, in reflect.Value, funs []reflect.Method) {
	rfs := make([]*remoteFun, len(funs))
	for i, fun := range funs {
		rf := &remoteFun{
			name: fun.Name,
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

func checkFun(m reflect.Method) bool {
	t := m.Type
	if t.NumIn() != 2 || t.In(1) != contextType || t.NumOut() != 0 {
		//TODO 处理
		return false
	}
	return true
}
