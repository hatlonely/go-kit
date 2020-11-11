package main

import (
	"encoding/json"
	"fmt"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/flag"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/strx"
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
	Key         string
	Val         string
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
  cfg --camel-name --in-base-file base_local.json -a get --key mysql
  cfg --camel-name --in-base-file base_local.json -a diff --out-base-file base_remote.json
  cfg --camel-name --in-base-file base_local.json -a put --out-base-file base_remote.json
  cfg --camel-name --in-base-file base_local.json -a diff --key mysql --val '{
	  "connMaxLifeTime": "60s",
	  "database": "testdb2",
	  "host": "127.0.0.1",
	  "maxIdleConns": 10,
	  "maxOpenConns": 20,
	  "password": "",
	  "port": 3306,
	  "username": "hatlonely"
	}'
  cfg --camel-name --in-base-file base_local.json -a set --key mysql --val '{
	  "connMaxLifeTime": "60s",
	  "database": "testdb2",
	  "host": "127.0.0.1",
	  "maxIdleConns": 10,
	  "maxOpenConns": 20,
	  "password": "",
	  "port": 3306,
	  "username": "hatlonely"
	}'`)
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
		val, ok := cfg.Get(options.Key)
		if !ok {
			fmt.Println("null")
		}
		fmt.Println(strx.JsonMarshalSortKeys(val))
	}

	if options.Action == "diff" {
		icfg, err := config.NewConfigWithOptions(&options.In)
		Must(err)
		ocfg, err := config.NewConfigWithOptions(&options.Out)
		Must(err)

		if options.Key != "" {
			var v interface{}
			if err := json.Unmarshal([]byte(options.Val), &v); err != nil {
				Must(ocfg.UnsafeSet(options.Key, options.Val))
			} else {
				Must(ocfg.UnsafeSet(options.Key, v))
			}
		}

		fmt.Println(icfg.Diff(ocfg, options.Key))
	}

	if options.Action == "set" {
		ocfg, err := config.NewConfigWithOptions(&options.Out)
		Must(err)
		if options.Key != "" {
			var v interface{}
			if err := json.Unmarshal([]byte(options.Val), &v); err != nil {
				Must(ocfg.UnsafeSet(options.Key, options.Val))
			} else {
				Must(ocfg.UnsafeSet(options.Key, v))
			}
		}
		Must(ocfg.Save())
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
