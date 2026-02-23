package consensus

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

// ProofOfWork implements the PoW consensus mechanism
type ProofOfWork struct {
	Target *big.Int
	MaxNonce uint64
}

// NewProofOfWork creates a new PoW engine
func NewProofOfWork(difficulty uint32) *ProofOfWork {
	// Difficulty target (lower = harder)
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	return &ProofOfWork{
		Target:   target,
		MaxNonce: ^uint64(0),
	}
}

// Mine performs proof-of-work mining
func (pow *ProofOfWork) Mine(blockHeader string, startNonce uint64) (uint64, string, error) {
	var nonce uint64 = startNonce
	var hash string
	var hashInt *big.Int

	startTime := time.Now()
	attempts := uint64(0)

	for nonce < pow.MaxNonce {
		attempts++

		// Every 1 second, check if we should abort
		if attempts%1000000 == 0 && time.Since(startTime) > 0 {
			elapsed := time.Since(startTime).Seconds()
			hashRate := float64(attempts) / elapsed
			fmt.Printf("[Mining] Nonce: %d, Hash/sec: %.0f, Attempts: %d\n", nonce, hashRate, attempts)
		}

		// Hash block header with current nonce
		data := fmt.Sprintf("%s%d", blockHeader, nonce)
		hashByte := sha256.Sum256([]byte(data))
		hash = hex.EncodeToString(hashByte[:])

		hashInt = new(big.Int)
		hashInt.SetString(hash, 16)

		// Check if hash is less than target (valid proof)
		if hashInt.Cmp(pow.Target) < 0 {
			fmt.Printf("[Mining SUCCESS] Nonce: %d, Hash: %s, Time: %.2fs\n", nonce, hash, time.Since(startTime).Seconds())
			return nonce, hash, nil
		}

		nonce++
	}

	return 0, "", fmt.Errorf("max nonce reached without finding valid proof")
}

// Validate checks if a hash meets the difficulty target
func (pow *ProofOfWork) Validate(hash string) bool {
	hashInt := new(big.Int)
	hashInt.SetString(hash, 16)
	return hashInt.Cmp(pow.Target) < 0
}

// GetHashesPerSecond estimates hashes per second (for UI)
func (pow *ProofOfWork) GetDifficulty() float64 {
	// Difficulty = max_target / current_target
	maxTarget := big.NewInt(1)
	maxTarget.Lsh(maxTarget, 256)

	difficulty := new(big.Int)
	difficulty.Div(maxTarget, pow.Target)

	fDiff := new(big.Float).SetInt(difficulty)
	result, _ := fDiff.Float64()
	return result
}
