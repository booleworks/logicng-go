package count

import (
	"math/big"

	"booleworks.com/logicng/assignment"
	"booleworks.com/logicng/dnnf"
	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/graph"
	"booleworks.com/logicng/normalform"
	"booleworks.com/logicng/transformation"
)

// Count computes the model count for the given formulas (interpreted as
// conjunction) and a set of relevant variables. This set can only be a
// superset of the original formulas' variables. No projected model counting is
// supported.  This is just used for don't care variable detection.
//
// Since the counter internally has to generate a CNF formula which does not
// alter the model count, only AMO and EXO cardinality constraints can be
// counted - if there are arbitrary cardinality constraints or pseudo-Boolean
// constraints in the formula, an error is returned.
func Count(fac f.Factory, variables []f.Variable, formulas ...f.Formula) (*big.Int, error) {
	vars := f.NewVarSet(variables...)
	if !vars.ContainsAll(f.Variables(fac, formulas...)) {
		panic(errorx.BadInput("variables must be a superset of the formulas' variables"))
	}
	if vars.Empty() {
		nonTrueCount := 0
		for _, formula := range formulas {
			if formula.Sort() != f.SortTrue {
				nonTrueCount++
				break
			}
		}
		if nonTrueCount == 0 {
			return big.NewInt(1), nil
		} else {
			return big.NewInt(0), nil
		}
	}
	cnfs, err := encodeAsCNF(fac, formulas)
	if err != nil {
		return nil, err
	}
	simplification := simplify(fac, cnfs)
	count := count(fac, simplification.simplifiedFormulas)
	dontCareVariables := simplification.getDontCareVariables(vars)
	exp := big.NewInt(2).Exp(big.NewInt(2), big.NewInt(int64(dontCareVariables.Size())), nil)
	return count.Mul(count, exp), nil
}

func encodeAsCNF(fac f.Factory, formulas []f.Formula) ([]f.Formula, error) {
	cnf := make([]f.Formula, len(formulas))
	for i, formula := range formulas {
		exp, err := transformation.ExpandAMOAndEXO(fac, formula)
		if err != nil {
			return nil, err
		}
		cnf[i] = exp
		cnf[i] = normalform.CNF(fac, cnf[i], normalform.DefaultCNFConfig())
	}
	return cnf, nil
}

func simplify(fac f.Factory, formulas []f.Formula) *simplificationResult {
	simpleBackbone := assignment.Empty()
	backboneVariables := f.NewVarSet()
	for _, formula := range formulas {
		if formula.Sort() == f.SortLiteral {
			literal := f.Literal(formula)
			_ = simpleBackbone.AddLit(fac, literal)
			variable := literal.Variable()
			backboneVariables.Add(variable)
		}
	}
	simplified := make([]f.Formula, 0, len(formulas))
	for _, formula := range formulas {
		restrict := assignment.Restrict(fac, formula, simpleBackbone)
		if restrict.Sort() != f.SortTrue {
			simplified = append(simplified, restrict)
		}
	}
	return &simplificationResult{fac, simplified, backboneVariables.Content()}
}

func count(fac f.Factory, formulas []f.Formula) *big.Int {
	constraintGraph := graph.GenerateConstraintGraph(fac, formulas...)
	ccs := graph.ComputeConnectedComponents(constraintGraph)
	components := graph.SplitFormulasByComponent(fac, formulas, ccs)
	count := big.NewInt(1)
	for _, component := range components {
		dnnf := dnnf.Compile(fac, fac.And(component...))
		dnnfCount := dnnf.ModelCount()
		count = count.Mul(count, dnnfCount)
	}
	return count
}

type simplificationResult struct {
	fac                f.Factory
	simplifiedFormulas []f.Formula
	backboneVariables  []f.Variable
}

func (s *simplificationResult) getDontCareVariables(variables *f.VarSet) *f.VarSet {
	dontCareVariables := f.NewVariableSetCopy(variables)
	dontCareVariables.RemoveAll(f.Variables(s.fac, s.simplifiedFormulas...))
	dontCareVariables.RemoveAllElements(&s.backboneVariables)
	return dontCareVariables
}
