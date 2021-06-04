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

	Action       string `flag:"-a; usage: actions, one of [dep/run/cmd/env/desc/listEnv/listTask]"`
	Playbook     string `flag:"usage: playbook file; default: .ops.yaml"`
	Variable     string `flag:"usage: variable file"`
	Env          string `flag:"usage: environment, one of key in env; default: default"`
	Task         string `flag:"usage: task, one of key in task"`
	StartStep    int    `flag:"usage: run task start step. start from 1; default: 1"`
	EndStep      int    `flag:"usage: run task end step. set -1 to the end; default: -1"`
	Command      string `flag:"usage: run command"`
	Force        bool   `flag:"usage: force update dependency"`
	VariableType string `flag:"usage: variable file type; default: Yaml"`
	BaseFile     string `flag:"usage: get variable from base file"`
	Base         config.Options
}

var Version string

func main() {
	var options Options
	refx.Must(flag.Struct(&options, refx.WithCamelName()))
	refx.Must(flag.Parse(flag.WithJsonVal()))
	if options.Help {
		strx.Trac(flag.Usage())
		strx.Trac(`
  ops --playbook .ops.yaml --variable ~/.ops/root.json -a listEnv
  ops --playbook .ops.yaml --variable ~/.ops/root.json -a listTask
  ops --playbook .ops.yaml --variable ~/.ops/root.json -a dep
  ops --playbook .ops.yaml --variable ~/.ops/root.json -a dep --force
  ops --playbook .ops.yaml --variable ~/.ops/root.json -a env --env test
  ops --playbook .ops.yaml --variable ~/.ops/root.json -a run --env test --task test
`)
		return
	}
	if options.Version {
		strx.Trac(Version)
		return
	}

	var cfg *config.Config
	var err error
	if len(options.Variable) != 0 {
		cfg, err = config.NewConfigWithSimpleFile(options.Variable, config.WithSimpleFileType(options.VariableType))
	} else if len(options.BaseFile) != 0 {
		cfg, err = config.NewConfigWithBaseFile(options.BaseFile)
	} else {
		cfg, err = config.NewConfigWithOptions(&options.Base)
	}
	if err != nil {
		strx.Warn(err.Error())
		return
	}

	v, _ := cfg.Get("")
	runner, err := ops.NewPlaybookRunnerWithVariable(options.Playbook, v)
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
			options.Env, flag.Instance(), options.Task,
			options.StartStep, options.EndStep,
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

	if options.Action == "desc" {
		envs, steps, err := runner.DescTask(options.Env, options.Task, flag.Instance())
		if err != nil {
			strx.Warn(err.Error())
			return
		}

		if options.StartStep < 1 {
			strx.Warn(fmt.Sprintf("start step [%v] should start 1", options.StartStep))
			return
		}
		length := len(steps)
		if options.EndStep == -1 {
			options.EndStep = len(steps)
		}
		if options.StartStep > options.EndStep {
			strx.Warn(fmt.Sprintf("start step [%v] should less than end step [%v]", options.StartStep, options.EndStep))
			return
		}
		if options.EndStep > length {
			strx.Warn(fmt.Sprintf("end step [%v] should less than length [%v]", options.StartStep, options.EndStep))
			return
		}

		for _, env := range envs {
			strx.Trac(env)
		}
		startIdx, endIdx := options.StartStep-1, options.EndStep
		for i, step := range steps {
			if i >= startIdx && i < endIdx {
				strx.Info(fmt.Sprintf("[%v/%v] step: [%v]", i+1, len(steps), step))
			} else {
				strx.Trac(fmt.Sprintf("[%v/%v] step: [%v]", i+1, len(steps), step))
			}
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
