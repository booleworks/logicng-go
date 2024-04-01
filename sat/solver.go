package sat

import (
	"sort"

	"github.com/booleworks/logicng-go/configuration"
	"github.com/booleworks/logicng-go/encoding"
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/normalform"
)

// A Solver is the main interface an external user should interact with a SAT
// solver.  It provides methods for adding and removing formulas to the solver,
// solving them, extracting models, generating proofs for unsatisfiable
// formulas, computing backbones, or optimize the formula on the solver.
type Solver struct {
	fac                            f.Factory
	result                         f.Tristate
	config                         Config
	core                           *CoreSolver
	lastComputationWithAssumptions bool
	pgTransformation               *pgOnSolver
	fullPgTransformation           *pgOnSolver
}

// A SolverState can be extracted from the solver by the SaveState method and be
// loaded again with the LoadState method.  It is used to mark certain states
// of the solver (the loaded formulas and learnt clauses) and be able to come
// back to them.
type SolverState struct {
	id    int32
	state []int
}

func newSolver(fac f.Factory, config *Config) *Solver {
	solver := NewCoreSolver(config, UncheckedEnqueue)
	return &Solver{
		fac:                  fac,
		config:               *config,
		core:                 solver,
		result:               f.TristateUndef,
		pgTransformation:     newPGOnSolver(fac, true, solver, config.InitialPhase),
		fullPgTransformation: newPGOnSolver(fac, false, solver, config.InitialPhase),
	}
}

// NewSolver creates a new SAT solver with the optional configuration.
func NewSolver(fac f.Factory, config ...*Config) *Solver {
	cfg := determineConfig(fac, config)
	return newSolver(fac, cfg)
}

func determineConfig(fac f.Factory, initConfig []*Config) *Config {
	if len(initConfig) > 0 {
		return initConfig[0]
	} else {
		configFromFactory, ok := fac.ConfigurationFor(configuration.Sat)
		if !ok {
			return DefaultConfig()
		} else {
			return configFromFactory.(*Config)
		}
	}
}

// Add adds the given formulas to the solver.  If the formulas are not already
// in CNF, they are converted by the CNFMethod configured in the solver's
// configuration.
func (s *Solver) Add(formulas ...f.Formula) {
	for _, formula := range formulas {
		s.addWithProp(formula, nil)
	}
}

// AddProposition adds the given propositions to the solver.  Propositions wrap
// formulas with some additional information.  If generating proofs for
// unsatisfiable formulas, it is a good idea to use propositions, since
// otherwise you just get clauses of the internal solver formulas as result.
func (s *Solver) AddProposition(propositions ...f.Proposition) {
	for _, prop := range propositions {
		s.addWithProp(prop.Formula(), prop)
	}
}

func (s *Solver) addWithProp(formula f.Formula, proposition f.Proposition) {
	s.result = f.TristateUndef
	if formula.Sort() == f.SortCC {
		if s.config.UseAtMostClauses {
			comparator, rhs, literals, _, _ := s.fac.PBCOps(formula)
			if comparator == f.LE {
				s.core.addAtMost(s.generateClauseVector(literals), rhs)
			} else if comparator == f.LT && rhs > 3 {
				s.core.addAtMost(s.generateClauseVector(literals), rhs-1)
			} else if comparator == f.EQ && rhs == 1 {
				s.core.addAtMost(s.generateClauseVector(literals), rhs)
				s.core.AddClause(s.generateClauseVector(literals), proposition)
			} else {
				s.addFormulaAsCNF(formula, proposition)
			}
		} else {
			result := resultForSolver(s.fac, s, proposition)
			_ = encoding.EncodeCCInResult(s.fac, formula, result) // we know it is a cardinality constraint
		}
	} else if formula.Sort() == f.SortPBC {
		result := resultForSolver(s.fac, s, proposition)
		err := encoding.EncodePBCInResult(s.fac, formula, result)
		if err != nil {
			panic(err)
		}
	} else {
		s.addFormulaAsCNF(formula, proposition)
	}
	s.addAllOriginalVars(formula)
}

func (s *Solver) addAllOriginalVars(originalFormula f.Formula) {
	for _, v := range f.Variables(s.fac, originalFormula).Content() {
		s.getOrAddIndex(v.AsLiteral())
	}
}

func (s *Solver) addFormulaAsCNF(formula f.Formula, proposition f.Proposition) {
	switch s.config.CNFMethod {
	case CNFFactorization:
		s.addClauseSet(normalform.CNF(s.fac, formula), proposition)
	case CNFPG:
		s.pgTransformation.addCNFToSolver(formula, proposition)
	case CNFFullPG:
		s.fullPgTransformation.addCNFToSolver(formula, proposition)
	default:
		panic(errorx.UnknownEnumValue(s.config.CNFMethod))
	}
}

func (s *Solver) addClauseSet(formula f.Formula, proposition f.Proposition) {
	switch formula.Sort() {
	case f.SortTrue:
		break
	case f.SortFalse, f.SortLiteral, f.SortOr:
		s.addClause(formula, proposition)
	case f.SortAnd:
		nary, ok := s.fac.NaryOperands(formula)
		if !ok {
			panic(errorx.UnknownFormula(formula))
		}
		for _, op := range nary {
			s.addClause(op, proposition)
		}
	default:
		panic(errorx.IllegalState("not a valid CNF"))
	}
}

