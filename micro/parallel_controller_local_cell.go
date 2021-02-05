package micro

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

func NewLocalParallelControllerCell(maxToken int) *LocalParallelControllerCell {
	c := &LocalParallelControllerCell{
		maxToken: maxToken,
		curToken: maxToken,
	}

	c.notFull = sync.NewCond(&c.mutex)
	c.notEmpty = sync.NewCond(&c.mutex)

	return c
}

type LocalParallelControllerCell struct {
	mutex    sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond

	maxToken int
	curToken int
}

func (l *LocalParallelControllerCell) TryGetToken() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.curToken <= 0 {
		return false
	}
	l.curToken--
	return true
}

func (l *LocalParallelControllerCell) TryPutToken() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for l.curToken >= l.maxToken {
		return false
	}
	l.curToken++
	return true
}

func (l *LocalParallelControllerCell) GetToken(ctx context.Context) error {
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

func (l *LocalParallelControllerCell) PutToken(ctx context.Context) error {
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

func (l *LocalParallelControllerCell) getToken() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	defer l.notFull.Signal()

	for l.curToken <= 0 {
		l.notEmpty.Wait()
	}
	l.curToken--
}

func (l *LocalParallelControllerCell) putToken() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	defer l.notEmpty.Signal()

	for l.curToken >= l.maxToken {
		l.notFull.Wait()
	}
	l.curToken++
}
