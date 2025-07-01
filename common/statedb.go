package common

import (
	"sync"
)

// StateDB 是账户状态管理器
type StateDB struct {
	accounts map[Address]*Account
	lock     sync.RWMutex
}

// NewStateDB 创建一个新的状态数据库
func NewStateDB() *StateDB {
	return &StateDB{
		accounts: make(map[Address]*Account),
	}
}

//mpt功能的

// StateDB 状态数据库接口
type statedb_mpt interface {
	// Get 获取指定地址关联的账户
	// 如果账户不存在则返回错误
	// 账户使用 common.Account 结构体表示
	// 地址使用 common.Address 结构体表示
	Get(address Address) (*Account, error)

	// Set 存储指定地址关联的账户
	// 如果账户已存在则返回错误
	Set(address Address, account *Account) error

	// Root 返回状态数据库的根哈希
	// 用于验证状态数据库的完整性
	// 根哈希使用字节切片表示
	Root() (Hash, error)
}
