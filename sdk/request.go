package begonia

import "mashiroc.fun/begonia/entity"

type Request struct {
	Service  string
	Function string
	Param    entity.Param
}
