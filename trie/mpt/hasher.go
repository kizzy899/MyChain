package trie

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Sha3_256 computes the Keccak-256 hash (Ethereum standard)
func Sha3_256(data []byte) common.Hash {
	return common.BytesToHash(crypto.Keccak256(data))
}

// HashNode is a helper to hash any Node implementation
func HashNode(node Node) common.Hash {
	if node == nil {
		return common.Hash{} // 空节点返回零值哈希
	}
	return Sha3_256(node.Serialize())
}
