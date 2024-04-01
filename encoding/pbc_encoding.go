package encoding

import (
	"math"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// EncodePBC encodes a pseudo-Boolean constraint to a CNF formula as a list of
// clauses.  This function can also be called with a cardinality constraint.
//
// It uses the encoding algorithm as specified either in the optional given
// config or if it not present the one configured in the formula factory.
//
// Returns an error if the input constraint is no valid pseudo-Boolean constraint.
func EncodePBC(fac f.Factory, pbc f.Formula, config ...*Config) ([]f.Formula, error) {
	result := ResultForFormula(fac)
	err := EncodePBCInResult(fac, pbc, result, config...)
	return result.result, err
}

// EncodePBCInResult encodes a pseudo-Boolean constraint into an encoding result.
// This result can be either a formula result and a formula is generated or it
// can be a solver encoding and the cnf is added directly to the solver without
// first generating formulas on the formula factory.
//
// It uses the encoding algorithm as specified either in the optional given
// config or if it not present the one configured in the formula factory.
//
// Returns an error if the input constraint is no valid pseudo-Boolean constraint.
func EncodePBCInResult(fac f.Factory, pbc f.Formula, result Result, config ...*Config) error {
	if pbc.Sort() == f.SortCC {
		return EncodeCCInResult(fac, pbc, result)
	}
	if pbc.Sort() != f.SortPBC {
		return errorx.BadFormulaSort(pbc.Sort())
	}
	cfg := determineConfig(fac, config)
	normalized := Normalize(fac, pbc)
	var err error
	switch normalized.Sort() {
	case f.SortFalse:
		result.AddClause()
	case f.SortCC, f.SortPBC:
		_, rhs, lits, coeffs, _ := fac.PBCOps(normalized)
		err = encodePBC(result, lits, coeffs, rhs, cfg)
	case f.SortAnd:
		for _, op := range fac.Operands(normalized) {
			err = EncodePBCInResult(fac, op, result, cfg)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func encodePBC(result Result, lits []f.Literal, coeffs []int, rhs int, config *Config) error {
	if rhs == math.MaxInt {
		return errorx.BadInput("overflow in the encoding")
	}
	if rhs < 0 {
		result.AddClause()
	}
	simplifiedLits := make([]f.Literal, 0, len(lits))
	simplifiedCoeffs := make([]int, 0, len(coeffs))
	if rhs == 0 {
		for _, lit := range lits {
			result.AddClause(lit.Negate(result.Factory()))
		}
		return nil
	}
	for i := 0; i < len(lits); i++ {
		if coeffs[i] <= rhs {
			simplifiedLits = append(simplifiedLits, lits[i])
			simplifiedCoeffs = append(simplifiedCoeffs, coeffs[i])
		} else {
			result.AddClause(lits[i].Negate(result.Factory()))
		}
	}
	if len(simplifiedLits) <= 1 {
		return nil
	}
	switch config.PBCEncoder {
	case PBCSWC, PBCBest:
		encodePBCSWC(result, simplifiedLits, simplifiedCoeffs, rhs)
	case PBCBinaryMerge:
		encodePBCBinaryMerge(result, simplifiedLits, simplifiedCoeffs, rhs, config)
	case PBCAdderNetworks:
		encodePBCAdder(result, simplifiedLits, simplifiedCoeffs, rhs)
	default:
		panic(errorx.UnknownEnumValue(config.PBCEncoder))
	}
	return nil
}
