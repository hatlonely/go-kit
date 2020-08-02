package kv

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKV(t *testing.T) {
	Convey("TestKV case 1", t, func() {
		store := NewMemStore(WithMemStoreMemoryByte(200*1024*1024), WithMemStoreExpiration(60*time.Second))
		cache := NewKV(store)

		var val1, val2 string
		So(cache.Set("key1", "val1"), ShouldBeNil)
		So(cache.Set("key2", ""), ShouldBeNil)
		So(cache.Get("key1", &val1), ShouldBeNil)
		So(cache.Get("key2", &val2), ShouldBeNil)
		So(val1, ShouldEqual, "val1")
		So(val2, ShouldEqual, "")
	})

	Convey("TestKV case 2", t, func() {
		store := NewMemStore(WithMemStoreMemoryByte(200*1024*1024), WithMemStoreExpiration(60*time.Second))
		cache := NewKV(store, WithLoadFunc(func(key interface{}) (interface{}, error) {
			val, ok := map[string]string{
				"key1": "val1",
				"key2": "",
			}[key.(string)]
			if !ok {
				return nil, nil
			}
			return &val, nil
		}), WithStringVal(), WithStringKey())

		var val1, val2, val3 string
		So(cache.Get("key1", &val1), ShouldBeNil)
		So(cache.Get("key2", &val2), ShouldEqual, ErrNotFound)
		So(cache.Get("key3", &val3), ShouldEqual, ErrNotFound)
		So(val1, ShouldEqual, "val1")
		So(val2, ShouldEqual, "")
	})
}
