## Feature

这个库提供一些任意类型转换的方法，这些方法被广泛应用在 flag/config 里面

1. `To<T>E` 任意类型转换，转化失败抛出错误
2. `To<T>D` 任意类型转换，转换失败返回默认值
3. `To<T>P` 任意类型转换，转换失败 panic
4. `To<T>`  任意类型转换，转换失败返回零值
5. `SetInterface` 任意类型设置

## Quick Start

```go
So(cast.ToString(123), ShouldEqual, "123")
So(cast.ToInt("123"), ShouldEqual, 123)
So(cast.ToIntSlice("1,2,3"), ShouldResemble, []int{1, 2, 3})
```