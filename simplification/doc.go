// Package simplification gaterns various simplification algorithms for Boolean
// formulas.
//
// The idea of the simplifiers is to simplify a given formula.
// But what is "simple" in terms of a formula? Since "simple" is no
// mathematically defined term and can alter from application to application,
// some simplifiers let the user provide their own definition of "simple". This
// is done via a rating function.
//
// A rating function is an interface which can be implemented by the user and
// computes a simplicity rating for a given formula. This could be for example
// the length of its string representation or the number of atoms. This rating
// function is then used to compare two formulas during the simplification
// process and thus deciding which of the formulas is the "simpler" one. There
// is a default rating function which is a rating function which compares
// formulas based on the length of their string representation (using the
// default string representation).
//
// Implemented simplifications are
// - propagating its backbone
// - propagating its unit literals
// - applying the distributive law
// - factoring out common parts of the formula
// - minimize negations
// - subsumption on CNF and DNF formulas
// - an advanced simplification combining various methods including minimal prime implicant covers
// - a Quine-McCluskey implementation based on the advanced simplifier
//
// To simplify a formula withe the advanced simplifier, you can simply call
//
//	simplified := simplification.Advanced(fac, formula)
package simplification
