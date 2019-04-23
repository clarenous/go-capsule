package consensus

import (
	"github.com/clarenous/go-capsule/protocol/types"
	"strings"
)

//consensus variables
const (
	// parameters for transactions
	// MaxTxInSequenceNum is the maximum sequence number the sequence field
	// of a transaction input can be.
	MaxTxInSequenceNum uint32 = 0xffffffff
	// SequenceLockTimeDisabled is a flag that if set on a transaction
	// input's sequence number, the sequence number will not be interpreted
	// as a relative locktime.
	SequenceLockTimeDisabled = 1 << 31
	// SequenceLockTimeIsSeconds is a flag that if set on a transaction
	// input's sequence number, the relative locktime has units of 512
	// seconds.
	SequenceLockTimeIsSeconds = 1 << 22
	// SequenceLockTimeMask is a mask that extracts the relative locktime
	// when masked against the transaction input sequence number.
	SequenceLockTimeMask = 0x0000ffff

	// config parameter for coinbase reward
	CoinbasePendingBlockNumber = uint64(100)
	subsidyReductionInterval   = uint64(840000)
	baseSubsidy                = uint64(41250000000)
	InitialBlockSubsidy        = uint64(140700041250000000)

	// config for pow mining
	BlocksPerRetarget     = uint64(2016)
	TargetSecondsPerBlock = uint64(150)

	// MaxTimeOffsetSeconds is the maximum number of seconds a block time is allowed to be ahead of the current time
	MaxTimeOffsetSeconds = uint64(60 * 60)
	MedianTimeBlocks     = 11

	PayToWitnessPubKeyHashDataSize = 20
	PayToWitnessScriptHashDataSize = 32
	CoinbaseArbitrarySizeLimit     = 128

	ProofType = "pow"
)

// BlockSubsidy calculate the coinbase rewards on given block height
func BlockSubsidy(height uint64) uint64 {
	if height == 0 {
		return InitialBlockSubsidy
	}
	return baseSubsidy >> uint(height/subsidyReductionInterval)
}

// IsBech32SegwitPrefix returns whether the prefix is a known prefix for segwit
// addresses on any default or registered network.  This is used when decoding
// an address string into a specific address type.
func IsBech32SegwitPrefix(prefix string, params *Params) bool {
	prefix = strings.ToLower(prefix)
	return prefix == params.Bech32HRPSegwit+"1"
}

// Checkpoint identifies a known good point in the block chain.  Using
// checkpoints allows a few optimizations for old blocks during initial download
// and also prevents forks from old blocks.
type Checkpoint struct {
	Height uint64
	Hash   types.Hash
}

// Params store the config for different network
type Params struct {
	// Name defines a human-readable identifier for the network.
	Name            string
	Bech32HRPSegwit string
	// DefaultPort defines the default peer-to-peer port for the network.
	DefaultPort string

	// DNSSeeds defines a list of DNS seeds for the network that are used
	// as one method to discover peers.
	DNSSeeds    []string
	Checkpoints []Checkpoint
}

// ActiveNetParams is ...
var ActiveNetParams = MainNetParams

// NetParams is the correspondence between chain_id and Params
var NetParams = map[string]Params{
	"mainnet": MainNetParams,
	"solonet": SoloNetParams,
}

// MainNetParams is the config for production
var MainNetParams = Params{
	Name:            "main",
	Bech32HRPSegwit: "bm",
	DefaultPort:     "46657",
	DNSSeeds:        []string{},
	Checkpoints:     []Checkpoint{
		//{10000, types.NewHash([32]byte{0x93, 0xe1, 0xeb, 0x78, 0x21, 0xd2, 0xb4, 0xad, 0x0f, 0x5b, 0x1c, 0xea, 0x82, 0xe8, 0x43, 0xad, 0x8c, 0x09, 0x9a, 0xb6, 0x5d, 0x8f, 0x70, 0xc5, 0x84, 0xca, 0xa2, 0xdd, 0xf1, 0x74, 0x65, 0x2c})},
	},
}

// SoloNetParams is the config for test-net
var SoloNetParams = Params{
	Name:            "solo",
	Bech32HRPSegwit: "sm",
	Checkpoints:     []Checkpoint{},
}
