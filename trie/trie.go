package trie

import (
	"CHAIN/common"
	"fmt"
)

type Trie interface {
	// Insert inserts a key-value pair into the trie.
	Insert(key, value []byte) error

	// Search searches for a value associated with the given key.
	Search(key []byte) ([]byte, error)

	// Root returns the root hash of the trie.
	// This is useful for verifying the integrity of the trie.
	Root() (common.Hash, error)
}

// MPT 实现了 Trie 接口
type MPT struct {
	// 你的内部存储，比如节点map
	nodes    map[string][]byte
	rootHash common.Hash
}

// Insert 实现 Trie 接口
func (m *MPT) Insert(key, value []byte) error {
	if m.nodes == nil {
		m.nodes = make(map[string][]byte)
	}
	m.nodes[string(key)] = value
	// 更新根哈希，这里简化直接用 key 做hash
	h := common.Hash(key)     // h 是一个 common.Hash 类型（[32]byte）
	copy(m.rootHash[:], h[:]) // h[:] 现在是 []byte，能用 copy 了

	return nil
}

// Search 实现 Trie 接口
func (m *MPT) Search(key []byte) ([]byte, error) {
	val, ok := m.nodes[string(key)]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return val, nil
}

// Root 实现 Trie 接口
func (m *MPT) Root() (common.Hash, error) {
	return m.rootHash, nil
}

// NewTrie 返回一个 Trie 接口实例
func NewTrie() Trie {
	return &MPT{
		nodes: make(map[string][]byte),
	}
}
