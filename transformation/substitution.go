package transformation

import (
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// A Substitution represents a mapping from variables to formulas.  When
// executing the Substitute function with such a substitution, all variables
// in a formula are replaced by their mapped formula.
type Substitution struct {
	subst map[f.Variable]f.Formula
}

// NewSubstitution generates a new empty substitution.
func NewSubstitution() *Substitution {
	return &Substitution{make(map[f.Variable]f.Formula, 16)}
}

// AddVar adds a single mapping from variable to formula to the substitution.
func (s *Substitution) AddVar(variable f.Variable, formula f.Formula) {
	s.subst[variable] = formula
}

// AddMapping adds the given mapping from variables to formulas to the
// substitution.
func (s *Substitution) AddMapping(mapping map[f.Variable]f.Formula) {
	for k, v := range mapping {
		s.AddVar(k, v)
	}
}

// Substitute performs the given substitution on the given formula and returns
// a new formula where all variables of the substitution are replaced by their
// mapped formulas. Variables not in the substitution are left as-is.  The
// function returns an error when you try to substitute a formula for a literal
// which is used in a cardinality or Pseudo-Boolean constraint since this is
// not possible.
func Substitute(fac f.Factory, formula f.Formula, subst *Substitution) (f.Formula, error) {
	switch fsort := formula.Sort(); fsort {
	case f.SortTrue, f.SortFalse:
		return formula, nil
	case f.SortLiteral:
		variable := f.Literal(formula).Variable()
		s, ok := subst.subst[variable]
		if !ok {
			return formula, nil
		} else {
			if formula.IsPos() {
				return s, nil
			} else {
				return s.Negate(fac), nil
			}
		}
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		subst, err := Substitute(fac, op, subst)
		if err != nil {
			return 0, err
		}
		return fac.Not(subst), nil
	case f.SortEquiv, f.SortImpl:
		left, right, _ := fac.BinaryLeftRight(formula)
		lSubst, err := Substitute(fac, left, subst)
		if err != nil {
			return 0, err
		}
		rSubst, err := Substitute(fac, right, subst)
		if err != nil {
			return 0, err
		}
		binOp, _ := fac.BinaryOperator(formula.Sort(), lSubst, rSubst)
		return binOp, nil
	case f.SortOr, f.SortAnd:
		ops, _ := fac.NaryOperands(formula)
		operands := make([]f.Formula, len(ops))
		var err error
		for i, op := range ops {
			operands[i], err = Substitute(fac, op, subst)
			if err != nil {
				return 0, err
			}
		}
		naryOp, _ := fac.NaryOperator(formula.Sort(), operands...)
		return naryOp, nil
	case f.SortCC, f.SortPBC:
		return substitutePbc(fac, formula, subst)
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
}

func substitutePbc(fac f.Factory, pbc f.Formula, substitution *Substitution) (f.Formula, error) {
	comparator, rhs, operands, coefficients, _ := fac.PBCOps(pbc)
	newLits := make([]f.Literal, 0, len(operands))
	newCoeffs := make([]int, 0, len(operands))
	lhsFixed := 0
	for i := range operands {
		variable := operands[i].Variable()
		subst, ok := substitution.subst[variable]
		if !ok {
			newLits = append(newLits, operands[i])
			newCoeffs = append(newCoeffs, coefficients[i])
		} else {
			switch fsort := subst.Sort(); fsort {
			case f.SortTrue:
				if operands[i].IsPos() {
					lhsFixed += coefficients[i]
				}
			case f.SortFalse:
				if operands[i].IsNeg() {
					lhsFixed += coefficients[i]
				}
			case f.SortLiteral:
				var newOp f.Literal
				if operands[i].IsPos() {
					newOp = f.Literal(subst)
				} else {
					newOp = f.Literal(subst).Negate(fac)
				}
				newLits = append(newLits, newOp)
				newCoeffs = append(newCoeffs, coefficients[i])
			default:
				return 0, errorx.BadInput("tried to substitute %s in a PBC", subst.Sprint(fac))
			}
		}
	}
	if len(newLits) == 0 {
		return fac.Constant(comparator.Evaluate(lhsFixed, rhs)), nil
	} else {
		return fac.PBC(comparator, rhs-lhsFixed, newLits, newCoeffs), nil
	}
}
