package test

import (
	"fmt"
	"os"
	"testing"
	"unicode"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/graphical"

	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestAstWriterConstants(t *testing.T) {
	fac := f.NewFactory()

	testFilesAst(t, "false", f.GenerateGraphicalFormulaAST(fac, fac.Falsum(), f.DefaultFormulaGraphicalGenerator()))
	testFilesAst(t, "true", f.GenerateGraphicalFormulaAST(fac, fac.Verum(), f.DefaultFormulaGraphicalGenerator()))
}

func TestAstWriterLiterals(t *testing.T) {
	fac := f.NewFactory()

	testFilesAst(t, "x", f.GenerateGraphicalFormulaAST(fac, fac.Variable("x"), f.DefaultFormulaGraphicalGenerator()))
	testFilesAst(t, "not_x", f.GenerateGraphicalFormulaAST(fac, fac.Literal("x", false), f.DefaultFormulaGraphicalGenerator()))
}

func TestAstWriterFormulas(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)

	f1 := p.ParseUnsafe("(a & b) <=> (~c => (x | z))")
	f2 := p.ParseUnsafe("a & b | b & ~c")
	f3 := p.ParseUnsafe("(a & b) <=> (~c => (a | b))")
	f4 := p.ParseUnsafe("~(a & b) | b & ~c")
	f5 := p.ParseUnsafe("a | ~b | (2*a + 3*~b + 4*c <= 23)")

	testFilesAst(t, "f1", f.GenerateGraphicalFormulaAST(fac, f1, f.DefaultFormulaGraphicalGenerator()))
	testFilesAst(t, "f2", f.GenerateGraphicalFormulaAST(fac, f2, f.DefaultFormulaGraphicalGenerator()))
	testFilesAst(t, "f3", f.GenerateGraphicalFormulaAST(fac, f3, f.DefaultFormulaGraphicalGenerator()))
	testFilesAst(t, "f4", f.GenerateGraphicalFormulaAST(fac, f4, f.DefaultFormulaGraphicalGenerator()))
	testFilesAst(t, "f5", f.GenerateGraphicalFormulaAST(fac, f5, f.DefaultFormulaGraphicalGenerator()))
}

func TestAstWriterDuplicateFormulaParts(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	f6 := p.ParseUnsafe("(a & b) | (c & ~(a & b))")
	f7 := p.ParseUnsafe("(c & d) | (a & b) | ((c & d) <=> (a & b))")

	testFilesAst(t, "f6", f.GenerateGraphicalFormulaAST(fac, f6, f.DefaultFormulaGraphicalGenerator()))
	testFilesAst(t, "f7", f.GenerateGraphicalFormulaAST(fac, f7, f.DefaultFormulaGraphicalGenerator()))
}

func TestAstWriterFixedStyle(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	f8 := p.ParseUnsafe("(A <=> B & (~A | C | X)) => a + b + c <= 2")
	generator := &graphical.Generator[f.Formula]{
		BackgroundColor:  graphical.Color("#020202"),
		DefaultEdgeStyle: graphical.Bold(graphical.ColorCyan),
		DefaultNodeStyle: graphical.Circle(graphical.ColorBlue, graphical.ColorWhite, graphical.ColorBlue),
		AlignTerminals:   true,
	}

	testFilesAst(t, "f8", f.GenerateGraphicalFormulaAST(fac, f8, generator))
}

func TestAstWriterDynamicStyle(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	f9 := p.ParseUnsafe("(A <=> B & (~A | C | X)) => a + b + c <= 2 & (~a | d => X & ~B)")

	style1 := graphical.Rectangle(graphical.ColorDarkGray, graphical.ColorDarkGray, graphical.ColorLightGray)
	style2 := graphical.Circle(graphical.ColorYellow, graphical.ColorBlack, graphical.ColorYellow)
	style3 := graphical.Circle(graphical.ColorTurquoise, graphical.ColorWhite, graphical.ColorTurquoise)
	style4 := graphical.Ellipse(graphical.ColorBlack, graphical.ColorBlack, "")

	computeNodeStyle := func(content f.Formula) *graphical.NodeStyle {
		switch content.Sort() {
		case f.SortCC, f.SortPBC:
			return style1
		case f.SortLiteral:
			name, _, _ := fac.LiteralNamePhase(content)
			if unicode.IsLower(rune(name[0])) {
				return style2
			}
			return style3
		default:
			return style4
		}
	}
	generator := &graphical.Generator[f.Formula]{
		BackgroundColor:  graphical.Color("#444444"),
		DefaultEdgeStyle: graphical.NoEdgeStyle(),
		ComputeNodeStyle: computeNodeStyle,
	}

	testFilesAst(t, "f9", f.GenerateGraphicalFormulaAST(fac, f9, generator))
}

