// Package normalform provides algorithms for converting Boolean formulas into
// normal forms.
//
// The two most important normal forms are the conjunctive normal form (CNF)
// and the disjunctive normal form (DNF). In particular, the CNF is of special
// importance, since it is the input form required for SAT Solving and many
// other operations or algorithms.
//
// Another important normal form is the negation normal form (NNF) where only
// the operators ~, &, and | are allowed and negations must only appear before
// variables. In LogicNG this means that the formula consists only of literals
// Ands, and Ors. The NNF is often used as pre-processing step before
// transforming the formula into another normal form. Also, some algorithms
// require a formula to be in NNF to work.
//
// Simple normal-forms can be computed like this:
//
//	normalform.NNF(fac, formula)               // negation normal form
//	normalform.TseitinCNFDefault(fac, formula) // CNF due to Tseitin
//	normalform.PGCNFDefault(fac, formula)      // CNF due to Plaisted & Greenbaum
//	normalform.FactorizedCNF(fac, formula)     // CNF by factorization
//	normalform.FactorizedDNF(fac, formula)     // DNF by factorization
package normalform
