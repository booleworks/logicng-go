package maxsat

import (
	"math"

	"booleworks.com/logicng/configuration"
)

// Algorithm encodes the different MAX-SAT algorithms.
type Algorithm byte

const (
	AlgWBO Algorithm = iota
	AlgIncWBO
	AlgLinearSU
	AlgLinearUS
	AlgMSU3
	AlgWMSU3
	AlgOLL
)

//go:generate stringer -type=Algorithm

// IncrementalStrategy encodes the different strategies for incremental encoding.
type IncrementalStrategy byte

const (
	IncNone IncrementalStrategy = iota
	IncIterative
)

//go:generate stringer -type=IncrementalStrategy

// WeightStrategy encodes the different strategies for handling weights.
type WeightStrategy byte

const (
	WeightNone WeightStrategy = iota
	WeightNormal
	WeightDiversify
)

//go:generate stringer -type=WeightStrategy

// Config describes the configuration of a MAX-SAT solver.  Incremental
// and weight strategy can be configured as well as flags for symmetry usage,
// and BMO as well as the symmetry limit.
type Config struct {
	IncrementalStrategy IncrementalStrategy
	WeightStrategy      WeightStrategy
	Symmetry            bool
	Limit               int
	BMO                 bool
}

// Sort returns the configuration sort (MaxSat).
func (Config) Sort() configuration.Sort {
	return configuration.MaxSat
}

// DefaultConfig returns the default configuration for a MAX-SAT
// configuration.
func (Config) DefaultConfig() configuration.Config {
	return DefaultConfig()
}

// DefaultConfig returns the default configuration for a
// MAX-SAT configuration.
func DefaultConfig() *Config {
	return &Config{
		IncrementalStrategy: IncNone,
		WeightStrategy:      WeightNone,
		Symmetry:            true,
		Limit:               math.MaxInt,
		BMO:                 true,
	}
}
