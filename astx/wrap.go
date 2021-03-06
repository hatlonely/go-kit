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
	SourceDir                string              `flag:"usage: source directory; default: vendor"`
	SourcePath               string              `flag:"usage: source path default is $SourceDir/$PackagePath"`
	PackagePath              string              `flag:"usage: package path"`
	PackageName              string              `flag:"usage: package name"`
	OutputPackage            string              `flag:"usage: output package name; default: wrap" dft:"wrap"`
	UseWrapPackagePrefix     bool                `flag:"usage: use wrap package prefix"`
	Classes                  []string            `flag:"usage: classes to wrap"`
	StarClasses              []string            `flag:"usage: star classes to wrap"`
	ClassPrefix              string              `flag:"usage: wrap class name prefix"`
	UnwrapFunc               string              `flag:"usage: unwrap function name; default: Unwrap" dft:"Unwrap"`
	ErrorField               string              `flag:"usage: function return no error, error is a filed in result"`
	Inherit                  map[string][]string `flag:"usage: inherit map"`
	EnableRuleForNoErrorFunc bool                `flag:"usage: enable trace for no error function"`
	EnableHystrix            bool                `flag:"usage: enable hystrix code"`

	Rule struct {
		MainClass                  Rule            // default rule for Constructor/OnWrapperChange/OnRetryChange/OnRateLimiterChange/OnParallelControllerChange
		Wrap                       map[string]Rule // default rule for Trace/Retry/Metric/RateLimiter/ParallelController
		Class                      Rule
		Interface                  Rule
		StarClass                  Rule
		Constructor                Rule
		OnWrapperChange            Rule
		OnRetryChange              Rule
		OnRateLimiterChange        Rule
		OnParallelControllerChange Rule
		CreateMetric               Rule
		ErrorInResult              Rule
		Function                   map[string]Rule
		Trace                      map[string]Rule
		Retry                      map[string]Rule
		Metric                     map[string]Rule
		RateLimiter                map[string]Rule
		ParallelController         map[string]Rule
		Hystrix                    map[string]Rule
	}
}

