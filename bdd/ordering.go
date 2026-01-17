package bdd

import (
	"slices"
	"sort"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/graph"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/emirpasic/gods/maps/linkedhashmap"
)

// BFSOrder generates a breadth-first-search variable ordering for the given
// formula.  It traverses the formula in a BFS manner and gathers all variables
// in the occurrence.
func BFSOrder(fac f.Factory, formula f.Formula) []f.Variable {
	variables := make([]f.Variable, 0, f.Variables(fac, formula).Size())
	queue := []f.Formula{formula}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		switch current.Sort() {
		case f.SortLiteral:
			if current.IsPos() {
				variable, _ := current.AsVariable()
				if !slices.Contains(variables, variable) {
					variables = append(variables, variable)
				}
			} else {
				variable := f.Literal(current).Variable()
				queue = append(queue, variable.AsFormula())
			}
		case f.SortNot:
			op, _ := fac.NotOperand(current)
			queue = append(queue, op)
		case f.SortImpl, f.SortEquiv:
			left, right, _ := fac.BinaryLeftRight(current)
			queue = append(queue, left, right)
		case f.SortAnd, f.SortOr:
			ops, _ := fac.NaryOperands(current)
			queue = append(queue, ops...)
		case f.SortCC, f.SortPBC:
			_, _, lits, _, _ := fac.PBCOps(current)
			for _, lit := range lits {
				variable := lit.Variable()
				if !slices.Contains(variables, variable) {
					variables = append(variables, variable)
				}
			}
		}
	}
	return variables
}

// DFSOrder generates a depth-first-search variable ordering for the given
// formula.  It traverses the formula in a DFS manner and gathers all variables
// in the occurrence.
func DFSOrder(fac f.Factory, formula f.Formula) []f.Variable {
	variables := make([]f.Variable, 0, f.Variables(fac, formula).Size())
	dfs(fac, formula, &variables)
	return variables
}

func dfs(fac f.Factory, formula f.Formula, variables *[]f.Variable) {
	switch formula.Sort() {
	case f.SortLiteral:
		variable := f.Literal(formula).Variable()
		if !slices.Contains(*variables, variable) {
			*variables = append(*variables, variable)
		}
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		dfs(fac, op, variables)
	case f.SortImpl, f.SortEquiv:
		left, right, _ := fac.BinaryLeftRight(formula)
		dfs(fac, left, variables)
		dfs(fac, right, variables)
	case f.SortAnd, f.SortOr:
		ops, _ := fac.NaryOperands(formula)
		for _, op := range ops {
			dfs(fac, op, variables)
		}
	case f.SortCC, f.SortPBC:
		_, _, lits, _, _ := fac.PBCOps(formula)
		for _, lit := range lits {
			variable := lit.Variable()
			if !slices.Contains(*variables, variable) {
				*variables = append(*variables, variable)
			}
		}
	}
}

// MaxToMinOrder generates a variable ordering sorting the variables from
// maximal to minimal occurrence in the input formula.  If two variables have
// the same number of occurrences, their ordering according to their DFS
// ordering will be considered.
func MaxToMinOrder(fac f.Factory, formula f.Formula) []f.Variable {
	return occurenceOrder(fac, formula, false)
}

// MinToMaxOrder generates a variable ordering sorting the variables from
// minimal to maximal occurrence in the input formula.  If two variables have
// the same number of occurrences, their ordering according to their DFS
// ordering will be considered.
func MinToMaxOrder(fac f.Factory, formula f.Formula) []f.Variable {
	return occurenceOrder(fac, formula, true)
}

func occurenceOrder(fac f.Factory, formula f.Formula, min2max bool) []f.Variable {
	profile := f.VariableProfile(fac, formula)
	dfs := DFSOrder(fac, formula)
	variables := make([]f.Variable, len(dfs))
	copy(variables, dfs)
	sort.Slice(variables, func(i, j int) bool {
		o1 := profile[variables[i]]
		o2 := profile[variables[j]]
		if o1 == o2 {
			return i < j
		} else if min2max {
			return o1 < o2
		}
		return o1 > o2
	})
	return variables
}

