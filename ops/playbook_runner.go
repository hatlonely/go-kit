package ops

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strings"
	"syscall"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/bind"
	"github.com/hatlonely/go-kit/cast"
	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/validator"
)

type PlaybookRunner struct {
	environment map[string][]string

	playbook *Playbook
}

type Playbook struct {
	Name string
	Tmp  string `dft:"tmp"`
	Env  map[string]map[string]string
	Dep  map[string]map[string]string
	Task map[string]struct {
		Args map[string]struct {
			Validation string
			Type       string `dft:"string"`
			Default    string
		}
		Step []string
	}
}

func NewPlaybookRunner(playbook string, varFile string) (*PlaybookRunner, error) {
	cfg, err := config.NewConfigWithSimpleFile(varFile, config.WithSimpleFileType("Json"))
	if err != nil {
		return nil, errors.Wrapf(err, "config.NewConfigWithSimpleFile failed. file: [%v]", varFile)
	}
	v, _ := cfg.Get("")
	return NewPlaybookRunnerWithVariable(playbook, v)
}

func NewPlaybookRunnerWithVariable(playbookFile string, v interface{}) (*PlaybookRunner, error) {
	buf, err := ioutil.ReadFile(playbookFile)
	if err != nil {
		return nil, errors.Wrapf(err, "ioutil.ReadFile failed. file: [%v]", playbookFile)
	}
	tpl, err := template.New("").Parse(string(buf))
	if err != nil {
		return nil, errors.Wrap(err, "template.Parse failed")
	}
	var renderBuf bytes.Buffer

	if err := tpl.Execute(&renderBuf, v); err != nil {
		return nil, errors.Wrap(err, "tpl.Execute failed")
	}

	cfg, err := config.NewConfigWithOptions(&config.Options{
		Decoder: config.DecoderOptions{
			Type: "Yaml",
		},
		Provider: config.ProviderOptions{
			Type: "Memory",
			Options: &config.MemoryProviderOptions{
				Buffer: renderBuf.String(),
			},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "config.NewConfigWithSimpleFile failed")
	}
	var playbook Playbook
	if err := cfg.Unmarshal(&playbook, refx.WithCamelName()); err != nil {
		return nil, errors.Wrap(err, "cfg.Unmarshal failed")
	}

	return &PlaybookRunner{
		playbook:    &playbook,
		environment: map[string][]string{},
	}, nil
}

type ExecCommandResult struct {
	Status int
	Stdout string
	Stderr string
	Error  error
}

func (r *PlaybookRunner) DownloadDependencyWithOutput(
	stdout io.Writer, stderr io.Writer,
	onStart func(idx int, length int, command string) error,
	onSuccess func(idx int, length int, command string, status int) error,
	onError func(idx int, length int, command string, err error),
	forceUpdate bool) error {
	length := len(r.playbook.Dep)
	i := 0
	for key, val := range r.playbook.Dep {
		output := path.Join(r.playbook.Tmp, "dep", key)
		if !forceUpdate {
			if _, err := os.Stat(output); err == nil {
				continue
			} else if !os.IsNotExist(err) {
				return errors.Wrap(err, "os.Stat failed")
			}
		}

		if forceUpdate {
			_ = os.RemoveAll(output)
		}

		var cmd string
		workDir := ""
		switch val["type"] {
		case "git", "":
			if len(val["url"]) == 0 {
				return errors.New("url is required")
			}
			version := "master"
			if len(val["version"]) == 0 {
				version = val["version"]
			}
			cmd = fmt.Sprintf(`git clone --depth=1 --branch "%s" "%s" "%s"`, version, val["url"], output)
		case "cmd":
			cmd = val["command"]
			workDir = output
			_ = os.MkdirAll(workDir, 0755)
		default:
			return errors.Errorf("unsupported dependency type [%v]", val["type"])
		}
		if err := onStart(i, length, cmd); err != nil {
			return errors.Wrap(err, "onStart failed")
		}
		status, err := ExecCommandWithOutput(cmd, nil, workDir, stdout, stderr)
		if err := onSuccess(i, length, cmd, status); err != nil {
			return errors.Wrap(err, "onSuccess failed")
		}
		if err != nil {
			onError(i, length, cmd, err)
			return errors.Wrap(err, "ExecCommand failed. cmd [%v]")
		}
		if status != 0 {
			return errors.Errorf("ExecCommand failed. cmd [%v], exit: [%v]", cmd, status)
		}
		i++
	}
	return nil
}

func (r *PlaybookRunner) Playbook() *Playbook {
	return r.playbook
}

func (r *PlaybookRunner) Environment(env string) ([]string, error) {
	if r.environment[env] == nil {
		environment, err := ParseEnvironment(r.playbook.Env, r.playbook.Tmp, env)
		if err != nil {
			return nil, errors.Wrap(err, "ParseEnvironment failed")
		}
		r.environment[env] = environment
	}

	return r.environment[env], nil
}

func (r *PlaybookRunner) ExecCmdWithOutput(env string, cmd string, stdout io.Writer, stderr io.Writer) (int, error) {
	_, err := r.Environment(env)
	if err != nil {
		return 0, err
	}
	return ExecCommandWithOutput(cmd, r.environment[env], "", stdout, stderr)
}

func (r *PlaybookRunner) taskEnvironment(env string, taskName string, getter bind.Getter) ([]string, error) {
	_, err := r.Environment(env)
	if err != nil {
		return nil, err
	}
	environment := r.environment[env]

	for key, info := range r.playbook.Task[taskName].Args {
		v, ok := getter.Get(key)
		if !ok {
			v = info.Default
		}
		switch info.Type {
		case "int":
			v, err = cast.ToIntE(v)
		case "bool":
			v, err = cast.ToBoolE(v)
		case "string":
			v, err = cast.ToStringE(v)
		default:
			err = errors.Errorf("unsupported type [%v]", info.Type)
		}
		if err != nil {
			return nil, errors.Wrapf(err, "cast failed. key: [%v], val: [%v], type: [%v]", key, v, info.Type)
		}
		if info.Validation != "" {
			eval, err := validator.Lang.NewEvaluable(info.Validation)
			if err != nil {
				return nil, errors.Wrapf(err, "create evaluable failed. key: [%v], validate: [%v]", key, info.Validation)
			}
			if ok, err := eval.EvalBool(context.Background(), map[string]interface{}{"x": v}); err != nil {
				return nil, errors.Wrapf(err, "eval.EvalBool failed")
			} else if !ok {
				return nil, errors.Errorf("validate failed. key: [%v], validate: [%v]", key, info.Validation)
			}
		}
		environment = append(environment, fmt.Sprintf("%s=%v", key, v))
	}
	return environment, nil
}

func (r *PlaybookRunner) DescTask(env string, taskName string, getter bind.Getter) (
	[]string, []string, error) {
	environment, err := r.taskEnvironment(env, taskName, getter)
	if err != nil {
		return nil, nil, err
	}
	return environment, r.playbook.Task[taskName].Step, nil
}

func (r *PlaybookRunner) RunTaskWithOutput(
	env string, getter bind.Getter, taskName string,
	startStep int, endStep int,
	stdout io.Writer, stderr io.Writer,
	onStart func(idx int, length int, command string) error,
	onSuccess func(idx int, length int, command string, status int) error,
	onError func(idx int, length int, command string, err error)) error {
	if startStep < 1 {
		return errors.Errorf("start step [%v] should start 1", startStep)
	}
	length := len(r.playbook.Task[taskName].Step)
	if endStep == -1 {
		endStep = length
	}
	if startStep > endStep {
		return errors.Errorf("start step [%v] should less than end step [%v]", startStep, endStep)
	}
	if endStep > length {
		return errors.Errorf("end step [%v] should less than length [%v]", endStep, length)
	}

	environment, err := r.taskEnvironment(env, taskName, getter)
	if err != nil {
		return err
	}

	startIdx, endIdx := startStep-1, endStep
	for i := startIdx; i < endIdx; i++ {
		cmd := r.playbook.Task[taskName].Step[i]
		if err := onStart(i, length, cmd); err != nil {
			return errors.Wrap(err, "onStart failed")
		}
		status, err := ExecCommandWithOutput(cmd, environment, "", stdout, stderr)
		if err := onSuccess(i, length, cmd, status); err != nil {
			return errors.Wrap(err, "onSuccess failed")
		}
		if err != nil {
			onError(i, length, cmd, err)
			return errors.Wrap(err, "ExecCommand failed. cmd [%v]")
		}
		if status != 0 {
			return errors.Errorf("ExecCommand failed. cmd [%v], exit: [%v]", cmd, status)
		}
	}
	return nil
}

func (r *PlaybookRunner) RunTask(env string, taskName string, callback func(result *ExecCommandResult) error) error {
	_, err := r.Environment(env)
	if err != nil {
		return err
	}

	for _, cmd := range r.playbook.Task[taskName].Step {
		status, stdout, stderr, err := ExecCommand(cmd, r.environment[env], "")
		if err := callback(&ExecCommandResult{
			Status: status,
			Stdout: stdout,
			Stderr: stderr,
			Error:  err,
		}); err != nil {
			return err
		}
		if err != nil {
			return errors.Wrap(err, "ExecCommand failed")
		}
	}
	return nil
}

var ShellVarRegex = regexp.MustCompile(`\${(\w+).*?}`)

func topOrderEnvMap(envMap map[string]string) ([]string, error) {
	allKey := map[string]bool{}
	for key := range envMap {
		allKey[key] = true
	}

	var keys []string
	for len(allKey) != 0 {
		length := len(allKey)
		for key := range allKey {
			find := false
			for _, vs := range ShellVarRegex.FindAllStringSubmatch(envMap[key], -1) {
				k := vs[1]
				if k == key {
					continue
				}
				if allKey[k] {
					find = true
				}
			}
			if find {
				continue
			}
			keys = append(keys, key)
			delete(allKey, key)
		}
		if length == len(allKey) {
			return nil, errors.Errorf("loop graph. keys: [%v]", allKey)
		}
	}

	return keys, nil
}

func ParseEnvironment(environmentMap map[string]map[string]string, tmp string, env string) ([]string, error) {
	envMap := map[string]string{}
	for key, val := range environmentMap["default"] {
		envMap[key] = val
	}
	if _, ok := environmentMap[env]; !ok {
		return nil, errors.Errorf("unknown environment. env: [%v]", env)
	}
	for key, val := range environmentMap[env] {
		envMap[key] = val
	}
	if tmp == "" {
		tmp = "tmp"
	}
	tmpDir := fmt.Sprintf(`%s/env/%s`, tmp, env)
	envMap["ENV"] = env
	envMap["TMP"] = tmpDir
	envMap["DEP"] = fmt.Sprintf("%s/dep", tmp)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, errors.Wrapf(err, "os.MkdirAll failed. dir: [%v]", tmpDir)
	}

	keys, err := topOrderEnvMap(envMap)
	if err != nil {
		return nil, errors.Wrap(err, "topOrderEnvMap failed")
	}
	var envs []string
	for _, key := range keys {
		status, stdout, _, err := ExecCommand(fmt.Sprintf("eval echo %v", envMap[key]), envs, "")
		if err != nil {
			return nil, errors.Wrap(err, "ExecCommand failed")
		}
		if status != 0 {
			return nil, errors.Errorf("ExecCommand failed. exit: [%v]", status)
		}
		envs = append(envs, fmt.Sprintf(`%s=%s`, key, strings.TrimSpace(stdout)))
	}

	sort.Strings(envs)
	return envs, nil
}

func ExecCommandWithOutput(command string, environment []string, workDir string, stdout io.Writer, stderr io.Writer) (int, error) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, environment...)

	if workDir != "" {
		cmd.Dir = workDir
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return -1, err
	}

	if err := cmd.Wait(); err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			if status, ok := e.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus(), nil
			}
		}

		return -1, err
	}

	return 0, nil
}

func ExecCommand(command string, environment []string, workDir string) (int, string, string, error) {
	var stdout = &bytes.Buffer{}
	var stderr = &bytes.Buffer{}
	status, err := ExecCommandWithOutput(command, environment, workDir, stdout, stderr)
	return status, stdout.String(), stderr.String(), err
}
