package validation

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/clarenous/go-capsule/consensus"
	"github.com/clarenous/go-capsule/consensus/algorithm/pow/difficulty"
	"github.com/clarenous/go-capsule/errors"
	"github.com/clarenous/go-capsule/protocol/types"

	"github.com/clarenous/go-capsule/protocol/state"
)

const logModule = "leveldb"

const BlockVersion = 1

var (
	errBadTimestamp          = errors.New("block timestamp is not in the valid range")
	errBadBits               = errors.New("block bits is invalid")
	errMismatchedBlock       = errors.New("mismatched block")
	errMismatchedMerkleRoot  = errors.New("mismatched merkle root")
	errMisorderedBlockHeight = errors.New("misordered block height")
	errOverBlockLimit        = errors.New("block's gas is over the limit")
	errWorkProof             = errors.New("invalid difficulty proof of work")
	errVersionRegression     = errors.New("version regression")
)

func checkBlockTime(b *types.Block, parent *state.BlockNode) error {
	if b.Timestamp > uint64(time.Now().Unix())+consensus.MaxTimeOffsetSeconds {
		return errBadTimestamp
	}

	if b.Timestamp <= parent.CalcPastMedianTime() {
		return errBadTimestamp
	}
	return nil
}

func checkCoinbaseAmount(b *types.Block, amount uint64) error {
	if len(b.Transactions) == 0 {
		return errors.Wrap(ErrWrongCoinbaseTransaction, "block is empty")
	}

	tx := b.Transactions[0]
	return nil
}

func ValidateProof(b *types.Block, parent *state.BlockNode) error {
	return b.Proof.ValidateProof(nil)

	if b.Proof.Target != parent.CalcNextBits() {
		return errBadBits
	}
	if !difficulty.CheckProofOfWork(b.Hash().Ptr(), b.Proof.Nonce) {
		return errWorkProof
	}
	return nil
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

// ValidateBlock validates a block and the transactions within.
func ValidateBlock(b *types.Block, parent *state.BlockNode) error {
	startTime := time.Now()
	if err := ValidateBlockHeader(b, parent); err != nil {
		return err
	}

	blockGasSum := uint64(0)
	coinbaseAmount := consensus.BlockSubsidy(b.BlockHeader.Height)

	for i, tx := range b.Transactions {
		// TODO:计算矿工收益
		err := ValidateTx(tx, b)
		if !gasStatus.GasValid {
			return errors.Wrapf(err, "validate of transaction %d of %d", i, len(b.Transactions))
		}

		coinbaseAmount += gasStatus.BTMValue
		if blockGasSum += uint64(gasStatus.GasUsed); blockGasSum > consensus.MaxBlockGas {
			return errOverBlockLimit
		}
	}

	if err := checkCoinbaseAmount(b, coinbaseAmount); err != nil {
		return err
	}

	txMerkleRoot, err := types.TxMerkleRoot(b.Transactions)
	if err != nil {
		return errors.Wrap(err, "computing transaction id merkle root")
	}
	if txMerkleRoot != *b.TransactionsRoot {
		return errors.WithDetailf(errMismatchedMerkleRoot, "transaction id merkle root")
	}

	log.WithFields(log.Fields{
		"module":   logModule,
		"height":   b.Height,
		"hash":     b.Hash().String(),
		"duration": time.Since(startTime),
	}).Debug("finish validate block")
	return nil
}
