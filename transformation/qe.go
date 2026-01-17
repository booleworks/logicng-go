package transformation

import (
	"github.com/booleworks/logicng-go/assignment"
	f "github.com/booleworks/logicng-go/formula"
)

// ExistentialQE eliminates a number of existentially quantified variables by
// replacing them with the Shannon expansion.  If x is eliminated from the
// formula, the resulting formula is formula[true/x] | formula[false/x].
func ExistentialQE(fac f.Factory, formula f.Formula, variable ...f.Variable) f.Formula {
	result := formula
	for _, v := range variable {
		pos, _ := assignment.New(fac, v.AsLiteral())
		neg, _ := assignment.New(fac, v.Negate(fac))
		result = fac.Or(assignment.Restrict(fac, result, pos), assignment.Restrict(fac, result, neg))
	}
	return result
}

// UniversalQE eliminates a number of universally quantified variables by
// replacing them with the Shannon expansion.  If x is eliminated from the
// formula, the resulting formula is formula[true/x] & formula[false/x].
func UniversalQE(fac f.Factory, formula f.Formula, variable ...f.Variable) f.Formula {
	result := formula
	for _, variable := range variable {
		pos, _ := assignment.New(fac, variable.AsLiteral())
		neg, _ := assignment.New(fac, variable.Negate(fac))
		result = fac.And(assignment.Restrict(fac, result, pos), assignment.Restrict(fac, result, neg))
	}
	return result
}
