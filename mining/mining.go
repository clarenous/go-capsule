package mining

import (
	"github.com/clarenous/go-capsule/consensus"
	"github.com/clarenous/go-capsule/errors"
	"github.com/clarenous/go-capsule/protocol"
	"github.com/clarenous/go-capsule/protocol/types"
	"time"
)

// createCoinbaseTx returns a coinbase transaction paying an appropriate subsidy
// based on the passed block height to the provided address.  When the address
// is nil, the coinbase transaction will instead be redeemable by anyone.
func createCoinbaseTx(amount uint64, blockHeight uint64) (tx *types.Tx, err error) {
	amount += consensus.BlockSubsidy(blockHeight)
	tx = types.MockTx()
	txOut := types.MockTxOut()
	txOut.Value = amount
	tx.Inputs = []types.TxIn{}
	tx.Outputs = []types.TxOut{*txOut}

	return tx, nil
}

// NewBlockTemplate returns a new block template that is ready to be solved
func NewBlockTemplate(c *protocol.Chain, txPool *protocol.TxPool) (b *types.Block, err error) {
	//view := state.NewUtxoViewpoint()
	//txEntries := []*bc.Tx{nil}
	//txFee := uint64(0)
	//
	//// get preblock info for generate next block
	b = types.MockBlock()

	preBlockHeader := c.BestBlockHeader()
	preBlockHash := preBlockHeader.Hash()
	nextBlockHeight := preBlockHeader.Height + 1

	b.Previous = preBlockHash
	b.Height = nextBlockHeight
	b.Timestamp = uint64(time.Now().Unix())
	//nextBits, err := c.CalcNextBits(&preBlockHash)
	//if err != nil {
	//	return nil, err
	//}
	//
	//b = &types.Block{
	//	BlockHeader: types.BlockHeader{
	//		Version:           1,
	//		Height:            nextBlockHeight,
	//		PreviousBlockHash: preBlockHash,
	//		Timestamp:         uint64(time.Now().Unix()),
	//		BlockCommitment:   types.BlockCommitment{},
	//		Bits:              nextBits,
	//	},
	//}
	//bcBlock := &bc.Block{BlockHeader: &bc.BlockHeader{Height: nextBlockHeight}}
	//b.Transactions = []*types.Tx{nil}
	//
	//txs := txPool.GetTransactions()
	//sort.Sort(byTime(txs))
	//for _, txDesc := range txs {
	//	tx := txDesc.Tx.Tx
	//	gasOnlyTx := false
	//
	//	if err := c.GetTransactionsUtxo(view, []*bc.Tx{tx}); err != nil {
	//		blkGenSkipTxForErr(txPool, &tx.ID, err)
	//		continue
	//	}
	//
	//	gasStatus, err := validation.ValidateTx(tx, bcBlock)
	//	if err != nil {
	//		if !gasStatus.GasValid {
	//			blkGenSkipTxForErr(txPool, &tx.ID, err)
	//			continue
	//		}
	//		gasOnlyTx = true
	//	}
	//
	//	if gasUsed+uint64(gasStatus.GasUsed) > consensus.MaxBlockGas {
	//		break
	//	}
	//
	//	if err := view.ApplyTransaction(bcBlock, tx, gasOnlyTx); err != nil {
	//		blkGenSkipTxForErr(txPool, &tx.ID, err)
	//		continue
	//	}
	//
	//	if err := txStatus.SetStatus(len(b.Transactions), gasOnlyTx); err != nil {
	//		return nil, err
	//	}
	//
	//	b.Transactions = append(b.Transactions, txDesc.Tx)
	//	txEntries = append(txEntries, tx)
	//	gasUsed += uint64(gasStatus.GasUsed)
	//	txFee += txDesc.Fee
	//
	//	if gasUsed == consensus.MaxBlockGas {
	//		break
	//	}
	//}

	// creater coinbase transaction
	b.Transactions[0], err = createCoinbaseTx(0, nextBlockHeight)
	if err != nil {
		return nil, errors.Wrap(err, "fail on createCoinbaseTx")
	}

	b.TransactionRoot, err = types.TxMerkleRoot(b.Transactions)
	b.WitnessRoot = b.TransactionRoot

	return b, err
}

//func blkGenSkipTxForErr(txPool *protocol.TxPool, txHash *bc.Hash, err error) {
//	log.WithField("error", err).Error("mining block generation: skip tx due to")
//	txPool.RemoveTransaction(txHash)
//}
