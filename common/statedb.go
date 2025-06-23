package common

import (
	"fmt"
	"math/big"
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

// GetAccount 获取某个地址的账户，如果不存在则返回 nil
func (db *StateDB) GetAccount(addr Address) *Account {
	db.lock.RLock()
	defer db.lock.RUnlock()
	return db.accounts[addr]
}

// CreateAccount 创建新账户，如果已存在则返回已存在账户
func (db *StateDB) CreateAccount(addr Address) *Account {
	db.lock.Lock()
	defer db.lock.Unlock()

	acct, exists := db.accounts[addr]
	if exists {
		return acct
	}

	newAcct := NewAccount(addr)
	db.accounts[addr] = newAcct
	return newAcct
}

// AddBalance 给指定账户增加余额（自动创建账户）
func (db *StateDB) AddBalance(addr Address, amount *big.Int) {
	acct := db.GetAccount(addr)
	if acct == nil {
		acct = db.CreateAccount(addr)
	}
	acct.AddBalance(amount)
}

// SubBalance 扣减余额
func (db *StateDB) SubBalance(addr Address, amount *big.Int) error {
	acct := db.GetAccount(addr)
	if acct == nil {
		return fmt.Errorf("account not found: %s", addr.String())
	}
	if acct.Balance.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance")
	}
	acct.SubBalance(amount)
	return nil
}

// GetBalance 查询账户余额
func (db *StateDB) GetBalance(addr Address) *big.Int {
	acct := db.GetAccount(addr)
	if acct == nil {
		return big.NewInt(0)
	}
	acct.lock.RLock()
	defer acct.lock.RUnlock()
	return new(big.Int).Set(acct.Balance)
}

// SetNonce 设置账户 nonce
func (db *StateDB) SetNonce(addr Address, nonce uint64) {
	acct := db.GetAccount(addr)
	if acct == nil {
		acct = db.CreateAccount(addr)
	}
	acct.SetNonce(nonce)
}

// GetNonce 获取账户 nonce
func (db *StateDB) GetNonce(addr Address) uint64 {
	acct := db.GetAccount(addr)
	if acct == nil {
		return 0
	}
	return acct.GetNonce()
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
