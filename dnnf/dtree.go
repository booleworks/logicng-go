package dnnf

import (
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/function"
	"booleworks.com/logicng/sat"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/utils"
)

type dtree interface {
	initialize(solver sat.DnnfSatSolver)
	size() int32
	staticVarSetArray() []int32
	staticVarSet() *bitset
	staticVariableSet(fac f.Factory) *f.VarSet
	dynamicSeparator() *bitset
	staticClauseIds() []int32
	countUnsubsumedOccurrences(occurrences []int32)
	depth() int
	widestSeparator() int
	leaves() []*dtreeLeaf
}

type dtreeLeaf struct {
	statVariables []int32
	statVarSet    *bitset
	statSeparator []int32

	id              int32
	clause          f.Formula
	clauseSize      int32
	literals        []int32
	separatorBitSet *bitset
	solver          sat.DnnfSatSolver
	statClauseIds   []int32
}

type dtreeNode struct {
	statVariables []int32
	statVarSet    *bitset
	statSeparator []int32

	left                dtree
	right               dtree
	sz                  int32
	solver              sat.DnnfSatSolver
	statVariableSet     *f.VarSet
	statSeparatorBitSet *bitset
	statClauseIds       []int32
	dpth                int
	wSeparator          int
	lvs                 []*dtreeLeaf
	leftLeaves          []*dtreeLeaf
	rightLeaves         []*dtreeLeaf
	clauseContents      []int32
	leftClauseContents  []int32
	rightClauseContents []int32
	localLeftVarSet     *bitset
	localRightVarSet    *bitset
}

func newDtreeLeaf(fac f.Factory, id int32, clause f.Formula) *dtreeLeaf {
	return &dtreeLeaf{
		id:              id,
		clause:          clause,
		statClauseIds:   []int32{id},
		clauseSize:      int32(function.NumberOfAtoms(fac, clause)),
		statSeparator:   []int32{},
		separatorBitSet: newBitset(),
	}
}

func (l *dtreeLeaf) initialize(solver sat.DnnfSatSolver) {
	l.solver = solver
	lits := f.Literals(solver.Factory(), l.clause)
	size := lits.Size()
	l.statVarSet = newBitset()
	l.statVariables = make([]int32, size)
	l.literals = make([]int32, size)
	for i, literal := range lits.Content() {
		variable := solver.VariableIndex(literal)
		l.statVarSet.set(variable)
		l.statVariables[i] = variable
		l.literals[i] = sat.MkLit(variable, literal.IsNeg())
	}
}

func (l *dtreeLeaf) size() int32                { return 1 }
func (l *dtreeLeaf) staticVarSetArray() []int32 { return l.statVariables }
func (l *dtreeLeaf) staticVarSet() *bitset      { return l.statVarSet }
func (l *dtreeLeaf) dynamicSeparator() *bitset  { return l.separatorBitSet }
func (l *dtreeLeaf) staticClauseIds() []int32   { return l.statClauseIds }
func (l *dtreeLeaf) depth() int                 { return 1 }
func (l *dtreeLeaf) widestSeparator() int       { return 0 }
func (l *dtreeLeaf) leaves() []*dtreeLeaf       { return []*dtreeLeaf{l} }

func (l *dtreeLeaf) staticVariableSet(fac f.Factory) *f.VarSet {
	return f.Variables(fac, l.clause)
}

func (l *dtreeLeaf) countUnsubsumedOccurrences(occurrences []int32) {
	if !l.isSubsumed() {
		for _, variable := range l.statVariables {
			occ := occurrences[variable]
			if occ != -1 {
				occurrences[variable]++
			}
		}
	}
}

func (l *dtreeLeaf) isSubsumed() bool {
	for _, lit := range l.literals {
		if l.solver.ValueOf(lit) == f.TristateTrue {
			return true
		}
	}
	return false
}

