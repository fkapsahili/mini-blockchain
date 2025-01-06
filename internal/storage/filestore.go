package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fkapsahili/mini-blockchain/internal/types"
)

type FileStore struct {
	dataDir     string
	latestBlock *types.Block
	mu          sync.RWMutex
}

// NewFileStore creates a new file-based store
func NewFileStore(dataDir string) (*FileStore, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	store := &FileStore{
		dataDir: dataDir,
	}

	latest, err := store.GetLatestBlock()
	if err == nil {
		store.latestBlock = latest
	}

	return store, nil
}

// blockPath returns the file path for a block at given height
func (f *FileStore) blockPath(height uint64) string {
	return filepath.Join(f.dataDir, fmt.Sprintf("block_%d.json", height))
}

func (f *FileStore) SaveBlock(block *types.Block) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	data, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("failed to marshal block: %w", err)
	}

	path := f.blockPath(block.Height)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write block file: %w", err)
	}

	f.latestBlock = block
	return nil
}

func (f *FileStore) GetBlock(height uint64) (*types.Block, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Read block file
	data, err := os.ReadFile(f.blockPath(height))
	if err != nil {
		return nil, fmt.Errorf("failed to read block file: %w", err)
	}

	// Deserialize
	var block types.Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block: %w", err)
	}

	return &block, nil
}

func (f *FileStore) GetBlockByHash(hash [32]byte) (*types.Block, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// For production, we'd want to maintain a hash -> height index
	files, err := os.ReadDir(f.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(f.dataDir, file.Name()))
		if err != nil {
			continue
		}

		var block types.Block
		if err := json.Unmarshal(data, &block); err != nil {
			continue
		}

		if block.Hash == hash {
			return &block, nil
		}
	}

	return nil, fmt.Errorf("block not found with hash %x", hash)
}

func (f *FileStore) GetLatestBlock() (*types.Block, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.latestBlock != nil {
		return f.latestBlock, nil
	}

	// Find the highest block number
	files, err := os.ReadDir(f.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	var highestHeight uint64
	for _, file := range files {
		var height uint64
		_, err := fmt.Sscanf(file.Name(), "block_%d.json", &height)
		if err != nil {
			continue
		}
		if height > highestHeight {
			highestHeight = height
		}
	}

	if highestHeight == 0 {
		return nil, fmt.Errorf("no blocks found")
	}

	return f.GetBlock(highestHeight)
}
