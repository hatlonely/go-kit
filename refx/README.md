## Feature

1. 封装一些利用反射操作 map/slice 嵌套结构的方法，是 config/flag 的底层存储
2. 通过 `dft` 支持结构体默认值

## Quick Start

```go
func SetDefaultValue(v interface{}) error

func InterfaceToStruct(src interface{}, dst interface{}, opts ...Option) error
func InterfaceDiff(v1 interface{}, v2 interface{}) ([]string, error)
func InterfaceGet(v interface{}, key string) (interface{}, error)
func InterfaceSet(pv *interface{}, key string, val interface{}) error
func InterfaceDel(pv *interface{}, key string) error
func InterfaceTravel(v interface{}, fun func(key string, val interface{}) error) error
```