package micro

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

type LocalLocker struct {
	mutex    sync.RWMutex
	mutexMap map[string]*Semaphore
}

func NewLocalLocker() *LocalLocker {
	return &LocalLocker{
		mutexMap: map[string]*Semaphore{},
	}
}

func (l *LocalLocker) Lock(ctx context.Context, key string) error {
	l.mutex.RLock()
	m, ok := l.mutexMap[key]
	l.mutex.RUnlock()
	if !ok {
		l.mutex.Lock()
		m, ok = l.mutexMap[key]
		if !ok {
			m, _ = NewSemaphore(1)
			l.mutexMap[key] = m
		}
		l.mutex.Unlock()
	}
	m.Acquire()
	return nil
}

func (l *LocalLocker) TryLock(ctx context.Context, key string) error {
	l.mutex.RLock()
	m, ok := l.mutexMap[key]
	l.mutex.RUnlock()
	if !ok {
		l.mutex.Lock()
		m, ok = l.mutexMap[key]
		if !ok {
			m, _ = NewSemaphore(1)
			l.mutexMap[key] = m
		}
		l.mutex.Unlock()
	}
	if !m.TryAcquire() {
		return ErrLocked
	}
	return nil
}

func (l *LocalLocker) Unlock(ctx context.Context, key string) error {
	l.mutex.RLock()
	m, ok := l.mutexMap[key]
	l.mutex.RUnlock()
	if !ok {
		return errors.New("release unlocked key")
	}
	m.Release()
	return nil
}
