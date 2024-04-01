package test

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/io"
	"github.com/booleworks/logicng-go/maxsat"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestWeightedMaxSat(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	fac := f.NewFactory()
	result := readResult("./data/longrunning_wms/result.txt")

	config := maxsat.DefaultConfig()
	config.WeightStrategy = maxsat.WeightDiversify
	solvers := [3]*maxsat.Solver{}
	solvers[0] = maxsat.OLL(fac)
	solvers[1] = maxsat.IncWBO(fac, config)
	solvers[2] = maxsat.IncWBO(fac)

	texts := [3]string{"OLL", "IncWBO (None)", "IncWBO (Diversify)"}

	folder := "./data/longrunning_wms/"
	items, _ := os.ReadDir(folder)
	for i, solver := range solvers {
		start := time.Now()
		for _, item := range items {
			if !item.IsDir() {
				item.Name()
				if strings.HasSuffix(item.Name(), "wcnf") {
					solver.Reset()
					maxsat.ReadDimacsToSolver(fac, solver, folder+item.Name())
					res := solver.Solve()
					assert.Equal(t, result[item.Name()], res.Optimum)
				}
			}
		}
		elapsed := time.Since(start) / 1_000_000
		t.Logf("%-18s: %.2f sec\n", texts[i], float64(elapsed)/1000.0)
	}
}

func TestOptimizerFunctionCompareWithMaxSat(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formulas, _ := io.ReadFormulas(fac, "../test/data/formulas/large2.txt")
	variables := f.Variables(fac, formulas...).Content()
	expected := 25

	assert.Equal(expected, solveMaxSat(fac, formulas, variables, maxsat.IncWBO(fac)))
	assert.Equal(expected, solveMaxSat(fac, formulas, variables, maxsat.LinearSU(fac)))
	assert.Equal(expected, solveMaxSat(fac, formulas, variables, maxsat.LinearUS(fac)))
	assert.Equal(expected, solveMaxSat(fac, formulas, variables, maxsat.MSU3(fac)))
	assert.Equal(expected, solveMaxSat(fac, formulas, variables, maxsat.WBO(fac)))
	assert.Equal(expected, solveMaxSat(fac, formulas, variables, maxsat.OLL(fac)))
	optimum := minimize(formulas, variables, sat.NewSolver(fac))
	assert.Equal(expected, len(optimum.PosVars()))
}

func solveMaxSat(fac f.Factory, formulas []f.Formula, variables []f.Variable, solver *maxsat.Solver) int {
	solver.AddHardFormula(formulas...)
	for _, variable := range variables {
		solver.AddSoftFormula(variable.Negate(fac).AsFormula(), 1)
	}
	return solver.Solve().Optimum
}

func minimize(formulas []f.Formula, literals []f.Variable, solver *sat.Solver) *model.Model {
	solver.Add(formulas...)
	return solver.Minimize(f.VariablesAsLiterals(literals))
}

func readResult(filename string) map[string]int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	result := make(map[string]int)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), ";")
		result[tokens[0]], _ = strconv.Atoi(tokens[1])
	}
	return result
}
