package begonia

import (
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/util/log"
	"sync"
)

type Context struct {
	lock  sync.Mutex
	uuid  string
	isRes bool
	Param entity.Param
	res   entity.Param
	err   error
}

func (c *Context) write(p entity.Param) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.isRes {
		log.Error("has result")
		return
	}
	c.res = p
}

func (c *Context) writeOne(in interface{}) {
	p := make(entity.Param, 1)
	p[c.uuid] = in
	c.write(p)
}

func (c *Context) String(s string) {
	c.writeOne(s)
}

func (c *Context) Int(i int) {
	c.writeOne(i)
}

func (c *Context) Error(err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.err = err
	c.isRes = true
}

func (c *Context) JSON(in interface{}) {
	if param, ok := in.(entity.Param); ok {
		c.write(param)
		return
	}

}

func newContext(req entity.Request) *Context {
	return &Context{
		lock:  sync.Mutex{},
		uuid:  req.UUID,
		isRes: false,
		Param: req.Data,
		res:   nil,
	}
}
