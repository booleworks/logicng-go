package transformation

import (
	"fmt"

	f "github.com/booleworks/logicng-go/formula"
)

// An Anonymizer can be used to replace variables in a formula with newly
// generated ones with a given prefix.  Thus, it can be used to anonymize a
// formula.  An Anonymizer can be used for multiple transformations and holds
// its internal substitution map for the variables.
type Anonymizer struct {
	fac          f.Factory
	prefix       string
	counter      int
	Substitution *Substitution
}

// NewAnonymizer generates a new anonymizer with an optional prefix for the
// generated variables.  The default Prefix is `v`.
func NewAnonymizer(fac f.Factory, prefix ...string) *Anonymizer {
	pfx := "v"
	if prefix != nil {
		pfx = prefix[0]
	}
	return &Anonymizer{
		fac:          fac,
		prefix:       pfx,
		Substitution: NewSubstitution(),
	}
}

// Anonymize anonymizes the given formula by replacing all variables in it by
// newly generated ones, starting with v0, v1, ...
func Anonymize(fac f.Factory, formula f.Formula) f.Formula {
	anonymizer := NewAnonymizer(fac)
	return anonymizer.Anonymize(formula)
}

// Anonymize anonymizes the given formula by replacing all variables in it by
// their mapping on the anonymizer.
func (a *Anonymizer) Anonymize(formula f.Formula) f.Formula {
	vars := f.Variables(a.fac, formula)
	if vars.Empty() {
		return formula
	}
	for _, variable := range vars.Content() {
		_, ok := a.Substitution.subst[variable]
		if !ok {
			a.counter++
			a.Substitution.AddVar(variable, a.fac.Variable(fmt.Sprintf("%s%d", a.prefix, a.counter)))
		}
	}
	subst, _ := Substitute(a.fac, formula, a.Substitution)
	return subst
}
