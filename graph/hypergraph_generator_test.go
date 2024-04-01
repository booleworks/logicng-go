package graph

import (
	"testing"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/parser"
	"github.com/stretchr/testify/assert"
)

func TestHypergraphFromCnf(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	a := fac.Var("A")
	b := fac.Var("B")
	c := fac.Var("C")
	d := fac.Var("D")
	e := fac.Var("E")
	x := fac.Var("X")
	y := fac.Var("Y")
	p := parser.New(fac)

	graph, _ := HypergraphFromCNF(fac, p.ParseUnsafe("$false"))
	assert.Equal(0, len(graph.Nodes))
	assert.Equal(0, len(graph.Edges))

	graph, _ = HypergraphFromCNF(fac, p.ParseUnsafe("$true"))
	assert.Equal(0, len(graph.Nodes))
	assert.Equal(0, len(graph.Edges))

	graph, _ = HypergraphFromCNF(fac, p.ParseUnsafe("A"))
	assert.Equal(1, len(graph.Nodes))
	assert.Equal(a, graph.Nodes[0].Content)
	assert.Equal(1, len(graph.Edges))
	assert.Equal(1, len(graph.Edges[0].Nodes))
	assert.Equal(a, graph.Edges[0].Nodes[0].Content)

	graph, _ = HypergraphFromCNF(fac, p.ParseUnsafe("A | B | ~C"))
	assert.Equal(3, len(graph.Nodes))
	assert.Equal(a, graph.Nodes[0].Content)
	assert.Equal(b, graph.Nodes[1].Content)
	assert.Equal(c, graph.Nodes[2].Content)
	assert.Equal(1, len(graph.Edges))
	assert.Equal(3, len(graph.Edges[0].Nodes))
	assert.Equal(a, graph.Edges[0].Nodes[0].Content)
	assert.Equal(b, graph.Edges[0].Nodes[1].Content)
	assert.Equal(c, graph.Edges[0].Nodes[2].Content)

	graph, _ = HypergraphFromCNF(fac, p.ParseUnsafe("(A | B | ~C) & (B | ~D) & (C | ~E) & (~B | ~D | E) & X & ~Y"))
	assert.Equal(7, len(graph.Nodes))
	assert.Equal(a, graph.Nodes[0].Content)
	assert.Equal(b, graph.Nodes[1].Content)
	assert.Equal(c, graph.Nodes[2].Content)
	assert.Equal(d, graph.Nodes[3].Content)
	assert.Equal(e, graph.Nodes[4].Content)
	assert.Equal(x, graph.Nodes[5].Content)
	assert.Equal(y, graph.Nodes[6].Content)

	assert.Equal(6, len(graph.Edges))
	assert.Equal(3, len(graph.Edges[0].Nodes))
	assert.Equal(a, graph.Edges[0].Nodes[0].Content)
	assert.Equal(b, graph.Edges[0].Nodes[1].Content)
	assert.Equal(c, graph.Edges[0].Nodes[2].Content)
	assert.Equal(2, len(graph.Edges[1].Nodes))
	assert.Equal(b, graph.Edges[1].Nodes[0].Content)
	assert.Equal(d, graph.Edges[1].Nodes[1].Content)
	assert.Equal(2, len(graph.Edges[2].Nodes))
	assert.Equal(c, graph.Edges[2].Nodes[0].Content)
	assert.Equal(e, graph.Edges[2].Nodes[1].Content)
	assert.Equal(3, len(graph.Edges[3].Nodes))
	assert.Equal(b, graph.Edges[3].Nodes[0].Content)
	assert.Equal(d, graph.Edges[3].Nodes[1].Content)
	assert.Equal(e, graph.Edges[3].Nodes[2].Content)
	assert.Equal(1, len(graph.Edges[4].Nodes))
	assert.Equal(x, graph.Edges[4].Nodes[0].Content)
	assert.Equal(1, len(graph.Edges[5].Nodes))
	assert.Equal(y, graph.Edges[5].Nodes[0].Content)
}
