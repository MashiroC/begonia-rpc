package begonia

// service.go 服务的实体和存服务的map

import (
	"errors"
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/util/log"
	"reflect"
	"sync"
)

// service 注册的服务
type service struct {
	name string
	fun  []*remoteFun
	in   reflect.Value
}

func (s *service) do(fun string, in []interface{}) (res entity.Param, err error) {
	// find func
	defer func() {
		if re := recover(); re != nil {
			if reErr, ok := re.(error); ok {
				log.Error("recover error:%s", reErr.Error())
				err = reErr
				return
			}
			log.Error("recover unknown:%s", re)
			err = errors.New("unknown error")
		}
	}()
	var f *remoteFun
	for _, v := range s.fun {
		if v.name == fun {
			f = v
			break
		}
	}

	return f.do(s.in, in)
}

type remoteFun struct {
	name string
	in   []reflect.Type
	fun  reflect.Value
}

func (rf *remoteFun) do(value reflect.Value, in []interface{}) (entity.Param, error) {
	if len(rf.in) != len(in) {
		return nil,entity.ParamsNumErr
	}

	for _,param:=range in {
		t:=reflect.TypeOf(param)

	}
	rf.fun.Call([]reflect.Value{value, reflect.ValueOf(c)})
	return c.res, c.err
}

// serviceMap 存服务的实体 并发安全
// key是服务名 value是服务实体
type serviceMap struct {
	data map[string]*service
	lock sync.Mutex
}

// newServiceMap 构造函数
func newServiceMap(len uint) *serviceMap {
	m := make(map[string]*service, len)
	return &serviceMap{
		data: m,
		lock: sync.Mutex{},
	}
}

// Get 拿服务
func (s *serviceMap) Get(key string) (v *service, ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	v, ok = s.data[key]
	return
}

// Set 加服务
func (s *serviceMap) Set(k string, v *service) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[k] = v
}

// Remove 移除服务
func (s *serviceMap) Remove(k string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, k)
}
