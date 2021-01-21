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

	astx.WrapperGeneratorOptions

	Output string `flag:"output path"`
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
  gen --goPath vendor --pkgPath "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore" --package tablestore --classes TableStoreClient --classPrefix OTS
  gen --goPath vendor --pkgPath "github.com/olivere/elastic/v7" --package elastic --classes Client --classPrefix ES
  gen --goPath vendor --pkgPath "github.com/aliyun/aliyun-oss-go-sdk/oss" --package oss --classes Client --classPrefix OSS
`)
	}

	generator := astx.NewWrapperGeneratorWithOptions(&options.WrapperGeneratorOptions)
	str, err := generator.Generate()
	if err != nil {
		strx.Warn(err.Error())
		os.Exit(1)
	}

	if options.Output == "" {
		options.Output = fmt.Sprintf("%s.go", options.Package)
	}
	if options.Output != "stdout" {
		_ = ioutil.WriteFile(options.Output, []byte(str), 0644)
	} else {
		fmt.Println(str)
	}
}
