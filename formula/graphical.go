package formula

import (
	"fmt"

	"github.com/booleworks/logicng-go/errorx"
	"github.com/booleworks/logicng-go/graphical"
)

const id = "id"

// DefaultFormulaGraphicalGenerator returns a graphical formula generator with
// sensible defaults.
func DefaultFormulaGraphicalGenerator() *graphical.Generator[Formula] {
	return &graphical.Generator[Formula]{
		DefaultNodeStyle: graphical.NoNodeStyle(),
		DefaultEdgeStyle: graphical.NoEdgeStyle(),
	}
}

// GenerateGraphicalFormulaAST generates a graphical representation of the
// formula with the configuration of the generator as an abstract syntax tree
// (AST).  The resulting representation can then be exported as mermaid or
// graphviz graph.
func GenerateGraphicalFormulaAST(
	fac Factory,
	formula Formula,
	generator *graphical.Generator[Formula],
) *graphical.Representation {
	astGenerator := astGenerator{
		Generator:      generator,
		fac:            fac,
		representation: graphical.NewGraphicalRepresentation(generator.AlignTerminals, true, generator.BackgroundColor),
	}
	astGenerator.walkFormula(formula)
	return astGenerator.representation
}

type astGenerator struct {
	*graphical.Generator[Formula]
	fac            Factory
	representation *graphical.Representation
}

func (g *astGenerator) walkFormula(formula Formula) *graphical.Node {
	switch formula.Sort() {
	case SortFalse, SortTrue, SortLiteral:
		return g.walkAtomicFormula(formula)
	case SortCC, SortPBC:
		return g.walkPBConstraint(formula)
	case SortNot:
		return g.walkNotFormula(formula)
	case SortImpl, SortEquiv:
		return g.walkBinaryFormula(formula)
	case SortAnd, SortOr:
		return g.walkNaryFormula(formula)
	default:
		panic(errorx.UnknownEnumValue(formula.Sort()))
	}
}

func (g *astGenerator) walkAtomicFormula(formula Formula) *graphical.Node {
	var label string
	if formula.Sort() == SortLiteral {
		label = litString(g.fac, Literal(formula))
	} else {
		label = formula.Sprint(g.fac)
	}
	return g.addNode(formula, label, true)
}

func (g *astGenerator) walkPBConstraint(pbc Formula) *graphical.Node {
	pbNode := g.addNode(pbc, pbc.Sprint(g.fac), false)
	_, _, literals, _, _ := g.fac.PBCOps(pbc)
	for _, operand := range literals {
		literalNode := g.addNode(operand.AsFormula(), litString(g.fac, operand), true)
		edge := graphical.NewEdge(pbNode, literalNode, g.EdgeStyle(pbc, operand.AsFormula()))
		g.representation.AddEdge(edge)
	}
	return pbNode
}

func (g *astGenerator) walkNotFormula(not Formula) *graphical.Node {
	op, _ := g.fac.NotOperand(not)
	node := g.addNode(not, "¬", false)
	operandNode := g.walkFormula(op)
	edge := graphical.NewEdge(node, operandNode, g.EdgeStyle(not, op))
	g.representation.AddEdge(edge)
	return node
}

func (g *astGenerator) walkBinaryFormula(op Formula) *graphical.Node {
	left, right, _ := g.fac.BinaryLeftRight(op)
	isImpl := op.Sort() == SortImpl
	var label string
	if isImpl {
		label = "⇒"
	} else {
		label = "⇔"
	}

	node := g.addNode(op, label, false)
	leftNode := g.walkFormula(left)
	rightNode := g.walkFormula(right)
	label = ""
	if isImpl {
		label = "l"
	}
	edge := graphical.NewEdge(node, leftNode, g.EdgeStyle(op, left), label)
	g.representation.AddEdge(edge)
	if isImpl {
		label = "r"
	}
	edge = graphical.NewEdge(node, rightNode, g.EdgeStyle(op, right), label)
	g.representation.AddEdge(edge)
	return node
}

func (g *astGenerator) walkNaryFormula(op Formula) *graphical.Node {
	ops, _ := g.fac.NaryOperands(op)
	var label string
	if op.Sort() == SortAnd {
		label = "∧"
	} else {
		label = "∨"
	}
	node := g.addNode(op, label, false)
	for _, operand := range ops {
		operandNode := g.walkFormula(operand)
		edge := graphical.NewEdge(node, operandNode, g.EdgeStyle(op, operand))
		g.representation.AddEdge(edge)
	}
	return node
}

func (g *astGenerator) addNode(formula Formula, defaultLabel string, terminal bool) *graphical.Node {
	nodeId := fmt.Sprintf("%s%d", id, len(g.representation.Nodes()))
	node := graphical.NewNode(nodeId, g.LabelOrDefault(formula, defaultLabel), g.NodeStyle(formula), terminal)
	g.representation.AddNode(node)
	return node
}

