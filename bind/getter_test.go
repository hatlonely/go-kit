package bind

import (
	"os"
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMapGetter(t *testing.T) {
	Convey("TestMapGetter", t, func() {
		g := NewMapGetter(map[string]interface{}{
			"key1.key2.key3": "val1",
		})
		{
			val, ok := g.Get("key1.key2.key3")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, "val1")
		}
		{
			_, ok := g.Get("key")
			So(ok, ShouldBeFalse)
		}
	})
}

func TestEnvGetter(t *testing.T) {
	Convey("TestEnvGetter", t, func() {
		_ = os.Setenv("PREFIX_KEY1_KEY2_KEY3", "val1")
		g := NewEnvGetter(WithEnvPrefix("PREFIX"), WithEnvSeparator("_"))

		{
			val, ok := g.Get("key1.key2.key3")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, "val1")
		}
		{
			val, ok := g.Get("key1.key2Key3")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, "val1")
		}
		{
			val, ok := g.Get("Key1.Key2.Key3")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, "val1")
		}
		{
			_, ok := g.Get("Key1.Key2key3")
			So(ok, ShouldBeFalse)
		}
	})
}

func TestGinCtxGetter(t *testing.T) {
	Convey("TestGinCtxGetter", t, func() {
		patches := ApplyMethod(reflect.TypeOf(&gin.Context{}), "GetQuery", func(ctx *gin.Context, key string) (string, bool) {
			val, ok := map[string]string{
				"key1.key2.key3": "val1",
			}[key]
			return val, ok
		}).ApplyMethod(reflect.TypeOf(&gin.Context{}), "Get", func(ctx *gin.Context, key string) (value interface{}, exists bool) {
			val, ok := map[string]string{
				"key1.key2.key4": "val2",
			}[key]
			return val, ok
		})
		defer patches.Reset()
		g := NewGinCtxGetter(&gin.Context{})
		{
			val, ok := g.Get("key1.key2.key3")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, "val1")
		}
		{
			val, ok := g.Get("key1.key2.key4")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, "val2")
		}
		{
			_, ok := g.Get("key1.key2.key5")
			So(ok, ShouldBeFalse)
		}
	})
}
