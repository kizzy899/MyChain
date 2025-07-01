// txpool_test.go
package txpool

import (
	"CHAIN/common"
	"CHAIN/statedb"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func generateTx(nonce uint64, gasPrice uint64, priv *ecdsa.PrivateKey) *common.Transaction {
	to := common.Address{9, 9, 9}
	msg := []byte("dummy")
	tx := &common.Transaction{
		Fro:       common.Address{}, // 会由 From() 动态生成
		To:        &to,
		Nonce:     nonce,
		GasLimit:  21000,
		GasPrice:  big.NewInt(int64(gasPrice)),
		Value:     big.NewInt(100),
		Input:     msg,
		Signature: []byte("placeholder"),
	}

	hash := tx.Hash()
	sig, err := crypto.Sign(hash, priv)
	if err != nil {
		panic(err)
	}
	tx.Signature = sig
	tx.R = new(big.Int).SetBytes(sig[:32])
	tx.S = new(big.Int).SetBytes(sig[32:64])
	tx.V = uint8(sig[64]) + 27
	tx.Fro = tx.From() // 补上真正 From 地址
	return tx
}

func TestDefaultPool_Behaviors(t *testing.T) {
	stateDB := statedb.NewInMemoryStateDB()
	privKey, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(privKey.PublicKey)
	var a common.Address
	copy(a[:], addr[:])
	stateDB.Store(a, &common.Account{Nonce: 0})

	pool := NewDefaultPool(nil)
	pool.State = stateDB

	tx1 := generateTx(1, 10, privKey)
	pool.NewTx(tx1)
	if len(pool.pendings[a]) != 1 {
		t.Fatalf("expected 1 pending tx block, got %d", len(pool.pendings[a]))
	}

	tx3 := generateTx(3, 15, privKey)
	pool.NewTx(tx3)
	if len(pool.queue[a]) != 1 {
		t.Fatalf("expected 1 queued tx, got %d", len(pool.queue[a]))
	}

	tx2 := generateTx(2, 20, privKey)
	pool.NewTx(tx2)
	if len(pool.queue[a]) != 0 {
		t.Fatal("expected queue empty after inserting nonce=2")
	}
	if len(pool.pendings[a]) == 0 {
		t.Fatal("expected pending txs after inserting nonce=2")
	}

	tx2replace := generateTx(2, 30, privKey)
	pool.NewTx(tx2replace)

	found := false
	for _, blk := range pool.pendings[a] {
		for _, tx := range *blk {
			if tx.Nonce == 2 && tx.GasPrice.Uint64() == 30 {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("expected nonce=2 with GasPrice=30 in pending txs after replacement")
	}

	if pool.txs[0].GasPrice() != 10 {
		t.Fatalf("expected lowest GasPrice=10 at txs[0], got %d", pool.txs[0].GasPrice())
	}

	popTx := pool.Pop()
	if popTx == nil || popTx.Nonce != 1 {
		t.Fatalf("expected nonce=1 popped, got %+v", popTx)
	}

	popTx2 := pool.Pop()
	if popTx2 == nil || popTx2.Nonce != 2 {
		t.Fatalf("expected nonce=2 popped, got %+v", popTx2)
	}

	pool.NotifyTxEvent([]*common.Transaction{
		generateTx(100, 50, privKey),
	})
	t.Log("TestDefaultPool_Behaviors completed successfully")
}
