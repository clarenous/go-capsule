package leveldb

import (
	"os"
	"testing"

	dbm "github.com/tendermint/tmlibs/db"

	"github.com/clarenous/go-capsule/database/storage"

	"github.com/clarenous/go-capsule/protocol/state"
	"github.com/clarenous/go-capsule/testutil"
)

func TestSaveUtxoView(t *testing.T) {
	testDB := dbm.NewDB("testdb", "leveldb", "temp")
	batch := testDB.NewBatch()
	defer os.RemoveAll("temp")

	cases := []struct {
		hash      types.Hash
		utxoEntry *storage.UtxoEntry
		exist     bool
	}{
		{
			hash:      types.Hash{V0: 0},
			utxoEntry: storage.NewUtxoEntry(true, 0, true),
			exist:     true,
		},
		{
			hash:      types.Hash{V0: 1},
			utxoEntry: storage.NewUtxoEntry(true, 0, false),
			exist:     true,
		},
		{
			hash:      types.Hash{V0: 2},
			utxoEntry: storage.NewUtxoEntry(false, 0, false),
			exist:     true,
		},
		{
			hash:      types.Hash{V0: 3},
			utxoEntry: storage.NewUtxoEntry(false, 0, true),
			exist:     false,
		},
	}

	view := state.NewUtxoViewpoint()
	for _, c := range cases {
		view.Entries[c.hash] = c.utxoEntry
	}

	saveUtxoView(batch, view)
	batch.Write()

	for _, c := range cases {
		entry, err := getUtxo(testDB, &c.hash)

		if !c.exist {
			if err == nil {
				t.Errorf("%v should be unexisted, but it's in the db", c)
			}
			continue
		}

		if !testutil.DeepEqual(entry, c.utxoEntry) {
			t.Errorf("%v utxo in the db isn't match", c)
		}
	}
}

func TestGetTransactionsUtxo(t *testing.T) {
	testDB := dbm.NewDB("testdb", "leveldb", "temp")
	defer os.RemoveAll("temp")

	batch := testDB.NewBatch()
	inputView := state.NewUtxoViewpoint()
	for i := 0; i <= 2; i++ {
		inputView.Entries[types.Hash{V0: uint64(i)}] = storage.NewUtxoEntry(false, uint64(i), false)
	}
	saveUtxoView(batch, inputView)
	batch.Write()

	cases := []struct {
		txs       []*types.Tx
		inputView *state.UtxoViewpoint
		fetchView *state.UtxoViewpoint
		err       bool
	}{

		{
			txs: []*types.Tx{
				&types.Tx{
					SpentOutputIDs: []types.Hash{types.Hash{V0: 10}},
				},
			},
			inputView: state.NewUtxoViewpoint(),
			fetchView: state.NewUtxoViewpoint(),
			err:       false,
		},
		{
			txs: []*types.Tx{
				&types.Tx{
					SpentOutputIDs: []types.Hash{types.Hash{V0: 0}},
				},
			},
			inputView: state.NewUtxoViewpoint(),
			fetchView: &state.UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 0, false),
				},
			},
			err: false,
		},
		{
			txs: []*types.Tx{
				&types.Tx{
					SpentOutputIDs: []types.Hash{
						types.Hash{V0: 0},
						types.Hash{V0: 1},
					},
				},
			},
			inputView: state.NewUtxoViewpoint(),
			fetchView: &state.UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 0, false),
					types.Hash{V0: 1}: storage.NewUtxoEntry(false, 1, false),
				},
			},
			err: false,
		},
		{
			txs: []*types.Tx{
				&types.Tx{
					SpentOutputIDs: []types.Hash{
						types.Hash{V0: 0},
						types.Hash{V0: 1},
					},
				},
				&types.Tx{
					SpentOutputIDs: []types.Hash{
						types.Hash{V0: 2},
					},
				},
			},
			inputView: state.NewUtxoViewpoint(),
			fetchView: &state.UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 0, false),
					types.Hash{V0: 1}: storage.NewUtxoEntry(false, 1, false),
					types.Hash{V0: 2}: storage.NewUtxoEntry(false, 2, false),
				},
			},
			err: false,
		},
		{
			txs: []*types.Tx{
				&types.Tx{
					SpentOutputIDs: []types.Hash{types.Hash{V0: 0}},
				},
			},
			inputView: &state.UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 1, false),
				},
			},
			fetchView: &state.UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 1, false),
				},
			},
			err: false,
		},
	}

	for i, c := range cases {
		if err := getTransactionsUtxo(testDB, c.inputView, c.txs); c.err != (err != nil) {
			t.Errorf("test case %d, want err = %v, get err = %v", i, c.err, err)
		}
		if !testutil.DeepEqual(c.inputView, c.fetchView) {
			t.Errorf("test case %d, want %v, get %v", i, c.fetchView, c.inputView)
		}
	}
}
