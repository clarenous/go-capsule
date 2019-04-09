package api

import (
	"github.com/clarenous/go-capsule/protocol/types"
	"regexp"
	"strconv"
)

const (
	blockIDHeight = iota
	blockIDHash
)

var (
	RegexpBlockHeightPattern = `^height-\d+$`
	RegexpBlockHashPattern   = `^hash-[a-fA-F0-9]{64}$`
	RegexpBlockHeight        *regexp.Regexp
	RegexpBlockHash          *regexp.Regexp
)

func (a *API) getBlockByID(id string) (*types.Block, error) {
	if height, err := decodeBlockID(id, blockIDHeight); err == nil {
		return a.Chain.GetBlockByHeight(height.(uint64))
	}

	if hash, err := decodeBlockID(id, blockIDHash); err == nil {
		return a.Chain.GetBlockByHash(hash.(types.Hash).Ptr())
	}

	return nil, ErrInvalidBlockID
}

func decodeBlockID(id string, typ int) (interface{}, error) {
	switch typ {
	case blockIDHeight:
		if !RegexpBlockHeight.MatchString(id) {
			return nil, ErrInvalidBlockID
		}
		height, _ := strconv.Atoi(id[7:]) // already make sure
		return uint64(height), nil

	case blockIDHash:
		if !RegexpBlockHash.MatchString(id) {
			return nil, ErrInvalidBlockID
		}
		return types.NewHashFromString(id[5:])

	default:
		return nil, ErrInvalidBlockID
	}
}

func init() {
	var err error

	RegexpBlockHeight, err = regexp.Compile(RegexpBlockHeightPattern)
	if err != nil {
		panic(err)
	}

	RegexpBlockHash, err = regexp.Compile(RegexpBlockHashPattern)
	if err != nil {
		panic(err)
	}
}
