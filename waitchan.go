package begoniarpc

import (
	"fmt"
	"github.com/MashiroC/begonia-rpc/entity"
	"sync"
)

// waitChan.go 根据uuid获得响应回调的map

// WaitChan 根据uuid获得响应的回调
// 并发安全
type WaitChan struct {
	data map[string]func(response entity.Response)
	lock sync.Mutex
}

// 构造函数
func NewWaitChan(len uint) *WaitChan {
	return &WaitChan{
		data: make(map[string]func(entity.Response), len),
		lock: sync.Mutex{},
	}
}

// Get 取
func (w *WaitChan) Get(k string) (callback func(entity.Response), ok bool) {
	w.lock.Lock()
	defer w.lock.Unlock()
	callback, ok = w.data[k]
	if !ok {
		fmt.Println(callback)
	}
	return
}

// Set 加
func (w *WaitChan) Set(k string, callback func(entity.Response)) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.data[k] = callback
}

// Remove 删
func (w *WaitChan) Remove(k string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	delete(w.data, k)
}
