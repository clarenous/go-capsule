package protocol

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/clarenous/go-capsule/consensus"
	"github.com/clarenous/go-capsule/database/storage"
	"github.com/clarenous/go-capsule/event"
	"github.com/clarenous/go-capsule/protocol/types"

	"github.com/clarenous/go-capsule/protocol/state"
	"github.com/clarenous/go-capsule/testutil"
)

var testTxs = []*types.Tx{
	//tx0
	types.NewTx(types.TxData{
		SerializedSize: 100,
		Inputs: []*types.TxInput{
			types.NewSpendInput(nil, types.NewHash([32]byte{0x01}), *consensus.BTMAssetID, 1, 1, []byte{0x51}),
		},
		Outputs: []*types.TxOutput{
			types.NewTxOutput(*consensus.BTMAssetID, 1, []byte{0x6a}),
		},
	}),
	//tx1
	types.NewTx(types.TxData{
		SerializedSize: 100,
		Inputs: []*types.TxInput{
			types.NewSpendInput(nil, types.NewHash([32]byte{0x01}), *consensus.BTMAssetID, 1, 1, []byte{0x51}),
		},
		Outputs: []*types.TxOutput{
			types.NewTxOutput(*consensus.BTMAssetID, 1, []byte{0x6b}),
		},
	}),
	//tx2
	types.NewTx(types.TxData{
		SerializedSize: 150,
		TimeRange:      0,
		Inputs: []*types.TxInput{
			types.NewSpendInput(nil, types.NewHash([32]byte{0x01}), *consensus.BTMAssetID, 1, 1, []byte{0x51}),
			types.NewSpendInput(nil, types.NewHash([32]byte{0x02}), types.NewAssetID([32]byte{0xa1}), 4, 1, []byte{0x51}),
		},
		Outputs: []*types.TxOutput{
			types.NewTxOutput(*consensus.BTMAssetID, 1, []byte{0x6b}),
			types.NewTxOutput(types.NewAssetID([32]byte{0xa1}), 4, []byte{0x61}),
		},
	}),
	//tx3
	types.NewTx(types.TxData{
		SerializedSize: 100,
		Inputs: []*types.TxInput{
			types.NewSpendInput(nil, testutil.MustDecodeHash("dbea684b5c5153ed7729669a53d6c59574f26015a3e1eb2a0e8a1c645425a764"), types.NewAssetID([32]byte{0xa1}), 4, 1, []byte{0x61}),
		},
		Outputs: []*types.TxOutput{
			types.NewTxOutput(types.NewAssetID([32]byte{0xa1}), 3, []byte{0x62}),
			types.NewTxOutput(types.NewAssetID([32]byte{0xa1}), 1, []byte{0x63}),
		},
	}),
	//tx4
	types.NewTx(types.TxData{
		SerializedSize: 100,
		Inputs: []*types.TxInput{
			types.NewSpendInput(nil, testutil.MustDecodeHash("d84d0be0fd08e7341f2d127749bb0d0844d4560f53bd54861cee9981fd922cad"), types.NewAssetID([32]byte{0xa1}), 3, 0, []byte{0x62}),
		},
		Outputs: []*types.TxOutput{
			types.NewTxOutput(types.NewAssetID([32]byte{0xa1}), 2, []byte{0x64}),
			types.NewTxOutput(types.NewAssetID([32]byte{0xa1}), 1, []byte{0x65}),
		},
	}),
}

type mockStore struct{}

func (s *mockStore) BlockExist(hash *types.Hash) bool                                { return false }
func (s *mockStore) GetBlock(*types.Hash) (*types.Block, error)                      { return nil, nil }
func (s *mockStore) GetStoreStatus() *BlockStoreState                             { return nil }
func (s *mockStore) GetTransactionStatus(*types.Hash) (*types.TransactionStatus, error) { return nil, nil }
func (s *mockStore) GetTransactionsUtxo(*state.UtxoViewpoint, []*types.Tx) error     { return nil }
func (s *mockStore) GetUtxo(*types.Hash) (*storage.UtxoEntry, error)                 { return nil, nil }
func (s *mockStore) LoadBlockIndex(uint64) (*state.BlockIndex, error)             { return nil, nil }
func (s *mockStore) SaveBlock(*types.Block, *types.TransactionStatus) error          { return nil }
func (s *mockStore) SaveChainStatus(*state.BlockNode, *state.UtxoViewpoint) error { return nil }

