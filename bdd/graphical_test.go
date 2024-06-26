package bdd

import (
	"os"
	"testing"

	"github.com/booleworks/logicng-go/graphical"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestWriteGraphicalBDDFormulas(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	ordering := []f.Variable{fac.Var("A"), fac.Var("B"), fac.Var("C"), fac.Var("D")}
	kernel := NewKernelWithOrdering(fac, ordering, 1000, 1000)

	testFiles(t, "false", GenerateGraphical(CompileWithKernel(fac, p.ParseUnsafe("$false"), kernel), DefaultGenerator()))
	testFiles(t, "true", GenerateGraphical(CompileWithKernel(fac, p.ParseUnsafe("$true"), kernel), DefaultGenerator()))
	testFiles(t, "a", GenerateGraphical(CompileWithKernel(fac, p.ParseUnsafe("A"), kernel), DefaultGenerator()))
	testFiles(t, "not_a", GenerateGraphical(CompileWithKernel(fac, p.ParseUnsafe("~A"), kernel), DefaultGenerator()))
	testFiles(t, "impl", GenerateGraphical(CompileWithKernel(fac, p.ParseUnsafe("A => ~C"), kernel), DefaultGenerator()))
	testFiles(t, "equiv", GenerateGraphical(CompileWithKernel(fac, p.ParseUnsafe("A <=> ~C"), kernel), DefaultGenerator()))
	testFiles(t, "or", GenerateGraphical(CompileWithKernel(fac, p.ParseUnsafe("A | B | ~C"), kernel), DefaultGenerator()))
	testFiles(t, "and", GenerateGraphical(CompileWithKernel(fac, p.ParseUnsafe("A & B & ~C"), kernel), DefaultGenerator()))
	testFiles(t, "not", GenerateGraphical(CompileWithKernel(fac, p.ParseUnsafe("~(A & B & ~C)"), kernel), DefaultGenerator()))
	bdd := CompileWithKernel(fac, p.ParseUnsafe("(A => (B|~C)) & (B => C & D) & (D <=> A)"), kernel)
	testFiles(t, "formula", GenerateGraphical(bdd, DefaultGenerator()))
}

func TestWriteGraphicalBDDFixedStyle(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	ordering := []f.Variable{fac.Var("A"), fac.Var("B"), fac.Var("C"), fac.Var("D")}
	kernel := NewKernelWithOrdering(fac, ordering, 1000, 1000)
	bdd := CompileWithKernel(fac, p.ParseUnsafe("(A => (B|~C)) & (B => C & D) & (D <=> A)"), kernel)

	defaultNodeStyle := graphical.Circle(graphical.ColorOrange, graphical.ColorBlack, graphical.ColorOrange)
	positiveNodeStyle := graphical.Rectangle(graphical.ColorCyan, graphical.ColorWhite, graphical.ColorCyan)
	negativeNodeStyle := graphical.Rectangle(graphical.ColorPurple, graphical.ColorWhite, graphical.ColorPurple)
	computeNodeStyle := func(index int32) *graphical.NodeStyle {
		switch index {
		case bddFalse:
			return negativeNodeStyle
		case bddTrue:
			return positiveNodeStyle
		default:
			return defaultNodeStyle
		}
	}
	generator := &GraphicalGenerator{
		Generator: &graphical.Generator[int32]{
			BackgroundColor:  graphical.ColorLightGray,
			AlignTerminals:   true,
			DefaultEdgeStyle: graphical.Bold(graphical.ColorCyan),
			ComputeNodeStyle: computeNodeStyle,
		},
		DefaultNegativeEdgeStyle: graphical.Dotted(graphical.ColorPurple),
	}

	testFiles(t, "formula-fixedStyle", GenerateGraphical(bdd, generator))
}

func testFiles(t *testing.T, fileName string, representation *graphical.Representation) {
	mermaid := graphical.WriteMermaidToString(representation)
	dot := graphical.WriteDotToString(representation)

	expectedDot, _ := os.ReadFile("../test/data/graphical/bdd/" + fileName + "_bdd.dot")
	expected := string(expectedDot)
	assert.Equal(t, expected, dot)

	expectedMermaid, _ := os.ReadFile("../test/data/graphical/bdd/" + fileName + "_bdd.txt")
	expected = string(expectedMermaid)
	assert.Equal(t, expected, mermaid)
}
