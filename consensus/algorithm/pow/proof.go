package pow

import "github.com/golang/protobuf/proto"

type WorkProof Proof

func (wp *WorkProof) FromProto(pb *proto.Message) error {
	return nil
}

func (wp *WorkProof) ToProto() (*proto.Message, error) {
	return nil, nil
}

func (wp *WorkProof) HintNextProof(args []interface{}) error {
	return nil
}

func (wp *WorkProof) ValidateProof(args []interface{}) error {
	return nil
}
