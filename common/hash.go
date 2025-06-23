package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

// Hash 表示一个32字节的哈希值
type Hash [32]byte

// BytesToHash 将字节切片转换为Hash
func BytesToHash(b []byte) Hash {
	var h Hash
	copy(h[:], b)
	return h
}

// String 实现Stringer接口
func (h Hash) String() string {
	return string(h[:])
}

func (tx *Transaction) Hash() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, tx.Nonce)
	binary.Write(buf, binary.BigEndian, tx.GasPrice)
	binary.Write(buf, binary.BigEndian, tx.GasLimit)

	if tx.To != nil {
		buf.Write(tx.To[:])
	} else {
		buf.Write(make([]byte, 20)) // Empty address
	}

	buf.Write(tx.Value.Bytes())
	buf.Write(tx.Input)

	hash := sha256.Sum256(buf.Bytes())
	return hash[:]
}

// Bytes 转换为字节切片
func (h Hash) Bytes() []byte {
	return h[:]
}

// Hex 返回16进制字符串表示
func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

// IsEmpty 判断是否为空哈希
func (h Hash) IsEmpty() bool {
	return h == Hash{}
}

// FromBytes 从字节切片创建Hash
func FromBytes(b []byte) Hash {
	var h Hash
	copy(h[:], b)
	return h
}
