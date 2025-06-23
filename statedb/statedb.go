package statedb

import (
	"CHAIN/common"
	"hash"
)

// 构造函数
func NewInMemoryStateDB() *InMemoryStateDB {
	return &InMemoryStateDB{
		accounts: make(map[common.Address]*common.Account),
	}
}

type InMemoryStateDB struct {
	root     hash.Hash
	accounts map[common.Address]*common.Account
}

type StateDB interface {
	SetRoot(root hash.Hash)
	Load(address common.Address) *common.Account
	Store(address common.Address, account *common.Account)
}

// 实现接口方法
// SetRoot 设置当前状态树的根哈希（可以用于快照、验证等）
func (db *InMemoryStateDB) SetRoot(root hash.Hash) {
	db.root = root
}

// Load 读取指定地址的账户信息，返回 *common.Account（不存在则返回 nil）
func (db *InMemoryStateDB) Load(address common.Address) *common.Account {
	if acc, ok := db.accounts[address]; ok {
		return acc
	}
	return nil
}

// Store 将账户信息写入状态数据库（更新或插入）
func (db *InMemoryStateDB) Store(address common.Address, account *common.Account) {
	db.accounts[address] = account
}
