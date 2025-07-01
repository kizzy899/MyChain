package BlockChain

import (
	"encoding/hex"
	common2 "github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
	"testing"

	"CHAIN/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// 简单从 hex 字符串转 common.Address
func HexToAddress(s string) common.Address {
	var addr common.Address
	b, _ := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	copy(addr[:], b)
	return addr
}

// 生成一个简单的以太坊转账交易 (go-ethereum types.Transaction)
func newGethTx() *types.Transaction {
	nonce := uint64(0)
	toAddr := HexToAddress("0x0000000000000000000000000000000000000001") // 有效地址
	amount := big.NewInt(1000)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(1)
	data := []byte{}

	// 注意 types.NewTransaction 要求 to 是 common.Address 类型，但你的 common.Address 和 go-ethereum 的可能不同
	// 所以这里暂时用 nil 强转或自己定义
	// 这里用 nil 仅为测试，真实项目请替换为正确的地址类型
	return types.NewTransaction(nonce, common2.Address(toAddr), amount, gasLimit, gasPrice, data)
}

func TestNewBlock(t *testing.T) {
	from := HexToAddress("0x0000000000000000000000000000000000000002")
	to := HexToAddress("0x0000000000000000000000000000000000000003")
	toPtr := &to

	tx := &common.Transaction{
		Transaction: newGethTx(),
		Fro:         from,
		To:          toPtr,
		Value:       big.NewInt(10),
	}

	txs := []*common.Transaction{tx}

	prevHash := []byte("prevHash")
	index := uint64(1)

	block := NewBlock(txs, prevHash, index)
	if block == nil {
		t.Fatal("NewBlock 返回了 nil")
	}
	if len(block.Hash) == 0 {
		t.Error("区块 Hash 为空，计算失败")
	}
}