func TestAddOrphan(t *testing.T) {
	cases := []struct {
		before         *TxPool
		after          *TxPool
		addOrphan      *TxDesc
		requireParents []*types.Hash
	}{
		{
			before: &TxPool{
				orphans:       map[types.Hash]*orphanTx{},
				orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{},
			},
			after: &TxPool{
				orphans: map[types.Hash]*orphanTx{
					testTxs[0].ID: {
						TxDesc: &TxDesc{
							Tx: testTxs[0],
						},
					},
				},
				orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{
					testTxs[0].SpentOutputIDs[0]: {
						testTxs[0].ID: {
							TxDesc: &TxDesc{
								Tx: testTxs[0],
							},
						},
					},
				},
			},
			addOrphan:      &TxDesc{Tx: testTxs[0]},
			requireParents: []*types.Hash{&testTxs[0].SpentOutputIDs[0]},
		},
		{
			before: &TxPool{
				orphans: map[types.Hash]*orphanTx{
					testTxs[0].ID: {
						TxDesc: &TxDesc{
							Tx: testTxs[0],
						},
					},
				},
				orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{
					testTxs[0].SpentOutputIDs[0]: {
						testTxs[0].ID: {
							TxDesc: &TxDesc{
								Tx: testTxs[0],
							},
						},
					},
				},
			},
			after: &TxPool{
				orphans: map[types.Hash]*orphanTx{
					testTxs[0].ID: {
						TxDesc: &TxDesc{
							Tx: testTxs[0],
						},
					},
					testTxs[1].ID: {
						TxDesc: &TxDesc{
							Tx: testTxs[1],
						},
					},
				},
				orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{
					testTxs[0].SpentOutputIDs[0]: {
						testTxs[0].ID: {
							TxDesc: &TxDesc{
								Tx: testTxs[0],
							},
						},
						testTxs[1].ID: {
							TxDesc: &TxDesc{
								Tx: testTxs[1],
							},
						},
					},
				},
			},
			addOrphan:      &TxDesc{Tx: testTxs[1]},
			requireParents: []*types.Hash{&testTxs[1].SpentOutputIDs[0]},
		},
		{
			before: &TxPool{
				orphans:       map[types.Hash]*orphanTx{},
				orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{},
			},
			after: &TxPool{
				orphans: map[types.Hash]*orphanTx{
					testTxs[2].ID: {
						TxDesc: &TxDesc{
							Tx: testTxs[2],
						},
					},
				},
				orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{
					testTxs[2].SpentOutputIDs[1]: {
						testTxs[2].ID: {
							TxDesc: &TxDesc{
								Tx: testTxs[2],
							},
						},
					},
				},
			},
			addOrphan:      &TxDesc{Tx: testTxs[2]},
			requireParents: []*types.Hash{&testTxs[2].SpentOutputIDs[1]},
		},
	}

	for i, c := range cases {
		c.before.addOrphan(c.addOrphan, c.requireParents)
		for _, orphan := range c.before.orphans {
			orphan.expiration = time.Time{}
		}
		for _, orphans := range c.before.orphansByPrev {
			for _, orphan := range orphans {
				orphan.expiration = time.Time{}
			}
		}
		if !testutil.DeepEqual(c.before, c.after) {
			t.Errorf("case %d: got %v want %v", i, c.before, c.after)
		}
	}
}