func TestAstWriterEdgeMapper(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	f10 := p.ParseUnsafe("(A <=> B & (~A | C | X)) => a + b + c <= 2 & (~a | d => X & ~B)")

	style1 := graphical.Dotted(graphical.ColorDarkGray)
	style2 := graphical.Solid(graphical.ColorBlack)

	computeEdgeStyle := func(src, _ f.Formula) *graphical.EdgeStyle {
		if src.Sort() == f.SortPBC || src.Sort() == f.SortCC {
			return style1
		}
		return style2
	}
	generator := &graphical.Generator[f.Formula]{
		BackgroundColor:  graphical.Color("#444444"),
		DefaultEdgeStyle: graphical.Solid(graphical.ColorPurple),
		ComputeEdgeStyle: computeEdgeStyle,
	}

	testFilesAst(t, "f10", f.GenerateGraphicalFormulaAST(fac, f10, generator))
}

func TestAstWriterLabelMapper(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	f8 := p.ParseUnsafe("(A <=> B & (~A | C | X)) => a + b + c <= 2")
	computeLabel := func(content f.Formula) string {
		return content.Sprint(fac)
	}
	generator := &graphical.Generator[f.Formula]{
		BackgroundColor:  graphical.Color("#020202"),
		DefaultEdgeStyle: graphical.Bold(graphical.ColorCyan),
		DefaultNodeStyle: graphical.Rectangle(graphical.ColorBlue, graphical.ColorWhite, graphical.ColorBlue),
		AlignTerminals:   true,
		ComputeLabel:     computeLabel,
	}

	testFilesAst(t, "f8-ownLabels", f.GenerateGraphicalFormulaAST(fac, f8, generator))
}

func TestDagWriterConstants(t *testing.T) {
	fac := f.NewFactory()

	testFilesDag(t, "false", f.GenerateGraphicalFormulaDAG(fac, fac.Falsum(), f.DefaultFormulaGraphicalGenerator()))
	testFilesDag(t, "true", f.GenerateGraphicalFormulaDAG(fac, fac.Verum(), f.DefaultFormulaGraphicalGenerator()))
}

func TestDagWriterLiterals(t *testing.T) {
	fac := f.NewFactory()

	testFilesDag(t, "x", f.GenerateGraphicalFormulaDAG(fac, fac.Variable("x"), f.DefaultFormulaGraphicalGenerator()))
	testFilesDag(t, "not_x", f.GenerateGraphicalFormulaDAG(fac, fac.Literal("x", false), f.DefaultFormulaGraphicalGenerator()))
}

func TestDagWriterFormulas(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)

	f1 := p.ParseUnsafe("(a & b) <=> (~c => (x | z))")
	f2 := p.ParseUnsafe("a & b | b & ~c")
	f3 := p.ParseUnsafe("(a & b) <=> (~c => (a | b))")
	f4 := p.ParseUnsafe("~(a & b) | b & ~c")
	f5 := p.ParseUnsafe("a | ~b | (2*a + 3*~b + 4*c <= 23)")

	testFilesDag(t, "f1", f.GenerateGraphicalFormulaDAG(fac, f1, f.DefaultFormulaGraphicalGenerator()))
	testFilesDag(t, "f2", f.GenerateGraphicalFormulaDAG(fac, f2, f.DefaultFormulaGraphicalGenerator()))
	testFilesDag(t, "f3", f.GenerateGraphicalFormulaDAG(fac, f3, f.DefaultFormulaGraphicalGenerator()))
	testFilesDag(t, "f4", f.GenerateGraphicalFormulaDAG(fac, f4, f.DefaultFormulaGraphicalGenerator()))
	testFilesDag(t, "f5", f.GenerateGraphicalFormulaDAG(fac, f5, f.DefaultFormulaGraphicalGenerator()))
}

