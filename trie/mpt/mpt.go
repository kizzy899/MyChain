package trie

import (
	"CHAIN/common"
	"CHAIN/kvstore"
	"bytes"
	"encoding/json"
	"fmt"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

// MPT 是Merkle Patricia Trie的主结构
// path类型都改成nibble

// MPT_KEY_NOT_FOUND 是键不存在时的错误
var MPT_KEY_NOT_FOUND = errors.New("MPT: Key not found")
var extNode *ExtensionNode

type MPT struct {
	Root Node
	db   kvstore.KVStore
}

// 添加获取根哈希的方法
func (m *MPT) RootHash() (common.Hash, error) {
	if m.Root == nil {
		return common.Hash{}, nil
	}
	return common.Hash(m.Root.GetHash()), nil
}

func (m *MPT) Commit() (common.Hash, error) {
	// 提交所有更改到数据库
	if m.Root == nil {
		return common.Hash{}, nil
	}
	rootHash := m.Root.GetHash()
	return common.Hash(rootHash), nil
}

func NewMPT(db kvstore.KVStore) *MPT {
	return &MPT{
		db: db,
	}
}

// FindLongestPrefix 查找与给定key有最长公共前缀的节点路径
func (m *MPT) FindLongestPrefix(key []Nibble) (paths [][]Nibble, nodes []Node) {
	var currentNode Node = m.Root

	for currentNode != nil {
		switch currentNode.GetType() {
		case LeafNodeType:
			leaf := currentNode.(*LeafNode)
			plength := prefixLength(key, leaf.Path)
			if plength == 0 {
				return
			}
			paths = append(paths, leaf.Path[:plength])
			nodes = append(nodes, leaf)
			return

		case ExtensionNodeType:
			ext := currentNode.(*ExtensionNode)
			plength := prefixLength(key, ext.Path)
			if plength == 0 {
				return
			}
			paths = append(paths, ext.Path[:plength])
			nodes = append(nodes, ext)
			currentNode = m.getNodeByHash(common.Hash(ext.Child))

		case BranchNodeType:
			return
		}
	}
	return
}

// Insert 非递归实现
func (m *MPT) Insert(key, value []byte) error {
	nibbles := convertToNibbles(key)
	valueHash := Sha3_256(value)

	// 新增：存储原始 value，key 为 valueHash
	if err := m.db.Put(valueHash.Bytes(), value); err != nil {
		return err
	}

	if m.Root == nil {
		leaf := NewLeafNode(nibbles, common.Hash(valueHash))
		m.Root = leaf
		m.db.Put(leaf.GetHash().Bytes(), leaf.Serialize())
		return nil
	}

	// 使用栈来管理待处理节点和路径
	type stackItem struct {
		node        Node
		parent      Node         // 父节点
		parentField *common.Hash // 父节点中指向当前节点的字段
		nibbles     []Nibble     // 剩余待处理的nibbles
	}

	var rootModified bool
	var newRoot Node
	stack := []stackItem{{
		node:    m.Root,
		nibbles: nibbles,
	}}

	for len(stack) > 0 {
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		currentNode := item.node
		remainingNibbles := item.nibbles

		switch node := currentNode.(type) {
		case *LeafNode:
			// 查找公共前缀
			commonPrefix := findCommonPrefix(node.Path, remainingNibbles)

			// 完全匹配则更新值
			if len(commonPrefix) == len(node.Path) && len(commonPrefix) == len(remainingNibbles) {
				node.Value = valueHash
				m.db.Put(node.GetHash().Bytes(), node.Serialize())
				if item.parent != nil {
					*item.parentField = m.storeNode(node)
				} else {
					rootModified = true
					newRoot = node
				}
				continue
			}

			// 创建分支节点处理部分匹配
			branch := &BranchNode{}

			// 处理现有叶子节点
			if len(node.Path) > len(commonPrefix) {
				nibble := node.Path[len(commonPrefix)]
				remainingPath := node.Path[len(commonPrefix)+1:]

				leaf := &LeafNode{
					Path:  remainingPath,
					Value: node.Value,
				}
				branch.Children[nibble] = common2.Hash(m.storeNode(leaf))
			} else {
				branch.Children[15] = common2.Hash(m.storeNode(&LeafNode{Value: node.Value}))
			}

			// 处理新插入的值
			if len(remainingNibbles) > len(commonPrefix) {
				nibble := remainingNibbles[len(commonPrefix)]
				remainingPath := remainingNibbles[len(commonPrefix)+1:]

				leaf := &LeafNode{
					Path:  remainingPath,
					Value: valueHash,
				}
				branch.Children[nibble] = common2.Hash(m.storeNode(leaf))
			} else {
				branch.Children[15] = common2.Hash(m.storeNode(&LeafNode{Value: valueHash}))
			}

			// 处理共同前缀
			resultNode := Node(branch)
			if len(commonPrefix) > 0 {
				extNode := &ExtensionNode{
					Path:  commonPrefix,
					Child: common2.Hash(m.storeNode(branch)),
				}
				resultNode = extNode
			}

			// 更新父节点或根节点
			if item.parent != nil {
				*item.parentField = m.storeNode(resultNode)
			} else {
				rootModified = true
				newRoot = resultNode
			}

			// 持久化新节点
			m.db.Put(branch.GetHash().Bytes(), branch.Serialize())
			if len(commonPrefix) > 0 {
				m.db.Put(resultNode.GetHash().Bytes(), resultNode.Serialize())
			}

		case *ExtensionNode:
			commonPrefix := findCommonPrefix(node.Path, remainingNibbles)

			// 完全匹配则继续处理子节点
			if len(commonPrefix) == len(node.Path) {
				child, err := m.loadNode(common.Hash(node.Child))
				if err != nil {
					return err
				}

				stack = append(stack, stackItem{
					node:        child,
					parent:      node,
					parentField: (*common.Hash)(&node.Child),
					nibbles:     remainingNibbles[len(commonPrefix):],
				})
				continue
			}

			// 部分匹配需要拆分扩展节点
			branch := &BranchNode{}

			// 处理原扩展节点剩余部分
			if len(node.Path) > len(commonPrefix) {
				nibble := node.Path[len(commonPrefix)]
				remainingPath := node.Path[len(commonPrefix)+1:]

				ext := &ExtensionNode{
					Path:  remainingPath,
					Child: node.Child,
				}
				branch.Children[nibble] = common2.Hash(m.storeNode(ext))
			} else {
				child, err := m.loadNode(common.Hash(node.Child))
				if err != nil {
					return err
				}
				branch.Children[15] = common2.Hash(m.storeNode(child))
			}

			// 处理新插入的值
			if len(remainingNibbles) > len(commonPrefix) {
				nibble := remainingNibbles[len(commonPrefix)]
				remainingPath := remainingNibbles[len(commonPrefix)+1:]

				leaf := &LeafNode{
					Path:  remainingPath,
					Value: valueHash,
				}
				branch.Children[nibble] = common2.Hash(m.storeNode(leaf))
			} else {
				branch.Children[15] = common2.Hash(m.storeNode(&LeafNode{Value: valueHash}))
			}

			// 处理共同前缀
			resultNode := Node(branch)
			if len(commonPrefix) > 0 {
				extNode := &ExtensionNode{
					Path:  commonPrefix,
					Child: common2.Hash(m.storeNode(branch)),
				}
				resultNode = extNode
			}

			// 更新父节点或根节点
			if item.parent != nil {
				*item.parentField = m.storeNode(resultNode)
			} else {
				rootModified = true
				newRoot = resultNode
			}

			// 持久化新节点
			m.db.Put(branch.GetHash().Bytes(), branch.Serialize())
			if len(commonPrefix) > 0 {
				m.db.Put(resultNode.GetHash().Bytes(), resultNode.Serialize())
			}

		case *BranchNode:
			if len(remainingNibbles) == 0 {
				// 更新value位置
				node.Children[15] = common2.Hash(m.storeNode(&LeafNode{Value: valueHash}))
				m.db.Put(node.GetHash().Bytes(), node.Serialize())

				if item.parent != nil {
					*item.parentField = m.storeNode(node)
				} else {
					rootModified = true
					newRoot = node
				}
				continue
			}

			nibble := remainingNibbles[0]
			var child Node
			var err error

			if childHash := node.Children[nibble]; !bytes.Equal(childHash[:], make([]byte, 32)) {
				child, err = m.loadNode(common.Hash(childHash))
				if err != nil {
					return err
				}
			} else {
				child = nil
			}

			stack = append(stack, stackItem{
				node:        child,
				parent:      node,
				parentField: (*common.Hash)(&node.Children[nibble]),
				nibbles:     remainingNibbles[1:],
			})
		}
	}

	if rootModified {
		m.Root = newRoot
		fmt.Printf("Set newRoot successfully: %#v\n", newRoot)
	}

	return nil
}

func (m *MPT) loadNode(hash common.Hash) (Node, error) {
	node := m.getNodeByHash(hash)
	if node == nil {
		return nil, fmt.Errorf("not found")
	}
	return node, nil
}

func (m *MPT) storeNode(node Node) common.Hash {
	hash := node.GetHash()
	m.db.Put(hash.Bytes(), node.Serialize())
	return common.Hash(hash)
}

func NewLeafNode(path []Nibble, valueHash common.Hash) *LeafNode {
	return &LeafNode{
		NodeType: LeafNodeType,
		Path:     path,
		Value:    common2.Hash(valueHash),
	}
}

// 通过节点的哈希值从底层存储中加载节点
func (m *MPT) getNodeByHash(hash common.Hash) Node {
	if hash == (common.Hash{}) {

		fmt.Println("getNodeByHash: empty hash")
		return nil
	}
	data, err := m.db.Get(hash.Bytes())
	if err != nil || len(data) == 0 {
		fmt.Printf("getNodeByHash: not found data for hash %x, err=%v\n", hash.Bytes(), err)
		return nil
	}
	fmt.Printf("getNodeByHash: got data len=%d for hash %x\n", len(data), hash.Bytes())

	var nodeType struct {
		NodeType NodeType `json:"type"`
	}
	if err := json.Unmarshal(data, &nodeType); err != nil {
		return nil
	}

	switch nodeType.NodeType {
	case LeafNodeType:
		var leaf LeafNode
		if err := json.Unmarshal(data, &leaf); err != nil {
			return nil
		}
		return &leaf
	case ExtensionNodeType:
		var ext ExtensionNode
		if err := json.Unmarshal(data, &ext); err != nil {
			return nil
		}
		return &ext
	case BranchNodeType:
		var branch BranchNode
		if err := json.Unmarshal(data, &branch); err != nil {
			return nil
		}
		return &branch
	default:
		return nil
	}
}

// prefixLength 计算两个nibble数组的公共前缀长度
func prefixLength(a, b []Nibble) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return minLen
}

