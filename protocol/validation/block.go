package validation

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/clarenous/go-capsule/consensus"
	"github.com/clarenous/go-capsule/errors"
	"github.com/clarenous/go-capsule/protocol/types"

	"github.com/clarenous/go-capsule/protocol/state"
)

const logModule = "leveldb"

const BlockVersion = 1

var (
	errBadTimestamp = errors.New("block timestamp is not in the valid range")

	errMismatchedBlock       = errors.New("mismatched block")
	errMismatchedMerkleRoot  = errors.New("mismatched merkle root")
	errMisorderedBlockHeight = errors.New("misordered block height")
	errOverBlockLimit        = errors.New("block's gas is over the limit")

	errVersionRegression        = errors.New("version regression")
	errWrongCoinbaseTransaction = errors.New("wrong coinbase transaction")
)

// TODO: check overflow! (19.03.24 gcy)
func CheckCoinbaseAmount(b *types.Block, amount uint64) error {
	if len(b.Transactions) == 0 {
		return errors.Wrap(errWrongCoinbaseTransaction, "block is empty")
	}

	var totalOuts uint64
	for _, out := range b.Transactions[0].Outputs {
		totalOuts += out.Value
	}

	if totalOuts > amount {
		return errors.Wrap(errWrongCoinbaseTransaction, "reward more than deserved")
	}
	return nil
}

func checkBlockTime(b *types.Block, parent *state.BlockNode) error {
	if b.Timestamp > uint64(time.Now().Unix())+consensus.MaxTimeOffsetSeconds {
		return errBadTimestamp
	}

	if b.Timestamp <= parent.CalcPastMedianTime() {
		return errBadTimestamp
	}
	return nil
}

func ValidateProof(b *types.Block, parent *state.BlockNode) error {
	return b.Proof.ValidateProof([]interface{}{b, parent})
}

// ValidateBlockHeader check the block's header
func ValidateBlockHeader(b *types.Block, parent *state.BlockNode) error {
	if b.Version < BlockVersion {
		return errors.WithDetailf(errVersionRegression, "previous block version %d, current block version %d", parent.Version, b.Version)
	}
	if b.Height != parent.Height+1 {
		return errors.WithDetailf(errMisorderedBlockHeight, "previous block height %d, current block height %d", parent.Height, b.Height)
	}
	if parent.Hash != b.Previous {
		return errors.WithDetailf(errMismatchedBlock, "previous block ID %x, current block wants %x", parent.Hash.Bytes(), b.Previous.Bytes())
	}
	if err := checkBlockTime(b, parent); err != nil {
		return err
	}
	if err := ValidateProof(b, parent); err != nil {
		return err
	}
	return nil
}

// Store provides storage interface for blockchain data
type Store interface {
	GetTransaction(hash *types.Hash) (*types.Tx, error)
}

// ValidateBlock validates a block and the transactions within.
func ValidateBlock(store Store, b *types.Block, parent *state.BlockNode) error {
	startTime := time.Now()
	if err := ValidateBlockHeader(b, parent); err != nil {
		return err
	}

	// Check coinBase value
	var fee, totalIn, totalOut uint64
	for _, tx := range b.Transactions[1:] {
		for _, in := range tx.Inputs {
			sourceTx, err := store.GetTransaction(&in.ValueSource.TxID)
			if err != nil {
				return err
			}
			if int(in.ValueSource.Index) >= len(sourceTx.Outputs) {
				return errors.New("invalid output index")
			}
			totalIn += sourceTx.Outputs[in.ValueSource.Index].Value
		}
		for _, out := range tx.Outputs {
			totalOut += out.Value
		}
	}
	if totalOut < totalIn {
		return errors.New("invalid fee in block")
	}
	fee = totalOut - totalIn
	coinbaseAmount := consensus.BlockSubsidy(b.Height)
	if err := CheckCoinbaseAmount(b, coinbaseAmount+fee); err != nil {
		return err
	}

	for i, tx := range b.Transactions {
		if err := ValidateTx(tx, b); err != nil {
			return errors.Wrapf(err, "validate of transaction %d of %d", i, len(b.Transactions))
		}
	}

	txMerkleRoot, err := types.TxMerkleRoot(b.Transactions)
	if err != nil {
		return errors.Wrap(err, "computing transaction id merkle root")
	}
	if txMerkleRoot != b.TransactionRoot {
		return errors.WithDetailf(errMismatchedMerkleRoot, "transaction id merkle root")
	}

	txWitnessRoot, err := types.TxWitnessRoot(b.Transactions)
	if err != nil {
		return errors.Wrap(err, "computing transaction id merkle root")
	}
	if txWitnessRoot != b.WitnessRoot {
		return errors.WithDetailf(errMismatchedMerkleRoot, "witness id merkle root")
	}

	log.WithFields(log.Fields{
		"module":   logModule,
		"height":   b.Height,
		"hash":     b.Hash().String(),
		"duration": time.Since(startTime),
	}).Debug("finish validate block")
	return nil
}
