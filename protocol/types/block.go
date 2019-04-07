package types

import (
	"bytes"
	"encoding/binary"
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

	proof, err := ca.NewProof("pow")
	if err != nil {
		return err
	}
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

func (blk *Block) MarshalTextForStore() ([]byte, []*TxLoc, error) {
	var blockHash = blk.Hash()
	var buf bytes.Buffer
	var b8 [8]byte
	var pos int
	var writeLength = func(length int) {
		binary.LittleEndian.PutUint64(b8[:], uint64(length))
		buf.Write(b8[:])
	}

	pb, err := blk.ToProto()
	if err != nil {
		return nil, nil, err
	}

	// Write blockHeader
	bhBytes, err := proto.Marshal(pb.BlockHeader)
	if err != nil {
		return nil, nil, err
	}
	writeLength(len(bhBytes))
	buf.Write(bhBytes)
	pos += 8 + len(bhBytes)

	// Write tx count
	writeLength(len(blk.Transactions))
	pos += 8

	// Write txs
	txLocs := make([]*TxLoc, len(pb.Transactions))
	for i, tx := range pb.Transactions {
		txbytes, err := proto.Marshal(tx)
		txLen := len(txbytes)
		if err != nil {
			return nil, nil, err
		}
		txLocs[i] = &TxLoc{
			TxHash:    blk.Transactions[i].Hash(),
			BlockHash: blockHash,
			Offset:    uint64(pos + 8),
			Length:    uint64(txLen),
		}
		writeLength(txLen)
		buf.Write(txbytes)
		pos += 8 + txLen
	}

	return buf.Bytes(), txLocs, nil
}

func (blk *Block) UnmarshalTextForStore(buf []byte) error {
	var r = bytes.NewReader(buf)
	var b8 [8]byte
	var readLength = func() int {
		r.Read(b8[:])
		return int(binary.LittleEndian.Uint64(b8[:]))
	}

	// Read blockHeader
	bhPb, header, bhBytes := new(typespb.BlockHeader), new(BlockHeader), make([]byte, readLength())
	if err := proto.Unmarshal(bhBytes, bhPb); err != nil {
		return err
	}
	if err := header.FromProto(bhPb); err != nil {
		return err
	}
	blk.BlockHeader = *header

	// Read tx count
	var txCount = readLength()

	// Read txs
	blk.Transactions = make([]*Tx, txCount)
	for i := 0; i < txCount; i++ {
		txPb, tx, txBytes := new(typespb.Tx), new(Tx), make([]byte, readLength())
		if err := proto.Unmarshal(txBytes, txPb); err != nil {
			return err
		}
		if err := tx.FromProto(txPb); err != nil {
			return err
		}
		blk.Transactions[i] = tx
	}

	return nil
}

type TxLoc struct {
	TxHash    Hash
	BlockHash Hash
	Offset    uint64
	Length    uint64
}

func (tl *TxLoc) Byte80() [80]byte {
	var b80 [80]byte
	var b8 [8]byte

	copy(b80[:], tl.TxHash[:])
	copy(b80[32:], tl.BlockHash[:])

	binary.LittleEndian.PutUint64(b8[:], tl.Offset)
	copy(b80[64:], b8[:])
	binary.LittleEndian.PutUint64(b8[:], tl.Length)
	copy(b80[72:], b8[:])

	return b80
}

func NewTxLocFromBytes(buf []byte) *TxLoc {
	var b80 [80]byte
	var txHash, blockHash Hash

	copy(b80[:], buf[:])
	copy(txHash[:], b80[:32])
	copy(blockHash[:], b80[32:64])

	loc := &TxLoc{
		TxHash:    txHash,
		BlockHash: blockHash,
		Offset:    binary.LittleEndian.Uint64(b80[64:72]),
		Length:    binary.LittleEndian.Uint64(b80[72:80]),
	}
	return loc
}
