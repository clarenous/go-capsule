package testutil

import (
	"github.com/clarenous/go-capsule/protocol/types"
)

var (
	MaxHash = &types.Hash{V0: 1<<64 - 1, V1: 1<<64 - 1, V2: 1<<64 - 1, V3: 1<<64 - 1}
	MinHash = &types.Hash{}
)
