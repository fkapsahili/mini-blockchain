package types

import "time"

// BlockHeaders contains metadata about a block
type BlockHeader struct {
	Version       int32
	PrevBlockHash [32]byte
	MerkleRoot    [32]byte
	Timestamp     time.Time
	Difficulty    uint32
	Nonce         uint32
}

// Block represents a complete block in the chain
type Block struct {
	Header       BlockHeader
	Transactions []Transaction
	Hash         [32]byte
	Height       uint64
}
