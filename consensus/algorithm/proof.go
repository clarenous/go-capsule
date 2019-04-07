package algorithm

type Proof interface {
	//FromProto(proto.Message) error
	//ToProto() (proto.Message, error)
	Bytes() []byte
	FromBytes([]byte) error
	HintNextProof(args []interface{}) error
	ValidateProof(args []interface{}) error
}

func NewProof(typ string) Proof {
	return &fakeProof{
		buf: []byte(typ),
	}
}

type fakeProof struct {
	buf []byte
}

func (fp *fakeProof) Bytes() []byte {
	return append([]byte{}, fp.buf...)
}

func (fp *fakeProof) FromBytes(bs []byte) error {
	fp.buf = bs
	return nil
}

func (fp *fakeProof) HintNextProof(args []interface{}) error {
	return nil
}

func (fp *fakeProof) ValidateProof(args []interface{}) error {
	return nil
}
