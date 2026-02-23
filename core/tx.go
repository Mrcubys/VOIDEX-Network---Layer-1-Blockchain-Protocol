package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Input represents a transaction input (previous output reference)
type Input struct {
	TxHash    string // Reference to previous transaction
	OutIndex  uint32 // Index of output in previous transaction
	Signature []byte // Signature from spender
	PublicKey []byte // Public key of spender
}

// Output represents a transaction output
type Output struct {
	Value      uint64 // Amount in satoshis (smallest unit)
	Address    string // Recipient address
	LockScript string // Locking script (simplified: "OP_CHECKSIG")
}

// Transaction represents a blockchain transaction
type Transaction struct {
	Version  uint32
	Inputs   []Input
	Outputs  []Output
	LockTime int64
	Timestamp int64
	TxHash   string // Calculated hash
}

// NewTransaction creates a new transaction
func NewTransaction(inputs []Input, outputs []Output) *Transaction {
	tx := &Transaction{
		Version:   1,
		Inputs:    inputs,
		Outputs:   outputs,
		LockTime:  0,
		Timestamp: time.Now().Unix(),
	}

	tx.TxHash = tx.CalculateHash()
	return tx
}

// CalculateHash computes the transaction hash
func (t *Transaction) CalculateHash() string {
	// Serialize transaction data
	data := []byte(
		fmt.Sprintf(
			"%d%s%d%d%d",
			t.Version,
			t.serializeInputs(),
			t.serializeOutputs(),
			t.LockTime,
			t.Timestamp,
		),
	)

	hash := sha256.Sum256(data)
	hash = sha256.Sum256(hash[:])
	return hex.EncodeToString(hash[:])
}

// serializeInputs converts inputs to string
func (t *Transaction) serializeInputs() string {
	result := ""
	for _, input := range t.Inputs {
		result += input.TxHash + fmt.Sprintf("%d", input.OutIndex)
	}
	return result
}

// serializeOutputs converts outputs to string
func (t *Transaction) serializeOutputs() string {
	result := ""
	for _, output := range t.Outputs {
		result += output.Address + fmt.Sprintf("%d", output.Value)
	}
	return result
}

// IsCoinbase checks if transaction is a block reward (coinbase)
func (t *Transaction) IsCoinbase() bool {
	return len(t.Inputs) == 1 && t.Inputs[0].TxHash == "" && t.Inputs[0].OutIndex == 0
}

// GetTotalInput calculates total input value
func (t *Transaction) GetTotalInput(utxoSet *UTXOSet) uint64 {
	total := uint64(0)
	for _, input := range t.Inputs {
		utxo := utxoSet.FindUTXO(input.TxHash, input.OutIndex)
		if utxo != nil {
			total += utxo.Value
		}
	}
	return total
}

// GetTotalOutput calculates total output value
func (t *Transaction) GetTotalOutput() uint64 {
	total := uint64(0)
	for _, output := range t.Outputs {
		total += output.Value
	}
	return total
}

// Validate performs basic transaction validation
func (t *Transaction) Validate(utxoSet *UTXOSet) bool {
	if len(t.Inputs) == 0 || len(t.Outputs) == 0 {
		return false
	}

	if !t.IsCoinbase() {
		inputTotal := t.GetTotalInput(utxoSet)
		outputTotal := t.GetTotalOutput()

		if inputTotal < outputTotal {
			return false
		}
	}

	return true
}
