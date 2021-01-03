package main

import (
	"fmt"
	"os"

	"github.com/hatlonely/go-kit/flag"
	"github.com/hatlonely/go-kit/ops"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/strx"
)

type Options struct {
	flag.Options

	Yaml string `flag:"usage: workflow file; default: .cicd.yaml"`
	Vars string `flag:"usage: variable file;"`
	Env  string `flag:"usage: environment, one of key in env; default: default"`
	Task string `flag:"usage: task, one of key in task"`
}

var Version string

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var options Options
	Must(flag.Struct(&options, refx.WithCamelName()))
	Must(flag.Parse(flag.WithJsonVal()))
	if options.Help {
		strx.Trac(flag.Usage())
		strx.Trac(`
  ops --yaml .cicd.yaml --vars ~/.gomplate/prod.json --env default --task test
`)
		return
	}
	if options.Version {
		strx.Trac(Version)
		return
	}

	runner, err := ops.NewCICDRunner(options.Yaml, options.Vars, options.Env)
	if err != nil {
		panic(err)
	}
	if err := runner.RunTaskWithOutput(
		options.Task, os.Stdout, os.Stderr,
		func(idx int, length int, command string) error {
			strx.Info(fmt.Sprintf("[%v/%v] step: [%v] start", idx+1, length, command))
			return nil
		}, func(idx int, length int, command string, status int) error {
			if status != 0 {
				strx.Warn(fmt.Sprintf("[%v/%v] step: [%v] failed. exit [%v]", idx+1, length, command, status))
			} else {
				strx.Info(fmt.Sprintf("[%v/%v] step: [%v] success", idx+1, length, command))
			}
			strx.Trac("")
			return nil
		}, func(idx int, length int, command string, err error) {
			strx.Warn(fmt.Sprintf("[%v/%v] step: [%v] failed. err: [%v]", idx+1, length, command, err.Error()))
		},
	); err != nil {
		strx.Warn(err.Error())
	}
}
