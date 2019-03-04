package leveldb

import (
	dbm "github.com/tendermint/tmlibs/db"

	"github.com/clarenous/go-capsule/database/storage"
	"github.com/clarenous/go-capsule/errors"

	"github.com/clarenous/go-capsule/protocol/state"
	"github.com/golang/protobuf/proto"
)

const utxoPreFix = "UT:"

func calcUtxoKey(hash *types.Hash) []byte {
	return []byte(utxoPreFix + hash.String())
}

func getTransactionsUtxo(db dbm.DB, view *state.UtxoViewpoint, txs []*types.Tx) error {
	for _, tx := range txs {
		for _, prevout := range tx.SpentOutputIDs {
			if view.HasUtxo(&prevout) {
				continue
			}

			data := db.Get(calcUtxoKey(&prevout))
			if data == nil {
				continue
			}

			var utxo storage.UtxoEntry
			if err := proto.Unmarshal(data, &utxo); err != nil {
				return errors.Wrap(err, "unmarshaling utxo entry")
			}

			view.Entries[prevout] = &utxo
		}
	}

	return nil
}

func getUtxo(db dbm.DB, hash *types.Hash) (*storage.UtxoEntry, error) {
	var utxo storage.UtxoEntry
	data := db.Get(calcUtxoKey(hash))
	if data == nil {
		return nil, errors.New("can't find utxo in db")
	}
	if err := proto.Unmarshal(data, &utxo); err != nil {
		return nil, errors.Wrap(err, "unmarshaling utxo entry")
	}
	return &utxo, nil
}

func saveUtxoView(batch dbm.Batch, view *state.UtxoViewpoint) error {
	for key, entry := range view.Entries {
		if entry.Spent && !entry.IsCoinBase {
			batch.Delete(calcUtxoKey(&key))
			continue
		}

		b, err := proto.Marshal(entry)
		if err != nil {
			return errors.Wrap(err, "marshaling utxo entry")
		}
		batch.Set(calcUtxoKey(&key), b)
	}
	return nil
}

func SaveUtxoView(batch dbm.Batch, view *state.UtxoViewpoint) error {
	return saveUtxoView(batch, view)
}
