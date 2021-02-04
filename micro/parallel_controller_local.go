package micro

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

func NewLocalParallelController(maxToken int) *LocalParallelController {
	c := &LocalParallelController{
		maxToken: maxToken,
		curToken: maxToken,
	}

	c.notFull = sync.NewCond(&c.mutex)
	c.notEmpty = sync.NewCond(&c.mutex)

	return c
}

type LocalParallelController struct {
	mutex    sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond

	maxToken int
	curToken int
}

func (l *LocalParallelController) PutToken(ctx context.Context) error {
	ch := make(chan struct{}, 1)
	go func() {
		l.putToken()
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return errors.New("context cancel")
	case <-ch:
		return nil
	}
}

func (l *LocalParallelController) GetToken(ctx context.Context) error {
	ch := make(chan struct{}, 1)
	go func() {
		l.getToken()
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return errors.New("context cancel")
	case <-ch:
		return nil
	}
}

func (l *LocalParallelController) putToken() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	defer l.notEmpty.Signal()

	for l.curToken >= l.maxToken {
		l.notFull.Wait()
	}
	l.curToken++
}

func (l *LocalParallelController) getToken() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	defer l.notFull.Signal()

	for l.curToken <= 0 {
		l.notEmpty.Wait()
	}
	l.curToken--
}
