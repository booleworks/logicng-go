package io

import (
	"fmt"
	"io"
	"os"

	f "github.com/booleworks/logicng-go/formula"
)

// WriteFormula writes the given formula to a file with the given filename. The
// flag splitAnd indicates whether - if the formula is a conjunction - the
// single operands should be written to different lines without a conjoining
// operator.  Returns an error if there was a problem writing the file.
func WriteFormula(fac f.Factory, filename string, formula f.Formula, splitAnd ...bool) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return WriteFormulaToWriter(fac, file, formula, splitAnd...)
}

// WriteFormulaToFile writes the given formula to the given file. The flag
// splitAnd indicates whether - if the formula is a conjunction - the single
// operands should be written to different lines without a conjoining operator.
// Returns an error if there was a problem writing the file.
func WriteFormulaToFile(fac f.Factory, file *os.File, formula f.Formula, splitAnd ...bool) error {
	return WriteFormulaToWriter(fac, file, formula, splitAnd...)
}

// WriteFormulaToWriter writes the given formula to the given writer. The flag
// splitAnd indicates whether - if the formula is a conjunction - the single
// operands should be written to different lines without a conjoining operator.
// Returns an error if there was a problem writing the writer.
func WriteFormulaToWriter(fac f.Factory, writer io.Writer, formula f.Formula, splitAnd ...bool) error {
	split := false
	if splitAnd != nil {
		split = splitAnd[0]
	}
	var err error
	if split && formula.Sort() == f.SortAnd {
		for _, op := range fac.Operands(formula) {
			_, err = fmt.Fprintln(writer, op.Sprint(fac))
			if err != nil {
				return err
			}
		}
	} else {
		_, err = fmt.Fprintln(writer, formula.Sprint(fac))
	}
	return err
}
