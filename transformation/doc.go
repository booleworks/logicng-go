// Package transformation gathers various transformations on formulas.  A
// formula transformations always takes a formula factory and formula as input
// and yields some transformed formula.
//
// Currently, the following transformations are implemented
//   - formula anonymization
//   - substitution from variable to formula
//   - substitution from literal to literal
//   - pure expansion of AMO and EXO constraints
//   - existential and universal quantifier elimination
//   - CNF and DNF subsumption
package transformation
