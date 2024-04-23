// Package enum provides algorithms to perform model enumeration on formulas in
// LogicNG.
//
// An enumeration can be performed on a formula:
//
//	assert := assert.New(t)
//	fac := f.NewFactory()
//	p := parser.New(fac)
//	formula := p.ParseUnsafe("A & (B | C)")
//	models := enum.OnFormula(fac, formula, fac.Vars("A", "B", "C")) // will produce 3 models
//
// or directly on a SAT solver
//
//	solver := sat.NewMiniSatSolver(fac)
//	models = enum.OnSolver(solver, fac.Vars("A", "B", "C")) // will produce 3 models
//
// Model enumeration is one use case of the model iteration in the [iter]
// package and can be configured as described there.
package enum
