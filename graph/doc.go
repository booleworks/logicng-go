// Package graph contains datastructures and algorithms for simple graph and
// hyper-graph implementations in LogicNG.
//
// Hyper-graphs are only used for the BDD Force variable ordering heuristic.
// Graphs however are used in a more versatile manner for constraint graphs of
// formulas.  With constraint graphs you often can split a problem in
// disjunctive sub-problems by computing their connected components.
//
// To generate a constraint graph from a set of formulas, you can call
//
//	constraintGraph := graph.GenerateConstraintGraph(fac, formulas)
//
// You can then compute the connected components and split the
// original formulas with respect to the constraint graph:
//
//	components := graph.ComputeConnectedComponents(constraintGraph)
//	splittedFormulas := graph.SplitFormulasByComponent(fac, formulas, components)
package graph
