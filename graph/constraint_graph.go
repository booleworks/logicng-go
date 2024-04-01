package graph

import f "booleworks.com/logicng/formula"

// GenerateConstraintGraph generates a constraint graph for the given formulas.
// The nodes of the constraint graph hold all variables of the given formulas.
// Two variable nodes are connected if they occur in the same formula.
func GenerateConstraintGraph(fac f.Factory, formulas ...f.Formula) *FormulaGraph {
	constraintGraph := NewFormulaGraph()
	for _, formula := range formulas {
		addSubformula(fac, formula, constraintGraph)
	}
	return constraintGraph
}

func addSubformula(fac f.Factory, formula f.Formula, graph *FormulaGraph) {
	variables := f.Variables(fac, formula).Content()
	if len(variables) == 1 {
		graph.AddNode(variables[0].AsFormula())
	}
	for i := 0; i < len(variables); i++ {
		for j := i + 1; j < len(variables); j++ {
			graph.Connect(variables[i].AsFormula(), variables[j].AsFormula())
		}
	}
}
