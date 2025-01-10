package blockchain

import (
	"testing"
	"time"

	"github.com/fkapsahili/mini-blockchain/internal/types"
)

func TestNewChain(t *testing.T) {
	tempDir := t.TempDir()
	chain, err := NewChain(tempDir)
	if err != nil {
		t.Fatal("Failed to create new chain: %w", err)
	}
	if chain == nil {
		t.Fatal("NewChain() returned nil")
	}
	if chain.GetHeight() != 0 {
		t.Errorf("New chain height = %d, want 0", chain.GetHeight())
	}
}

func TestAddBlock(t *testing.T) {
	tempDir := t.TempDir()
	chain, err := NewChain(tempDir)
	if err != nil {
		t.Fatal("Failed to create new chain: %w", err)
	}

	genesis, err := chain.GetLatestBlock()
	if err != nil {
		t.Fatal("Failed to retrieve latest block: %w", err)
	}

	newBlock := &types.Block{
		Header: types.BlockHeader{
			Version:       1,
			PrevBlockHash: genesis.Hash,
			MerkleRoot:    [32]byte{},
			Timestamp:     time.Now(),
			Difficulty:    1,
			Nonce:         0,
		},
		Height:       genesis.Height + 1,
		Transactions: []types.Transaction{},
	}
	newBlock.UpdateHash()

	err = chain.AddBlock(newBlock)
	if err != nil {
		t.Errorf("AddBlock() error = %v", err)
	}

	if chain.GetHeight() != 1 {
		t.Errorf("Chain height = %d, want 1", chain.GetHeight())
	}

	savedBlock, err := chain.GetBlock(1)
	if err != nil {
		t.Errorf("Failed to retrieve saved block: %v", err)
		return
	}

	if savedBlock.Hash != newBlock.Hash {
		t.Errorf("Saved block hash = %x, want %x", savedBlock.Hash, newBlock.Hash)
	}
}

func TestValidateBlock(t *testing.T) {
	tempDir := t.TempDir()
	chain, err := NewChain(tempDir)
	if err != nil {
		t.Fatal("Failed to create new chain: %w", err)
	}

	tests := []struct {
		name    string
		block   *types.Block
		wantErr bool
	}{
		{
			name:    "nil block",
			block:   nil,
			wantErr: true,
		},
		{
			name: "invalid height",
			block: &types.Block{
				Height: 5,
				Header: types.BlockHeader{
					Version: 1,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := chain.ValidateBlock(tt.block)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
