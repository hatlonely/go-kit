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
	Version    bool   `flag:"-v; usage: show version"`
	SubCommand string `flag:"sub; usage: sub command" rule:"x in ['', 'wrap']"`
	Output     string `flag:"usage: output path"`
}

var Version string

func main() {
	var options Options
	refx.Must(flag.Struct(&options, refx.WithCamelName()))
	refx.Must(flag.Parse(flag.WithJsonVal()))
	refx.Must(bind.Bind(&options, []bind.Getter{flag.Instance()}, refx.WithCamelName(), refx.WithDefaultValidator()))

	if options.Version {
		strx.Trac(Version)
		return
	}

	if options.Help {
		switch options.SubCommand {
		case "wrap":
			f := flag.NewFlag("gen wrap")
			var options astx.WrapperGeneratorOptions
			refx.Must(f.Struct(&options, refx.WithCamelName()))
			strx.Trac(f.Usage())
			strx.Trac(`Example:
  gen wrap --sourcePath vendor \
      --packagePath "github.com/aliyun/aliyun-tablestore-go-sdk/v5/tablestore" \
      --packageName tablestore \
      --classPrefix OTS \
      --starClasses TableStoreClient \
      --rule.mainClass.include "^TableStoreClient$$" \
      --output $@
  gen wrap --sourcePath vendor \
      --packagePath "github.com/go-redis/redis" \
      --packageName redis \
      --classPrefix Redis \
      --errorField "Err()" \
      --inherit '{"Client": ["cmdable", "baseClient"], "ClusterClient": ["cmdable"]}' \
      --starClasses Client,ClusterClient \
      --rule.mainClass.include "^(?i:(Client)|(ClusterClient))$$" \
      --rule.errorInResult.include "^\*redis\..*Cmd$$" \
      --output $@
  gen wrap --sourcePath vendor \
      --packagePath "github.com/olivere/elastic/v7" \
      --packageName elastic \
      --classPrefix ES \
      --rule.starClass '{"include": "^(?i:(Client)|(.*Service))$$", "exclude": ".*"}' \
      --rule.mainClass.include "^Client$$" \
      --rule.wrap '{"default": {"exclude": ".*", "include": "^(Do)|(DoAsync)$$"}, "Client": {"exclude": ".*"}}' \
      --output $@
`)
		default:
			strx.Trac(flag.Usage())
		}
		return
	}

	if options.SubCommand == "wrap" {
		f := flag.NewFlag("gen wrap")
		var subOptions astx.WrapperGeneratorOptions
		refx.Must(f.Struct(&subOptions, refx.WithCamelName()))
		refx.Must(f.ParseArgs(os.Args[2:], flag.WithJsonVal()))
		refx.Must(bind.Bind(&subOptions, []bind.Getter{f}, refx.WithCamelName(), refx.WithDefaultValidator()))

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
