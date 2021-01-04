package ops

import (
	"bytes"
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

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
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
	Task map[string][]string
}

func NewPlaybookRunner(yaml string, varFile string) (*PlaybookRunner, error) {
	cfg, err := config.NewConfigWithSimpleFile(varFile, config.WithSimpleFileType("Json"))
	if err != nil {
		return nil, errors.Wrapf(err, "config.NewConfigWithSimpleFile failed. file: [%v]", varFile)
	}
	v, _ := cfg.Get("")
	return NewPlaybookRunnerWithVariable(yaml, v)
}

func NewPlaybookRunnerWithVariable(yaml string, v interface{}) (*PlaybookRunner, error) {
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
		default:
			return errors.Errorf("unsupported dependency type [%v]", val["type"])
		}
		if err := onStart(i, length, cmd); err != nil {
			return errors.Wrap(err, "onStart failed")
		}
		status, err := ExecCommandWithOutput(cmd, nil, stdout, stderr)
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
		environment, err := ParseEnvironment(r.playbook.Env, env)
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
	return ExecCommandWithOutput(cmd, r.environment[env], stdout, stderr)
}

func (r *PlaybookRunner) RunTaskWithOutput(
	env string, taskName string, stdout io.Writer, stderr io.Writer,
	onStart func(idx int, length int, command string) error,
	onSuccess func(idx int, length int, command string, status int) error,
	onError func(idx int, length int, command string, err error)) error {
	_, err := r.Environment(env)
	if err != nil {
		return err
	}

	length := len(r.playbook.Task[taskName])
	for i, cmd := range r.playbook.Task[taskName] {
		if err := onStart(i, length, cmd); err != nil {
			return errors.Wrap(err, "onStart failed")
		}
		status, err := ExecCommandWithOutput(cmd, r.environment[env], stdout, stderr)
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

	for _, cmd := range r.playbook.Task[taskName] {
		status, stdout, stderr, err := ExecCommand(cmd, r.environment[env])
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

var ShellVarRegex = regexp.MustCompile(`\${(\w+)}`)
var ShellCmdRegex = regexp.MustCompile(`\s*\$\((.*)\)\s*`)

func evaluate(envMap map[string]string, key string) (string, error) {
	val, ok := envMap[key]
	if !ok {
		return "", errors.Errorf("evalute failed. no such key: [%v]", key)
	}
	if ShellVarRegex.MatchString(val) {
		var err error
		val := ShellVarRegex.ReplaceAllStringFunc(val, func(item string) string {
			key := ShellVarRegex.FindStringSubmatch(item)[1]

			val, e := evaluate(envMap, key)
			if e != nil {
				err = errors.Wrapf(e, "evaluate variable failed. key: [%v]", key)
				return item
			}

			return val
		})
		if err != nil {
			return "", err
		}
		envMap[key] = val
		return val, nil
	}
	if ShellCmdRegex.MatchString(val) {
		status, stdout, _, err := ExecCommand(ShellCmdRegex.FindStringSubmatch(val)[1], nil)
		if err != nil {
			return "", errors.Wrap(err, "ExecCommand failed")
		}
		if status != 0 {
			return "", errors.Errorf("ExecCommand failed. exit: [%v]", status)
		}
		envMap[key] = strings.TrimSpace(stdout)
		return strings.TrimSpace(stdout), nil
	}

	return val, nil
}

func ParseEnvironment(environmentMap map[string]map[string]string, env string) ([]string, error) {
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
	tmpDir := fmt.Sprintf(`tmp/env/%s`, env)
	envMap["ENV"] = env
	envMap["TMP"] = tmpDir
	envMap["DEP"] = "tmp/dep"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, errors.Wrapf(err, "os.MkdirAll failed. dir: [%v]", tmpDir)
	}
	for key := range envMap {
		_, err := evaluate(envMap, key)
		if err != nil {
			return nil, errors.Wrapf(err, "evaluate failed. key: [%v], val: [%v]", key, envMap[key])
		}
	}
	var envs []string
	for key, val := range envMap {
		envs = append(envs, fmt.Sprintf(`%s=%s`, key, val))
	}
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
