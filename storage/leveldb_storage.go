package storage

import (
    "encoding/json"
    "fmt"
    "github.com/Mrcubys/VOIDEX-Network---Layer-1-Blockchain-Protocol/core"
    "github.com/syndtr/goleveldb/leveldb"
)

type LevelDBStorage struct {
    db *leveldb.DB
}

func NewLevelDBStorage(path string) (*LevelDBStorage, error) {
    db, err := leveldb.OpenFile(path, nil)
    if err != nil {
        return nil, err
    }
    return &LevelDBStorage{db: db}, nil
}

func (ldb *LevelDBStorage) StoreBlock(block *core.Block) error {
    data, err := json.Marshal(block)
    if err != nil {
        return err
    }
    key := fmt.Sprintf("block-%d", block.Height)
    return ldb.db.Put([]byte(key), data, nil)
}

func (ldb *LevelDBStorage) GetBlock(height uint64) (*core.Block, error) {
    key := fmt.Sprintf("block-%d", height)
    data, err := ldb.db.Get([]byte(key), nil)
    if err != nil {
        return nil, err
    }
    var block core.Block
    if err := json.Unmarshal(data, &block); err != nil {
        return nil, err
    }
    return &block, nil
}

func (ldb *LevelDBStorage) Close() error {
    return ldb.db.Close()
}
