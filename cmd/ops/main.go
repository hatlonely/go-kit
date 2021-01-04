package main

import (
	"fmt"
	"os"

	"github.com/ghodss/yaml"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/flag"
	"github.com/hatlonely/go-kit/ops"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/strx"
)

type Options struct {
	flag.Options

	Action   string `flag:"-a; usage: actions, one of [dep/run/cmd/env/list/listTask]"`
	Playbook string `flag:"usage: playbook file; default: .ops.yaml"`
	Variable string `flag:"usage: variable file; default: ~/.gomplate/root.json"`
	Env      string `flag:"usage: environment, one of key in env; default: default"`
	Task     string `flag:"usage: task, one of key in task"`
	Command  string `flag:"usage: run command"`
	Force    bool   `flag:"usage: force update dependency"`
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
  ops --playbook .ops.yaml --variable ~/.gomplate/root.json -a listEnv
  ops --playbook .ops.yaml --variable ~/.gomplate/root.json -a listTask
  ops --playbook .ops.yaml --variable ~/.gomplate/root.json -a dep
  ops --playbook .ops.yaml --variable ~/.gomplate/root.json -a dep --force
  ops --playbook .ops.yaml --variable ~/.gomplate/root.json -a env --env test
  ops --playbook .ops.yaml --variable ~/.gomplate/root.json -a run --env test --task test
`)
		return
	}
	if options.Version {
		strx.Trac(Version)
		return
	}

	cfg, err := config.NewConfigWithSimpleFile(options.Playbook, config.WithSimpleFileType("Yaml"))
	if err != nil {
		strx.Warn(fmt.Sprintf("open yaml failed. err: [%v]", err.Error()))
		return
	}
	if options.Action == "list" {
		cfgMap, err := cfg.SubMap("env")
		if err != nil {
			strx.Warn(fmt.Sprintf("parse workflow failed. err: [%v]", err.Error()))
		}
		for key := range cfgMap {
			strx.Trac(key)
		}
		return
	}

	if options.Action == "listTask" {
		v, _ := cfg.Sub("task").Get(options.Task)
		buf, err := yaml.Marshal(v)
		if err != nil {
			panic(err)
		}
		strx.Trac(string(buf))
		return
	}

	runner, err := ops.NewPlaybookRunner(options.Playbook, options.Variable)
	if err != nil {
		strx.Warn(err.Error())
		return
	}

	if options.Action == "dep" {
		if err := runner.DownloadDependencyWithOutput(
			os.Stdout, os.Stderr,
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
			options.Force,
		); err != nil {
			strx.Warn(err.Error())
		}
		return
	}

	if options.Action == "cmd" {
		strx.Info(fmt.Sprintf("[1/1] step: [%v] start", options.Command))
		status, err := runner.ExecCmdWithOutput(options.Env, options.Command, os.Stdout, os.Stderr)
		if err != nil {
			strx.Warn(fmt.Sprintf("[1/1] step: [%v] failed. err: [%v]", options.Command, err.Error()))
		}
		if status != 0 {
			strx.Warn(fmt.Sprintf("[1/1] step: [%v] failed. exit [%v]", options.Command, status))
		}
		strx.Info(fmt.Sprintf("[1/1] step: [%v] success", options.Command))
		return
	}

	if options.Action == "run" {
		if options.Task == "" {
			strx.Warn("task is required")
			return
		}

		if err := runner.RunTaskWithOutput(
			options.Env, options.Task, os.Stdout, os.Stderr,
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
		return
	}

	if options.Action == "env" {
		envs, err := runner.Environment(options.Env)
		if err != nil {
			strx.Warn(err.Error())
			return
		}
		for _, env := range envs {
			strx.Trac(env)
		}
		return
	}

	if options.Action == "listTask" {
		buf, err := yaml.Marshal(runner.Playbook().Task)
		if err != nil {
			panic(err)
		}
		strx.Trac(string(buf))
		return
	}

	if options.Action == "listEnv" {
		for key := range runner.Playbook().Env {
			strx.Trac(key)
		}
		return
	}

	strx.Warn("unknown action. action: [%v]", options.Action)
}
