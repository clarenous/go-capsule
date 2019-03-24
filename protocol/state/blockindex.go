package state

import (
	"errors"
	ca "github.com/clarenous/go-capsule/consensus/algorithm"
	"github.com/clarenous/go-capsule/consensus/algorithm/pow"
	"math/big"
	"sort"
	"sync"

	"github.com/clarenous/go-capsule/common"
	"github.com/clarenous/go-capsule/consensus"
	"github.com/clarenous/go-capsule/consensus/algorithm/pow/difficulty"
	"github.com/clarenous/go-capsule/protocol/types"
)

// approxNodesPerDay is an approximation of the number of new blocks there are
// in a day on average.
const approxNodesPerDay = 24 * 24

// BlockNode represents a block within the block chain and is primarily used to
// aid in selecting the best chain to be the main chain.
type BlockNode struct {
	Parent  *BlockNode // parent is the parent block for this node.
	Hash    types.Hash // hash of the block.
	WorkSum *big.Int   // total amount of work in the chain up to

	Version         uint64
	Height          uint64
	Timestamp       uint64
	Proof           ca.Proof
	TransactionRoot types.Hash
	WitnessRoot     types.Hash
}

func NewBlockNode(bh *types.BlockHeader, parent *BlockNode) (*BlockNode, error) {
	if bh.Height != 0 && parent == nil {
		return nil, errors.New("parent node can not be nil")
	}

	node := &BlockNode{
		Parent:          parent,
		Hash:            bh.Hash(),
		WorkSum:         difficulty.CalcWork(bh.Proof.(*pow.WorkProof).Target),
		Version:         bh.Version,
		Height:          bh.Height,
		Timestamp:       bh.Timestamp,
		Proof:           bh.Proof,
		TransactionRoot: bh.TransactionRoot,
		WitnessRoot:     bh.WitnessRoot,
	}

	if bh.Height != 0 {
		node.WorkSum = node.WorkSum.Add(parent.WorkSum, node.WorkSum)
	}
	return node, nil
}

// blockHeader convert a node to the header struct
func (node *BlockNode) BlockHeader() *types.BlockHeader {
	previousBlockHash := types.Hash{}
	if node.Parent != nil {
		previousBlockHash = node.Parent.Hash
	}
	return &types.BlockHeader{
		Version:         node.Version,
		Height:          node.Height,
		Timestamp:       node.Timestamp,
		Previous:        previousBlockHash,
		TransactionRoot: node.TransactionRoot,
		Proof:           node.Proof,
	}
}

func (node *BlockNode) CalcPastMedianTime() uint64 {
	timestamps := []uint64{}
	iterNode := node
	for i := 0; i < consensus.MedianTimeBlocks && iterNode != nil; i++ {
		timestamps = append(timestamps, iterNode.Timestamp)
		iterNode = iterNode.Parent
	}

	sort.Sort(common.TimeSorter(timestamps))
	return timestamps[len(timestamps)/2]
}

// HintNextProof calculate the bits for next block
func (node *BlockNode) HintNextProof() (ca.Proof, error) {
	proof := new(pow.WorkProof)
	err := proof.HintNextProof([]interface{}{node})
	return proof, err
}

// BlockIndex is the struct for help chain trace block chain as tree
type BlockIndex struct {
	sync.RWMutex

	index     map[types.Hash]*BlockNode
	mainChain []*BlockNode
}

// NewBlockIndex will create a empty BlockIndex
func NewBlockIndex() *BlockIndex {
	return &BlockIndex{
		index:     make(map[types.Hash]*BlockNode),
		mainChain: make([]*BlockNode, 0, approxNodesPerDay),
	}
}

// AddNode will add node to the index map
func (bi *BlockIndex) AddNode(node *BlockNode) {
	bi.Lock()
	bi.index[node.Hash] = node
	bi.Unlock()
}

// GetNode will search node from the index map
func (bi *BlockIndex) GetNode(hash *types.Hash) *BlockNode {
	bi.RLock()
	defer bi.RUnlock()
	return bi.index[*hash]
}

func (bi *BlockIndex) BestNode() *BlockNode {
	bi.RLock()
	defer bi.RUnlock()
	return bi.mainChain[len(bi.mainChain)-1]
}

// BlockExist check does the block existed in blockIndex
func (bi *BlockIndex) BlockExist(hash *types.Hash) bool {
	bi.RLock()
	_, ok := bi.index[*hash]
	bi.RUnlock()
	return ok
}

// TODO: THIS FUNCTION MIGHT BE DELETED
func (bi *BlockIndex) InMainchain(hash types.Hash) bool {
	bi.RLock()
	defer bi.RUnlock()

	node, ok := bi.index[hash]
	if !ok {
		return false
	}
	return bi.nodeByHeight(node.Height) == node
}

func (bi *BlockIndex) nodeByHeight(height uint64) *BlockNode {
	if height >= uint64(len(bi.mainChain)) {
		return nil
	}
	return bi.mainChain[height]
}

// NodeByHeight returns the block node at the specified height.
func (bi *BlockIndex) NodeByHeight(height uint64) *BlockNode {
	bi.RLock()
	defer bi.RUnlock()
	return bi.nodeByHeight(height)
}

// SetMainChain will set the the mainChain array
func (bi *BlockIndex) SetMainChain(node *BlockNode) {
	bi.Lock()
	defer bi.Unlock()

	needed := node.Height + 1
	if uint64(cap(bi.mainChain)) < needed {
		nodes := make([]*BlockNode, needed, needed+approxNodesPerDay)
		copy(nodes, bi.mainChain)
		bi.mainChain = nodes
	} else {
		i := uint64(len(bi.mainChain))
		bi.mainChain = bi.mainChain[0:needed]
		for ; i < needed; i++ {
			bi.mainChain[i] = nil
		}
	}

	for node != nil && bi.mainChain[node.Height] != node {
		bi.mainChain[node.Height] = node
		node = node.Parent
	}
}
