package astx

import (
	"fmt"
	"testing"

	"github.com/hatlonely/go-kit/strx"
)

func TestParseFunction(t *testing.T) {
	fs, err := ParseFunction("../vendor/github.com/aliyun/aliyun-tablestore-go-sdk/tablestore", "tablestore")
	fmt.Println(err)
	fmt.Println(strx.JsonMarshalIndent(fs))
}

func TestGenerateWrapper(t *testing.T) {
	code, err := GenerateWrapper("../vendor", "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore", "tablestore", "TableStoreClient")
	fmt.Println(err)
	fmt.Println(code)
}
