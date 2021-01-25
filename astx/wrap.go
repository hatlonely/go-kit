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
	starClassSet map[string]bool

	wrapPackagePrefix string
}

type Rule struct {
	Include *regexp.Regexp
	Exclude *regexp.Regexp
}

type WrapperGeneratorOptions struct {
	SourcePath             string   `flag:"usage: gopath; default: vendor"`
	PackagePath            string   `flag:"usage: package path"`
	PackageName            string   `flag:"usage: package name"`
	OutputPackage          string   `flag:"usage: output package name; default: wrap" dft:"wrap"`
	Classes                []string `flag:"usage: classes to wrap"`
	StarClasses            []string `flag:"usage: star classes to wrap"`
	ClassPrefix            string   `flag:"usage: wrap class name"`
	UnwrapFunc             string   `flag:"usage: unwrap function name; default: Unwrap" dft:"Unwrap"`
	EnableRuleForChainFunc bool     `flag:"usage: enable trace on chain function"`

	Rule struct {
		Class           Rule
		StarClass       Rule
		OnWrapperChange Rule
		OnRetryChange   Rule
		CreateMetric    Rule
		Function        map[string]Rule
		Trace           map[string]Rule
		Retry           map[string]Rule
		Metric          map[string]Rule
	}
}

func NewWrapperGeneratorWithOptions(options *WrapperGeneratorOptions) *WrapperGenerator {
	wrapClassMap := map[string]string{}
	starClassSet := map[string]bool{}
	for _, cls := range options.Classes {
		wrapClassMap[cls] = fmt.Sprintf("%s%sWrapper", options.ClassPrefix, cls)
	}
	for _, cls := range options.StarClasses {
		wrapClassMap[cls] = fmt.Sprintf("%s%sWrapper", options.ClassPrefix, cls)
		starClassSet[cls] = true
	}

	excludeAllRegex := regexp.MustCompile(`.*`)
	if options.Rule.OnWrapperChange.Exclude == nil {
		options.Rule.OnWrapperChange.Exclude = excludeAllRegex
	}
	if options.Rule.OnRetryChange.Exclude == nil {
		options.Rule.OnRetryChange.Exclude = excludeAllRegex
	}
	if options.Rule.StarClass.Exclude == nil {
		options.Rule.StarClass.Exclude = excludeAllRegex
	}
	if options.Rule.Class.Exclude == nil {
		options.Rule.Class.Exclude = excludeAllRegex
	}
	if options.Rule.CreateMetric.Exclude == nil {
		options.Rule.CreateMetric.Exclude = excludeAllRegex
	}

	var wrapPackagePrefix string
	if options.OutputPackage != "wrap" {
		wrapPackagePrefix = "wrap."
	}

	return &WrapperGenerator{
		options:           options,
		wrapClassMap:      wrapClassMap,
		starClassSet:      starClassSet,
		wrapPackagePrefix: wrapPackagePrefix,
	}
}

