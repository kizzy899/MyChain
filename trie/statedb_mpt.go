package trie

import (
	"CHAIN/common"
	"CHAIN/trie/mpt"
)

type statedbmpt struct {
	// The root hash of the state databasetrie trie.Trie
	trie *mpt.MPT
}

// Get 获取指定地址的账户信息
func (s *statedbmpt) Get(address common.Address) (*common.Account, error) {
	key := address[:]
	value, err := s.trie.Search(key)
	if err != nil {
		return nil, err
	}
	return common.BytesToAccount(value)
}

// Set 设置/更新指定地址的账户信息
func (s *statedbmpt) Set(address common.Address, account *common.Account) error {
	key := address[:]
	value, _ := account.Bytes()
	return s.trie.Insert(key, value)
}

// Root 获取当前状态树的根哈希
func (s *statedbmpt) Root() (common.Hash, error) {
	return s.trie.RootHash()
}