func (s *Solver) addClause(formula f.Formula, proposition f.Proposition) {
	s.result = f.TristateUndef
	ps := s.generateClauseVector(f.Literals(s.fac, formula).Content())
	s.core.AddClause(ps, proposition)
}

func (s *Solver) generateClauseVector(literals []f.Literal) []int32 {
	clause := make([]int32, len(literals))
	sort.Slice(literals, func(i, j int) bool { return literals[i] < literals[j] })
	for i, lit := range literals {
		_, phase, _ := s.fac.LitNamePhase(lit)
		index := s.getOrAddIndex(lit)
		var litNum int32
		if phase {
			litNum = index * 2
		} else {
			litNum = (index * 2) ^ 1
		}
		clause[i] = litNum
	}
	return clause
}

func (s *Solver) getOrAddIndex(lit f.Literal) int32 {
	name, _, _ := s.fac.LitNamePhase(lit)
	index, ok := s.core.name2idx[name]
	if !ok {
		index = s.core.NewVar(!s.config.InitialPhase, true)
		s.core.addName(name, index)
	}
	return index
}

// Sat solves the formula on the solver and returns whether it is satisfiable.
// If assumptions are given,  they are assumed before solving the formula.  In
// this case the satisfiable result has to be interpreted with regard to the
// assumptions.
func (s *Solver) Sat(assumptions ...f.Literal) bool {
	result, _ := s.SatWithHandler(nil, assumptions...)
	return result
}

// SatWithHandler solves the formula on the solver and returns whether it is
// satisfiable. If assumptions are given,  they are assumed before solving the
// formula.  In this case the satisfiable result has to be interpreted with
// regard to the assumptions.  The handler can be used to abort the solving
// process.  If the computation was aborted, the ok flag will be false.
func (s *Solver) SatWithHandler(handler Handler, assumptions ...f.Literal) (result, ok bool) {
	if len(assumptions) == 0 {
		if s.lastResultIsUsable() {
			return s.result == f.TristateTrue, true
		}
		s.result, ok = s.core.Solve(handler)
	} else {
		assVec := s.generateClauseVector(assumptions)
		s.result, ok = s.core.SolveWithAssumptions(handler, assVec)
	}
	s.lastComputationWithAssumptions = len(assumptions) > 0
	return s.result == f.TristateTrue, ok
}

func (s *Solver) lastResultIsUsable() bool {
	return s.result != f.TristateUndef && !s.lastComputationWithAssumptions
}

// Model returns a model of the formula currently on the solver.  The model
// will include only the given variables.  Returns with an error if the solver
// was not yet solved or the formula was unsatisfiable.
func (s *Solver) Model(variables []f.Variable) (*model.Model, error) {
	if s.result == f.TristateUndef {
		return nil, errorx.IllegalState("SAT solver is not yet solved")
	}
	var relevantIndices []int32
	if variables != nil {
		relevantIndices = make([]int32, len(variables))
		for i, v := range variables {
			name, _ := s.fac.VarName(v)
			relevantIndices[i] = s.core.IdxForName(name)
		}
	}
	if s.result == f.TristateTrue {
		return s.core.CreateModel(s.fac, s.core.model, relevantIndices), nil
	} else {
		return nil, errorx.IllegalState("SAT problem was not satisfiable")
	}
}

// SaveState saves and returns the current solver state.
func (s *Solver) SaveState() *SolverState {
	return s.core.saveState()
}

// LoadState loads the given state to the solver.  Returns with an error if the
// state is not valid on the solver.
func (s *Solver) LoadState(state *SolverState) error {
	err := s.core.loadState(state)
	if err != nil {
		return err
	}
	s.result = f.TristateUndef
	s.pgTransformation.clearCache()
	s.fullPgTransformation.clearCache()
	return nil
}

// AddIncrementalCC adds the given constraint as an incremental cardinality
// constraint to the solver.  It returns the incremental data used to tighten
// the bound of the formula on the solver.  Returns with an error if the
// incremental constraint could not be generated.
func (s *Solver) AddIncrementalCC(cc f.Formula) (*encoding.CCIncrementalData, error) {
	result := resultForSolver(s.fac, s, nil)
	return encoding.EncodeIncremental(s.fac, cc, result)
}

// KnownVariables returns the variables currently known by the solver.
func (s *Solver) KnownVariables() *f.VarSet {
	result := f.NewVarSet()
	nVars := s.core.NVars()
	for k, v := range s.core.name2idx {
		if v < nVars {
			result.Add(s.fac.Var(k))
		}
	}
	return result
}

// SetResult can be used to set the solver's result from outside.  You should
// never have a reason to use this from the outside except you develop your own
// solver functions.
func (s *Solver) SetResult(result f.Tristate) {
	s.result = result
}

// Factory returns the solver's formula factory.
func (s *Solver) Factory() f.Factory {
	return s.fac
}

// CoreSolver returns the core solver.  You should not need this from the outside.
func (s *Solver) CoreSolver() *CoreSolver {
	return s.core
}

// Reset resets the solver to its initial state.
func (s *Solver) Reset() {
	s.result = f.TristateUndef
	s.core.reset()
	s.lastComputationWithAssumptions = false
	s.pgTransformation.clearCache()
	s.fullPgTransformation.clearCache()
}
