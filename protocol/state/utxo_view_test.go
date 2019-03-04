package state

import (
	"testing"

	"github.com/clarenous/go-capsule/consensus"
	"github.com/clarenous/go-capsule/database/storage"
	"github.com/clarenous/go-capsule/protocol/types"
	"github.com/clarenous/go-capsule/testutil"
)

var defaultEntry = map[types.Hash]types.Entry{
	types.Hash{V0: 0}: &types.Output{
		Source: &types.ValueSource{
			Value: &types.AssetAmount{
				AssetId: &types.AssetID{V0: 0},
			},
		},
	},
}

var gasOnlyTxEntry = map[types.Hash]types.Entry{
	types.Hash{V1: 0}: &types.Output{
		Source: &types.ValueSource{
			Value: &types.AssetAmount{
				AssetId: consensus.BTMAssetID,
			},
		},
	},
	types.Hash{V1: 1}: &types.Output{
		Source: &types.ValueSource{
			Value: &types.AssetAmount{
				AssetId: &types.AssetID{V0: 999},
			},
		},
	},
}

func TestApplyBlock(t *testing.T) {
	cases := []struct {
		block     *types.Block
		inputView *UtxoViewpoint
		fetchView *UtxoViewpoint
		gasOnlyTx bool
		err       bool
	}{
		{
			// can't find prevout in tx entries
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						SpentOutputIDs: []types.Hash{
							types.Hash{V0: 1},
						},
						Entries: defaultEntry,
					},
				},
			},
			inputView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 0, false),
				},
			},
			fetchView: NewUtxoViewpoint(),
			err:       true,
		},
		{
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						SpentOutputIDs: []types.Hash{
							types.Hash{V0: 0},
						},
						Entries: defaultEntry,
					},
				},
			},
			inputView: NewUtxoViewpoint(),
			fetchView: NewUtxoViewpoint(),
			err:       true,
		},
		{
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						SpentOutputIDs: []types.Hash{
							types.Hash{V0: 0},
						},
						Entries: defaultEntry,
					},
				},
			},
			inputView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 0, true),
				},
			},
			err: true,
		},
		{
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						TxHeader: &types.TxHeader{
							ResultIds: []*types.Hash{},
						},
						SpentOutputIDs: []types.Hash{
							types.Hash{V0: 0},
						},
						Entries: defaultEntry,
					},
				},
			},
			inputView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 0, false),
				},
			},
			fetchView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 0, true),
				},
			},
			err: false,
		},
		{
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					Height:            101,
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						TxHeader: &types.TxHeader{
							ResultIds: []*types.Hash{},
						},
						SpentOutputIDs: []types.Hash{
							types.Hash{V0: 0},
						},
						Entries: defaultEntry,
					},
				},
			},
			inputView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(true, 0, false),
				},
			},
			fetchView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(true, 0, true),
				},
			},
			err: false,
		},
		{
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					Height:            0,
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						TxHeader: &types.TxHeader{
							ResultIds: []*types.Hash{},
						},
						SpentOutputIDs: []types.Hash{
							types.Hash{V0: 0},
						},
						Entries: defaultEntry,
					},
				},
			},
			inputView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(true, 0, false),
				},
			},
			fetchView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(true, 0, true),
				},
			},
			err: true,
		},
		{
			// output will be store
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						TxHeader: &types.TxHeader{
							ResultIds: []*types.Hash{
								&types.Hash{V0: 0},
							},
						},
						SpentOutputIDs: []types.Hash{},
						Entries:        defaultEntry,
					},
				},
			},
			inputView: NewUtxoViewpoint(),
			fetchView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(true, 0, false),
				},
			},
			err: false,
		},
		{
			// apply gas only tx, non-btm asset spent input will not be spent
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						TxHeader: &types.TxHeader{
							ResultIds: []*types.Hash{},
						},
						SpentOutputIDs: []types.Hash{
							types.Hash{V1: 0},
							types.Hash{V1: 1},
						},
						Entries: gasOnlyTxEntry,
					},
				},
			},
			inputView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V1: 0}: storage.NewUtxoEntry(false, 0, false),
					types.Hash{V1: 1}: storage.NewUtxoEntry(false, 0, false),
				},
			},
			fetchView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V1: 0}: storage.NewUtxoEntry(false, 0, true),
					types.Hash{V1: 1}: storage.NewUtxoEntry(false, 0, false),
				},
			},
			gasOnlyTx: true,
			err:       false,
		},
		{
			// apply gas only tx, non-btm asset spent output will not be store
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						TxHeader: &types.TxHeader{
							ResultIds: []*types.Hash{
								&types.Hash{V1: 0},
								&types.Hash{V1: 1},
							},
						},
						SpentOutputIDs: []types.Hash{},
						Entries:        gasOnlyTxEntry,
					},
				},
			},
			inputView: NewUtxoViewpoint(),
			fetchView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V1: 0}: storage.NewUtxoEntry(true, 0, false),
				},
			},
			gasOnlyTx: true,
			err:       false,
		},
	}

	for i, c := range cases {
		c.block.TransactionStatus.SetStatus(0, c.gasOnlyTx)
		if err := c.inputView.ApplyBlock(c.block, c.block.TransactionStatus); c.err != (err != nil) {
			t.Errorf("case #%d want err = %v, get err = %v", i, c.err, err)
		}
		if c.err {
			continue
		}
		if !testutil.DeepEqual(c.inputView, c.fetchView) {
			t.Errorf("test case %d, want %v, get %v", i, c.fetchView, c.inputView)
		}
	}
}

