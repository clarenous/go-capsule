package state

import (
	"github.com/clarenous/go-capsule/consensus"
	"github.com/clarenous/go-capsule/database/storage"
	"github.com/clarenous/go-capsule/errors"
	"github.com/clarenous/go-capsule/protocol/types"
)

var (
	ErrWrongCoinbaseTransaction = errors.New("wrong coinbase transaction")
)

// UtxoViewpoint represents a view into the set of unspent transaction outputs
type UtxoViewpoint struct {
	Entries map[types.Hash]*storage.UtxoEntry
}

// NewUtxoViewpoint returns a new empty unspent transaction output view.
func NewUtxoViewpoint() *UtxoViewpoint {
	return &UtxoViewpoint{
		Entries: make(map[types.Hash]*storage.UtxoEntry),
	}
}

func (view *UtxoViewpoint) ApplyTransaction(block *types.Block, tx *types.Tx) error {
	for _, in := range tx.Inputs {
		entry, ok := view.Entries[in.ValueSource.Hash()]
		if !ok {
			return errors.New("fail to find utxo entry")
		}
		if entry.Spent {
			return errors.New("utxo has been spent")
		}
		if entry.IsCoinBase && entry.BlockHeight+consensus.CoinbasePendingBlockNumber > block.Height {
			return errors.New("coinbase utxo is not ready for use")
		}
		entry.SpendOutput()
	}

	for i, _ := range tx.Outputs {
		isCoinbase := false
		if block != nil && len(block.Transactions) > 0 && block.Transactions[0].Hash() == tx.Hash() {
			isCoinbase = true
		}
		view.Entries[tx.OutHash(i)] = storage.NewUtxoEntry(isCoinbase, block.Height, false)
	}
	return nil
}

func (view *UtxoViewpoint) ApplyBlock(block *types.Block) error {
	// Check coinbase value
	var fee uint64
	for _, tx := range block.Transactions {
		for _, in := range tx.Inputs {
			entry, ok := view.Entries[in.ValueSource.Hash()]
			if !ok {
				return errors.New("fail to find utxo entry")
			}
			fee += entry.Value
		}
	}
	coinbaseAmount := consensus.BlockSubsidy(block.Height)
	if err := CheckCoinbaseAmount(block, coinbaseAmount+fee); err != nil {
		return err
	}

	// Check Inputs
	for _, tx := range block.Transactions {
		if err := view.ApplyTransaction(block, tx); err != nil {
			return err
		}
	}

	return nil
}

func (view *UtxoViewpoint) CanSpend(hash *types.Hash) bool {
	entry := view.Entries[*hash]
	return entry != nil && !entry.Spent
}

func (view *UtxoViewpoint) DetachTransaction(tx *types.Tx) error {
	for _, in := range tx.Inputs {
		utxoHash := in.ValueSource.Hash()
		entry, ok := view.Entries[utxoHash]
		if ok && !entry.Spent {
			return errors.New("try to revert an unspent utxo")
		}
		if !ok {
			view.Entries[utxoHash] = storage.NewUtxoEntry(false, 0, false)
			continue
		}
		entry.UnspendOutput()
	}

	for i, _ := range tx.Outputs {
		view.Entries[tx.OutHash(i)] = storage.NewUtxoEntry(false, 0, true)
	}
	return nil
}

func (view *UtxoViewpoint) DetachBlock(block *types.Block) error {
	for i := len(block.Transactions) - 1; i >= 0; i-- {
		if err := view.DetachTransaction(block.Transactions[i]); err != nil {
			return err
		}
	}
	return nil
}

func (view *UtxoViewpoint) HasUtxo(hash *types.Hash) bool {
	_, ok := view.Entries[*hash]
	return ok
}

// TODO: check overflow! (19.03.24 gcy)
func CheckCoinbaseAmount(b *types.Block, amount uint64) error {
	if len(b.Transactions) == 0 {
		return errors.Wrap(ErrWrongCoinbaseTransaction, "block is empty")
	}

	var totalOuts uint64
	for _, out := range b.Transactions[0].Outputs {
		totalOuts += out.Value
	}

	if totalOuts > amount {
		return errors.Wrap(ErrWrongCoinbaseTransaction, "reward more than deserved")
	}
	return nil
}
