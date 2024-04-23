// Package sat provides LogicNG's SAT solver with a rich interface.
//
// The SAT problem is the problem of deciding whether a formula in Boolean
// logic is satisfiable or not. In other words, does there exist a variable
// assignment under which the formula evaluates to true?
//
// A SAT solver is a tool that, given a formula f can compute its
// satisfiability.
//
// A small example for using the solver:
//
//	fac := formula.NewFactory()
//	p := parser.New(fac)
//	solver := sat.NewSolver(fac)
//	solver.Add(p.ParseUnsafe("A & B & (C | D)"))
//	solver.Add(p.ParseUnsafe("A => X"))
//	result := solver.Sat() // true
//	model := solver.Model(fac.Vars("A", "B", "C")) // model on the solver
//
// LogicNG's SAT solver also has an incremental/decremental interface which can
// be used like this:
//
//	f1 := p.ParseUnsafe("A & B & C")
//	solver.Add(f1)
//	solver.Sat()                       // true
//	initialState := solver.SaveState() // save the initial state
//	solver.Add(p.ParseUnsafe("~A"))    // add another formula
//	solver.Sat()                       // false
//	solver.LoadState(initialState)     // load the initial state again
//	solver.Add(p.ParseUnsafe("D"))     // add another formula
//	solver.Sat()                       // true
package sat
