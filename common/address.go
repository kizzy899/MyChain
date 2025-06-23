package common

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"math/big"
)

type Address [20]byte

func (a Address) String() string {
	return "0x" + hex.EncodeToString(a[:])
}

func RecoverAddress(hash []byte, r, s *big.Int, v byte) Address {
	// 构建 65 字节的签名 (r || s || v)
	sig := make([]byte, 65)
	copy(sig[32-len(r.Bytes()):32], r.Bytes())
	copy(sig[64-len(s.Bytes()):64], s.Bytes())
	sig[64] = v - 27 // v should be 0 or 1 for recovery in secp256k1

	pubKeyBytes, err := secp256k1.RecoverPubkey(hash, sig)
	if err != nil {
		panic(err) // 你可以更优雅地处理错误
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		panic(err)
	}

	// 计算地址：Keccak256(pubkey[1:])[12:]
	pubBytes := crypto.FromECDSAPub(pubKey)[1:] // 去掉前缀 0x04
	addressBytes := crypto.Keccak256(pubBytes)[12:]

	var addr Address
	copy(addr[:], addressBytes)
	return addr
}
