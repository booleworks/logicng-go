package graph

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"booleworks.com/logicng/encoding"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/io"
	"booleworks.com/logicng/normalform"
	"github.com/stretchr/testify/assert"
)

func TestFormulaGraph(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	v1 := fac.Variable("v1")
	v2 := fac.Variable("v2")
	v3 := fac.Variable("v3")
	v4 := fac.Variable("v4")
	v5 := fac.Variable("v5")

	graph := NewFormulaGraph()
	n1 := graph.AddNode(v1)
	n2 := graph.AddNode(v2)
	n3 := graph.AddNode(v3)
	n4 := graph.AddNode(v4)
	n5 := graph.AddNode(v5)

	graph.Connect(v1, v2)
	graph.Connect(v1, v3)
	graph.Connect(v4, v5)
	graph.Connect(v5, v1)
	graph.Connect(v5, v1)
	graph.Connect(v5, v2)
	graph.Connect(v1, v2)
	graph.Connect(v3, v1)

	assert.Equal(3, len(graph.adjList[n1]))
	assert.Equal(2, len(graph.adjList[n2]))
	assert.Equal(1, len(graph.adjList[n3]))
	assert.Equal(1, len(graph.adjList[n4]))
	assert.Equal(3, len(graph.adjList[n5]))
}

func TestConnectedComponents30(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	g := readGraph(fac, "30")

	assert.Equal(27, len(g.nodes))

	ccs := ComputeConnectedComponents(g)
	assert.Equal(4, len(ccs))
	assert.Equal(7, len(ccs[0]))
	assert.Equal(5, len(ccs[1]))
	assert.Equal(8, len(ccs[2]))
	assert.Equal(7, len(ccs[3]))
}

func TestConnectedComponents50(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	g := readGraph(fac, "50")

	assert.Equal(50, len(g.nodes))

	ccs := ComputeConnectedComponents(g)
	assert.Equal(1, len(ccs))
	assert.Equal(50, len(ccs[0]))
}

func TestSplitFormulasWithGraph(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	config := encoding.DefaultConfig()
	config.AMOEncoder = encoding.AMOPure
	fac.PutConfiguration(config)
	parsed, _ := io.ReadFormula(fac, "../test/data/formulas/small.txt")
	formulas := make([]f.Formula, 0)
	originalFormulas := make([]f.Formula, 0)
	for _, formula := range fac.Operands(parsed) {
		originalFormulas = append(originalFormulas, formula)
		if formula.Sort() == f.SortCC || formula.Sort() == f.SortPBC {
			formulas = append(formulas, formula)
		} else {
			formulas = append(formulas, normalform.FactorizedCNF(fac, formula))
		}
	}
	constraintGraph := GenerateConstraintGraph(fac, formulas...)
	ccs := ComputeConnectedComponents(constraintGraph)
	split := SplitFormulasByComponent(fac, originalFormulas, ccs)
	assert.Equal(4, len(split))
	assert.Equal(1899, len(split[0]))
	assert.Equal(3, len(split[1]))
	assert.Equal(3, len(split[2]))
	assert.Equal(3, len(split[3]))
}

func readGraph(fac f.Factory, id string) *FormulaGraph {
	g := NewFormulaGraph()
	filename := "../test/data/graphs/graph" + id + ".txt"
	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		nodePair := strings.Split(line, ":")
		g.Connect(fac.Variable(nodePair[0]), fac.Variable(nodePair[1]))
	}
	return g
}
