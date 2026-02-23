package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Block represents a blockchain block
type Block struct {
	Version       uint32
	PrevBlockHash string
	MerkleRoot    string
	Timestamp     int64
	Difficulty    uint32
	Nonce         uint64
	Transactions  []*Transaction
	Height        uint64
	BlockHash     string
	Miner         string // Miner's address
}

// NewBlock creates a new block
func NewBlock(prevHash string, transactions []*Transaction, difficulty uint32, height uint64, minerAddress string) *Block {
	block := &Block{
		Version:       1,
		PrevBlockHash: prevHash,
		Timestamp:     time.Now().Unix(),
		Difficulty:    difficulty,
		Nonce:         0,
		Transactions:  transactions,
		Height:        height,
		Miner:         minerAddress,
	}

	block.MerkleRoot = block.CalculateMerkleRoot()
	block.BlockHash = block.CalculateHash()
	return block
}

// CalculateMerkleRoot computes merkle tree root
func (b *Block) CalculateMerkleRoot() string {
	if len(b.Transactions) == 0 {
		return hex.EncodeToString(make([]byte, 32))
	}

	var txHashes []string
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.TxHash)
	}

	for len(txHashes) > 1 {
		if len(txHashes)%2 != 0 {
			txHashes = append(txHashes, txHashes[len(txHashes)-1])
		}

		var nextLevel []string
		for i := 0; i < len(txHashes); i += 2 {
			combined := txHashes[i] + txHashes[i+1]
			hash := sha256.Sum256([]byte(combined))
			nextLevel = append(nextLevel, hex.EncodeToString(hash[:]))
		}

		txHashes = nextLevel
	}

	return txHashes[0]
}

// CalculateHash computes block header hash
func (b *Block) CalculateHash() string {
	header := fmt.Sprintf(
		"%d%s%s%d%d%d%d",
		b.Version,
		b.PrevBlockHash,
		b.MerkleRoot,
		b.Timestamp,
		b.Difficulty,
		b.Nonce,
		b.Height,
	)

	hash := sha256.Sum256([]byte(header))
	hash = sha256.Sum256(hash[:])
	return hex.EncodeToString(hash[:])
}

// IsGenesisBlock checks if this is the genesis block
func (b *Block) IsGenesisBlock() bool {
	return b.Height == 0 && b.PrevBlockHash == "0"
}

// GetTransactionCount returns number of transactions in block
func (b *Block) GetTransactionCount() int {
	return len(b.Transactions)
}

// FindTransaction finds a transaction in the block
func (b *Block) FindTransaction(txHash string) *Transaction {
	for _, tx := range b.Transactions {
		if tx.TxHash == txHash {
			return tx
		}
	}
	return nil
}
