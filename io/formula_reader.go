package io

import (
	"bufio"
	"os"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
)

// ReadFormula reads a formula from the given file.  If the file contains more
// than one line, the conjunction of all lines is returned.  Returns the
// formula and an optional error if there was a problem reading the file.
func ReadFormula(fac f.Factory, filename string) (f.Formula, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return ReadFormulaFile(fac, file)
}

// ReadFormulas reads a list of formulas from the given file.  Each line in the
// file is one formula in the result.  Returns the list of formulas ad an
// optional error if there was a problem reading the file.
func ReadFormulas(fac f.Factory, filename string) ([]f.Formula, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ReadFormulasFile(fac, file)
}

// ReadFormulaFile reads a formula from the given file.  If the file contains
// more than one line, the conjunction of all lines is returned.  Returns the
// formula and an optional error if there was a problem reading the file.
func ReadFormulaFile(fac f.Factory, file *os.File) (f.Formula, error) {
	ops, err := ReadFormulasFile(fac, file)
	if err != nil {
		return 0, err
	}
	return fac.And(ops...), nil
}

// ReadFormulasFile reads a list of formulas from the given file.  Each line in
// the file is one formula in the result.  Returns the list of formulas ad an
// optional error if there was a problem reading the file.
func ReadFormulasFile(fac f.Factory, file *os.File) ([]f.Formula, error) {
	parser := parser.New(fac)
	ops := make([]f.Formula, 0, 64)
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		ops = append(ops, parser.ParseUnsafe(line))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ops, nil
}
