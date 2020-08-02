package kv

type Store interface {
	Get(key []byte) (val []byte, err error)
	Set(key []byte, val []byte) (err error)
}
