package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

// PrivateKey wraps an ECDSA private key
type PrivateKey struct {
	key *ecdsa.PrivateKey
}

// PublicKey wraps an ECDSA public key
type PublicKey struct {
	key *ecdsa.PublicKey
}

// NewPrivateKey generates a new random ECDSA private key using secp256k1
func NewPrivateKey() (*PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	return &PrivateKey{key: privateKey}, nil
}

// PublicKey returns the public key associated with this private key
func (pk *PrivateKey) PublicKey() *PublicKey {
	return &PublicKey{key: &pk.key.PublicKey}
}

// Bytes returns the private key as bytes (32 bytes for P256)
func (pk *PrivateKey) Bytes() []byte {
	return pk.key.D.Bytes()
}

// Bytes returns the public key in uncompressed format (65 bytes)
func (pub *PublicKey) Bytes() []byte {
	// X and Y coordinates concatenated (each 32 bytes) with 0x04 prefix
	x := pub.key.X.Bytes()
	y := pub.key.Y.Bytes()
	
	// Pad to 32 bytes if needed
	xBytes := make([]byte, 32)
	yBytes := make([]byte, 32)
	copy(xBytes[32-len(x):], x)
	copy(yBytes[32-len(y):], y)
	
	result := make([]byte, 65)
	result[0] = 0x04 // uncompressed format prefix
	copy(result[1:33], xBytes)
	copy(result[33:65], yBytes)
	
	return result
}

// String returns the public key as hex string
func (pub *PublicKey) String() string {
	return fmt.Sprintf("%x", pub.Bytes())
}

// GetRawKey returns the underlying ECDSA key for signing/verification
func (pk *PrivateKey) GetRawKey() *ecdsa.PrivateKey {
	return pk.key
}

// GetRawKey returns the underlying ECDSA public key
func (pub *PublicKey) GetRawKey() *ecdsa.PublicKey {
	return pub.key
}