func (g *WrapperGenerator) Generate() (string, error) {
	functions, err := ParseFunction(path.Join(g.options.SourcePath, g.options.PackagePath), g.options.PackageName)
	if err != nil {
		return "", errors.Wrap(err, "ParseFunction failed")
	}

	var buf bytes.Buffer

	buf.WriteString(g.generateWrapperHeader())

	classes := append(g.options.Classes, g.options.StarClasses...)
	if len(classes) == 0 {
		for _, function := range functions {
			if !function.IsMethod {
				continue
			}
			if _, ok := g.wrapClassMap[function.Class]; ok {
				continue
			}
			if g.MatchRule(function.Class, g.options.Rule.StarClass) {
				g.wrapClassMap[function.Class] = fmt.Sprintf("%s%sWrapper", g.options.ClassPrefix, function.Class)
				classes = append(classes, function.Class)
				g.starClassSet[function.Class] = true
				continue
			}
			if g.MatchRule(function.Class, g.options.Rule.Class) {
				g.wrapClassMap[function.Class] = fmt.Sprintf("%s%sWrapper", g.options.ClassPrefix, function.Class)
				classes = append(classes, function.Class)
				continue
			}
		}
	}

	sort.Strings(classes)

	vals := map[string]interface{}{
		"package":           g.options.PackageName,
		"wrapPackagePrefix": g.wrapPackagePrefix,
		"unwrapFunc":        g.options.UnwrapFunc,
	}

	for _, cls := range classes {
		vals["class"] = cls
		vals["wrapClass"] = g.wrapClassMap[cls]

		buf.WriteString(g.generateWrapperStruct(vals))
		buf.WriteString(g.generateWrapperGet(vals))

		if g.MatchRule(cls, g.options.Rule.OnWrapperChange) {
			buf.WriteString(g.generateWrapperOnWrapperChange(vals))
		}
		if g.MatchRule(cls, g.options.Rule.OnRetryChange) {
			buf.WriteString(g.generateWrapperOnRetryChange(vals))
		}
		if g.MatchRule(cls, g.options.Rule.CreateMetric) {
			buf.WriteString(g.generateWrapperCreateMetric(vals))
		}
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
	var buf bytes.Buffer

	buf.WriteString("// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!\n")
	buf.WriteString(fmt.Sprintf("package %s\n\n", g.options.OutputPackage))
	buf.WriteString("import (\n")
	buf.WriteString("\t\"context\"\n")
	buf.WriteString("")
	buf.WriteString(fmt.Sprintf("\t\"%s\"\n", g.options.PackagePath))
	buf.WriteString(fmt.Sprintf("\t\"github.com/opentracing/opentracing-go\"\n\n"))

	if g.options.OutputPackage != "wrap" {
		buf.WriteString(fmt.Sprintf("\t\"github.com/hatlonely/go-kit/wrap\"\n"))
	}
	if g.options.Rule.OnWrapperChange.Include != nil || g.options.Rule.OnRetryChange.Include != nil {
		buf.WriteString(fmt.Sprintf("\t\"github.com/hatlonely/go-kit/config\"\n"))
		buf.WriteString(fmt.Sprintf("\t\"github.com/hatlonely/go-kit/rpcx\"\n"))
	}

	buf.WriteString(")\n")

	return buf.String()
}

func (g *WrapperGenerator) generateWrapperStruct(vals map[string]interface{}) string {
	const tplStr = `
type {{.wrapClass}} struct {
	obj            *{{.package}}.{{.class}}
	retry          *{{.wrapPackagePrefix}}Retry
	options        *{{.wrapPackagePrefix}}WrapperOptions
	durationMetric *prometheus.HistogramVec
	totalMetric    *prometheus.CounterVec
}
`
	tpl, _ := template.New("").Parse(tplStr)

	var buf bytes.Buffer
	_ = tpl.Execute(&buf, vals)
	return buf.String()
}

func (g *WrapperGenerator) generateWrapperGet(vals map[string]interface{}) string {
	const tplStr1 = `
func (w *{{.wrapClass}}) {{.unwrapFunc}}() *{{.package}}.{{.class}} {
	return w.obj
}
`
	const tplStr2 = `
func (w {{.wrapClass}}) {{.unwrapFunc}}() *{{.package}}.{{.class}} {
	return w.obj
}
`

	var tpl *template.Template
	if g.starClassSet[vals["class"].(string)] {
		tpl, _ = template.New("").Parse(tplStr1)
	} else {
		tpl, _ = template.New("").Parse(tplStr2)
	}

	var buf bytes.Buffer
	_ = tpl.Execute(&buf, vals)
	return buf.String()
}

func (g *WrapperGenerator) generateWrapperOnWrapperChange(vals map[string]interface{}) string {
	const tplStr = `
func (w *{{.wrapClass}}) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options {{.wrapPackagePrefix}}WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}
`

	tpl, _ := template.New("").Parse(tplStr)
	var buf bytes.Buffer
	_ = tpl.Execute(&buf, vals)
	return buf.String()
}

func (g *WrapperGenerator) generateWrapperOnRetryChange(vals map[string]interface{}) string {
	const tplStr = `
func (w *{{.wrapClass}}) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options {{.wrapPackagePrefix}}RetryOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		retry, err := {{.wrapPackagePrefix}}NewRetryWithOptions(&options)
		if err != nil {
			return errors.Wrap(err, "NewRetryWithOptions failed")
		}
		w.retry = retry
		return nil
	}
}
`

	tpl, _ := template.New("").Parse(tplStr)
	var buf bytes.Buffer
	_ = tpl.Execute(&buf, vals)
	return buf.String()
}

func (g *WrapperGenerator) generateWrapperCreateMetric(vals map[string]interface{}) string {
	const tplStr = `
func (w *{{.wrapClass}}) CreateMetric(options *{{.wrapPackagePrefix}}WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "{{.package}}_{{.class}}_durationMs",
		Help:        "{{.package}} {{.class}} response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "{{.package}}_{{.class}}_total",
		Help:        "{{.package}} {{.class}} request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode"})
}
`

	tpl, _ := template.New("").Parse(tplStr)
	var buf bytes.Buffer
	_ = tpl.Execute(&buf, vals)
	return buf.String()
}

func (g *WrapperGenerator) generateWrapperFunctionDeclare(function *Function) string {
	var buf bytes.Buffer

	buf.WriteString("func ")
	if function.Recv != nil {
		if g.starClassSet[function.Class] {
			buf.WriteString(fmt.Sprintf("(w *%s)", g.wrapClassMap[function.Class]))
		} else {
			buf.WriteString(fmt.Sprintf("(w %s)", g.wrapClassMap[function.Class]))
		}
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
		if !strings.HasPrefix(cls, g.options.PackageName) {
			results = append(results, i.Type)
			continue
		}

		cls = strings.TrimPrefix(cls, g.options.PackageName+".")
		if wrapCls, ok := g.wrapClassMap[cls]; ok {
			if g.starClassSet[cls] {
				results = append(results, fmt.Sprintf(`*%s`, wrapCls))
			} else {
				results = append(results, fmt.Sprintf(`%s`, wrapCls))
			}
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
`, g.options.PackageName, function.Class, function.Name)
}

func (g *WrapperGenerator) generateWrapperMetric(function *Function) string {
	const tplStr = `
	if w.options.EnableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("{{.package}}.{{.class}}.{{.function.name}}", ErrCode(err)).Inc()
			w.durationMetric.WithLabelValues("{{.package}}.{{.class}}.{{.function.name}}", ErrCode(err)).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
`

	tpl, _ := template.New("").Parse(tplStr)

	var buf bytes.Buffer
	_ = tpl.Execute(&buf, map[string]interface{}{
		"package":           g.options.PackageName,
		"class":             function.Class,
		"wrapPackagePrefix": g.wrapPackagePrefix,
		"function": map[string]string{
			"class": function.Class,
			"name":  function.Name,
		},
	})

	return buf.String()
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
		if !strings.HasPrefix(cls, g.options.PackageName) {
			results = append(results, fmt.Sprintf("%s", i.Name))
			continue
		}

		cls = strings.TrimPrefix(cls, g.options.PackageName+".")
		if wrapCls, ok := g.wrapClassMap[cls]; ok {
			if g.starClassSet[cls] {
				results = append(results, fmt.Sprintf(`&%s{obj: %s, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric}`, wrapCls, i.Name))
			} else {
				results = append(results, fmt.Sprintf(`%s{obj: %s, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric}`, wrapCls, i.Name))
			}
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
		if g.options.EnableRuleForChainFunc {
			if g.MatchFunctionRule(function, g.options.Rule.Trace) {
				buf.WriteString(g.generateWrapperOpentracing(function))
				buf.WriteString("\n")
			}
		}
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