func litString(fac Factory, literal Literal) string {
	name, phase, _ := fac.LitNamePhase(literal)
	if phase {
		return name
	}
	return "¬" + name
}

// GenerateGraphicalFormulaDAG generates a graphical representation of the
// formula with the configuration of the generator as a graph (DAG).  The
// resulting representation can then be exported as mermaid or graphviz graph.
func GenerateGraphicalFormulaDAG(
	fac Factory,
	formula Formula,
	generator *graphical.Generator[Formula],
) *graphical.Representation {
	dagGenerator := dagGenerator{
		Generator:      generator,
		fac:            fac,
		representation: graphical.NewGraphicalRepresentation(generator.AlignTerminals, true, generator.BackgroundColor),
		nodes:          map[Formula]*graphical.Node{},
	}
	dagGenerator.initNodes(formula)
	dagGenerator.walkFormula(formula)
	return dagGenerator.representation
}

type dagGenerator struct {
	*graphical.Generator[Formula]
	fac            Factory
	nodes          map[Formula]*graphical.Node
	representation *graphical.Representation
}

func (g *dagGenerator) initNodes(formula Formula) {
	for _, lit := range Literals(g.fac, formula).Content() {
		label := litString(g.fac, lit)
		nodeId := fmt.Sprintf("%s%d", id, len(g.nodes))
		literalNode := graphical.NewNode(nodeId, label, g.NodeStyle(lit.AsFormula()), true)
		g.representation.AddNode(literalNode)
		g.nodes[lit.AsFormula()] = literalNode
	}
}

func (g *dagGenerator) walkFormula(formula Formula) *graphical.Node {
	switch formula.Sort() {
	case SortFalse, SortTrue:
		node, _ := g.addNode(formula, formula.Sprint(g.fac), true)
		return node
	case SortLiteral:
		// since this is a literal, it has to be already present
		return g.nodes[formula]
	case SortCC, SortPBC:
		return g.walkPBConstraint(formula)
	case SortNot:
		return g.walkNotFormula(formula)
	case SortImpl, SortEquiv:
		return g.walkBinaryFormula(formula)
	case SortAnd, SortOr:
		return g.walkNaryFormula(formula)
	default:
		panic(errorx.UnknownEnumValue(formula.Sort()))
	}
}

func (g *dagGenerator) walkPBConstraint(pbc Formula) *graphical.Node {
	node, present := g.addNode(pbc, pbc.Sprint(g.fac), false)
	if !present {
		_, _, ops, _, _ := g.fac.PBCOps(pbc)
		for _, operand := range ops {
			literalNode := g.nodes[operand.AsFormula()]
			g.representation.AddEdge(graphical.NewEdge(node, literalNode, g.EdgeStyle(pbc, operand.AsFormula())))
		}
	}
	return node
}

func (g *dagGenerator) walkNotFormula(not Formula) *graphical.Node {
	node, present := g.addNode(not, "¬", false)
	if !present {
		op, _ := g.fac.NotOperand(not)
		operandNode := g.walkFormula(op)
		g.representation.AddEdge(graphical.NewEdge(node, operandNode, g.EdgeStyle(not, op)))
	}
	return node
}

func (g *dagGenerator) walkBinaryFormula(op Formula) *graphical.Node {
	isImpl := op.Sort() == SortImpl
	var label string
	if isImpl {
		label = "⇒"
	} else {
		label = "⇔"
	}
	node, present := g.addNode(op, label, false)
	if !present {
		left, right, _ := g.fac.BinaryLeftRight(op)
		leftNode := g.walkFormula(left)
		rightNode := g.walkFormula(right)
		label = ""
		if isImpl {
			label = "l"
		}
		g.representation.AddEdge(graphical.NewEdge(node, leftNode, g.EdgeStyle(op, left), label))
		if isImpl {
			label = "r"
		}
		g.representation.AddEdge(graphical.NewEdge(node, rightNode, g.EdgeStyle(op, right), label))
	}
	return node
}

func (g *dagGenerator) walkNaryFormula(op Formula) *graphical.Node {
	var label string
	if op.Sort() == SortAnd {
		label = "∧"
	} else {
		label = "∨"
	}
	node, present := g.addNode(op, label, false)
	if !present {
		ops, _ := g.fac.NaryOperands(op)
		for _, operand := range ops {
			operandNode := g.walkFormula(operand)
			g.representation.AddEdge(graphical.NewEdge(node, operandNode, g.EdgeStyle(op, operand)))
		}
	}
	return node
}

func (g *dagGenerator) addNode(formula Formula, defaultLabel string, terminal bool) (*graphical.Node, bool) {
	node, ok := g.nodes[formula]
	if !ok {
		nodeId := fmt.Sprintf("%s%d", id, len(g.nodes))
		node = graphical.NewNode(nodeId, g.LabelOrDefault(formula, defaultLabel), g.NodeStyle(formula), terminal)
		g.representation.AddNode(node)
		g.nodes[formula] = node
		return node, false
	}
	return node, true
}
