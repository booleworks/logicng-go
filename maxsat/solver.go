package maxsat

import (
	"fmt"
	"slices"

	"github.com/booleworks/logicng-go/configuration"
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/sat"
)

// Result represents the result of a MAX-SAT computation.  It holds a flag
// whether the problem was satisfiable or not.  In case it was satisfiable, the
// final lower bound of the solver is stored as the Optimum.
type Result struct {
	Satisfiable bool
	Optimum     int
}

// A SolverState can be extracted from the solver by the SaveState method and be
// loaded again with the LoadState method.  It is used to mark certain states
// of the solver and be able to come back to them.
type SolverState struct {
	id            int32
	nbVars        int
	nbHard        int
	nbSoft        int
	ubCost        int
	currentWeight int
	softWeights   []int
}

const selPrefix = "@SEL_SOFT_"

// A Solver can be used to solve the MAX-SAT problem.  Depending on the
// underlying solving algorithm it supports also partial and/or weighted
// MAX-SAT problems.
type Solver struct {
	configuration     *Config
	algorithm         Algorithm
	fac               f.Factory
	result            *Result
	solver            algorithm
	pgTransformation  *pgOnSolver
	selectorVariables []f.Variable
}

func newSolver(fac f.Factory, algorithm Algorithm, config ...*Config) *Solver {
	solver := &Solver{
		fac:           fac,
		algorithm:     algorithm,
		configuration: determineConfig(fac, config),
	}
	solver.result = nil
	solver.selectorVariables = []f.Variable{}
	switch solver.algorithm {
	case AlgLinearSU:
		solver.solver = newLinearSU(fac, solver.configuration)
	case AlgLinearUS:
		solver.solver = newLinearUS(fac, solver.configuration)
	case AlgMSU3:
		solver.solver = newMSU3(fac, solver.configuration)
	case AlgWMSU3:
		solver.solver = newWMSU3(fac, solver.configuration)
	case AlgWBO:
		solver.solver = newWBO(fac, solver.configuration)
	case AlgIncWBO:
		solver.solver = newIncWBO(fac, solver.configuration)
	case AlgOLL:
		solver.solver = newOLL(fac)
	}
	if solver.configuration.CNFMethod != sat.CNFFactory {
		withNNF := solver.configuration.CNFMethod == sat.CNFPG
		solver.pgTransformation = newPGOnSolver(fac, withNNF, solver.solver)
	}
	return solver
}

func determineConfig(fac f.Factory, initConfig []*Config) *Config {
	if len(initConfig) > 0 {
		return initConfig[0]
	} else {
		configFromFactory, ok := fac.ConfigurationFor(configuration.MaxSat)
		if !ok {
			return DefaultConfig()
		} else {
			return configFromFactory.(*Config)
		}
	}
}

// LinearSU generates a new MAX-SAT solver with the Linear Sat-Unsat algorithm.
// This algorithm is based on linear search and supports both weighted and
// partial MAX-SAT problems.
func LinearSU(fac f.Factory, config ...*Config) *Solver {
	var cfg *Config
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultConfig()
	}
	return newSolver(fac, AlgLinearSU, cfg)
}

// LinearUS generates a new MAX-SAT solver with the Linear Unsat-Sat algorithm.
// This algorithm is based on linear search and supports partial MAX-SAT
// problems but no weighted problems.
func LinearUS(fac f.Factory, config ...*Config) *Solver {
	return newSolver(fac, AlgLinearUS, config...)
}

// MSU3 generates a new MAX-SAT solver with the MSU3 algorithm, a seminal-core
// guided algorithm. This algorithm is based on unsat cores and supports
// partial MAX-SAT problems but no weighted problems.
func MSU3(fac f.Factory, config ...*Config) *Solver {
	return newSolver(fac, AlgMSU3, config...)
}

// WMSU3 generates a new MAX-SAT solver with the weighted MSU3 algorithm, a
// seminal-core guided algorithm. This algorithm is based on unsat cores and
// supports both partial and weighted MAX-SAT problems.
func WMSU3(fac f.Factory, config ...*Config) *Solver {
	var cfg *Config
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultConfig()
		cfg.IncrementalStrategy = IncIterative
	}
	return newSolver(fac, AlgWMSU3, cfg)
}

// WBO generates a new MAX-SAT solver with the Weighted Boolean Optimization
// algorithm. This algorithm is based on unsat cores and supports both partial
// and weighted MAX-SAT problems.
func WBO(fac f.Factory, config ...*Config) *Solver {
	return newSolver(fac, AlgWBO, config...)
}

// IncWBO generates a new MAX-SAT solver with the Incremental Weighted Boolean
// Optimization algorithm. This algorithm is based on unsat cores and supports
// both partial and weighted MAX-SAT problems.
func IncWBO(fac f.Factory, config ...*Config) *Solver {
	return newSolver(fac, AlgIncWBO, config...)
}

// OLL generates a new MAX-SAT solver with the OLL algorithm. This algorithm is
// based on unsat cores and supports both partial and weighted MAX-SAT
// problems.
func OLL(fac f.Factory, config ...*Config) *Solver {
	var cfg *Config
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultConfig()
		cfg.IncrementalStrategy = IncIterative
	}
	return newSolver(fac, AlgOLL, cfg)
}

