// Package dnnf provides datastructures and algorithms for compiling and
// manipulation deterministic, decomposable negation normal forms, or short
// [DNNF] in LogicNG.
//
// The d-DNNF of a formula has proven to be more succinct than it's BDD.
// Further, it helps to alleviate the ubiquitous memory explosion problem of
// BDDs with large formulas.
//
// A simple definition is:
//   - A formula in NNF is in decomposable negation normal form (DNNF) if the
//     decompositional property holds, that is, the operands of a conjunction do
//     not share variables.
//   - A DNNF is called deterministic (d-DNNF) if operands of a disjunction do
//     not share models.
//
// The following example compiles a DNNF from a given formula:
//
//	fac := f.NewFactory()
//	parser := parser.NewPropositionalParser(fac)
//	formula := parser.Parse("a | ((b & ~c) | (c & (~d | ~a & b)) & e)")
//	dnnf := dnnf.Compile(fac, formula)
//
// The dnnf itself is just a regular formula.  In LogicNG DNNFs are primarily
// used for model counting.
//
// [DNNF]: https://dl.acm.org/doi/10.1145/502090.502091
package dnnf
