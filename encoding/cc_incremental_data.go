package encoding

import (
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// CCIncrementalData gathers data for an incremental at-most-k cardinality
// constraint.  When an at-most-k cardinality constraint is constructed, it is
// possible to save incremental data with it.  Then one can modify the
// constraint after it was created by tightening the original bound.
type CCIncrementalData struct {
	Result     Result // encoding result of the incremental constraint
	amkEncoder AMKEncoder
	alkEncoder ALKEncoder
	vector1    []f.Literal
	vector2    []f.Literal
	mod        int
	nVars      int
	currentRhs int
}

// NewUpperBound tightens the upper bound of an at-most-k constraint and
// returns the resulting formula.
//
// Usage constraints:
//   - New right-hand side must be smaller than current right-hand side.
//   - Cannot be used for at-least-k constraints.
func (cc *CCIncrementalData) NewUpperBound(rhs int) []f.Formula {
	cc.computeUbConstraint(rhs)
	return cc.Result.Formulas()
}

// NewUpperBoundForSolver tightens the upper bound of an at-most-k constraint
// and encodes it on the solver of the result.
//
// Usage constraints:
//   - New right-hand side must be smaller than current right-hand side.
//   - Cannot be used for at-least-k constraints.
func (cc *CCIncrementalData) NewUpperBoundForSolver(rhs int) {
	cc.computeUbConstraint(rhs)
}

func (cc *CCIncrementalData) computeUbConstraint(rhs int) {
	fac := cc.Result.Factory()
	if rhs >= cc.currentRhs {
		panic(errorx.BadInput("new upper bound %d does not tighten the current bound of %d", rhs, cc.currentRhs))
	}
	cc.currentRhs = rhs
	switch cc.amkEncoder {
	case AMKModularTotalizer:
		ulimit := (rhs + 1) / cc.mod
		llimit := (rhs + 1) - ulimit*cc.mod
		for i := ulimit; i < len(cc.vector1); i++ {
			cc.Result.AddClause(cc.vector1[i].Negate(fac))
		}
		if ulimit != 0 && llimit != 0 {
			for i := llimit - 1; i < len(cc.vector2); i++ {
				cc.Result.AddClause(cc.vector1[ulimit-1].Negate(fac), cc.vector2[i].Negate(fac))
			}
		} else {
			if ulimit == 0 {
				for i := llimit - 1; i < len(cc.vector2); i++ {
					cc.Result.AddClause(cc.vector2[i].Negate(fac))
				}
			} else {
				cc.Result.AddClause(cc.vector1[ulimit-1].Negate(fac))
			}
		}
	case AMKTotalizer:
		for i := rhs; i < len(cc.vector1); i++ {
			cc.Result.AddClause(cc.vector1[i].Negate(fac))
		}
	case AMKCardinalityNetwork:
		if len(cc.vector1) > rhs {
			cc.Result.AddClause(cc.vector1[rhs].Negate(fac))
		}
	default:
		panic(errorx.UnknownEnumValue(cc.amkEncoder))
	}
}

// NewLowerBound tightens the lower bound of an at-least-k constraint and
// returns the resulting formula.
//
// Usage constraints:
//   - New right-hand side must be greater than current right-hand side.
//   - Cannot be used for at-most-k constraints.
func (cc *CCIncrementalData) NewLowerBound(rhs int) []f.Formula {
	cc.computeLbConstraint(rhs)
	return cc.Result.Formulas()
}

// NewLowerBoundForSolver tightens the lower bound of an at-least-k constraint
// and encodes it on the solver of the result.
//
// Usage constraints:
//   - New right-hand side must be greater than current right-hand side.
//   - Cannot be used for at-most-k constraints.
func (cc *CCIncrementalData) NewLowerBoundForSolver(rhs int) {
	cc.computeLbConstraint(rhs)
}

func (cc *CCIncrementalData) computeLbConstraint(rhs int) {
	fac := cc.Result.Factory()
	if rhs <= cc.currentRhs {
		panic(errorx.BadInput("new lower bound %d does not tighten the current bound of %d", rhs, cc.currentRhs))
	}
	cc.currentRhs = rhs
	switch cc.alkEncoder {
	case ALKTotalizer:
		for i := range rhs {
			cc.Result.AddClause(cc.vector1[i])
		}
	case ALKModularTotalizer:
		newRhs := cc.nVars - rhs
		ulimit := (newRhs + 1) / cc.mod
		llimit := (newRhs + 1) - ulimit*cc.mod
		for i := ulimit; i < len(cc.vector1); i++ {
			cc.Result.AddClause(cc.vector1[i].Negate(fac))
		}
		if ulimit != 0 && llimit != 0 {
			for i := llimit - 1; i < len(cc.vector2); i++ {
				cc.Result.AddClause(cc.vector1[ulimit-1].Negate(fac), cc.vector2[i].Negate(fac))
			}
		} else {
			if ulimit == 0 {
				for i := llimit - 1; i < len(cc.vector2); i++ {
					cc.Result.AddClause(cc.vector2[i].Negate(fac))
				}
			} else {
				cc.Result.AddClause(cc.vector1[ulimit-1].Negate(fac))
			}
		}
	case ALKCardinalityNetwork:
		newRhs := cc.nVars - rhs
		if len(cc.vector1) > newRhs {
			cc.Result.AddClause(cc.vector1[newRhs].Negate(fac))
		}
	default:
		panic(errorx.UnknownEnumValue(cc.alkEncoder))
	}
}
