package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/hatlonely/go-kit/refx"
)

func CreateTestFile() {
	fp, _ := os.Create("test.json")
	_, _ = fp.WriteString(`{
  "http": {
    "port": 80
  },
  "grpc": {
    "port": 6080
  },
  "service": {
    "accountExpiration": "5m",
    "captchaExpiration": "30m"
  },
  "redis": {
    "addr": "127.0.0.1:6379",
    "dialTimeout": "200ms",
    "readTimeout": "200ms",
    "writeTimeout": "200ms",
    "maxRetries": 3,
    "poolSize": 20,
    "db": 0,
    "password": ""
  },
  "mysql": {
    "username": "root",
    "password": "",
    "database": "account",
    "host": "127.0.0.1",
    "port": 3306,
    "connMaxLifeTime": "60s",
    "maxIdleConns": 10,
    "maxOpenConns": 20
  },
  "email": {
    "from": "hatlonely@foxmail.com",
    "password": "123456",
    "server": "smtp.qq.com",
    "port": 25
  },
  "logger": {
    "grpc": {
      "level": "Info",
      "writers": [{
        "type": "RotateFile",
        "rotateFileWriter": {
          "filename": "log/account.rpc",
          "maxAge": "24h",
          "formatter": {
            "type": "Json"
          }
        }
      }]
    },
    "info": {
      "level": "Info",
      "writers": [{
        "type": "RotateFile",
        "rotateFileWriter": {
          "filename": "log/account.log",
          "maxAge": "48h",
          "formatter": {
            "type": "Json"
          }
        }
      }]
    }
  }
}`)
	_ = fp.Close()
}

func CreateBaseFile() {
	fp, _ := os.Create("base.json")
	_, _ = fp.WriteString(`{
  "decoder": {
    "type": "Json"
  },
  "provider": {
    "type": "Local",
	"localProvider": {
      "filename": "test.json"
    }
  },
  "cipher": {
    "type": "Group",
    "cipherGroup": [{
      "type": "AES",
      "aesCipher": {
        "base64Key": "IrjXy4vx7iwgCLaUeu5TVUA9TkgMwSw3QWcgE/IW5W0="
      }
    }, {
      "type": "Base64",
    }]
  }
}`)
}

func DeleteTestFile() {
	_ = os.Remove("test.json")
}

func DeleteBaseFile() {
	_ = os.Remove("base.json")
}

func TestConfig_Get(t *testing.T) {
	Convey("TestConfig_Get", t, func() {
		CreateTestFile()

		cfg, err := NewConfigWithSimpleFile("test.json")
		So(err, ShouldBeNil)
		So(cfg.GetString("http.port"), ShouldEqual, "80")
		So(cfg.GetInt("http.port"), ShouldEqual, 80)
		So(cfg.GetInt64("http.port"), ShouldEqual, 80)
		So(cfg.GetInt("grpc.port"), ShouldEqual, 6080)
		So(cfg.GetDuration("service.accountExpiration"), ShouldEqual, 5*time.Minute)
		So(cfg.GetDuration("service.captchaExpiration"), ShouldEqual, 30*time.Minute)
		So(cfg.GetString("service.accountExpiration"), ShouldEqual, "5m")
		So(cfg.GetString("service.captchaExpiration"), ShouldEqual, "30m")
		So(cfg.GetString("logger.grpc.level"), ShouldEqual, "Info")
		So(cfg.GetString("logger.grpc.writers[0].type"), ShouldEqual, "RotateFile")
		So(cfg.GetString("logger.grpc.writers[0].rotateFileWriter.filename"), ShouldEqual, "log/account.rpc")
		So(cfg.GetDuration("logger.grpc.writers[0].rotateFileWriter.maxAge"), ShouldEqual, 24*time.Hour)
		So(cfg.GetString("logger.grpc.writers[0].rotateFileWriter.formatter.type"), ShouldEqual, "Json")

		Convey("Get ok", func() {
			val, ok := cfg.Get("http.port")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, 80)
		})

		Convey("Get not ok", func() {
			val, ok := cfg.Get("xxx")
			So(ok, ShouldBeFalse)
			So(val, ShouldBeNil)
		})

		DeleteTestFile()
	})
}

