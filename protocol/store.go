package protocol

import (
	"github.com/clarenous/go-capsule/database/storage"
	"github.com/clarenous/go-capsule/protocol/types"

	"github.com/clarenous/go-capsule/protocol/state"
)

// Store provides storage interface for blockchain data
type Store interface {
	BlockExist(*types.Hash) bool

	GetBlock(*types.Hash) (*types.Block, error)
	GetStoreStatus() *BlockStoreState
	GetTransactionStatus(*types.Hash) (*types.TransactionStatus, error)
	GetTransactionsUtxo(*state.UtxoViewpoint, []*types.Tx) error
	GetUtxo(*types.Hash) (*storage.UtxoEntry, error)

	LoadBlockIndex(uint64) (*state.BlockIndex, error)
	SaveBlock(*types.Block, *types.TransactionStatus) error
	SaveChainStatus(*state.BlockNode, *state.UtxoViewpoint) error
}

// BlockStoreState represents the core's db status
type BlockStoreState struct {
	Height uint64
	Hash   *types.Hash
}