func fillDefaultRule(dst *Rule, dft Rule) {
	if dst.Include == nil {
		dst.Include = dft.Include
	}
	if dst.Exclude == nil {
		dst.Exclude = dft.Exclude
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
	if options.Rule.MainClass.Exclude == nil {
		options.Rule.MainClass.Exclude = excludeAllRegex
	}
	fillDefaultRule(&options.Rule.Constructor, options.Rule.MainClass)
	fillDefaultRule(&options.Rule.OnWrapperChange, options.Rule.MainClass)
	fillDefaultRule(&options.Rule.OnRetryChange, options.Rule.MainClass)
	fillDefaultRule(&options.Rule.OnRateLimiterChange, options.Rule.MainClass)
	fillDefaultRule(&options.Rule.OnParallelControllerChange, options.Rule.MainClass)
	fillDefaultRule(&options.Rule.CreateMetric, options.Rule.MainClass)

	if options.Rule.StarClass.Exclude == nil {
		options.Rule.StarClass.Exclude = excludeAllRegex
	}
	if options.Rule.Class.Exclude == nil {
		options.Rule.Class.Exclude = excludeAllRegex
	}
	if options.Rule.ErrorInResult.Exclude == nil {
		options.Rule.ErrorInResult.Exclude = excludeAllRegex
	}
	if options.Rule.Interface.Exclude == nil {
		options.Rule.Interface.Exclude = excludeAllRegex
	}
	if len(options.Rule.Hystrix) == 0 {
		options.Rule.Hystrix = map[string]Rule{
			"default": {
				Exclude: regexp.MustCompile(`^.*$`),
			},
		}
	}

	if len(options.Rule.Retry) == 0 {
		options.Rule.Retry = options.Rule.Wrap
	}
	if len(options.Rule.Trace) == 0 {
		options.Rule.Trace = options.Rule.Wrap
	}
	if len(options.Rule.Metric) == 0 {
		options.Rule.Metric = options.Rule.Wrap
	}
	if len(options.Rule.RateLimiter) == 0 {
		options.Rule.RateLimiter = options.Rule.Wrap
	}
	if len(options.Rule.ParallelController) == 0 {
		options.Rule.ParallelController = options.Rule.Wrap
	}

	var wrapPackagePrefix string
	if options.OutputPackage != "wrap" || options.UseWrapPackagePrefix {
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
	Debug                bool
	Package              string
	WrapPackagePrefix    string
	UnwrapFunc           string
	Class                string
	WrapClass            string
	OutputPackage        string
	UseWrapPackagePrefix bool
	PackagePath          string
	ErrorField           string
	Interface            string
	ObjectType           string
	EnableHystrix        bool
	IsStarClass          bool
	NewLine              string
	Function             struct {
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
		Constructor                bool
		OnWrapperChange            bool
		OnRetryChange              bool
		OnRateLimiterChange        bool
		OnParallelControllerChange bool
		CreateMetric               bool
		Interface                  bool
		Trace                      bool
		Retry                      bool
		Metric                     bool
		Hystrix                    bool
		RateLimiter                bool
		ParallelController         bool
		ErrorInResult              bool
	}
}

func (g *WrapperGenerator) Generate() (string, error) {
	sourcePath := g.options.SourcePath
	if sourcePath == "" {
		sourcePath = path.Join(g.options.SourceDir, g.options.PackagePath)
	}
	functions, err := ParseFunction(sourcePath, g.options.PackageName)
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
		UseWrapPackagePrefix:     g.options.UseWrapPackagePrefix,
	}

	buf.WriteString(renderTemplate(WrapperImportTpl, info, "WrapperImportTpl"))

	for _, cls := range classes {
		info.Class = cls
		info.WrapClass = g.wrapClassMap[cls]
		info.Rule.Constructor = g.MatchRule(cls, g.options.Rule.Constructor)
		info.Rule.OnWrapperChange = g.MatchRule(cls, g.options.Rule.OnWrapperChange)
		info.Rule.OnRetryChange = g.MatchRule(cls, g.options.Rule.OnRetryChange)
		info.Rule.OnRateLimiterChange = g.MatchRule(cls, g.options.Rule.OnRateLimiterChange)
		info.Rule.OnParallelControllerChange = g.MatchRule(cls, g.options.Rule.OnParallelControllerChange)
		info.Rule.CreateMetric = g.MatchRule(cls, g.options.Rule.CreateMetric)
		info.Rule.Interface = g.MatchRule(cls, g.options.Rule.Interface)
		info.IsStarClass = g.starClassSet[cls]
		info.Interface = fmt.Sprintf("I%s%s", g.options.ClassPrefix, cls)
		if info.Rule.Interface {
			info.ObjectType = fmt.Sprintf("%s.%s", info.Package, info.Class)
		} else {
			info.ObjectType = fmt.Sprintf("*%s.%s", info.Package, info.Class)
		}

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
			info.Rule.OnRateLimiterChange = g.MatchRule(cls, g.options.Rule.OnRateLimiterChange)
			info.Rule.OnParallelControllerChange = g.MatchRule(cls, g.options.Rule.OnParallelControllerChange)
			info.Rule.Trace = g.MatchFunctionRule(function.Name, cls, g.options.Rule.Trace)
			info.Rule.Metric = g.MatchFunctionRule(function.Name, cls, g.options.Rule.Metric)
			info.Rule.Retry = g.MatchFunctionRule(function.Name, cls, g.options.Rule.Retry)
			info.Rule.RateLimiter = g.MatchFunctionRule(function.Name, cls, g.options.Rule.RateLimiter)
			info.Rule.ParallelController = g.MatchFunctionRule(function.Name, cls, g.options.Rule.ParallelController)
			info.Rule.Hystrix = g.MatchFunctionRule(function.Name, cls, g.options.Rule.Hystrix)
			if len(function.Results) == 1 {
				info.Rule.ErrorInResult = g.MatchRule(function.Results[0].Type, g.options.Rule.ErrorInResult)
			} else {
				info.Rule.ErrorInResult = false
			}

			if !function.IsReturnError && !(info.Rule.ErrorInResult && len(function.Results) == 1) {
				info.Function.ErrCode = `"OK"`
			} else if info.OutputPackage == "wrap" && !info.UseWrapPackagePrefix {
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
	"github.com/hatlonely/go-kit/micro"
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
{{- if .Rule.Constructor}}

func New{{.WrapClass}}(
	obj {{.ObjectType}},
	retry *micro.Retry,
	options *{{.WrapPackagePrefix}}WrapperOptions,
	durationMetric *prometheus.HistogramVec,
	inflightMetric *prometheus.GaugeVec,
	rateLimiter micro.RateLimiter,
	parallelController micro.ParallelController) *{{.WrapClass}} {
	return &{{.WrapClass}}{
		obj:                obj,
		retry:              retry,
		options:            options,
		durationMetric:     durationMetric,
		inflightMetric:     inflightMetric,
		rateLimiter:        rateLimiter,
		parallelController: parallelController,
	}
}
{{- end}}

type {{.WrapClass}} struct {
	obj                {{.ObjectType}}
{{- if .EnableHystrix}}
	backupObj          {{.Interface}}
{{- end}}
	retry              *micro.Retry
	options            *{{.WrapPackagePrefix}}WrapperOptions
	durationMetric     *prometheus.HistogramVec
	inflightMetric     *prometheus.GaugeVec
	rateLimiter        micro.RateLimiter
	parallelController micro.ParallelController
}

{{- if .IsStarClass}}
func (w *{{.WrapClass}}) {{.UnwrapFunc}}() {{.ObjectType}} {
	return w.obj
}
{{- else}}
func (w {{.WrapClass}}) {{.UnwrapFunc}}() {{.ObjectType}} {
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
		var options micro.RetryOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		retry, err := micro.NewRetryWithOptions(&options)
		if err != nil {
			return errors.Wrap(err, "NewRetryWithOptions failed")
		}
		w.retry = retry
		return nil
	}
}
{{- end}}

{{- if .Rule.OnRateLimiterChange}}

func (w *{{.WrapClass}}) OnRateLimiterChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options micro.RateLimiterOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		rateLimiter, err := micro.NewRateLimiterWithOptions(&options, opts...)
		if err != nil {
			return errors.Wrap(err, "NewRateLimiterWithOptions failed")
		}
		w.rateLimiter = rateLimiter
		return nil
	}
}
{{- end}}

{{- if .Rule.OnParallelControllerChange}}

func (w *{{.WrapClass}}) OnParallelControllerChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options micro.ParallelControllerOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		parallelController, err := micro.NewParallelControllerWithOptions(&options, opts...)
		if err != nil {
			return errors.Wrap(err, "NewParallelControllerWithOptions failed")
		}
		w.parallelController = parallelController
		return nil
	}
}
{{- end}}

