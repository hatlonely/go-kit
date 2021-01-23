package astx

import (
	"fmt"
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/strx"
)

func TestParseFunction(t *testing.T) {
	fs, err := ParseFunction("../vendor/github.com/aliyun/aliyun-tablestore-go-sdk/tablestore", "tablestore")
	fmt.Println(err)
	fmt.Println(strx.JsonMarshalIndent(fs))
}

func TestRegex(t *testing.T) {
	Convey("case 1", t, func() {
		re := regexp.MustCompile(`(\.*)\bDB\b`)
		f := func(s string) string {
			fmt.Println(s)
			if s[0] == '.' && s[1] != '.' {
				return s
			}
			return s[:len(s)-2] + "gorm.DB"
		}
		So(re.ReplaceAllStringFunc(`sql.DB`, f), ShouldEqual, `sql.DB`)
		So(re.ReplaceAllStringFunc(`XDB`, f), ShouldEqual, `XDB`)
		So(re.ReplaceAllStringFunc(`DB`, f), ShouldEqual, `gorm.DB`)
		So(re.ReplaceAllStringFunc(`...DB`, f), ShouldEqual, `...gorm.DB`)
	})
}
