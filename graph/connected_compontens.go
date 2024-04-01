package graph

import (
	"slices"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// ComputeConnectedComponents returns the strongly connected components of a
// constraint graph.  Each component is represented as a slice of variables.
func ComputeConnectedComponents(graph *FormulaGraph) [][]f.Variable {
	var connectedComponents [][]f.Variable
	if len(graph.nodes) == 0 {
		return connectedComponents
	}
	marked := make([]bool, len(graph.nodes))
	for idxFalse := 0; idxFalse != -1; idxFalse = slices.Index(marked, false) {
		var connectedComp []f.Variable
		deepFirstSearch(graph, idxFalse, &connectedComp, &marked)
		connectedComponents = append(connectedComponents, connectedComp)
	}
	return connectedComponents
}

func deepFirstSearch(graph *FormulaGraph, v int, component *[]f.Variable, marked *[]bool) {
	*component = append(*component, f.Variable(graph.nodes[v]))
	(*marked)[v] = true
	for _, neigh := range graph.adjList[v] {
		if !(*marked)[neigh] {
			deepFirstSearch(graph, neigh, component, marked)
		}
	}
}

// SplitFormulasByComponent splits a given list of formulas with respect to a
// given list of components, usually from the constraint graph of the given
// formulas.  The result contains one list of formulas for each of the given
// components.  All formulas in this list have only variables from the
// respective component.
func SplitFormulasByComponent(fac f.Factory, formulas []f.Formula, components [][]f.Variable) [][]f.Formula {
	result := make([][]f.Formula, len(components))
	varMap := make(map[f.Variable]int)
	for i, component := range components {
		result[i] = make([]f.Formula, 0)
		for _, variable := range component {
			varMap[variable] = i
		}
	}
	for _, formula := range formulas {
		variables := f.Variables(fac, formula)
		if !variables.Empty() {
			anyVar, _ := variables.Any()
			componentId, ok := varMap[anyVar]
			if !ok {
				panic(errorx.BadInput("no component for variable %s in the graph", anyVar.Sprint(fac)))
			}
			result[componentId] = append(result[componentId], formula)
		} else if formula.IsNeg() {
			result = append(result, []f.Formula{fac.Falsum()})
		}
	}
	return result
}
