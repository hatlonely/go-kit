package micro

import (
	"sync"

	"github.com/pkg/errors"
)

func NewSemaphore(maxToken int) (*Semaphore, error) {
	if maxToken <= 0 {
		return nil, errors.New("max token should be positive")
	}

	c := &Semaphore{
		maxToken: maxToken,
	}
	c.notEmpty = sync.NewCond(&c.mutex)
	c.notFull = sync.NewCond(&c.mutex)
	return c, nil
}

type Semaphore struct {
	mutex    sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond

	maxToken int
	curToken int
}

func (l *Semaphore) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.curToken >= l.maxToken {
		return false
	}
	l.curToken++
	l.notEmpty.Signal()
	return true
}

func (l *Semaphore) Acquire() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	defer l.notEmpty.Signal()

	for l.curToken >= l.maxToken {
		l.notFull.Wait()
	}
	l.curToken++
}

func (l *Semaphore) Release() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	defer l.notFull.Signal()

	for l.curToken <= 0 {
		l.notEmpty.Wait()
	}
	l.curToken--
}

func (l *Semaphore) ForceRelease() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.curToken <= 0 {
		return
	}
	l.curToken--
	l.notFull.Signal()
}
