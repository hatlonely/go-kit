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
	Help       bool   `flag:"-h; usage: show help info"`
	Version    string `flag:"-v; usage: show version"`
	SubCommand string `flag:"sub; usage: sub command" rule:"x in ['wrap']"`
	Output     string `flag:"usage: output path"`
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
	Must(bind.Bind(&options, []bind.Getter{flag.Instance()}, refx.WithCamelName(), refx.WithDefaultValidator()))

	if options.Help {
		switch options.SubCommand {
		case "wrap":
			f := flag.NewFlag("gen wrap")
			var options astx.WrapperGeneratorOptions
			Must(f.Struct(&options, refx.WithCamelName()))
			strx.Trac(f.Usage())
			strx.Trac(`
    gen wrap --sourcePath vendor \
        --packagePath "go.mongodb.org/mongo-driver/mongo" \
        --packageName mongo \
        --classPrefix Mongo \
        --rule.starClass '{"include": "^(?i:(Client)|(Database)|(Collection))$$", "exclude": ".*"}' \
        --rule.createMetric.include "^Client$$" \
        --rule.onWrapperChange.include "^Client$$" \
        --rule.onRetryChange.include "^Client$$" \
        --rule.trace '{"Client": {"exclude": "^Database$$"}, "Database": {"exclude": "^Collection$$"}}' \
        --rule.metric '{"Client": {"exclude": "^Database$$"}, "Database": {"exclude": "^Collection$$"}}' \
        --output wrap/autogen_mongo.go
`)
		default:
			strx.Trac(flag.Usage())
		}
		return
	}

	if options.SubCommand == "wrap" {
		f := flag.NewFlag("gen wrap")
		var subOptions astx.WrapperGeneratorOptions
		Must(f.Struct(&subOptions, refx.WithCamelName()))
		Must(f.ParseArgs(os.Args[2:], flag.WithJsonVal()))
		Must(bind.Bind(&subOptions, []bind.Getter{f}, refx.WithCamelName(), refx.WithDefaultValidator()))

		generator := astx.NewWrapperGeneratorWithOptions(&subOptions)
		str, err := generator.Generate()
		if err != nil {
			strx.Warn(err.Error())
			os.Exit(1)
		}

		if options.Output == "" {
			options.Output = fmt.Sprintf("autogen_%s.go", subOptions.PackageName)
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

}
