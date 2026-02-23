package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/Mrcubys/VOIDEX-Network---Layer-1-Blockchain-Protocol/tree/main/consensus"
	"github.com/Mrcubys/VOIDEX-Network---Layer-1-Blockchain-Protocol/tree/main/storage"
)

// Blockchain manages the chain of blocks
type Blockchain struct {
	mutex              sync.RWMutex
	Blocks             []*Block
	UTXOSet            *UTXOSet
	Difficulty         uint32
	PendingTransactions *Mempool
	Chain              storage.Storage
	RewardCalculator   *consensus.BlockRewardCalculator
	DifficultyAdjuster *consensus.DifficultyAdjuster
	GenesisConfig      *GenesisConfig
	BlockTimestamps    []int64
}

// NewBlockchain creates a new blockchain
func NewBlockchain(store storage.Storage, minerAddress string) (*Blockchain, error) {
	genesis := GetMainnetGenesis()

	bc := &Blockchain{
		Blocks:             make([]*Block, 0),
		UTXOSet:            NewUTXOSet(),
		Difficulty:         0x00000FFF,
		PendingTransactions: NewMempool(10000, 24*time.Hour),
		Chain:              store,
		RewardCalculator:   consensus.NewBlockRewardCalculator(genesis.InitialReward, genesis.RewardHalvingInterval, genesis.MaxSupply),
		DifficultyAdjuster: consensus.NewDifficultyAdjuster(genesis.DifficultyWindow, genesis.TargetBlockTime),
		GenesisConfig:      genesis,
		BlockTimestamps:    make([]int64, 0),
	}

	// Check if genesis block exists in storage
	genesisFromStore, err := store.GetBlock(0)
	if err != nil || genesisFromStore == nil {
		// Create and add genesis block
		genesisBlock := CreateGenesisBlock(minerAddress)
		err := bc.AddBlock(genesisBlock)
		if err != nil {
			return nil, fmt.Errorf("failed to create genesis block: %v", err)
		}
	} else {
		// Load blockchain from storage
		bc.Blocks = append(bc.Blocks, genesisFromStore)
		bc.BlockTimestamps = append(bc.BlockTimestamps, genesisFromStore.Timestamp)
		bc.rebuildUTXOSet()
	}

	return bc, nil
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(block *Block) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Validate block
	if err := bc.validateBlock(block); err != nil {
		return fmt.Errorf("block validation failed: %v", err)
	}

	// Add transactions to UTXO set
	for _, tx := range block.Transactions {
		for _, output := range tx.Outputs {
			utxo := &UTXO{
				TxHash:    tx.TxHash,
				OutIndex:  uint32(len(tx.Outputs)),
				Value:     output.Value,
				Address:   output.Address,
				LockScript: output.LockScript,
			}
			bc.UTXOSet.AddUTXO(utxo)
		}

		// Remove spent outputs
		for _, input := range tx.Inputs {
			if input.TxHash != "" {
				bc.UTXOSet.RemoveUTXO(input.TxHash, input.OutIndex)
			}
		}
	}

	// Add to chain
	bc.Blocks = append(bc.Blocks, block)
	bc.BlockTimestamps = append(bc.BlockTimestamps, block.Timestamp)

	// Update total mined
	coinbaseReward := block.Transactions[0].Outputs[0].Value
	bc.RewardCalculator.UpdateTotalMined(coinbaseReward)

	// Store block
	if err := bc.Chain.StoreBlock(block); err != nil {
		return fmt.Errorf("failed to store block: %v", err)
	}

	// Remove transactions from mempool
	for _, tx := range block.Transactions {
		bc.PendingTransactions.RemoveTransaction(tx.TxHash)
	}

	// Adjust difficulty if needed
	if bc.DifficultyAdjuster.ShouldAdjustDifficulty(block.Height) {
		bc.Difficulty = bc.DifficultyAdjuster.AdjustDifficulty(
			bc.Difficulty,
			[]uint64{block.Height},
			bc.BlockTimestamps,
		)
	}

	fmt.Printf("[Blockchain] Block #%d added: %s, Txs: %d, UTXOs: %d\n",
		block.Height, block.BlockHash[:16], len(block.Transactions), bc.UTXOSet.Count())

	return nil
}

