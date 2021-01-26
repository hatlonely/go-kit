package astx

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"text/template"

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

	buf.WriteString(g.generateWrapperImport())

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

		buf.WriteString(renderTemplate(WrapperStructTpl, vals))
		if g.starClassSet[cls] {
			buf.WriteString(renderTemplate(WrapperFunctionGetWithStarTpl, vals))
		} else {
			buf.WriteString(renderTemplate(WrapperFunctionGetTpl, vals))
		}

		if g.MatchRule(cls, g.options.Rule.OnWrapperChange) {
			buf.WriteString(renderTemplate(WrapperFunctionOnWrapperChangeTpl, vals))
		}
		if g.MatchRule(cls, g.options.Rule.OnRetryChange) {
			buf.WriteString(renderTemplate(WrapperFunctionOnRetryChangeTpl, vals))
		}
		if g.MatchRule(cls, g.options.Rule.CreateMetric) {
			buf.WriteString(renderTemplate(WrapperFunctionCreateMetricTpl, vals))
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

		vals["class"] = function.Class
		vals["wrapClass"] = g.wrapClassMap[function.Class]
		fmap := map[string]string{
			"name":    function.Name,
			"errCode": "ErrCode(err)",
		}

		if !function.IsReturnError {
			fmap["errCode"] = `"OK"`
		}

		var params []string
		for _, i := range function.Params {
			if strings.HasPrefix(i.Type, "...") {
				params = append(params, fmt.Sprintf("%s...", i.Name))
			} else {
				params = append(params, i.Name)
			}
		}
		fmap["paramList"] = strings.Join(params, ", ")
		var results []string
		for _, i := range function.Results {
			results = append(results, fmt.Sprintf("%s", i.Name))
		}
		fmap["resultList"] = strings.Join(results, ", ")
		if len(results) != 0 {
			fmap["lastResult"] = results[len(results)-1]
		}

		vals["function"] = fmap
		buf.WriteString("\n")
		buf.WriteString(g.generateWrapperFunctionDeclare(function))
		buf.WriteString(" {")
		buf.WriteString(g.generateWrapperFunctionBody(vals, function))
		buf.WriteString("}\n")
	}

	return buf.String(), nil
}

func (g *WrapperGenerator) generateWrapperImport() string {
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

func renderTemplate(tplStr string, vals interface{}) string {
	tpl, err := template.New("").Parse(tplStr)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, vals); err != nil {
		panic(err)
	}

	return buf.String()
}

const WrapperStructTpl = `
type {{.wrapClass}} struct {
	obj            *{{.package}}.{{.class}}
	retry          *{{.wrapPackagePrefix}}Retry
	options        *{{.wrapPackagePrefix}}WrapperOptions
	durationMetric *prometheus.HistogramVec
	totalMetric    *prometheus.CounterVec
}
`
const WrapperFunctionGetWithStarTpl = `
func (w *{{.wrapClass}}) {{.unwrapFunc}}() *{{.package}}.{{.class}} {
	return w.obj
}
`
const WrapperFunctionGetTpl = `
func (w {{.wrapClass}}) {{.unwrapFunc}}() *{{.package}}.{{.class}} {
	return w.obj
}
`
const WrapperFunctionOnWrapperChangeTpl = `
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
const WrapperFunctionOnRetryChangeTpl = `
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
const WrapperFunctionCreateMetricTpl = `
func (w *{{.wrapClass}}) CreateMetric(options *{{.wrapPackagePrefix}}WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "{{.package}}_{{.class}}_durationMs",
		Help:        "{{.package}} {{.class}} response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "{{.package}}_{{.class}}_total",
		Help:        "{{.package}} {{.class}} request total",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
}
`

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

