package main

import (
	"CHAIN/BlockChain"
	"CHAIN/common"
	"CHAIN/statedb"
	"CHAIN/txpool"
	"fmt"
	"math/big"
)

func main() {
	fmt.Println("🚀 启动简易区块链...")

	// 初始化状态数据库
	stateDB := statedb.NewInMemoryStateDB()

	addrA := common.Address{1, 2, 3}
	addrB := common.Address{4, 5, 6}

	// 初始化账户
	stateDB.Store(addrA, &common.Account{
		Address: addrA,
		Balance: big.NewInt(1000),
		Nonce:   0,
	})
	stateDB.Store(addrB, &common.Account{
		Address: addrB,
		Balance: big.NewInt(0),
		Nonce:   0,
	})

	// 初始化交易池
	pool := txpool.NewDefaultPool(nil)
	pool.State = stateDB

	// 创建交易 1（nonce = 1）
	tx1 := &common.Transaction{
		Fro:      addrA,
		To:       &addrB,
		Value:    big.NewInt(100),
		GasLimit: 21000,
		GasPrice: big.NewInt(1),
		Nonce:    1,
		R:        big.NewInt(1),
		S:        big.NewInt(2),
		V:        27,
		Input:    []byte{},
	}
	pool.NewTx(tx1)

	// 创建交易 2（nonce = 2）
	tx2 := &common.Transaction{
		Fro:      addrA,
		To:       &addrB,
		Value:    big.NewInt(200),
		GasLimit: 21000,
		GasPrice: big.NewInt(2),
		Nonce:    2,
		R:        big.NewInt(3),
		S:        big.NewInt(4),
		V:        28,
		Input:    []byte("data"),
	}
	pool.NewTx(tx2)

	// 初始化区块链（创世块）
	var chain []*BlockChain.Block
	genesis := BlockChain.NewBlock(nil, nil, 0)
	chain = append(chain, genesis)

	// 从交易池获取所有待打包交易
	var txs []*common.Transaction
	for {
		t := pool.Pop()
		if t == nil {
			break
		}
		txs = append(txs, t)
		// 应用交易结果（简单模拟转账逻辑）
		_ = stateDB.SubBalance(t.Fro, t.Value)
		stateDB.AddBalance(*t.To, t.Value)
		stateDB.SetNonce(t.Fro, stateDB.GetNonce(t.Fro)+1)
	}

	// 打包新区块
	prev := chain[len(chain)-1]
	block := BlockChain.NewBlock(txs, prev.Hash, prev.Index+1)
	chain = append(chain, block)

	fmt.Println("✅ 区块链当前高度：", block.Index)
	fmt.Println("🧾 当前区块交易数量：", len(block.Transactions))
	fmt.Println("📦 当前链长度：", len(chain))

	// 输出账户状态
	fmt.Println("账户 A 余额:", stateDB.GetBalance(addrA))
	fmt.Println("账户 A Nonce:", stateDB.GetNonce(addrA))
	fmt.Println("账户 B 余额:", stateDB.GetBalance(addrB))
}
