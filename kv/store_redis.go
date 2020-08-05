package kv

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisStore struct {
	client     *redis.Client
	expiration time.Duration

	getFunc func(s *RedisStore, key []byte) ([]byte, error)
	setFunc func(s *RedisStore, key []byte, val []byte) error
	delFunc func(s *RedisStore, key []byte) error

	keyIdx int
	keyLen int
}

func NewRedisStore(client *redis.Client, opts ...RedisStoreOption) (*RedisStore, error) {
	options := defaultRedisStoreOptions
	for _, opt := range opts {
		opt(&options)
	}

	return NewRedisStoreWithOptions(client, &options)
}

func NewRedisStoreWithOptions(client *redis.Client, options *RedisStoreOptions) (*RedisStore, error) {
	store := &RedisStore{
		client:     client,
		expiration: options.Expiration,
	}
	if options.Mode == RedisStoreModeString {
		store.getFunc = redisStringGet
		store.setFunc = redisStringSet
		store.delFunc = redisStringDel
	} else {
		store.getFunc = redisHashGet
		store.setFunc = redisHashSet
		store.delFunc = redisHashDel
	}

	return store, nil
}

func (s *RedisStore) Get(key []byte) ([]byte, error) {
	return s.getFunc(s, key)
}

func (s *RedisStore) Set(key []byte, val []byte) error {
	return s.setFunc(s, key, val)
}

func (s *RedisStore) Del(key []byte) error {
	return s.delFunc(s, key)
}

func redisStringGet(s *RedisStore, key []byte) ([]byte, error) {
	val, err := s.client.Get(string(key)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func redisStringSet(s *RedisStore, key []byte, val []byte) error {
	return s.client.Set(string(key), val, s.expiration).Err()
}

func redisStringDel(s *RedisStore, key []byte) error {
	return s.client.Del(string(key)).Err()
}

func redisHashGet(s *RedisStore, key []byte) ([]byte, error) {
	k, f := s.parseKey(string(key))
	val, err := s.client.HGet(k, f).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func redisHashSet(s *RedisStore, key []byte, val []byte) error {
	k, f := s.parseKey(string(key))
	return s.client.HSet(k, f, val).Err()
}

func redisHashDel(s *RedisStore, key []byte) error {
	k, f := s.parseKey(string(key))
	return s.client.HDel(k, f).Err()
}

func (s *RedisStore) parseKey(key string) (string, string) {
	if len(key) > s.keyIdx+s.keyLen {
		return key[s.keyIdx : s.keyIdx+s.keyLen], key[:s.keyIdx] + key[s.keyIdx+s.keyLen:]
	} else if len(key) > s.keyIdx {
		return key[s.keyIdx:], key[:s.keyIdx]
	} else {
		return "", key
	}
}

type RedisStoreOptions struct {
	Expiration time.Duration
	Mode       RedisStoreMode
	KeyIdx     int
	KeyLen     int
}

var defaultRedisStoreOptions = RedisStoreOptions{
	Expiration: 30 * time.Minute,
	Mode:       RedisStoreModeString,
}

type RedisStoreMode int

const (
	RedisStoreModeString RedisStoreMode = 1
	RedisStoreModeHash   RedisStoreMode = 2
)

type RedisStoreOption func(options *RedisStoreOptions)

func WithRedisExpiration(expiration time.Duration) RedisStoreOption {
	return func(options *RedisStoreOptions) {
		options.Expiration = expiration
	}
}

func WithRedisStoreModeString() RedisStoreOption {
	return func(options *RedisStoreOptions) {
		options.Mode = RedisStoreModeString
	}
}

func WithRedisStoreModeHash() RedisStoreOption {
	return func(options *RedisStoreOptions) {
		options.Mode = RedisStoreModeHash
	}
}

func WithRedisHashKey(keyIdx int, keyLen int) RedisStoreOption {
	return func(options *RedisStoreOptions) {
		options.KeyIdx = keyIdx
		options.KeyLen = keyLen
	}
}
