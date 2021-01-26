package astx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/strx"
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

func (r Rule) MarshalJSON() ([]byte, error) {
	var include *string
	var exclude *string
	if r.Include != nil {
		s := r.Include.String()
		include = &s
	}
	if r.Exclude != nil {
		s := r.Exclude.String()
		exclude = &s
	}

	return json.Marshal(struct {
		Include *string
		Exclude *string
	}{
		Include: include,
		Exclude: exclude,
	})
}

type WrapperGeneratorOptions struct {
	Debug                  bool     `flag:"usage: debug mode"`
	SourcePath             string   `flag:"usage: source path; default: vendor"`
	PackagePath            string   `flag:"usage: package path"`
	PackageName            string   `flag:"usage: package name"`
	OutputPackage          string   `flag:"usage: output package name; default: wrap" dft:"wrap"`
	Classes                []string `flag:"usage: classes to wrap"`
	StarClasses            []string `flag:"usage: star classes to wrap"`
	ClassPrefix            string   `flag:"usage: wrap class name prefix"`
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
		RateLimiter     map[string]Rule
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

	if options.Debug {
		fmt.Println(strx.JsonMarshalIndentSortKeys(options))
	}

	return &WrapperGenerator{
		options:           options,
		wrapClassMap:      wrapClassMap,
		starClassSet:      starClassSet,
		wrapPackagePrefix: wrapPackagePrefix,
	}
}

