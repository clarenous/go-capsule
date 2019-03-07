package types

import ca "github.com/clarenous/go-capsule/consensus/algorithm"

// Block
type Block struct {
	BlockHeader
	Transactions []*Tx
}

// Block Header
type BlockHeader struct {
	ChainID         Hash
	Version         uint64
	Height          uint64
	Timestamp       uint64
	Previous        Hash
	TransactionRoot Hash
	WitnessRoot     Hash
	Proof           ca.Proof
}

func (bh *BlockHeader) bytesForID() []byte {
	return []byte{}
}

func (bh *BlockHeader) Hash() Hash {
	return GetID(bh)
}