{{- if .Rule.CreateMetric}}

func (w *{{.WrapClass}}) CreateMetric(options *{{.WrapPackagePrefix}}WrapperOptions) {
	w.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        fmt.Sprintf("%s_{{.Package}}_{{.Class}}_durationMs", options.Name),
		Help:        "{{.Package}} {{.Class}} response time milliseconds",
		Buckets:     options.Metric.Buckets,
		ConstLabels: options.Metric.ConstLabels,
	}, []string{"method", "errCode", "custom"})
	w.inflightMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        fmt.Sprintf("%s_{{.Package}}_{{.Class}}_inflight", options.Name),
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
				results = append(results, fmt.Sprintf(`&%s{obj: %s, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}`, wrapCls, i.Name))
			} else {
				results = append(results, fmt.Sprintf(`%s{obj: %s, retry: w.retry, options: w.options, durationMetric: w.durationMetric, inflightMetric: w.inflightMetric, rateLimiter: w.rateLimiter, parallelController: w.parallelController}`, wrapCls, i.Name))
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
	if w.rateLimiter != nil {
		_ = w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.{{.Class}}.{{.Function.Name}}", w.options.Name))
	}
{{- end}}
{{- if .Rule.ParallelController}}
	if w.parallelController != nil {
		token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.{{.Class}}.{{.Function.Name}}", w.options.Name))
		if err == nil {
			defer w.parallelController.Release(ctx, fmt.Sprintf("%s.{{.Class}}.{{.Function.Name}}", w.options.Name), token)
		}
	}
{{- end}}
{{- if .Rule.Trace}}
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "{{.Package}}.{{.Class}}.{{.Function.Name}}", ctxOptions.StartSpanOpts...)
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
	if w.rateLimiter != nil {
		if {{.Function.LastResult}} = w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.{{.Class}}.{{.Function.Name}}", w.options.Name)); {{.Function.LastResult}} != nil {
			return {{.Function.ReturnList}}
		}
	}
{{- end}}
{{- if .Rule.ParallelController}}
	if w.parallelController != nil {
		if token, {{.Function.LastResult}} = w.parallelController.Acquire(ctx, fmt.Sprintf("%s.{{.Class}}.{{.Function.Name}}", w.options.Name)); {{.Function.LastResult}} != nil {
			return {{.Function.ReturnList}}
		} else {
			defer w.parallelController.Release(ctx, fmt.Sprintf("%s.{{.Class}}.{{.Function.Name}}", w.options.Name), token)
		}
	}
{{- end}}
{{- if .Rule.Trace}}
	var span opentracing.Span
	if w.options.EnableTrace && !ctxOptions.DisableTrace {
		span, _ = opentracing.StartSpanFromContext(ctx, "{{.Package}}.{{.Class}}.{{.Function.Name}}", ctxOptions.StartSpanOpts...)
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
		if w.rateLimiter != nil {
			if err := w.rateLimiter.Wait(ctx, fmt.Sprintf("%s.{{.Class}}.{{.Function.Name}}", w.options.Name)); err != nil {
				return err
			}
		}
{{- end}}
{{- if .Rule.ParallelController}}
		if w.parallelController != nil {
			if token, err := w.parallelController.Acquire(ctx, fmt.Sprintf("%s.{{.Class}}.{{.Function.Name}}", w.options.Name)); err != nil {
				return err
			} else {
				defer w.parallelController.Release(ctx, fmt.Sprintf("%s.{{.Class}}.{{.Function.Name}}", w.options.Name), token)
			}
		}
{{- end}}
{{- if .Rule.Trace}}
		var span opentracing.Span
		if w.options.EnableTrace && !ctxOptions.DisableTrace {
			span, _ = opentracing.StartSpanFromContext(ctx, "{{.Package}}.{{.Class}}.{{.Function.Name}}", ctxOptions.StartSpanOpts...)
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
