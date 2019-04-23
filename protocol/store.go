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
	GetTransactionsUtxo(*state.UtxoViewpoint, []*types.Tx) error
	GetUtxo(*types.Hash) (*storage.UtxoEntry, error)

	LoadBlockIndex(uint64) (*state.BlockIndex, error)
	SaveBlock(*types.Block) error
	SaveChainStatus(*state.BlockNode, *state.UtxoViewpoint) error

	GetTransaction(hash *types.Hash) (*types.Tx, error)
	GetEvidence(hash *types.Hash) (*types.Evidence, error)
}

// BlockStoreState represents the core's db status
type BlockStoreState struct {
	Height uint64
	Hash   *types.Hash
}
