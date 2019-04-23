package api

import (
	"github.com/clarenous/go-capsule/consensus/algorithm/pow"
	"github.com/clarenous/go-capsule/protocol/types"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

func (a *API) GetBestBlock(ctx context.Context, in *empty.Empty) (*GetBestBlockResponse, error) {
	bestHeader := a.Chain.BestBlockHeader()
	resp := &GetBestBlockResponse{
		Height: bestHeader.Height,
		Hash:   bestHeader.Hash().String(),
	}
	return resp, nil
}

func (a *API) GetBlock(ctx context.Context, in *GetBlockRequest) (*GetBlockResponse, error) {
	block, err := a.getBlockByID(in.Id)
	if err != nil {
		return nil, err
	}

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

	var evidIndex int
	for i, tx := range block.Transactions {
		txid := tx.Hash()
		resp.Transactions[i] = txid.String()
		for j, evid := range tx.Evidences {
			resp.Evidences[evidIndex] = evid.Hash(txid, uint64(j)).String()
			evidIndex++
		}
	}

	return resp, nil
}

func (a *API) GetBlockHeader(ctx context.Context, in *GetBlockHeaderRequest) (*GetBlockHeaderResponse, error) {
	block, err := a.getBlockByID(in.Id)
	if err != nil {
		return nil, err
	}

	resp := &GetBlockHeaderResponse{
		Proof: &Proof{},
	}

	constructBlockHeaderResp(resp, &block.BlockHeader)

	return resp, nil
}

func (a *API) GetBlockVerboseV0(ctx context.Context, in *GetBlockVerboseRequest) (*GetBlockVerboseV0Response, error) {
	block, err := a.getBlockByID(in.Id)
	if err != nil {
		return nil, err
	}

	resp := &GetBlockVerboseV0Response{
		Proof:        &Proof{},
		Transactions: make([]*GetBlockVerboseV0Response_Transaction, len(block.Transactions)),
	}

	constructBlockHeaderResp(resp, &block.BlockHeader)
	for i, tx := range block.Transactions {
		txid := tx.Hash()
		resp.Transactions[i] = &GetBlockVerboseV0Response_Transaction{
			Txid:      txid.String(),
			Evidences: make([]string, len(tx.Evidences)),
		}
		for j, evid := range tx.Evidences {
			resp.Transactions[i].Evidences[j] = evid.Hash(txid, uint64(j)).String()
		}
	}

	return resp, nil
}

func (a *API) GetBlockVerboseV1(ctx context.Context, in *GetBlockVerboseRequest) (*GetBlockVerboseV1Response, error) {
	block, err := a.getBlockByID(in.Id)
	if err != nil {
		return nil, err
	}

	resp := &GetBlockVerboseV1Response{
		Proof:        &Proof{},
		Transactions: make([]*Tx, len(block.Transactions)),
	}

	constructBlockHeaderResp(resp, &block.BlockHeader)
	for i, tx := range block.Transactions {
		respTx := new(Tx)
		constructTxResp(respTx, tx)
		resp.Transactions[i] = respTx
	}

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
		e.Proof.Nonce = header.Proof.(*pow.WorkProof).Nonce
		e.Proof.Target = header.Proof.(*pow.WorkProof).Target

	case *GetBlockHeaderResponse:
		e.Hash = header.Hash().String()
		e.ChainId = header.ChainID.String()
		e.Version = header.Version
		e.Height = header.Height
		e.Timestamp = header.Timestamp
		e.Previous = header.Previous.String()
		e.TransactionRoot = header.TransactionRoot.String()
		e.WitnessRoot = header.WitnessRoot.String()
		e.Proof.Nonce = header.Proof.(*pow.WorkProof).Nonce
		e.Proof.Target = header.Proof.(*pow.WorkProof).Target

	case *GetBlockVerboseV0Response:
		e.Hash = header.Hash().String()
		e.ChainId = header.ChainID.String()
		e.Version = header.Version
		e.Height = header.Height
		e.Timestamp = header.Timestamp
		e.Previous = header.Previous.String()
		e.TransactionRoot = header.TransactionRoot.String()
		e.WitnessRoot = header.WitnessRoot.String()
		e.Proof.Nonce = header.Proof.(*pow.WorkProof).Nonce
		e.Proof.Target = header.Proof.(*pow.WorkProof).Target

	case *GetBlockVerboseV1Response:
		e.Hash = header.Hash().String()
		e.ChainId = header.ChainID.String()
		e.Version = header.Version
		e.Height = header.Height
		e.Timestamp = header.Timestamp
		e.Previous = header.Previous.String()
		e.TransactionRoot = header.TransactionRoot.String()
		e.WitnessRoot = header.WitnessRoot.String()
		e.Proof.Nonce = header.Proof.(*pow.WorkProof).Nonce
		e.Proof.Target = header.Proof.(*pow.WorkProof).Target
	default:

	}
}
