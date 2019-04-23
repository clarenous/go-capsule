package api

import "github.com/clarenous/go-capsule/errors"

var (
	ErrInvalidBlockID       = errors.New("invalid id for block")
	ErrInvalidTransactionID = errors.New("invalid id for transaction")
	ErrInvalidEvidenceID    = errors.New("invalid id for evidence")
)
