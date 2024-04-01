// Package mus provides data structures and algorithms for minimal
// unsatisfiable sets (MUS).
//
// A MUS contains only those formulas which lead to the given set of formulas
// being unsatisfiable. In other words: If you remove at least one of the
// formulas in the MUS from the given set of formulas, your set of formulas is
// satisfiable. This means a MUS is locally minimal. Thus, given a set of
// formulas which is unsatisfiable, you can compute its MUS and have one
// locally minimal conflict description why it is unsatisfiable.
//
// To compute a MUS on a list of propositions, you can call
//
//	mus.ComputeInsertionBased(fac, propositions)
//
// or
//
//	mus.DeletionInsertionBased(fac, propositions)
package mus
