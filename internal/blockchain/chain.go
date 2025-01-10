package blockchain

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/fkapsahili/mini-blockchain/internal/crypto"
	"github.com/fkapsahili/mini-blockchain/internal/storage"
	"github.com/fkapsahili/mini-blockchain/internal/types"
)

type Chain struct {
	store         storage.ChainStore
	currentHeight uint64
	mu            sync.RWMutex
	latestHash    [32]byte
}

// NewChain creates a new blockchain
func NewChain(dataDir string) (*Chain, error) {
	store, err := storage.NewFileStore(dataDir)
	if err != nil {
		return nil, err
	}

	chain := &Chain{
		store: store,
	}

	if latest, err := store.GetLatestBlock(); err != nil {
		// Create genesis
		genesis := chain.CreateGenesisBlock()
		if err := chain.AddBlock(genesis); err != nil {
			return nil, err
		}
	} else {
		// Restore from storage
		chain.currentHeight = latest.Height
		chain.latestHash = latest.Hash
	}

	return chain, nil
}

// AddBlock adds a new block to the chain
func (c *Chain) AddBlock(block *types.Block) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.ValidateBlock(block); err != nil {
		return fmt.Errorf("block validation failed: %w", err)
	}

	block.UpdateHash()

	// Save
	if err := c.store.SaveBlock(block); err != nil {
		return fmt.Errorf("failed to save block: %w", err)
	}

	// Update the chain state
	c.currentHeight = block.Height
	c.latestHash = block.Hash

	return nil
}

// GetBlock retrieves a block by height
func (c *Chain) GetBlock(height uint64) (*types.Block, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.store.GetBlock(height)
}

// GetLatestBlock returns the most recent block of the chain
func (c *Chain) GetLatestBlock() (*types.Block, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.store.GetBlock(c.currentHeight)
}

// ValidateBlock checks if a new block can be added
func (c *Chain) ValidateBlock(block *types.Block) error {
	if block == nil {
		return errors.New("block cannot be nil")
	}

	if c.currentHeight == 0 && block.Height == 0 {
		// First block must be genesis
		if block.Header.PrevBlockHash != [32]byte{} {
			return errors.New("genesis block must have zero previous hash")
		}
		return nil
	}

	prevBlock, err := c.store.GetBlock(c.currentHeight)
	if err != nil {
		return fmt.Errorf("failed to get previous block: %w", err)
	}

	// Check height
	if block.Height != prevBlock.Height+1 {
		return errors.New("invalid block height")
	}

	// Check previous hash
	if block.Header.PrevBlockHash != prevBlock.Hash {
		return errors.New("invalid previous block hash")

	}

	// Verify Merkle root
	expectedRoot := crypto.CalculateMerkleRoot(getTransactionHashes(block.Transactions))
	if block.Header.MerkleRoot != expectedRoot {
		return errors.New("invalid merkle root")
	}

	// Verify block hash
	expectedHash := crypto.Hash(block.GetHeaderBytes())
	if block.Hash != expectedHash {
		return errors.New("invalid block hash")
	}

	if !crypto.CheckProofOfWork(block.Hash, block.Header.Difficulty) {
		return errors.New("proof of work verification failed")
	}

	for _, tx := range block.Transactions {
		if err := validateTransaction(&tx); err != nil {
			return fmt.Errorf("invalid transaction: %w", err)
		}
	}

	return nil
}

// Helper function to get transaction hashes
func getTransactionHashes(txs []types.Transaction) [][32]byte {
	hashes := make([][32]byte, len(txs))
	for i, tx := range txs {
		hashes[i] = tx.Hash
	}
	return hashes
}

// validateTransaction verifies a single transaction
func validateTransaction(tx *types.Transaction) error {
	// Verify each input
	for _, input := range tx.Inputs {
		pubKey, err := crypto.BytesToPublicKey(input.PublicKey)
		if err != nil {
			return fmt.Errorf("invalid public key in transaction: %w", err)
		}

		// Verify signature
		if !crypto.Verify(pubKey, tx.Hash[:], input.Signature) {
			return errors.New("invalid transaction signature")
		}
	}
	return nil
}

// GetHeight returns the chain's current height
func (c *Chain) GetHeight() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentHeight
}

// GetBlockByHash returns the block by its hash
func (c *Chain) GetBlockByHash(hash [32]byte) (*types.Block, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.store.GetBlockByHash(hash)
}

// CreateGenesisBlock creates the first block
func (c *Chain) CreateGenesisBlock() *types.Block {
	header := types.BlockHeader{
		Version:       1,
		PrevBlockHash: [32]byte{}, // all zeros
		MerkleRoot:    [32]byte{},
		Timestamp:     time.Now(),
		Difficulty:    1,
		Nonce:         0,
	}

	genesisBlock := &types.Block{
		Header:       header,
		Transactions: []types.Transaction{},
		Hash:         [32]byte{},
		Height:       0,
	}

	genesisBlock.UpdateHash()
	return genesisBlock
}

// IsEmpty checks if chain has any blocks
func (c *Chain) IsEmpty() bool {
	return c.currentHeight == 0
}