// validateBlock performs comprehensive block validation
func (bc *Blockchain) validateBlock(block *Block) error {
	// Check block height
	expectedHeight := uint64(len(bc.Blocks))
	if block.Height != expectedHeight {
		return fmt.Errorf("invalid block height: expected %d, got %d", expectedHeight, block.Height)
	}

	// Check previous block hash
	if len(bc.Blocks) > 0 {
		lastBlock := bc.Blocks[len(bc.Blocks)-1]
		if block.PrevBlockHash != lastBlock.BlockHash {
			return fmt.Errorf("invalid previous block hash")
		}
	} else {
		if !block.IsGenesisBlock() {
			return fmt.Errorf("first block must be genesis block")
		}
	}

	// Check transactions
	if len(block.Transactions) == 0 {
		return fmt.Errorf("block must have at least coinbase transaction")
	}

	// First transaction must be coinbase
	if !block.Transactions[0].IsCoinbase() {
		return fmt.Errorf("first transaction must be coinbase")
	}

	// Validate each transaction
	for i, tx := range block.Transactions {
		if i > 0 && tx.IsCoinbase() {
			return fmt.Errorf("only first transaction can be coinbase")
		}

		if !tx.Validate(bc.UTXOSet) {
			return fmt.Errorf("transaction validation failed: %s", tx.TxHash)
		}
	}

	// Verify Merkle root
	if block.MerkleRoot != block.CalculateMerkleRoot() {
		return fmt.Errorf("invalid merkle root")
	}

	// Verify PoW
	pow := consensus.NewProofOfWork(block.Difficulty)
	if !pow.Validate(block.BlockHash) {
		return fmt.Errorf("invalid proof of work")
	}

	return nil
}

// GetLatestBlock returns the most recent block
func (bc *Blockchain) GetLatestBlock() *Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	if len(bc.Blocks) == 0 {
		return nil
	}
	return bc.Blocks[len(bc.Blocks)-1]
}

// GetBlock retrieves a block by height
func (bc *Blockchain) GetBlock(height uint64) *Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	if height >= uint64(len(bc.Blocks)) {
		return nil
	}
	return bc.Blocks[height]
}

// GetBlockByHash retrieves a block by hash
func (bc *Blockchain) GetBlockByHash(hash string) *Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	for _, block := range bc.Blocks {
		if block.BlockHash == hash {
			return block
		}
	}
	return nil
}

// GetHeight returns the current blockchain height
func (bc *Blockchain) GetHeight() uint64 {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	return uint64(len(bc.Blocks))
}

// AddPendingTransaction adds a transaction to mempool
func (bc *Blockchain) AddPendingTransaction(tx *Transaction, fee uint64) error {
	return bc.PendingTransactions.AddTransaction(tx, fee)
}

// GetPendingTransactions retrieves transactions for mining
func (bc *Blockchain) GetPendingTransactions(limit int) []*Transaction {
	return bc.PendingTransactions.GetTransactions(limit)
}

// GetBalance returns address balance
func (bc *Blockchain) GetBalance(address string) uint64 {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	return bc.UTXOSet.GetBalance(address)
}

// GetUnspentOutputs retrieves UTXOs for an address
func (bc *Blockchain) GetUnspentOutputs(address string) []*UTXO {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	return bc.UTXOSet.FindUTXOsByAddress(address)
}

// rebuildUTXOSet reconstructs UTXO set from blocks
func (bc *Blockchain) rebuildUTXOSet() {
	bc.UTXOSet = NewUTXOSet()

	for _, block := range bc.Blocks {
		for _, tx := range block.Transactions {
			// Add outputs
			for i, output := range tx.Outputs {
				utxo := &UTXO{
					TxHash:    tx.TxHash,
					OutIndex:  uint32(i),
					Value:     output.Value,
					Address:   output.Address,
					LockScript: output.LockScript,
				}
				bc.UTXOSet.AddUTXO(utxo)
			}

			// Remove spent inputs
			for _, input := range tx.Inputs {
				if input.TxHash != "" {
					bc.UTXOSet.RemoveUTXO(input.TxHash, input.OutIndex)
				}
			}
		}
	}
}

// GetChainInfo returns chain statistics
func (bc *Blockchain) GetChainInfo() map[string]interface{} {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	latestBlock := bc.Blocks[len(bc.Blocks)-1]
	nextHalving := bc.RewardCalculator.GetBlocksUntilNextHalving(latestBlock.Height)

	return map[string]interface{}{
		"height":              bc.GetHeight(),
		"blocks":              len(bc.Blocks),
		"difficulty":          bc.Difficulty,
		"last_block_hash":     latestBlock.BlockHash,
		"last_block_time":     latestBlock.Timestamp,
		"total_transactions":  bc.countAllTransactions(),
		"pending_txs":         bc.PendingTransactions.Size(),
		"utxo_count":          bc.UTXOSet.Count(),
		"supply_mined":        bc.RewardCalculator.TotalMinedCoins,
		"mined_percentage":    bc.RewardCalculator.GetMinedPercentage(),
		"blocks_until_halving": nextHalving,
		"current_reward":      bc.RewardCalculator.GetBlockReward(latestBlock.Height),
	}
}

// countAllTransactions counts transactions across all blocks
func (bc *Blockchain) countAllTransactions() int {
	count := 0
	for _, block := range bc.Blocks {
		count += len(block.Transactions)
	}
	return count
}
