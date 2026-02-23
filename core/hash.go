package core

import (
	"crypto/sha256"
	"fmt"
)

// Hash represents a 32-byte SHA256 hash used throughout the blockchain
type Hash [32]byte

// NewHash creates a Hash from a byte slice
func NewHash(data []byte) Hash {
	return sha256.Sum256(data)
}

// String returns the hexadecimal string representation of the hash
func (h Hash) String() string {
	return fmt.Sprintf("%x", h[:])
}

// Bytes returns the hash as a byte slice
func (h Hash) Bytes() []byte {
	return h[:]
}

// IsZero checks if the hash is all zeros (genesis hash)
func (h Hash) IsZero() bool {
	return h == [32]byte{}
}

// DoubleHash computes SHA256(SHA256(data)) - standard for blockchain
func DoubleHash(data []byte) Hash {
	first := sha256.Sum256(data)
	return sha256.Sum256(first[:])
}
