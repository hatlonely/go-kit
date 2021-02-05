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

var Version = "1.0.0"

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
	Val           string   `flag:"usage: val for set or diff"`
	JsonVal       string   `flag:"usage: json val to set or diff"`
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
	if options.Version {
		strx.Trac(Version)
		return
	}
	if options.Help || options.Action == "" {
		strx.Trac(flag.Usage())
		strx.Trac(`examples:
  cfg --camelName --inBaseFile base.json -a get
  cfg --camelName --inBaseFile base.json -a get --key mysql
  cfg --camelName --inBaseFile base.json -a get --key mysql.password
  cfg --camelName --inBaseFile base.json -a diff --key mysql.username --val hatlonely
  cfg --camelName --inBaseFile base.json -a diff --key mysql.password --val 12345678
  cfg --camelName --inBaseFile base.json -a diff --key mysql --jsonVal '{
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
  cfg --camelName --inBaseFile base.json -a set --key mysql.password --val 12345678
  cfg --camelName --inBaseFile base.json -a set --key mysql --jsonVal '{
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
			if options.JsonVal != "" {
				var v interface{}
				Must(json.Unmarshal([]byte(options.JsonVal), &v))
				Must(ocfg.UnsafeSet(options.Key, v))
			} else {
				Must(icfg.UnsafeSet(options.Key, options.Val))
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
			strx.Warn("[key] is required in set action")
			return
		}

		if options.JsonVal != "" {
			var v interface{}
			Must(json.Unmarshal([]byte(options.JsonVal), &v))
			Must(ocfg.UnsafeSet(options.Key, v))
		} else {
			Must(ocfg.UnsafeSet(options.Key, options.Val))
		}
		Must(ocfg.Save())

		strx.Trac("Save success. Use follow command to rollback:")
		strx.Info(RollbackCommand(&options))

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

		strx.Trac("Save success. Use follow command to rollback:")
		strx.Info(RollbackCommand(&options))

		return
	}

	if options.Action == "rollback" {
		icfg, err := config.NewConfigWithSimpleFile(options.BackupFile)
		Must(err)
		ocfg, err := icfg.Transform(&outOptions)
		Must(err)
		Must(ocfg.Save())

		strx.Info("Rollback success")

		return
	}

	strx.Warn("Unknown action %v", options.Action)
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
			Options: &config.LocalProviderOptions{
				Filename: name,
			},
		},
	})
	Must(err)
	Must(cfg.Save())
}
