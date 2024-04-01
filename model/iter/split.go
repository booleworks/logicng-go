package iter

import (
	"github.com/booleworks/logicng-go/sat"
	"math"
	"sort"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// A SplitProvider provides functionality to generate a new set of variables
// for model iteration functions.
type SplitProvider interface {
	Vars(solver *sat.Solver, variables *f.VarSet) *f.VarSet
}

// A FixedVarProvider always returns the same split vars.
type FixedVarProvider struct {
	vars *f.VarSet
}

// NewFixedVarProvider generates a new split variable provider which
// always returns the given variables.
func NewFixedVarProvider(variables *f.VarSet) *FixedVarProvider {
	return &FixedVarProvider{variables}
}

// Vars returns the split variables for the given solver and subset of
// variables.
func (p *FixedVarProvider) Vars(_ *sat.Solver, _ *f.VarSet) *f.VarSet {
	return p.vars
}

// A LeastCommonVarProvider is a split variable provider which provides split
// variables which occur particularly seldom in the formulas on the solver. The
// variables occurring in the formulas are sorted by their occurrence. This
// provider returns those variables with the smallest occurrence.
type LeastCommonVarProvider struct {
	takeRate float64
	maxVars  int
}

// NewLeastCommonVarProvider returns a new split variable provider which
// returns the variables with the smallest occurrence on the solver.  The take
// rate must be between 0 and 1 otherwise the function panics.
func NewLeastCommonVarProvider(takeRate float64, maxVars int) *LeastCommonVarProvider {
	if takeRate < 0 || takeRate > 1 {
		panic(errorx.BadInput("take rate must be between 0 and 1"))
	}
	return &LeastCommonVarProvider{takeRate, maxVars}
}

// DefaultLeastCommonVarProvider returns a new split variable provider which
// returns the variables with the smallest occurrence on the solver with
// default values.
func DefaultLeastCommonVarProvider() *LeastCommonVarProvider {
	return &LeastCommonVarProvider{0.5, 18}
}

// Vars returns the split variables for the given solver and subset of
// variables.
func (p *LeastCommonVarProvider) Vars(solver *sat.Solver, variables *f.VarSet) *f.VarSet {
	return chooseByOccs(solver, variables, false, p.takeRate, p.maxVars)
}

// A MostCommonVarProvider is a split variable provider which provides split
// variables which occur particularly often in the formulas on the solver. The
// variables occurring in the formulas are sorted by their occurrence. This
// provider returns those variables with the largest occurrence.
type MostCommonVarProvider struct {
	takeRate float64
	maxVars  int
}

// NewMostCommonVarProvider returns a new split variable provider which
// returns the variables with the largest occurrence on the solver.  The take
// rate must be between 0 and 1 otherwise the function panics.
func NewMostCommonVarProvider(takeRate float64, maxVars int) *MostCommonVarProvider {
	if takeRate < 0 || takeRate > 1 {
		panic(errorx.BadInput("take rate must be between 0 and 1"))
	}
	return &MostCommonVarProvider{takeRate, maxVars}
}

// DefaultMostCommonVarProvider returns a new split variable provider which
// returns the variables with the largest occurrence on the solver with
// default values.
func DefaultMostCommonVarProvider() *MostCommonVarProvider {
	return &MostCommonVarProvider{0.5, 18}
}

// Vars returns the split variables for the given solver and subset of
// variables.
func (p *MostCommonVarProvider) Vars(solver *sat.Solver, variables *f.VarSet) *f.VarSet {
	return chooseByOccs(solver, variables, true, p.takeRate, p.maxVars)
}

func chooseByOccs(
	solver *sat.Solver, variables *f.VarSet, mostCommon bool, takeRate float64, maxVars int,
) *f.VarSet {
	occs := solver.VarOccurrences(variables)
	vars := make([]f.Variable, 0, len(occs))
	for v := range occs {
		vars = append(vars, v)
	}
	if mostCommon {
		sort.SliceStable(vars, func(i, j int) bool {
			if diff := occs[vars[i]] - occs[vars[j]]; diff != 0 {
				return diff > 0
			} else {
				return vars[i] > vars[j]
			}
		})
	} else {
		sort.SliceStable(vars, func(i, j int) bool {
			if diff := occs[vars[i]] - occs[vars[j]]; diff != 0 {
				return diff < 0
			} else {
				return vars[i] < vars[j]
			}
		})
	}
	var limit int
	if variables != nil {
		limit = numbOfVarsToChoose(variables, takeRate, maxVars)
	} else {
		limit = numbOfVarsToChoose(f.NewVarSet(vars...), takeRate, maxVars)
	}
	return f.NewVarSet(vars[:limit]...)
}

func numbOfVarsToChoose(variables *f.VarSet, takeRate float64, maxNumOfVars int) int {
	return min(maxNumOfVars, int(math.Ceil(float64(variables.Size())*takeRate)))
}
