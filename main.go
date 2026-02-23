package main

import (
	"fmt"
	"log"

	"github.com/Mrcubys/VOIDEX-Network---Layer-1-Blockchain-Protocol/core"
	"github.com/Mrcubys/VOIDEX-Network---Layer-1-Blockchain-Protocol/storage"
)

func main() {
	fmt.Println("Starting VOIDEX Network Layer-1 Blockchain Protocol...")

	// Initialize storage
	dbPath := "./blockchain_data"
	store, err := storage.NewLevelDBStorage(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.Close()

	// Create blockchain instance
	minerAddress := "VDX1ABC123..." // Replace with your address
	blockchain, err := core.NewBlockchain(store, minerAddress)
	if err != nil {
		log.Fatalf("Failed to create blockchain: %v", err)
	}

	fmt.Printf("Blockchain initialized successfully!\n")
	fmt.Printf("Total blocks: %d\n", len(blockchain.Blocks))
	fmt.Printf("Current difficulty: 0x%X\n", blockchain.Difficulty)

	// Your blockchain logic here...
	fmt.Println("VOIDEX Network is running...")
}
