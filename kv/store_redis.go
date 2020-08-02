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

	keyIdx int
	keyLen int
}

type RedisStoreOptions struct {
	Expiration time.Duration
	Mode       string
	KeyIdx     int
	KeyLen     int
}

var defaultRedisStoreOptions = RedisStoreOptions{
	Mode: "string",
}

type RedisStoreOption func(options *RedisStoreOptions)

func NewRedisStore(client *redis.Client, opts ...RedisStoreOption) (*RedisStore, error) {
	options := defaultRedisStoreOptions
	for _, opt := range opts {
		opt(&options)
	}

	store := &RedisStore{
		client:     client,
		expiration: options.Expiration,
	}
	if options.Mode == "string" {
		store.getFunc = redisStringGet
		store.setFunc = redisStringSet
	} else {
		store.getFunc = redisHashGet
		store.setFunc = redisHashSet
	}

	return store, nil
}

func (s *RedisStore) Get(key []byte) ([]byte, error) {
	return s.getFunc(s, key)
}

func (s *RedisStore) Set(key []byte, val []byte) error {
	return s.setFunc(s, key, val)
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

func (s *RedisStore) parseKey(key string) (string, string) {
	if len(key) > s.keyIdx+s.keyLen {
		return key[s.keyIdx : s.keyIdx+s.keyLen], key[:s.keyIdx] + key[s.keyIdx+s.keyLen:]
	} else if len(key) > s.keyIdx {
		return key[s.keyIdx:], key[:s.keyIdx]
	} else {
		return "", key
	}
}
