package core

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"

    "github.com/Mrcubys/VOIDEX-Network---Layer-1-Blockchain-Protocol/storage"
)

type Blockchain struct {
    Blocks     []*Block
    Store      *storage.LevelDBStorage
    Difficulty int
}

func NewBlockchain(store *storage.LevelDBStorage, minerAddress string) (*Blockchain, error) {
    bc := &Blockchain{
        Blocks:     []*Block{},
        Store:      store,
        Difficulty: 1,
    }

    // Tambahkan genesis block
    genesis := &Block{
        Height:    0,
        Timestamp: time.Now().Unix(),
        PrevHash:  "",
        Hash:      "GENESIS",
        Data:      "Genesis Block",
        Miner:     minerAddress,
    }

    bc.Blocks = append(bc.Blocks, genesis)
    if err := store.StoreBlock(genesis); err != nil {
        return nil, err
    }

    return bc, nil
}

func (bc *Blockchain) AddBlock(data string, miner string) (*Block, error) {
    prev := bc.Blocks[len(bc.Blocks)-1]
    b := &Block{
        Height:    prev.Height + 1,
        Timestamp: time.Now().Unix(),
        PrevHash:  prev.Hash,
        Data:      data,
        Miner:     miner,
    }
    b.Hash = calculateHash(b)
    bc.Blocks = append(bc.Blocks, b)
    if err := bc.Store.StoreBlock(b); err != nil {
        return nil, err
    }
    return b, nil
}

func calculateHash(b *Block) string {
    record := fmt.Sprintf("%d%d%s%s%s", b.Height, b.Timestamp, b.PrevHash, b.Data, b.Miner)
    h := sha256.New()
    h.Write([]byte(record))
    return hex.EncodeToString(h.Sum(nil))
}
