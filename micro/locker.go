package micro

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

var ErrLocked = errors.New("locked")

type Locker interface {
	TryLock(ctx context.Context, key string) error
	Lock(ctx context.Context, key string) error
	Unlock(ctx context.Context, key string) error
}

func RegisterLocker(key string, constructor interface{}) {
	if _, ok := lockerConstructorMap[key]; ok {
		panic(fmt.Sprintf("locker type [%v] is already registered", key))
	}

	info, err := refx.NewConstructor(constructor, reflect.TypeOf((*Locker)(nil)).Elem())
	refx.Must(err)

	lockerConstructorMap[key] = info
}

var lockerConstructorMap = map[string]*refx.Constructor{}

func NewLockerWithOptions(options *LockerOptions, opts ...refx.Option) (Locker, error) {
	if options.Type == "" {
		return nil, nil
	}

	constructor, ok := lockerConstructorMap[options.Type]
	if !ok {
		return nil, errors.Errorf("unregistered Locker type: [%v]", options.Type)
	}

	result, err := constructor.Call(options.Options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "constructor.Call failed")
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, errors.Wrapf(result[1].Interface().(error), "NewLockerWithOptions failed. type: [%v]", options.Type)
		}
		return result[0].Interface().(Locker), nil
	}

	return result[0].Interface().(Locker), nil
}

type LockerOptions struct {
	Type    string
	Options interface{}
}
