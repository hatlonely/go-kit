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
	NoBackup      bool     `flag:"usage: not generate backup file"`
	InBaseFile    string   `flag:"usage: base file name"`
	OutBaseFile   string   `flag:"usage: put/set target config, it will use in-base-file if not set"`
	BackupFile    string   `flag:"usage: file name to backup or rollback"`
	InFile        string   `flag:"usage: input config file"`
	InFileType    string   `flag:"usage: input config file type; default: Json"`
	OutFile       string   `flag:"usage: output config file"`
	OutFileType   string   `flag:"usage: output config file type; default: Json"`
}

func main() {
	var options Options
	refx.Must(flag.Struct(&options, refx.WithCamelName()))
	refx.Must(flag.Parse())
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
		refx.Must(err)
		refx.Must(cfg.Unmarshal(&inOptions, opts...))
	}
	if options.OutBaseFile != "" {
		cfg, err := config.NewConfigWithSimpleFile(options.OutBaseFile)
		refx.Must(err)
		refx.Must(cfg.Unmarshal(&outOptions, opts...))
	}
	if options.InFile != "" {
		inOptions = config.Options{
			Decoder: config.DecoderOptions{
				Type: options.InFileType,
			},
			Provider: config.ProviderOptions{
				Type: "Local",
				Options: &config.LocalProviderOptions{
					Filename: options.InFile,
				},
			},
		}
	}
	if options.OutFile != "" {
		if _, err := os.Stat(options.OutFile); err != nil {
			if os.IsNotExist(err) {
				_, err := os.Create(options.OutFile)
				refx.Must(err)
			} else {
				refx.Must(err)
			}
		}
		outOptions = config.Options{
			Decoder: config.DecoderOptions{
				Type: options.OutFileType,
			},
			Provider: config.ProviderOptions{
				Type: "Local",
				Options: &config.LocalProviderOptions{
					Filename: options.OutFile,
				},
			},
		}
	}
	if options.OutBaseFile == "" && options.OutFile == "" {
		outOptions = inOptions
	}

	if options.Action == "get" {
		cfg, err := config.NewConfigWithOptions(&inOptions)
		refx.Must(err)
		val, ok := cfg.Get(options.Key)
		if !ok {
			fmt.Println("null")
		}
		fmt.Println(strx.JsonMarshalSortKeys(val))

		return
	}

	if options.Action == "diff" {
		icfg, err := config.NewConfigWithOptions(&inOptions)
		refx.Must(err)
		ocfg, err := config.NewConfigWithOptions(&outOptions)
		refx.Must(err)

		if options.Key != "" {
			if options.JsonVal != "" {
				var v interface{}
				refx.Must(json.Unmarshal([]byte(options.JsonVal), &v))
				refx.Must(ocfg.UnsafeSet(options.Key, v))
			} else {
				refx.Must(icfg.UnsafeSet(options.Key, options.Val))
			}
		}

		fmt.Println(ocfg.Diff(icfg, options.Key))

		return
	}

	if options.Action == "set" {
		ocfg, err := config.NewConfigWithOptions(&outOptions)
		refx.Must(err)

		if !options.NoBackup {
			BackUpCurrentConfig(ocfg, options.BackupFile)
		}

		if options.Key == "" {
			strx.Warn("[key] is required in set action")
			return
		}

		if options.JsonVal != "" {
			var v interface{}
			refx.Must(json.Unmarshal([]byte(options.JsonVal), &v))
			refx.Must(ocfg.UnsafeSet(options.Key, v))
		} else {
			refx.Must(ocfg.UnsafeSet(options.Key, options.Val))
		}
		refx.Must(ocfg.Save())

		if !options.NoBackup {
			strx.Trac("Save success. Use follow command to rollback:")
			strx.Info(RollbackCommand(&options))
		}

		return
	}

	if options.Action == "put" {
		icfg, err := config.NewConfigWithOptions(&inOptions)
		refx.Must(err)
		ocfg, err := config.NewConfigWithOptions(&outOptions)
		refx.Must(err)

		if !options.NoBackup {
			BackUpCurrentConfig(ocfg, options.BackupFile)
		}

		ocfg, err = icfg.TransformWithOptions(&outOptions, &config.TransformOptions{
			CipherKeysToSet: options.SetCipherKeys,
			CipherKeysToAdd: options.AddCipherKeys,
			NoCipher:        options.NoCipher,
		})
		refx.Must(err)
		refx.Must(ocfg.Save())

		if !options.NoBackup {
			strx.Trac("Save success. Use follow command to rollback:")
			strx.Info(RollbackCommand(&options))
		}

		return
	}

	if options.Action == "rollback" {
		icfg, err := config.NewConfigWithSimpleFile(options.BackupFile)
		refx.Must(err)
		ocfg, err := icfg.Transform(&outOptions)
		refx.Must(err)
		refx.Must(ocfg.Save())

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
	refx.Must(err)
	refx.Must(cfg.Save())
}
