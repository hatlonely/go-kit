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
	GoPath    string
	PkgPath   string
	Package   string
	Class     string
	WrapClass string
	Output    string
}

func NewWrapperGeneratorWithOptions(options *WrapperGeneratorOptions) *WrapperGenerator {
	if options.WrapClass == "" {
		options.WrapClass = options.Class + "Wrapper"
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
		buf.WriteString(fmt.Sprintf("(%s %s)", function.Recv.Name, g.options.WrapClass))
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

func (g *WrapperGenerator) generateWrapperFunctionBody(function *Function) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(`
	span, _ := opentracing.StartSpanFromContext(ctx, "%s.%s.%s")
	defer span.Finish()
`, g.options.Package, g.options.Class, function.Name))

	if len(function.Results) != 0 && function.Results[len(function.Results)-1].Type == "error" {
		function.Results[len(function.Results)-1].Name = "err"
		for _, field := range function.Results {
			buf.WriteString(fmt.Sprintf("\n	var %s %s", field.Name, field.Type))
		}
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

	if len(results) == 0 {
		buf.WriteString(fmt.Sprintf("	c.client.%s(%s)\n", function.Name, strings.Join(params, ", ")))
	} else if function.Results[len(function.Results)-1].Type == "error" {
		buf.WriteString(fmt.Sprintf(`
	err = retry.Do(func() error {
		%s = %s.client.%s(%s)
		return %s
	})
`, strings.Join(results, ", "), function.Recv.Name, function.Name, strings.Join(params, ", "), results[len(results)-1]))
		buf.WriteString(fmt.Sprintf("	return %s\n", strings.Join(results, ", ")))
	} else {
		buf.WriteString(fmt.Sprintf("	return c.client.%s(%s)\n", function.Name, strings.Join(params, ", ")))
	}

	return buf.String()
}
