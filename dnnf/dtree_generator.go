package dnnf

import (
	"math"

	"github.com/booleworks/logicng-go/errorx"
	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/normalform"
)

func generateMinFillDtree(fac f.Factory, cnf f.Formula, hdl handler.Handler) (dtree, handler.State) {
	graph := newGraph(fac, cnf)
	if e := event.DnnfDtreeMinFillGraphInitialized; !hdl.ShouldResume(e) {
		return nil, handler.Cancelation(e)
	}
	ordering, state := graph.getMinFillOrdering(hdl)
	if !state.Success {
		return nil, state
	}
	return generateWithEliminatingOrder(fac, cnf, ordering, hdl)
}

func generateWithEliminatingOrder(
	fac f.Factory,
	cnf f.Formula,
	ordering []f.Variable,
	hdl handler.Handler,
) (dtree, handler.State) {
	if !normalform.IsCNF(fac, cnf) || cnf.IsAtomic() {
		panic(errorx.IllegalState("cannot generate DTree from a non-cnf or atomic formula"))
	} else if cnf.Sort() != f.SortAnd {
		return newDtreeLeaf(fac, 0, cnf), succ
	}

	ops := fac.Operands(cnf)
	sigma := make([]dtree, len(ops))
	for id, clause := range ops {
		sigma[id] = newDtreeLeaf(fac, int32(id), clause)
	}

	for _, variable := range ordering {
		if e := event.DnnfDtreeProcessingNextOrderVar; !hdl.ShouldResume(e) {
			return nil, handler.Cancelation(e)
		}
		var gamma []dtree
		newSigma := make([]dtree, 0, len(sigma))
		for _, tree := range sigma {
			if tree.staticVariableSet(fac).Contains(variable) {
				gamma = append(gamma, tree)
			} else {
				newSigma = append(newSigma, tree)
			}
		}
		sigma = newSigma
		if len(gamma) > 0 {
			sigma = append(sigma, compose(fac, gamma))
		}
	}
	return compose(fac, sigma), succ
}

func compose(fac f.Factory, trees []dtree) dtree {
	if len(trees) == 1 {
		return trees[0]
	} else {
		left := compose(fac, trees[0:len(trees)/2])
		right := compose(fac, trees[len(trees)/2:])
		return newDtreeNode(fac, left, right)
	}
}

type graph struct {
	numberOfVertices int
	adjMatrix        [][]bool
	vertices         []f.Variable
	edgeList         [][]int32
}

func newGraph(fac f.Factory, cnf f.Formula) *graph {
	graph := graph{}
	vars := f.Variables(fac, cnf)
	graph.numberOfVertices = vars.Size()
	graph.vertices = make([]f.Variable, vars.Size())
	varToIndex := make(map[f.Variable]int32)
	for i, v := range vars.Content() {
		graph.vertices[i] = v
		varToIndex[v] = int32(i)
	}

	graph.adjMatrix = make([][]bool, graph.numberOfVertices)
	for i := 0; i < len(graph.adjMatrix); i++ {
		graph.adjMatrix[i] = make([]bool, graph.numberOfVertices)
	}
	edgeList := make([]map[int32]bool, graph.numberOfVertices)
	for i := 0; i < len(edgeList); i++ {
		edgeList[i] = make(map[int32]bool)
	}

	ops, _ := fac.NaryOperands(cnf)
	for _, clause := range ops {
		variables := f.Variables(fac, clause)
		varNums := make([]int32, variables.Size())
		for i, v := range variables.Content() {
			varNums[i] = varToIndex[v]
		}
		for i := 0; i < len(varNums); i++ {
			for j := i + 1; j < len(varNums); j++ {
				edgeList[varNums[i]][varNums[j]] = true
				edgeList[varNums[j]][varNums[i]] = true
				graph.adjMatrix[varNums[i]][varNums[j]] = true
				graph.adjMatrix[varNums[j]][varNums[i]] = true
			}
		}
	}
	graph.edgeList = make([][]int32, graph.numberOfVertices)
	for i := 0; i < len(edgeList); i++ {
		edges := edgeList[i]
		graph.edgeList[i] = make([]int32, len(edges))
		j := 0
		for edge := range edges {
			graph.edgeList[i][j] = edge
			j += 1
		}
	}
	return &graph
}

