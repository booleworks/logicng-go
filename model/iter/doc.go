// Package iter gathers functionality for model iteration.
//
// All (potentially projected) models on a SAT solver are iterated.  Iteration
// can be configured with different strategies.  The default strategy is to
// recursively split the models and iterate over sub-sets.  Especially for large
// enumerations this should yield better performance than iterating in one go.
package iter
