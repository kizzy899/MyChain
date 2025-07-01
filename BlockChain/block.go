package BlockChain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"

	"CHAIN/common" // 根据你的项目路径导入 common 包
)

type Block struct {
	Index        uint64
	Timestamp    int64
	PrevHash     []byte
	Hash         []byte
	Nonce        uint64
	Transactions []*common.Transaction // 修改这里
}

// NewBlock 创建新区块
func NewBlock(transactions []*common.Transaction, prevHash []byte, index uint64) *Block {
	block := &Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		PrevHash:     prevHash,
		Transactions: transactions,
	}
	block.Hash = block.CalculateHash()
	return block
}

// CalculateHash 计算区块哈希
func (b *Block) CalculateHash() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(b.Index)
	_ = enc.Encode(b.Timestamp)
	_ = enc.Encode(b.PrevHash)
	_ = enc.Encode(b.Transactions) // 这里自动序列化 []*common.Transaction
	_ = enc.Encode(b.Nonce)
	hash := sha256.Sum256(buf.Bytes())
	return hash[:]
}
