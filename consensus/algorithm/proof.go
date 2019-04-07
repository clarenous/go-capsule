package algorithm

import "github.com/golang/protobuf/proto"

type Proof interface {
	FromProto(proto.Message) error
	ToProto() (proto.Message, error)
	Bytes() []byte
	FromBytes([]byte) error
	HintNextProof(args []interface{}) error
	ValidateProof(args []interface{}) error
}

func NewProof(typ string) Proof {
	return nil
}
