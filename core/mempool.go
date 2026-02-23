package core

import (
	"fmt"
	"sync"
	"time"
)

// MempoolTx wraps a transaction with metadata
type MempoolTx struct {
	Tx          *Transaction
	AddedTime   time.Time
	Fee         uint64
	Priority    int64 // Higher = more likely to be included
}

// Mempool manages pending transactions
type Mempool struct {
	mutex     sync.RWMutex
	txs       map[string]*MempoolTx
	maxSize   int
	maxAge    time.Duration
	feeTierSize int
}

// NewMempool creates a new mempool
func NewMempool(maxSize int, maxAge time.Duration) *Mempool {
	mp := &Mempool{
		txs:     make(map[string]*MempoolTx),
		maxSize: maxSize,
		maxAge:  maxAge,
		feeTierSize: 1000, // transactions per fee tier
	}

	// Start cleanup goroutine
	go mp.cleanupExpired()
	return mp
}

// AddTransaction adds a transaction to mempool
func (mp *Mempool) AddTransaction(tx *Transaction, fee uint64) error {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()

	if len(mp.txs) >= mp.maxSize {
		return fmt.Errorf("mempool full, max size: %d", mp.maxSize)
	}

	if _, exists := mp.txs[tx.TxHash]; exists {
		return fmt.Errorf("transaction already in mempool: %s", tx.TxHash)
	}

	priority := int64(fee) / int64(len(tx.Inputs)+len(tx.Outputs))

	mp.txs[tx.TxHash] = &MempoolTx{
		Tx:        tx,
		AddedTime: time.Now(),
		Fee:       fee,
		Priority:  priority,
	}

	return nil
}

// RemoveTransaction removes transaction from mempool
func (mp *Mempool) RemoveTransaction(txHash string) {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()
	delete(mp.txs, txHash)
}

// GetTransaction retrieves a transaction from mempool
func (mp *Mempool) GetTransaction(txHash string) *Transaction {
	mp.mutex.RLock()
	defer mp.mutex.RUnlock()

	if mempoolTx, exists := mp.txs[txHash]; exists {
		return mempoolTx.Tx
	}
	return nil
}

// GetTransactions returns transactions ordered by fee (highest first)
func (mp *Mempool) GetTransactions(limit int) []*Transaction {
	mp.mutex.RLock()
	defer mp.mutex.RUnlock()

	// Collect all transactions
	var txList []*MempoolTx
	for _, mempoolTx := range mp.txs {
		txList = append(txList, mempoolTx)
	}

	// Sort by priority/fee (simplified - use heap sort in production)
	// Placeholder: return in arbitrary order
	var result []*Transaction
	for i, mempoolTx := range txList {
		if i >= limit {
			break
		}
		result = append(result, mempoolTx.Tx)
	}

	return result
}

// Size returns number of transactions in mempool
func (mp *Mempool) Size() int {
	mp.mutex.RLock()
	defer mp.mutex.RUnlock()
	return len(mp.txs)
}

// cleanupExpired removes old transactions
func (mp *Mempool) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		mp.mutex.Lock()

		now := time.Now()
		for txHash, mempoolTx := range mp.txs {
			if now.Sub(mempoolTx.AddedTime) > mp.maxAge {
				delete(mp.txs, txHash)
			}
		}

		mp.mutex.Unlock()
	}
}

// HasTransaction checks if transaction exists in mempool
func (mp *Mempool) HasTransaction(txHash string) bool {
	mp.mutex.RLock()
	defer mp.mutex.RUnlock()
	_, exists := mp.txs[txHash]
	return exists
}
