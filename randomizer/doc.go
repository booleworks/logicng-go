// Package randomizer provides a formula randomizer which can be used for
// fuzzing in testing in LogicNG.
//
// To generate a randomized formula with a nested size of 4, you can call
//
//	randomizer := randomizer.New(fac)
//	randomFormula := randomizer.Formula(4)
package randomizer
