## Feature

基于 `github.com/PaesslerAG/gval` 实现的 Validator 库，设计初衷是为了解决 `gopkg.in/go-playground/validator.v9` tag 语法晦涩难懂的问题

- 目前性能较差，和 `validator.v9` 有几倍差距
- 支持正则表达式
- 语法自然
- 支持更复杂的条件表达式

## Quick Start

```go
type A struct {
    Key1 string `rule:"x in ['world', 'hello']"`
    Key2 int    `rule:"x>=5 && x<=6"`
    Key3 string `rule:"x =~ '^[0-9]{6}$'"`
    Key4 string `rule:"isEmail(x)"`
    Key5 int64  `rule:"x in [0, 1, 2]"`
    key6 int
    Key7 struct {
        Key8 int `rule:"x in [5, 7]"`
    }
}

obj := &A{
    Key1: "hello",
    Key2: 5,
    Key3: "123456",
    Key4: "hatlonely@foxmail.com",
    Key5: 1,
}
obj.Key7.Key8 = 7


So(validator.Validate(obj), ShouldNotBeNil)
```
