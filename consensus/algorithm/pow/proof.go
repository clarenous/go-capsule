package pow

import (
	"encoding/binary"
	"errors"
	"github.com/clarenous/go-capsule/consensus"
	ca "github.com/clarenous/go-capsule/consensus/algorithm"
	"github.com/clarenous/go-capsule/protocol/state"
	"github.com/clarenous/go-capsule/protocol/types"
	"math/big"
)

const (
	TypePoW     = "pow"
	ProofLength = 16
)

type WorkProof Proof

var (
	ErrBadBits      = errors.New("block bits is invalid")
	ErrBadWork      = errors.New("invalid difficulty proof of work")
	ErrInvalidBytes = errors.New("invalid bytes to deserialize work proof")
)

func (wp *WorkProof) Bytes() []byte {
	var b16 [ProofLength]byte
	binary.LittleEndian.PutUint64(b16[:8], wp.Target)
	binary.LittleEndian.PutUint64(b16[8:ProofLength], wp.Nonce)
	return b16[:]
}

func (wp *WorkProof) FromBytes(buf []byte) error {
	if len(buf) != ProofLength {
		return ErrInvalidBytes
	}
	wp.Target = binary.LittleEndian.Uint64(buf[:8])
	wp.Nonce = binary.LittleEndian.Uint64(buf[8:])
	return nil
}

func (wp *WorkProof) HintNextProof(args []interface{}) error {
	node, ok := args[0].(*state.BlockNode)
	if !ok {
		return errors.New("wrong type for *BlockNode")
	}
	parentProof, ok := node.Proof.(*WorkProof)
	if !ok {
		return errors.New("wrong type for *WorkProof")
	}

	if node.Height%consensus.BlocksPerRetarget != 0 || node.Height == 0 {
		wp.Target = parentProof.Target
		return nil
	}

	compareNode := node.Parent
	for compareNode.Height%consensus.BlocksPerRetarget != 0 {
		compareNode = compareNode.Parent
	}
	wp.Target = CalcNextRequiredDifficulty(node.BlockHeader(), compareNode.BlockHeader())
	return nil
}

func (wp *WorkProof) ValidateProof(args []interface{}) error {
	b, ok := args[0].(*types.Block)
	if !ok {
		return errors.New("wrong type for *Block")
	}
	parent, ok := args[1].(*state.BlockNode)
	if !ok {
		return errors.New("wrong type for *BlockNode")
	}

	expectedProof, err := parent.HintNextProof()
	if err != nil {
		return nil
	}
	if expectedProof.(*WorkProof).Target != wp.Target {
		return errors.New(" validate proof target not equal")
	}

	if !CheckProofOfWork(b.Hash().Ptr(), b.Proof.(*WorkProof).Nonce) {
		return ErrBadWork
	}
	return nil
}

func (wp *WorkProof) CalcWeight() *big.Int {
	return CalcWork(wp.Target)
}

func NewProof(args ...interface{}) (ca.Proof, error) {
	return &WorkProof{}, nil
}

func init() {
	ca.AddDBBackend(ca.CABackend{
		Typ:      TypePoW,
		NewProof: NewProof,
	})
}
