package astx

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
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
						Name: strx.CamelName(strings.TrimLeft(kvs[len(kvs)-1], "*[]")),
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
						Name: strx.CamelName(strings.TrimLeft(kvs[len(kvs)-1], "*[]")),
					})
				}
			}
		}

		functions = append(functions, &f)
	}

	return functions, nil
}
