package kvstore

import (
	"errors"
	"sync"
)

// MemoryKVStore 是基于内存的键值存储实现
type MemoryKVStore struct {
	data map[string][]byte
	mu   sync.RWMutex
}

// NewMemoryKVStore 创建新的内存键值存储
func NewMemoryKVStore() KVStore {
	return &MemoryKVStore{
		data: make(map[string][]byte),
	}
}

func (m *MemoryKVStore) Put(key []byte, value []byte) error {
	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[string(key)] = value
	return nil
}

func (m *MemoryKVStore) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New("key cannot be empty")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.data[string(key)]
	if !exists {
		return nil, errors.New("key not found")
	}
	return value, nil
}

func (m *MemoryKVStore) Delete(key []byte) error {
	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, string(key))
	return nil
}

func (m *MemoryKVStore) Has(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, errors.New("key cannot be empty")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.data[string(key)]
	return exists, nil
}

func (m *MemoryKVStore) Close() error {
	// 清空数据
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = nil
	return nil
}
