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
	Debug                    bool                `flag:"usage: debug mode"`
	SourcePath               string              `flag:"usage: source path; default: vendor"`
	PackagePath              string              `flag:"usage: package path"`
	PackageName              string              `flag:"usage: package name"`
	OutputPackage            string              `flag:"usage: output package name; default: wrap" dft:"wrap"`
	Classes                  []string            `flag:"usage: classes to wrap"`
	StarClasses              []string            `flag:"usage: star classes to wrap"`
	ClassPrefix              string              `flag:"usage: wrap class name prefix"`
	UnwrapFunc               string              `flag:"usage: unwrap function name; default: Unwrap" dft:"Unwrap"`
	ErrorField               string              `flag:"usage: function return no error, error is a filed in result"`
	Inherit                  map[string][]string `flag:"usage: inherit map"`
	EnableRuleForNoErrorFunc bool                `flag:"usage: enable trace for no error function"`
	EnableHystrix            bool                `flag:"usage: enable hystrix code"`

	Rule struct {
		Class                    Rule
		StarClass                Rule
		OnWrapperChange          Rule
		OnRetryChange            Rule
		OnRateLimiterGroupChange Rule
		CreateMetric             Rule
		ErrorInResult            Rule
		Function                 map[string]Rule
		Trace                    map[string]Rule
		Retry                    map[string]Rule
		Metric                   map[string]Rule
		RateLimiter              map[string]Rule
		Hystrix                  map[string]Rule
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
	if options.Rule.OnRateLimiterGroupChange.Exclude == nil {
		options.Rule.OnRateLimiterGroupChange.Exclude = excludeAllRegex
	}
	if options.Rule.CreateMetric.Exclude == nil {
		options.Rule.CreateMetric.Exclude = excludeAllRegex
	}
	if options.Rule.StarClass.Exclude == nil {
		options.Rule.StarClass.Exclude = excludeAllRegex
	}
	if options.Rule.Class.Exclude == nil {
		options.Rule.Class.Exclude = excludeAllRegex
	}
	if options.Rule.ErrorInResult.Exclude == nil {
		options.Rule.ErrorInResult.Exclude = excludeAllRegex
	}
	if len(options.Rule.Hystrix) == 0 {
		options.Rule.Hystrix = map[string]Rule{
			"default": {
				Exclude: regexp.MustCompile(`^.*$`),
			},
		}
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
	OutputPackage     string
	PackagePath       string
	ErrorField        string
	Interface         string
	EnableHystrix     bool
	IsStarClass       bool
	NewLine           string
	Function          struct {
		Name             string
		ErrCode          string
		ParamList        string
		ResultList       string
		LastResult       string
		DeclareVariables string
		ReturnList       string
		IsChain          bool
		IsReturnVoid     bool
		IsReturnError    bool
	}

	EnableRuleForNoErrorFunc bool
	Rule                     struct {
		OnWrapperChange          bool
		OnRetryChange            bool
		OnRateLimiterGroupChange bool
		CreateMetric             bool
		Trace                    bool
		Retry                    bool
		Metric                   bool
		Hystrix                  bool
		RateLimiter              bool
		ErrorInResult            bool
	}
}

func (g *WrapperGenerator) Generate() (string, error) {
	functions, err := ParseFunction(path.Join(g.options.SourcePath, g.options.PackagePath), g.options.PackageName)
	if err != nil {
		return "", errors.Wrap(err, "ParseFunction failed")
	}

	var buf bytes.Buffer

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
		Package:                  g.options.PackageName,
		WrapPackagePrefix:        g.wrapPackagePrefix,
		UnwrapFunc:               g.options.UnwrapFunc,
		Debug:                    g.options.Debug,
		EnableRuleForNoErrorFunc: g.options.EnableRuleForNoErrorFunc,
		OutputPackage:            g.options.OutputPackage,
		PackagePath:              g.options.PackagePath,
		ErrorField:               g.options.ErrorField,
		EnableHystrix:            g.options.EnableHystrix,
		NewLine:                  "\n",
	}

	buf.WriteString(renderTemplate(WrapperImportTpl, info, "WrapperImportTpl"))

	for _, cls := range classes {
		info.Class = cls
		info.WrapClass = g.wrapClassMap[cls]
		info.Rule.OnWrapperChange = g.MatchRule(cls, g.options.Rule.OnWrapperChange)
		info.Rule.OnRetryChange = g.MatchRule(cls, g.options.Rule.OnRetryChange)
		info.Rule.OnRateLimiterGroupChange = g.MatchRule(cls, g.options.Rule.OnRateLimiterGroupChange)
		info.Rule.CreateMetric = g.MatchRule(cls, g.options.Rule.CreateMetric)
		info.IsStarClass = g.starClassSet[cls]
		info.Interface = fmt.Sprintf("I%s%s", g.options.ClassPrefix, cls)

		if g.options.Debug {
			fmt.Printf("process class: %v, info: %v\n", cls, strx.JsonMarshalIndent(info))
		}

		buf.WriteString(renderTemplate(WrapperClassTpl, info, "WrapperClassTpl"))
	}

	inheritMap := map[string][]string{}
	for key, vals := range g.options.Inherit {
		for _, val := range vals {
			inheritMap[val] = append(inheritMap[val], key)
		}
	}

	traveled := map[string]bool{}
	for _, cls := range classes {
		clsset := map[string]bool{
			cls: true,
		}
		for _, c := range g.options.Inherit[cls] {
			clsset[c] = true
		}

		info.Interface = fmt.Sprintf("I%s%s", g.options.ClassPrefix, cls)
		var interfaceBuffer bytes.Buffer
		interfaceBuffer.WriteString(fmt.Sprintf("\ntype %s interface {", info.Interface))
		for _, function := range functions {
			if !function.IsMethod {
				continue
			}
			if !clsset[function.Class] {
				continue
			}
			if traveled[cls+function.Name] {
				continue
			}
			traveled[cls+function.Name] = true
			if !g.MatchFunctionRule(function.Name, cls, g.options.Rule.Function) {
				continue
			}

			var paramTypes []string
			for _, i := range function.Params {
				paramTypes = append(paramTypes, i.Type)
			}
			var resultTypes []string
			for _, i := range function.Results {
				resultTypes = append(resultTypes, i.Type)
			}
			if len(resultTypes) >= 2 {
				interfaceBuffer.WriteString(fmt.Sprintf("\t%s(%s) (%s)\n", function.Name, strings.Join(paramTypes, ", "), strings.Join(resultTypes, ",")))
			} else {
				interfaceBuffer.WriteString(fmt.Sprintf("\t%s(%s) %s\n", function.Name, strings.Join(paramTypes, ", "), strings.Join(resultTypes, ",")))
			}

			info.Class = cls
			info.WrapClass = g.wrapClassMap[cls]
			info.Function.Name = function.Name
			info.Function.IsChain = function.IsChain
			info.Function.IsReturnVoid = function.IsReturnVoid
			info.Function.IsReturnError = function.IsReturnError
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
			info.Function.ReturnList = g.generateWrapperReturnList(function)
			info.Function.DeclareVariables = g.generateWrapperDeclareReturnVariables(function)
			info.Rule.OnWrapperChange = g.MatchRule(cls, g.options.Rule.OnWrapperChange)
			info.Rule.OnRetryChange = g.MatchRule(cls, g.options.Rule.OnRetryChange)
			info.Rule.CreateMetric = g.MatchRule(cls, g.options.Rule.CreateMetric)
			info.Rule.OnRateLimiterGroupChange = g.MatchRule(cls, g.options.Rule.OnRateLimiterGroupChange)
			info.Rule.Trace = g.MatchFunctionRule(function.Name, cls, g.options.Rule.Trace)
			info.Rule.Metric = g.MatchFunctionRule(function.Name, cls, g.options.Rule.Metric)
			info.Rule.Retry = g.MatchFunctionRule(function.Name, cls, g.options.Rule.Retry)
			info.Rule.RateLimiter = g.MatchFunctionRule(function.Name, cls, g.options.Rule.RateLimiter)
			info.Rule.Hystrix = g.MatchFunctionRule(function.Name, cls, g.options.Rule.Hystrix)
			if len(function.Results) == 1 {
				info.Rule.ErrorInResult = g.MatchRule(function.Results[0].Type, g.options.Rule.ErrorInResult)
			} else {
				info.Rule.ErrorInResult = false
			}

			if !function.IsReturnError && !(info.Rule.ErrorInResult && len(function.Results) == 1) {
				info.Function.ErrCode = `"OK"`
			} else if info.OutputPackage == "wrap" {
				if info.Rule.ErrorInResult && len(function.Results) == 1 {
					info.Function.ErrCode = fmt.Sprintf("ErrCode(res0.%s)", info.ErrorField)
				} else {
					info.Function.ErrCode = "ErrCode(err)"
				}
			} else {
				if info.Rule.ErrorInResult && len(function.Results) == 1 {
					info.Function.ErrCode = fmt.Sprintf("wrap.ErrCode(res0.%s)", info.ErrorField)
				} else {
					info.Function.ErrCode = "wrap.ErrCode(err)"
				}
			}

			buf.WriteString("\n")
			buf.WriteString(g.generateWrapperFunctionDeclare(info, function))
			buf.WriteString(" {")
			buf.WriteString(g.generateWrapperFunctionBody(info, function))
			buf.WriteString("}\n")
		}
		interfaceBuffer.WriteString("}\n")

		if info.EnableHystrix {
			buf.WriteString(interfaceBuffer.String())
		}
	}

	return buf.String(), nil
}

