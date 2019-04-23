package api

import (
	"encoding/hex"
	"github.com/clarenous/go-capsule/protocol/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (a *API) GetTransaction(ctx context.Context, in *GetTransactionRequest) (*GetTransactionResponse, error) {
	id, err := types.NewHashFromString(in.Txid)
	if err != nil {
		logrus.Error(err)
		return nil, ErrInvalidTransactionID
	}

	tx, err := a.Chain.GetTransaction(&id)
	if err != nil {
		return nil, err
	}

	resp := new(GetTransactionResponse)

	constructTxResp(resp, tx)

	return resp, nil
}

func constructTxResp(resp interface{}, tx *types.Tx) {
	txid := tx.Hash()

	switch e := resp.(type) {
	case *Tx:
		e.Txid = txid.String()
		e.Version = tx.Version
		e.Inputs = make([]*Tx_TxIn, len(tx.Inputs))
		e.Outputs = make([]*Tx_TxOut, len(tx.Outputs))
		e.Evidences = make([]*Evidence, len(tx.Evidences))
		e.LockTime = tx.LockTime

		for i, in := range tx.Inputs {
			e.Inputs[i] = &Tx_TxIn{
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
			e.Outputs[i] = &Tx_TxOut{
				Value:      out.Value,
				ScriptHash: out.ScriptHash.String(),
			}
		}

		for i, evid := range tx.Evidences {
			evidResp := new(Evidence)
			constructEvidenceResp(evidResp, &evid, txid, uint64(i))
			e.Evidences[i] = evidResp
		}

	case *GetTransactionResponse:
		e.Txid = txid.String()
		e.Version = tx.Version
		e.Inputs = make([]*GetTransactionResponse_TxIn, len(tx.Inputs))
		e.Outputs = make([]*GetTransactionResponse_TxOut, len(tx.Outputs))
		e.Evidences = make([]*Evidence, len(tx.Evidences))
		e.LockTime = tx.LockTime

		for i, in := range tx.Inputs {
			e.Inputs[i] = &GetTransactionResponse_TxIn{
				ValueSource: &GetTransactionResponse_TxIn_ValueSource{
					Txid:  in.ValueSource.TxID.String(),
					Index: in.ValueSource.Index,
				},
				RedeemScript: hex.EncodeToString(in.RedeemScript),
				UnlockScript: hex.EncodeToString(in.UnlockScript),
				Sequence:     in.Sequence,
			}
		}

		for i, out := range tx.Outputs {
			e.Outputs[i] = &GetTransactionResponse_TxOut{
				Value:      out.Value,
				ScriptHash: out.ScriptHash.String(),
			}
		}

		for i, evid := range tx.Evidences {
			evidResp := new(Evidence)
			constructEvidenceResp(evidResp, &evid, txid, uint64(i))
			e.Evidences[i] = evidResp
		}

	default:

	}
}
