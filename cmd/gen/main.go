package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hatlonely/go-kit/astx"
	"github.com/hatlonely/go-kit/bind"
	"github.com/hatlonely/go-kit/flag"
	"github.com/hatlonely/go-kit/ops"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/strx"
)

type Options struct {
	Help    bool   `flag:"-h; usage: show help info"`
	Version string `flag:"-v; usage: show version"`
	Output  string `flag:"usage: output path"`

	astx.WrapperGeneratorOptions
}

var Version string

func Must(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}

func main() {
	var options Options
	Must(flag.Struct(&options, refx.WithCamelName()))
	Must(flag.Parse(flag.WithJsonVal()))
	Must(bind.Bind(&options, []bind.Getter{flag.Instance()}, refx.WithCamelName()))

	if options.Help {
		strx.Trac(flag.Usage())
		strx.Trac(`
    gen --sourcePath vendor \
        --packagePath "go.mongodb.org/mongo-driver/mongo" \
        --packageName mongo \
        --classPrefix Mongo \
        --rule.starClass '{"include": "^(?i:(Client)|(Database)|(Collection))$$", "exclude": ".*"}' \
        --rule.createMetric.include "^Client$$" \
        --rule.onWrapperChange.include "^Client$$" \
        --rule.onRetryChange.include "^Client$$" \
        --rule.trace '{"Client": {"exclude": "^Database$$"}, "Database": {"exclude": "^Collection$$"}}' \
        --rule.metric '{"Client": {"exclude": "^Database$$"}, "Database": {"exclude": "^Collection$$"}}' \
        --output $@
`)
		return
	}

	generator := astx.NewWrapperGeneratorWithOptions(&options.WrapperGeneratorOptions)
	str, err := generator.Generate()
	if err != nil {
		strx.Warn(err.Error())
		os.Exit(1)
	}

	if options.Output == "" {
		options.Output = fmt.Sprintf("%s.go", options.PackageName)
	}
	if options.Output != "stdout" {
		_ = ioutil.WriteFile(options.Output, []byte(str), 0644)
		code, err := ops.ExecCommandWithOutput(fmt.Sprintf("goimports -w %v", options.Output), nil, ".", os.Stdout, os.Stderr)
		if code != 0 {
			strx.Warn(fmt.Sprintf("exit [%v]", code))
		}
		if err != nil {
			strx.Warn(err.Error())
		}
	} else {
		fmt.Println(str)
	}
}
