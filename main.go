package main

import (
    "fmt"
    "log"

    "github.com/Mrcubys/VOIDEX-Network---Layer-1-Blockchain-Protocol/core"
    "github.com/Mrcubys/VOIDEX-Network---Layer-1-Blockchain-Protocol/storage"
)

func main() {
    fmt.Println("Starting VOIDEX Network Node...")

    dbPath := "./blockchain_data"
    store, err := storage.NewLevelDBStorage(dbPath)
    if err != nil {
        log.Fatalf("Failed to initialize storage: %v", err)
    }
    defer store.Close()

    minerAddress := "VDX1ABC123..."
    blockchain, err := core.NewBlockchain(store, minerAddress)
    if err != nil {
        log.Fatalf("Failed to create blockchain: %v", err)
    }

    fmt.Printf("Blockchain initialized! Total blocks: %d\n", len(blockchain.Blocks))

    // Tambah block contoh
    _, err = blockchain.AddBlock("Transaksi pertama", minerAddress)
    if err != nil {
        log.Fatalf("Failed to add block: %v", err)
    }

    fmt.Printf("Total blocks after adding one: %d\n", len(blockchain.Blocks))
    fmt.Println("VOIDEX Node is running with persistent storage!")
}
