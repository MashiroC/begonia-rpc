package begonia

type Request struct {
	Service  string
	Function string
	Param    []interface{}
}
