package consensus

// BlockRewardCalculator manages block rewards and halving
type BlockRewardCalculator struct {
	InitialReward       uint64
	HalvingInterval     uint64
	MaxSupply           uint64
	TotalMinedCoins     uint64
}

// NewBlockRewardCalculator creates reward calculator
func NewBlockRewardCalculator(initialReward, halvingInterval, maxSupply uint64) *BlockRewardCalculator {
	return &BlockRewardCalculator{
		InitialReward:   initialReward,
		HalvingInterval: halvingInterval,
		MaxSupply:       maxSupply,
		TotalMinedCoins: 0,
	}
}

// GetBlockReward calculates reward for a given block height
func (brc *BlockRewardCalculator) GetBlockReward(blockHeight uint64) uint64 {
	halvings := blockHeight / brc.HalvingInterval
	if halvings >= 64 { // Prevent shift overflow
		return 0
	}

	reward := brc.InitialReward >> halvings

	// Check if total supply would exceed max
	if brc.TotalMinedCoins+reward > brc.MaxSupply {
		return brc.MaxSupply - brc.TotalMinedCoins
	}

	return reward
}

// GetHalvingHeight returns the block height of the next halving
func (brc *BlockRewardCalculator) GetHalvingHeight(blockHeight uint64) uint64 {
	halvings := blockHeight / brc.HalvingInterval
	return (halvings + 1) * brc.HalvingInterval
}

// GetTotalSupply returns total coins that will ever exist
func (brc *BlockRewardCalculator) GetTotalSupply() uint64 {
	return brc.MaxSupply
}

// GetCoinbaseReward calculates total coinbase (block reward + fees)
func (brc *BlockRewardCalculator) GetCoinbaseReward(blockHeight uint64, fees uint64) uint64 {
	blockReward := brc.GetBlockReward(blockHeight)
	return blockReward + fees
}

// UpdateTotalMined records mined coins
func (brc *BlockRewardCalculator) UpdateTotalMined(amount uint64) {
	brc.TotalMinedCoins += amount
	if brc.TotalMinedCoins > brc.MaxSupply {
		brc.TotalMinedCoins = brc.MaxSupply
	}
}

// GetMinedPercentage returns percentage of total supply mined
func (brc *BlockRewardCalculator) GetMinedPercentage() float64 {
	if brc.MaxSupply == 0 {
		return 0
	}
	return float64(brc.TotalMinedCoins) / float64(brc.MaxSupply) * 100
}

// GetBlocksUntilNextHalving returns blocks until next reward halving
func (brc *BlockRewardCalculator) GetBlocksUntilNextHalving(blockHeight uint64) uint64 {
	nextHalving := brc.GetHalvingHeight(blockHeight)
	if nextHalving > blockHeight {
		return nextHalving - blockHeight
	}
	return 0
}