func TestConfig_Unmarshal(t *testing.T) {
	type Options struct {
		Http struct {
			Port int
		}
		Grpc struct {
			Port int
		}

		Service struct {
			AccountExpiration time.Duration
			CaptchaExpiration time.Duration
		}

		Logger struct {
			Info struct {
				Level   string `rule:"x in ['Info', 'Warn', 'Debug', 'Error']"`
				Writers []struct {
					Type             string `dft:"Stdout" rule:"x in ['RotateFile', 'Stdout']"`
					RotateFileWriter struct {
						Filename  string        `bind:"required"`
						MaxAge    time.Duration `dft:"24h"`
						Formatter struct {
							Type string `dft:"Json" rule:"x in ['Json', 'Text']"`
						}
					}
				}
			}
		}
	}

	Convey("TestConfig_Unmarshal", t, func() {
		CreateTestFile()

		cfg, err := NewConfigWithSimpleFile("test.json")
		So(err, ShouldBeNil)

		var options Options
		So(cfg.Unmarshal(&options, refx.WithCamelName()), ShouldBeNil)

		So(options.Http.Port, ShouldEqual, 80)
		So(options.Grpc.Port, ShouldEqual, 6080)
		So(options.Grpc.Port, ShouldEqual, 6080)
		So(options.Service.AccountExpiration, ShouldEqual, 5*time.Minute)
		So(options.Service.CaptchaExpiration, ShouldEqual, 30*time.Minute)
		So(options.Logger.Info.Level, ShouldEqual, "Info")
		So(options.Logger.Info.Writers[0].Type, ShouldEqual, "RotateFile")
		So(options.Logger.Info.Writers[0].RotateFileWriter.Filename, ShouldEqual, "log/account.log")
		So(options.Logger.Info.Writers[0].RotateFileWriter.MaxAge, ShouldEqual, 48*time.Hour)
		So(options.Logger.Info.Writers[0].RotateFileWriter.Formatter.Type, ShouldEqual, "Json")

		DeleteTestFile()
	})
}

