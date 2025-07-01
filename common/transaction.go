package common

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"math/big"
)

type Transaction struct {
	*types.Transaction // 嵌入 go-ethereum 的 Transaction，
	R, S               *big.Int
	V                  uint8
	GasPrice           *big.Int
	// 基础字段
	Fro Address  // 发送方地址
	To  *Address // 接收方地址(合约创建时为nil)

	// 新增的核心字段
	GasLimit uint64   // 交易消耗的Gas上限
	Value    *big.Int // 转账金额(wei)
	Input    []byte   // 交易输入数据(合约调用时使用)

	// 其他原有字段...
	Nonce     uint64 // 交易序号
	Signature []byte // 交易签
}

// NewTransaction 构造函数：从 types.Transaction 复制构造 common.Transaction
func NewTransaction(tx *types.Transaction, R, S *big.Int, V uint8) *Transaction {
	return &Transaction{
		Transaction: tx,
		R:           R,
		S:           S,
		V:           V,
	}
}

// From 返回发送者地址，通过签名恢复公钥再转地址
// From 返回发送者地址，通过签名恢复公钥再转地址
func (tx *Transaction) From() Address {
	hash := tx.Hash() // 1. 获取交易哈希（不包含签名部分）

	// 2. 构造 65 字节签名数据：r||s||v
	sig := make([]byte, 65)
	copy(sig[32-len(tx.R.Bytes()):32], tx.R.Bytes()) // r 填充到前 32 字节
	copy(sig[64-len(tx.S.Bytes()):64], tx.S.Bytes()) // s 填充到中间 32 字节

	// 支持 EIP-155 计算 recovery id
	var recID byte
	if tx.V == 27 || tx.V == 28 {
		recID = byte(tx.V - 27)
	} else if tx.V >= 35 {
		recID = byte((uint64(tx.V) - 35) % 2)
	} else {
		panic(fmt.Sprintf("invalid V value for signature recovery: %d", tx.V))
	}

	sig[64] = recID

	// 3. 恢复未压缩公钥
	pubKeyBytes, err := secp256k1.RecoverPubkey(hash, sig)
	if err != nil {
		panic(fmt.Sprintf("signature recovery failed: %v", err))
	}

	// 4. 解码为 ecdsa 公钥结构
	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		panic(fmt.Sprintf("invalid pubkey bytes: %v", err))
	}

	// 5. 获取地址（keccak256(pubkey[1:])[12:]）
	pubBytes := crypto.FromECDSAPub(pubKey)[1:] // 跳过前缀 0x04
	addressBytes := crypto.Keccak256(pubBytes)[12:]

	var addr Address
	copy(addr[:], addressBytes)
	return addr
}

func (tx *Transaction) GasPriceUint64() uint64 {
	if tx.GasPrice == nil {
		return 0
	}
	return tx.GasPrice.Uint64()
}

// Hex 返回交易哈希的十六进制字符串表示
func (tx *Transaction) Hex() string {
	return hex.EncodeToString(tx.Hash())
}
