package common

import (
	"encoding/json"
	"math/big"
	"sync"
)

// Account 表示区块链上的一个账户
type Account struct {
	Address Address           // 账户地址
	Balance *big.Int          // 账户余额
	Nonce   uint64            // 发送过的交易数量，用于防止重放
	Code    []byte            // 合约账户代码（EOA则为nil）
	Storage map[string]string // 合约存储（可选，做MPT时常用）

	lock sync.RWMutex // 并发读写保护
}

// NewAccount 创建一个新账户
func NewAccount(address Address) *Account {
	return &Account{
		Address: address,
		Balance: big.NewInt(0),
		Nonce:   0,
		Code:    nil,
		Storage: make(map[string]string),
	}
}

// AddBalance 增加余额
func (a *Account) AddBalance(amount *big.Int) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Balance.Add(a.Balance, amount)
}

// SubBalance 扣减余额（不检查余额是否足够）
func (a *Account) SubBalance(amount *big.Int) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Balance.Sub(a.Balance, amount)
}

// SetNonce 设置账户的nonce
func (a *Account) SetNonce(nonce uint64) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Nonce = nonce
}

// GetNonce 获取账户的nonce
func (a *Account) GetNonce() uint64 {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.Nonce
}

// SetCode 设置合约代码
func (a *Account) SetCode(code []byte) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Code = code
}

// IsContract 判断账户是否是合约账户
func (a *Account) IsContract() bool {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return len(a.Code) > 0
}

// BytesToAccount 将字节数据反序列化为Account对象
func BytesToAccount(data []byte) (*Account, error) {
	var account Account
	err := json.Unmarshal(data, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// Bytes 将Account对象序列化为字节数组
func (a *Account) Bytes() ([]byte, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	// 使用结构体拷贝避免竞态条件
	accountCopy := struct {
		Address Address           `json:"address"`
		Balance *big.Int          `json:"balance"`
		Nonce   uint64            `json:"nonce"`
		Code    []byte            `json:"code,omitempty"`
		Storage map[string]string `json:"storage,omitempty"`
	}{
		Address: a.Address,
		Balance: new(big.Int).Set(a.Balance), // 深拷贝Balance
		Nonce:   a.Nonce,
		Code:    a.Code,
		Storage: a.Storage,
	}

	return json.Marshal(accountCopy)
}
