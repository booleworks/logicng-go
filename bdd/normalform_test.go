package bdd

import (
	"testing"

	"github.com/booleworks/logicng-go/event"
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

func TestBddNormalformCnfWithHandler(t *testing.T) {
	fac := f.NewFactory()
	rand := randomizer.NewWithSeed(fac, int64(42))
	formula := rand.Formula(5)
	hdl := HandlerWithNodes(5)
	cnf, state := CNFWithHandler(fac, formula, hdl)
	assert.False(t, state.Success)
	assert.Equal(t, event.BddNewRefAdded, state.CancelCause)
	assert.Equal(t, fac.Falsum(), cnf)
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

func TestBddNormalformDnfWithHandler(t *testing.T) {
	fac := f.NewFactory()
	rand := randomizer.NewWithSeed(fac, int64(42))
	formula := rand.Formula(5)
	hdl := HandlerWithNodes(5)
	dnf, state := DNFWithHandler(fac, formula, hdl)
	assert.False(t, state.Success)
	assert.Equal(t, event.BddNewRefAdded, state.CancelCause)
	assert.Equal(t, fac.Falsum(), dnf)
}
