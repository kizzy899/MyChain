package trie

import (
	"CHAIN/common"
	"CHAIN/kvstore"
	trie "CHAIN/trie/mpt"
)

// StateDBMPT 是基于 MPT 实现的状态数据库
type StateDBMPT struct {
	trie *trie.MPT
}

// NewStateDBMPT 创建一个新的 StateDBMPT，底层初始化 MPT
func NewStateDBMPT(db kvstore.KVStore) *StateDBMPT {
	return &StateDBMPT{
		trie: trie.NewMPT(db),
	}
}

// Get 根据地址获取账户信息
func (s *StateDBMPT) Get(address common.Address) (*common.Account, error) {
	key := address[:]
	value, err := s.trie.Search(key)
	if err != nil {
		return nil, err
	}
	return common.BytesToAccount(value)
}

// Set 设置或更新地址对应的账户信息
func (s *StateDBMPT) Set(address common.Address, account *common.Account) error {
	key := address[:]
	value, err := account.Bytes()
	if err != nil {
		return err
	}
	return s.trie.Insert(key, value)
}

// Root 返回当前状态树的根哈希
func (s *StateDBMPT) Root() (common.Hash, error) {
	return s.trie.RootHash()
}
