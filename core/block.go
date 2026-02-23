package core

type Block struct {
    Height     uint64
    Timestamp  int64
    PrevHash   string
    Hash       string
    Data       string
    Miner      string
}