// ForceOrder generates a variable ordering for the given formula according to
// the FORCE ordering due to Aloul, Markov, and Sakallah.  This ordering only
// works for CNF formulas.  This method converts the formula to CNF before this
// ordering is called which can have side effects on the generated CNF
// variables and/or the formula caches in the factory.
func ForceOrder(fac f.Factory, formula f.Formula) []f.Variable {
	vars := f.Variables(fac, formula)
	originalVariables := f.NewMutableVarSet()
	originalVariables.AddAll(vars)
	nnf := normalform.NNF(fac, formula)
	originalVariables.AddAll(f.Variables(fac, nnf))
	cnf := normalform.PGCNFDefault(fac, nnf)
	hypergraph, _ := graph.HypergraphFromCNF(fac, cnf)
	nodes := make(map[f.Variable]*graph.HypergraphNode)
	for _, node := range hypergraph.Nodes {
		nodes[node.Content] = node
	}
	ordering := force(fac, cnf, hypergraph, &nodes)
	finalOrdering := make([]f.Variable, 0, len(ordering))
	for _, k := range ordering {
		if originalVariables.Contains(k) {
			finalOrdering = append(finalOrdering, k)
		}
	}
	for _, v := range originalVariables.Content() {
		if !slices.Contains(ordering, v) {
			finalOrdering = append(finalOrdering, v)
		}
	}
	return finalOrdering
}

func force(
	fac f.Factory,
	formula f.Formula,
	hypergraph *graph.Hypergraph,
	nodes *map[f.Variable]*graph.HypergraphNode,
) []f.Variable {
	initialOrdering := createInitialOrdering(fac, formula, nodes)
	var lastOrdering *linkedhashmap.Map
	currentOrdering := initialOrdering

	for ok := true; ok; ok = shouldProceed(lastOrdering, currentOrdering) {
		lastOrdering = currentOrdering
		newLocations := linkedhashmap.New()
		for _, node := range hypergraph.Nodes {
			newLocations.Put(node, node.ComputeTentativeNewLocation(lastOrdering))
		}
		currentOrdering = orderingFromTentativeNewLocations(newLocations)
	}

	ordering := make([]f.Variable, currentOrdering.Size())
	count := 0
	currentOrdering.Each(func(k any, _ any) {
		ordering[count] = k.(*graph.HypergraphNode).Content
		count++
	})
	return ordering
}

func createInitialOrdering(
	fac f.Factory,
	formula f.Formula,
	nodes *map[f.Variable]*graph.HypergraphNode,
) *linkedhashmap.Map {
	initialOrdering := linkedhashmap.New()
	dfsOrder := DFSOrder(fac, formula)
	for i := 0; i < len(dfsOrder); i++ {
		initialOrdering.Put((*nodes)[dfsOrder[i]], i)
	}
	return initialOrdering
}

func orderingFromTentativeNewLocations(newLocations *linkedhashmap.Map) *linkedhashmap.Map {
	list := sortedlocPairList(newLocations)
	ordering := linkedhashmap.New()
	count := 0
	for _, k := range *list {
		ordering.Put(k.node, count)
		count++
	}
	return ordering
}

func sortedlocPairList(mapping *linkedhashmap.Map) *[]locPair {
	list := make([]locPair, mapping.Size())
	count := 0
	mapping.Each(func(k any, v any) {
		list[count] = locPair{k.(*graph.HypergraphNode), v.(float64)}
		count++
	})
	sort.Slice(list, func(i, j int) bool {
		o1, _ := mapping.Get(list[i].node)
		o2, _ := mapping.Get(list[j].node)
		return o1.(float64) < o2.(float64)
	})
	return &list
}

func shouldProceed(lastOrdering, currentOrdering *linkedhashmap.Map) bool {
	if lastOrdering.Size() != currentOrdering.Size() {
		return true
	}
	return lastOrdering.Any(func(k any, v any) bool {
		val, _ := currentOrdering.Get(k)
		return val != v
	})
}

type locPair struct {
	node     *graph.HypergraphNode
	location float64
}