func newDtreeNode(fac f.Factory, left, right dtree) *dtreeNode {
	node := dtreeNode{}
	node.left = left
	node.right = right
	node.sz = left.size() + right.size()

	ll := left.leaves()
	node.excludeUnitLeaves(&ll)
	node.leftLeaves = ll
	rl := right.leaves()
	node.excludeUnitLeaves(&rl)
	node.rightLeaves = rl
	node.lvs = make([]*dtreeLeaf, 0, len(node.leftLeaves)+len(node.rightLeaves))
	node.lvs = append(node.lvs, node.leftLeaves...)
	node.lvs = append(node.lvs, node.rightLeaves...)

	node.statVariableSet = f.NewVarSet()
	node.statVariableSet.AddAll(left.staticVariableSet(fac))
	node.statVariableSet.AddAll(right.staticVariableSet(fac))
	node.statSeparatorBitSet = newBitset()
	leftClauseIds := left.staticClauseIds()
	rightClauseIds := right.staticClauseIds()
	node.statClauseIds = make([]int32, 0, len(leftClauseIds)+len(rightClauseIds))
	node.statClauseIds = append(node.statClauseIds, leftClauseIds...)
	node.statClauseIds = append(node.statClauseIds, rightClauseIds...)
	node.dpth = 1 + max(left.depth(), right.depth())
	return &node
}

func (n *dtreeNode) initialize(solver sat.DnnfSatSolver) {
	n.solver = solver
	n.left.initialize(solver)
	n.right.initialize(solver)
	n.statVarSet = n.left.staticVarSet()
	n.statVarSet.or(n.right.staticVarSet())
	n.statVariables = toArray(n.statVarSet)
	n.statSeparator = sortedIntersect(n.left.staticVarSetArray(), n.right.staticVarSetArray())
	for _, i := range n.statSeparator {
		n.statSeparatorBitSet.set(i)
	}
	n.wSeparator = max(len(n.statSeparator), n.left.widestSeparator(), n.right.widestSeparator())
	n.localLeftVarSet = newBitset(n.statVariables[len(n.statVariables)-1])
	n.localRightVarSet = newBitset(n.statVariables[len(n.statVariables)-1])

	var lClauseContents []int32
	for _, leaf := range n.leftLeaves {
		lClauseContents = append(lClauseContents, leaf.literals...)
		lClauseContents = append(lClauseContents, -leaf.id-1)
	}
	n.leftClauseContents = lClauseContents
	var rClauseContents []int32
	for _, leaf := range n.rightLeaves {
		rClauseContents = append(rClauseContents, leaf.literals...)
		rClauseContents = append(rClauseContents, -leaf.id-1)
	}
	n.rightClauseContents = rClauseContents
	n.clauseContents = make([]int32, 0, len(n.leftClauseContents)+len(n.rightClauseContents))
	n.clauseContents = append(n.clauseContents, n.leftClauseContents...)
	n.clauseContents = append(n.clauseContents, n.rightClauseContents...)
}

func (n *dtreeNode) size() int32                { return n.sz }
func (n *dtreeNode) staticVarSetArray() []int32 { return n.statVariables }
func (n *dtreeNode) staticVarSet() *bitset      { return n.statVarSet }
func (n *dtreeNode) staticClauseIds() []int32   { return n.statClauseIds }
func (n *dtreeNode) depth() int                 { return n.dpth }
func (n *dtreeNode) widestSeparator() int       { return n.wSeparator }

func (n *dtreeNode) staticVariableSet(_ f.Factory) *f.VarSet {
	return n.statVariableSet
}

func (n *dtreeNode) dynamicSeparator() *bitset {
	n.localLeftVarSet.clear()
	n.localRightVarSet.clear()
	n.varSet(n.leftClauseContents, n.localLeftVarSet)
	n.varSet(n.rightClauseContents, n.localRightVarSet)
	n.localLeftVarSet.and(n.localRightVarSet)
	return n.localLeftVarSet
}

func (n *dtreeNode) countUnsubsumedOccurrences(occurrences []int32) {
	for _, leaf := range n.lvs {
		leaf.countUnsubsumedOccurrences(occurrences)
	}
}

func (n *dtreeNode) leaves() []*dtreeLeaf {
	result := n.left.leaves()
	result = append(result, n.right.leaves()...)
	return result
}

