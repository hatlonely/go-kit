package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

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

	Action      string   `flag:"--action,-a; usage: actions, one of [get/set/put/diff/rollback]"`
	CamelName   bool     `flag:"usage: base file key format, for example [redisExpiration]"`
	SnakeName   bool     `flag:"usage: base file key format, for example [redis_expiration]"`
	KebabName   bool     `flag:"usage: base file key format, for example [redis-expiration]"`
	PascalName  bool     `flag:"usage: base file key format: for example [RedisExpiration]"`
	Key         string   `flag:"usage: key for set or diff"`
	Val         string   `flag:"usage: val for set or diff"`
	CipherKeys  []string `flag:"usage: change cipher keys when put"`
	NoCipher    bool     `flag:"usage: decrypt all keys when put"`
	InBaseFile  string   `flag:"usage: base file name"`
	OutBaseFile string   `flag:"usage: put/set target config, it will use in-base-file if not set"`
	BackupFile  string   `flag:"usage: file name to backup or rollback; default: cfg.backup.json"`
}

func main() {
	var options Options
	Must(flag.Struct(&options))
	Must(flag.Parse())
	if options.Help {
		fmt.Println(flag.Usage())
		fmt.Println(`examples:
  cfg --camel-name --in-base-file base.json -a get --key mysql
  cfg --camel-name --in-base-file base.json -a diff --out-base-file base_remote.json
  cfg --camel-name --in-base-file base.json -a put --out-base-file base_remote.json
  cfg --camel-name --in-base-file base.json -a diff --key mysql --val '{
    "connMaxLifeTime": "60s",
    "database": "testdb2",
    "host": "127.0.0.1",
    "maxIdleConns": 10,
    "maxOpenConns": 20,
    "password": "",
    "port": 3306,
    "username": "hatlonely"
  }'
  cfg --camel-name --in-base-file base.json -a set --key mysql --val '{
    "connMaxLifeTime": "60s",
    "database": "testdb2",
    "host": "127.0.0.1",
    "maxIdleConns": 10,
    "maxOpenConns": 20,
    "password": "",
    "port": 3306,
    "username": "hatlonely"
  }'
  cfg --camel-name --in-base-file base.json -a rollback --backup-file cfg.backup.json`)
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
	options.BackupFile = fmt.Sprintf("%v.%v", options.BackupFile, time.Now().Format("20060102.150405"))

	var inOptions config.Options
	var outOptions config.Options
	if options.InBaseFile != "" {
		cfg, err := config.NewConfigWithSimpleFile(options.InBaseFile)
		Must(err)
		Must(cfg.Unmarshal(&inOptions, opts...))
		if options.OutBaseFile == "" {
			options.OutBaseFile = options.InBaseFile
		}
	}
	if options.OutBaseFile != "" {
		cfg, err := config.NewConfigWithSimpleFile(options.OutBaseFile)
		Must(err)
		Must(cfg.Unmarshal(&outOptions, opts...))
	}

	if options.Action == "get" {
		cfg, err := config.NewConfigWithOptions(&inOptions)
		Must(err)
		cfg, err = cfg.TransformWithOptions(&inOptions, &config.TransformOptions{
			CipherKeys: options.CipherKeys,
			NoCipher:   options.NoCipher,
		})
		Must(err)
		val, ok := cfg.Get(options.Key)
		if !ok {
			fmt.Println("null")
		}
		fmt.Println(strx.JsonMarshalSortKeys(val))

		return
	}

	if options.Action == "diff" {
		icfg, err := config.NewConfigWithOptions(&inOptions)
		Must(err)
		ocfg, err := config.NewConfigWithOptions(&outOptions)
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

		return
	}

	if options.Action == "set" {
		ocfg, err := config.NewConfigWithOptions(&outOptions)
		Must(err)

		BackUpCurrentConfig(ocfg, options.BackupFile)

		if options.Key != "" {
			var v interface{}
			if err := json.Unmarshal([]byte(options.Val), &v); err != nil {
				Must(ocfg.UnsafeSet(options.Key, options.Val))
			} else {
				Must(ocfg.UnsafeSet(options.Key, v))
			}
		}
		Must(ocfg.Save())

		fmt.Println(strx.Render(fmt.Sprintf("Save success. Use follow command to rollback:")))
		fmt.Println(strx.Render(RollbackCommand(&options), strx.FormatSetBold, strx.ForegroundGreen))

		return
	}

	if options.Action == "put" {
		icfg, err := config.NewConfigWithOptions(&inOptions)
		Must(err)
		ocfg, err := config.NewConfigWithOptions(&outOptions)
		Must(err)
		BackUpCurrentConfig(ocfg, options.BackupFile)

		ocfg, err = icfg.TransformWithOptions(&outOptions, &config.TransformOptions{
			CipherKeys: options.CipherKeys,
			NoCipher:   options.NoCipher,
		})
		Must(err)
		Must(ocfg.Save())

		fmt.Println(strx.Render(fmt.Sprintf("Save success. Use follow command to rollback:")))
		fmt.Println(strx.Render(RollbackCommand(&options), strx.FormatSetBold, strx.ForegroundGreen))

		return
	}

	if options.Action == "rollback" {
		icfg, err := config.NewConfigWithSimpleFile(options.BackupFile)
		Must(err)
		ocfg, err := icfg.Transform(&outOptions)
		Must(err)
		Must(ocfg.Save())

		fmt.Println(strx.Render("Rollback success", strx.FormatSetBold, strx.ForegroundGreen))
		buf, _ := ioutil.ReadFile(options.OutBaseFile)
		fmt.Println(string(buf))

		return
	}

	fmt.Println(strx.Render(strx.Render(fmt.Sprintf("Unknown action %v", options.Action), strx.FormatSetBold, strx.ForegroundRed)))
	fmt.Println(flag.Usage())
	os.Exit(1)
}

func RollbackCommand(options *Options) string {
	if options.CamelName {
		return fmt.Sprintf("cfg --camel-name --in-base-file %v -a rollback --backup-file %v", options.OutBaseFile, options.BackupFile)
	}
	if options.PascalName {
		return fmt.Sprintf("cfg --pascal-name --in-base-file %v -a rollback --backup-file %v", options.OutBaseFile, options.BackupFile)
	}
	if options.KebabName {
		return fmt.Sprintf("cfg --kebab-name --in-base-file %v -a rollback --backup-file %v", options.OutBaseFile, options.BackupFile)
	}
	if options.SnakeName {
		return fmt.Sprintf("cfg --snake-name --in-base-file %v -a rollback --backup-file %v", options.OutBaseFile, options.BackupFile)
	}
	return fmt.Sprintf("cfg --in-base-file %v -a rollback --backup-file %v", options.OutBaseFile, options.BackupFile)
}

func BackUpCurrentConfig(cfg *config.Config, name string) {
	cfg, err := cfg.Transform(&config.Options{
		Provider: config.ProviderOptions{
			Type: "Local",
			LocalProvider: config.LocalProviderOptions{
				Filename: name,
			},
		},
	})
	Must(err)
	Must(cfg.Save())
}
