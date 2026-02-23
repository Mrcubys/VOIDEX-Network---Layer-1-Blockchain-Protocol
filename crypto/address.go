package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// Address represents a VOIDEX network address
// Format: VDX + Base58(Hash160(PublicKey) + Checksum)
// Example: VDX1ABC123...
type Address struct {
	hash160 [20]byte // RIPEMD160(SHA256(PublicKey))
}

// NewAddressFromPublicKey creates a VOIDEX address from a public key
func NewAddressFromPublicKey(pubKey *PublicKey) (*Address, error) {
	pubKeyBytes := pubKey.Bytes()
	
	// Step 1: SHA256 hash of public key
	sha256Hash := sha256.Sum256(pubKeyBytes)
	
	// Step 2: RIPEMD160 hash of SHA256 result
	h := ripemd160.New()
	h.Write(sha256Hash[:])
	hash160 := [20]byte{}
	copy(hash160[:], h.Sum(nil))
	
	return &Address{hash160: hash160}, nil
}

// NewAddressFromHash160 creates an address from an existing hash160
func NewAddressFromHash160(hash160 [20]byte) *Address {
	return &Address{hash160: hash160}
}

// String returns the VOIDEX address in standard format
func (a *Address) String() string {
	// Create payload: version byte (0x01 for mainnet) + hash160
	payload := make([]byte, 21)
	payload[0] = 0x01 // VOIDEX mainnet version byte
	copy(payload[1:], a.hash160[:])
	
	// Calculate checksum: first 4 bytes of SHA256(SHA256(payload))
	checksum := calculateChecksum(payload)
	
	// Combine payload + checksum
	addressBytes := append(payload, checksum...)
	
	// Encode to Base58 and prepend "VDX" prefix
	base58Encoded := base58.Encode(addressBytes)
	return "VDX" + base58Encoded
}

// FromString parses a VOIDEX address string back to an Address object
func FromString(addressStr string) (*Address, error) {
	// Remove "VDX" prefix
	if len(addressStr) < 3 || addressStr[:3] != "VDX" {
		return nil, fmt.Errorf("invalid address prefix: expected VDX, got %s", addressStr[:3])
	}
	
	base58Part := addressStr[3:]
	decoded := base58.Decode(base58Part)
	
	// Should be 25 bytes: 1 (version) + 20 (hash160) + 4 (checksum)
	if len(decoded) != 25 {
		return nil, fmt.Errorf("invalid address length: expected 25 bytes, got %d", len(decoded))
	}
	
	// Verify checksum
	payload := decoded[:21]
	storedChecksum := decoded[21:]
	expectedChecksum := calculateChecksum(payload)
	
	if !bytesEqual(storedChecksum, expectedChecksum) {
		return nil, fmt.Errorf("invalid address checksum")
	}
	
	// Verify version byte
	if payload[0] != 0x01 {
		return nil, fmt.Errorf("invalid address version: expected 0x01, got %02x", payload[0])
	}
	
	// Extract hash160
	hash160 := [20]byte{}
	copy(hash160[:], payload[1:])
	
	return &Address{hash160: hash160}, nil
}

// Hash160 returns the 20-byte RIPEMD160 hash
func (a *Address) Hash160() [20]byte {
	return a.hash160
}

// Bytes returns the raw hash160 bytes
func (a *Address) Bytes() []byte {
	return a.hash160[:]
}

// calculateChecksum computes the 4-byte checksum for address validation
func calculateChecksum(payload []byte) []byte {
	hash := sha256.Sum256(payload)
	hash = sha256.Sum256(hash[:])
	return hash[:4]
}

// bytesEqual compares two byte slices for equality (constant time)
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	result := 0
	for i := 0; i < len(a); i++ {
		result |= int(a[i]) ^ int(b[i])
	}
	return result == 0
}
