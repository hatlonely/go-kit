package astx

import (
	"bytes"
	"fmt"
	"html/template"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type WrapperGenerator struct {
	options *WrapperGeneratorOptions

	wrapClassMap map[string]string
}

type Rule struct {
	Include *regexp.Regexp
	Exclude *regexp.Regexp
}

type WrapperGeneratorOptions struct {
	GoPath      string   `flag:"usage: gopath; default: vendor"`
	PkgPath     string   `flag:"usage: source path"`
	Package     string   `flag:"usage: package name"`
	Classes     []string `flag:"usage: classes to wrap"`
	ClassPrefix string   `flag:"usage: wrap class name"`

	Rule struct {
		Class    Rule
		Function map[string]Rule
		Trace    map[string]Rule
		Retry    map[string]Rule
	}
}

func NewWrapperGeneratorWithOptions(options *WrapperGeneratorOptions) *WrapperGenerator {
	wrapClassMap := map[string]string{}
	for _, cls := range options.Classes {
		wrapClassMap[cls] = fmt.Sprintf("%s%sWrapper", options.ClassPrefix, cls)
	}

	return &WrapperGenerator{
		options:      options,
		wrapClassMap: wrapClassMap,
	}
}

func (g *WrapperGenerator) Generate() (string, error) {
	functions, err := ParseFunction(path.Join(g.options.GoPath, g.options.PkgPath), g.options.Package)
	if err != nil {
		return "", errors.Wrap(err, "ParseFunction failed")
	}

	var buf bytes.Buffer

	buf.WriteString(g.generateWrapperHeader())

	var classes []string
	if len(g.options.Classes) == 0 {
		for _, function := range functions {
			if !function.IsMethod {
				continue
			}
			if _, ok := g.wrapClassMap[function.Class]; ok {
				continue
			}
			if !g.MatchRule(function.Class, g.options.Rule.Class) {
				continue
			}
			g.wrapClassMap[function.Class] = fmt.Sprintf("%s%sWrapper", g.options.ClassPrefix, function.Class)
			classes = append(classes, function.Class)
		}
	} else {
		classes = g.options.Classes
	}

	sort.Strings(classes)
	for _, cls := range classes {
		buf.WriteString(g.generateWrapperStruct(cls))
	}

	for _, function := range functions {
		if !function.IsMethod {
			continue
		}
		if _, ok := g.wrapClassMap[function.Class]; !ok {
			continue
		}
		if !g.MatchFunctionRule(function, g.options.Rule.Function) {
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
		if g.MatchFunctionRule(function, g.options.Rule.Trace) || g.MatchFunctionRule(function, g.options.Rule.Retry) {
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

func (g *WrapperGenerator) generateWrapperReturnChain(function *Function) string {
	var params []string
	for _, i := range function.Params {
		if strings.HasPrefix(i.Type, "...") {
			params = append(params, fmt.Sprintf("%s...", i.Name))
		} else {
			params = append(params, i.Name)
		}
	}

	return fmt.Sprintf("\tw.obj = w.obj.%s(%s)\n", function.Name, strings.Join(params, ", "))
}

func (g *WrapperGenerator) generateWrapperFunctionBody(function *Function) string {
	var buf bytes.Buffer
	if function.IsChain {
		buf.WriteString(g.generateWrapperReturnChain(function))
		buf.WriteString("\treturn w")
		return buf.String()
	}

	if g.MatchFunctionRule(function, g.options.Rule.Trace) {
		buf.WriteString(g.generateWrapperOpentracing(function))
	}

	if function.IsReturnVoid {
		buf.WriteString(g.generateWrapperReturnVoid(function))
		return buf.String()
	}

	if function.IsReturnError && g.MatchFunctionRule(function, g.options.Rule.Retry) {
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

func (g *WrapperGenerator) MatchFunctionRule(function *Function, rules map[string]Rule) bool {
	fun := function.Name
	cls := function.Class

	if _, ok := rules[cls]; ok {
		return g.MatchRule(fun, rules[cls])
	}
	return g.MatchRule(fun, rules["default"])
}

func (g *WrapperGenerator) MatchRule(key string, rule Rule) bool {
	if rule.Include != nil {
		if rule.Include.MatchString(key) {
			return true
		}
	}

	if rule.Exclude != nil {
		if rule.Exclude.MatchString(key) {
			return false
		}
	}

	return true
}
