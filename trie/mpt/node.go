package trie

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
)

// NodeType 表示MPT节点类型
type NodeType int

const (
	LeafNodeType NodeType = iota
	ExtensionNodeType
	BranchNodeType
)

// Node 是所有MPT节点的基础接口
type Node interface {
	GetType() NodeType
	Serialize() []byte
	GetHash() common.Hash
}

// LeafNode 存储最终键值对的叶子节点
type LeafNode struct {
	NodeType NodeType    `json:"type"`
	Value    common.Hash `json:"value"`
	Path     []Nibble    `json:"path"`
}

func (n *LeafNode) GetType() NodeType { return LeafNodeType }

func (n *LeafNode) Serialize() []byte {
	data, _ := json.Marshal(n)
	return data
}

func (n *LeafNode) GetHash() common.Hash {
	return Sha3_256(n.Serialize())
}

// ExtensionNode 存储路径前缀和子节点哈希的扩展节点
type ExtensionNode struct {
	NodeType NodeType    `json:"type"`
	Path     []Nibble    `json:"path"`
	Child    common.Hash `json:"child"`
}

func (n *ExtensionNode) GetType() NodeType { return ExtensionNodeType }

func (n *ExtensionNode) Serialize() []byte {
	data, _ := json.Marshal(n)
	return data
}

func (n *ExtensionNode) GetHash() common.Hash {
	return Sha3_256(n.Serialize())
}

// BranchNode 包含16个子节点的分支节点
type BranchNode struct {
	NodeType NodeType        `json:"type"`
	Children [16]common.Hash `json:"children"`
}

func (n *BranchNode) GetType() NodeType { return BranchNodeType }

func (n *BranchNode) Serialize() []byte {
	data, _ := json.Marshal(n)
	return data
}

func (n *BranchNode) GetHash() common.Hash {
	return Sha3_256(n.Serialize())
}
