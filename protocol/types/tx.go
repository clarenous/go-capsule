package types

import (
	"bytes"
	"encoding/binary"
	"github.com/clarenous/go-capsule/protocol/types/pb"
	"github.com/golang/protobuf/proto"
)

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
	var b8 [8]byte
	var b40 [40]byte
	binary.LittleEndian.PutUint64(b8[:], vs.Index)
	copy(b40[:], vs.TxID[:])
	copy(b40[32:], b8[:])
	return b40[:]
}

// TxOut
type TxOut struct {
	Value      uint64
	ScriptHash Hash160
}

func (tx *Tx) bytesForID() []byte {
	pb, _ := tx.ToProto()
	var buf bytes.Buffer
	if err := typespb.WriteElement(&buf, pb); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (tx *Tx) Hash() Hash {
	return GetID(tx)
}

func (tx *Tx) WitnessHash() Hash {
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
	pb, _ := tx.ToProto()
	buf, _ := proto.Marshal(pb)
	return uint64(len(buf))
}

func (tx *Tx) ToProto() (*typespb.Tx, error) {
	pb := &typespb.Tx{
		Version:   tx.Version,
		Inputs:    make([]*typespb.Tx_TxIn, len(tx.Inputs)),
		Outputs:   make([]*typespb.Tx_TxOut, len(tx.Outputs)),
		Evidences: make([]*typespb.Tx_Evidence, len(tx.Evidences)),
		LockTime:  tx.LockTime,
	}

	for i, in := range tx.Inputs {
		pb.Inputs[i] = &typespb.Tx_TxIn{
			ValueSource: &typespb.Tx_TxIn_ValueSource{
				Txid:  in.ValueSource.TxID.ToProto(),
				Index: in.ValueSource.Index,
			},
			RedeemScript: append([]byte{}, in.RedeemScript...),
			UnlockScript: append([]byte{}, in.UnlockScript...),
			Sequence:     in.Sequence,
		}
	}

	for i, out := range tx.Outputs {
		pb.Outputs[i] = &typespb.Tx_TxOut{
			Value:      out.Value,
			ScriptHash: out.ScriptHash.ToProto(),
		}
	}

	for i, evid := range tx.Evidences {
		pb.Evidences[i] = &typespb.Tx_Evidence{
			Digest:      append([]byte{}, evid.Digest...),
			Source:      append([]byte{}, evid.Source...),
			ValidScript: append([]byte{}, evid.ValidScript...),
		}
	}

	return pb, nil
}

func (tx *Tx) FromProto(pb *typespb.Tx) error {
	tx.Version = pb.Version
	tx.Inputs = make([]TxIn, len(pb.Inputs))
	tx.Outputs = make([]TxOut, len(pb.Outputs))
	tx.Evidences = make([]Evidence, len(pb.Evidences))
	tx.LockTime = pb.LockTime

	for i, inPb := range pb.Inputs {
		tx.Inputs[i] = TxIn{
			ValueSource: ValueSource{
				TxID:  NewHashFromProto(inPb.ValueSource.Txid).Value(),
				Index: inPb.ValueSource.Index,
			},
			RedeemScript: append([]byte{}, inPb.RedeemScript...),
			UnlockScript: append([]byte{}, inPb.UnlockScript...),
			Sequence:     inPb.Sequence,
		}
	}

	for i, outPb := range pb.Outputs {
		tx.Outputs[i] = TxOut{
			Value:      outPb.Value,
			ScriptHash: NewHash160FromProto(outPb.ScriptHash).Value(),
		}
	}

	for i, evidPb := range pb.Evidences {
		tx.Evidences[i] = Evidence{
			Digest:      append([]byte{}, evidPb.Digest...),
			Source:      append([]byte{}, evidPb.Source...),
			ValidScript: append([]byte{}, evidPb.ValidScript...),
		}
	}

	return nil
}
