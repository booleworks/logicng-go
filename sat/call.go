package sat

import (
	"github.com/booleworks/logicng-go/errorx"
	"github.com/booleworks/logicng-go/explanation"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/model"
)

// CallParams describe the parameters for a single SAT solver call.
type CallParams struct {
	handler     Handler
	addProps    []f.Proposition
	modelVars   []f.Variable
	modelIfSat  bool
	coreIfUnsat bool
	upZeroIfSat bool
}

// Params generates a new empty call parameter struct with the following default setting:
//   - no handler
//   - no additional formulas or propositions for the SAT call
//   - no model generation for satisfiable formulas
//   - no unsat core computation for unsatisfiable formulas
//   - no computation of propagated literals at decision level 0
func Params() *CallParams {
	return &CallParams{}
}

// WithModel generates a new parameter struct with the following setting:
//   - model generation for satisfiable formulas for the given variables
//   - no handler
//   - no additional formulas or propositions for the SAT call
//   - no unsat core computation for unsatisfiable formulas
//   - no computation of propagated literals at decision level 0
func WithModel(variables []f.Variable) *CallParams {
	return &CallParams{modelIfSat: true, modelVars: variables}
}

// WithCore generates a new parameter struct with the following setting:
//   - unsat core computation for unsatisfiable formulas
//   - no handler
//   - no additional formulas or propositions for the SAT call
//   - no model generation for satisfiable formulas
//   - no computation of propagated literals at decision level 0
func WithCore() *CallParams {
	return &CallParams{coreIfUnsat: true}
}

// WithAssumptions generates a new parameter struct with the following setting:
//   - additional assumption literals for the SAT call
//   - no unsat core computation for unsatisfiable formulas
//   - no handler
//   - no model generation for satisfiable formulas
//   - no computation of propagated literals at decision level 0
func WithAssumptions(literals []f.Literal) *CallParams {
	params := &CallParams{}
	params.Literal(literals...)
	return params
}

// Handler sets a handler for the SAT call
func (p *CallParams) Handler(handler Handler) *CallParams {
	p.handler = handler
	return p
}

// WithModel activates model generation after the SAT solver call.  The model
// will be generated only for the given variables.  If the solver is
// unsatisfiable, no model will be generated.
func (p *CallParams) WithModel(variables []f.Variable) *CallParams {
	p.modelIfSat = true
	p.modelVars = variables
	return p
}

// WithCore activates an unsat core computation after the SAT solver call.
// If the solver is satisfiable, no core will be computed.
func (p *CallParams) WithCore() *CallParams {
	p.coreIfUnsat = true
	return p
}

// WithUPZeros activates computation of literals propagated at decision level 0
// after the SAT solver call.  If the solver is unsatisfiable, no literals will
// be computed.
func (p *CallParams) WithUPZeros() *CallParams {
	p.upZeroIfSat = true
	return p
}

// Variable sets additional variables which will be added to the SAT solver
// before solving.  Results like satisfiability, model, or unsat core are with
// respect to these additional variables.
func (p *CallParams) Variable(variable ...f.Variable) *CallParams {
	for _, v := range variable {
		p.addProps = append(p.addProps, f.NewStandardProposition(v.AsFormula()))
	}
	return p
}

// Literal sets additional literals which will be added to the SAT solver
// before solving.  Results like satisfiability, model, or unsat core are with
// respect to these additional literals.
func (p *CallParams) Literal(literal ...f.Literal) *CallParams {
	for _, l := range literal {
		p.addProps = append(p.addProps, f.NewStandardProposition(l.AsFormula()))
	}
	return p
}

// Formula sets additional formulas which will be added to the SAT solver
// before solving.  Results like satisfiability, model, or unsat core are with
// respect to these additional formulas.
func (p *CallParams) Formula(formula ...f.Formula) *CallParams {
	for _, form := range formula {
		p.addProps = append(p.addProps, f.NewStandardProposition(form))
	}
	return p
}

// Proposition sets additional propositions which will be added to the SAT solver
// before solving.  Results like satisfiability, model, or unsat core are with
// respect to these additional propositions.
func (p *CallParams) Proposition(proposition []f.Proposition) *CallParams {
	p.addProps = append(p.addProps, proposition...)
	return p
}

// CallResult represents the result of a single call to the SAT solver.
type CallResult struct {
	ok          bool
	satisfiable bool
	model       *model.Model
	core        *explanation.UnsatCore
	upZeroLits  []f.Literal
}

