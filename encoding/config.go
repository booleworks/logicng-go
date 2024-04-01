package encoding

import "github.com/booleworks/logicng-go/configuration"

type (
	AMOEncoder        byte
	AMKEncoder        byte
	ALKEncoder        byte
	EXKEncoder        byte
	PBCEncoder        byte
	BimanderGroupSize byte
)

//go:generate stringer -type=AMOEncoder
//go:generate stringer -type=AMKEncoder
//go:generate stringer -type=ALKEncoder
//go:generate stringer -type=EXKEncoder
//go:generate stringer -type=PBCEncoder
//go:generate stringer -type=BimanderGroupSize

// CcAmoEncoder represents the different algorithms for encoding an at-most-one
// constraint (AMO) to a CNF.
const (
	AMOPure AMOEncoder = iota
	AMOLadder
	AMOProduct
	AMONested
	AMOCommander
	AMOBinary
	AMOBimander
	AMOBest
)

// CcAmkEncoder represents the different algorithms for encoding an at-most-k
// constraint (AMK) to a CNF.
const (
	AMKTotalizer AMKEncoder = iota
	AMKModularTotalizer
	AMKCardinalityNetwork
	AMKBest
)

// CcAlkEncoder represents the different algorithms for encoding an at-least-k
// constraint (ALK) to a CNF.
const (
	ALKTotalizer ALKEncoder = iota
	ALKModularTotalizer
	ALKCardinalityNetwork
	ALKBest
)

// CcExkEncoder represents the different algorithms for encoding an exactly-k
// constraint (EXK) to a CNF.
const (
	EXKTotalizer EXKEncoder = iota
	EXKCardinalityNetwork
	EXKBest
)

// BimanderGroupSize is a parameter which can be used for defining the group
// size in the Bimander encoding for AMO constraints.
const (
	BimanderHalf BimanderGroupSize = iota
	BimanderSqrt
	BimanderFixed
)

// PbcEncoder represents the different algorithms for encoding a pseudo-Boolean
// constraint to a CNF.
const (
	PBCSWC PBCEncoder = iota
	PBCBinaryMerge
	PBCAdderNetworks
	PBCBest
)

// Config describes the configuration for a cardinality or
// pseudo-Boolean encoding.  This configuration struct defines the encoding
// algorithms for each type of constraint and their respective configuration
// parameters.
type Config struct {
	AMOEncoder AMOEncoder
	AMKEncoder AMKEncoder
	ALKEncoder ALKEncoder
	EXKEncoder EXKEncoder
	PBCEncoder PBCEncoder

	BimanderGroupSize      BimanderGroupSize
	BimanderFixedGroupSize int
	NestingGroupSize       int
	ProductRecursiveBound  int
	CommanderGroupSize     int

	BinaryMergeUseGAC                bool
	BinaryMergeNoSupportForSingleBit bool
	BinaryMergeUseWatchDog           bool
}

// Sort returns the configuration sort (Encoder).
func (Config) Sort() configuration.Sort {
	return configuration.Encoder
}

// DefaultConfig returns the default configuration for an encoder
// configuration.
func (Config) DefaultConfig() configuration.Config {
	return DefaultConfig()
}

// DefaultConfig returns the default configuration for an
// encoder configuration.
func DefaultConfig() *Config {
	return &Config{
		AMOEncoder: AMOBest,
		AMKEncoder: AMKBest,
		ALKEncoder: ALKBest,
		EXKEncoder: EXKBest,
		PBCEncoder: PBCBest,

		BimanderGroupSize:      BimanderSqrt,
		BimanderFixedGroupSize: 3,
		NestingGroupSize:       4,
		ProductRecursiveBound:  20,
		CommanderGroupSize:     3,

		BinaryMergeUseGAC:                true,
		BinaryMergeNoSupportForSingleBit: false,
		BinaryMergeUseWatchDog:           true,
	}
}