func (g *graph) getMinFillOrdering(hdl handler.Handler) ([]f.Variable, handler.State) {
	fillAdjMatrix := g.getCopyOfAdjMatrix()
	fillEdgeList := g.getCopyOfEdgeList()

	ordering := make([]f.Variable, g.numberOfVertices)
	processed := make([]bool, g.numberOfVertices)
	treewidth := 0

	for iteration := 0; iteration < g.numberOfVertices; iteration++ {
		if e := event.DnnfDtreeMinFillNewIteration; !hdl.ShouldResume(e) {
			return nil, handler.Cancelation(e)
		}
		var possiblyBestVertices []int32
		minEdges := math.MaxInt
		for currentVertex := 0; currentVertex < g.numberOfVertices; currentVertex++ {
			if processed[currentVertex] {
				continue
			}
			edgesAdded := 0
			neighborList := fillEdgeList[currentVertex]
			for i := 0; i < len(neighborList); i++ {
				firstNeighbor := neighborList[i]
				if processed[firstNeighbor] {
					continue
				}
				for j := i + 1; j < len(neighborList); j++ {
					secondNeighbor := neighborList[j]
					if processed[secondNeighbor] {
						continue
					}
					if !fillAdjMatrix[firstNeighbor][secondNeighbor] {
						edgesAdded++
					}
				}
			}
			if edgesAdded < minEdges {
				minEdges = edgesAdded
				possiblyBestVertices = []int32{}
				possiblyBestVertices = append(possiblyBestVertices, int32(currentVertex))
			} else if edgesAdded == minEdges {
				possiblyBestVertices = append(possiblyBestVertices, int32(currentVertex))
			}
		}

		// random choice
		bestVertex := possiblyBestVertices[0]

		neighborList := fillEdgeList[bestVertex]
		for i := 0; i < len(neighborList); i++ {
			firstNeighbor := neighborList[i]
			if processed[firstNeighbor] {
				continue
			}
			for j := i + 1; j < len(neighborList); j++ {
				secondNeighbor := neighborList[j]
				if processed[secondNeighbor] {
					continue
				}
				if !fillAdjMatrix[firstNeighbor][secondNeighbor] {
					fillAdjMatrix[firstNeighbor][secondNeighbor] = true
					fillAdjMatrix[secondNeighbor][firstNeighbor] = true
					fillEdgeList[firstNeighbor] = append(fillEdgeList[firstNeighbor], secondNeighbor)
					fillEdgeList[secondNeighbor] = append(fillEdgeList[secondNeighbor], firstNeighbor)
				}
			}
		}

		currentNumberOfEdges := 0
		for k := 0; k < g.numberOfVertices; k++ {
			if fillAdjMatrix[bestVertex][k] && !processed[k] {
				currentNumberOfEdges++
			}
		}
		if treewidth < currentNumberOfEdges {
			treewidth = currentNumberOfEdges
		}

		processed[bestVertex] = true
		ordering[iteration] = g.vertices[bestVertex]
	}
	return ordering, succ
}

func (g *graph) getCopyOfAdjMatrix() [][]bool {
	result := make([][]bool, g.numberOfVertices)
	for i := 0; i < g.numberOfVertices; i++ {
		cpy := make([]bool, len(g.adjMatrix[i]))
		copy(cpy, g.adjMatrix[i])
		result[i] = cpy
	}
	return result
}

func (g *graph) getCopyOfEdgeList() [][]int32 {
	result := make([][]int32, len(g.edgeList))
	for i := 0; i < len(g.edgeList); i++ {
		cpy := make([]int32, len(g.edgeList[i]))
		copy(cpy, g.edgeList[i])
		result[i] = cpy
	}
	return result
}
