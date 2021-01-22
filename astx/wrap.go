package astx

import (
	"bytes"
	"fmt"
	"html/template"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type WrapperGenerator struct {
	options *WrapperGeneratorOptions

	wrapClassMap map[string]string
	rule         struct {
		trace map[string]rule
		retry map[string]rule
	}
}

type rule struct {
	include *regexp.Regexp
	exclude *regexp.Regexp
}

type Rule struct {
	Include string
	Exclude string
}

type WrapperGeneratorOptions struct {
	GoPath      string   `flag:"usage: gopath; default: vendor"`
	PkgPath     string   `flag:"usage: source path"`
	Package     string   `flag:"usage: package name"`
	Classes     []string `flag:"usage: classes to wrap"`
	ClassPrefix string   `flag:"usage: wrap class name"`

	Rule struct {
		Trace map[string]Rule
		Retry map[string]Rule
	}
}

func parseRule(rules map[string]Rule) map[string]rule {
	m := map[string]rule{}
	for key, val := range rules {
		var include *regexp.Regexp
		var exclude *regexp.Regexp
		if val.Include != "" {
			include = regexp.MustCompile(val.Include)
		}
		if val.Exclude != "" {
			exclude = regexp.MustCompile(val.Exclude)
		}

		m[key] = rule{
			include: include,
			exclude: exclude,
		}
	}
	return m
}

func NewWrapperGeneratorWithOptions(options *WrapperGeneratorOptions) *WrapperGenerator {
	wrapClassMap := map[string]string{}
	for _, cls := range options.Classes {
		wrapClassMap[cls] = fmt.Sprintf("%s%sWrapper", options.ClassPrefix, cls)
	}

	return &WrapperGenerator{
		options:      options,
		wrapClassMap: wrapClassMap,
		rule: struct {
			trace map[string]rule
			retry map[string]rule
		}{
			trace: parseRule(options.Rule.Trace),
			retry: parseRule(options.Rule.Retry),
		},
	}
}

func (g *WrapperGenerator) Generate() (string, error) {
	functions, err := ParseFunction(path.Join(g.options.GoPath, g.options.PkgPath), g.options.Package)
	if err != nil {
		return "", errors.Wrap(err, "ParseFunction failed")
	}

	var buf bytes.Buffer

	buf.WriteString(g.generateWrapperHeader())

	for _, cls := range g.options.Classes {
		buf.WriteString(g.generateWrapperStruct(cls))
	}

	for _, function := range functions {
		if !function.IsMethod {
			continue
		}
		if _, ok := g.wrapClassMap[function.Class]; !ok {
			continue
		}

		buf.WriteString("\n")
		buf.WriteString(g.generateWrapperFunctionDeclare(function))
		buf.WriteString(" {")
		buf.WriteString(g.generateWrapperFunctionBody(function))
		buf.WriteString("}\n")
	}

	return buf.String(), nil
}

func (g *WrapperGenerator) generateWrapperHeader() string {
	const tplStr = `// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"

	"{{.pkgPath}}"
	"github.com/opentracing/opentracing-go"
)
`

	tpl, _ := template.New("").Parse(tplStr)

	var buf bytes.Buffer
	_ = tpl.Execute(&buf, map[string]string{
		"pkgPath": g.options.PkgPath,
	})

	return buf.String()
}

func (g *WrapperGenerator) generateWrapperStruct(cls string) string {
	const tplStr = `
type {{.wrapClass}} struct {
	obj     *{{.package}}.{{.class}}
	retry   *Retry
	options *WrapperOptions
}
`

	tpl, _ := template.New("").Parse(tplStr)

	var buf bytes.Buffer
	_ = tpl.Execute(&buf, map[string]string{
		"package":   g.options.Package,
		"class":     cls,
		"wrapClass": g.wrapClassMap[cls],
	})

	return buf.String()
}

func (g *WrapperGenerator) generateWrapperFunctionDeclare(function *Function) string {
	var buf bytes.Buffer

	buf.WriteString("func ")
	if function.Recv != nil {
		buf.WriteString(fmt.Sprintf("(w *%s)", g.wrapClassMap[function.Class]))
	}

	buf.WriteString(" ")
	buf.WriteString(function.Name)

	buf.WriteString("(")

	var params []string
	if len(function.Params) == 0 || function.Params[0].Type != "context.Context" {
		if g.MatchRule(function, g.rule.trace[function.Class]) || g.MatchRule(function, g.rule.retry[function.Class]) {
			params = append(params, "ctx context.Context")
		}
	}

	for _, i := range function.Params {
		params = append(params, fmt.Sprintf("%s %s", i.Name, i.Type))
	}
	buf.WriteString(strings.Join(params, ", "))
	buf.WriteString(") ")

	var results []string
	for _, i := range function.Results {
		cls := strings.TrimLeft(i.Type, "*")
		if !strings.HasPrefix(cls, g.options.Package) {
			results = append(results, i.Type)
			continue
		}

		cls = strings.TrimPrefix(cls, g.options.Package+".")
		if wrapCls, ok := g.wrapClassMap[cls]; ok {
			results = append(results, fmt.Sprintf(`*%s`, wrapCls))
			continue
		}
		results = append(results, i.Type)
	}

	if len(function.Results) >= 2 {
		buf.WriteString("(")
		buf.WriteString(strings.Join(results, ", "))
		buf.WriteString(")")
	} else {
		buf.WriteString(strings.Join(results, ", "))
	}

	return buf.String()
}

func (g *WrapperGenerator) generateWrapperOpentracing(function *Function) string {
	return fmt.Sprintf(`
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "%s.%s.%s")
		defer span.Finish()
	}
`, g.options.Package, function.Class, function.Name)
}

func (g *WrapperGenerator) generateWrapperDeclareReturnVariables(function *Function) string {
	var buf bytes.Buffer
	for _, field := range function.Results {
		buf.WriteString(fmt.Sprintf("\n	var %s %s", field.Name, field.Type))
	}
	return buf.String()
}

func (g *WrapperGenerator) generateWrapperCallWithRetry(function *Function) string {
	if !function.IsReturnError {
		panic(fmt.Sprintf("generateWrapperCallWithRetry with no error function. function: [%v]", function.Name))
	}

	var params []string
	for _, i := range function.Params {
		if strings.HasPrefix(i.Type, "...") {
			params = append(params, fmt.Sprintf("%s...", i.Name))
		} else {
			params = append(params, i.Name)
		}
	}

	var results []string
	for _, i := range function.Results {
		results = append(results, fmt.Sprintf("%s", i.Name))
	}

	return fmt.Sprintf(`
	err = w.retry.Do(func() error {
		%s = w.obj.%s(%s)
		return %s
	})
`, strings.Join(results, ", "), function.Name, strings.Join(params, ", "), results[len(results)-1])
}

func (g *WrapperGenerator) generateWrapperCallWithoutRetry(function *Function) string {
	var params []string
	for _, i := range function.Params {
		if strings.HasPrefix(i.Type, "...") {
			params = append(params, fmt.Sprintf("%s...", i.Name))
		} else {
			params = append(params, i.Name)
		}
	}

	var results []string
	for _, i := range function.Results {
		results = append(results, fmt.Sprintf("%s", i.Name))
	}

	return fmt.Sprintf("\t%s := w.obj.%s(%s)", strings.Join(results, ", "), function.Name, strings.Join(params, ", "))
}

func (g *WrapperGenerator) generateWrapperReturnVariables(function *Function) string {
	if function.IsReturnVoid {
		panic(fmt.Sprintf("generateWrapperReturnVariables with void function. function [%v]", function.Name))
	}

	var results []string
	for _, i := range function.Results {
		cls := strings.TrimLeft(i.Type, "*")
		if !strings.HasPrefix(cls, g.options.Package) {
			results = append(results, fmt.Sprintf("%s", i.Name))
			continue
		}

		cls = strings.TrimPrefix(cls, g.options.Package+".")
		if wrapCls, ok := g.wrapClassMap[cls]; ok {
			results = append(results, fmt.Sprintf(`&%s{obj: %s, retry: w.retry, options: w.options}`, wrapCls, i.Name))
			continue
		}

		results = append(results, fmt.Sprintf("%s", i.Name))
	}

	return fmt.Sprintf("\treturn %s\n", strings.Join(results, ", "))
}

func (g *WrapperGenerator) generateWrapperReturnFunction(function *Function) string {
	if function.IsReturnVoid {
		panic(fmt.Sprintf("generateWrapperReturnFunction with void function. function [%v]", function.Name))
	}

	var params []string
	for _, i := range function.Params {
		if strings.HasPrefix(i.Type, "...") {
			params = append(params, fmt.Sprintf("%s...", i.Name))
		} else {
			params = append(params, i.Name)
		}
	}

	return fmt.Sprintf("\treturn w.obj.%s(%s)\n", function.Name, strings.Join(params, ", "))
}

func (g *WrapperGenerator) generateWrapperReturnVoid(function *Function) string {
	var params []string
	for _, i := range function.Params {
		if strings.HasPrefix(i.Type, "...") {
			params = append(params, fmt.Sprintf("%s...", i.Name))
		} else {
			params = append(params, i.Name)
		}
	}

	return fmt.Sprintf("\tw.obj.%s(%s)\n", function.Name, strings.Join(params, ", "))
}

func (g *WrapperGenerator) generateWrapperFunctionBody(function *Function) string {
	var buf bytes.Buffer
	if function.IsChain {
		buf.WriteString(g.generateWrapperReturnVoid(function))
		buf.WriteString("\treturn w")
		return buf.String()
	}

	if g.MatchRule(function, g.rule.trace[function.Class]) {
		buf.WriteString(g.generateWrapperOpentracing(function))
	}

	if function.IsReturnVoid {
		buf.WriteString(g.generateWrapperReturnVoid(function))
		return buf.String()
	}

	if function.IsReturnError && g.MatchRule(function, g.rule.retry[function.Class]) {
		buf.WriteString(g.generateWrapperDeclareReturnVariables(function))
		buf.WriteString(g.generateWrapperCallWithRetry(function))
		buf.WriteString(g.generateWrapperReturnVariables(function))
		return buf.String()
	}

	buf.WriteString("\n")
	buf.WriteString(g.generateWrapperCallWithoutRetry(function))
	buf.WriteString("\n")
	buf.WriteString(g.generateWrapperReturnVariables(function))

	return buf.String()
}

func (g *WrapperGenerator) MatchRule(function *Function, rule rule) bool {
	if rule.include != nil {
		if rule.include.MatchString(function.Name) {
			return true
		}
	}

	if rule.exclude != nil {
		if rule.exclude.MatchString(function.Name) {
			return false
		}
	}

	return true
}