type RenderInfo struct {
	Debug             bool
	Package           string
	WrapPackagePrefix string
	UnwrapFunc        string
	Class             string
	WrapClass         string
	Function          struct {
		Name             string
		ErrCode          string
		ParamList        string
		ResultList       string
		LastResult       string
		DeclareVariables string
		ReturnVariables  string
		IsChain          bool
		IsReturnVoid     bool
		IsReturnError    bool
	}

	EnableRuleForChainFunc bool
	Rule                   struct {
		Trace       bool
		Retry       bool
		Metric      bool
		RateLimiter bool
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

	if g.options.Debug {
		fmt.Printf("classes to wrap: %v\n", strx.JsonMarshal(classes))
		fmt.Printf("star class set: %v\n", strx.JsonMarshal(g.starClassSet))
		fmt.Printf("wrap class map: %v\n", strx.JsonMarshalIndent(g.wrapClassMap))
	}

	info := &RenderInfo{
		Package:                g.options.PackageName,
		WrapPackagePrefix:      g.wrapPackagePrefix,
		UnwrapFunc:             g.options.UnwrapFunc,
		Debug:                  g.options.Debug,
		EnableRuleForChainFunc: g.options.EnableRuleForChainFunc,
	}

	for _, cls := range classes {
		info.Class = cls
		info.WrapClass = g.wrapClassMap[cls]

		if g.options.Debug {
			fmt.Printf("process class: %v, info: %v\n", cls, strx.JsonMarshalIndent(info))
		}
		buf.WriteString(renderTemplate(WrapperStructTpl, info, "WrapperStructTpl"))
		if g.starClassSet[cls] {
			buf.WriteString(renderTemplate(WrapperFunctionGetWithStarTpl, info, "WrapperFunctionGetWithStarTpl"))
		} else {
			buf.WriteString(renderTemplate(WrapperFunctionGetTpl, info, "WrapperFunctionGetTpl"))
		}

		if g.MatchRule(cls, g.options.Rule.OnWrapperChange) {
			buf.WriteString(renderTemplate(WrapperFunctionOnWrapperChangeTpl, info, "WrapperFunctionOnWrapperChangeTpl"))
		}
		if g.MatchRule(cls, g.options.Rule.OnRetryChange) {
			buf.WriteString(renderTemplate(WrapperFunctionOnRetryChangeTpl, info, "WrapperFunctionOnRetryChangeTpl"))
		}
		if g.MatchRule(cls, g.options.Rule.CreateMetric) {
			buf.WriteString(renderTemplate(WrapperFunctionCreateMetricTpl, info, "WrapperFunctionCreateMetricTpl"))
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

		info.Class = function.Class
		info.WrapClass = g.wrapClassMap[function.Class]
		info.Function.Name = function.Name
		info.Function.IsChain = function.IsChain
		info.Function.IsReturnVoid = function.IsReturnVoid
		info.Function.IsReturnError = function.IsReturnError
		info.Function.ErrCode = "ErrCode(err)"
		if !function.IsReturnError {
			info.Function.ErrCode = `"OK"`
		}
		var params []string
		for _, i := range function.Params {
			if strings.HasPrefix(i.Type, "...") {
				params = append(params, fmt.Sprintf("%s...", i.Name))
			} else {
				params = append(params, i.Name)
			}
		}
		info.Function.ParamList = strings.Join(params, ", ")
		var results []string
		for _, i := range function.Results {
			results = append(results, fmt.Sprintf("%s", i.Name))
		}
		info.Function.ResultList = strings.Join(results, ", ")
		if len(results) != 0 {
			info.Function.LastResult = results[len(results)-1]
		}
		info.Function.ReturnVariables = g.generateWrapperReturnVariables(function)
		info.Function.DeclareVariables = g.generateWrapperDeclareReturnVariables(function)
		info.Rule.Trace = g.MatchFunctionRule(function, g.options.Rule.Trace)
		info.Rule.Metric = g.MatchFunctionRule(function, g.options.Rule.Metric)
		info.Rule.Retry = g.MatchFunctionRule(function, g.options.Rule.Retry)
		info.Rule.RateLimiter = g.MatchFunctionRule(function, g.options.Rule.RateLimiter)

		buf.WriteString("\n")
		buf.WriteString(g.generateWrapperFunctionDeclare(function))
		buf.WriteString(" {")
		buf.WriteString(g.generateWrapperFunctionBody(info, function))
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

func renderTemplate(tplStr string, vals interface{}, tplName string) string {
	v := vals.(*RenderInfo)
	if v.Debug {
		if v.Function.Name != "" {
			fmt.Printf("render template function. class: [%v], function: [%v], tpl[%v]\n", v.Class, v.Function.Name, tplName)
		} else {
			fmt.Printf("render template class. class: [%v], tpl: [%v]\n", v.Class, tplName)
		}
	}
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
type {{.WrapClass}} struct {
	obj              *{{.Package}}.{{.Class}}
	retry            *{{.WrapPackagePrefix}}Retry
	options          *{{.WrapPackagePrefix}}WrapperOptions
	durationMetric   *prometheus.HistogramVec
	totalMetric      *prometheus.CounterVec
	rateLimiterGroup RateLimiterGroup
}
`
const WrapperFunctionGetWithStarTpl = `
func (w *{{.WrapClass}}) {{.UnwrapFunc}}() *{{.Package}}.{{.Class}} {
	return w.obj
}
`
const WrapperFunctionGetTpl = `
func (w {{.WrapClass}}) {{.UnwrapFunc}}() *{{.Package}}.{{.Class}} {
	return w.obj
}
`
const WrapperFunctionOnWrapperChangeTpl = `
func (w *{{.WrapClass}}) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options {{.WrapPackagePrefix}}WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}
`
const WrapperFunctionOnRetryChangeTpl = `
func (w *{{.WrapClass}}) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options {{.WrapPackagePrefix}}RetryOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		retry, err := {{.WrapPackagePrefix}}NewRetryWithOptions(&options)
		if err != nil {
			return errors.Wrap(err, "NewRetryWithOptions failed")
		}
		w.retry = retry
		return nil
	}
}
`
const WrapperFunctionCreateMetricTpl = `
func (w *{{.WrapClass}}) CreateMetric(options *{{.WrapPackagePrefix}}WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "{{.Package}}_{{.Class}}_durationMs",
		Help:        "{{.Package}} {{.Class}} response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.totalMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "{{.Package}}_{{.Class}}_total",
		Help:        "{{.Package}} {{.Class}} request total",
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
	var vars []string
	for _, field := range function.Results {
		vars = append(vars, fmt.Sprintf("var %s %s", field.Name, field.Type))
	}
	return strings.Join(vars, "\n")
}

func (g *WrapperGenerator) generateWrapperReturnVariables(function *Function) string {
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
				results = append(results, fmt.Sprintf(`&%s{obj: %s, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric, rateLimiterGroup: w.rateLimiterGroup}`, wrapCls, i.Name))
			} else {
				results = append(results, fmt.Sprintf(`%s{obj: %s, retry: w.retry, options: w.options, durationMetric: w.durationMetric, totalMetric: w.totalMetric, rateLimiterGroup: w.rateLimiterGroup}`, wrapCls, i.Name))
			}
			continue
		}

		results = append(results, fmt.Sprintf("%s", i.Name))
	}

	return fmt.Sprintf("return %s", strings.Join(results, ", "))
}

