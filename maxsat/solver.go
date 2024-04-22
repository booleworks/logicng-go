package maxsat

import (
	"fmt"
	"slices"

	"github.com/booleworks/logicng-go/configuration"
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/function"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/normalform"
)

// Result represents the result of a MAX-SAT computation.  It holds a flag
// whether the problem was satisfiable or not.  In case it was satisfiable, the
// final lower bound of the solver is stored as the Optimum.
type Result struct {
	Satisfiable bool
	Optimum     int
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
	var2index         map[f.Variable]int32
	index2var         map[int32]f.Variable
	selectorVariables []f.Variable
}

func newSolver(fac f.Factory, algorithm Algorithm, config ...*Config) *Solver {
	solver := &Solver{
		fac:           fac,
		algorithm:     algorithm,
		configuration: determineConfig(fac, config),
	}
	solver.Reset()
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

// Reset resets the MAX-SAT solver by clearing all internal data structures.
func (m *Solver) Reset() {
	m.result = nil
	m.var2index = make(map[f.Variable]int32)
	m.index2var = make(map[int32]f.Variable)
	m.selectorVariables = []f.Variable{}
	switch m.algorithm {
	case AlgLinearSU:
		m.solver = newLinearSU(m.configuration)
	case AlgLinearUS:
		m.solver = newLinearUS(m.configuration)
	case AlgMSU3:
		m.solver = newMSU3(m.configuration)
	case AlgWMSU3:
		m.solver = newWMSU3(m.configuration)
	case AlgWBO:
		m.solver = newWBO(m.configuration)
	case AlgIncWBO:
		m.solver = newIncWBO(m.configuration)
	case AlgOLL:
		m.solver = newOLL()
	}
}

// AddHardFormula adds the given formulas as hard formulas to the solver which
// must always be satisfied.  Since MAX-SAT solvers in LogicNG do not support
// an incremental interface, this function returns an error if the solver was
// already solved once.
func (m *Solver) AddHardFormula(formula ...f.Formula) error {
	if m.result != nil {
		return errorx.IllegalState("MAX-SAT solver does not support an incremental interface")
	}
	for _, formula := range formula {
		m.addCNF(normalform.CNF(m.fac, formula), -1)
	}
	return nil
}

// AddSoftFormula adds the given formulas as soft formulas with the given
// weight to the solver.  The weight must be >0 otherwise an error is returned.
// Since MAX-SAT solvers in LogicNG do not support an incremental interface,
// this function returns an error if the solver was already solved once.
func (m *Solver) AddSoftFormula(formula f.Formula, weight int) error {
	if m.result != nil {
		return errorx.IllegalState("MAX-SAT solver does not support an incremental interface")
	}
	if weight < 1 {
		return errorx.BadInput("the weight of a formula must be > 0")
	}
	selVar := m.fac.Var(fmt.Sprintf("%s%d", selPrefix, len(m.selectorVariables)))
	m.selectorVariables = append(m.selectorVariables, selVar)
	_ = m.AddHardFormula(m.fac.Or(selVar.Negate(m.fac).AsFormula(), formula))
	_ = m.AddHardFormula(m.fac.Or(formula.Negate(m.fac), selVar.AsFormula()))
	m.addClause(selVar.AsFormula(), weight)
	return nil
}

func (m *Solver) addCNF(formula f.Formula, weight int) {
	switch formula.Sort() {
	case f.SortTrue:
		break
	case f.SortFalse, f.SortLiteral, f.SortOr:
		m.addClause(formula, weight)
	case f.SortAnd:
		for _, op := range m.fac.Operands(formula) {
			m.addClause(op, weight)
		}
	default:
		panic(errorx.IllegalState("input formula is not a valid CNF: %s", formula.Sprint(m.fac)))
	}
}

func (m *Solver) addClause(formula f.Formula, weight int) {
	clauseVec := make([]int32, function.NumberOfAtoms(m.fac, formula))
	for i, lit := range f.Literals(m.fac, formula).Content() {
		variable := lit.Variable()
		index, ok := m.var2index[variable]
		if !ok {
			index = m.solver.newLiteral(false) >> 1
			m.var2index[variable] = index
			m.index2var[index] = variable
		}
		var litNum int32
		if lit.IsPos() {
			litNum = index * 2
		} else {
			litNum = (index * 2) ^ 1
		}
		clauseVec[i] = litNum
	}
	if weight == -1 {
		m.solver.addHardClause(clauseVec)
	} else {
		m.solver.setCurrentWeight(weight)
		m.solver.updateSumWeights(weight)
		m.solver.addSoftClause(weight, clauseVec)
	}
}

// Solve solves the MAX-SAT problem currently on the solver and returns the
// computation result.
func (m *Solver) Solve() Result {
	result, _ := m.SolveWithHandler(nil)
	return result
}

// SolveWithHandler solves the MAX-SAT problem currently on the solver.  The
// computation can be aborted with the given handler.  The computation result
// is returned and an ok flag which is false if the computation was aborted by
// the handler.
func (m *Solver) SolveWithHandler(maxsatHandler Handler) (result Result, ok bool) {
	if m.result != nil {
		return *m.result, true
	}
	if m.solver.getCurrentWeight() == 1 {
		m.solver.setProblemType(unweighted)
	} else {
		m.solver.setProblemType(weighted)
	}
	res, ok := m.solver.search(maxsatHandler)
	if !ok || res == resUndef {
		return Result{}, false
	}
	if res == resUnsat {
		m.result = &Result{Satisfiable: false, Optimum: -1}
	} else {
		m.result = &Result{Satisfiable: true, Optimum: m.solver.result()}
	}
	return *m.result, true
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
		variable, ok := m.index2var[int32(i)]
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
