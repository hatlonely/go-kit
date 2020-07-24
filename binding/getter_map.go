package binding

type MapGetter map[string]interface{}

func (g MapGetter) Get(key string) (interface{}, bool) {
	v, ok := g[key]
	return v, ok
}
