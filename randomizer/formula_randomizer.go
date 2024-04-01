package randomizer

import (
	"fmt"
	"math"
	"math/rand"

	"booleworks.com/logicng/configuration"
	f "booleworks.com/logicng/formula"
)

// Config describes the configuration of a formula randomizer with
// weights for all different formula types.
type Config struct {
	Seed           int64
	Variables      []f.Variable
	NumVars        int
	WeightConstant float64
	WeightPosLit   float64
	WeightNegLit   float64
	WeightOr       float64
	WeightAnd      float64
	WeightNot      float64
	WeightImpl     float64
	WeightEquiv    float64
	MaximumOpsAnd  int
	MaximumOpsOr   int

	WeightPBC         float64
	WeightPBCCoeffPos float64
	WeightPBCCoeffNeg float64
	WeightPBCLE       float64
	WeightPBCLT       float64
	WeightPBCGE       float64
	WeightPBCGT       float64
	WeightPBCEQ       float64
	MaximumOpsPBC     int
	MaximumCoeffPBC   int

	WeightCC     float64
	WeightAMO    float64
	WeightEXO    float64
	MaximumOpsCC int
}

// Sort returns the configuration sort (Randomizer).
func (Config) Sort() configuration.Sort {
	return configuration.FormulaRandomizer
}

// DefaultConfig returns the default configuration for a randomizer
// configuration.
func (Config) DefaultConfig() configuration.Config {
	return DefaultConfig()
}

// DefaultConfig returns the default configuration for a randomizer
// configuration.
func DefaultConfig() *Config {
	return &Config{
		Seed:              0,
		Variables:         nil,
		NumVars:           25,
		WeightConstant:    0.1,
		WeightPosLit:      1.0,
		WeightNegLit:      1.0,
		WeightOr:          30.0,
		WeightAnd:         30.0,
		WeightNot:         1.0,
		WeightImpl:        1.0,
		WeightEquiv:       1.0,
		MaximumOpsAnd:     5,
		MaximumOpsOr:      5,
		WeightPBC:         0.0,
		WeightPBCCoeffPos: 1.0,
		WeightPBCCoeffNeg: 0.2,
		WeightPBCLE:       0.2,
		WeightPBCLT:       0.2,
		WeightPBCGE:       0.2,
		WeightPBCGT:       0.2,
		WeightPBCEQ:       0.2,
		MaximumOpsPBC:     5,
		MaximumCoeffPBC:   10,
		WeightCC:          0.0,
		WeightAMO:         0.0,
		WeightEXO:         0.0,
		MaximumOpsCC:      5,
	}
}

// A FormulaRandomizer can be used to generate randomized formulas.
type FormulaRandomizer struct {
	fac                      f.Factory
	config                   Config
	random                   rand.Rand
	variables                []f.Variable
	fTypeProbabilities       fTypeProbabilities
	cTypeProbabilities       cTypeProbabilities
	phaseProbability         float64
	coeffNegativeProbability float64
}

// New generates a new formula randomizer with the given
// formula factory and an optional randomizer configuration.
func New(fac f.Factory, config ...*Config) *FormulaRandomizer {
	return newFormulaRandomizer(fac, determineConfig(fac, config))
}

// NewWithSeed generates a new formula randomizer with the
// given formula factory and the default randomizer configuration with the
// given seed.
func NewWithSeed(fac f.Factory, seed int64) *FormulaRandomizer {
	config := DefaultConfig()
	config.Seed = seed
	return newFormulaRandomizer(fac, config)
}

func newFormulaRandomizer(fac f.Factory, config *Config) *FormulaRandomizer {
	return &FormulaRandomizer{
		fac:                fac,
		config:             *config,
		random:             *rand.New(rand.NewSource(config.Seed)),
		variables:          generateVars(fac, config),
		fTypeProbabilities: newFTypeProbabilities(config),
		cTypeProbabilities: newCTypeProbabilities(config),
		phaseProbability: config.WeightPosLit /
			(config.WeightPosLit + config.WeightNegLit),
		coeffNegativeProbability: config.WeightPBCCoeffNeg /
			(config.WeightPBCCoeffPos + config.WeightPBCCoeffNeg),
	}
}

func determineConfig(fac f.Factory, initConfig []*Config) *Config {
	if len(initConfig) > 0 {
		return initConfig[0]
	} else {
		configFromFactory, ok := fac.ConfigurationFor(configuration.FormulaRandomizer)
		if !ok {
			return DefaultConfig()
		} else {
			return configFromFactory.(*Config)
		}
	}
}

// Constant returns a random constant.
func (r *FormulaRandomizer) Constant() f.Formula {
	return r.fac.Constant(r.random.Intn(2) > 0)
}

// Variable returns a random variable.
func (r *FormulaRandomizer) Variable() f.Variable {
	return r.variables[r.random.Intn(len(r.variables))]
}

// Literal returns a random literal. The probability of whether it is positive
// or negative depends on the configuration.
func (r *FormulaRandomizer) Literal() f.Literal {
	randVar := r.variables[r.random.Intn(len(r.variables))]
	name, _ := r.fac.VarName(randVar)
	return r.fac.Lit(name, r.random.Float64() < r.phaseProbability)
}

