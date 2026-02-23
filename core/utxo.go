package core

import (
	"fmt"
	"sync"
)

// UTXO represents an Unspent Transaction Output
type UTXO struct {
	TxHash    string
	OutIndex  uint32
	Value     uint64
	Address   string
	LockScript string
}

// UTXOKey creates a unique key for a UTXO
func (u *UTXO) Key() string {
	return fmt.Sprintf("%s:%d", u.TxHash, u.OutIndex)
}

// UTXOSet manages all unspent transaction outputs
type UTXOSet struct {
	mutex sync.RWMutex
	utxos map[string]*UTXO
}

// NewUTXOSet creates a new UTXO set
func NewUTXOSet() *UTXOSet {
	return &UTXOSet{
		utxos: make(map[string]*UTXO),
	}
}

// AddUTXO adds a UTXO to the set
func (us *UTXOSet) AddUTXO(utxo *UTXO) {
	us.mutex.Lock()
	defer us.mutex.Unlock()
	us.utxos[utxo.Key()] = utxo
}

// RemoveUTXO removes a UTXO from the set
func (us *UTXOSet) RemoveUTXO(txHash string, outIndex uint32) {
	us.mutex.Lock()
	defer us.mutex.Unlock()
	key := fmt.Sprintf("%s:%d", txHash, outIndex)
	delete(us.utxos, key)
}

// FindUTXO finds a UTXO by transaction hash and output index
func (us *UTXOSet) FindUTXO(txHash string, outIndex uint32) *UTXO {
	us.mutex.RLock()
	defer us.mutex.RUnlock()
	key := fmt.Sprintf("%s:%d", txHash, outIndex)
	return us.utxos[key]
}

// FindUTXOsByAddress finds all UTXOs belonging to an address
func (us *UTXOSet) FindUTXOsByAddress(address string) []*UTXO {
	us.mutex.RLock()
	defer us.mutex.RUnlock()

	var result []*UTXO
	for _, utxo := range us.utxos {
		if utxo.Address == address {
			result = append(result, utxo)
		}
	}
	return result
}

// GetBalance calculates balance for an address
func (us *UTXOSet) GetBalance(address string) uint64 {
	balance := uint64(0)
	for _, utxo := range us.FindUTXOsByAddress(address) {
		balance += utxo.Value
	}
	return balance
}

// Count returns total number of UTXOs
func (us *UTXOSet) Count() int {
	us.mutex.RLock()
	defer us.mutex.RUnlock()
	return len(us.utxos)
}

// GetAll returns all UTXOs (for syncing)
func (us *UTXOSet) GetAll() []*UTXO {
	us.mutex.RLock()
	defer us.mutex.RUnlock()

	result := make([]*UTXO, 0, len(us.utxos))
	for _, utxo := range us.utxos {
		result = append(result, utxo)
	}
	return result
}
