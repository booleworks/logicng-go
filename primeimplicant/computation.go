package primeimplicant

import (
	"slices"

	"github.com/booleworks/logicng-go/handler"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/sat"
	"github.com/booleworks/logicng-go/transformation"
)

// CoverSort encodes the sort of the cover: Implicants or Implicates.
type CoverSort byte

const (
	CoverImplicants CoverSort = iota
	CoverImplicates
)

// PrimeResult gathers the result of a prime implicant computation and holds
// implicants, implicates, and the cover type.
type PrimeResult struct {
	Implicants [][]f.Literal
	Implicates [][]f.Literal
	CoverSort  CoverSort
}

const (
	pos = "_POS"
	neg = "_NEG"
)

// CoverMin computes prime implicants and prime implicates for a given
// formula using minimal models. The cover type specifies if the implicants or
// the implicates will be complete, the other one will still be a cover of the
// given formula.
func CoverMin(fac f.Factory, formula f.Formula, coverSort CoverSort) *PrimeResult {
	min, _ := compute(fac, formula, coverSort, false, nil)
	return min
}

// CoverMax computes prime implicants and prime implicates for a given
// formula using maximal models. The cover type specifies if the implicants or
// the implicates will be complete, the other one will still be a cover of the
// given formula.
func CoverMax(fac f.Factory, formula f.Formula, coverSort CoverSort) *PrimeResult {
	max, _ := compute(fac, formula, coverSort, true, nil)
	return max
}

// CoverMinWithHandler computes prime implicants and prime implicates
// for a given formula using minimal models. The cover type specifies if the
// implicants or the implicates will be complete, the other one will still be a
// cover of the given formula.  The given handler can be used to abort the
// SAT-Solver based optimization during the computation.
func CoverMinWithHandler(
	fac f.Factory, formula f.Formula, coverSort CoverSort, optimizationHandler sat.OptimizationHandler,
) (*PrimeResult, bool) {
	return compute(fac, formula, coverSort, false, optimizationHandler)
}

// CoverMaxWithHandler computes prime implicants and prime implicates
// for a given formula using maximal models. The cover type specifies if the
// implicants or the implicates will be complete, the other one will still be a
// cover of the given formula.  The given handler can be used to abort the
// SAT-Solver based optimization during the computation.
func CoverMaxWithHandler(
	fac f.Factory, formula f.Formula, coverSort CoverSort, optimizationHandler sat.OptimizationHandler,
) (*PrimeResult, bool) {
	return compute(fac, formula, coverSort, true, optimizationHandler)
}

func compute(
	fac f.Factory,
	formula f.Formula,
	coverSort CoverSort,
	maximize bool,
	optimizationHandler sat.OptimizationHandler,
) (*PrimeResult, bool) {
	handler.Start(optimizationHandler)
	completeImplicants := coverSort == CoverImplicants
	var formulaForComputation f.Formula
	if completeImplicants {
		formulaForComputation = formula
	} else {
		formulaForComputation = formula.Negate(fac)
	}
	implicants, implicates, ok := computeGeneric(fac, formulaForComputation, maximize, optimizationHandler)
	if !ok {
		return nil, false
	}
	if completeImplicants {
		return &PrimeResult{implicants, implicates, coverSort}, true
	} else {
		return &PrimeResult{negateAll(fac, implicates), negateAll(fac, implicants), coverSort}, true
	}
}

func computeGeneric(
	fac f.Factory,
	formula f.Formula,
	maximize bool,
	optimizationHandler sat.OptimizationHandler,
) ([][]f.Literal, [][]f.Literal, bool) {
	sub := createSubstitution(fac, formula)
	literals := make([]f.Literal, 0, len(sub.newVar2oldLit))
	for key := range sub.newVar2oldLit {
		literals = append(literals, key.AsLiteral())
	}
	slices.Sort(literals)
	hSolver := sat.NewSolver(fac)
	hSolver.Add(sub.constraintFormula)
	fSolver := sat.NewSolver(fac)
	fSolver.Add(formula.Negate(fac))
	primeReduction := newPrimeReduction(fac, formula)
	var primeImplicants, primeImplicates [][]f.Literal
	for {
		var hModel *model.Model
		var ok bool
		if maximize {
			hModel, ok = hSolver.MaximizeWithHandler(literals, optimizationHandler)
		} else {
			hModel, ok = hSolver.MinimizeWithHandler(literals, optimizationHandler)
		}
		if !ok {
			return nil, nil, false
		}
		if hModel == nil {
			return primeImplicants, primeImplicates, true
		}
		fModel := transformModel(hModel, &sub.newVar2oldLit)
		var satHandler sat.Handler
		if optimizationHandler != nil {
			satHandler = optimizationHandler.SatHandler()
		}
		fSat, ok := fSolver.SatWithHandler(satHandler, fModel.Literals...)
		if !ok {
			return nil, nil, false
		}
		if !fSat {
			var primeImplicant []f.Literal
			var ok bool
			if maximize {
				primeImplicant, ok = primeReduction.reduceImplicant(fModel.Literals, satHandler)
			} else {
				primeImplicant, ok = fModel.Literals, true
			}
			if !ok {
				return nil, nil, false
			}
			primeImplicants = append(primeImplicants, primeImplicant)
			blockingClause := make([]f.Formula, len(primeImplicant))
			for i, lit := range primeImplicant {
				blockingClause[i] = transformation.SubstituteLiterals(fac, lit.AsFormula(), &sub.substitution).Negate(fac)
			}
			hSolver.Add(fac.Or(blockingClause...))
		} else {
			var ls []f.Literal
			if maximize {
				ls = fModel.Literals
			} else {
				mdl, _ := fSolver.Model(f.Variables(fac, formula).Content())
				ls = mdl.Literals
			}
			implicate := make([]f.Literal, len(ls))
			for i, lit := range ls {
				implicate[i] = lit.Negate(fac)
			}
			primeImplicate, ok := primeReduction.reduceImplicate(fac, implicate, satHandler)
			if !ok {
				return nil, nil, false
			}
			primeImplicates = append(primeImplicates, primeImplicate)
			hSolver.Add(transformation.SubstituteLiterals(fac, fac.Or(f.LiteralsAsFormulas(primeImplicate)...), &sub.substitution))
		}
	}
}

