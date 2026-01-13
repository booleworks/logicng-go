package sat

import "github.com/booleworks/logicng-go/configuration"

// ClauseMinimization encodes the different algorithms for minimizing learnt
// clauses on the solver.
type ClauseMinimization byte

const (
	ClauseMinNone  ClauseMinimization = iota // no clause minimization
	ClauseMinBasic                           // simple minimization
	ClauseMinDeep                            // recursive deep minimization
)

//go:generate stringer -type=ClauseMinimization

// CNFMethod encodes the different methods for adding formulas as CNF to
// the solver.
type CNFMethod byte

const (
	CNFFactory CNFMethod = iota // formula factories CNF method
	CNFPG                       // Plaisted-Greenbaum with NNF generation directly on solver
	CNFFullPG                   // Plaisted-Greenbaum without NNF generation directly on solver
)

//go:generate stringer -type=CNFMethod

// Config describes the configuration of a SAT solver.
type Config struct {
	ProofGeneration    bool               // record proof generation information on-the-fly
	UseAtMostClauses   bool               // use a special representation of at-most-one clauses
	CNFMethod          CNFMethod          // method for adding CNFs
	ClauseMinimization ClauseMinimization // algorithm for minimizing learnt clauses
	InitialPhase       bool               // initial phase for assigning literals
	LowLevelConfig     *LowLevelConfig    // low level config
}

// Sort returns the configuration sort (Sat).
func (c *Config) Sort() configuration.Sort {
	return configuration.Sat
}

// DefaultConfig returns the default configuration for a SAT solver
// configuration.
func (c *Config) DefaultConfig() configuration.Config {
	return DefaultConfig()
}

// CNF sets the CNF method on this configuration and returns the config.
func (c *Config) CNF(cnfMethod CNFMethod) *Config {
	c.CNFMethod = cnfMethod
	return c
}

// ClauseMin sets the clause minimization method on this configuration and
// returns the config.
func (c *Config) ClauseMin(clauseMin ClauseMinimization) *Config {
	c.ClauseMinimization = clauseMin
	return c
}

// InitPhase sets the initial phase on this configuration and returns the
// config.
func (c *Config) InitPhase(initPhase bool) *Config {
	c.InitialPhase = initPhase
	return c
}

// Proofs sets whether proofs should be generated and returns the config.
func (c *Config) Proofs(proofs bool) *Config {
	c.ProofGeneration = proofs
	return c
}

// UseAtMost sets the flag whether at-most clauses should be used and returns
// the config.
func (c *Config) UseAtMost(useAtMost bool) *Config {
	c.UseAtMostClauses = useAtMost
	return c
}

// DefaultConfig returns the default configuration for a SAT solver
// configuration.
func DefaultConfig() *Config {
	return &Config{
		ProofGeneration:    false,
		UseAtMostClauses:   false,
		CNFMethod:          CNFPG,
		ClauseMinimization: ClauseMinDeep,
		InitialPhase:       false,
		LowLevelConfig:     DefaultLowLevelConfig(),
	}
}

// LowLevelConfig describes the low-level parameters of the SAT solver.
// Usually there is no need to manually change these parameters.
type LowLevelConfig struct {
	VarDecay         float64
	VarInc           float64
	RestartFirst     int
	RestartInc       float64
	ClauseDecay      float64
	LearntsizeFactor float64
	LearntsizeInc    float64

	LBLBDMinimizingClause  int
	LBLBDFrozenClause      int
	LBSizeMinimizingClause int
	FirstReduceDB          int
	SpecialIncReduceDB     int
	IncReduceDB            int
	FactorK                float64
	FactorR                float64
	SizeLBDQueue           int
	SizeTrailQueue         int
	ReduceOnSize           bool
	ReduceOnSizeSize       int
	MaxVarDecay            float64
}

// DefaultLowLevelConfig returns a new default configuration of the low-level
// parameters of the SAT solver.
func DefaultLowLevelConfig() *LowLevelConfig {
	return &LowLevelConfig{
		VarDecay:               0.95,
		VarInc:                 1.0,
		RestartFirst:           100,
		RestartInc:             2.0,
		ClauseDecay:            0.999,
		LearntsizeFactor:       1.0 / 3.0,
		LearntsizeInc:          1.1,
		LBLBDMinimizingClause:  6,
		LBLBDFrozenClause:      30,
		LBSizeMinimizingClause: 30,
		FirstReduceDB:          2000,
		SpecialIncReduceDB:     1000,
		IncReduceDB:            300,
		FactorK:                0.8,
		FactorR:                1.4,
		SizeLBDQueue:           50,
		SizeTrailQueue:         5000,
		ReduceOnSize:           false,
		ReduceOnSizeSize:       12,
		MaxVarDecay:            0.95,
	}
}