// OK reports whether the call to the SAT solver yielded a result and was not
// aborted.
func (r CallResult) OK() bool {
	return r.ok
}

// Aborted reports whether the SAT solver call was aborted by the given
// handler.
func (r CallResult) Aborted() bool {
	return !r.ok
}

// Sat reports whether the SAT solver call returned SAT or UNSAT.
func (r CallResult) Sat() bool {
	return r.satisfiable
}

// Model returns the model of the last SAT call if the formula was satisfiable
// and model generation was requested in the call.
func (r CallResult) Model() *model.Model {
	return r.model
}

// UnsatCore returns the unsatisfiable core of the last SAT call if the formula
// was unsatisfiable and unsat core computation was requested in the call.
func (r CallResult) UnsatCore() *explanation.UnsatCore {
	return r.core
}

// UpZeroLits returns the propagated literals at decision level 0 of the last
// SAT call if the formula was satisfiable and UpZero computation was requested
// in the call.
func (r CallResult) UpZeroLits() []f.Literal {
	return r.upZeroLits
}

// Call calls the SAT solver with the given call parameters.  Such a call
// always performs a solving process.  If additional variables / literals /
// formulas were set, these are added to the solver before solving.  When there
// are only variables and literals this is done by assumption solving,
// otherwise the solver's save and load state capabilities are used.  Depending
// on the request a model, unsat core, or propagated literals are also
// computed.
func (s *Solver) Call(params ...*CallParams) CallResult {
	var model *model.Model
	var core *explanation.UnsatCore
	var upZeroLits []f.Literal
	var param *CallParams
	if params == nil {
		param = Params()
	} else {
		param = params[0]
	}
	if param.coreIfUnsat && !s.config.ProofGeneration {
		panic(errorx.IllegalState("core computation on a SAT solver without proof tracing"))
	}
	call := initCall(s, param.handler, param.addProps)
	if call.ok && call.sat && param.modelIfSat {
		model = s.computeModel(param.modelVars)
	}
	if call.ok && !call.sat && param.coreIfUnsat {
		core = s.computeUnsatCore()
	}
	if call.ok && call.sat && param.upZeroIfSat {
		upZeroLits = s.computeUpZeroLits()
	}
	call.close()
	return CallResult{call.ok, call.sat, model, core, upZeroLits}
}

type call struct {
	solver            *Solver
	pgOriginalClauses int
	initialState      *SolverState
	sat               bool
	ok                bool
}

func initCall(solver *Solver, handler Handler, addProps []f.Proposition) *call {
	c := call{solver: solver}
	c.solver.core.startCall()
	if c.solver.config.ProofGeneration {
		c.pgOriginalClauses = len(c.solver.core.pgOriginalClauses)
	}
	adds := splitPropsIntoLitsAndFormulas(addProps)
	if len(adds.lits) > 0 {
		c.solver.core.assumptions = c.solver.generateClauseVector(adds.lits)
		c.solver.core.assumptionProps = adds.propsForLits
	}
	if len(adds.props) > 0 {
		c.initialState = c.solver.SaveState()
		for _, p := range adds.props {
			c.solver.AddProposition(p)
		}
	}
	res, ok := c.solver.core.Solve(handler)
	c.ok = ok
	c.sat = res == f.TristateTrue
	return &c
}

func splitPropsIntoLitsAndFormulas(additionalPropositions []f.Proposition) additionals {
	var lits []f.Literal
	var propsForLits []f.Proposition
	var props []f.Proposition
	for _, prop := range additionalPropositions {
		if prop.Formula().Sort() == f.SortLiteral {
			lits = append(lits, f.Literal(prop.Formula()))
			propsForLits = append(propsForLits, prop)
		} else {
			props = append(props, prop)
		}
	}
	return additionals{lits, propsForLits, props}
}

func (c *call) close() {
	c.solver.core.assumptions = nil
	c.solver.core.assumptionProps = nil
	if c.solver.config.ProofGeneration {
		shrinkTo(&c.solver.core.pgOriginalClauses, c.pgOriginalClauses)
	}
	if c.initialState != nil {
		err := c.solver.LoadState(c.initialState)
		if err != nil {
			panic(err)
		}
	}
	c.solver.core.satHandler = nil
	c.solver.core.finishCall()
}

type additionals struct {
	lits         []f.Literal
	propsForLits []f.Proposition
	props        []f.Proposition
}
