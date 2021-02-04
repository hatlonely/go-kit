package micro

import (
	"context"

	"github.com/pkg/errors"
)

type LocalParallelControllerGroupOptions map[string]int

type LocalParallelControllerGroup struct {
	options         *LocalParallelControllerGroupOptions
	controllerGroup map[string]*LocalParallelController
}

func NewLocalParallelControllerGroupWithOptions(options *LocalParallelControllerGroupOptions) (*LocalParallelControllerGroup, error) {
	if options == nil || len(*options) == 0 {
		return nil, nil
	}

	controllerGroup := map[string]*LocalParallelController{}
	for key, val := range *options {
		if val <= 0 {
			return nil, errors.New("max parallel should be positive")
		}
		controllerGroup[key] = NewLocalParallelController(val)
	}

	c := &LocalParallelControllerGroup{
		options:         options,
		controllerGroup: controllerGroup,
	}

	return c, nil
}

func (l *LocalParallelControllerGroup) PutToken(ctx context.Context, key string) error {
	c, ok := l.controllerGroup[key]
	if !ok {
		return nil
	}
	return c.PutToken(ctx)
}

func (l *LocalParallelControllerGroup) GetToken(ctx context.Context, key string) error {
	c, ok := l.controllerGroup[key]
	if !ok {
		return nil
	}
	return c.GetToken(ctx)
}