const WrapperImportTpl = `
// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package {{.OutputPackage}}

import (
	"context"
	"time"

	"{{.PackagePath}}"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hatlonely/go-kit/config"
{{- if not (eq .OutputPackage "wrap")}}
	"github.com/hatlonely/go-kit/wrap"
{{- end}}
	"github.com/hatlonely/go-kit/refx"
)
`

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

const WrapperClassTpl = `
type {{.WrapClass}} struct {
	obj              *{{.Package}}.{{.Class}}
{{- if .EnableHystrix}}
	backupObj        {{.Interface}}
{{- end}}
	retry            *{{.WrapPackagePrefix}}Retry
	options          *{{.WrapPackagePrefix}}WrapperOptions
	durationMetric   *prometheus.HistogramVec
	inflightMetric   *prometheus.GaugeVec
	rateLimiterGroup {{.WrapPackagePrefix}}RateLimiterGroup
}

{{- if .IsStarClass}}
func (w *{{.WrapClass}}) {{.UnwrapFunc}}() *{{.Package}}.{{.Class}} {
	return w.obj
}
{{- else}}
func (w {{.WrapClass}}) {{.UnwrapFunc}}() *{{.Package}}.{{.Class}} {
	return w.obj
}
{{- end}}


{{- if .Rule.OnWrapperChange}}

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
{{- end}}


{{- if .Rule.OnRetryChange}}

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
{{- end}}

{{- if .Rule.OnRateLimiterGroupChange}}

func (w *{{.WrapClass}}) OnRateLimiterGroupChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options {{.WrapPackagePrefix}}RateLimiterGroupOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		rateLimiterGroup, err := {{.WrapPackagePrefix}}NewRateLimiterGroupWithOptions(&options, opts...)
		if err != nil {
			return errors.Wrap(err, "NewRateLimiterGroupWithOptions failed")
		}
		w.rateLimiterGroup = rateLimiterGroup
		return nil
	}
}
{{- end}}

{{- if .Rule.CreateMetric}}

func (w *{{.WrapClass}}) CreateMetric(options *{{.WrapPackagePrefix}}WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "{{.Package}}_{{.Class}}_durationMs",
		Help:        "{{.Package}} {{.Class}} response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "{{.Package}}_{{.Class}}_inflight",
		Help:        "{{.Package}} {{.Class}} inflight",
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "custom"})
}
{{- end}}
`

