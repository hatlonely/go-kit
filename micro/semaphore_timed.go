package micro

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/emirpasic/gods/sets/treeset"
)

type TimedSemaphore struct {
	mutex    sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond

	tokenSet *treeset.Set
	maxToken int
	timeout  time.Duration
}

func NewTimedSemaphore(maxToken int, timeout time.Duration) (*TimedSemaphore, error) {
	if maxToken <= 0 {
		return nil, errors.New("maxToken should be positive")
	}
	if timeout < 0 {
		return nil, errors.New("timeout should greater than or equal to 0")
	}

	c := &TimedSemaphore{
		tokenSet: treeset.NewWithIntComparator(),
		maxToken: maxToken,
		timeout:  timeout,
	}
	c.notFull = sync.NewCond(&c.mutex)
	c.notEmpty = sync.NewCond(&c.mutex)

	if timeout != 0 {
		go func() {
			for {
				c.mutex.Lock()
				for c.tokenSet.Size() == 0 {
					c.notEmpty.Wait()
				}
				it := c.tokenSet.Iterator()
				it.Next()
				val := it.Value().(int)
				sleepTimeNs := int64(val) - time.Now().Add(-c.timeout).UnixNano()
				if sleepTimeNs > 0 {
					c.mutex.Unlock()
					time.Sleep(time.Duration(sleepTimeNs) * time.Nanosecond)
				} else {
					c.tokenSet.Remove(val)
					c.mutex.Unlock()
					c.notFull.Signal()
				}
			}
		}()
	}

	return c, nil
}

var tokenSequence = uint64(0)

func (c *TimedSemaphore) TryAcquire() (int, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.tokenSet.Size() >= c.maxToken {
		return 0, false
	}
	atomic.AddUint64(&tokenSequence, 1)
	token := int(time.Now().UnixNano()/1000*1000) + int(tokenSequence%1000)
	c.tokenSet.Add(token)
	c.notEmpty.Signal()
	return token, true
}

func (c *TimedSemaphore) Acquire() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	defer c.notEmpty.Signal()

	for c.tokenSet.Size() >= c.maxToken {
		c.notFull.Wait()
	}
	atomic.AddUint64(&tokenSequence, 1)
	// 这里不直接使用 time.Now().UnixNano 是因为 treemap.Set 不支持可重复元素，但是 time.Now() 是可能重复的
	// 这里直接在后面拼上一个序列，如果在千分之一毫秒内颁发超过 1000 个 token，可能会出现重复
	token := int(time.Now().UnixNano()/1000*1000) + int(tokenSequence%1000)
	c.tokenSet.Add(token)
	return token
}

func (c *TimedSemaphore) Release(token int) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.tokenSet.Size() == 0 {
		return false
	}
	if !c.tokenSet.Contains(token) {
		return false
	}
	c.tokenSet.Remove(token)
	c.notFull.Signal()
	return true
}
