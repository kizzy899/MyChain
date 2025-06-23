package txpool

import (
	"github.com/ethereum/go-ethereum/core/types"
	"hash"
)

type TxPool interface {
	NewTx(tx *types.Transaction)            //接收一个 `*types.Transaction` 类型的参数 `tx`，用于将新的交易加入到交易池中。
	PoP() *types.Transaction                //返回一个 `*types.Transaction` 类型的指针，可能是从交易池中弹出的交易。
	SetStatRoot(root hash.Hash)             //接收一个 `hash.Hash` 类型的参数 `root`，用于设置状态根。
	NotifyTxEvent(txs []*types.Transaction) //接收一个 `[]*types.Transaction` 类型的切片 `txs`，用于通知交易事件。
}
