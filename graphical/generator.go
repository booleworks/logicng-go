package graphical

// A Generator configures the graphical representation of graphs and trees for
// exporting them to mermaid.js or Graphviz.
type Generator[T any] struct {
	BackgroundColor  Color
	AlignTerminals   bool
	DefaultEdgeStyle *EdgeStyle
	DefaultNodeStyle *NodeStyle
	ComputeNodeStyle func(content T) *NodeStyle
	ComputeLabel     func(content T) string
	ComputeEdgeStyle func(src, dst T) *EdgeStyle
}

// NodeStyle returns the node style for the given node.
func (g *Generator[T]) NodeStyle(node T) *NodeStyle {
	if g.ComputeNodeStyle == nil && g.DefaultNodeStyle == nil {
		return NoNodeStyle()
	} else if g.ComputeNodeStyle == nil {
		return g.DefaultNodeStyle
	} else {
		return g.ComputeNodeStyle(node)
	}
}

// LabelOrDefault returns the label for the given node or a given defaultLabel
// when there is no label function in the generator.
func (g *Generator[T]) LabelOrDefault(node T, defaultLabel string) string {
	if g.ComputeLabel == nil {
		return defaultLabel
	} else {
		return g.ComputeLabel(node)
	}
}

// EdgeStyle returns the edge style for the edge between the given src and dst
// node.
func (g *Generator[T]) EdgeStyle(src, dst T) *EdgeStyle {
	if g.ComputeEdgeStyle == nil && g.DefaultEdgeStyle == nil {
		return NoEdgeStyle()
	} else if g.ComputeEdgeStyle == nil {
		return g.DefaultEdgeStyle
	} else {
		return g.ComputeEdgeStyle(src, dst)
	}
}
