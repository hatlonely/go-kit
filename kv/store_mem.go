package kv

import (
	"time"

	"github.com/coocood/freecache"

	"github.com/hatlonely/go-kit/config"
)

type MemStore struct {
	cache      *freecache.Cache
	Expiration time.Duration
}

type MemStoreOptions struct {
	MemoryByte int
	Expiration time.Duration
}

var defaultMemStoreOptions = MemStoreOptions{
	MemoryByte: 200 * 1024 * 1024,
	Expiration: 15 * time.Minute,
}

type MemStoreOption func(options *MemStoreOptions)

func WithMemStoreMemoryByte(memoryByte int) MemStoreOption {
	return func(options *MemStoreOptions) {
		options.MemoryByte = memoryByte
	}
}

func WithMemStoreExpiration(expiration time.Duration) MemStoreOption {
	return func(options *MemStoreOptions) {
		options.Expiration = expiration
	}
}

func NewMemStore(opts ...MemStoreOption) *MemStore {
	options := defaultMemStoreOptions
	for _, opt := range opts {
		opt(&options)
	}
	return NewMemStoreWithOptions(&options)
}

func NewMemStoreWithOptions(options *MemStoreOptions) *MemStore {
	cache := freecache.NewCache(options.MemoryByte)
	return &MemStore{
		cache:      cache,
		Expiration: options.Expiration,
	}
}

func NewMemStoreWithConfig(conf *config.Config) (*MemStore, error) {
	var options = defaultMemStoreOptions
	if err := conf.Unmarshal(&options); err != nil {
		return nil, err
	}
	return NewMemStoreWithOptions(&options), nil
}

func (s *MemStore) Get(key []byte) ([]byte, error) {
	val, err := s.cache.Get(key)
	if err == freecache.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (s *MemStore) Set(key []byte, val []byte) (err error) {
	return s.cache.Set(key, val, int(s.Expiration/time.Second))
}

func (s *MemStore) Del(key []byte) error {
	s.cache.Del(key)
	return nil
}
