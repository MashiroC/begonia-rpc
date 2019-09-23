package begonia

import (
	"mashiroc.fun/begonia/entity"
	"mashiroc.fun/begonia/util/log"
	"sync"
)

type Context struct {
	lock  sync.Mutex
	uuid  string
	isRes bool
	Param entity.Param
	res   entity.Param
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

func (c *Context) JSON(param entity.Param) {
	c.write(param)
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
