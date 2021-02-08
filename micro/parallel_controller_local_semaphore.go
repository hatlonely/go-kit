package micro

import (
	"context"

	"github.com/pkg/errors"
)

type LocalSemaphoreParallelControllerOptions map[string]int

type LocalSemaphoreParallelController struct {
	options       *LocalSemaphoreParallelControllerOptions
	controllerMap map[string]*Semaphore
}

func NewLocalSemaphoreParallelControllerWithOptions(options *LocalSemaphoreParallelControllerOptions) (*LocalSemaphoreParallelController, error) {
	if options == nil || len(*options) == 0 {
		return nil, nil
	}

	controllerGroup := map[string]*Semaphore{}
	for key, val := range *options {
		if val <= 0 {
			return nil, errors.New("max parallel should be positive")
		}
		var err error
		controllerGroup[key], err = NewSemaphore(val)
		if err != nil {
			return nil, errors.Wrap(err, "NewSemaphore failed")
		}
	}

	c := &LocalSemaphoreParallelController{
		options:       options,
		controllerMap: controllerGroup,
	}

	return c, nil
}

func (l *LocalSemaphoreParallelController) TryAcquire(ctx context.Context, key string) (int, error) {
	c, ok := l.controllerMap[key]
	if !ok {
		return 0, nil
	}
	if c.TryAcquire() {
		return 0, nil
	}
	return 0, ErrParallelControl
}

func (l *LocalSemaphoreParallelController) Acquire(ctx context.Context, key string) (int, error) {
	c, ok := l.controllerMap[key]
	if !ok {
		return 0, nil
	}

	ch := make(chan struct{}, 1)
	go func() {
		c.Acquire()
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return 0, ErrContextCancel
	case <-ch:
		return 0, nil
	}
}

func (l *LocalSemaphoreParallelController) Release(ctx context.Context, key string, token int) error {
	c, ok := l.controllerMap[key]
	if !ok {
		return nil
	}

	ch := make(chan struct{}, 1)
	go func() {
		c.ForceRelease()
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ErrContextCancel
	case <-ch:
		return nil
	}
}
