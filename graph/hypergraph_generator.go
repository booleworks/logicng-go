package graph

import (
	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/normalform"
)

// HypergraphFromClauses generates a hyper-graph from a CNF given as a list of
// clauses. Each variable is represented by a node in the hyper-graph, each
// clause is represented by a hyper-edge between all variables of the clause.
// Returns an error if the input is not in CNF.
func HypergraphFromClauses(fac f.Factory, clauses ...f.Formula) (*Hypergraph, error) {
	hypergraph := NewHypergraph()
	nodes := make(map[f.Variable]*HypergraphNode)
	for _, clause := range clauses {
		switch clause.Sort() {
		case f.SortCC, f.SortPBC, f.SortEquiv, f.SortImpl, f.SortNot, f.SortAnd:
			return nil, errorx.BadInput("not a clause %s", clause.Sprint(fac))
		case f.SortLiteral, f.SortOr:
			addClause(fac, clause, hypergraph, &nodes)
		}
	}
	return hypergraph, nil
}

// HypergraphFromCNF generates a hyper-graph from a CNF formula. Each variable
// is represented by a node in the hyper-graph, each clause is represented by a
// hyper-edge between all variables of the clause. Returns an error if the
// input is not in CNF.
func HypergraphFromCNF(fac f.Factory, cnf f.Formula) (*Hypergraph, error) {
	if !normalform.IsCNF(fac, cnf) {
		return nil, errorx.BadInput("formula is not in CNF")
	}
	hypergraph := NewHypergraph()
	nodes := make(map[f.Variable]*HypergraphNode)
	switch cnf.Sort() {
	case f.SortLiteral, f.SortOr:
		addClause(fac, cnf, hypergraph, &nodes)
	case f.SortAnd:
		ops, _ := fac.NaryOperands(cnf)
		for _, clause := range ops {
			addClause(fac, clause, hypergraph, &nodes)
		}
	}
	return hypergraph, nil
}

func addClause(fac f.Factory, formula f.Formula, hypergraph *Hypergraph, nodes *map[f.Variable]*HypergraphNode) {
	variables := f.Variables(fac, formula)
	clause := make([]*HypergraphNode, variables.Size())
	for i, variable := range variables.Content() {
		node, ok := (*nodes)[variable]
		if !ok {
			node = NewHypergraphNode(hypergraph, variable)
			(*nodes)[variable] = node
		}
		clause[i] = node
	}
	hypergraph.AddEdge(clause)
}
