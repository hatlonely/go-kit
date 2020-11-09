package main

import (
	"fmt"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/flag"
	"github.com/hatlonely/go-kit/refx"
)

var Version string

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

type Options struct {
	flag.Options

	Action      string `flag:"--action,-a; usage: actions"`
	CamelName   bool
	SnakeName   bool
	KebabName   bool
	PascalName  bool
	CipherKeys  []string
	NoCipher    bool
	InBaseFile  string
	OutBaseFile string
	In          config.Options
	Out         config.Options
}

func main() {
	var options Options
	Must(flag.Struct(&options))
	Must(flag.Parse())
	if options.Help {
		fmt.Println(flag.Usage())
		fmt.Println(`examples:
cfg -a get --in-base-file base_local.json
cfg -a diff --in-base-file base_local.json --out-base-file base_remote.json
cfg -a put --in-base-file base_local.json --out-base-file base_remote.json`)
		return
	}
	if options.Version {
		fmt.Println(Version)
		return
	}

	var opts []refx.Option
	if options.CamelName {
		opts = append(opts, refx.WithCamelName())
	}
	if options.PascalName {
		opts = append(opts, refx.WithPascalName())
	}
	if options.KebabName {
		opts = append(opts, refx.WithKebabName())
	}
	if options.SnakeName {
		opts = append(opts, refx.WithSnakeName())
	}

	if options.InBaseFile != "" {
		cfg, err := config.NewConfigWithSimpleFile(options.InBaseFile)
		Must(err)
		Must(cfg.Unmarshal(&options.In, opts...))
		options.Out = options.In
	}
	if options.OutBaseFile != "" {
		cfg, err := config.NewConfigWithSimpleFile(options.OutBaseFile)
		Must(err)
		Must(cfg.Unmarshal(&options.Out, opts...))
	}

	if options.Action == "get" {
		cfg, err := config.NewConfigWithOptions(&options.In)
		Must(err)
		cfg, err = cfg.TransformWithOptions(&options.In, &config.TransformOptions{
			CipherKeys: options.CipherKeys,
			NoCipher:   options.NoCipher,
		})
		Must(err)
		fmt.Println(cfg.ToString())
	}

	if options.Action == "diff" {
		icfg, err := config.NewConfigWithOptions(&options.In)
		Must(err)
		ocfg, err := config.NewConfigWithOptions(&options.Out)
		Must(err)
		icfg.Diff(ocfg)
	}

	if options.Action == "put" {
		icfg, err := config.NewConfigWithOptions(&options.In)
		Must(err)
		ocfg, err := icfg.TransformWithOptions(&options.Out, &config.TransformOptions{
			CipherKeys: options.CipherKeys,
			NoCipher:   options.NoCipher,
		})
		Must(err)
		Must(ocfg.Save())
	}
}
