package validator

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/PaesslerAG/gval"
	"github.com/generikvault/gvalstrings"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/strex"
)

func Validate(v interface{}) error {
	rule, err := Compile(v)
	if err != nil {
		return err
	}
	return rule.Validate(v)
}

var lang = gval.Full(
	gvalstrings.SingleQuoted(),
	gval.InfixOperator("match", func(x, pattern interface{}) (interface{}, error) {
		re, err := regexp.Compile(pattern.(string))
		if err != nil {
			return nil, err
		}
		return re.MatchString(x.(string)), nil
	}),
	gval.Function("isEmail", func(x interface{}) (bool, error) {
		return strex.ReEmail.MatchString(x.(string)), nil
	}),
	gval.Function("isPhone", func(x interface{}) (bool, error) {
		return strex.RePhone.MatchString(x.(string)), nil
	}),
	gval.Function("isIdentifier", func(x interface{}) (bool, error) {
		return strex.ReIdentifier.MatchString(x.(string)), nil
	}),
)

func RegisterFunction(name string, fun interface{}) {
	lang = gval.NewLanguage(lang, gval.Function(name, fun))
}

func MustCompile(v interface{}) *Validator {
	rule, err := Compile(v)
	if err != nil {
		panic(err)
	}
	return rule
}

func Compile(v interface{}) (*Validator, error) {
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	rules := map[string]gval.Evaluable{}
	tags := map[string]string{}
	if err := interfaceToRuleRecursive(rules, tags, rt, ""); err != nil {
		return nil, errors.Wrap(err, "compile failed")
	}

	return &Validator{rules: rules, tags: tags}, nil
}

type Validator struct {
	rules map[string]gval.Evaluable
	tags  map[string]string
}

func (r *Validator) Validate(v interface{}) error {
	return evaluateInterfaceRecursive(r.rules, r.tags, v, "")
}

func evaluateInterfaceRecursive(rules map[string]gval.Evaluable, tags map[string]string, v interface{}, prefix string) error {
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	if rt.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rt = rt.Elem()
		rv = rv.Elem()
	}

	if rt.Kind() == reflect.Struct {
		for i := 0; i < rt.NumField(); i++ {
			key := rt.Field(i).Name
			if rt.Field(i).Anonymous {
				if err := evaluateInterfaceRecursive(rules, tags, rv.Field(i).Interface(), prefix); err != nil {
					return err
				}
			} else {
				if err := evaluateInterfaceRecursive(rules, tags, rv.Field(i).Interface(), prefixAppendKey(prefix, key)); err != nil {
					return err
				}
			}
		}
		return nil
	}

	eval, ok := rules[prefix]
	if !ok {
		return nil
	}
	b, err := eval.EvalBool(context.Background(), map[string]interface{}{"x": rv.Interface()})
	if err != nil {
		return &Error{Code: ErrEvalFailed, Err: errors.New("eval failed"), Key: prefix, Val: rv.Interface(), Tag: tags[prefix]}
	}
	if !b {
		return &Error{Code: ErrRuleNotMatch, Err: errors.New("rule not match"), Key: prefix, Val: rv.Interface(), Tag: tags[prefix]}
	}
	return nil
}

func interfaceToRuleRecursive(rules map[string]gval.Evaluable, tags map[string]string, rt reflect.Type, prefix string) error {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	for i := 0; i < rt.NumField(); i++ {
		t := rt.Field(i).Type
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		tag := rt.Field(i).Tag.Get("rule")
		if tag == "-" {
			continue
		}

		key := rt.Field(i).Name
		if t.Kind() == reflect.Struct {
			if rt.Field(i).Anonymous {
				if err := interfaceToRuleRecursive(rules, tags, t, prefix); err != nil {
					return err
				}
			} else {
				if err := interfaceToRuleRecursive(rules, tags, t, prefixAppendKey(prefix, key)); err != nil {
					return err
				}
			}
			continue
		}

		if tag == "" {
			continue
		}
		eval, err := lang.NewEvaluable(tag)
		if err != nil {
			return errors.Wrapf(err, "create evaluable failed. key [%v], tag [%v]", key, tag)
		}
		rules[prefixAppendKey(prefix, key)] = eval
		tags[prefixAppendKey(prefix, key)] = tag
	}

	return nil
}

func prefixAppendKey(prefix string, key string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%v.%v", prefix, key)
}
