package bind

import (
	"sort"

	"github.com/hatlonely/go-kit/refx"
)

func NewMapGetter(m map[string]interface{}) MapGetter {
	var v interface{}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if err := refx.InterfaceSet(&v, key, m[key]); err != nil {
			panic(err)
		}
	}
	return MapGetter{v: v}
}

type MapGetter struct {
	v interface{}
}

func (g MapGetter) Get(key string) (interface{}, bool) {
	v, err := refx.InterfaceGet(g.v, key)
	if err != nil {
		return nil, false
	}
	return v, true
}
