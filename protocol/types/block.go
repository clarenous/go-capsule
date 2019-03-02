package types

// Block
type Block struct {
	BlockHeader  *BlockHeader
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
	Proof           BlockProof
}

// Block Proof
type BlockProof struct {
	Target uint64
	Nonce  uint64
}
