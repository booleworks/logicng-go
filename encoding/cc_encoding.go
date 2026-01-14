package encoding

import (
	"github.com/booleworks/logicng-go/configuration"
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// EncodeCC encodes a cardinality constraint to a CNF formula as a list of
// clauses.
//
// Depending on the type of cardinality constraint it uses the encoding
// algorithm as specified either in the optional given config or if it not
// present the one configured in the formula factory.
//
// Returns an error if the input constraint is no valid cardinality constraint.
func EncodeCC(fac f.Factory, constraint f.Formula, config ...*Config) ([]f.Formula, error) {
	result := ResultForFormula(fac)
	err := EncodeCCInResult(fac, constraint, result, config...)
	return result.result, err
}

// EncodeCCInResult encodes a cardinality constraint into an encoding result.
// This result can be either a formula result and a formula is generated or it
// can be a solver encoding and the cnf is added directly to the solver without
// first generating formulas on the formula factory.
//
// Depending on the type of cardinality constraint it uses the encoding
// algorithm as specified either in the optional given config or if it not
// present the one configured in the formula factory.
//
// Returns an error if the input constraint is no valid cardinality constraint.
func EncodeCCInResult(
	fac f.Factory,
	constraint f.Formula,
	result Result,
	config ...*Config,
) error {
	if constraint.Sort() != f.SortCC {
		return errorx.BadFormulaSort(constraint.Sort())
	}
	cfg := determineConfig(fac, config)
	comparator, rhs, lits, _, found := fac.PBCOps(constraint)
	if !found {
		panic(errorx.UnknownFormula(constraint))
	}
	ops, _ := f.LiteralsAsVariables(lits)
	switch comparator {
	case f.LE:
		if rhs == 1 {
			amo(result, cfg, ops)
		} else {
			amk(result, cfg, rhs, ops)
		}
	case f.LT:
		if rhs == 2 {
			amo(result, cfg, ops)
		} else {
			amk(result, cfg, rhs-1, ops)
		}
	case f.GE:
		alk(result, cfg, rhs, ops)
	case f.GT:
		alk(result, cfg, rhs+1, ops)
	case f.EQ:
		if rhs == 1 {
			exo(result, cfg, ops)
		} else {
			exk(result, cfg, rhs, ops)
		}
	default:
		panic(errorx.UnknownEnumValue(comparator))
	}
	return nil
}

// EncodeIncremental encodes an incremental cardinalityConstraint into an
// encoding result and return the incremental data with this result. This
// result can be either a formula result and a formula is generated or it can
// be a solver encoding and the cnf is added directly to the solver without
// first generating formulas on the formula factory.  The returned incremental
// data can then be used to tighten the bound of the constraint (either as
// formula or also directly on the solver).
//
// Depending on the type of cardinality constraint it uses the encoding
// algorithm as specified either in the optional given config or if it not
// present the one configured in the formula factory.
//
// Returns an error if the input constraint cannot be converted to an
// incremental constraint which is the case if it is an exo-constraint, a
// tautology or contradiction, or a trivial constraint (all literals are set
// to true or false)
func EncodeIncremental(
	fac f.Factory,
	constraint f.Formula,
	result Result,
	config ...*Config,
) (*CCIncrementalData, error) {
	cfg := determineConfig(fac, config)
	comparator, rhs, lits, _, found := fac.PBCOps(constraint)
	if !found {
		panic(errorx.UnknownFormula(constraint))
	}
	ops, _ := f.LiteralsAsVariables(lits)
	switch comparator {
	case f.LE:
		return amkIncremental(result, cfg, ops, rhs)
	case f.LT:
		return amkIncremental(result, cfg, ops, rhs-1)
	case f.GE:
		return alkIncremental(result, cfg, ops, rhs)
	case f.GT:
		return alkIncremental(result, cfg, ops, rhs+1)
	default:
		panic(errorx.BadInput("incremental encodings are only supported for at-most-k and at-least k constraints"))
	}
}

func determineConfig(fac f.Factory, initConfig []*Config) *Config {
	if len(initConfig) > 0 {
		return initConfig[0]
	}
	configFromFactory, ok := fac.ConfigurationFor(configuration.Encoder)
	if !ok {
		return DefaultConfig()
	}
	return configFromFactory.(*Config)
}

func amo(result Result, config *Config, vars []f.Variable) {
	if len(vars) <= 1 {
		return
	}
	switch config.AMOEncoder {
	case AMOPure:
		amoPure(result, vars)
	case AMOLadder:
		amoLadder(result, vars)
	case AMOProduct:
		amoProduct(result, config.ProductRecursiveBound, vars)
	case AMONested:
		amoNested(result, config.NestingGroupSize, f.VariablesAsLiterals(vars))
	case AMOCommander:
		amoCommander(result, config.CommanderGroupSize, f.VariablesAsLiterals(vars))
	case AMOBinary:
		amoBinary(result, vars)
	case AMOBimander:
		amoBimander(result, config.BimanderGroupSize, config.BimanderFixedGroupSize, vars)
	case AMOBest:
		bestAMO(result, config, vars)
	default:
		panic(errorx.UnknownEnumValue(config.AMOEncoder))
	}
}

func bestAMO(result Result, config *Config, vars []f.Variable) {
	if len(vars) <= 10 {
		amoPure(result, vars)
	} else {
		amoProduct(result, config.ProductRecursiveBound, vars)
	}
}

func exo(result Result, config *Config, vars []f.Variable) {
	if len(vars) == 0 {
		result.AddClause()
		return
	}
	if len(vars) == 1 {
		result.AddClause(vars[0].AsLiteral())
		return
	}
	amo(result, config, vars)
	result.AddClause(f.VariablesAsLiterals(vars)...)
}

func amk(result Result, config *Config, rhs int, vars []f.Variable) {
	if rhs >= len(vars) { // there is no constraint
		return
	}
	if rhs == 0 { // no variable can be true
		for _, v := range vars {
			result.AddClause(v.Negate(result.Factory()))
		}
		return
	}
	switch config.AMKEncoder {
	case AMKTotalizer:
		totalizerAMK(result, vars, rhs)
	case AMKModularTotalizer:
		modtotalizerAMK(result, f.VariablesAsLiterals(vars), rhs)
	case AMKCardinalityNetwork:
		cnAmk(result, vars, rhs)
	case AMKBest:
		modtotalizerAMK(result, f.VariablesAsLiterals(vars), rhs)
	default:
		panic(errorx.UnknownEnumValue(config.AMKEncoder))
	}
}

func alk(result Result, config *Config, rhs int, vars []f.Variable) {
	if rhs > len(vars) {
		result.AddClause()
		return
	}
	if rhs == 0 {
		return
	}
	if rhs == 1 {
		result.AddClause(f.VariablesAsLiterals(vars)...)
		return
	}
	if rhs == len(vars) {
		for _, v := range vars {
			result.AddClause(v.AsLiteral())
		}
		return
	}
	switch config.ALKEncoder {
	case ALKTotalizer:
		totalizerALK(result, vars, rhs)
	case ALKModularTotalizer:
		modtotalizerALK(result, f.VariablesAsLiterals(vars), rhs)
	case ALKCardinalityNetwork:
		cnAlk(result, vars, rhs)
	case ALKBest:
		modtotalizerALK(result, f.VariablesAsLiterals(vars), rhs)
	default:
		panic(errorx.UnknownEnumValue(config.ALKEncoder))
	}
}

func exk(result Result, config *Config, rhs int, vars []f.Variable) {
	if rhs > len(vars) {
		result.AddClause()
		return
	}
	if rhs == 0 {
		for _, v := range vars {
			result.AddClause(v.Negate(result.Factory()))
		}
		return
	}
	if rhs == len(vars) {
		for _, v := range vars {
			result.AddClause(v.AsLiteral())
		}
		return
	}
	switch config.EXKEncoder {
	case EXKTotalizer:
		totalizerEXK(result, vars, rhs)
	case EXKCardinalityNetwork:
		cnExk(result, vars, rhs)
	case EXKBest:
		totalizerEXK(result, vars, rhs)
	default:
		panic(errorx.UnknownEnumValue(config.EXKEncoder))
	}
}

func amkIncremental(
	result Result,
	config *Config,
	vars []f.Variable,
	rhs int,
) (*CCIncrementalData, error) {
	if rhs >= len(vars) {
		return nil, errorx.BadInput("tautology AMK with rhs > len(vars)")
	}
	if rhs == 0 { // no variable can be true
		for _, variable := range vars {
			result.AddClause(variable.Negate(result.Factory()))
		}
		return nil, errorx.BadInput("trivial AMK with rhs = 0")
	}
	switch config.AMKEncoder {
	case AMKTotalizer:
		return totalizerAMK(result, vars, rhs), nil
	case AMKModularTotalizer, AMKBest:
		return modtotalizerAMK(result, f.VariablesAsLiterals(vars), rhs), nil
	case AMKCardinalityNetwork:
		return cnAmkForIncremental(result, vars, rhs), nil
	default:
		panic(errorx.BadFormulaSort(config.AMKEncoder))
	}
}

func alkIncremental(
	result Result,
	config *Config,
	vars []f.Variable,
	rhs int,
) (*CCIncrementalData, error) {
	if rhs > len(vars) {
		result.AddClause()
		return nil, errorx.BadInput("contradiction ALK with rhs > len(vars)")
	}
	if rhs == 0 {
		return nil, errorx.BadInput("tautology ALK with rhs = 0")
	}
	if rhs == 1 {
		result.AddClause(f.VariablesAsLiterals(vars)...)
		return nil, errorx.BadInput("trivial ALK with rhs = 1")
	}
	if rhs == len(vars) {
		for _, variable := range vars {
			result.AddClause(variable.AsLiteral())
		}
		return nil, errorx.BadInput("trivial ALK with rhs = len(vars)")
	}
	switch config.ALKEncoder {
	case ALKTotalizer:
		return totalizerALK(result, vars, rhs), nil
	case ALKModularTotalizer, ALKBest:
		return modtotalizerALK(result, f.VariablesAsLiterals(vars), rhs), nil
	case ALKCardinalityNetwork:
		return cnAlkForIncremental(result, vars, rhs), nil
	default:
		panic(errorx.UnknownEnumValue(config.ALKEncoder))
	}
}
