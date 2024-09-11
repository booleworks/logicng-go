package primeimplicant

import (
	"slices"

	"github.com/booleworks/logicng-go/event"
	"github.com/booleworks/logicng-go/handler"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/sat"
	"github.com/booleworks/logicng-go/transformation"
)

// CoverSort encodes the sort of the cover: Implicants or Implicates.
type CoverSort byte

var succ = handler.Success()

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
	min, _ := compute(fac, formula, coverSort, false, handler.NopHandler)
	return min
}

// CoverMax computes prime implicants and prime implicates for a given
// formula using maximal models. The cover type specifies if the implicants or
// the implicates will be complete, the other one will still be a cover of the
// given formula.
func CoverMax(fac f.Factory, formula f.Formula, coverSort CoverSort) *PrimeResult {
	max, _ := compute(fac, formula, coverSort, true, handler.NopHandler)
	return max
}

// CoverMinWithHandler computes prime implicants and prime implicates
// for a given formula using minimal models. The cover type specifies if the
// implicants or the implicates will be complete, the other one will still be a
// cover of the given formula.  The given handler can be used to cancel the
// SAT-Solver based optimization during the computation.
func CoverMinWithHandler(
	fac f.Factory, formula f.Formula, coverSort CoverSort, hdl handler.Handler,
) (*PrimeResult, handler.State) {
	return compute(fac, formula, coverSort, false, hdl)
}

// CoverMaxWithHandler computes prime implicants and prime implicates
// for a given formula using maximal models. The cover type specifies if the
// implicants or the implicates will be complete, the other one will still be a
// cover of the given formula.  The given handler can be used to cancel the
// SAT-Solver based optimization during the computation.
func CoverMaxWithHandler(
	fac f.Factory, formula f.Formula, coverSort CoverSort, hdl handler.Handler,
) (*PrimeResult, handler.State) {
	return compute(fac, formula, coverSort, true, hdl)
}

func compute(
	fac f.Factory,
	formula f.Formula,
	coverSort CoverSort,
	maximize bool,
	hdl handler.Handler,
) (*PrimeResult, handler.State) {
	if e := event.PrimeComputationStarted; !hdl.ShouldResume(e) {
		return nil, handler.Cancelation(e)
	}
	completeImplicants := coverSort == CoverImplicants
	var formulaForComputation f.Formula
	if completeImplicants {
		formulaForComputation = formula
	} else {
		formulaForComputation = formula.Negate(fac)
	}
	implicants, implicates, state := computeGeneric(fac, formulaForComputation, maximize, hdl)
	if !state.Success {
		return nil, state
	}
	if completeImplicants {
		return &PrimeResult{implicants, implicates, coverSort}, succ
	} else {
		return &PrimeResult{negateAll(fac, implicates), negateAll(fac, implicants), coverSort}, succ
	}
}

func computeGeneric(
	fac f.Factory,
	formula f.Formula,
	maximize bool,
	hdl handler.Handler,
) ([][]f.Literal, [][]f.Literal, handler.State) {
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
		var state handler.State
		if maximize {
			hModel, state = hSolver.MaximizeWithHandler(literals, hdl)
		} else {
			hModel, state = hSolver.MinimizeWithHandler(literals, hdl)
		}
		if !state.Success {
			return nil, nil, state
		}
		if hModel == nil {
			return primeImplicants, primeImplicates, succ
		}
		fModel := transformModel(hModel, &sub.newVar2oldLit)
		params := sat.Params().
			Handler(hdl).
			WithModel(f.Variables(fac, formula).Content()).
			Literal(fModel.Literals...)
		fResult := fSolver.Call(params)
		if fResult.Canceled() {
			return nil, nil, fResult.State()
		}
		if !fResult.Sat() {
			var primeImplicant []f.Literal
			var state handler.State
			if maximize {
				primeImplicant, state = primeReduction.reduceImplicant(fModel.Literals, hdl)
			} else {
				primeImplicant, state = fModel.Literals, succ
			}
			if !state.Success {
				return nil, nil, state
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
				ls = fResult.Model().Literals
			}
			implicate := make([]f.Literal, len(ls))
			for i, lit := range ls {
				implicate[i] = lit.Negate(fac)
			}
			primeImplicate, state := primeReduction.reduceImplicate(fac, implicate, hdl)
			if !state.Success {
				return nil, nil, state
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

func (p *primeReduction) reduceImplicant(implicant []f.Literal, hdl handler.Handler) ([]f.Literal, handler.State) {
	if e := event.ImplicantReductionStarted; !hdl.ShouldResume(e) {
		return nil, handler.Cancelation(e)
	}
	primeImplicant := f.NewMutableLitSet(implicant...)
	for _, lit := range implicant {
		primeImplicant.Remove(lit)
		sResult := p.implicantSolver.Call(sat.Params().Handler(hdl).Literal(primeImplicant.Content()...))
		if sResult.Canceled() {
			return nil, sResult.State()
		}
		if sResult.Sat() {
			primeImplicant.Add(lit)
		}
	}
	return primeImplicant.Content(), succ
}

func (p *primeReduction) reduceImplicate(
	fac f.Factory, implicate []f.Literal, hdl handler.Handler,
) ([]f.Literal, handler.State) {
	if e := event.ImplicateReductionStarted; !hdl.ShouldResume(e) {
		return nil, handler.Cancelation(e)
	}
	primeImplicate := f.NewMutableLitSet(implicate...)
	for _, lit := range implicate {
		primeImplicate.Remove(lit)
		assumptions := make([]f.Literal, primeImplicate.Size())
		for i, lit := range primeImplicate.Content() {
			assumptions[i] = lit.Negate(fac)
		}
		sResult := p.implicateSolver.Call(sat.Params().Handler(hdl).Literal(assumptions...))
		if sResult.Canceled() {
			return nil, sResult.State()
		}
		if sResult.Sat() {
			primeImplicate.Add(lit)
		}
	}
	return primeImplicate.Content(), succ
}

type substitutionResult struct {
	newVar2oldLit     map[f.Variable]f.Literal
	substitution      map[f.Literal]f.Literal
	constraintFormula f.Formula
}
