package normalform

import (
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// IsMinterm reports whether the given formula is a minterm, i.e. a conjunction
// of literals.
func IsMinterm(fac f.Factory, formula f.Formula) bool {
	return testTerm(fac, formula, true)
}

// IsMaxterm reports whether the given formula is a maxterm, i.e. a disjunction
// of literals.
func IsMaxterm(fac f.Factory, formula f.Formula) bool {
	return testTerm(fac, formula, false)
}

func testTerm(fac f.Factory, formula f.Formula, minterm bool) bool {
	switch fsort := formula.Sort(); fsort {
	case f.SortTrue, f.SortFalse, f.SortLiteral:
		return true
	case f.SortImpl, f.SortEquiv, f.SortNot, f.SortCC, f.SortPBC:
		return false
	case f.SortOr:
		return !minterm && onlyLiterals(fac, formula)
	case f.SortAnd:
		return minterm && onlyLiterals(fac, formula)
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
}

func onlyLiterals(fac f.Factory, nary f.Formula) bool {
	ops, _ := fac.NaryOperands(nary)
	for _, op := range ops {
		if op.Sort() != f.SortLiteral {
			return false
		}
	}
	return true
}
