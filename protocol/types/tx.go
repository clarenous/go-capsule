package types

// Transaction
type Tx struct {
	Version   uint64
	Inputs    []TxIn
	Outputs   []TxOut
	Evidences []Evidence
	LockTime  uint64
}

// TxIn
type TxIn struct {
	ValueSource  ValueSource
	RedeemScript []byte
	UnlockScript []byte
	Sequence     uint64
}

// ValueSource
type ValueSource struct {
	TxID  Hash
	Index uint64
}

func (vs *ValueSource) Hash() Hash {
	return GetID(vs)
}

func (vs *ValueSource) bytesForID() []byte {
	return []byte{}
}

// TxOut
type TxOut struct {
	Value      uint64
	ScriptHash Hash160
}

func (tx *Tx) bytesForID() []byte {
	return []byte{}
}

func (tx *Tx) Hash() Hash {
	return GetID(tx)
}

func (tx *Tx) OutHash(index int) Hash {
	if index >= len(tx.Outputs) {
		panic("out of index for tx OutHash")
	}
	vs := ValueSource{
		TxID:  tx.Hash(),
		Index: uint64(index),
	}
	return vs.Hash()
}

func (tx *Tx) SerializedSize() uint64 {
	// TODO: Calc Size
	return 0
}
