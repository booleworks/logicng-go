package simplification

import (
	"github.com/emirpasic/gods/lists/arraylist"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/emirpasic/gods/sets/treeset"
)

// CNFSubsumption performs subsumption on a given CNF formula and returns a new
// CNF. I.e. it performs as many subsumptions as possible. A subsumption in a
// CNF means, that e.g. a clause A | B | C is subsumed by another clause A | B
// and can therefore be deleted for an equivalent CNF.  Returns with an error
// if the input formula was not in CNF.
func CNFSubsumption(fac f.Factory, formula f.Formula) (f.Formula, error) {
	if !normalform.IsCNF(fac, formula) {
		return 0, errorx.BadInput("Formula not in CNF")
	}
	if formula.Sort() <= f.SortLiteral || formula.Sort() == f.SortOr {
		return formula, nil
	}
	ubTree := generateSubsumedUBTree(fac, formula)
	sets := ubTree.allSets()
	return combine(sets, fac.Or, fac.And), nil
}

// DNFSubsumption performs subsumption on a given DNF formula and returns a new
// DNF. I.e. it performs as many subsumptions as possible. A subsumption in a
// DNF means, that e.g. a minterm A & B is subsumed by another clause A & B & C
// and can therefore be deleted for an equivalent DNF.  Returns with an error
// if the input formula was not in DNF.
func DNFSubsumption(fac f.Factory, formula f.Formula) (f.Formula, error) {
	if !normalform.IsDNF(fac, formula) {
		return 0, errorx.BadInput("Formula not in DNF")
	}
	if formula.Sort() <= f.SortLiteral || formula.Sort() == f.SortAnd {
		return formula, nil
	}
	ubTree := generateSubsumedUBTree(fac, formula)
	sets := ubTree.allSets()
	return combine(sets, fac.And, fac.Or), nil
}

func combine(sets *linkedhashset.Set, innerFunc, outerFunc func(...f.Formula) f.Formula) f.Formula {
	clauses := make([]f.Formula, sets.Size())
	sets.Each(func(i int, _lits interface{}) {
		lits := _lits.(*treeset.Set)
		literals := make([]f.Literal, lits.Size())
		lits.Each(func(i int, val interface{}) { literals[i] = val.(f.Literal) })
		clauses[i] = innerFunc(f.LiteralsAsFormulas(literals)...)
	})
	return outerFunc(clauses...)
}

func generateSubsumedUBTree(fac f.Factory, formula f.Formula) *ubtree {
	mapping := treemap.NewWithIntComparator()
	for _, term := range fac.Operands(formula) {
		lits := f.Literals(fac, term)
		terms, ok := mapping.Get(lits.Size())
		if !ok {
			terms = arraylist.New()
			mapping.Put(lits.Size(), terms)
		}
		terms.(*arraylist.List).Add(lits)
	}
	ubTree := newUbtree()
	mapping.Each(func(_ interface{}, value interface{}) {
		value.(*arraylist.List).Each(func(_ int, _set interface{}) {
			set := _set.(*f.LitSet)
			if ubTree.firstSubset(set) == nil {
				ubTree.addSet(set)
			}
		})
	})
	return ubTree
}

type ubnode struct {
	element  f.Literal
	children *treemap.Map
	endSet   *treeset.Set
}

func newUbnode(element f.Literal) *ubnode {
	return &ubnode{
		element:  element,
		children: treemap.NewWith(f.Comparator),
	}
}

func (n *ubnode) isEndOfPath() bool {
	return n.endSet != nil
}

type ubtree struct {
	rootNodes *treemap.Map
}

func newUbtree() *ubtree {
	return &ubtree{treemap.NewWith(f.Comparator)}
}

func (u *ubtree) addSet(formulas *f.LitSet) {
	nodes := u.rootNodes
	var node *ubnode
	set := convertSet(formulas)
	set.Each(func(_ int, element interface{}) {
		res, ok := nodes.Get(element)
		if !ok {
			node = newUbnode(element.(f.Literal))
			nodes.Put(element, node)
		} else {
			node = res.(*ubnode)
		}
		nodes = node.children
	})
	if node != nil {
		node.endSet = set
	}
}

func (u *ubtree) firstSubset(formulas *f.LitSet) *treeset.Set {
	if u.rootNodes.Empty() || formulas.Empty() {
		return nil
	}
	set := convertSet(formulas)
	return u.firstSubsetRec(set, u.rootNodes)
}

func (u *ubtree) allSets() *linkedhashset.Set {
	allEndOfPathNodes := u.getAllEndOfPathNodes(u.rootNodes)
	allSets := linkedhashset.New()
	for _, node := range allEndOfPathNodes {
		allSets.Add(node.endSet)
	}
	return allSets
}

func (u *ubtree) firstSubsetRec(set *treeset.Set, forest *treemap.Map) *treeset.Set {
	nodes := u.getAllNodesContainingElements(set, forest)
	var foundSubset *treeset.Set
	nodes.Each(func(_ int, _node interface{}) {
		node := _node.(*ubnode)
		if foundSubset != nil {
			return
		}
		if node.isEndOfPath() {
			foundSubset = node.endSet
			return
		}
		remainingSet := treeset.NewWith(f.Comparator)
		set.Each(func(index int, node interface{}) {
			if index > 0 {
				remainingSet.Add(node)
			}
		})
		foundSubset = u.firstSubsetRec(remainingSet, node.children)
	})
	return foundSubset
}

func (u *ubtree) getAllNodesContainingElements(set *treeset.Set, forest *treemap.Map) *linkedhashset.Set {
	nodes := linkedhashset.New()
	set.Each(func(_ int, element interface{}) {
		node, ok := forest.Get(element)
		if ok {
			nodes.Add(node)
		}
	})
	return nodes
}

func (u *ubtree) getAllEndOfPathNodes(forest *treemap.Map) []*ubnode {
	var endOfPathNodes []*ubnode
	u.getAllEndOfPathNodesRec(forest, &endOfPathNodes)
	return endOfPathNodes
}

func (u *ubtree) getAllEndOfPathNodesRec(forest *treemap.Map, endOfPathNodes *[]*ubnode) {
	for _, _node := range forest.Values() {
		node := _node.(*ubnode)
		if node.isEndOfPath() {
			*endOfPathNodes = append(*endOfPathNodes, node)
		}
		u.getAllEndOfPathNodesRec(node.children, endOfPathNodes)
	}
}

func convertSet(formulas *f.LitSet) *treeset.Set {
	set := treeset.NewWith(f.Comparator)
	for _, formula := range formulas.Content() {
		set.Add(formula)
	}
	return set
}