// Atom returns a random atom. This includes constants, literals, pseudo
// boolean constraints, and cardinality constraints (including AMO and EXO).
func (r *FormulaRandomizer) Atom() f.Formula {
	n := r.random.Float64() * r.fTypeProbabilities.exo
	switch {
	case n < r.fTypeProbabilities.constant:
		return r.Constant()
	case n < r.fTypeProbabilities.literal:
		return r.Literal().AsFormula()
	case n < r.fTypeProbabilities.pbc:
		return r.PBC()
	case n < r.fTypeProbabilities.cc:
		return r.CC()
	case n < r.fTypeProbabilities.amo:
		return r.AMO()
	default:
		return r.EXO()
	}
}

// Not returns a random negation with a given maximal depth.
func (r *FormulaRandomizer) Not(maxDepth int) f.Formula {
	if maxDepth == 0 {
		return r.Atom()
	}
	not := r.fac.Not(r.Formula(maxDepth - 1))
	if maxDepth >= 2 && not.Sort() != f.SortNot {
		return r.Not(maxDepth)
	}
	return not
}

// Impl returns a random implication with a given maximal depth.
func (r *FormulaRandomizer) Impl(maxDepth int) f.Formula {
	if maxDepth == 0 {
		return r.Atom()
	}
	implication := r.fac.Implication(r.Formula(maxDepth-1), r.Formula(maxDepth-1))
	if implication.Sort() != f.SortImpl {
		return r.Impl(maxDepth)
	}
	return implication
}

// Equiv returns a random equivalence with a given maximal depth.
func (r *FormulaRandomizer) Equiv(maxDepth int) f.Formula {
	if maxDepth == 0 {
		return r.Atom()
	}
	equiv := r.fac.Equivalence(r.Formula(maxDepth-1), r.Formula(maxDepth-1))
	if equiv.Sort() != f.SortEquiv {
		return r.Equiv(maxDepth)
	}
	return equiv
}

// And returns a random conjunction with a given maximal depth.
func (r *FormulaRandomizer) And(maxDepth int) f.Formula {
	if maxDepth == 0 {
		return r.Atom()
	}
	operands := make([]f.Formula, 2+r.random.Intn(r.config.MaximumOpsAnd-2))
	for i := 0; i < len(operands); i++ {
		operands[i] = r.Formula(maxDepth - 1)
	}
	formula := r.fac.And(operands...)
	if formula.Sort() != f.SortAnd {
		return r.And(maxDepth)
	}
	return formula
}

// Or returns a random disjunction with a given maximal depth.
func (r *FormulaRandomizer) Or(maxDepth int) f.Formula {
	if maxDepth == 0 {
		return r.Atom()
	}
	operands := make([]f.Formula, 2+r.random.Intn(r.config.MaximumOpsOr-2))
	for i := 0; i < len(operands); i++ {
		operands[i] = r.Formula(maxDepth - 1)
	}
	formula := r.fac.Or(operands...)
	if formula.Sort() != f.SortOr {
		return r.Or(maxDepth)
	}
	return formula
}

// CC returns a random cardinality constraint.
func (r *FormulaRandomizer) CC() f.Formula {
	variables := r.ccVariables()
	sort := r.cSort()
	rhsBound := len(variables)
	switch sort {
	case f.GT:
		rhsBound = len(variables) + 1
	case f.LT:
		rhsBound = len(variables) + 1
	}
	rhsOffset := 0
	switch sort {
	case f.GT:
		rhsOffset = -1
	case f.LT:
		rhsOffset = 1
	}
	if rhsBound == 0 {
		return r.CC()
	}
	rhs := uint32(max(0, rhsOffset+r.random.Intn(rhsBound)))
	cc := r.fac.CC(sort, rhs, variables...)
	if cc.Sort() <= f.SortTrue {
		return r.CC()
	}
	return cc
}

// AMO returns a random at-most-one constraint.
func (r *FormulaRandomizer) AMO() f.Formula {
	return r.fac.AMO(r.ccVariables()...)
}

// EXO returns a random exactly-one constraint.
func (r *FormulaRandomizer) EXO() f.Formula {
	return r.fac.EXO(r.ccVariables()...)
}

// PBC returns a random pseudo boolean constraint.
func (r *FormulaRandomizer) PBC() f.Formula {
	numOps := r.random.Intn(r.config.MaximumOpsPBC)
	literals := make([]f.Literal, numOps)
	coefficients := make([]int, numOps)
	minSum := 0 // (positive) sum of all negative coefficients
	maxSum := 0 // sum of all positive coefficients
	for i := 0; i < numOps; i++ {
		literals[i] = r.Literal()
		coefficients[i] = r.random.Intn(r.config.MaximumCoeffPBC) + 1
		if r.random.Float64() < r.coeffNegativeProbability {
			minSum += coefficients[i]
			coefficients[i] = -coefficients[i]
		} else {
			maxSum += coefficients[i]
		}
	}
	sort := r.cSort()
	rhs := r.random.Intn(maxSum+minSum+1) - minSum
	pbc := r.fac.PBC(sort, rhs, literals, coefficients)
	if pbc.Sort() <= f.SortTrue {
		return r.PBC()
	}
	return pbc
}