// AddHardFormula adds the given formulas as hard formulas to the solver which
// must always be satisfied.  Since MAX-SAT solvers in LogicNG do not support
// an incremental interface, this function returns an error if the solver was
// already solved once.
func (m *Solver) AddHardFormula(formula ...f.Formula) error {
	for _, formula := range formula {
		m.addFormulaAsCNF(formula, -1)
	}
	return nil
}

// AddSoftFormula adds the given formulas as soft formulas with the given
// weight to the solver.  The weight must be >0 otherwise an error is returned.
// Since MAX-SAT solvers in LogicNG do not support an incremental interface,
// this function returns an error if the solver was already solved once.
func (m *Solver) AddSoftFormula(formula f.Formula, weight int) error {
	if weight < 1 {
		return errorx.BadInput("the weight of a formula must be > 0")
	}
	selVar := m.fac.Var(fmt.Sprintf("%s%d", selPrefix, len(m.selectorVariables)))
	m.selectorVariables = append(m.selectorVariables, selVar)
	m.addFormulaAsCNF(m.fac.Or(selVar.Negate(m.fac).AsFormula(), formula), -1)
	m.addFormulaAsCNF(m.fac.Or(formula.Negate(m.fac), selVar.AsFormula()), -1)
	m.addFormulaAsCNF(selVar.AsFormula(), weight)
	return nil
}

// SaveState saves and returns the current solver state.
func (m *Solver) SaveState() *SolverState {
	return m.solver.saveState()
}

// LoadState loads the given state to the solver. ATTENTION: You can only load
// a state which was created by this instance of the solver before the current
// state. Only the sizes of the internal data structures are stored, meaning
// you can go back in time and restore a solver state with fewer variables
// and/or fewer clauses. It is not possible to import a solver state from
// another solver or another solving execution.  Returns with an error if the
// state is not valid on the solver.
func (m *Solver) LoadState(state *SolverState) error {
	err := m.solver.loadState(state)
	if err != nil {
		return err
	}
	return nil
}

func (m *Solver) addFormulaAsCNF(formula f.Formula, weight int) {
	m.result = nil
	if m.configuration.CNFMethod == sat.CNFFactory {
		m.addCNF(normalform.CNF(m.fac, formula), weight)
	} else {
		m.pgTransformation.addCNFToSolver(formula, weight)
	}
}

func (m *Solver) addCNF(formula f.Formula, weight int) {
	switch formula.Sort() {
	case f.SortTrue:
		break
	case f.SortFalse, f.SortLiteral, f.SortOr:
		m.solver.addClause(formula, weight)
	case f.SortAnd:
		for _, op := range m.fac.Operands(formula) {
			m.solver.addClause(op, weight)
		}
	default:
		panic(errorx.IllegalState("input formula is not a valid CNF: %s", formula.Sprint(m.fac)))
	}
}

// Solve solves the MAX-SAT problem currently on the solver and returns the
// computation result.
func (m *Solver) Solve() Result {
	result, _ := m.SolveWithHandler(handler.NopHandler)
	return result
}

// SolveWithHandler solves the MAX-SAT problem currently on the solver.  The
// computation can be canceled with the given handler.  The computation result
// is returned and handler state.
func (m *Solver) SolveWithHandler(hdl handler.Handler) (Result, handler.State) {
	if m.result != nil {
		return *m.result, succ
	}
	if m.solver.getCurrentWeight() == 1 {
		m.solver.setProblemType(unweighted)
	} else {
		m.solver.setProblemType(weighted)
	}
	res, state := m.solver.search(hdl)
	if !state.Success {
		return Result{}, state
	}
	if res == resUnsat {
		m.result = &Result{Satisfiable: false, Optimum: -1}
	} else {
		m.result = &Result{Satisfiable: true, Optimum: m.solver.result()}
	}
	return *m.result, succ
}

// Model returns the model for the last MAX-SAT computation.  It returns an
// error if the problem is not yet solved or it was unsatisfiable.
func (m *Solver) Model() (*model.Model, error) {
	if m.result == nil {
		return nil, errorx.IllegalState("MAX-SAT solver is not yet solved")
	}
	if !m.result.Satisfiable {
		return nil, errorx.IllegalState("MAX-SAT problem was not satisfiable")
	} else {
		return m.createModel(m.solver.getModel()), nil
	}
}

// SupportsWeighted reports whether the solver supports weighted problems.
func (m *Solver) SupportsWeighted() bool {
	return m.algorithm != AlgLinearUS && m.algorithm != AlgMSU3
}

// SupportsUnweighted reports whether the solver supports unweighted problems.
func (m *Solver) SupportsUnweighted() bool {
	return m.algorithm != AlgWMSU3
}

func (m *Solver) createModel(vec []bool) *model.Model {
	var mdl []f.Literal
	for i := 0; i < len(vec); i++ {
		variable, ok := m.solver.varForIndex(i)
		if ok && !slices.Contains(m.selectorVariables, variable) {
			if vec[i] {
				mdl = append(mdl, variable.AsLiteral())
			} else {
				mdl = append(mdl, variable.Negate(m.fac))
			}
		}
	}
	return model.New(mdl...)
}
