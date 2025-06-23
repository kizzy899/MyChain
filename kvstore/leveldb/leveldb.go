package leveldb

import (
	"Trie/kvstore"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDBStore struct {
	db *leveldb.DB
}

// NewLevelDBStore 创建并打开一个 LevelDB 实例
func NewLevelDBStore(path string) (kvstore.KVStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBStore{db: db}, nil
}

func (l *LevelDBStore) Get(key []byte) ([]byte, error) {
	return l.db.Get(key, nil)
}

func (l *LevelDBStore) Put(key, value []byte) error {
	return l.db.Put(key, value, nil)
}

func (l *LevelDBStore) Delete(key []byte) error {
	return l.db.Delete(key, nil)
}

func (l *LevelDBStore) Has(key []byte) (bool, error) {
	return l.db.Has(key, nil)
}

func (l *LevelDBStore) Close() error {
	return l.db.Close()
}
