package bind

import (
	"github.com/hatlonely/go-kit/refx"
)

func NewMapGetter(m map[string]interface{}) MapGetter {
	var v interface{}
	for key, val := range m {
		if err := refx.InterfaceSet(&v, key, val); err != nil {
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
