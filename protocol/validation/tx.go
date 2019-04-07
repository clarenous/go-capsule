package validation

import (
	"github.com/clarenous/go-capsule/errors"
	"github.com/clarenous/go-capsule/protocol/types"
)

// validate transaction error
var (
	ErrTxVersion                 = errors.New("invalid transaction version")
	ErrWrongTransactionSize      = errors.New("invalid transaction size")
	ErrBadLockTime               = errors.New("invalid transaction lock time")
	ErrEmptyInputIDs             = errors.New("got the empty InputIDs")
	ErrNotStandardTx             = errors.New("not standard transaction")
	ErrWrongCoinbaseAsset        = errors.New("wrong coinbase assetID")
	ErrCoinbaseArbitraryOversize = errors.New("coinbase arbitrary size is larger than limit")
	ErrEmptyScriptHash           = errors.New("transaction has output with empty script hash")
	ErrEmptyResults              = errors.New("transaction has no results")
	ErrMismatchedAssetID         = errors.New("mismatched assetID")
	ErrMismatchedPosition        = errors.New("mismatched value source/dest position")
	ErrMismatchedReference       = errors.New("mismatched reference")
	ErrMismatchedValue           = errors.New("mismatched value")
	ErrMissingField              = errors.New("missing required field")
	ErrNoSource                  = errors.New("no source for value")
	ErrOverflow                  = errors.New("arithmetic overflow/underflow")
	ErrPosition                  = errors.New("invalid source or destination position")
	ErrUnbalanced                = errors.New("unbalanced asset amount between input and output")
	ErrOverGasCredit             = errors.New("all gas credit has been spend")
	ErrGasCalculate              = errors.New("gas usage calculate got a math error")
)

// validationState contains the context that must propagate through
// the transaction graph when validating entries.
type validationState struct {
	block     *types.Block
	tx        *types.Tx
	entryID   types.Hash           // The ID of the nearest enclosing entry
	sourcePos uint64               // The source position, for validate ValueSources
	destPos   uint64               // The destination position, for validate ValueDestinations
	cache     map[types.Hash]error // Memoized per-entry validation results
}

func checkValidTx(vs *validationState, tx *types.Tx) (err error) {
	var ok bool
	txID := tx.Hash()
	if err, ok = vs.cache[txID]; ok {
		return err
	}

	defer func() {
		vs.cache[txID] = err
	}()

	// check tx version
	if err = checkValidTxVersion(vs, tx.Version); err != nil {
		return err
	}

	// check tx inputs
	for _, in := range tx.Inputs {
		if err = checkValidTxIn(vs, &in); err != nil {
			return err
		}
	}

	// check tx outputs
	for _, out := range tx.Outputs {
		if err = checkValidTxOut(vs, &out); err != nil {
			return err
		}
	}

	// currently not checking evidence

	return nil
}

func checkValidTxVersion(vs *validationState, version uint64) error {
	return nil
}

func checkValidTxIn(vs *validationState, in *types.TxIn) error {
	if in == nil {
		return errors.Wrap(ErrMissingField, "empty value txIn")
	}
	if in.RedeemScript == nil {
		return errors.Wrap(ErrMissingField, "missing redeem script in value txIn")
	}
	if in.UnlockScript == nil {
		return errors.Wrap(ErrMissingField, "missing unlock script in value txIn")
	}

	return nil
}

func checkValidTxOut(vs *validationState, out *types.TxOut) error {
	if out == nil {
		return errors.Wrap(ErrMissingField, "empty value txOut")
	}

	return nil
}

func checkStandardTx(tx *types.Tx) error {
	for _, in := range tx.Inputs {
		if in.ValueSource.TxID.IsZero() {
			return ErrEmptyInputIDs
		}
	}

	for _, out := range tx.Outputs {
		if out.ScriptHash.IsZero() {
			return ErrEmptyScriptHash
		}
	}
	return nil
}

func checkLockTime(tx *types.Tx, block *types.Block) error {
	if tx.LockTime == 0 {
		return nil
	}

	if tx.LockTime >= block.Height {
		return ErrBadLockTime
	}
	return nil
}

// ValidateTx validates a transaction.
func ValidateTx(tx *types.Tx, block *types.Block) error {
	var err error
	if tx.SerializedSize() == 0 {
		return ErrWrongTransactionSize
	}
	if err = checkLockTime(tx, block); err != nil {
		return err
	}
	if err = checkStandardTx(tx); err != nil {
		return err
	}

	vs := &validationState{
		block:   block,
		tx:      tx,
		entryID: tx.Hash(),
		cache:   make(map[types.Hash]error),
	}
	return checkValidTx(vs, tx)
}
