// Package bdd provides data structures and algorithms on Binary Decision
// Diagrams (BDD) in LogicNG.
//
// A BDD is a directed acyclic graph of a given formula. It has a single root;
// Each inner node is labeled with a propositional variable and has one
// outgoing edge for a positive assignment, and one edge for a negative
// assignment of the respective variable. The leaves are labeled with 1 and 0
// representing true and false. An assignment is represented by a path from the
// root node to a leaf and its evaluation is the respective value of the leaf.
// Therefore, all paths to a 1-leaf are valid (possibly partial) models for the
// formula.
//
// Crucial for a small BDD representation of a formula is a good variable
// ordering.  There are different orderings implemented in LogicNG: based on
// variable occurrence, DFS- or BFS based, and the FORCE heuristic.
//
// A BDD is always compiled by a kernel.  This kernel holds all internal data
// structures used during the compilation, especially the node cache.
// Algorithms on BDDs always require the kernel which was used to compile the
// BDD.
//
// The following example creates an BDD of a formula without configuring a
// specific kernel or variable ordering:
//
//	fac := formula.NewFactory()
//	formula := fac.Or(fac.Variable("a"), fac.Variable("b")) // formula a / b
//	bdd := bdd.Build(fac, formula)
//
// You can also configure the node table size and the cache size of the kernel
// by hand.  The BDD kernel internally holds a table with all nodes in the BDD.
// This table can be extended dynamically, but this is an expensive operation.
// On the other hand, one wants to avoid reserving too much space for nodes,
// since this costs unnecessary memory. 30 * x proved to be efficient in
// practice for medium-sized formulas.  The following example manually creates
// a kernel with a node table size of 10 and a cache size of 100.
//
//	fac := formula.NewFactory()
//	formula := fac.Or(fac.Variable("a"), fac.Variable("b")) // formula a / b
//	kernel := bdd.NewKernel(fac, 2, 10, 100)
//	bdd := bdd.BuildWithKernel(fac, formula, kernel)
//
// Finally, you can also compute a variable ordering first and create the BDD
// with the given ordering.
//
//	fac := formula.NewFactory()
//	p := parser.NewPropositionalParser(fac)
//	formula := p.Parse("(A => ~B) & ((A & C) | (D & ~C)) & (A | Y | X) & (Y <=> (X | (W + A + F < 1)))")
//	ordering := bdd.DfsOrder(fac, formula)
//	kernel := bdd.NewKernelWithOrdering(fac, ordering, 10, 100)
//	bdd := bdd.BuildWithKernel(fac, formula, kernel)
package bdd
