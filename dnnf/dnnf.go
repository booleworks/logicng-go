package dnnf

import (
	"math/big"

	f "github.com/booleworks/logicng-go/formula"
)

// A DNNF holds the formula of the d-DNNF and all the original variables (since
// some of them might have been eliminated by the compiler).
type DNNF struct {
	Fac          f.Factory // factory which was used to generate the DNNF
	Formula      f.Formula // formula of the DNNF
	OriginalVars *f.VarSet // set of original variables
}

// ModelCount returns the number of models for the DNNF.
func (d *DNNF) ModelCount() *big.Int {
	cached, ok := f.LookupFunctionCache(d.Fac, f.FuncDNNFModelCount, d.Formula)
	var result *big.Int
	if ok {
		result = new(big.Int).Set(cached.(*big.Int))
	} else {
		result = d.count(d.Formula)
	}
	countDontCares := 0
	dnnfVariables := f.Variables(d.Fac, d.Formula)
	for _, originalVariable := range d.OriginalVars.Content() {
		if !dnnfVariables.Contains(originalVariable) {
			countDontCares++
		}
	}
	factor := big.NewInt(2)
	factor = factor.Exp(factor, big.NewInt(int64(countDontCares)), nil)
	return result.Mul(result, factor)
}

func (d *DNNF) count(dnnf f.Formula) *big.Int {
	cached, ok := f.LookupFunctionCache(d.Fac, f.FuncDNNFModelCount, dnnf)
	if ok {
		return new(big.Int).Set(cached.(*big.Int))
	}
	var c *big.Int
	switch dnnf.Sort() {
	case f.SortLiteral, f.SortTrue:
		c = big.NewInt(1)
	case f.SortAnd:
		c = big.NewInt(1)
		for _, op := range d.Fac.Operands(dnnf) {
			c.Mul(c, d.count(op))
		}
	case f.SortOr:
		allVariables := f.Variables(d.Fac, dnnf).Size()
		c = big.NewInt(0)
		for _, op := range d.Fac.Operands(dnnf) {
			opCount := d.count(op)
			factor := big.NewInt(2)
			factor = factor.Exp(factor, big.NewInt(int64(allVariables-f.Variables(d.Fac, op).Size())), nil)
			mul := new(big.Int).Mul(opCount, factor)
			c.Add(c, mul)
		}
	case f.SortFalse:
		c = big.NewInt(0)
	}
	f.SetFunctionCache(d.Fac, f.FuncDNNFModelCount, dnnf, new(big.Int).Set(c))
	return c
}
