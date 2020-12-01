## Feature

1. 提供不同风格标识符转换函数
2. 提供字符串的 diff 功能
3. 提供 Json 转字符串函数
4. 简单字符串断言
5. 常见正则表达式
6. 终端颜色输出

## Quick Start

```go
func Info(format string, args ...interface{})
func Warn(format string, args ...interface{})
func Trac(format string, args ...interface{})

func Render(str string, effects ...Effect) string
func RenderLine(str string, effects ...Effect) string

func Diff(text1, text2 string) string

func JsonMarshal(v interface{}) string
func JsonMarshalIndent(v interface{}) string
func JsonMarshalSortKeys(v interface{}) string
func JsonMarshalIndentSortKeys(v interface{}) string
func FormatSpace(str string) string

func ToLower(ch uint8) uint8
func ToUpper(ch uint8) uint8
func CamelName(str string) string
func PascalName(str string) string
func SnakeName(str string) string
func KebabName(str string) string
func SnakeNameAllCaps(str string) string
func KebabNameAllCaps(str string) string

func IsUpper(ch uint8) bool
func IsLower(ch uint8) bool
func IsDigit(ch uint8) bool
func IsXdigit(ch uint8) bool
func IsAlpha(ch uint8) bool
func IsAlnum(ch uint8) bool
func All(str string, op func(uint8) bool) bool
func Any(str string, op func(uint8) bool) bool
```