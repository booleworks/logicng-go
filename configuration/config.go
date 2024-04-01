package configuration

// Sort describes the different sorts for configurations.
type Sort byte

const (
	FormulaFactory Sort = iota
	CNF
	Sat
	MaxSat
	Encoder
	FormulaRandomizer
	AdvancedSimplifier
	ModelIteration
)

//go:generate stringer -type=Sort

// Config is an abstraction over all LogicNG configuration structs.
type Config interface {
	Sort() Sort            // returns the sort of the configuration
	DefaultConfig() Config // creates the default configuration for the sort
}