func TestAddTransaction(t *testing.T) {
	dispatcher := event.NewDispatcher()
	cases := []struct {
		before *TxPool
		after  *TxPool
		addTx  *TxDesc
	}{
		{
			before: &TxPool{
				pool:            map[types.Hash]*TxDesc{},
				utxo:            map[types.Hash]*types.Tx{},
				eventDispatcher: dispatcher,
			},
			after: &TxPool{
				pool: map[types.Hash]*TxDesc{
					testTxs[2].ID: {
						Tx:         testTxs[2],
						StatusFail: false,
					},
				},
				utxo: map[types.Hash]*types.Tx{
					*testTxs[2].ResultIds[0]: testTxs[2],
					*testTxs[2].ResultIds[1]: testTxs[2],
				},
			},
			addTx: &TxDesc{
				Tx:         testTxs[2],
				StatusFail: false,
			},
		},
		{
			before: &TxPool{
				pool:            map[types.Hash]*TxDesc{},
				utxo:            map[types.Hash]*types.Tx{},
				eventDispatcher: dispatcher,
			},
			after: &TxPool{
				pool: map[types.Hash]*TxDesc{
					testTxs[2].ID: {
						Tx:         testTxs[2],
						StatusFail: true,
					},
				},
				utxo: map[types.Hash]*types.Tx{
					*testTxs[2].ResultIds[0]: testTxs[2],
				},
			},
			addTx: &TxDesc{
				Tx:         testTxs[2],
				StatusFail: true,
			},
		},
	}

	for i, c := range cases {
		c.before.addTransaction(c.addTx)
		for _, txD := range c.before.pool {
			txD.Added = time.Time{}
		}
		if !testutil.DeepEqual(c.before.pool, c.after.pool) {
			t.Errorf("case %d: got %v want %v", i, c.before.pool, c.after.pool)
		}
		if !testutil.DeepEqual(c.before.utxo, c.after.utxo) {
			t.Errorf("case %d: got %v want %v", i, c.before.utxo, c.after.utxo)
		}
	}
}

func TestExpireOrphan(t *testing.T) {
	before := &TxPool{
		orphans: map[types.Hash]*orphanTx{
			testTxs[0].ID: {
				expiration: time.Unix(1533489701, 0),
				TxDesc: &TxDesc{
					Tx: testTxs[0],
				},
			},
			testTxs[1].ID: {
				expiration: time.Unix(1633489701, 0),
				TxDesc: &TxDesc{
					Tx: testTxs[1],
				},
			},
		},
		orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{
			testTxs[0].SpentOutputIDs[0]: {
				testTxs[0].ID: {
					expiration: time.Unix(1533489701, 0),
					TxDesc: &TxDesc{
						Tx: testTxs[0],
					},
				},
				testTxs[1].ID: {
					expiration: time.Unix(1633489701, 0),
					TxDesc: &TxDesc{
						Tx: testTxs[1],
					},
				},
			},
		},
	}

	want := &TxPool{
		orphans: map[types.Hash]*orphanTx{
			testTxs[1].ID: {
				expiration: time.Unix(1633489701, 0),
				TxDesc: &TxDesc{
					Tx: testTxs[1],
				},
			},
		},
		orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{
			testTxs[0].SpentOutputIDs[0]: {
				testTxs[1].ID: {
					expiration: time.Unix(1633489701, 0),
					TxDesc: &TxDesc{
						Tx: testTxs[1],
					},
				},
			},
		},
	}

	before.ExpireOrphan(time.Unix(1633479701, 0))
	if !testutil.DeepEqual(before, want) {
		t.Errorf("got %v want %v", before, want)
	}
}

