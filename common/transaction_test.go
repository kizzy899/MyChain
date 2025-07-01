package common_test

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"CHAIN/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// HexToAddress 同前
func HexToAddress(s string) common.Address {
	var addr common.Address
	b, _ := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	copy(addr[:], b)
	return addr
}

// 生成并签名交易，返回 *types.Transaction 和私钥对应的 common.Address
// 生成签名后的 common.Transaction，返回交易对象和对应地址
func createSignedCommonTx(t *testing.T) (*common.Transaction, common.Address) {
	privKey, err := ethcrypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	pubBytes := ethcrypto.FromECDSAPub(&privKey.PublicKey)[1:]
	fromBytes := ethcrypto.Keccak256(pubBytes)[12:]
	var from common.Address
	copy(from[:], fromBytes)

	// 1. 先构造未签名的 tx
	tx := &common.Transaction{
		Fro:   from,
		Value: big.NewInt(12345),
		Nonce: 0,
		Input: []byte{},
	}

	// 2. 获取签名哈希
	hash := tx.Hash()

	// 3. 对哈希签名
	sig, err := ethcrypto.Sign(hash, privKey)
	if err != nil {
		t.Fatal(err)
	}

	// 4. 拆解签名成 R,S,V
	R := new(big.Int).SetBytes(sig[:32])
	S := new(big.Int).SetBytes(sig[32:64])
	V := uint8(sig[64]) + 27

	// 5. 补充签名后的字段
	tx.R = R
	tx.S = S
	tx.V = V

	return tx, from
}

func TestCommonTransactionFrom(t *testing.T) {
	tx, from := createSignedCommonTx(t)

	recovered := tx.From()
	if recovered != from {
		t.Fatalf("From() 地址不匹配，期望 %x 实际 %x", from, recovered)
	}

	if len(tx.Hex()) == 0 {
		t.Error("Hex() 返回空字符串")
	}
}
