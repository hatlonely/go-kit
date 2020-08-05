package kv

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type LevelDBStore struct {
	db       *leveldb.DB
	roptions *opt.ReadOptions
	woptions *opt.WriteOptions
}

type LevelDBStoreOptions struct {
	Directory     string
	DontFillCache bool
	Strict        int
	NoWriteMerge  bool
	Sync          bool
}

var defaultLevelDBStoreOptions = LevelDBStoreOptions{
	Directory:     "leveldb/",
	DontFillCache: false,
	Strict:        0,
	NoWriteMerge:  false,
	Sync:          false,
}

type LevelDBStoreOption func(options *LevelDBStoreOptions)

func NewLevelDBStore(opts ...LevelDBStoreOption) (*LevelDBStore, error) {
	options := defaultLevelDBStoreOptions
	for _, opt := range opts {
		opt(&options)
	}

	db, err := leveldb.OpenFile(options.Directory, nil)
	if err != nil {
		return nil, err
	}

	return &LevelDBStore{
		db: db,
		roptions: &opt.ReadOptions{
			DontFillCache: options.DontFillCache,
			Strict:        opt.Strict(options.Strict),
		},
		woptions: &opt.WriteOptions{
			NoWriteMerge: options.NoWriteMerge,
			Sync:         options.Sync,
		},
	}, nil
}

func (l *LevelDBStore) Get(key []byte) ([]byte, error) {
	val, err := l.db.Get(key, l.roptions)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (l *LevelDBStore) Set(key []byte, val []byte) error {
	return l.db.Put(key, val, l.woptions)
}

func (l *LevelDBStore) Del(key []byte) error {
	return l.db.Delete(key, l.woptions)
}
