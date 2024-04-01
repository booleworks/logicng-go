package graph

import (
	"fmt"

	"github.com/booleworks/logicng-go/graphical"

	f "github.com/booleworks/logicng-go/formula"
)

const id = "id"

// GenerateGraphicalFormulaGraph generates a graphical representation of the
// graph with the configuration of the generator.  The resulting representation
// can then be exported as mermaid or graphviz graph.
func GenerateGraphicalFormulaGraph(
	fac f.Factory, graph *FormulaGraph, generator *graphical.Generator[f.Formula],
) *graphical.Representation {
	graphGenerator := graphGenerator{
		Generator:      generator,
		fac:            fac,
		representation: graphical.NewGraphicalRepresentation(false, false, generator.BackgroundColor),
		nodes:          map[f.Formula]*graphical.Node{},
		visited:        map[f.Formula]present{},
	}
	graphGenerator.walkGraph(graph)
	return graphGenerator.representation
}

type graphGenerator struct {
	*graphical.Generator[f.Formula]
	fac            f.Factory
	representation *graphical.Representation
	nodes          map[f.Formula]*graphical.Node
	visited        map[f.Formula]present
}

func (g *graphGenerator) walkGraph(graph *FormulaGraph) {
	for _, node := range graph.Nodes() {
		graphicalNode := g.addNode(node)
		for _, neighbour := range graph.Neighbours(node) {
			graphicalNeighbourNode := g.addNode(neighbour)
			if _, ok := g.visited[neighbour]; !ok {
				g.representation.AddEdge(graphical.NewEdge(graphicalNode, graphicalNeighbourNode, g.EdgeStyle(node, neighbour)))
			}
		}
		g.visited[node] = present{}
	}
}

func (g *graphGenerator) addNode(node f.Formula) *graphical.Node {
	graphicalNode, ok := g.nodes[node]
	if !ok {
		nodeId := fmt.Sprintf("%s%d", id, len(g.nodes))
		graphicalNode = graphical.NewNode(nodeId, g.LabelOrDefault(node, node.Sprint(g.fac)), g.NodeStyle(node))
		g.representation.AddNode(graphicalNode)
		g.nodes[node] = graphicalNode
	}
	return graphicalNode
}

type present struct{}
