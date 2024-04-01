package bdd

import (
	"testing"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/normalform"
	"booleworks.com/logicng/randomizer"
	"booleworks.com/logicng/sat"
	"github.com/stretchr/testify/assert"
)

func TestBddNormalformCnf(t *testing.T) {
	fac := f.NewFactory()
	for i := 0; i < 100; i++ {
		rand := randomizer.NewWithSeed(fac, int64(i))
		formula := rand.Formula(3)
		cnf := CNF(fac, formula)
		assert.True(t, normalform.IsCNF(fac, cnf))
		assert.True(t, sat.IsEquivalent(fac, formula, cnf))
	}
}

func TestBddNormalformDnf(t *testing.T) {
	fac := f.NewFactory()
	for i := 0; i < 100; i++ {
		rand := randomizer.NewWithSeed(fac, int64(i))
		formula := rand.Formula(3)
		cnf := DNF(fac, formula)
		assert.True(t, normalform.IsDNF(fac, cnf))
		assert.True(t, sat.IsEquivalent(fac, formula, cnf))
	}
}