func TestProcessOrphans(t *testing.T) {
	dispatcher := event.NewDispatcher()
	cases := []struct {
		before    *TxPool
		after     *TxPool
		processTx *TxDesc
	}{
		{
			before: &TxPool{
				pool:            map[types.Hash]*TxDesc{},
				utxo:            map[types.Hash]*types.Tx{},
				eventDispatcher: dispatcher,
				orphans: map[types.Hash]*orphanTx{
					testTxs[3].ID: {
						TxDesc: &TxDesc{
							Tx: testTxs[3],
						},
					},
				},
				orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{
					testTxs[3].SpentOutputIDs[0]: {
						testTxs[3].ID: {
							TxDesc: &TxDesc{
								Tx: testTxs[3],
							},
						},
					},
				},
			},
			after: &TxPool{
				pool: map[types.Hash]*TxDesc{
					testTxs[3].ID: {
						Tx:         testTxs[3],
						StatusFail: false,
					},
				},
				utxo: map[types.Hash]*types.Tx{
					*testTxs[3].ResultIds[0]: testTxs[3],
					*testTxs[3].ResultIds[1]: testTxs[3],
				},
				eventDispatcher: dispatcher,
				orphans:         map[types.Hash]*orphanTx{},
				orphansByPrev:   map[types.Hash]map[types.Hash]*orphanTx{},
			},
			processTx: &TxDesc{Tx: testTxs[2]},
		},
		{
			before: &TxPool{
				pool:            map[types.Hash]*TxDesc{},
				utxo:            map[types.Hash]*types.Tx{},
				eventDispatcher: dispatcher,
				orphans: map[types.Hash]*orphanTx{
					testTxs[3].ID: {
						TxDesc: &TxDesc{
							Tx: testTxs[3],
						},
					},
					testTxs[4].ID: {
						TxDesc: &TxDesc{
							Tx: testTxs[4],
						},
					},
				},
				orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{
					testTxs[3].SpentOutputIDs[0]: {
						testTxs[3].ID: {
							TxDesc: &TxDesc{
								Tx: testTxs[3],
							},
						},
					},
					testTxs[4].SpentOutputIDs[0]: {
						testTxs[4].ID: {
							TxDesc: &TxDesc{
								Tx: testTxs[4],
							},
						},
					},
				},
			},
			after: &TxPool{
				pool: map[types.Hash]*TxDesc{
					testTxs[3].ID: {
						Tx:         testTxs[3],
						StatusFail: false,
					},
					testTxs[4].ID: {
						Tx:         testTxs[4],
						StatusFail: false,
					},
				},
				utxo: map[types.Hash]*types.Tx{
					*testTxs[3].ResultIds[0]: testTxs[3],
					*testTxs[3].ResultIds[1]: testTxs[3],
					*testTxs[4].ResultIds[0]: testTxs[4],
					*testTxs[4].ResultIds[1]: testTxs[4],
				},
				eventDispatcher: dispatcher,
				orphans:         map[types.Hash]*orphanTx{},
				orphansByPrev:   map[types.Hash]map[types.Hash]*orphanTx{},
			},
			processTx: &TxDesc{Tx: testTxs[2]},
		},
	}

	for i, c := range cases {
		c.before.store = &mockStore{}
		c.before.addTransaction(c.processTx)
		c.before.processOrphans(c.processTx)
		c.before.RemoveTransaction(&c.processTx.Tx.ID)
		c.before.store = nil
		c.before.lastUpdated = 0
		for _, txD := range c.before.pool {
			txD.Added = time.Time{}
		}

		if !testutil.DeepEqual(c.before, c.after) {
			t.Errorf("case %d: got %v want %v", i, c.before, c.after)
		}
	}
}

func TestRemoveOrphan(t *testing.T) {
	cases := []struct {
		before       *TxPool
		after        *TxPool
		removeHashes []*types.Hash
	}{
		{
			before: &TxPool{
				orphans: map[types.Hash]*orphanTx{
					testTxs[0].ID: {
						expiration: time.Unix(1533489701, 0),
						TxDesc: &TxDesc{
							Tx: testTxs[0],
						},
					},
				},
				orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{
					testTxs[0].SpentOutputIDs[0]: {
						testTxs[0].ID: {
							expiration: time.Unix(1533489701, 0),
							TxDesc: &TxDesc{
								Tx: testTxs[0],
							},
						},
					},
				},
			},
			after: &TxPool{
				orphans:       map[types.Hash]*orphanTx{},
				orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{},
			},
			removeHashes: []*types.Hash{
				&testTxs[0].ID,
			},
		},
		{
			before: &TxPool{
				orphans: map[types.Hash]*orphanTx{
					testTxs[0].ID: {
						expiration: time.Unix(1533489701, 0),
						TxDesc: &TxDesc{
							Tx: testTxs[0],
						},
					},
					testTxs[1].ID: {
						expiration: time.Unix(1533489701, 0),
						TxDesc: &TxDesc{
							Tx: testTxs[1],
						},
					},
				},
				orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{
					testTxs[0].SpentOutputIDs[0]: {
						testTxs[0].ID: {
							expiration: time.Unix(1533489701, 0),
							TxDesc: &TxDesc{
								Tx: testTxs[0],
							},
						},
						testTxs[1].ID: {
							expiration: time.Unix(1533489701, 0),
							TxDesc: &TxDesc{
								Tx: testTxs[1],
							},
						},
					},
				},
			},
			after: &TxPool{
				orphans: map[types.Hash]*orphanTx{
					testTxs[0].ID: {
						expiration: time.Unix(1533489701, 0),
						TxDesc: &TxDesc{
							Tx: testTxs[0],
						},
					},
				},
				orphansByPrev: map[types.Hash]map[types.Hash]*orphanTx{
					testTxs[0].SpentOutputIDs[0]: {
						testTxs[0].ID: {
							expiration: time.Unix(1533489701, 0),
							TxDesc: &TxDesc{
								Tx: testTxs[0],
							},
						},
					},
				},
			},
			removeHashes: []*types.Hash{
				&testTxs[1].ID,
			},
		},
	}

	for i, c := range cases {
		for _, hash := range c.removeHashes {
			c.before.removeOrphan(hash)
		}
		if !testutil.DeepEqual(c.before, c.after) {
			t.Errorf("case %d: got %v want %v", i, c.before, c.after)
		}
	}
}

