package maxsat

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	f "github.com/booleworks/logicng-go/formula"
)

// ReadDimacsToSolver reads a Dimacs file for weighted MAX-SAT problems from
// the given filename and loads it directly to the given solver.  Returns
// an error if the file could not be read.
func ReadDimacsToSolver(fac f.Factory, solver *Solver, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	pureMaxSat := false
	hardWeight := -1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "p wcnf") {
			header := strings.Split(strings.TrimSpace(line), " ")
			if len(header) > 4 {
				hardWeight, _ = strconv.Atoi(header[4])
			}
			break
		} else if strings.HasPrefix(line, "p cnf") {
			pureMaxSat = true
			break
		}
	}
	for scanner.Scan() {
		tokens := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		weight, err := strconv.Atoi(tokens[0])
		if err != nil {
			return err
		}
		var start int
		if pureMaxSat {
			start = 0
		} else {
			start = 1
		}
		literals := make([]f.Formula, len(tokens)-start-1)
		for i := start; i < len(tokens)-1; i++ {
			if len(tokens[i]) > 0 {
				parsedLit, err := strconv.Atoi(tokens[i])
				if err != nil {
					return err
				}
				variable := fmt.Sprintf("v%d", int(math.Abs(float64(parsedLit))))
				literals[i-start] = fac.Literal(variable, parsedLit > 0)
			}
		}
		if pureMaxSat {
			solver.AddSoftFormula(fac.Or(literals...), 1)
		} else if weight == hardWeight {
			solver.AddHardFormula(fac.Or(literals...))
		} else {
			solver.AddSoftFormula(fac.Or(literals...), weight)
		}
	}
	return nil
}
