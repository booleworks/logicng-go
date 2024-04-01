package function

import f "booleworks.com/logicng/formula"

// LiteralProfile returns the number of occurrences of each literal in the
// given formula.
func LiteralProfile(fac f.Factory, formula f.Formula) map[f.Literal]int {
	cache := make(map[f.Literal]int)
	litProfileRec(fac, formula, cache)
	return cache
}

// VariableProfile  returns the number of occurrences of each variable in the
// given formula.
func VariableProfile(fac f.Factory, formula f.Formula) map[f.Variable]int {
	cache := make(map[f.Variable]int)
	varProfileRec(fac, formula, cache)
	return cache
}

func litProfileRec(fac f.Factory, formula f.Formula, cache map[f.Literal]int) {
	if formula.Sort() == f.SortLiteral {
		cache[f.Literal(formula)]++
	} else if formula.Sort() == f.SortPBC || formula.Sort() == f.SortCC {
		_, _, lits, _, _ := fac.PBCOps(formula)
		for _, l := range lits {
			litProfileRec(fac, l.AsFormula(), cache)
		}
	} else {
		for _, op := range fac.Operands(formula) {
			litProfileRec(fac, op, cache)
		}
	}
}

func varProfileRec(fac f.Factory, formula f.Formula, cache map[f.Variable]int) {
	if formula.Sort() == f.SortLiteral {
		variable := f.Literal(formula).Variable()
		cache[variable]++
	} else if formula.Sort() == f.SortPBC || formula.Sort() == f.SortCC {
		_, _, lits, _, _ := fac.PBCOps(formula)
		for _, l := range lits {
			variable := l.Variable()
			varProfileRec(fac, variable.AsFormula(), cache)
		}
	} else {
		for _, op := range fac.Operands(formula) {
			varProfileRec(fac, op, cache)
		}
	}
}
