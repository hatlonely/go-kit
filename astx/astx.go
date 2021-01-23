package astx

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"regexp"
	"sort"
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

// 给包内部的结构增加包名前缀
func calculateTypeName(fset *token.FileSet, field *ast.Field, typeRegexMap map[string]*regexp.Regexp, pkg string) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, field.Type)
	if err != nil {
		return "", errors.Wrap(err, "printer.Fprint failed")
	}

	t := buf.String()

	for key, re := range typeRegexMap {
		// DB => gorm.DB
		// sql.DB => sql.DB
		// XDB => XDB
		// ...Options => ...oss.Options
		t = re.ReplaceAllStringFunc(t, func(s string) string {
			if s[0] == '.' && s[1] != '.' {
				return s
			}
			return fmt.Sprintf("%s%s.%s", s[:len(s)-len(key)], pkg, key)
		})
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
				typeRegexMap[spec.Name.String()] = regexp.MustCompile(fmt.Sprintf(`\.*\b%s\b`, spec.Name.String()))
			}
		}
	}

	var functions []*Function
	for _, fn := range funcDecls {
		// 跳过私有函数
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
			var name string
			if len(fn.Recv.List[0].Names) != 0 {
				name = fn.Recv.List[0].Names[0].String()
			}
			f.Recv = &Field{
				Type: t,
				Name: name,
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

				// 如果有返回值名字，直接使用返回值的名字，如果没有，命名为 res{i}
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
		// 如果最后一个返回值是 error，并且参数中没有 error，把返回命名为 err
		if len(f.Results) != 0 && f.Results[len(f.Results)-1].Type == "error" {
			var hasErrParam bool
			for _, param := range f.Params {
				if param.Name == "err" {
					hasErrParam = true
					break
				}
			}
			if !hasErrParam {
				f.Results[len(f.Results)-1].Name = "err"
				f.IsReturnError = true
			}
		}
		f.IsChain = f.IsMethod && len(f.Results) == 1 && f.Recv.Type[0] == '*' && f.Recv.Type == f.Results[0].Type

		functions = append(functions, &f)
	}

	// 排序，以保证生成代码是一致的
	sort.Slice(functions, func(i, j int) bool {
		if functions[i].Class < functions[j].Class {
			return true
		}
		if functions[i].Class > functions[j].Class {
			return false
		}
		if functions[i].Name < functions[j].Name {
			return true
		}
		if functions[i].Name > functions[j].Name {
			return false
		}
		return true
	})
	return functions, nil
}
