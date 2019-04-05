package api

import (
	"github.com/clarenous/go-capsule/protocol/types"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"math/rand"
)

func (a *API) GetBestBlock(ctx context.Context, in *empty.Empty) (*GetBestBlockResponse, error) {
	resp := &GetBestBlockResponse{
		Height: 2,
		Hash:   "0x123dac",
	}
	return resp, nil
}

func (a *API) GetBlock(ctx context.Context, in *GetBlockRequest) (*GetBlockResponse, error) {
	block := types.MockBlock()
	var evidenceCount int
	for _, tx := range block.Transactions {
		evidenceCount += len(tx.Evidences)
	}

	resp := &GetBlockResponse{
		Proof:        &Proof{},
		Transactions: make([]string, len(block.Transactions)),
		Evidences:    make([]string, evidenceCount),
	}

	constructBlockHeaderResp(resp, &block.BlockHeader)

	for i := range resp.Transactions {
		resp.Transactions[i] = types.MockHash().String()
	}

	for i := range resp.Evidences {
		resp.Evidences[i] = types.MockHash().String()
	}

	return resp, nil
}

func (a *API) GetBlockHeader(ctx context.Context, in *GetBlockHeaderRequest) (*GetBlockHeaderResponse, error) {
	block := types.MockBlock()

	resp := &GetBlockHeaderResponse{
		Proof: &Proof{},
	}

	constructBlockHeaderResp(resp, &block.BlockHeader)

	return resp, nil
}

func constructBlockHeaderResp(resp interface{}, header *types.BlockHeader) {
	switch e := resp.(type) {
	case *GetBlockResponse:
		e.Hash = header.Hash().String()
		e.ChainId = header.ChainID.String()
		e.Version = header.Version
		e.Height = header.Height
		e.Timestamp = header.Timestamp
		e.Previous = header.Previous.String()
		e.TransactionRoot = header.TransactionRoot.String()
		e.WitnessRoot = header.WitnessRoot.String()
		e.Proof.Nonce = rand.Uint64()
		e.Proof.Target = rand.Uint64()

	case *GetBlockHeaderResponse:
		e.Hash = header.Hash().String()
		e.ChainId = header.ChainID.String()
		e.Version = header.Version
		e.Height = header.Height
		e.Timestamp = header.Timestamp
		e.Previous = header.Previous.String()
		e.TransactionRoot = header.TransactionRoot.String()
		e.WitnessRoot = header.WitnessRoot.String()
		e.Proof.Nonce = rand.Uint64()
		e.Proof.Target = rand.Uint64()

	default:

	}
}
