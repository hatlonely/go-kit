package micro

import (
	"context"

	"github.com/pkg/errors"
)

type LocalParallelControllerOptions map[string]int

type LocalParallelController struct {
	options       *LocalParallelControllerOptions
	controllerMap map[string]*LocalParallelControllerCell
}

func NewLocalParallelControllerGroupWithOptions(options *LocalParallelControllerOptions) (*LocalParallelController, error) {
	if options == nil || len(*options) == 0 {
		return nil, nil
	}

	controllerGroup := map[string]*LocalParallelControllerCell{}
	for key, val := range *options {
		if val <= 0 {
			return nil, errors.New("max parallel should be positive")
		}
		controllerGroup[key] = NewLocalParallelControllerCell(val)
	}

	c := &LocalParallelController{
		options:       options,
		controllerMap: controllerGroup,
	}

	return c, nil
}

func (l *LocalParallelController) TryGetToken(ctx context.Context, key string) bool {
	c, ok := l.controllerMap[key]
	if !ok {
		return true
	}
	return c.TryGetToken(ctx)
}

func (l *LocalParallelController) TryPutToken(ctx context.Context, key string) bool {
	c, ok := l.controllerMap[key]
	if !ok {
		return true
	}
	return c.TryPutToken(ctx)
}

func (l *LocalParallelController) PutToken(ctx context.Context, key string) error {
	c, ok := l.controllerMap[key]
	if !ok {
		return nil
	}
	return c.PutToken(ctx)
}

func (l *LocalParallelController) GetToken(ctx context.Context, key string) error {
	c, ok := l.controllerMap[key]
	if !ok {
		return nil
	}
	return c.GetToken(ctx)
}
