package storage

import "github.com/Mrcubys/VOIDEX-Network---Layer-1-Blockchain-Protocol/core"

// Storage interface defines database operations
type Storage interface {
	// Block operations
	StoreBlock(block *core.Block) error
	GetBlock(height uint64) (*core.Block, error)
	GetBlockByHash(hash string) (*core.Block, error)
	DeleteBlock(height uint64) error

	// UTXO operations
	StoreUTXO(utxo *core.UTXO) error
	GetUTXO(txHash string, outIndex uint32) (*core.UTXO, error)
	DeleteUTXO(txHash string, outIndex uint32) error
	GetAllUTXOs() ([]*core.UTXO, error)

	// Transaction operations
	StoreTx(tx *core.Transaction) error
	GetTx(txHash string) (*core.Transaction, error)
	GetTxsByAddress(address string) ([]*core.Transaction, error)

	// State operations
	StoreState(key string, value []byte) error
	GetState(key string) ([]byte, error)
	DeleteState(key string) error

	// Maintenance
	Close() error
	Backup() error
	GetStats() map[string]interface{}
}