func TestConfig_Sub(t *testing.T) {
	Convey("TestConfig_Sub", t, func() {
		CreateTestFile()
		cfg, err := NewConfigWithSimpleFile("test.json")
		So(err, ShouldBeNil)

		Convey("sub", func() {
			cfg := cfg.Sub("logger")
			So(cfg.GetString("grpc.writers[0].type"), ShouldEqual, "RotateFile")
			So(cfg.GetString("grpc.writers[0].rotateFileWriter.filename"), ShouldEqual, "log/account.rpc")
			So(cfg.GetDuration("grpc.writers[0].rotateFileWriter.maxAge"), ShouldEqual, 24*time.Hour)
			So(cfg.GetString("grpc.writers[0].rotateFileWriter.formatter.type"), ShouldEqual, "Json")
		})

		Convey("sub.sub", func() {
			So(cfg.Sub("logger").Sub("grpc").Sub("writers").GetString("[0].type"), ShouldEqual, "RotateFile")
			So(cfg.Sub("logger").Sub("grpc").Sub("writers").Sub("[0]").GetString("type"), ShouldEqual, "RotateFile")
			So(cfg.Sub("logger").Sub("grpc").Sub("writers").Sub("[0]").Sub("type").GetString(""), ShouldEqual, "RotateFile")
			So(cfg.Sub("logger").Sub("grpc.writers").Sub("[0]").Sub("type").GetString(""), ShouldEqual, "RotateFile")
		})

		Convey("sub unmarshal", func() {
			type Options struct {
				Level   string `rule:"x in ['Info', 'Warn', 'Debug', 'Error']"`
				Writers []struct {
					Type             string `dft:"Stdout" rule:"x in ['RotateFile', 'Stdout']"`
					RotateFileWriter struct {
						Filename  string        `bind:"required"`
						MaxAge    time.Duration `dft:"24h"`
						Formatter struct {
							Type string `dft:"Json" rule:"x in ['Json', 'Text']"`
						}
					}
				}
			}
			var options Options
			cfg := cfg.Sub("logger.info")
			So(cfg.Unmarshal(&options, refx.WithCamelName()), ShouldBeNil)
			So(options.Writers[0].Type, ShouldEqual, "RotateFile")
			So(options.Writers[0].RotateFileWriter.Filename, ShouldEqual, "log/account.log")
			So(options.Writers[0].RotateFileWriter.MaxAge, ShouldEqual, 48*time.Hour)
			So(options.Writers[0].RotateFileWriter.Formatter.Type, ShouldEqual, "Json")
		})

		Convey("sub array", func() {
			cfgArr, err := cfg.SubArr("logger.grpc.writers")
			So(err, ShouldBeNil)
			So(len(cfgArr), ShouldEqual, 1)
			So(cfgArr[0].GetString("type"), ShouldEqual, "RotateFile")
			So(cfgArr[0].GetString("rotateFileWriter.filename"), ShouldEqual, "log/account.rpc")
			So(cfgArr[0].GetDuration("rotateFileWriter.maxAge"), ShouldEqual, 24*time.Hour)
			So(cfgArr[0].GetString("rotateFileWriter.formatter.type"), ShouldEqual, "Json")

			type Options struct {
				Type             string `dft:"Stdout" rule:"x in ['RotateFile', 'Stdout']"`
				RotateFileWriter struct {
					Filename  string        `bind:"required"`
					MaxAge    time.Duration `dft:"24h"`
					Formatter struct {
						Type string `dft:"Json" rule:"x in ['Json', 'Text']"`
					}
				}
			}
			var options Options
			So(cfgArr[0].Unmarshal(&options, refx.WithCamelName()), ShouldBeNil)
			So(options.Type, ShouldEqual, "RotateFile")
			So(options.RotateFileWriter.Filename, ShouldEqual, "log/account.rpc")
			So(options.RotateFileWriter.MaxAge, ShouldEqual, 24*time.Hour)
			So(options.RotateFileWriter.Formatter.Type, ShouldEqual, "Json")
		})

		Convey("sub map", func() {
			cfg, err := cfg.SubMap("logger")
			So(err, ShouldBeNil)
			So(cfg["grpc"].GetString("level"), ShouldEqual, "Info")
			So(cfg["grpc"].GetString("writers[0].type"), ShouldEqual, "RotateFile")
			So(cfg["grpc"].GetString("writers[0].rotateFileWriter.filename"), ShouldEqual, "log/account.rpc")
			So(cfg["grpc"].GetDuration("writers[0].rotateFileWriter.maxAge"), ShouldEqual, 24*time.Hour)
			So(cfg["grpc"].GetString("writers[0].rotateFileWriter.formatter.type"), ShouldEqual, "Json")
		})

		DeleteTestFile()
	})
}

func TestConfig_ToString(t *testing.T) {
	Convey("TestConfig_ToString", t, func() {
		CreateTestFile()
		cfg, _ := NewConfigWithSimpleFile("test.json")
		fmt.Println(cfg.ToString())
		fmt.Println(cfg.ToJsonString())
		DeleteTestFile()
	})
}

func TestNewConfigWithBaseFile(t *testing.T) {
	Convey("TestNewConfigWithBaseFile", t, func() {
		CreateTestFile()
		CreateBaseFile()

		cfg, err := NewConfigWithBaseFile("base.json", refx.WithCamelName())
		So(err, ShouldBeNil)
		So(cfg.GetInt("grpc.port"), ShouldEqual, 6080)
		So(cfg.GetInt("http.port"), ShouldEqual, 80)

		DeleteTestFile()
		DeleteBaseFile()
	})
}
