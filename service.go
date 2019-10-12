package begoniarpc

// service.go 服务的实体和存服务的map

import (
	"github.com/MashiroC/begonia-rpc/conn"
	"github.com/MashiroC/begonia-rpc/entity"
	"sync"
)

// service 注册的服务
type service struct {
	name string
	fun  []entity.FunEntity
	c    conn.Conn
}

// serviceMap 存服务的实体 并发安全
// key是服务名 value是服务实体
type serviceMap struct {
	data map[string]service
	lock sync.Mutex
}

// newServiceMap 构造函数
func newServiceMap(len uint) *serviceMap {
	m := make(map[string]service, len)
	return &serviceMap{
		data: m,
		lock: sync.Mutex{},
	}
}

// Get 拿服务
func (s *serviceMap) Get(key string) (v service, ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	v, ok = s.data[key]
	return
}

// Set 加服务
func (s *serviceMap) Set(k string, v service) {
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

// Unbind 根据value来删除服务
func (s *serviceMap) Unbind(conn conn.Conn) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, c := range s.data {
		if c.c == conn {
			delete(s.data, c.name)
			return
		}
	}
}
