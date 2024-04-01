package sat

import (
	f "github.com/booleworks/logicng-go/formula"
)

type solverEncoding struct {
	fac         f.Factory
	proposition f.Proposition
	solver      *Solver
}

func resultForSolver(fac f.Factory, solver *Solver, proposition f.Proposition) *solverEncoding {
	return &solverEncoding{fac, proposition, solver}
}

func (r *solverEncoding) Reset() {}

func (r *solverEncoding) AddClause(literals ...f.Literal) {
	clause := make([]int32, len(literals))
	for i, literal := range literals {
		r.addLiteral(&clause, i, literal)
	}
	r.solver.core.AddClause(clause, r.proposition)
	r.solver.result = f.TristateUndef
}

func (r *solverEncoding) addLiteral(clauseVec *[]int32, idx int, lit f.Literal) {
	solver := r.solver.core
	name, phase, _ := r.fac.LitNamePhase(lit)
	index := solver.IdxForName(name)
	if index == -1 {
		index = solver.NewVar(!r.solver.config.InitialPhase, true)
		solver.addName(name, index)
	}
	var litNum int32
	if phase {
		litNum = index * 2
	} else {
		litNum = (index * 2) ^ 1
	}
	(*clauseVec)[idx] = litNum
}

func (r *solverEncoding) NewCcVariable() f.Variable {
	return r.addVarToSolver(r.fac.NewCCVariable())
}

func (r *solverEncoding) NewPbVariable() f.Variable {
	return r.addVarToSolver(r.fac.NewPBCVariable())
}

func (r *solverEncoding) addVarToSolver(variable f.Variable) f.Variable {
	solver := r.solver.core
	index := solver.NewVar(!r.solver.config.InitialPhase, true)
	name, _ := r.fac.VarName(variable)
	solver.addName(name, index)
	return variable
}

func (r *solverEncoding) Factory() f.Factory {
	return r.fac
}

func (r *solverEncoding) Formulas() []f.Formula {
	return nil
}
