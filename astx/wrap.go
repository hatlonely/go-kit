package astx

import (
	"bytes"
	"fmt"
	"html/template"
	"path"
	"strings"

	"github.com/pkg/errors"
)

type WrapperGenerator struct {
	options *WrapperGeneratorOptions
}

type WrapperGeneratorOptions struct {
	GoPath      string
	PkgPath     string
	Package     string
	Class       string
	WrapClass   string
	ClassPrefix string
	Output      string
}

func NewWrapperGeneratorWithOptions(options *WrapperGeneratorOptions) *WrapperGenerator {
	if options.WrapClass == "" {
		options.WrapClass = fmt.Sprintf("%s%sWrapper", options.ClassPrefix, options.Class)
	}
	return &WrapperGenerator{
		options: options,
	}
}

func (g *WrapperGenerator) Generate() (string, error) {
	functions, err := ParseFunction(path.Join(g.options.GoPath, g.options.PkgPath), g.options.Package)
	if err != nil {
		return "", errors.Wrap(err, "ParseFunction failed")
	}

	var buf bytes.Buffer

	buf.WriteString(g.generateWrapperHeader())
	buf.WriteString(g.generateWrapperStruct())

	for _, function := range functions {
		if function.Recv == nil {
			continue
		}
		if function.Recv.Type != fmt.Sprintf("%s.%s", g.options.Package, g.options.Class) &&
			function.Recv.Type != fmt.Sprintf("*%s.%s", g.options.Package, g.options.Class) {
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
package wrapper

import (
	"context"

	"{{.pkgPath}}"
	"github.com/avast/retry-go"
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

func (g *WrapperGenerator) generateWrapperStruct() string {
	const tplStr = `
type {{.wrapClass}} struct {
	client *{{.package}}.{{.class}}
	retry  *Retry
}
`

	tpl, _ := template.New("").Parse(tplStr)

	var buf bytes.Buffer
	_ = tpl.Execute(&buf, map[string]string{
		"package":   g.options.Package,
		"class":     g.options.Class,
		"wrapClass": g.options.WrapClass,
	})

	return buf.String()
}

func (g *WrapperGenerator) generateWrapperFunctionDeclare(function *Function) string {
	var buf bytes.Buffer

	buf.WriteString("func ")
	if function.Recv != nil {
		buf.WriteString(fmt.Sprintf("(w *%s)", g.options.WrapClass))
	}

	buf.WriteString(" ")
	buf.WriteString(function.Name)

	buf.WriteString("(")

	var params []string
	if len(function.Params) == 0 || function.Params[0].Type != "context.Context" {
		params = append(params, "ctx context.Context")
	}

	for _, i := range function.Params {
		params = append(params, fmt.Sprintf("%s %s", i.Name, i.Type))
	}
	buf.WriteString(strings.Join(params, ", "))
	buf.WriteString(") ")

	var results []string
	for _, i := range function.Results {
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
	span, _ := opentracing.StartSpanFromContext(ctx, "%s.%s.%s")
	defer span.Finish()
`, g.options.Package, g.options.Class, function.Name)
}

func (g *WrapperGenerator) generateWrapperDeclareReturnVariables(function *Function) string {
	var buf bytes.Buffer
	for _, field := range function.Results {
		buf.WriteString(fmt.Sprintf("\n	var %s %s", field.Name, field.Type))
	}
	return buf.String()
}

func (g *WrapperGenerator) generateWrapperRetry(function *Function) string {
	if !function.IsReturnError {
		panic(fmt.Sprintf("generateWrapperRetry with no error function. function: [%v]", function.Name))
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
		%s = w.client.%s(%s)
		return %s
	})
`, strings.Join(results, ", "), function.Name, strings.Join(params, ", "), results[len(results)-1])
}

func (g *WrapperGenerator) generateWrapperReturnVariables(function *Function) string {
	if function.IsReturnVoid {
		panic(fmt.Sprintf("generateWrapperReturnVariables with void function. function [%v]", function.Name))
	}

	var results []string
	for _, i := range function.Results {
		results = append(results, fmt.Sprintf("%s", i.Name))
	}

	return fmt.Sprintf("	return %s\n", strings.Join(results, ", "))
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

	return fmt.Sprintf("	return c.client.%s(%s)\n", function.Name, strings.Join(params, ", "))
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

	return fmt.Sprintf("	c.client.%s(%s)\n", function.Name, strings.Join(params, ", "))
}

func (g *WrapperGenerator) generateWrapperFunctionBody(function *Function) string {
	var buf bytes.Buffer

	buf.WriteString(g.generateWrapperOpentracing(function))

	if function.IsReturnError {
		buf.WriteString(g.generateWrapperDeclareReturnVariables(function))
	}

	if function.IsReturnVoid {
		buf.WriteString(g.generateWrapperReturnVoid(function))
	} else if function.IsReturnError {
		buf.WriteString(g.generateWrapperRetry(function))
		buf.WriteString(g.generateWrapperReturnVariables(function))
	} else {
		buf.WriteString(g.generateWrapperReturnFunction(function))
	}

	return buf.String()
}
