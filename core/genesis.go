package core

import (
	"time"
)

// Genesis configuration for VOIDEX mainnet
const (
	MAINNET_CHAIN_ID           = "voidex-mainnet"
	MAINNET_VERSION            = uint32(1)
	MAINNET_BLOCK_TIME         = 60 // seconds
	MAINNET_INITIAL_SUPPLY     = uint64(50000000 * 100000000) // 50M coins in satoshis
	MAINNET_INITIAL_REWARD     = uint64(50 * 100000000) // 50 coins per block
	MAINNET_REWARD_HALVING     = uint64(210000) // Halve reward every 210k blocks
	MAINNET_MAX_BLOCK_SIZE     = uint32(4000000) // 4MB
	MAINNET_MAX_TX_PER_BLOCK   = uint32(2000)
	MAINNET_DIFFICULTY_WINDOW  = uint32(2016) // Adjust every 2 weeks
	MAINNET_TARGET_BLOCK_TIME  = uint32(60) // seconds
)

// GenesisConfig holds genesis block parameters
type GenesisConfig struct {
	ChainID              string
	Version              uint32
	Timestamp            int64
	InitialDifficulty    uint32
	InitialReward        uint64
	MaxSupply            uint64
	RewardHalvingInterval uint64
	BlockTime            int64
	MaxBlockSize         uint32
	MaxTxPerBlock        uint32
	DifficultyWindow     uint32
	TargetBlockTime      uint32
}

// GetMainnetGenesis returns mainnet genesis configuration
func GetMainnetGenesis() *GenesisConfig {
	return &GenesisConfig{
		ChainID:               MAINNET_CHAIN_ID,
		Version:               MAINNET_VERSION,
		Timestamp:             time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		InitialDifficulty:     0x00000FFF, // Difficulty target
		InitialReward:         MAINNET_INITIAL_REWARD,
		MaxSupply:             MAINNET_INITIAL_SUPPLY,
		RewardHalvingInterval: MAINNET_REWARD_HALVING,
		BlockTime:             MAINNET_BLOCK_TIME,
		MaxBlockSize:          MAINNET_MAX_BLOCK_SIZE,
		MaxTxPerBlock:         MAINNET_MAX_TX_PER_BLOCK,
		DifficultyWindow:      MAINNET_DIFFICULTY_WINDOW,
		TargetBlockTime:       MAINNET_TARGET_BLOCK_TIME,
	}
}

// CreateGenesisBlock creates the genesis block
func CreateGenesisBlock(minerAddress string) *Block {
	genesis := GetMainnetGenesis()

	// Genesis coinbase transaction
	coinbaseTx := &Transaction{
		Version:   1,
		Inputs:    []Input{{TxHash: "", OutIndex: 0}},
		Outputs:   []Output{{Value: genesis.InitialReward, Address: minerAddress, LockScript: "OP_CHECKSIG"}},
		LockTime:  0,
		Timestamp: genesis.Timestamp,
	}
	coinbaseTx.TxHash = coinbaseTx.CalculateHash()

	// Create genesis block
	genesisBlock := &Block{
		Version:       genesis.Version,
		PrevBlockHash: "0",
		Timestamp:     genesis.Timestamp,
		Difficulty:    0x00000FFF,
		Nonce:         0,
		Transactions:  []*Transaction{coinbaseTx},
		Height:        0,
		Miner:         minerAddress,
	}

	genesisBlock.MerkleRoot = genesisBlock.CalculateMerkleRoot()
	genesisBlock.BlockHash = genesisBlock.CalculateHash()

	return genesisBlock
}

// TestnetGenesis returns testnet configuration
func GetTestnetGenesis() *GenesisConfig {
	cfg := GetMainnetGenesis()
	cfg.ChainID = "voidex-testnet"
	cfg.InitialReward = 10 * 100000000 // 10 coins for testing
	return cfg
}