// Formula returns a random formula with a given maximal depth.
func (r *FormulaRandomizer) Formula(maxDepth int) f.Formula {
	if maxDepth == 0 {
		return r.Atom()
	} else {
		n := r.random.Float64()
		switch {
		case n < r.fTypeProbabilities.constant:
			return r.Constant()
		case n < r.fTypeProbabilities.literal:
			return r.Literal().AsFormula()
		case n < r.fTypeProbabilities.pbc:
			return r.PBC()
		case n < r.fTypeProbabilities.cc:
			return r.CC()
		case n < r.fTypeProbabilities.amo:
			return r.AMO()
		case n < r.fTypeProbabilities.exo:
			return r.EXO()
		case n < r.fTypeProbabilities.or:
			return r.Or(maxDepth)
		case n < r.fTypeProbabilities.and:
			return r.And(maxDepth)
		case n < r.fTypeProbabilities.not:
			return r.Not(maxDepth)
		case n < r.fTypeProbabilities.impl:
			return r.Impl(maxDepth)
		default:
			return r.Equiv(maxDepth)
		}
	}
}

// ConstraintSet returns a list of numConstraints random formula with a given
// maximal depth.
func (r *FormulaRandomizer) ConstraintSet(numConstraints, maxDepth int) []f.Formula {
	formulas := make([]f.Formula, numConstraints)
	for i := 0; i < len(formulas); i++ {
		formulas[i] = r.Formula(maxDepth)
	}
	return formulas
}

func (r *FormulaRandomizer) ccVariables() []f.Variable {
	variables := make(map[f.Variable]present)
	bound := r.random.Intn(r.config.MaximumOpsCC-1) + 2
	for i := 0; i < bound; i++ {
		variables[r.Variable()] = present{}
	}
	vars := make([]f.Variable, 0, len(variables))
	for v := range variables {
		vars = append(vars, v)
	}
	return vars
}

func generateVars(fac f.Factory, config *Config) []f.Variable {
	if config.Variables != nil {
		return config.Variables
	} else {
		variables := make([]f.Variable, config.NumVars)
		decimalPlaces := int(math.Ceil(math.Log10(float64(config.NumVars))))
		formatter := fmt.Sprintf("v%s%dd", "%0", decimalPlaces)
		for i := 0; i < len(variables); i++ {
			variables[i] = fac.Var(fmt.Sprintf(formatter, i))
		}
		return variables
	}
}

func (r *FormulaRandomizer) cSort() f.CSort {
	var sort f.CSort
	n := r.random.Float64()
	switch {
	case n < r.cTypeProbabilities.le:
		sort = f.LE
	case n < r.cTypeProbabilities.lt:
		sort = f.LT
	case n < r.cTypeProbabilities.ge:
		sort = f.GE
	case n < r.cTypeProbabilities.gt:
		sort = f.GT
	default:
		sort = f.EQ
	}
	return sort
}

type fTypeProbabilities struct {
	constant float64
	literal  float64
	pbc      float64
	cc       float64
	amo      float64
	exo      float64
	or       float64
	and      float64
	not      float64
	impl     float64
	equiv    float64
}

func newFTypeProbabilities(config *Config) fTypeProbabilities {
	total := config.WeightConstant + config.WeightPosLit + config.WeightNegLit +
		config.WeightOr + config.WeightAnd + config.WeightNot + config.WeightImpl + config.WeightEquiv +
		config.WeightPBC + config.WeightCC + config.WeightAMO + config.WeightEXO
	constant := config.WeightConstant / total
	literal := constant + (config.WeightPosLit+config.WeightNegLit)/total
	pbc := literal + config.WeightPBC/total
	cc := pbc + config.WeightCC/total
	amo := cc + config.WeightAMO/total
	exo := amo + config.WeightEXO/total
	or := exo + config.WeightOr/total
	and := or + config.WeightAnd/total
	not := and + config.WeightNot/total
	impl := not + config.WeightImpl/total
	equiv := impl + config.WeightEquiv/total
	return fTypeProbabilities{constant, literal, pbc, cc, amo, exo, or, and, not, impl, equiv}
}

type cTypeProbabilities struct {
	le float64
	lt float64
	ge float64
	gt float64
	eq float64
}

func newCTypeProbabilities(config *Config) cTypeProbabilities {
	total := config.WeightPBCLE + config.WeightPBCLT + config.WeightPBCGE +
		config.WeightPBCGT + config.WeightPBCEQ
	le := config.WeightPBCLE / total
	lt := le + config.WeightPBCLT/total
	ge := lt + config.WeightPBCGE/total
	gt := ge + config.WeightPBCGT/total
	eq := gt + config.WeightPBCEQ/total
	return cTypeProbabilities{le, lt, ge, gt, eq}
}

type present struct{}
