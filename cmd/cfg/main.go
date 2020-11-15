package main

import (
	"encoding/json"
	"fmt"
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

	Action        string   `flag:"-a; usage: actions, one of [get/set/put/diff/rollback]"`
	CamelName     bool     `flag:"usage: base file key format, for example [redisExpiration]"`
	SnakeName     bool     `flag:"usage: base file key format, for example [redis_expiration]"`
	KebabName     bool     `flag:"usage: base file key format, for example [redis-expiration]"`
	PascalName    bool     `flag:"usage: base file key format: for example [RedisExpiration]"`
	Key           string   `flag:"usage: key for set or diff"`
	Val           string   `flag:"usage: val for set or diff, val will auto convert to json"`
	RawVal        string   `flag:"usage: raw string value"`
	SetCipherKeys []string `flag:"usage: set cipher keys when put"`
	AddCipherKeys []string `flag:"usage: add cipher keys when put"`
	NoCipher      bool     `flag:"usage: decrypt all keys when put"`
	InBaseFile    string   `flag:"usage: base file name; default: base.json"`
	OutBaseFile   string   `flag:"usage: put/set target config, it will use in-base-file if not set"`
	BackupFile    string   `flag:"usage: file name to backup or rollback"`
}

func main() {
	var options Options
	Must(flag.Struct(&options, refx.WithCamelName()))
	Must(flag.Parse())
	if options.Help || options.Action == "" {
		Trac(flag.Usage())
		Trac(`examples:
  cfg --camelName --inBaseFile base.json -a get
  cfg --camelName --inBaseFile base.json -a get --key mysql
  cfg --camelName --inBaseFile base.json -a get --key mysql.password
  cfg --camelName --inBaseFile base.json -a diff --key mysql.username --val hatlonely
  cfg --camelName --inBaseFile base.json -a diff --key mysql.password --rawVal 12345678
  cfg --camelName --inBaseFile base.json -a diff --key mysql --val '{
      "connMaxLifeTime": "60s",
      "database": "testdb2",
      "host": "127.0.0.1",
      "maxIdleConns": 10,
      "maxOpenConns": 20,
      "password": "",
      "port": 3306,
      "username": "hatlonely"
  }'
  cfg --camelName --inBaseFile base.json -a set --key mysql.username --val hatlonely
  cfg --camelName --inBaseFile base.json -a set --key mysql.password --rawVal 12345678
  cfg --camelName --inBaseFile base.json -a set --key mysql --val '{
      "connMaxLifeTime": "60s",
      "database": "testdb2",
      "host": "127.0.0.1",
      "maxIdleConns": 10,
      "maxOpenConns": 20,
      "password": "",
      "port": 3306,
      "username": "hatlonely"
  }'
  cfg --camelName --inBaseFile base.json -a rollback --backupFile cfg.backup.json.20201113.224234
  cfg --camelName --inBaseFile base.json -a put --setCipherKeys mysql.password,redis.password
  cfg --camelName --inBaseFile base.json -a put --addCipherKeys email.password
  cfg --camelName --inBaseFile base.json -a put --noCipher
  cfg --camelName --inBaseFile base.json -a diff --outBaseFile base_ots.json
  cfg --camelName --inBaseFile base.json -a put --outBaseFile base_ots.json`)
		return
	}
	if options.Version {
		Trac(Version)
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

	if options.BackupFile == "" {
		options.BackupFile = fmt.Sprintf("cfg.backup.json.%v", time.Now().Format("20060102.150405"))
	}

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
			if options.RawVal != "" {
				Must(icfg.UnsafeSet(options.Key, options.RawVal))
			} else {
				var v interface{}
				if err := json.Unmarshal([]byte(options.Val), &v); err != nil {
					Must(icfg.UnsafeSet(options.Key, options.Val))
				} else {
					Must(icfg.UnsafeSet(options.Key, v))
				}
			}
		}

		fmt.Println(ocfg.Diff(icfg, options.Key))

		return
	}

	if options.Action == "set" {
		ocfg, err := config.NewConfigWithOptions(&outOptions)
		Must(err)

		BackUpCurrentConfig(ocfg, options.BackupFile)

		if options.Key == "" {
			Warn("[key] is required in set action")
			return
		}

		if options.RawVal != "" {
			Must(ocfg.UnsafeSet(options.Key, options.RawVal))
		} else {
			var v interface{}
			if err := json.Unmarshal([]byte(options.Val), &v); err != nil {
				Must(ocfg.UnsafeSet(options.Key, options.Val))
			} else {
				Must(ocfg.UnsafeSet(options.Key, v))
			}
		}
		Must(ocfg.Save())

		Trac("Save success. Use follow command to rollback:")
		Info(RollbackCommand(&options))

		return
	}

	if options.Action == "put" {
		icfg, err := config.NewConfigWithOptions(&inOptions)
		Must(err)
		ocfg, err := config.NewConfigWithOptions(&outOptions)
		Must(err)
		BackUpCurrentConfig(ocfg, options.BackupFile)

		ocfg, err = icfg.TransformWithOptions(&outOptions, &config.TransformOptions{
			CipherKeysToSet: options.SetCipherKeys,
			CipherKeysToAdd: options.AddCipherKeys,
			NoCipher:        options.NoCipher,
		})
		Must(err)
		Must(ocfg.Save())

		Trac("Save success. Use follow command to rollback:")
		Info(RollbackCommand(&options))

		return
	}

	if options.Action == "rollback" {
		icfg, err := config.NewConfigWithSimpleFile(options.BackupFile)
		Must(err)
		ocfg, err := icfg.Transform(&outOptions)
		Must(err)
		Must(ocfg.Save())

		Info("Rollback success")

		return
	}

	Warn("Unknown action %v", options.Action)
	os.Exit(1)
}

func RollbackCommand(options *Options) string {
	if options.CamelName {
		return fmt.Sprintf("cfg --camelName --inBaseFile %v -a rollback --backupFile %v", options.OutBaseFile, options.BackupFile)
	}
	if options.PascalName {
		return fmt.Sprintf("cfg --pascalName --inBaseFile %v -a rollback --backupFile %v", options.OutBaseFile, options.BackupFile)
	}
	if options.KebabName {
		return fmt.Sprintf("cfg --kebabName --inBaseFile %v -a rollback --backupFile %v", options.OutBaseFile, options.BackupFile)
	}
	if options.SnakeName {
		return fmt.Sprintf("cfg --snakeName --inBaseFile %v -a rollback --backupFile %v", options.OutBaseFile, options.BackupFile)
	}
	return fmt.Sprintf("cfg --inBaseFile %v -a rollback --backupFile %v", options.OutBaseFile, options.BackupFile)
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

func Info(format string, args ...interface{}) {
	fmt.Println(strx.Render(fmt.Sprintf(format, args...), strx.FormatSetBold, strx.ForegroundGreen))
}

func Warn(format string, args ...interface{}) {
	fmt.Println(strx.Render(fmt.Sprintf(format, args...), strx.FormatSetBold, strx.ForegroundRed))
}

func Trac(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println()
}
