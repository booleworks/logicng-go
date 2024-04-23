// Package io gathers readers and writers for formulas in different formats in
// LogicNG.
//
// To write a formula to a file, you can use the
//
//	io.WriteFormula(fac, "filename", formula)
//
// To read the file again, simply use
//
//	read, err := io.ReadFormula(fac, "filename")
package io
