package core

type Params map[string]interface{}

func (ps Params) Get(name string) (val interface{}, ok bool) {
	val, ok = ps[name]
	return
}
