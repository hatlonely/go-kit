package strx

import (
	"bytes"
	"fmt"
	"strings"
)

// 相关 diff 库比较
// github.com/sergi/go-diff/diffmatchpatch
//   显示不太直观
// github.com/yudai/gojsondiff
//   只支持 json，且 json 不能为空，必须为 map 或者 slice
//   如果 target 为空，显示不出 diff
// github.com/cj1128/myers-diff
//   只有工具

// 下面这个实现来自于： https://github.com/cj1128/myers-diff/blob/master/main.go
func Diff(text1, text2 string) string {
	return generateDiff(strings.Split(text1, "\n"), strings.Split(text2, "\n"))
}

type operation uint

const (
	INSERT operation = 1
	DELETE           = 2
	MOVE             = 3
)

func (op operation) String() string {
	switch op {
	case INSERT:
		return "INS"
	case DELETE:
		return "DEL"
	case MOVE:
		return "MOV"
	default:
		return "UNKNOWN"
	}
}

var colors = map[operation]string{
	INSERT: "\033[32m",
	DELETE: "\033[31m",
	MOVE:   "\033[39m",
}

func generateDiff(src, dst []string) string {
	script := shortestEditScript(src, dst)
	srcIndex, dstIndex := 0, 0
	var buf bytes.Buffer
	for _, op := range script {
		switch op {
		case INSERT:
			buf.WriteString(colors[op] + "+" + dst[dstIndex] + "\n")
			dstIndex += 1
		case MOVE:
			buf.WriteString(colors[op] + " " + src[srcIndex] + "\n")
			srcIndex += 1
			dstIndex += 1
		case DELETE:
			buf.WriteString(colors[op] + "-" + src[srcIndex] + "\n")
			srcIndex += 1
		}
	}
	return buf.String()
}

// 生成最短的编辑脚本
func shortestEditScript(src, dst []string) []operation {
	n := len(src)
	m := len(dst)
	max := n + m
	var trace []map[int]int
	var x, y int

loop:
	for d := 0; d <= max; d++ {
		// 最多只有 d+1 个 k
		v := make(map[int]int, d+2)
		trace = append(trace, v)

		// 需要注意处理对角线
		if d == 0 {
			t := 0
			for len(src) > t && len(dst) > t && src[t] == dst[t] {
				t++
			}
			v[0] = t
			if t == len(src) && t == len(dst) {
				break loop
			}
			continue
		}

		lastV := trace[d-1]

		for k := -d; k <= d; k += 2 {
			// 向下
			if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
				x = lastV[k+1]
			} else { // 向右
				x = lastV[k-1] + 1
			}

			y = x - k

			// 处理对角线
			for x < n && y < m && src[x] == dst[y] {
				x, y = x+1, y+1
			}

			v[k] = x

			if x == n && y == m {
				break loop
			}
		}
	}

	// 反向回溯
	var script []operation

	x = n
	y = m
	var k, prevK, prevX, prevY int

	for d := len(trace) - 1; d > 0; d-- {
		k = x - y
		lastV := trace[d-1]

		if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		prevX = lastV[prevK]
		prevY = prevX - prevK

		for x > prevX && y > prevY {
			script = append(script, MOVE)
			x -= 1
			y -= 1
		}

		if x == prevX {
			script = append(script, INSERT)
		} else {
			script = append(script, DELETE)
		}

		x, y = prevX, prevY
	}

	if trace[0][0] != 0 {
		for i := 0; i < trace[0][0]; i++ {
			script = append(script, MOVE)
		}
	}

	return reverse(script)
}

func printTrace(trace []map[int]int) {
	for d := 0; d < len(trace); d++ {
		fmt.Printf("d = %d:\n", d)
		v := trace[d]
		for k := -d; k <= d; k += 2 {
			x := v[k]
			y := x - k
			fmt.Printf("  k = %2d: (%d, %d)\n", k, x, y)
		}
	}
}

func reverse(s []operation) []operation {
	result := make([]operation, len(s))

	for i, v := range s {
		result[len(s)-1-i] = v
	}

	return result
}