// Search 在MPT中查找指定key对应的value
func (m *MPT) Search(key []byte) ([]byte, error) {
	nibbles := convertToNibbles(key)
	node := m.Root

	for node != nil {
		switch n := node.(type) {
		case *LeafNode:
			fmt.Printf("At LeafNode with path: %v\n", n.Path)
			fmt.Printf("Searching for nibbles: %v\n", nibbles)

			// 先检查是否完全相等
			if nibblesEqual(n.Path, nibbles) {
				val, err := m.db.Get(n.Value[:])
				if err != nil {
					return nil, err
				}
				return val, nil
			}

			// 如果不等，尝试打印公共前缀长度，帮助调试
			prefixLen := prefixLength(n.Path, nibbles)
			fmt.Printf("Prefix length between node path and key: %d\n", prefixLen)

			return nil, MPT_KEY_NOT_FOUND

		case *ExtensionNode:
			if len(nibbles) < len(n.Path) || !nibblesEqual(nibbles[:len(n.Path)], n.Path) {
				return nil, MPT_KEY_NOT_FOUND
			}
			childNode := m.getNodeByHash(common.Hash(n.Child))
			if childNode == nil {
				return nil, MPT_KEY_NOT_FOUND
			}
			node = childNode
			nibbles = nibbles[len(n.Path):]

		case *BranchNode:
			if len(nibbles) == 0 {
				return m.db.Get(n.Children[15][:]) // value 在 branch 的第 15 槽
			}
			next := n.Children[nibbles[0]]
			childNode := m.getNodeByHash(common.Hash(next))
			if childNode == nil {
				return nil, MPT_KEY_NOT_FOUND
			}
			node = childNode
			nibbles = nibbles[1:]

		default:
			return nil, MPT_KEY_NOT_FOUND
		}
	}

	return nil, MPT_KEY_NOT_FOUND
}

// 辅助函数
func joinNibblePaths(pathes [][]Nibble) []Nibble {
	var result []Nibble
	for _, path := range pathes {
		result = append(result, path...)
	}
	return result
}

func nibblesEqual(a, b []Nibble) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// findCommonPrefix 查找两个nibble数组的公共前缀
func findCommonPrefix(a, b []Nibble) []Nibble {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return a[:i]
		}
	}
	return a[:minLen]
}