func (g *WrapperGenerator) generateWrapperDeclareReturnVariables(function *Function) string {
	var buf bytes.Buffer
	for _, field := range function.Results {
		buf.WriteString(fmt.Sprintf("\n	var %s %s", field.Name, field.Type))
	}
	return buf.String()
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

const WrapperFunctionBodyOpentracingTpl = `
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "{{.package}}.{{.class}}.{{.function.name}}")
		defer span.Finish()
	}
`
const WrapperFunctionBodyMetricTpl = `
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("{{.package}}.{{.class}}.{{.function.name}}", {{.function.errCode}}, ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("{{.package}}.{{.class}}.{{.function.name}}", {{.function.errCode}}, ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
`
const WrapperFunctionBodyRetryTpl = `
	err = w.retry.Do(func() error {
		{{.function.resultList}} = w.obj.{{.function.name}}({{.function.paramList}})
		return {{.function.lastResult}}
	})
`
const WrapperFunctionBodyRetryWithMetricTpl = `
	err = w.retry.Do(func() error {
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("{{.package}}.{{.class}}.{{.function.name}}", {{.function.errCode}}, ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("{{.package}}.{{.class}}.{{.function.name}}", {{.function.errCode}}, ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}

		{{.function.resultList}} = w.obj.{{.function.name}}({{.function.paramList}})
		return {{.function.lastResult}}
	})
`
const WrapperFunctionBodyCallTpl = `
	{{.function.resultList}} = w.obj.{{.function.name}}({{.function.paramList}})
`

const WrapperFunctionBodyReturnVoidTpl = `
	w.obj.{{.function.name}}({{.function.paramList}})
`
const WrapperFunctionBodyReturnChainTpl = `
	w.obj = w.obj.{{.function.name}}({{.function.paramList}})
	return w
`

func (g *WrapperGenerator) generateWrapperFunctionBody(vals map[string]interface{}, function *Function) string {
	var buf bytes.Buffer
	if function.IsChain {
		if g.options.EnableRuleForChainFunc {
			if g.MatchFunctionRule(function, g.options.Rule.Trace) || g.MatchFunctionRule(function, g.options.Rule.Metric) {
				buf.WriteString("\tctxOptions := FromContext(ctx)\n")
			}
			if g.MatchFunctionRule(function, g.options.Rule.Trace) {
				buf.WriteString(renderTemplate(WrapperFunctionBodyOpentracingTpl, vals))
			}
			if g.MatchFunctionRule(function, g.options.Rule.Metric) {
				buf.WriteString(renderTemplate(WrapperFunctionBodyMetricTpl, vals))
			}
		}
		buf.WriteString(renderTemplate(WrapperFunctionBodyReturnChainTpl, vals))
		return buf.String()
	}

	if g.MatchFunctionRule(function, g.options.Rule.Trace) || g.MatchFunctionRule(function, g.options.Rule.Metric) {
		buf.WriteString("\tctxOptions := FromContext(ctx)\n")
	}

	if g.MatchFunctionRule(function, g.options.Rule.Trace) {
		buf.WriteString(renderTemplate(WrapperFunctionBodyOpentracingTpl, vals))
	}

	if function.IsReturnVoid {
		buf.WriteString(renderTemplate(WrapperFunctionBodyReturnVoidTpl, vals))
		return buf.String()
	}

	buf.WriteString(g.generateWrapperDeclareReturnVariables(function))

	if function.IsReturnError && g.MatchFunctionRule(function, g.options.Rule.Retry) {
		if g.MatchFunctionRule(function, g.options.Rule.Metric) {
			buf.WriteString(renderTemplate(WrapperFunctionBodyRetryWithMetricTpl, vals))
		} else {
			buf.WriteString(renderTemplate(WrapperFunctionBodyRetryTpl, vals))
		}
		buf.WriteString(g.generateWrapperReturnVariables(function))
		return buf.String()
	}

	if g.MatchFunctionRule(function, g.options.Rule.Metric) {
		buf.WriteString(renderTemplate(WrapperFunctionBodyMetricTpl, vals))
	}
	buf.WriteString(renderTemplate(WrapperFunctionBodyCallTpl, vals))
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
