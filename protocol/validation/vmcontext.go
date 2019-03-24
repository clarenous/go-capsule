package validation

import (
	"bytes"

	"github.com/clarenous/go-capsule/consensus/segwit"
	"github.com/clarenous/go-capsule/crypto/sha3pool"
	"github.com/clarenous/go-capsule/errors"
	"github.com/clarenous/go-capsule/protocol/types"
	"github.com/clarenous/go-capsule/protocol/vm"
)

// NewTxVMContext generates the vm.Context for BVM
func NewTxVMContext(vs *validationState, entry types.Entry, prog *types.Program, args [][]byte) *vm.Context {
	var (
		tx          = vs.tx
		blockHeight = vs.block.BlockHeader.GetHeight()
		numResults  = uint64(len(tx.ResultIds))
		entryID     = types.EntryID(entry) // TODO(bobg): pass this in, don't recompute it

		assetID       *[]byte
		amount        *uint64
		destPos       *uint64
		spentOutputID *[]byte
	)

	switch e := entry.(type) {
	case *types.Issuance:
		a1 := e.Value.AssetId.Bytes()
		assetID = &a1
		amount = &e.Value.Amount
		destPos = &e.WitnessDestination.Position

	case *types.Spend:
		spentOutput := tx.Entries[*e.SpentOutputId].(*types.Output)
		a1 := spentOutput.Source.Value.AssetId.Bytes()
		assetID = &a1
		amount = &spentOutput.Source.Value.Amount
		destPos = &e.WitnessDestination.Position
		s := e.SpentOutputId.Bytes()
		spentOutputID = &s
	}

	var txSigHash *[]byte
	txSigHashFn := func() []byte {
		if txSigHash == nil {
			hasher := sha3pool.Get256()
			defer sha3pool.Put256(hasher)

			entryID.WriteTo(hasher)
			tx.ID.WriteTo(hasher)

			var hash types.Hash
			hash.ReadFrom(hasher)
			hashBytes := hash.Bytes()
			txSigHash = &hashBytes
		}
		return *txSigHash
	}

	ec := &entryContext{
		entry:   entry,
		entries: tx.Entries,
	}

	result := &vm.Context{
		VMVersion: prog.VmVersion,
		Code:      witnessProgram(prog.Code),
		Arguments: args,

		EntryID: entryID.Bytes(),

		TxVersion:   &tx.Version,
		BlockHeight: &blockHeight,

		TxSigHash:     txSigHashFn,
		NumResults:    &numResults,
		AssetID:       assetID,
		Amount:        amount,
		DestPos:       destPos,
		SpentOutputID: spentOutputID,
		CheckOutput:   ec.checkOutput,
	}

	return result
}

func witnessProgram(prog []byte) []byte {
	if segwit.IsP2WPKHScript(prog) {
		if witnessProg, err := segwit.ConvertP2PKHSigProgram([]byte(prog)); err == nil {
			return witnessProg
		}
	} else if segwit.IsP2WSHScript(prog) {
		if witnessProg, err := segwit.ConvertP2SHProgram([]byte(prog)); err == nil {
			return witnessProg
		}
	}
	return prog
}

type entryContext struct {
	entry   types.Entry
	entries map[types.Hash]types.Entry
}

func (ec *entryContext) checkOutput(index uint64, amount uint64, assetID []byte, vmVersion uint64, code []byte, expansion bool) (bool, error) {
	checkEntry := func(e types.Entry) (bool, error) {
		check := func(prog *types.Program, value *types.AssetAmount) bool {
			return (prog.VmVersion == vmVersion &&
				bytes.Equal(prog.Code, code) &&
				bytes.Equal(value.AssetId.Bytes(), assetID) &&
				value.Amount == amount)
		}

		switch e := e.(type) {
		case *types.Output:
			return check(e.ControlProgram, e.Source.Value), nil

		case *types.Retirement:
			var prog types.Program
			if expansion {
				// The spec requires prog.Code to be the empty string only
				// when !expansion. When expansion is true, we prepopulate
				// prog.Code to give check() a freebie match.
				//
				// (The spec always requires prog.VmVersion to be zero.)
				prog.Code = code
			}
			return check(&prog, e.Source.Value), nil
		}

		return false, vm.ErrContext
	}

	checkMux := func(m *types.Mux) (bool, error) {
		if index >= uint64(len(m.WitnessDestinations)) {
			return false, errors.Wrapf(vm.ErrBadValue, "index %d >= %d", index, len(m.WitnessDestinations))
		}
		eID := m.WitnessDestinations[index].Ref
		e, ok := ec.entries[*eID]
		if !ok {
			return false, errors.Wrapf(types.ErrMissingEntry, "entry for mux destination %d, id %x, not found", index, eID.Bytes())
		}
		return checkEntry(e)
	}

	switch e := ec.entry.(type) {
	case *types.Mux:
		return checkMux(e)

	case *types.Issuance:
		d, ok := ec.entries[*e.WitnessDestination.Ref]
		if !ok {
			return false, errors.Wrapf(types.ErrMissingEntry, "entry for issuance destination %x not found", e.WitnessDestination.Ref.Bytes())
		}
		if m, ok := d.(*types.Mux); ok {
			return checkMux(m)
		}
		if index != 0 {
			return false, errors.Wrapf(vm.ErrBadValue, "index %d >= 1", index)
		}
		return checkEntry(d)

	case *types.Spend:
		d, ok := ec.entries[*e.WitnessDestination.Ref]
		if !ok {
			return false, errors.Wrapf(types.ErrMissingEntry, "entry for spend destination %x not found", e.WitnessDestination.Ref.Bytes())
		}
		if m, ok := d.(*types.Mux); ok {
			return checkMux(m)
		}
		if index != 0 {
			return false, errors.Wrapf(vm.ErrBadValue, "index %d >= 1", index)
		}
		return checkEntry(d)
	}

	return false, vm.ErrContext
}
