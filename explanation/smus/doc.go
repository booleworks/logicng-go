// Package smus provides an algorithm to compute a smallest minimal
// unsatisfiable set (SMUS) of a set of propositions or formulas in LogicNG.
//
// A smallest minimal unsatisfiable set (SMUS) is a smallest MUS based on the
// number of formulas it contains. So in contrast to a regular MUS it is not
// only locally minimal, but globally minimal.
//
// The implementation in LogicNG is based on "Smallest MUS extraction with
// minimal hitting set dualization" (Ignatiev, Previti, Liffiton &
// Marques-Silva, 2015).
//
// To compute a SMUS on a list of propositions you can call
//
//	smus.Compute(fac, propositions)
package smus
