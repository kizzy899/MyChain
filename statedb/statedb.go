package statedb

import (
	"CHAIN/common"
	"fmt"
	"hash"
	"math/big"
	"sync"
)

// InMemoryStateDB 是状态数据库的内存实现
type InMemoryStateDB struct {
	root     hash.Hash
	accounts map[common.Address]*common.Account
	lock     sync.RWMutex
}

// 构造函数
func NewInMemoryStateDB() *InMemoryStateDB {
	return &InMemoryStateDB{
		accounts: make(map[common.Address]*common.Account),
	}
}

// 接口定义
type StateDB interface {
	SetRoot(root hash.Hash)
	Load(address common.Address) *common.Account
	Store(address common.Address, account *common.Account)
}

// SetRoot 设置根哈希
func (db *InMemoryStateDB) SetRoot(root hash.Hash) {
	db.root = root
}

// Load 读取账户
func (db *InMemoryStateDB) Load(address common.Address) *common.Account {
	db.lock.RLock()
	defer db.lock.RUnlock()
	return db.accounts[address]
}

// Store 存储账户
func (db *InMemoryStateDB) Store(address common.Address, account *common.Account) {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.accounts[address] = account
}

// 获取账户（内部方法）
func (db *InMemoryStateDB) GetAccount(addr common.Address) *common.Account {
	db.lock.RLock()
	defer db.lock.RUnlock()
	return db.accounts[addr]
}

// 创建账户
func (db *InMemoryStateDB) CreateAccount(addr common.Address) *common.Account {
	db.lock.Lock()
	defer db.lock.Unlock()

	acct, exists := db.accounts[addr]
	if exists {
		return acct
	}

	newAcct := common.NewAccount(addr)
	db.accounts[addr] = newAcct
	return newAcct
}

// 增加余额
func (db *InMemoryStateDB) AddBalance(addr common.Address, amount *big.Int) {
	acct := db.GetAccount(addr)
	if acct == nil {
		acct = db.CreateAccount(addr)
	}
	acct.AddBalance(amount)
}

// 扣减余额
func (db *InMemoryStateDB) SubBalance(addr common.Address, amount *big.Int) error {
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

// 查询余额
func (db *InMemoryStateDB) GetBalance(addr common.Address) *big.Int {
	acct := db.GetAccount(addr)
	if acct == nil {
		return big.NewInt(0)
	}
	acct.Lock()
	defer acct.Unlock()
	return new(big.Int).Set(acct.Balance)
}

// 设置 Nonce
func (db *InMemoryStateDB) SetNonce(addr common.Address, nonce uint64) {
	acct := db.GetAccount(addr)
	if acct == nil {
		acct = db.CreateAccount(addr)
	}
	acct.SetNonce(nonce)
}

// 获取 Nonce
func (db *InMemoryStateDB) GetNonce(addr common.Address) uint64 {
	acct := db.GetAccount(addr)
	if acct == nil {
		return 0
	}
	return acct.GetNonce()
}
