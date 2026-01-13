package graph

import (
	"slices"

	f "github.com/booleworks/logicng-go/formula"
)

// A FormulaGraph represents a graph with nodes holding formulas.
type FormulaGraph struct {
	nodeCache map[f.Formula]int
	nodes     []f.Formula
	adjList   [][]int
}

// NewFormulaGraph generates a new empty formula graph.
func NewFormulaGraph() *FormulaGraph {
	return &FormulaGraph{make(map[f.Formula]int), make([]f.Formula, 0, 8), make([][]int, 0, 8)}
}

// Nodes returns all nodes - and thus formulas - of the graph.
func (g *FormulaGraph) Nodes() []f.Formula {
	return g.nodes
}

// Neighbours returns all formulas connected to the given node in the graph.
func (g *FormulaGraph) Neighbours(node f.Formula) []f.Formula {
	index, ok := g.nodeCache[node]
	if !ok {
		return []f.Formula{}
	}
	nbghs := g.adjList[index]
	result := make([]f.Formula, len(nbghs))
	for i, n := range nbghs {
		result[i] = g.nodes[n]
	}
	return result
}

// AddNode adds a new formula to the graph if it is not already present and
// returns its internal index.
func (g *FormulaGraph) AddNode(formula f.Formula) int {
	idx, ok := g.nodeCache[formula]
	if !ok {
		idx = len(g.nodes)
		g.nodes = append(g.nodes, formula)
		g.nodeCache[formula] = idx
		g.adjList = append(g.adjList, []int{})
	}
	return idx
}

// Connect connects two formula nodes in the graph with an undirected edge.
func (g *FormulaGraph) Connect(n1, n2 f.Formula) {
	node1 := g.AddNode(n1)
	node2 := g.AddNode(n2)
	if node1 != node2 {
		idx, found := slices.BinarySearch(g.adjList[node1], node2)
		if !found {
			g.adjList[node1] = slices.Insert(g.adjList[node1], idx, node2)
		}
		idx, found = slices.BinarySearch(g.adjList[node2], node1)
		if !found {
			g.adjList[node2] = slices.Insert(g.adjList[node2], idx, node1)
		}
	}
}
