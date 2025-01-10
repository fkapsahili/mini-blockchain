package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"time"

	"github.com/fkapsahili/mini-blockchain/internal/crypto"
)

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

// getHeaderBytes serializes block header fields in bytes for hashing
func (b *Block) GetHeaderBytes() []byte {
	buf := new(bytes.Buffer)

	// Write all header fields in order
	binary.Write(buf, binary.LittleEndian, b.Header.Version)
	buf.Write(b.Header.PrevBlockHash[:])
	buf.Write(b.Header.MerkleRoot[:])
	binary.Write(buf, binary.LittleEndian, b.Header.Timestamp.Unix())
	binary.Write(buf, binary.LittleEndian, b.Header.Difficulty)
	binary.Write(buf, binary.LittleEndian, b.Header.Nonce)

	return buf.Bytes()
}

// ComputeHash calculates the hash of the block
func (b *Block) ComputeHash() [32]byte {
	headerBytes := b.GetHeaderBytes()
	return crypto.Hash(headerBytes)
}

// ComputeMerkleRoot calculates the Merkle root of transactions
func (b *Block) ComputeMerkleRoot() [32]byte {
	if len(b.Transactions) == 0 {
		// Empty byte slice for empty block
		return sha256.Sum256([]byte{})
	}

	hashes := make([][32]byte, len(b.Transactions))
	return crypto.CalculateMerkleRoot(hashes)
}

// UpdateHash updates both Merkle root and block hash
func (b *Block) UpdateHash() {
	b.Header.MerkleRoot = b.ComputeMerkleRoot()
	b.Hash = b.ComputeHash()
}
