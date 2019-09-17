package begonia_rpc

import (
	"mashiroc.fun/begoniarpc/conn"
	"sync"
)

type serviceMap struct {
	data map[string]service
	lock sync.Mutex
}

func NewServiceMap(len uint) *serviceMap {
	m := make(map[string]service, len)
	return &serviceMap{
		data: m,
		lock: sync.Mutex{},
	}
}

func (s *serviceMap) Get(key string) (v service, ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	v, ok = s.data[key]
	return
}

func (s *serviceMap) Set(k string, v service) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[k] = v
}

func (s *serviceMap) Remove(k string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, k)
}

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
