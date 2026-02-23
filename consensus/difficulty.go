package consensus

import (
	"fmt"
)

// DifficultyAdjuster handles difficulty recalculation
type DifficultyAdjuster struct {
	DifficultyWindow   uint32
	TargetBlockTime    uint32 // seconds
	AdjustmentInterval int64  // nanoseconds
}

// NewDifficultyAdjuster creates a new adjuster
func NewDifficultyAdjuster(window, targetTime uint32) *DifficultyAdjuster {
	return &DifficultyAdjuster{
		DifficultyWindow:   window,
		TargetBlockTime:    targetTime,
		AdjustmentInterval: int64(targetTime) * int64(window) * int64(1e9),
	}
}

// AdjustDifficulty recalculates difficulty based on actual block times
func (da *DifficultyAdjuster) AdjustDifficulty(
	previousDifficulty uint32,
	blockHeights []uint64,
	blockTimestamps []int64,
) uint32 {
	// Need enough blocks to measure
	if len(blockHeights) < int(da.DifficultyWindow) {
		return previousDifficulty
	}

	// Get timestamps for difficulty window
	firstTime := blockTimestamps[len(blockTimestamps)-int(da.DifficultyWindow)]
	lastTime := blockTimestamps[len(blockTimestamps)-1]

	actualTime := lastTime - firstTime
	if actualTime < 0 {
		return previousDifficulty
	}

	targetTime := int64(da.TargetBlockTime) * int64(da.DifficultyWindow)

	// Limit adjustment to 4x (prevent wild swings)
	maxAdjustment := int64(4)
	minAdjustment := int64(1)

	adjustment := int64(actualTime) / targetTime
	if adjustment > maxAdjustment {
		adjustment = maxAdjustment
	} else if adjustment < minAdjustment {
		adjustment = minAdjustment
	}

	newDifficulty := uint32(int64(previousDifficulty) / adjustment)

	fmt.Printf("[Difficulty Adjustment] Previous: %d, New: %d, Actual Time: %d, Target Time: %d\n",
		previousDifficulty, newDifficulty, actualTime, targetTime)

	return newDifficulty
}

// ShouldAdjustDifficulty checks if we're at an adjustment point
func (da *DifficultyAdjuster) ShouldAdjustDifficulty(blockHeight uint64) bool {
	return blockHeight > 0 && blockHeight%uint64(da.DifficultyWindow) == 0
}
