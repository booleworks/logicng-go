// Package function gathers various functions on formulas.  A formula function
// always takes a formula factory and formula as input and yields some
// computation result on this formula.
//
// Currently, the following functions are implemented
//   - depth of a formula
//   - number of atoms
//   - number of nodes
//   - variable profile (how often does each variable occur in the formula)
//   - literal profile (how often does each literal occur in the formula)
//   - sub-formulas of a formula
package function