func TestDagWriterDuplicateFormulaParts(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	f6 := p.ParseUnsafe("(a & b) | (c & ~(a & b))")
	f7 := p.ParseUnsafe("(c & d) | (a & b) | ((c & d) <=> (a & b))")

	testFilesDag(t, "f6", f.GenerateGraphicalFormulaDAG(fac, f6, f.DefaultFormulaGraphicalGenerator()))
	testFilesDag(t, "f7", f.GenerateGraphicalFormulaDAG(fac, f7, f.DefaultFormulaGraphicalGenerator()))
}

func TestDagWriterFixedStyle(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	f8 := p.ParseUnsafe("(A <=> B & (~A | C | X)) => a + b + c <= 2")
	generator := &graphical.Generator[f.Formula]{
		BackgroundColor:  graphical.Color("#020202"),
		DefaultEdgeStyle: graphical.Bold(graphical.ColorCyan),
		DefaultNodeStyle: graphical.Circle(graphical.ColorBlue, graphical.ColorWhite, graphical.ColorBlue),
		AlignTerminals:   true,
	}

	testFilesDag(t, "f8", f.GenerateGraphicalFormulaDAG(fac, f8, generator))
}

func TestDagWriterDynamicStyle(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	f9 := p.ParseUnsafe("(A <=> B & (~A | C | X)) => a + b + c <= 2 & (~a | d => X & ~B)")

	style1 := graphical.Rectangle(graphical.ColorDarkGray, graphical.ColorDarkGray, graphical.ColorLightGray)
	style2 := graphical.Circle(graphical.ColorYellow, graphical.ColorBlack, graphical.ColorYellow)
	style3 := graphical.Circle(graphical.ColorTurquoise, graphical.ColorWhite, graphical.ColorTurquoise)
	style4 := graphical.Ellipse(graphical.ColorBlack, graphical.ColorBlack, graphical.ColorWhite)

	computeNodeStyle := func(content f.Formula) *graphical.NodeStyle {
		switch content.Sort() {
		case f.SortCC, f.SortPBC:
			return style1
		case f.SortLiteral:
			name, _, _ := fac.LiteralNamePhase(content)
			if unicode.IsLower(rune(name[0])) {
				return style2
			}
			return style3
		default:
			return style4
		}
	}

	computeLabel := func(content f.Formula) string { return fmt.Sprintf("Formula Type: %s", content.Sort()) }

	generator := &graphical.Generator[f.Formula]{
		BackgroundColor:  graphical.Color("#444444"),
		DefaultEdgeStyle: graphical.Solid(graphical.ColorPurple),
		ComputeNodeStyle: computeNodeStyle,
		ComputeLabel:     computeLabel,
	}

	testFilesDag(t, "f9", f.GenerateGraphicalFormulaDAG(fac, f9, generator))
}

func TestDagWriterEdgeMapper(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	f10 := p.ParseUnsafe("(A <=> B & (~A | C | X)) => a + b + c <= 2 & (~a | d => X & ~B)")

	style1 := graphical.Dotted(graphical.ColorDarkGray)

	computeEdgeStyle := func(src, _ f.Formula) *graphical.EdgeStyle {
		if src.Sort() == f.SortPBC || src.Sort() == f.SortCC {
			return style1
		}
		return graphical.NoEdgeStyle()
	}
	generator := &graphical.Generator[f.Formula]{
		DefaultEdgeStyle: graphical.Solid(graphical.ColorPurple),
		ComputeEdgeStyle: computeEdgeStyle,
	}

	testFilesDag(t, "f10", f.GenerateGraphicalFormulaDAG(fac, f10, generator))
}

func testFilesDag(t *testing.T, fileName string, representation *graphical.Representation) {
	testFiles(t, fileName, representation, "dag")
}

func testFilesAst(t *testing.T, fileName string, representation *graphical.Representation) {
	testFiles(t, fileName, representation, "ast")
}

func testFiles(t *testing.T, fileName string, representation *graphical.Representation, sort string) {
	mermaid := graphical.WriteMermaidToString(representation)
	dot := graphical.WriteDotToString(representation)

	expectedDot, _ := os.ReadFile("./data/graphical/formulas-" + sort + "/" + fileName + ".dot")
	expected := string(expectedDot)
	assert.Equal(t, expected, dot)

	expectedMermaid, _ := os.ReadFile("./data/graphical/formulas-" + sort + "/" + fileName + ".txt")
	expected = string(expectedMermaid)
	assert.Equal(t, expected, mermaid)
}