const WrapperFunctionBodyWithoutErrorTpl = `
{{- if .EnableRuleForChainFunc}}
{{- if or .Rule.Trace .Rule.Metric}}
	ctxOptions := FromContext(ctx)
{{- end}}
{{- if .Rule.Trace}}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "{{.Package}}.{{.Class}}.{{.Function.Name}}")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
{{- end}}
{{- if .Rule.RateLimiter}}
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "{{.Class}}.{{.Function.Name}}")
	}
{{- end}}
{{- if .Rule.Metric}}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", {{.Function.ErrCode}}, ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", {{.Function.ErrCode}}, ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
{{- end}}
{{- end}}
{{- if .Function.IsChain}}
	w.obj = w.obj.{{.Function.Name}}({{.Function.ParamList}})
	return w
{{- else if .Function.IsReturnVoid}}
	w.obj.{{.Function.Name}}({{.Function.ParamList}})
{{- else if not .Function.IsReturnError}}
	{{.Function.ResultList}} := w.obj.{{.Function.Name}}({{.Function.ParamList}})
	{{.Function.ReturnVariables}}
{{- end}}
`

const WrapperFunctionBodyWithErrorWithoutRetryTpl = `
{{- if or .Rule.Trace .Rule.Metric}}
	ctxOptions := FromContext(ctx)
{{- end}}
{{- if .Rule.Trace}}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "{{.Package}}.{{.Class}}.{{.Function.Name}}")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
{{- end}}
{{.Function.DeclareVariables}}
{{- if .Rule.RateLimiter}}
	if w.rateLimiterGroup != nil {
		if err := w.rateLimiterGroup.Wait(ctx, "{{.Class}}.{{.Function.Name}}"); err != nil {
			{{.Function.ReturnVariables}}
		}
	}
{{- end}}
{{- if .Rule.Metric}}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		defer func() {
			w.totalMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", {{.Function.ErrCode}}, ctxOptions.MetricCustomLabelValue).Inc()
			w.durationMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", {{.Function.ErrCode}}, ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
{{- end}}
	{{.Function.ResultList}} = w.obj.{{.Function.Name}}({{.Function.ParamList}})
	{{.Function.ReturnVariables}}
`

const WrapperFunctionBodyWithErrorWithRetryTpl = `
{{- if or .Rule.Trace .Rule.Metric}}
	ctxOptions := FromContext(ctx)
{{- end}}
{{- if .Rule.Trace}}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "{{.Package}}.{{.Class}}.{{.Function.Name}}")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
{{- end}}
{{.Function.DeclareVariables}}
	err = w.retry.Do(func() error {
{{- if .Rule.RateLimiter}}
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "{{.Class}}.{{.Function.Name}}"); err != nil {
				return err
			}
		}
{{- end}}
{{- if .Rule.Metric}}
		if w.options.EnableMetric {
			ts := time.Now()
			defer func() {
				w.totalMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", {{.Function.ErrCode}}, ctxOptions.MetricCustomLabelValue).Inc()
				w.durationMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", {{.Function.ErrCode}}, ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
{{- end}}
		{{.Function.ResultList}} = w.obj.{{.Function.Name}}({{.Function.ParamList}})
		return {{.Function.LastResult}}
	})
{{.Function.ReturnVariables}}
`

func (g *WrapperGenerator) generateWrapperFunctionBody(info *RenderInfo, function *Function) string {
	var buf bytes.Buffer
	if !function.IsReturnError {
		buf.WriteString(renderTemplate(WrapperFunctionBodyWithoutErrorTpl, info, "WrapperFunctionBodyWithoutErrorTpl"))
		return buf.String()
	}
	if !info.Rule.Retry {
		buf.WriteString(renderTemplate(WrapperFunctionBodyWithErrorWithoutRetryTpl, info, "WrapperFunctionBodyWithoutErrorTpl"))
		return buf.String()
	}
	buf.WriteString(renderTemplate(WrapperFunctionBodyWithErrorWithRetryTpl, info, "WrapperFunctionBodyWithErrorWithRetryTpl"))
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
