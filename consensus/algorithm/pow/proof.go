package pow

import (
	"errors"
	"github.com/clarenous/go-capsule/consensus"
	"github.com/clarenous/go-capsule/consensus/algorithm/pow/difficulty"
	"github.com/clarenous/go-capsule/protocol/state"
	"github.com/clarenous/go-capsule/protocol/types"
	"github.com/golang/protobuf/proto"
)

type WorkProof Proof

var (
	errBadBits   = errors.New("block bits is invalid")
	errWorkProof = errors.New("invalid difficulty proof of work")
)

func (wp *WorkProof) FromProto(pb *proto.Message) error {
	return nil
}

func (wp *WorkProof) ToProto() (*proto.Message, error) {
	return nil, nil
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
	wp.Target = difficulty.CalcNextRequiredDifficulty(node.BlockHeader(), compareNode.BlockHeader())
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

	if !difficulty.CheckProofOfWork(b.Hash().Ptr(), b.Proof.(*WorkProof).Nonce) {
		return errWorkProof
	}
	return nil
}
