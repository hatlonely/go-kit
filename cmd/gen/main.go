package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hatlonely/go-kit/astx"
	"github.com/hatlonely/go-kit/flag"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/strx"
)

type Options struct {
	flag.Options

	GoPath  string
	PkgPath string
	Package string
	Class   string
	Output  string
}

var Version string

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var options Options
	Must(flag.Struct(&options, refx.WithCamelName()))
	Must(flag.Parse(flag.WithJsonVal()))

	if options.Help {
		strx.Trac(flag.Usage())
		strx.Trac(`
  gen --goPath vendor --pkgPath "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore" --package tablestore --class TableStoreClient
`)
	}

	str, err := astx.GenerateWrapper(options.GoPath, options.PkgPath, options.Package, options.Class)
	if err != nil {
		strx.Warn(err.Error())
		os.Exit(1)
	}

	if options.Output == "" {
		options.Output = fmt.Sprintf("%s_%s", options.Package, options.Class)
	}
	if options.Output != "stdout" {
		_ = ioutil.WriteFile(options.Output, []byte(str), 0644)
	} else {
		fmt.Println(str)
	}
}
