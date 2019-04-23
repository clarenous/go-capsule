package protocol

import (
	log "github.com/sirupsen/logrus"

	"github.com/clarenous/go-capsule/errors"
	"github.com/clarenous/go-capsule/protocol/types"

	"github.com/clarenous/go-capsule/protocol/state"
	"github.com/clarenous/go-capsule/protocol/validation"
)

// ErrBadTx is returned for transactions failing validation
var ErrBadTx = errors.New("invalid transaction")

// GetTransactionsUtxo return all the utxos that related to the txs' inputs
func (c *Chain) GetTransactionsUtxo(view *state.UtxoViewpoint, txs []*types.Tx) error {
	return c.store.GetTransactionsUtxo(view, txs)
}

// ValidateTx validates the given transaction. A cache holds
// per-transaction validation results and is consulted before
// performing full validation.
func (c *Chain) ValidateTx(tx *types.Tx) (bool, error) {
	if ok := c.txPool.HaveTransaction(tx.Hash().Ptr()); ok {
		return false, c.txPool.GetErrCache(tx.Hash().Ptr())
	}

	if c.txPool.IsDust(tx) {
		c.txPool.AddErrCache(tx.Hash().Ptr(), ErrDustTx)
		return false, ErrDustTx
	}

	bh := c.BestBlockHeader()
	err := validation.ValidateTx(tx, &types.Block{BlockHeader: *bh})

	if err != nil {
		log.WithFields(log.Fields{"module": logModule, "tx_id": tx.Hash().String(), "error": err}).Info("transaction status fail")
	}

	return c.txPool.ProcessTransaction(tx, bh.Height)
}

func (c *Chain) GetTransaction(hash *types.Hash) (*types.Tx, error) {
	return c.store.GetTransaction(hash)
}

func (c *Chain) GetEvidence(hash *types.Hash) (*types.Evidence, *types.Tx, int, error) {
	return c.store.GetEvidence(hash)
}
