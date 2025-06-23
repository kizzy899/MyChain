package txpool

import (
	"CHAIN/common"
	"CHAIN/statedb"
	"fmt"
	"github.com/ethereum/go-ethereum/trie"
	"hash"
	"sort"
)

type SortedTxs interface { // 定义接口 SortedTxs，用于处理排序后的交易
	GasPrice() uint64               // GasPrice 方法，返回交易的 Gas 价格
	Push(tx *common.Transaction)    // Push 方法，向交易列表中添加交易
	Replace(tx *common.Transaction) // Replace 方法，替换交易列表中的交易
	Pop() *common.Transaction       // Pop 方法，从交易列表中弹出交易
	Nonce() uint64                  // Nonce 方法，返回交易的 Nonce 值
}

type pendingTxs []*DefaultSortedTxs // 定义类型 pendingTxs，代表待处理交易的列表集合

// 快排接口实现
func (p pendingTxs) Len() int           { return len(p) }
func (p pendingTxs) Less(i, j int) bool { return p[i].GasPrice() < p[j].GasPrice() }
func (p pendingTxs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// 定义结构体 DefaultPool，代表默认交易池
type DefaultPool struct {
	State    statedb.StateDB
	Stat     *trie.StateTrie
	all      map[hash.Hash]bool
	txs      pendingTxs
	pendings map[common.Address]pendingTxs
	queue    map[common.Address]QueueSortedTxs
}

func NewDefaultPool(state *trie.StateTrie) *DefaultPool { // 创建并返回一个新的 DefaultPool 实例
	return &DefaultPool{
		Stat:     state,
		all:      make(map[hash.Hash]bool),
		pendings: make(map[common.Address]pendingTxs),
		queue:    make(map[common.Address]QueueSortedTxs),
	}
}

type QueueSortedTxs []*common.Transaction

func (q QueueSortedTxs) Len() int           { return len(q) }
func (q QueueSortedTxs) Less(i, j int) bool { return q[i].Nonce < q[j].Nonce }
func (q QueueSortedTxs) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

type DefaultSortedTxs []*common.Transaction

func (sorted *DefaultSortedTxs) Push(tx *common.Transaction) {
	*sorted = append(*sorted, tx)
}

func (sorted DefaultSortedTxs) Replace(tx *common.Transaction) {
	for key, value := range sorted {
		if value.Nonce == tx.Nonce && tx.GasPrice() > 0 {
			sorted[key] = tx
		}
	}
}

func (sorted *DefaultSortedTxs) Pop() *common.Transaction {
	if len(*sorted) > 0 {
		tx := (*sorted)[0]
		*sorted = (*sorted)[1:]
		return tx
	}
	return nil
}

func (sorted DefaultSortedTxs) Nonce() uint64 {
	return sorted[len(sorted)-1].Nonce
}

func (sorted DefaultSortedTxs) GasPrice() uint64 {
	return sorted[0].GasPrice()
}

func (pool *DefaultPool) PrintfPool() {
	for _, txs := range pool.txs {
		fmt.Println("tx block")
		for _, tx := range *txs {
			fmt.Println("tx:", tx)
		}
	}
	for _, txs := range pool.pendings {
		for _, tx := range txs {
			for _, t := range *tx {
				fmt.Println("pending:", t)
			}
		}
	}
	for _, txs := range pool.queue {
		for _, tx := range txs {
			fmt.Println("queue:", tx)
		}
	}
}

func (pool *DefaultPool) SetStatRoot(root hash.Hash) {
	// 设置状态树的根哈希值（占位）
	// pool.Stat.SetRoot(root)
}

func (pool *DefaultPool) NewTx(tx *common.Transaction) {
	account := pool.State.Load(tx.From())
	if account.Nonce >= tx.Nonce {
		return
	}

	nonce := account.Nonce
	blks := pool.pendings[tx.From()]
	if len(blks) > 0 {
		last := blks[len(blks)-1]
		nonce = last.Nonce()
	}

	if tx.Nonce > nonce+1 {
		pool.addQueueTx(tx)
	} else if tx.Nonce == nonce+1 {
		pool.pushPendingTx(tx)
	} else {
		pool.replacePendingTx(tx)
	}
}

func (pool *DefaultPool) replacePendingTx(tx *common.Transaction) {
	blks := pool.pendings[tx.From()]
	for _, blk := range blks {
		if blk.Nonce() >= tx.Nonce {
			blk.Replace(tx)
			sort.Sort(blks)
			break
		}
	}
}

func (pool *DefaultPool) pushPendingTx(tx *common.Transaction) {
	blks := pool.pendings[tx.From()]
	if len(blks) == 0 {
		blk := &DefaultSortedTxs{tx}
		blks = append(blks, blk)
		pool.pendings[tx.From()] = blks
		pool.txs = append(pool.txs, blk)
		sort.Sort(pool.txs)
	} else {
		last := blks[len(blks)-1]
		if last.GasPrice() <= tx.GasPrice() {
			*last = append(*last, tx)
		} else {
			blk := &DefaultSortedTxs{tx}
			blks = append(blks, blk)
			pool.pendings[tx.From()] = blks
			pool.txs = append(pool.txs, blk)
			sort.Sort(pool.txs)
		}
	}

	queueTxs := pool.queue[tx.From()]
	nonce := tx.Nonce
	for i := 0; i < len(queueTxs); i++ {
		if queueTxs[i].Nonce == nonce+1 {
			nonce++
			nextTx := queueTxs[i]
			queueTxs = append(queueTxs[:i], queueTxs[i+1:]...)
			i--
			pool.queue[tx.From()] = queueTxs
			pool.pushPendingTx(nextTx)
		}
	}
}

func (pool *DefaultPool) addQueueTx(tx *common.Transaction) {
	txs := pool.queue[tx.From()]
	txs = append(txs, tx)
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Nonce < txs[j].Nonce
	})
	pool.queue[tx.From()] = txs
}

func (pool *DefaultPool) Pop() *common.Transaction {
	if len(pool.txs) == 0 || pool.txs[0] == nil {
		return nil
	}
	tx := pool.txs[0].Pop()
	if len(*pool.txs[0]) == 0 {
		pool.txs = pool.txs[1:]
	}
	return tx
}

func (pool *DefaultPool) NotifyTxEvent(txs []*common.Transaction) {
	for _, tx := range txs {
		fmt.Printf("NotifyTxEvent: New tx from %s with nonce %d and gas price %d\n",
			tx.Hex(),
			tx.Nonce,
			tx.GasPrice(),
		)
	}
}
