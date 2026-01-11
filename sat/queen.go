package sat

import (
	"fmt"

	e "github.com/booleworks/logicng-go/encoding"
	f "github.com/booleworks/logicng-go/formula"
)

// GenerateNQueens generates an n-queens problem of size n and returns it as a
// formula.
func GenerateNQueens(fac f.Factory, n int) f.Formula {
	ec := &e.Config{AMOEncoder: e.AMOPure}
	kk := 1
	varNames := make([][]f.Variable, n)
	for i := range n {
		varNames[i] = make([]f.Variable, n)
		for j := range n {
			varNames[i][j] = fac.Var(fmt.Sprintf("v%d", kk))
			kk++
		}
	}

	operands := make([]f.Formula, 0)

	for i := range n {
		vars := varNames[i]
		cc, _ := e.EncodeCC(fac, fac.EXO(vars...), ec)
		encoding := fac.And(cc...)
		operands = append(operands, encoding)
	}
	for i := range n {
		vars := make([]f.Variable, n)
		for j := range n {
			vars[j] = varNames[j][i]
		}
		cc, _ := e.EncodeCC(fac, fac.EXO(vars...), ec)
		encoding := fac.And(cc...)
		operands = append(operands, encoding)
	}
	for i := 0; i < n-1; i++ {
		vars := make([]f.Variable, n-i)
		for j := 0; j < n-i; j++ {
			vars[j] = varNames[j][i+j]
		}
		cc, _ := e.EncodeCC(fac, fac.AMO(vars...), ec)
		encoding := fac.And(cc...)
		operands = append(operands, encoding)
	}
	for i := 1; i < n-1; i++ {
		vars := make([]f.Variable, n-i)
		for j := 0; j < n-i; j++ {
			vars[j] = varNames[j+i][j]
		}
		cc, _ := e.EncodeCC(fac, fac.AMO(vars...), ec)
		encoding := fac.And(cc...)
		operands = append(operands, encoding)
	}
	for i := 0; i < n-1; i++ {
		vars := make([]f.Variable, n-i)
		for j := 0; j < n-i; j++ {
			vars[j] = varNames[j][n-1-(i+j)]
		}
		cc, _ := e.EncodeCC(fac, fac.AMO(vars...), ec)
		encoding := fac.And(cc...)
		operands = append(operands, encoding)
	}
	for i := 1; i < n-1; i++ {
		vars := make([]f.Variable, n-i)
		for j := 0; j < n-i; j++ {
			vars[j] = varNames[j+i][n-1-j]
		}
		cc, _ := e.EncodeCC(fac, fac.AMO(vars...), ec)
		encoding := fac.And(cc...)
		operands = append(operands, encoding)
	}
	return fac.And(operands...)
}
