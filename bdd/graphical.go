package bdd

import (
	"fmt"

	"github.com/booleworks/logicng-go/graphical"
)

var (
	id                    = "id"
	defaultTrueNodeStyle  = graphical.Rectangle(graphical.ColorGreen, graphical.ColorWhite, graphical.ColorGreen)
	defaultFalseNodeStyle = graphical.Rectangle(graphical.ColorRed, graphical.ColorWhite, graphical.ColorRed)
	defaultTrueEdgeStyle  = graphical.Solid(graphical.ColorGreen)
	defaultFalseEdgeStyle = graphical.Dotted(graphical.ColorRed)
)

// A GraphicalGenerator is used to configure the graphical representation of
// a BDD as mermaid or Graphviz graphic.  It inherits all fields of
// [graphical.Generator].  The [graphical.Generator.DefaultEdgeStyle] is used
// as the default edge style for the positive BDD edges.  The
// DefaultNegativeEdgeStyle is specific to the BDD generator and is used for
// negative BDD edges.  The same holds for
// [graphical.Generator.ComputeEdgeStyle] and ComputeNegativeEdgeStyle.
type GraphicalGenerator struct {
	*graphical.Generator[int32]
	DefaultNegativeEdgeStyle *graphical.EdgeStyle
	ComputeNegativeEdgeStyle func(src, dst int32) *graphical.EdgeStyle
}

// GenerateGraphical generates a graphical representation of the given bdd
// with the configuration of the generator.  The resulting representation can
// then be exported as mermaid or graphviz graph.
func GenerateGraphical(bdd *BDD, generator *GraphicalGenerator) *graphical.Representation {
	representation := graphical.NewGraphicalRepresentation(generator.AlignTerminals, true, generator.BackgroundColor)
	bddGenerator := bddGenerator{
		GraphicalGenerator: generator,
		representation:     representation,
		index2Node:         map[int32]*graphical.Node{},
	}
	bddGenerator.walkBDD(bdd)
	return bddGenerator.representation
}

type bddGenerator struct {
	*GraphicalGenerator
	representation *graphical.Representation
	index2Node     map[int32]*graphical.Node
}

// DefaultGenerator returns a BDD generator with sensible defaults. Positive
// edges are solid green lines whereas negative edges are dotted red lines.
func DefaultGenerator() *GraphicalGenerator {
	computeNodeStyle := func(index int32) *graphical.NodeStyle {
		switch index {
		case bddFalse:
			return defaultFalseNodeStyle
		case bddTrue:
			return defaultTrueNodeStyle
		default:
			return graphical.NoNodeStyle()
		}
	}
	generator := &graphical.Generator[int32]{
		DefaultNodeStyle: defaultTrueNodeStyle,
		DefaultEdgeStyle: defaultTrueEdgeStyle,
		ComputeNodeStyle: computeNodeStyle,
	}

	return &GraphicalGenerator{
		Generator:                generator,
		DefaultNegativeEdgeStyle: defaultFalseEdgeStyle,
	}
}

func (g *GraphicalGenerator) negativeEdgeStyle(src, dst int32) *graphical.EdgeStyle {
	if g.ComputeNegativeEdgeStyle == nil && g.DefaultNegativeEdgeStyle == nil {
		return graphical.NoEdgeStyle()
	} else if g.ComputeNegativeEdgeStyle == nil {
		return g.DefaultNegativeEdgeStyle
	}
	return g.ComputeNegativeEdgeStyle(src, dst)
}

func (g *bddGenerator) walkBDD(bdd *BDD) {
	if !bdd.IsTautology() {
		falseNode := graphical.NewNode("id0", g.LabelOrDefault(bddFalse, "false"), g.NodeStyle(bddFalse), true)
		g.representation.AddNode(falseNode)
		g.index2Node[bddFalse] = falseNode
	}
	if !bdd.IsContradiction() {
		trueNode := graphical.NewNode("id1", g.LabelOrDefault(bddTrue, "true"), g.NodeStyle(bddTrue), true)
		g.representation.AddNode(trueNode)
		g.index2Node[bddTrue] = trueNode
	}
	for _, internalNode := range bdd.Kernel.allNodes(bdd.Index) {
		index := internalNode[0]
		variable, _ := bdd.Kernel.getVariableForIndex(internalNode[1])
		defaultLabel, _ := bdd.Kernel.fac.VarName(variable)
		g.addNode(index, g.LabelOrDefault(index, defaultLabel))
	}
	for _, internalNode := range bdd.Kernel.allNodes(bdd.Index) {
		index := internalNode[0]
		lowIndex := internalNode[2]
		highIndex := internalNode[3]
		node := g.index2Node[index]
		lowNode := g.index2Node[lowIndex]
		highNode := g.index2Node[highIndex]
		g.representation.AddEdge(graphical.NewEdge(node, lowNode, g.negativeEdgeStyle(index, lowIndex)))
		g.representation.AddEdge(graphical.NewEdge(node, highNode, g.EdgeStyle(index, highIndex)))
	}
}

func (g *bddGenerator) addNode(index int32, label string) {
	_, ok := g.index2Node[index]
	if !ok {
		nodeId := fmt.Sprintf("%s%d", id, index)
		node := graphical.NewNode(nodeId, label, g.NodeStyle(index), false)
		g.representation.AddNode(node)
		g.index2Node[index] = node
	}
}
