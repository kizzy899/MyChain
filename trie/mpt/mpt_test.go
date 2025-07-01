package trie

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSimpleInsertAndSearch(t *testing.T) {
	// 这里假设有一个内存kvstore实现，满足kvstore.KVStore接口
	db := NewInMemoryKVStore()

	trie := NewMPT(db)

	key := []byte("key1")
	value := []byte("value1")

	fmt.Printf("Insert key: %s, value: %s\n", key, value)

	err := trie.Insert(key, value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	rootHash, _ := trie.RootHash()
	fmt.Printf("Root hash after insert: %x\n", rootHash.Bytes())

	// 直接用相同key查找
	result, err := trie.Search(key)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	fmt.Printf("Search result: %s\n", string(result))

	// 简单断言
	if !bytes.Equal(result, value) {
		t.Fatalf("Search result mismatch, want %s, got %s", value, result)
	}
}

// 以下是一个非常简单的内存KVStore实现，只做排查用
type InMemoryKVStore struct {
	store map[string][]byte
}

func (m *InMemoryKVStore) Close() error {
	//TODO implement me
	panic("implement me")
}

func NewInMemoryKVStore() *InMemoryKVStore {
	return &InMemoryKVStore{
		store: make(map[string][]byte),
	}
}

func (m *InMemoryKVStore) Get(key []byte) ([]byte, error) {
	val, ok := m.store[string(key)]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return val, nil
}

func (m *InMemoryKVStore) Put(key, value []byte) error {
	m.store[string(key)] = value
	return nil
}

func (m *InMemoryKVStore) Delete(key []byte) error {
	delete(m.store, string(key))
	return nil
}

func (m *InMemoryKVStore) Has(key []byte) (bool, error) {
	_, ok := m.store[string(key)]
	return ok, nil
}
