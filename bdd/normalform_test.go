package bdd

import (
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/randomizer"
	"github.com/booleworks/logicng-go/sat"
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
