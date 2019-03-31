package api

import (
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

func (a *API) GetBestBlock(ctx context.Context, in *empty.Empty) (*GetBestBlockResponse, error) {
	resp := &GetBestBlockResponse{
		Height: 2,
		Hash:   "0x123dac",
	}
	return resp, nil
}
