// Package assignment provides a data structure for Boolean assignments in
// LogicNG, mapping Boolean variables to truth values.
//
// In contrast to the model found in the [model] package, an assignment stores
// the variables internally in such a way that it can be efficiently used for
// evaluation and restriction.
//
// The primary use case for assignments is their usage in the functions
// [Evaluate] and [Restrict].
//
// The following example creates an assignment and restricts and evaluates a
// formula with it.
//
//	fac := formula.NewFactory()
//	formula := fac.And(fac.Variable("a"), fac.Variable("b")) // formula a & b
//	ass := assignment.New(fac, fac.Variable("a"))            // assigns a to true
//	eval := Evaluate(fac, formula, ass)                      // evaluates to false (b is false, since it is not in the assignment)
//	restrict := Restrict(fac, formula, ass)                  // restricts to b
package assignment
