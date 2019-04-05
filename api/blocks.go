package api

import (
	"encoding/hex"
	"github.com/clarenous/go-capsule/protocol/types"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"math/rand"
)

func (a *API) GetBestBlock(ctx context.Context, in *empty.Empty) (*GetBestBlockResponse, error) {
	resp := &GetBestBlockResponse{
		Height: rand.Uint64() % 1000,
		Hash:   types.MockHash().String(),
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

func (a *API) GetBlockVerboseV0(ctx context.Context, in *GetBlockVerboseRequest) (*GetBlockVerboseV0Response, error) {
	block := types.MockBlock()

	resp := &GetBlockVerboseV0Response{
		Proof:        &Proof{},
		Transactions: make([]*GetBlockVerboseV0Response_Transaction, len(block.Transactions)),
	}

	constructBlockHeaderResp(resp, &block.BlockHeader)
	for i, tx := range block.Transactions {
		resp.Transactions[i] = &GetBlockVerboseV0Response_Transaction{
			Txid:      types.MockHash().String(),
			Evidences: make([]string, len(tx.Evidences)),
		}
		for j := range tx.Evidences {
			resp.Transactions[i].Evidences[j] = types.MockHash().String()
		}
	}

	return resp, nil
}

func (a *API) GetBlockVerboseV1(ctx context.Context, in *GetBlockVerboseRequest) (*GetBlockVerboseV1Response, error) {
	block := types.MockBlock()

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

	case *GetBlockVerboseV0Response:
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

	case *GetBlockVerboseV1Response:
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

func constructTxResp(resp *Tx, tx *types.Tx) {
	resp.Txid = types.MockHash().String()
	resp.Version = tx.Version
	resp.Inputs = make([]*Tx_TxIn, len(tx.Inputs))
	resp.Outputs = make([]*Tx_TxOut, len(tx.Outputs))
	resp.Evidences = make([]*Evidence, len(tx.Evidences))
	resp.LockTime = tx.LockTime

	for i, in := range tx.Inputs {
		resp.Inputs[i] = &Tx_TxIn{
			ValueSource: &Tx_TxIn_ValueSource{
				Txid:  in.ValueSource.TxID.String(),
				Index: in.ValueSource.Index,
			},
			RedeemScript: hex.EncodeToString(in.RedeemScript),
			UnlockScript: hex.EncodeToString(in.UnlockScript),
			Sequence:     in.Sequence,
		}
	}

	for i, out := range tx.Outputs {
		resp.Outputs[i] = &Tx_TxOut{
			Value:      out.Value,
			ScriptHash: out.ScriptHash.String(),
		}
	}

	for i, evid := range tx.Evidences {
		evidResp := new(Evidence)
		constructEvidenceResp(evidResp, &evid)
		resp.Evidences[i] = evidResp
	}
}

func constructEvidenceResp(resp *Evidence, evid *types.Evidence) {
	resp.Evid = types.MockHash().String()
	resp.Digest = hex.EncodeToString(evid.Digest)
	resp.Source = hex.EncodeToString(evid.Source)
	resp.ValidScript = hex.EncodeToString(evid.ValidScript)
}