type mockStore1 struct{}

func (s *mockStore1) BlockExist(hash *types.Hash) bool                                { return false }
func (s *mockStore1) GetBlock(*types.Hash) (*types.Block, error)                      { return nil, nil }
func (s *mockStore1) GetStoreStatus() *BlockStoreState                             { return nil }
func (s *mockStore1) GetTransactionStatus(*types.Hash) (*types.TransactionStatus, error) { return nil, nil }
func (s *mockStore1) GetTransactionsUtxo(utxoView *state.UtxoViewpoint, tx []*types.Tx) error {
	for _, hash := range testTxs[2].SpentOutputIDs {
		utxoView.Entries[hash] = &storage.UtxoEntry{IsCoinBase: false, Spent: false}
	}
	return nil
}
func (s *mockStore1) GetUtxo(*types.Hash) (*storage.UtxoEntry, error)                 { return nil, nil }
func (s *mockStore1) LoadBlockIndex(uint64) (*state.BlockIndex, error)             { return nil, nil }
func (s *mockStore1) SaveBlock(*types.Block, *types.TransactionStatus) error          { return nil }
func (s *mockStore1) SaveChainStatus(*state.BlockNode, *state.UtxoViewpoint) error { return nil }

func TestProcessTransaction(t *testing.T) {
	txPool := &TxPool{
		pool:            make(map[types.Hash]*TxDesc),
		utxo:            make(map[types.Hash]*types.Tx),
		orphans:         make(map[types.Hash]*orphanTx),
		orphansByPrev:   make(map[types.Hash]map[types.Hash]*orphanTx),
		store:           &mockStore1{},
		eventDispatcher: event.NewDispatcher(),
	}
	cases := []struct {
		want  *TxPool
		addTx *TxDesc
	}{
		//Dust tx
		{
			want: &TxPool{},
			addTx: &TxDesc{
				Tx:         testTxs[3],
				StatusFail: false,
			},
		},
		//normal tx
		{
			want: &TxPool{
				pool: map[types.Hash]*TxDesc{
					testTxs[2].ID: {
						Tx:         testTxs[2],
						StatusFail: false,
						Weight:     150,
					},
				},
				utxo: map[types.Hash]*types.Tx{
					*testTxs[2].ResultIds[0]: testTxs[2],
					*testTxs[2].ResultIds[1]: testTxs[2],
				},
			},
			addTx: &TxDesc{
				Tx:         testTxs[2],
				StatusFail: false,
			},
		},
	}

	for i, c := range cases {
		txPool.ProcessTransaction(c.addTx.Tx, c.addTx.StatusFail, 0, 0)
		for _, txD := range txPool.pool {
			txD.Added = time.Time{}
		}
		for _, txD := range txPool.orphans {
			txD.Added = time.Time{}
			txD.expiration = time.Time{}
		}

		if !testutil.DeepEqual(txPool.pool, c.want.pool) {
			t.Errorf("case %d: test ProcessTransaction pool mismatch got %s want %s", i, spew.Sdump(txPool.pool), spew.Sdump(c.want.pool))
		}
		if !testutil.DeepEqual(txPool.utxo, c.want.utxo) {
			t.Errorf("case %d: test ProcessTransaction utxo mismatch got %s want %s", i, spew.Sdump(txPool.utxo), spew.Sdump(c.want.utxo))
		}
		if !testutil.DeepEqual(txPool.orphans, c.want.orphans) {
			t.Errorf("case %d: test ProcessTransaction orphans mismatch got %s want %s", i, spew.Sdump(txPool.orphans), spew.Sdump(c.want.orphans))
		}
		if !testutil.DeepEqual(txPool.orphansByPrev, c.want.orphansByPrev) {
			t.Errorf("case %d: test ProcessTransaction orphansByPrev mismatch got %s want %s", i, spew.Sdump(txPool.orphansByPrev), spew.Sdump(c.want.orphansByPrev))
		}
	}
}