func TestDetachBlock(t *testing.T) {
	cases := []struct {
		block     *types.Block
		inputView *UtxoViewpoint
		fetchView *UtxoViewpoint
		gasOnlyTx bool
		err       bool
	}{
		{
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						TxHeader: &types.TxHeader{
							ResultIds: []*types.Hash{},
						},
						SpentOutputIDs: []types.Hash{
							types.Hash{V0: 0},
						},
						Entries: defaultEntry,
					},
				},
			},
			inputView: NewUtxoViewpoint(),
			fetchView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 0, false),
				},
			},
			err: false,
		},
		{
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						TxHeader: &types.TxHeader{
							ResultIds: []*types.Hash{
								&types.Hash{V0: 0},
							},
						},
						SpentOutputIDs: []types.Hash{},
						Entries:        defaultEntry,
					},
				},
			},
			inputView: NewUtxoViewpoint(),
			fetchView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 0, true),
				},
			},
			err: false,
		},
		{
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						TxHeader: &types.TxHeader{
							ResultIds: []*types.Hash{},
						},
						SpentOutputIDs: []types.Hash{
							types.Hash{V0: 0},
						},
						Entries: defaultEntry,
					},
				},
			},
			inputView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 0, false),
				},
			},
			err: true,
		},
		{
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						TxHeader: &types.TxHeader{
							ResultIds: []*types.Hash{},
						},
						SpentOutputIDs: []types.Hash{
							types.Hash{V0: 0},
						},
						Entries: defaultEntry,
					},
				},
			},
			inputView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 0, true),
				},
			},
			fetchView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V0: 0}: storage.NewUtxoEntry(false, 0, false),
				},
			},
			err: false,
		},
		{
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						TxHeader: &types.TxHeader{
							ResultIds: []*types.Hash{},
						},
						SpentOutputIDs: []types.Hash{
							types.Hash{V1: 0},
							types.Hash{V1: 1},
						},
						Entries: gasOnlyTxEntry,
					},
				},
			},
			inputView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V1: 0}: storage.NewUtxoEntry(false, 0, true),
					types.Hash{V1: 1}: storage.NewUtxoEntry(false, 0, true),
				},
			},
			fetchView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V1: 0}: storage.NewUtxoEntry(false, 0, false),
					types.Hash{V1: 1}: storage.NewUtxoEntry(false, 0, true),
				},
			},
			gasOnlyTx: true,
			err:       false,
		},
		{
			block: &types.Block{
				BlockHeader: &types.BlockHeader{
					TransactionStatus: types.NewTransactionStatus(),
				},
				Transactions: []*types.Tx{
					&types.Tx{
						TxHeader: &types.TxHeader{
							ResultIds: []*types.Hash{
								&types.Hash{V1: 0},
								&types.Hash{V1: 1},
							},
						},
						SpentOutputIDs: []types.Hash{},
						Entries:        gasOnlyTxEntry,
					},
				},
			},
			inputView: NewUtxoViewpoint(),
			fetchView: &UtxoViewpoint{
				Entries: map[types.Hash]*storage.UtxoEntry{
					types.Hash{V1: 0}: storage.NewUtxoEntry(false, 0, true),
				},
			},
			gasOnlyTx: true,
			err:       false,
		},
	}

	for i, c := range cases {
		c.block.TransactionStatus.SetStatus(0, c.gasOnlyTx)
		if err := c.inputView.DetachBlock(c.block, c.block.TransactionStatus); c.err != (err != nil) {
			t.Errorf("case %d want err = %v, get err = %v", i, c.err, err)
		}
		if c.err {
			continue
		}
		if !testutil.DeepEqual(c.inputView, c.fetchView) {
			t.Errorf("test case %d, want %v, get %v", i, c.fetchView, c.inputView)
		}
	}
}
