// @see https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
// @see https://en.wikipedia.org/wiki/ANSI_escape_code

package strx

import (
	"bytes"
	"strconv"
	"strings"
)

func Render(str string, effects ...Effect) string {
	buf := &bytes.Buffer{}
	lines := strings.Split(str, "\n")
	for i, line := range strings.Split(str, "\n") {
		buf.WriteString(RenderLine(line, effects...))
		if i < len(lines)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

func RenderLine(str string, effects ...Effect) string {
	if len(effects) == 0 {
		return str
	}
	buf := &bytes.Buffer{}
	buf.WriteString("\033")
	buf.WriteString("[")
	for i := 0; i < len(effects)-1; i++ {
		buf.WriteString(strconv.Itoa(int(effects[i])))
		buf.WriteString(";")
	}
	buf.WriteString(strconv.Itoa(int(effects[len(effects)-1])))
	buf.WriteString("m")
	buf.WriteString(str)
	buf.WriteString("\033[0m")
	return buf.String()
}

type Effect int

const (
	FormatSetClose     Effect = 0
	FormatSetBold      Effect = 1
	FormatSetDim       Effect = 2
	FormatSetUnderline Effect = 4
	FormatSetBlink     Effect = 5
	FormatSetReverse   Effect = 7
	FormatSetHidden    Effect = 8
)

const (
	FormatResetAll       Effect = 0
	FormatResetBold      Effect = 21
	FormatResetDim       Effect = 22
	FormatResetUnderline Effect = 24
	FormatResetBlink     Effect = 25
	FormatResetReverse   Effect = 27
	FormatResetHidden    Effect = 28
)

const (
	ForegroundDefault      Effect = 39
	ForegroundBlack        Effect = 30
	ForegroundRed          Effect = 31
	ForegroundGreen        Effect = 32
	ForegroundYellow       Effect = 33
	ForegroundBlue         Effect = 34
	ForegroundMagenta      Effect = 35
	ForegroundCyan         Effect = 36
	ForegroundLightGray    Effect = 37
	ForegroundDarkGray     Effect = 90
	ForegroundLightRed     Effect = 91
	ForegroundLightGreen   Effect = 92
	ForegroundLightYellow  Effect = 93
	ForegroundLightBlue    Effect = 94
	ForegroundLightMagenta Effect = 95
	ForegroundLightCyan    Effect = 96
	ForegroundWhite        Effect = 97
)

const (
	BackgroundDefault      Effect = 49
	BackgroundBlack        Effect = 40
	BackgroundRed          Effect = 41
	BackgroundGreen        Effect = 42
	BackgroundYellow       Effect = 43
	BackgroundBlue         Effect = 44
	BackgroundMagenta      Effect = 45
	BackgroundCyan         Effect = 46
	BackgroundLightGray    Effect = 47
	BackgroundDarkGray     Effect = 100
	BackgroundLightRed     Effect = 101
	BackgroundLightGreen   Effect = 102
	BackgroundLightYellow  Effect = 103
	BackgroundLightBlue    Effect = 104
	BackgroundLightMagenta Effect = 105
	BackgroundLightCyan    Effect = 106
	BackgroundWhite        Effect = 107
)
