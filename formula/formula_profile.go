package formula

// LiteralProfile returns the number of occurrences of each literal in the
// given formula.
func LiteralProfile(fac Factory, formula Formula) map[Literal]int {
	cache := make(map[Literal]int)
	litProfileRec(fac, formula, cache)
	return cache
}

// VariableProfile returns the number of occurrences of each variable in the
// given formula.
func VariableProfile(fac Factory, formula Formula) map[Variable]int {
	cache := make(map[Variable]int)
	varProfileRec(fac, formula, cache)
	return cache
}

func litProfileRec(fac Factory, formula Formula, cache map[Literal]int) {
	if formula.Sort() == SortLiteral {
		cache[Literal(formula)]++
	} else if formula.Sort() == SortPBC || formula.Sort() == SortCC {
		_, _, lits, _, _ := fac.PBCOps(formula)
		for _, l := range lits {
			cache[l]++
		}
	} else {
		for _, op := range fac.Operands(formula) {
			litProfileRec(fac, op, cache)
		}
	}
}

func varProfileRec(fac Factory, formula Formula, cache map[Variable]int) {
	if formula.Sort() == SortLiteral {
		variable := Literal(formula).Variable()
		cache[variable]++
	} else if formula.Sort() == SortPBC || formula.Sort() == SortCC {
		_, _, lits, _, _ := fac.PBCOps(formula)
		for _, l := range lits {
			cache[l.Variable()]++
		}
	} else {
		for _, op := range fac.Operands(formula) {
			varProfileRec(fac, op, cache)
		}
	}
}