func (n *dtreeNode) cacheKey(key *bitset, numberOfVariables int32) {
	i := 0
	for i < len(n.clauseContents) {
		j := i
		subsumed := false
		for n.clauseContents[j] >= 0 {
			if !subsumed && n.solver.ValueOf(n.clauseContents[j]) == f.TristateTrue {
				subsumed = true
			}
			j++
		}
		if !subsumed {
			key.set(-n.clauseContents[j] + 1 + numberOfVariables)
			for m := i; m < j; m++ {
				if n.solver.ValueOf(n.clauseContents[m]) == f.TristateUndef {
					key.set(sat.Vari(n.clauseContents[m]))
				}
			}
		}
		i = j + 1
	}
}

func (n *dtreeNode) excludeUnitLeaves(leaves *[]*dtreeLeaf) {
	nonUnit := make([]*dtreeLeaf, 0, len(*leaves))
	for _, leaf := range *leaves {
		if leaf.clauseSize > 1 {
			nonUnit = append(nonUnit, leaf)
		}
	}
	*leaves = nonUnit
}

func toArray(bits *bitset) []int32 {
	result := make([]int32, bits.cardinality())
	n := 0
	for i := bits.nextSetBit(0); i != -1; i = bits.nextSetBit(i + 1) {
		result[n] = i
		n++
	}
	return result
}

func sortedIntersect(left, right []int32) []int32 {
	l := treeset.NewWith(utils.Int32Comparator)
	intersection := treeset.NewWith(utils.Int32Comparator)
	for _, i := range left {
		l.Add(i)
	}
	for _, i := range right {
		if l.Contains(i) {
			intersection.Add(i)
		}
	}
	result := make([]int32, intersection.Size())
	intersection.Each(func(i int, value interface{}) {
		result[i] = value.(int32)
	})
	return result
}

func (n *dtreeNode) varSet(clausesContents []int32, localVarSet *bitset) {
	i := 0
	for i < len(clausesContents) {
		j := i
		subsumed := false
		for clausesContents[j] >= 0 {
			if !subsumed && n.solver.ValueOf(clausesContents[j]) == f.TristateTrue {
				subsumed = true
			}
			j++
		}
		if !subsumed {
			for m := i; m < j; m++ {
				if n.solver.ValueOf(clausesContents[m]) == f.TristateUndef {
					localVarSet.set(sat.Vari(clausesContents[m]))
				}
			}
		}
		i = j + 1
	}
}

// Bitset stuff

type bitset struct {
	bits []bool
}

func newBitset(size ...int32) *bitset {
	capacity := int32(0)
	if len(size) > 0 {
		capacity = size[0]
	}
	return &bitset{make([]bool, capacity)}
}

func (b *bitset) set(index int32) {
	if int(index) < len(b.bits) {
		b.bits[index] = true
	} else {
		b.ensureSize(int(index) + 1)
		b.bits[index] = true
	}
}

func (b *bitset) get(index int) bool {
	return index < len(b.bits) && b.bits[index]
}

func (b *bitset) or(other *bitset) {
	b.ensureSize(len(other.bits))
	for i := 0; i < len(b.bits); i++ {
		b.bits[i] = b.bits[i] || i < len(other.bits) && other.bits[i]
	}
}

func (b *bitset) and(other *bitset) {
	b.ensureSize(len(other.bits))
	for i := 0; i < len(b.bits); i++ {
		b.bits[i] = b.bits[i] && i < len(other.bits) && other.bits[i]
	}
}

func (b *bitset) ensureSize(size int) {
	if len(b.bits) >= size {
		return
	}
	newBits := make([]bool, size)
	copy(newBits, b.bits)
	b.bits = newBits
}

func (b *bitset) cardinality() int {
	count := 0
	for _, b := range b.bits {
		if b {
			count++
		}
	}
	return count
}

func (b *bitset) nextSetBit(fromIndex int32) int32 {
	for i := fromIndex; i < int32(len(b.bits)); i++ {
		if b.bits[i] {
			return i
		}
	}
	return -1
}

func (b *bitset) clear() {
	for i := 0; i < len(b.bits); i++ {
		b.bits[i] = false
	}
}

func (b *bitset) clone() *bitset {
	cln := make([]bool, len(b.bits))
	copy(cln, b.bits)
	return &bitset{cln}
}
