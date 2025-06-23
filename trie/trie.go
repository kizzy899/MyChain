package trie

import "CHAIN/common"

type Trie interface {
	// Insert inserts a key-value pair into the trie.
	Insert(key, value []byte) error

	// Search searches for a value associated with the given key.
	Search(key []byte) ([]byte, error)

	// Root returns the root hash of the trie.
	// This is useful for verifying the integrity of the trie.
	Root() (common.Hash, error)
}
