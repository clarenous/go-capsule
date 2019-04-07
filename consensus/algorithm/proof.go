package algorithm

import (
	"errors"
	"math/big"
)

type Proof interface {
	Bytes() []byte
	FromBytes([]byte) error
	HintNextProof(args []interface{}) error
	ValidateProof(args []interface{}) error
	CalcWeight() *big.Int
}

var (
	ErrInvalidCAType = errors.New("invalid consensus algorithm type")
	ErrInvalidCAArgs = errors.New("invalid consensus algorithm args")
)

var (
	CABackendList []CABackend
)

type CABackend struct {
	Typ      string
	NewProof func(args ...interface{}) (Proof, error)
}

func AddDBBackend(ins CABackend) {
	for _, dbb := range CABackendList {
		if dbb.Typ == ins.Typ {
			return
		}
	}
	CABackendList = append(CABackendList, ins)
}

func NewProof(dbType string, args ...interface{}) (Proof, error) {
	for _, cab := range CABackendList {
		if cab.Typ == dbType {
			return cab.NewProof(args...)
		}
	}
	return nil, ErrInvalidCAType
}
