package begonia

import "github.com/MashiroC/begonia-rpc/entity"

type Request struct {
	Service  string
	Function string
	Param    entity.Param
}
