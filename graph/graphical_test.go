package graph

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/booleworks/logicng-go/graphical"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestGraphWriterSmallDefault(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	f1 := p.ParseUnsafe("A | B")
	f2 := p.ParseUnsafe("C")
	g := GenerateConstraintGraph(fac, f1, f2)

	testGraphFiles(t, "small", GenerateGraphicalFormulaGraph(fac, g, f.DefaultFormulaGraphicalGenerator()))
}

func TestGraphWriterSmallFixedStyle(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	f1 := p.ParseUnsafe("A | B")
	f2 := p.ParseUnsafe("C")
	g := GenerateConstraintGraph(fac, f1, f2)
	generator := &graphical.Generator[f.Formula]{
		BackgroundColor:  graphical.Color("#4f4f4f"),
		DefaultEdgeStyle: graphical.Dotted(graphical.ColorWhite),
		DefaultNodeStyle: graphical.NewNodeStyle(graphical.ShapeDefault, graphical.ColorRed, graphical.ColorGreen, ""),
		AlignTerminals:   true,
	}

	testGraphFiles(t, "small-fixedStyle", GenerateGraphicalFormulaGraph(fac, g, generator))
}

func TestGraphWriter30(t *testing.T) {
	fac := f.NewFactory()
	g := readGraph(fac, "30")
	for i := 0; i < 30; i++ {
		g.AddNode(fac.Variable(fmt.Sprintf("%d", i)))
	}

	testGraphFiles(t, "30", GenerateGraphicalFormulaGraph(fac, g, f.DefaultFormulaGraphicalGenerator()))
}

func TestWriteGraph30DynamicStyle(t *testing.T) {
	fac := f.NewFactory()
	g := readGraph(fac, "30")
	for i := 0; i < 30; i++ {
		g.AddNode(fac.Variable(fmt.Sprintf("%d", i)))
	}

	style1 := graphical.Rectangle(graphical.ColorGreen, graphical.ColorBlack, graphical.ColorGreen)
	style2 := graphical.Ellipse(graphical.ColorOrange, graphical.ColorBlack, graphical.ColorOrange)
	style3 := graphical.Circle(graphical.ColorRed, graphical.ColorWhite, graphical.ColorRed)

	computeNodeStyle := func(content f.Formula) *graphical.NodeStyle {
		name, _, _ := fac.LiteralNamePhase(content)
		l, _ := strconv.Atoi(name)
		if l <= 10 {
			return style1
		} else if l <= 20 {
			return style2
		} else {
			return style3
		}
	}

	eStyle1 := graphical.NewEdgeStyle(graphical.EdgeDefault, graphical.ColorGreen)
	eStyle2 := graphical.Solid(graphical.ColorOrange)
	eStyle3 := graphical.Dotted(graphical.ColorLightGray)

	computeEdgeStyle := func(src, dst f.Formula) *graphical.EdgeStyle {
		nameSrc, _, _ := fac.LiteralNamePhase(src)
		nameDst, _, _ := fac.LiteralNamePhase(dst)
		l1, _ := strconv.Atoi(nameSrc)
		l2, _ := strconv.Atoi(nameDst)
		if l1 <= 10 && l2 <= 10 {
			return eStyle1
		} else if l1 <= 20 && l2 <= 20 {
			return eStyle2
		} else {
			return eStyle3
		}
	}

	computeLabel := func(content f.Formula) string { return fmt.Sprintf("value: %s", content.Sprint(fac)) }

	generator := &graphical.Generator[f.Formula]{
		ComputeLabel:     computeLabel,
		ComputeNodeStyle: computeNodeStyle,
		ComputeEdgeStyle: computeEdgeStyle,
	}

	testGraphFiles(t, "30-dynamic", GenerateGraphicalFormulaGraph(fac, g, generator))
}

func testGraphFiles(t *testing.T, fileName string, representation *graphical.Representation) {
	mermaid := graphical.WriteMermaidToString(representation)
	dot := graphical.WriteDotToString(representation)

	expectedDot, _ := os.ReadFile("../test/data/graphical/graph/" + fileName + ".dot")
	expected := string(expectedDot)
	assert.Equal(t, expected, dot)

	expectedMermaid, _ := os.ReadFile("../test/data/graphical/graph/" + fileName + ".txt")
	expected = string(expectedMermaid)
	assert.Equal(t, expected, mermaid)
}
