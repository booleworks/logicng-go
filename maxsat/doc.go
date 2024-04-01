// Package maxsat provides an implementation of a MAX-SAT solver with different
// solving algorithms.
//
// Given an unsatisfiable formula in CNF, the MAX-SAT problem is the problem of
// finding an assignment which satisfies the maximum number of clauses and
// therefore solves an optimization problem rather than a decision problem.
// There are two extensions to MAX-SAT Solving: 1) the distinction of hard/soft
// clauses, and 2) additional weighted clauses, yielding four different
// flavours of MAX-SAT solving:
//  1. Pure MaxSAT
//  2. Partial MaxSAT
//  3. Weighted MaxSAT
//  4. Weighted Partial MaxSAT
//
// In a Partial MAX-SAT problem you can distinguish between hard and soft
// clauses. A hard clause must be satisfied whereas a soft clause should be
// satisfied but can be left unsatisfied. This means the solver only optimizes
// over the soft clauses. If the hard clauses themselves are unsatisfiable, no
// solution can be found.
//
// In a Weighted MAX-SAT problem clauses can have a positive weight. The solver
// then does not optimize the number of satisfied clauses but the sum of the
// weights of the satisfied clauses.
//
// The Weighted Partial MAX-SAT problem is the combination of Partial MaxSAT and
// weighted MAX-SAT. I.e. you can add hard clauses and weighted soft clauses to
// the MAX-SAT solver.
//
// Note two important points:
//   - MAX-SAT can be defined as weighted MAX-SAT restricted to formulas whose
//     clauses have weight 1, and as Partial MAX-SAT in the case that all the
//     clauses are declared to be soft.
//   - The above definitions talk about clauses on the solver, not arbitrary
//     formulas. In real-world use cases you often want to weight whole formulas
//     and not just clauses. LogicNG's MAX-SAT solver API gives you this
//     possibility and internally translates the formulas and their weights
//     accordingly.
//
// A small example for using the solver:
//
//	fac := formula.NewFactory()
//	p := parser.New(fac)
//	solver := maxsat.OLL(fac)
//	solver.AddHardFormula(p.ParseUnsafe("A & B & (C | D)"))
//	solver.AddSoftFormula(p.ParseUnsafe("A => ~B"), 2)
//	solver.AddSoftFormula(p.ParseUnsafe("~C"), 4)
//	solver.AddSoftFormula(p.ParseUnsafe("~D"), 8)
//	result := solver.Solve() // {Satisfiable: true, Optimum: 6}
package maxsat
