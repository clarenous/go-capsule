package types

import (
	"bytes"
	ca "github.com/clarenous/go-capsule/consensus/algorithm"
	"github.com/clarenous/go-capsule/protocol/types/pb"
	"github.com/golang/protobuf/proto"
)

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
	pb, _ := bh.ToProto()
	var buf bytes.Buffer
	if err := typespb.WriteElement(&buf, pb); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (bh *BlockHeader) Hash() Hash {
	return GetID(bh)
}

func (bh *BlockHeader) ToProto() (*typespb.BlockHeader, error) {
	pb := new(typespb.BlockHeader)
	pb.ChainId = bh.ChainID.ToProto()
	pb.Version = bh.Version
	pb.Height = bh.Height
	pb.Timestamp = bh.Timestamp
	pb.Previous = bh.Previous.ToProto()
	pb.TransactionRoot = bh.TransactionRoot.ToProto()
	pb.WitnessRoot = bh.WitnessRoot.ToProto()
	pb.Proof = bh.Proof.Bytes()

	return pb, nil
}

func (bh *BlockHeader) FromProto(pb *typespb.BlockHeader) error {
	bh.ChainID = NewHashFromProto(pb.ChainId).Value()
	bh.Version = pb.Version
	bh.Height = pb.Height
	bh.Timestamp = pb.Timestamp
	bh.Previous = NewHashFromProto(pb.Previous).Value()
	bh.TransactionRoot = NewHashFromProto(pb.TransactionRoot).Value()
	bh.WitnessRoot = NewHashFromProto(pb.WitnessRoot).Value()

	proof := ca.NewProof("pow")
	if err := proof.FromBytes(pb.Proof); err != nil {
		return err
	}
	bh.Proof = proof

	return nil
}

func (blk *Block) ToProto() (*typespb.Block, error) {
	pb := new(typespb.Block)
	bhPb, _ := blk.BlockHeader.ToProto()
	pb.BlockHeader = bhPb

	pb.Transactions = make([]*typespb.Tx, len(blk.Transactions))
	for i, tx := range blk.Transactions {
		txPb, err := tx.ToProto()
		if err != nil {
			return nil, err
		}
		pb.Transactions[i] = txPb
	}

	return pb, nil
}

func (blk *Block) FromProto(pb *typespb.Block) error {
	bh := new(BlockHeader)
	if err := bh.FromProto(pb.BlockHeader); err != nil {
		return err
	}
	blk.BlockHeader = *bh

	blk.Transactions = make([]*Tx, len(pb.Transactions))
	for i, txPb := range pb.Transactions {
		tx := new(Tx)
		if err := tx.FromProto(txPb); err != nil {
			return err
		}
		blk.Transactions[i] = tx
	}

	return nil
}

func (bh *BlockHeader) MarshalText() ([]byte, error) {
	pb, _ := bh.ToProto()
	return proto.Marshal(pb)
}

func (bh *BlockHeader) UnmarshalText(buf []byte) error {
	pb := new(typespb.BlockHeader)
	if err := proto.Unmarshal(buf, pb); err != nil {
		return err
	}
	return bh.FromProto(pb)
}

func (blk *Block) MarshalText() ([]byte, error) {
	pb, _ := blk.ToProto()
	return proto.Marshal(pb)
}

func (blk *Block) UnmarshalText(buf []byte) error {
	pb := new(typespb.Block)
	if err := proto.Unmarshal(buf, pb); err != nil {
		return err
	}
	return blk.FromProto(pb)
}
