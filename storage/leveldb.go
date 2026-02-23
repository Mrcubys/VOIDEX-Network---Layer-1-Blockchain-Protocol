package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Mrcubys/VOIDEX-Network/core"
	"github.com/syndtr/goleveldb/leveldb"
)

// LevelDBStorage implements Storage interface using LevelDB
type LevelDBStorage struct {
	db    *leveldb.DB
	path  string
	mutex sync.RWMutex
}

// NewLevelDBStorage creates a new LevelDB storage instance
func NewLevelDBStorage(dbPath string) (*LevelDBStorage, error) {
	// Ensure directory exists
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %v", err)
	}

	// Open database
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open leveldb: %v", err)
	}

	return &LevelDBStorage{
		db:   db,
		path: dbPath,
	}, nil
}

// StoreBlock saves a block to database
func (ls *LevelDBStorage) StoreBlock(block *core.Block) error {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	key := fmt.Sprintf("block:%d", block.Height)
	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	return ls.db.Put([]byte(key), data, nil)
}

// GetBlock retrieves a block from database
func (ls *LevelDBStorage) GetBlock(height uint64) (*core.Block, error) {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	key := fmt.Sprintf("block:%d", height)
	data, err := ls.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	var block core.Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, err
	}

	return &block, nil
}

// GetBlockByHash retrieves a block by hash
func (ls *LevelDBStorage) GetBlockByHash(hash string) (*core.Block, error) {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	key := fmt.Sprintf("blockhash:%s", hash)
	data, err := ls.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	var block core.Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, err
	}

	return &block, nil
}

// DeleteBlock removes a block from database
func (ls *LevelDBStorage) DeleteBlock(height uint64) error {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	key := fmt.Sprintf("block:%d", height)
	return ls.db.Delete([]byte(key), nil)
}

// StoreUTXO saves a UTXO to database
func (ls *LevelDBStorage) StoreUTXO(utxo *core.UTXO) error {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	key := fmt.Sprintf("utxo:%s:%d", utxo.TxHash, utxo.OutIndex)
	data, err := json.Marshal(utxo)
	if err != nil {
		return err
	}

	return ls.db.Put([]byte(key), data, nil)
}

// GetUTXO retrieves a UTXO from database
func (ls *LevelDBStorage) GetUTXO(txHash string, outIndex uint32) (*core.UTXO, error) {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	key := fmt.Sprintf("utxo:%s:%d", txHash, outIndex)
	data, err := ls.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	var utxo core.UTXO
	if err := json.Unmarshal(data, &utxo); err != nil {
		return nil, err
	}

	return &utxo, nil
}

// DeleteUTXO removes a UTXO from database
func (ls *LevelDBStorage) DeleteUTXO(txHash string, outIndex uint32) error {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	key := fmt.Sprintf("utxo:%s:%d", txHash, outIndex)
	return ls.db.Delete([]byte(key), nil)
}

// GetAllUTXOs retrieves all UTXOs (for syncing)
func (ls *LevelDBStorage) GetAllUTXOs() ([]*core.UTXO, error) {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	var utxos []*core.UTXO
	iter := ls.db.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		key := string(iter.Key())
		if key[:5] == "utxo:" {
			var utxo core.UTXO
			if err := json.Unmarshal(iter.Value(), &utxo); err != nil {
				continue
			}
			utxos = append(utxos, &utxo)
		}
	}

	return utxos, nil
}

// StoreTx saves a transaction to database
func (ls *LevelDBStorage) StoreTx(tx *core.Transaction) error {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	key := fmt.Sprintf("tx:%s", tx.TxHash)
	data, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	return ls.db.Put([]byte(key), data, nil)
}

// GetTx retrieves a transaction from database
func (ls *LevelDBStorage) GetTx(txHash string) (*core.Transaction, error) {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	key := fmt.Sprintf("tx:%s", txHash)
	data, err := ls.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	var tx core.Transaction
	if err := json.Unmarshal(data, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}

// GetTxsByAddress retrieves all transactions for an address
func (ls *LevelDBStorage) GetTxsByAddress(address string) ([]*core.Transaction, error) {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	var txs []*core.Transaction
	iter := ls.db.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		key := string(iter.Key())
		if key[:3] == "tx:" {
			var tx core.Transaction
			if err := json.Unmarshal(iter.Value(), &tx); err != nil {
				continue
			}

			// Check if transaction involves address
			for _, output := range tx.Outputs {
				if output.Address == address {
					txs = append(txs, &tx)
					break
				}
			}
		}
	}

	return txs, nil
}

// StoreState saves arbitrary state key-value pair
func (ls *LevelDBStorage) StoreState(key string, value []byte) error {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	stateKey := fmt.Sprintf("state:%s", key)
	return ls.db.Put([]byte(stateKey), value, nil)
}

// GetState retrieves state value
func (ls *LevelDBStorage) GetState(key string) ([]byte, error) {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	stateKey := fmt.Sprintf("state:%s", key)
	return ls.db.Get([]byte(stateKey), nil)
}

// DeleteState removes state key
func (ls *LevelDBStorage) DeleteState(key string) error {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	stateKey := fmt.Sprintf("state:%s", key)
	return ls.db.Delete([]byte(stateKey), nil)
}

// Close closes the database
func (ls *LevelDBStorage) Close() error {
	return ls.db.Close()
}

// Backup creates database backup
func (ls *LevelDBStorage) Backup() error {
	// In production, implement proper backup strategy
	return nil
}

// GetStats returns database statistics
func (ls *LevelDBStorage) GetStats() map[string]interface{} {
	stats := map[string]interface{}{
		"db_path": ls.path,
		"status":  "ok",
	}

	// Get approximate size
	// In production, use ls.db.Stats() if available
	return stats
}