func (g *WrapperGenerator) generateWrapperFunctionDeclare(info *RenderInfo, function *Function) string {
	var buf bytes.Buffer

	buf.WriteString("func ")
	if function.Recv != nil {
		if g.starClassSet[info.Class] {
			buf.WriteString(fmt.Sprintf("(w *%s)", g.wrapClassMap[info.Class]))
		} else {
			buf.WriteString(fmt.Sprintf("(w %s)", g.wrapClassMap[info.Class]))
		}
	}

	buf.WriteString(" ")
	buf.WriteString(function.Name)

	buf.WriteString("(")

	var params []string
	if len(function.Params) == 0 || function.Params[0].Type != "context.Context" {
		if !function.IsReturnError {
			if (info.EnableRuleForNoErrorFunc || info.Rule.ErrorInResult) && (info.Rule.Trace || info.Rule.Retry) {
				params = append(params, "ctx context.Context")
			}
		} else {
			if info.Rule.Trace || info.Rule.Retry {
				params = append(params, "ctx context.Context")
			}
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

func (g *WrapperGenerator) generateWrapperReturnList(function *Function) string {
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
				results = append(results, fmt.Sprintf(`&%s{obj: %s, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}`, wrapCls, i.Name))
			} else {
				results = append(results, fmt.Sprintf(`%s{obj: %s, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiterGroup: w.rateLimiterGroup}`, wrapCls, i.Name))
			}
			continue
		}

		results = append(results, fmt.Sprintf("%s", i.Name))
	}

	return strings.Join(results, ", ")
}

