package storage

import "github.com/fkapsahili/mini-blockchain/internal/types"

type ChainStore interface {
	SaveBlock(block *types.Block) error
	GetBlock(height uint64) (*types.Block, error)
	GetBlockByHash(hash [32]byte) (*types.Block, error)
	GetLatestBlock() (*types.Block, error)
}
