package begonia

import (
	"fmt"
	"mashiroc.fun/begonia/entity"
	"sync"
)

type WaitChan struct {
	data map[string]func(response entity.Response)
	lock sync.Mutex
}

func NewWaitChan(len uint) *WaitChan {
	return &WaitChan{
		data: make(map[string]func(entity.Response), len),
		lock: sync.Mutex{},
	}
}

func (w *WaitChan) Get(k string) (callback func(entity.Response), ok bool) {
	w.lock.Lock()
	defer w.lock.Unlock()
	callback, ok = w.data[k]
	if !ok {
		fmt.Println(callback)
	}
	return
}

func (w *WaitChan) Set(k string, callback func(entity.Response)) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.data[k] = callback
}

func (w *WaitChan) Remove(k string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	delete(w.data, k)
}
