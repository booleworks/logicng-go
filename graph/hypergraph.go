package graph

import (
	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
	"github.com/emirpasic/gods/maps/linkedhashmap"
)

// A HypergraphNode represents a node holding a variable in a hyper-graph.  It
// holds also a list of edges on which it occurs.
type HypergraphNode struct {
	Content f.Variable
	Edges   []*HypergraphEdge
}

// A HypergraphEdge connects nodes in a hyper-graph.
type HypergraphEdge struct {
	Nodes []*HypergraphNode
}

// A Hypergraph holds its nodes and edges.
type Hypergraph struct {
	Nodes []*HypergraphNode
	Edges []*HypergraphEdge
}

// NewHypergraph generates a new empty hyper-graph.
func NewHypergraph() *Hypergraph {
	return &Hypergraph{
		Nodes: make([]*HypergraphNode, 0),
		Edges: make([]*HypergraphEdge, 0),
	}
}

// AddEdge adds an edge between the given nodes.
func (g *Hypergraph) AddEdge(nodes []*HypergraphNode) {
	edge := NewHypergraphEdge(nodes)
	g.Edges = append(g.Edges, edge)
}

// NewHypergraphNode returns a new node in the given graph with the given
// variable as content.
func NewHypergraphNode(graph *Hypergraph, variable f.Variable) *HypergraphNode {
	node := &HypergraphNode{
		Content: variable,
		Edges:   make([]*HypergraphEdge, 0),
	}
	graph.Nodes = append(graph.Nodes, node)
	return node
}

// NewHypergraphEdge returns a new edge connecting the given nodes.
func NewHypergraphEdge(nodes []*HypergraphNode) *HypergraphEdge {
	edge := &HypergraphEdge{nodes}
	for _, node := range nodes {
		node.Edges = append(node.Edges, edge)
	}
	return edge
}

// ComputeTentativeNewLocation is used in the Force orderings for BDDs.
func (n *HypergraphNode) ComputeTentativeNewLocation(nodeOrdering *linkedhashmap.Map) float64 {
	newLocation := .0
	for _, edge := range n.Edges {
		newLocation += edge.centerOfGravity(nodeOrdering)
	}
	return newLocation / float64(len(n.Edges))
}

func (e *HypergraphEdge) centerOfGravity(nodeOrdering *linkedhashmap.Map) float64 {
	cog := 0
	for _, node := range e.Nodes {
		level, ok := nodeOrdering.Get(node)
		if !ok {
			panic(errorx.IllegalState("could not find node %v in the node ordering", node))
		}
		cog += level.(int)
	}
	return float64(cog) / float64(len(e.Nodes))
}
