package bdd

import (
	"slices"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/model"
)

func bddModelEnum(bdd *BDD, variables []f.Variable) []*model.Model {
	kernel := bdd.Kernel
	relevantIndices := make([]int32, 0, len(variables))
	for _, v := range variables {
		idx, ok := kernel.var2idx[v]
		if ok {
			relevantIndices = append(relevantIndices, idx)
		}
	}
	models := kernel.allSat(bdd.Index)
	res := make([]*model.Model, 0, len(models))
	for _, model := range models {
		generateAllModels(kernel, &res, model, relevantIndices, 0)
	}
	return res
}

func generateAllModels(kernel *Kernel, models *[]*model.Model, mdl []byte, relevantIndices []int32, position int) {
	if position == len(relevantIndices) {
		lits := make([]f.Literal, len(relevantIndices))
		for i, idx := range relevantIndices {
			var lit f.Literal
			if mdl[idx] == 0 {
				v, _ := kernel.getVariableForIndex(idx)
				lit = v.Negate(kernel.fac)
			} else {
				v, _ := kernel.getVariableForIndex(idx)
				lit = v.AsLiteral()
			}
			lits[i] = lit
		}
		mdl := model.New(lits...)
		if !containsModel(models, mdl) {
			*models = append(*models, mdl)
		}
	} else if mdl[relevantIndices[position]] != 2 {
		generateAllModels(kernel, models, mdl, relevantIndices, position+1)
	} else {
		mdl[relevantIndices[position]] = 0
		generateAllModels(kernel, models, mdl, relevantIndices, position+1)
		mdl[relevantIndices[position]] = 1
		generateAllModels(kernel, models, mdl, relevantIndices, position+1)
		mdl[relevantIndices[position]] = 2
	}
}

func containsModel(models *[]*model.Model, mdl *model.Model) bool {
	for _, m := range *models {
		if slices.Equal(m.Literals, mdl.Literals) {
			return true
		}
	}
	return false
}

func cnf(bdd *BDD) f.Formula {
	return bddNormalform(bdd, true)
}

func dnf(bdd *BDD) f.Formula {
	return bddNormalform(bdd, false)
}

func bddNormalform(bdd *BDD, cnf bool) f.Formula {
	kernel := bdd.Kernel
	var pathsToConstant [][]byte
	if cnf {
		pathsToConstant = kernel.allUnsat(bdd.Index)
	} else {
		pathsToConstant = kernel.allSat(bdd.Index)
	}
	terms := make([]f.Formula, 0)
	for _, path := range pathsToConstant {
		literals := make([]f.Literal, 0, len(path))
		for i := 0; i < len(path); i++ {
			variable, _ := kernel.getVariableForIndex(int32(i))
			switch path[i] {
			case 0:
				if cnf {
					literals = append(literals, variable.AsLiteral())
				} else {
					literals = append(literals, variable.Negate(kernel.fac))
				}
			case 1:
				if cnf {
					literals = append(literals, variable.Negate(kernel.fac))
				} else {
					literals = append(literals, variable.AsLiteral())
				}
			}
		}
		if cnf {
			terms = append(terms, kernel.fac.Clause(literals...))
		} else {
			terms = append(terms, kernel.fac.Minterm(literals...))
		}
	}
	if cnf {
		return kernel.fac.And(terms...)
	} else {
		return kernel.fac.Or(terms...)
	}
}