const WrapperFunctionBodyWithoutErrorTpl = `
{{- if .EnableRuleForNoErrorFunc}}
{{- if or .Rule.Trace .Rule.Metric}}
	ctxOptions := {{.WrapPackagePrefix}}FromContext(ctx)
{{- end}}
{{- if .Rule.RateLimiter}}
	if w.rateLimiterGroup != nil {
		_ = w.rateLimiterGroup.Wait(ctx, "{{.Class}}.{{.Function.Name}}")
	}
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
{{- if .Rule.Metric}}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", ctxOptions.MetricCustomLabelValue).Dec()
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
	return {{.Function.ReturnList}}
{{- end}}
`

const WrapperFunctionBodyWithErrorWithoutRetryTpl = `
{{- if or .Rule.Trace .Rule.Metric}}
	ctxOptions := {{.WrapPackagePrefix}}FromContext(ctx)
{{- end}}
	{{.Function.DeclareVariables}}
{{- if .Rule.RateLimiter}}
	if w.rateLimiterGroup != nil {
		if {{.Function.LastResult}} = w.rateLimiterGroup.Wait(ctx, "{{.Class}}.{{.Function.Name}}"); {{.Function.LastResult}} != nil {
			return {{.Function.ReturnList}}
		}
	}
{{- end}}
{{- if .Rule.Trace}}
	var span opentracing.Span
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ = opentracing.StartSpanFromContext(ctx, "{{.Package}}.{{.Class}}.{{.Function.Name}}")
		for key, val := range w.options.Trace.ConstTags {
			span.SetTag(key, val)
		}
		for key, val := range ctxOptions.TraceTags {
			span.SetTag(key, val)
		}
		defer span.Finish()
	}
{{- end}}
{{- if .Rule.Metric}}
	if w.options.EnableMetric && !ctxOptions.DisableMetric {
		ts := time.Now()
		w.inflightMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", ctxOptions.MetricCustomLabelValue).Inc()
		defer func() {
			w.inflightMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", ctxOptions.MetricCustomLabelValue).Dec()
			w.durationMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", {{.Function.ErrCode}}, ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
		}()
	}
{{- end}}
	{{.Function.ResultList}} = w.obj.{{.Function.Name}}({{.Function.ParamList}})
{{- if .Rule.Trace}}
	if {{.Function.LastResult}}{{- if .Rule.ErrorInResult}}.{{.ErrorField}}{{- end}} != nil && span != nil {
		span.SetTag("error", {{.Function.LastResult}}{{- if .Rule.ErrorInResult}}.{{.ErrorField}}{{- end}}.Error())
	}
{{- end}}
	return {{.Function.ReturnList}}
`

