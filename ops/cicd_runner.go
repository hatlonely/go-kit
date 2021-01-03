package ops

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"syscall"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type CICDRunner struct {
	environment []string
	tasks       map[string][]string
}

type Playbook struct {
	Name string
	Env  map[string]map[string]string
	Task map[string][]string
}

func NewCICDRunner(yaml string, varFile string, envName string) (*CICDRunner, error) {
	cfg, err := config.NewConfigWithSimpleFile(varFile, config.WithSimpleFileType("Json"))
	if err != nil {
		return nil, errors.Wrapf(err, "config.NewConfigWithSimpleFile failed. file: [%v]", varFile)
	}
	v, _ := cfg.Get("")
	return NewCICDRunnerWithVariable(yaml, v, envName)
}

func NewCICDRunnerWithVariable(yaml string, v interface{}, envName string) (*CICDRunner, error) {
	buf, err := ioutil.ReadFile(yaml)
	if err != nil {
		return nil, errors.Wrapf(err, "ioutil.ReadFile failed. file: [%v]", yaml)
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
			MemoryProvider: config.MemoryProviderOptions{
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
	environment, err := ParseEnvironment(playbook.Env, envName)
	if err != nil {
		return nil, errors.Wrap(err, "ParseEnvironment failed")
	}

	return &CICDRunner{
		environment: environment,
		tasks:       playbook.Task,
	}, nil
}

type ExecCommandResult struct {
	Status int
	Stdout string
	Stderr string
	Error  error
}

func (y *CICDRunner) Environment() []string {
	return y.environment
}

func (y *CICDRunner) Task() map[string][]string {
	return y.tasks
}

func (y *CICDRunner) RunTaskWithOutput(
	name string, stdout io.Writer, stderr io.Writer,
	onStart func(idx int, length int, command string) error,
	onSuccess func(idx int, length int, command string, status int) error,
	onError func(idx int, length int, command string, err error)) error {
	length := len(y.tasks[name])
	for i, cmd := range y.tasks[name] {
		if err := onStart(i, length, cmd); err != nil {
			return errors.Wrap(err, "onStart failed")
		}
		status, err := ExecCommandWithOutput(cmd, y.environment, stdout, stderr)
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

func (y *CICDRunner) RunTask(name string, callback func(result *ExecCommandResult) error) error {
	for _, cmd := range y.tasks[name] {
		status, stdout, stderr, err := ExecCommand(cmd, y.environment)
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

var ShellVarRegex = regexp.MustCompile(`^\s*\${(.*)}\s*$`)
var ShellCmdRegex = regexp.MustCompile(`^\s*\$\((.*)\)\s*$`)

func evaluate(envMap map[string]string, key string) (string, error) {
	val := envMap[key]
	if ShellVarRegex.MatchString(val) {
		return evaluate(envMap, ShellVarRegex.FindStringSubmatch(val)[1])
	}
	if ShellCmdRegex.MatchString(val) {
		status, stdout, _, err := ExecCommand(ShellCmdRegex.FindStringSubmatch(val)[1], nil)
		if err != nil {
			return "", errors.Wrap(err, "ExecCommand failed")
		}
		if status != 0 {
			return "", errors.Errorf("ExecCommand failed. exit: [%v]", status)
		}
		return strings.TrimSpace(stdout), nil
	}

	return val, nil
}

func ParseEnvironment(environmentMap map[string]map[string]string, name string) ([]string, error) {
	envMap := map[string]string{}
	for key, val := range environmentMap["default"] {
		envMap[key] = val
	}
	if _, ok := environmentMap[name]; !ok {
		return nil, errors.Errorf("unknown environment. name: [%v]", name)
	}
	for key, val := range environmentMap[name] {
		envMap[key] = val
	}
	for key := range envMap {
		val, err := evaluate(envMap, key)
		if err != nil {
			return nil, errors.Wrapf(err, "evaluate failed. key: [%v], val: [%v]", key, envMap[key])
		}
		envMap[key] = val
	}
	var envs []string
	for key, val := range envMap {
		envs = append(envs, fmt.Sprintf(`%s=%s`, key, val))
	}
	envs = append(envs, fmt.Sprintf(`ENVIRONMENT=%s`, name))
	envs = append(envs, fmt.Sprintf(`TMP=tmp/%s`, name))
	sort.Strings(envs)
	return envs, nil
}

func ExecCommandWithOutput(command string, environment []string, stdout io.Writer, stderr io.Writer) (int, error) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, environment...)

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

func ExecCommand(command string, environment []string) (int, string, string, error) {
	var stdout = &bytes.Buffer{}
	var stderr = &bytes.Buffer{}
	status, err := ExecCommandWithOutput(command, environment, stdout, stderr)
	return status, stdout.String(), stderr.String(), err
}
