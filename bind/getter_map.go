package bind

func NewMapGetter(m map[string]interface{}) MapGetter {
	return m
}

type MapGetter map[string]interface{}

func (g MapGetter) Get(key string) (interface{}, bool) {
	v, ok := g[key]
	return v, ok
}