func createSubstitution(fac f.Factory, formula f.Formula) *substitutionResult {
	vars := f.Variables(fac, formula).Content()
	newVar2oldLit := make(map[f.Variable]f.Literal, len(vars)*2)
	substitution := make(map[f.Literal]f.Literal)
	constraintOps := make([]f.Formula, len(vars))
	for i, variable := range vars {
		name, _ := fac.VarName(variable)
		posVar := fac.Var(name + pos)
		newVar2oldLit[posVar] = variable.AsLiteral()
		substitution[variable.AsLiteral()] = posVar.AsLiteral()
		negVar := fac.Var(name + neg)
		newVar2oldLit[negVar] = variable.Negate(fac)
		substitution[variable.Negate(fac)] = negVar.AsLiteral()
		constraintOps[i] = fac.AMO(posVar, negVar)
	}
	return &substitutionResult{newVar2oldLit, substitution, fac.And(constraintOps...)}
}

func transformModel(mdl *model.Model, mapping *map[f.Variable]f.Literal) *model.Model {
	mapped := model.New()
	for _, variable := range mdl.PosVars() {
		mapped.AddLiteral((*mapping)[variable])
	}
	return mapped
}

func negateAll(fac f.Factory, literals [][]f.Literal) [][]f.Literal {
	result := make([][]f.Literal, len(literals))
	for i, lits := range literals {
		negated := make([]f.Literal, len(lits))
		for j, lit := range lits {
			negated[j] = lit.Negate(fac)
		}
		result[i] = negated
	}
	return result
}

type primeReduction struct {
	implicantSolver *sat.Solver
	implicateSolver *sat.Solver
}

func newPrimeReduction(fac f.Factory, formula f.Formula) *primeReduction {
	implicantSolver := sat.NewSolver(fac)
	implicantSolver.Add(formula.Negate(fac))
	implicateSolver := sat.NewSolver(fac)
	implicateSolver.Add(formula)
	return &primeReduction{implicantSolver, implicateSolver}
}

func (p *primeReduction) reduceImplicant(
	implicant []f.Literal, satHandler sat.Handler,
) ([]f.Literal, bool) {
	handler.Start(satHandler)
	primeImplicant := f.NewMutableLitSet(implicant...)
	for _, lit := range implicant {
		primeImplicant.Remove(lit)
		sat, ok := p.implicantSolver.SatWithHandler(satHandler, primeImplicant.Content()...)
		if !ok {
			return nil, false
		}
		if sat {
			primeImplicant.Add(lit)
		}
	}
	return primeImplicant.Content(), true
}

func (p *primeReduction) reduceImplicate(
	fac f.Factory, implicate []f.Literal, satHandler sat.Handler,
) ([]f.Literal, bool) {
	handler.Start(satHandler)
	primeImplicate := f.NewMutableLitSet(implicate...)
	for _, lit := range implicate {
		primeImplicate.Remove(lit)
		assumptions := make([]f.Literal, primeImplicate.Size())
		for i, lit := range primeImplicate.Content() {
			assumptions[i] = lit.Negate(fac)
		}
		sat, ok := p.implicateSolver.SatWithHandler(satHandler, assumptions...)
		if !ok {
			return nil, false
		}
		if sat {
			primeImplicate.Add(lit)
		}
	}
	return primeImplicate.Content(), true
}

type substitutionResult struct {
	newVar2oldLit     map[f.Variable]f.Literal
	substitution      map[f.Literal]f.Literal
	constraintFormula f.Formula
}
