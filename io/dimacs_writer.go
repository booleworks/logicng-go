package io

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/normalform"
)

const (
	cnfExtension = ".cnf"
)

// WriteDimacs writes the given formula in CNF to a file with the given
// filename in the
// http://www.satcompetition.org/2009/format-benchmarks2009.html. Returns a
// mapping from each variable of the original problem to its index in the CNF
// file an optional error if there was a problem writing the file or the
// formula was not in CNF.
func WriteDimacs(fac f.Factory, filename string, formula f.Formula) (map[f.Variable]int, error) {
	var name string
	if strings.HasSuffix(filename, cnfExtension) {
		name = filename
	} else {
		name = filename + cnfExtension
	}
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return WriteDimacsToWriter(fac, file, formula)
}

// WriteDimacsToFile writes the given formula in CNF to the given file in the
// http://www.satcompetition.org/2009/format-benchmarks2009.html. Returns a
// mapping from each variable of the original problem to its index in the CNF
// file an optional error if there was a problem writing the file or the
// formula was not in CNF.
func WriteDimacsToFile(fac f.Factory, file *os.File, formula f.Formula) (map[f.Variable]int, error) {
	return WriteDimacsToWriter(fac, file, formula)
}

// WriteDimacsToWriter writes the given formula to the given writer in the
// http://www.satcompetition.org/2009/format-benchmarks2009.html. Returns a
// mapping from each variable of the original problem to its index in the CNF
// file an optional error if there was a problem writing to the writer.
func WriteDimacsToWriter(fac f.Factory, writer io.Writer, formula f.Formula) (map[f.Variable]int, error) {
	if !normalform.IsCNF(fac, formula) {
		return nil, errorx.BadInput("formula is not in CNF")
	}

	var2id := map[f.Variable]int{}
	for i, variable := range f.Variables(fac, formula).Content() {
		var2id[variable] = i + 1
	}

	parts := make([]f.Formula, 0)
	if formula.Sort() == f.SortLiteral || formula.Sort() == f.SortOr {
		parts = append(parts, formula)
	} else {
		parts = append(parts, fac.Operands(formula)...)
	}

	partsSize := 1
	if formula.Sort() != f.SortFalse {
		partsSize = len(parts)
	}
	_, err := fmt.Fprintf(writer, "p cnf %d %d\n", len(var2id), partsSize)
	if err != nil {
		return nil, err
	}
	for _, part := range parts {
		lits := f.Literals(fac, part).Content()
		literals := make([]string, len(lits))
		for i, lit := range lits {
			_, phase, _ := fac.LitNamePhase(lit)
			prefix := ""
			if !phase {
				prefix = "-"
			}
			variable := lit.Variable()
			literals[i] = fmt.Sprintf("%s%d", prefix, var2id[variable])
		}
		_, err = fmt.Fprintf(writer, "%s 0\n", strings.Join(literals, " "))
		if err != nil {
			return nil, err
		}
	}
	if formula.Sort() == f.SortFalse {
		_, err := fmt.Fprintln(writer, "0")
		if err != nil {
			return nil, err
		}
	}
	return var2id, nil
}
