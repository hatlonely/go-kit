package micro

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type LocalTimedSemaphoreParallelControllerOptions map[string]struct {
	MaxToken int
	Timeout  time.Duration
}

type LocalTimedSemaphoreParallelController struct {
	options       *LocalTimedSemaphoreParallelControllerOptions
	controllerMap map[string]*TimedSemaphore
}

func NewLocalTimedSemaphoreParallelControllerWithOptions(options *LocalTimedSemaphoreParallelControllerOptions) (*LocalTimedSemaphoreParallelController, error) {
	if options == nil || len(*options) == 0 {
		return nil, nil
	}

	controllerGroup := map[string]*TimedSemaphore{}
	for key, val := range *options {
		if val.MaxToken <= 0 {
			return nil, errors.New("max parallel should be positive")
		}
		var err error
		controllerGroup[key], err = NewTimedSemaphore(val.MaxToken, val.Timeout)
		if err != nil {
			return nil, errors.Wrap(err, "NewSemaphore failed")
		}
	}

	c := &LocalTimedSemaphoreParallelController{
		options:       options,
		controllerMap: controllerGroup,
	}

	return c, nil
}

func (l *LocalTimedSemaphoreParallelController) TryAcquire(ctx context.Context, key string) (int, error) {
	c, ok := l.controllerMap[key]
	if !ok {
		return 0, nil
	}

	if token, ok := c.TryAcquire(); ok {
		return token, nil
	}
	return 0, ErrParallelControl
}

func (l *LocalTimedSemaphoreParallelController) Acquire(ctx context.Context, key string) (int, error) {
	c, ok := l.controllerMap[key]
	if !ok {
		return 0, nil
	}

	ch := make(chan int, 1)
	go func() {
		ch <- c.Acquire()
	}()

	select {
	case <-ctx.Done():
		return 0, ErrContextCancel
	case token := <-ch:
		return token, nil
	}
}

func (l *LocalTimedSemaphoreParallelController) Release(ctx context.Context, key string, token int) error {
	c, ok := l.controllerMap[key]
	if !ok {
		return nil
	}

	ch := make(chan struct{}, 1)
	go func() {
		c.Release(token)
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ErrContextCancel
	case <-ch:
		return nil
	}
}