const WrapperFunctionBodyWithErrorWithRetryTpl = `
{{- if or .Rule.Trace .Rule.Metric}}
	ctxOptions := {{.WrapPackagePrefix}}FromContext(ctx)
{{- end}}
	{{.Function.DeclareVariables}}{{.NewLine}}
	{{- if or .Function.IsReturnError}}err{{else}}_{{- end}} = w.retry.Do(func() error {
{{- if .Rule.RateLimiter}}
		if w.rateLimiterGroup != nil {
			if err := w.rateLimiterGroup.Wait(ctx, "{{.Class}}.{{.Function.Name}}"); err != nil {
				return err
			}
		}
{{- end}}
{{- if .Rule.Trace}}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "{{.Package}}.{{.Class}}.{{.Function.Name}}")
			for key, val := range w.options.Trace.ConstTags {
				span.SetTag(key, val)
			}
			for key, val := range ctxOptions.TraceTags {
				span.SetTag(key, val)
			}
			defer span.Finish()
		}
{{- end}}
{{- if .Rule.Metric}}
		if w.options.EnableMetric && !ctxOptions.DisableMetric {
			ts := time.Now()
			w.inflightMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", ctxOptions.MetricCustomLabelValue).Inc()
			defer func() {
				w.inflightMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", ctxOptions.MetricCustomLabelValue).Dec()
				w.durationMetric.WithLabelValues("{{.Package}}.{{.Class}}.{{.Function.Name}}", {{.Function.ErrCode}}, ctxOptions.MetricCustomLabelValue).Observe(float64(time.Now().Sub(ts).Milliseconds()))
			}()
		}
{{- end}}
{{- if .Rule.Hystrix}}
		if w.options.EnableHystrix {
			{{- if or .Function.IsReturnError}}err{{else}}_{{- end}} = hystrix.DoC(ctx, "{{.Class}}.{{.Function.Name}}", func(ctx context.Context) error {
				{{.Function.ResultList}} = w.obj.{{.Function.Name}}({{.Function.ParamList}})
				return {{.Function.LastResult}}{{- if .Rule.ErrorInResult}}.{{.ErrorField}}{{- end}}
			}, func(ctx context.Context, err error) error {
				if w.backupObj == nil {
					return {{.Function.LastResult}}{{- if .Rule.ErrorInResult}}.{{.ErrorField}}{{- end}}
				}
				{{.Function.ResultList}} = w.backupObj.{{.Function.Name}}({{.Function.ParamList}})
				return {{.Function.LastResult}}{{- if .Rule.ErrorInResult}}.{{.ErrorField}}{{- end}}
			})
		} else {
			{{.Function.ResultList}} = w.obj.{{.Function.Name}}({{.Function.ParamList}})
		}
{{- else}}
		{{.Function.ResultList}} = w.obj.{{.Function.Name}}({{.Function.ParamList}})
{{- end}}
{{- if .Rule.Trace}}
	if {{.Function.LastResult}}{{- if .Rule.ErrorInResult}}.{{.ErrorField}}{{- end}} != nil && span != nil {
		span.SetTag("error", {{.Function.LastResult}}{{- if .Rule.ErrorInResult}}.{{.ErrorField}}{{- end}}.Error())
	}
{{- end}}
		return {{.Function.LastResult}}{{- if .Rule.ErrorInResult}}.{{.ErrorField}}{{- end}}
	})
	return {{.Function.ReturnList}}
`

func (g *WrapperGenerator) generateWrapperFunctionBody(info *RenderInfo, function *Function) string {
	var buf bytes.Buffer
	if !function.IsReturnError {
		if info.Rule.ErrorInResult && info.Rule.Retry && len(function.Results) == 1 {
			buf.WriteString(renderTemplate(WrapperFunctionBodyWithErrorWithRetryTpl, info, "WrapperFunctionBodyWithErrorWithRetryTpl"))
		} else {
			buf.WriteString(renderTemplate(WrapperFunctionBodyWithoutErrorTpl, info, "WrapperFunctionBodyWithoutErrorTpl"))
		}
		return buf.String()
	}
	if !info.Rule.Retry {
		buf.WriteString(renderTemplate(WrapperFunctionBodyWithErrorWithoutRetryTpl, info, "WrapperFunctionBodyWithoutErrorTpl"))
		return buf.String()
	}
	buf.WriteString(renderTemplate(WrapperFunctionBodyWithErrorWithRetryTpl, info, "WrapperFunctionBodyWithErrorWithRetryTpl"))
	return buf.String()
}

func (g *WrapperGenerator) MatchFunctionRule(fun string, cls string, rules map[string]Rule) bool {
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
