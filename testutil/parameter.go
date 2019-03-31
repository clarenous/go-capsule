package testutil

import (
	"github.com/clarenous/go-capsule/protocol/types"
)

var (
	MaxHash = &types.Hash{S0: 1<<64 - 1, S1: 1<<64 - 1, S2: 1<<64 - 1, S3: 1<<64 - 1}
	MinHash = &types.Hash{}
)
