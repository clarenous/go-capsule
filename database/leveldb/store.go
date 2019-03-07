package leveldb

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"

	"github.com/clarenous/go-capsule/database/storage"
	"github.com/clarenous/go-capsule/errors"
	"github.com/clarenous/go-capsule/protocol"

	"github.com/clarenous/go-capsule/protocol/state"
	"github.com/clarenous/go-capsule/protocol/types"
)

const logModule = "leveldb"

var (
	blockStoreKey     = []byte("blockStore")
	blockPrefix       = []byte("B:")
	blockHeaderPrefix = []byte("BH:")
	txStatusPrefix    = []byte("BTS:")
)

func loadBlockStoreStateJSON(db dbm.DB) *protocol.BlockStoreState {
	bytes := db.Get(blockStoreKey)
	if bytes == nil {
		return nil
	}
	bsj := &protocol.BlockStoreState{}
	if err := json.Unmarshal(bytes, bsj); err != nil {
		common.PanicCrisis(common.Fmt("Could not unmarshal bytes: %X", bytes))
	}
	return bsj
}

// A Store encapsulates storage for blockchain validation.
// It satisfies the interface protocol.Store, and provides additional
// methods for querying current data.
type Store struct {
	db    dbm.DB
	cache blockCache
}

func calcBlockKey(hash *types.Hash) []byte {
	return append(blockPrefix, hash.Bytes()...)
}

func calcBlockHeaderKey(height uint64, hash *types.Hash) []byte {
	buf := [8]byte{}
	binary.BigEndian.PutUint64(buf[:], height)
	key := append(blockHeaderPrefix, buf[:]...)
	return append(key, hash.Bytes()...)
}

// GetBlock return the block by given hash
func GetBlock(db dbm.DB, hash *types.Hash) *types.Block {
	bytez := db.Get(calcBlockKey(hash))
	if bytez == nil {
		return nil
	}

	block := &types.Block{}
	block.UnmarshalText(bytez)
	return block
}

// NewStore creates and returns a new Store object.
func NewStore(db dbm.DB) *Store {
	cache := newBlockCache(func(hash *types.Hash) *types.Block {
		return GetBlock(db, hash)
	})
	return &Store{
		db:    db,
		cache: cache,
	}
}

// GetUtxo will search the utxo in db
func (s *Store) GetUtxo(hash *types.Hash) (*storage.UtxoEntry, error) {
	return getUtxo(s.db, hash)
}

// BlockExist check if the block is stored in disk
func (s *Store) BlockExist(hash *types.Hash) bool {
	block, err := s.cache.lookup(hash)
	return err == nil && block != nil
}

// GetBlock return the block by given hash
func (s *Store) GetBlock(hash *types.Hash) (*types.Block, error) {
	return s.cache.lookup(hash)
}

// GetTransactionsUtxo will return all the utxo that related to the input txs
func (s *Store) GetTransactionsUtxo(view *state.UtxoViewpoint, txs []*types.Tx) error {
	return getTransactionsUtxo(s.db, view, txs)
}

// GetStoreStatus return the BlockStoreStateJSON
func (s *Store) GetStoreStatus() *protocol.BlockStoreState {
	return loadBlockStoreStateJSON(s.db)
}

func (s *Store) LoadBlockIndex(stateBestHeight uint64) (*state.BlockIndex, error) {
	startTime := time.Now()
	blockIndex := state.NewBlockIndex()
	bhIter := s.db.IteratorPrefix(blockHeaderPrefix)
	defer bhIter.Release()

	var lastNode *state.BlockNode
	for bhIter.Next() {
		bh := &types.BlockHeader{}
		if err := bh.UnmarshalText(bhIter.Value()); err != nil {
			return nil, err
		}

		// If a block with a height greater than the best height of state is added to the index,
		// It may cause a bug that the new block cant not be process properly.
		if bh.Height > stateBestHeight {
			break
		}

		var parent *state.BlockNode
		if lastNode == nil || lastNode.Hash == bh.PreviousBlockHash {
			parent = lastNode
		} else {
			parent = blockIndex.GetNode(&bh.PreviousBlockHash)
		}

		node, err := state.NewBlockNode(bh, parent)
		if err != nil {
			return nil, err
		}

		blockIndex.AddNode(node)
		lastNode = node
	}

	log.WithFields(log.Fields{
		"module":   logModule,
		"height":   stateBestHeight,
		"duration": time.Since(startTime),
	}).Debug("initialize load history block index from database")
	return blockIndex, nil
}

// SaveBlock persists a new block in the protocol.
func (s *Store) SaveBlock(block *types.Block) error {
	startTime := time.Now()
	binaryBlock, err := block.MarshalText()
	if err != nil {
		return errors.Wrap(err, "Marshal block meta")
	}

	binaryBlockHeader, err := block.BlockHeader.MarshalText()
	if err != nil {
		return errors.Wrap(err, "Marshal block header")
	}

	blockHash := block.Hash()
	batch := s.db.NewBatch()
	batch.Set(calcBlockKey(&blockHash), binaryBlock)
	batch.Set(calcBlockHeaderKey(block.Height, &blockHash), binaryBlockHeader)
	batch.Write()

	log.WithFields(log.Fields{
		"module":   logModule,
		"height":   block.Height,
		"hash":     blockHash.String(),
		"duration": time.Since(startTime),
	}).Info("block saved on disk")
	return nil
}

// SaveChainStatus save the core's newest status && delete old status
func (s *Store) SaveChainStatus(node *state.BlockNode, view *state.UtxoViewpoint) error {
	batch := s.db.NewBatch()
	if err := saveUtxoView(batch, view); err != nil {
		return err
	}

	bytes, err := json.Marshal(protocol.BlockStoreState{Height: node.Height, Hash: &node.Hash})
	if err != nil {
		return err
	}

	batch.Set(blockStoreKey, bytes)
	batch.Write()
	return nil
}
