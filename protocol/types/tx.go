package types

// Transaction
type Tx struct {
	Version  uint64
	Inputs   []TxIn
	Outputs  []TxOut
	Evidence []Evidence
	LockTime uint64
}

// TxIn
type TxIn struct {
	PreviousOutPoint OutPoint
	RedeemScript     []byte
	UnlockScript     []byte
	Sequence         uint64
}

// OutPoint
type OutPoint struct {
	TxID  Hash
	Index uint64
}

// TxOut
type TxOut struct {
	Value      uint64
	LockScript []byte
}

func (tx *Tx) bytesForID() []byte {
	return []byte{}
}

func (tx *Tx) Hash() Hash {
	return GetID(tx)
}
