package astx

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"regexp"
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

	IsReturnVoid  bool
	IsReturnError bool
	IsMethod      bool
	IsChain       bool
	Class         string
}

func calculateTypeName(fset *token.FileSet, field *ast.Field, typeRegexMap map[string]*regexp.Regexp, pkg string) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, field.Type)
	if err != nil {
		return "", errors.Wrap(err, "printer.Fprint failed")
	}

	t := buf.String()

	for key, re := range typeRegexMap {
		t = re.ReplaceAllString(t, fmt.Sprintf("%s.%s", pkg, key))
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

	typeRegexMap := map[string]*regexp.Regexp{}
	for _, decl := range genDecls {
		for _, spec := range decl.Specs {
			switch spec := spec.(type) {
			case *ast.TypeSpec:
				typeRegexMap[spec.Name.String()] = regexp.MustCompile(fmt.Sprintf("\\b%s\\b", spec.Name.String()))
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
			t, err := calculateTypeName(fset, fn.Recv.List[0], typeRegexMap, pkg)
			if err != nil {
				return nil, errors.Wrap(err, "printer.Fprint failed")
			}
			f.Recv = &Field{
				Type: t,
				Name: fn.Recv.List[0].Names[0].String(),
			}

			f.IsMethod = true
			var buf bytes.Buffer
			_ = printer.Fprint(&buf, fset, fn.Recv.List[0].Type)
			f.Class = strings.TrimLeft(buf.String(), "*")
		}

		if fn.Type.Params != nil {
			for _, i := range fn.Type.Params.List {
				t, err := calculateTypeName(fset, i, typeRegexMap, pkg)
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
			for i, val := range fn.Type.Results.List {
				t, err := calculateTypeName(fset, val, typeRegexMap, pkg)
				if err != nil {
					return nil, errors.Wrap(err, "printer.Fprint failed")
				}

				if len(val.Names) != 0 {
					for _, name := range val.Names {
						f.Results = append(f.Results, &Field{
							Type: t,
							Name: name.String(),
						})
					}
				} else {
					f.Results = append(f.Results, &Field{
						Type: t,
						Name: fmt.Sprintf("res%v", i),
					})
				}
			}
		}

		f.IsReturnVoid = len(f.Results) == 0
		if len(f.Results) != 0 && f.Results[len(f.Results)-1].Type == "error" {
			f.Results[len(f.Results)-1].Name = "err"
			f.IsReturnError = true
		}
		f.IsChain = f.IsMethod && len(f.Results) == 1 && f.Recv.Type[0] == '*' && f.Recv.Type == f.Results[0].Type

		functions = append(functions, &f)
	}

	return functions, nil
}
