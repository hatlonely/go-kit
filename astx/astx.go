package astx

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"html/template"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/strx"
)

type Field struct {
	Type string
	Name string
}

type Function struct {
	Name    string
	Recv    *Field
	Params  []*Field
	Results []*Field
}

func calculateTypeName(fset *token.FileSet, field *ast.Field, typeset map[string]bool, pkg string) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, field.Type)
	if err != nil {
		return "", errors.Wrap(err, "printer.Fprint failed")
	}

	t := buf.String()
	if strings.HasPrefix(t, "*") {
		if typeset[t[1:]] {
			return fmt.Sprintf("*%s.%s", pkg, t[1:]), nil
		}
	} else {
		if typeset[t] {
			return fmt.Sprintf("%s.%s", pkg, t), nil
		}
	}

	return t, nil
}

func ParseFunction(path string, pkg string) ([]*Function, error) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, path, nil, 0)
	if err != nil {
		return nil, errors.Wrap(err, "parser.ParseDir failed")
	}

	var funcDecls []*ast.FuncDecl
	var genDecls []*ast.GenDecl
	for _, pack := range packages {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				switch decl := d.(type) {
				case *ast.FuncDecl:
					funcDecls = append(funcDecls, decl)
				case *ast.GenDecl:
					genDecls = append(genDecls, decl)
				}
			}
		}
	}

	typeset := map[string]bool{}
	for _, decl := range genDecls {
		for _, spec := range decl.Specs {
			switch spec := spec.(type) {
			case *ast.TypeSpec:
				typeset[spec.Name.String()] = true
			}
		}
	}

	var functions []*Function
	for _, fn := range funcDecls {
		if strx.IsLower(fn.Name.String()[0]) {
			continue
		}
		var f Function
		f.Name = fn.Name.String()
		if fn.Recv != nil {
			t, err := calculateTypeName(fset, fn.Recv.List[0], typeset, pkg)
			if err != nil {
				return nil, errors.Wrap(err, "printer.Fprint failed")
			}
			f.Recv = &Field{
				Type: t,
				Name: fn.Recv.List[0].Names[0].String(),
			}
		}

		if fn.Type.Params != nil {
			for _, i := range fn.Type.Params.List {
				t, err := calculateTypeName(fset, i, typeset, pkg)
				if err != nil {
					return nil, errors.Wrap(err, "printer.Fprint failed")
				}

				if len(i.Names) != 0 {
					for _, name := range i.Names {
						f.Params = append(f.Params, &Field{
							Type: t,
							Name: name.String(),
						})
					}
				} else {
					kvs := strings.Split(t, ".")
					f.Params = append(f.Params, &Field{
						Type: t,
						Name: strx.CamelName(strings.TrimLeft(kvs[len(kvs)-1], "*")),
					})
				}
			}
		}

		if fn.Type.Results != nil {
			for _, i := range fn.Type.Results.List {
				t, err := calculateTypeName(fset, i, typeset, pkg)
				if err != nil {
					return nil, errors.Wrap(err, "printer.Fprint failed")
				}

				if len(i.Names) != 0 {
					for _, name := range i.Names {
						f.Results = append(f.Results, &Field{
							Type: t,
							Name: name.String(),
						})
					}
				} else {
					kvs := strings.Split(t, ".")
					f.Results = append(f.Results, &Field{
						Type: t,
						Name: strx.CamelName(strings.TrimLeft(kvs[len(kvs)-1], "*")),
					})
				}
			}
		}

		functions = append(functions, &f)
	}

	return functions, nil
}

func GenerateWrapper(goPath string, pkgPath string, pkg string, cls string) (string, error) {
	functions, err := ParseFunction(path.Join(goPath, pkgPath), pkg)
	if err != nil {
		return "", errors.Wrap(err, "ParseFunction failed")
	}

	var buf bytes.Buffer

	buf.WriteString(generateWrapperHeader(pkgPath))
	buf.WriteString(generateWrapperStruct(pkg, cls))

	for _, function := range functions {
		if function.Recv == nil {
			continue
		}
		if function.Recv.Type != fmt.Sprintf("%s.%s", pkg, cls) &&
			function.Recv.Type != fmt.Sprintf("*%s.%s", pkg, cls) {
			continue
		}

		buf.WriteString("\n")
		buf.WriteString(generateWrapperFunctionDeclare(function, cls))
		buf.WriteString(" {")
		buf.WriteString(generateWrapperFunctionBody(function, pkg, cls))
		buf.WriteString("}\n")
	}

	return buf.String(), nil
}

func generateWrapperHeader(pkgPath string) string {
	const tplStr = `
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
		"pkgPath": pkgPath,
	})

	return buf.String()
}

func generateWrapperStruct(pkg string, cls string) string {
	const tplStr = `
type {{.cls}}Wrapper struct {
	client *{{.pkg}}.{{.cls}}
	retry  *Retry
}
`

	tpl, _ := template.New("").Parse(tplStr)

	var buf bytes.Buffer
	_ = tpl.Execute(&buf, map[string]string{
		"pkg": pkg,
		"cls": cls,
	})

	return buf.String()
}

func generateWrapperFunctionDeclare(function *Function, cls string) string {
	var buf bytes.Buffer

	buf.WriteString("func ")
	if function.Recv != nil {
		buf.WriteString(fmt.Sprintf("(%s %sWrapper)", function.Recv.Name, cls))
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
	buf.WriteString(")")

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

func generateWrapperFunctionBody(function *Function, pkg string, cls string) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(`
	span, _ := opentracing.StartSpanFromContext(ctx, "%s.%s.%s")
	defer span.Finish()
`, pkg, cls, function.Name))

	if len(function.Results) != 0 && function.Results[len(function.Results)-1].Type == "error" {
		function.Results[len(function.Results)-1].Name = "err"
	}

	for _, field := range function.Results {
		buf.WriteString(fmt.Sprintf("\n	var %s %s", field.Name, field.Type))
	}

	var params []string
	for _, i := range function.Params {
		params = append(params, fmt.Sprintf("%s", i.Name))
	}

	var results []string
	for _, i := range function.Results {
		results = append(results, fmt.Sprintf("%s", i.Name))
	}

	if len(results) != 0 && function.Results[len(function.Results)-1].Type == "error" {
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
