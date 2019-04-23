package api

import (
	"encoding/hex"
	"github.com/clarenous/go-capsule/protocol/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (a *API) GetEvidence(ctx context.Context, in *GetEvidenceRequest) (*GetEvidenceResponse, error) {
	id, err := types.NewHashFromString(in.Evid)
	if err != nil {
		logrus.Error(err)
		return nil, ErrInvalidEvidenceID
	}

	evid, tx, index, err := a.Chain.GetEvidence(&id)
	if err != nil {
		return nil, err
	}

	resp := new(GetEvidenceResponse)

	constructEvidenceResp(resp, evid, tx.Hash(), uint64(index))

	return resp, nil
}

func constructEvidenceResp(resp interface{}, evid *types.Evidence, txid types.Hash, index uint64) {
	switch e := resp.(type) {
	case *Evidence:
		e.Evid = evid.Hash(txid, index).String()
		e.Digest = hex.EncodeToString(evid.Digest)
		e.Source = hex.EncodeToString(evid.Source)
		e.ValidScript = hex.EncodeToString(evid.ValidScript)

	case *GetEvidenceResponse:
		e.Txid = txid.String()
		e.Index = index
		e.Evid = evid.Hash(txid, index).String()
		e.Digest = hex.EncodeToString(evid.Digest)
		e.Source = hex.EncodeToString(evid.Source)
		e.ValidScript = hex.EncodeToString(evid.ValidScript)

	default:

	}
}
