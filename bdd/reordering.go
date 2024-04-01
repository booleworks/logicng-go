package bdd

import (
	"math"
	"math/rand"
	"sort"

	"github.com/booleworks/logicng-go/errorx"
)

type ReorderingMethod byte

const (
	ReorderNone    ReorderingMethod = iota // no reordering
	ReorderWin2                            // sliding window of size 2
	ReorderWin2Ite                         // sliding window of size 2 iterative
	ReorderSift                            // sifting
	ReorderSiftIte                         // iterative sifting
	ReorderWin3                            // sliding window of size 3
	ReorderWin3Ite                         // sliding window of size 3 iterative
	ReorderRandom                          // random reordering (should only be used for testing)
)

//go:generate stringer -type=ReorderingMethod

type reordering struct {
	k *Kernel

	reorderMethod   ReorderingMethod
	bddreorderTimes int32
	reorderDisabled bool
	varTree         *bddTree
	blockId         int32

	extRoots          []int32
	extRootSize       int32
	levels            []levelData
	interactionMatrix *interactionMatrix

	usednumBefore        int32
	usednumAfter         int32
	resizedInMakenode    bool
	usedNodesNextReorder int32
}

func newReordering(k *Kernel) *reordering {
	r := &reordering{
		k:               k,
		reorderDisabled: false,
		varTree:         nil,
		usednumBefore:   0,
		usednumAfter:    0,
		blockId:         0,
	}
	r.clrVarBlocks()
	r.setReorderDuringConstruction(ReorderNone, 0)
	return r
}

func (r *reordering) swapVariables(v1, v2 int32) {
	// Do not swap when variable-blocks are used
	if r.varTree != nil {
		panic(errorx.IllegalState("swapping variables is not allowed with variable blocks"))
	}

	// Don't bother swapping x with x
	if v1 == v2 {
		return
	}

	// Make sure the variable exists
	if v1 < 0 || v1 >= r.k.varnum {
		panic(errorx.IllegalState("unknown variable number: %d ", v1))
	}
	if v2 < 0 || v2 >= r.k.varnum {
		panic(errorx.IllegalState("unknown variable number: %d ", v2))
	}

	l1 := r.k.var2level[v1]
	l2 := r.k.var2level[v2]

	// Make sure v1 is before v2
	if l1 > l2 {
		v1, v2 = v2, v1
		l1 = r.k.var2level[v1]
		l2 = r.k.var2level[v2]
	}

	r.reorderInit()
	// Move v1 to v2's position
	for r.k.var2level[v1] < l2 {
		r.reorderVardown(v1)
	}
	// Move v2 to v1's position
	for r.k.var2level[v2] > l1 {
		r.reorderVarup(v2)
	}
	r.reorderDone()
}

func (r *reordering) reorder(method ReorderingMethod) {
	savemethod := r.reorderMethod
	savetimes := r.bddreorderTimes
	r.reorderMethod = method
	r.bddreorderTimes = 1
	top := newBddTree(-1)
	if r.reorderInit() < 0 {
		return
	}
	r.usednumBefore = r.k.nodesize - r.k.freenum
	top.first = 0
	top.last = r.k.varnum - 1
	top.fixed = false
	top.next = nil
	top.nextlevel = r.varTree

	r.reorderBlock(top, method)
	r.varTree = top.nextlevel
	r.usednumAfter = r.k.nodesize - r.k.freenum
	r.reorderDone()
	r.reorderMethod = savemethod
	r.bddreorderTimes = savetimes
}

func (r *reordering) setReorderDuringConstruction(method ReorderingMethod, num int32) {
	r.reorderMethod = method
	r.bddreorderTimes = num
}

func (r *reordering) addVariableBlock(first, last int32, fixed bool) {
	if first < 0 || first >= r.k.varnum || last < 0 || last >= r.k.varnum {
		panic(errorx.IllegalState("invalid var range from %d to %d", first, last))
	}
	t := addRange(r.varTree, first, last, fixed, r.blockId, r.k.level2var)
	if t == nil {
		panic(errorx.IllegalState("could not add range to tree"))
	}
	r.varTree = t
	r.blockId++
}

func (r *reordering) addVariableBlockAll() {
	for n := int32(0); n < r.k.varnum; n++ {
		r.addVariableBlock(n, n, false)
	}
}

func (r *reordering) vari(n int32) int32 {
	return r.k.level(n)
}

func (r *reordering) reorderNodenum() int32 {
	return r.k.nodesize - r.k.freenum
}

func (r *reordering) nodehashReorder(variable, l, h int32) int32 {
	return int32(math.Abs(float64(pair(l, h)%r.levels[variable].size)) + float64(r.levels[variable].start))
}

func (r *reordering) reorderBlock(t *bddTree, method ReorderingMethod) {
	if t == nil {
		return
	}
	if !t.fixed && t.nextlevel != nil {
		switch method {
		case ReorderWin2:
			t.nextlevel = r.reorderWin2(t.nextlevel)
		case ReorderWin2Ite:
			t.nextlevel = r.reorderWin2ite(t.nextlevel)
		case ReorderSift:
			t.nextlevel = r.reorderSift(t.nextlevel)
		case ReorderSiftIte:
			t.nextlevel = r.reorderSiftite(t.nextlevel)
		case ReorderWin3:
			t.nextlevel = r.reorderWin3(t.nextlevel)
		case ReorderWin3Ite:
			t.nextlevel = r.reorderWin3ite(t.nextlevel)
		case ReorderRandom:
			t.nextlevel = r.reorderRandom(t.nextlevel)
		}
	}
	for thisTree := t.nextlevel; thisTree != nil; thisTree = thisTree.next {
		r.reorderBlock(thisTree, method)
	}
	if t.seq != nil {
		newSeq := make([]int32, min(len(t.seq), int(t.last-t.first+1)))
		copy(newSeq, t.seq)
		sort.Slice(newSeq, func(i, j int) bool {
			a := r.k.var2level[newSeq[i]]
			b := r.k.var2level[newSeq[j]]
			return a < b
		})
		t.seq = newSeq
	}
}

func (r *reordering) reorderDone() {
	for n := int32(0); n < r.extRootSize; n++ {
		r.k.setMark(r.extRoots[n])
	}
	for n := int32(2); n < r.k.nodesize; n++ {
		if r.k.marked(n) {
			r.k.unmark(n)
		} else {
			r.k.setRefcou(n, 0)
		}
		// This is where we go from .var to .level again! - Do NOT use the LEVEL macro here.
		r.k.setLevel(n, r.k.var2level[r.k.level(n)])
	}
	r.k.gbc()
}

func (r *reordering) reorderWin2(t *bddTree) *bddTree {
	thisTree := t
	first := t
	if t == nil {
		return nil
	}
	for thisTree.next != nil {
		best := r.reorderNodenum()
		r.blockdown(thisTree)
		if best < r.reorderNodenum() {
			r.blockdown(thisTree.prev)
			thisTree = thisTree.next
		} else if first == thisTree {
			first = thisTree.prev
		}
	}
	return first
}

func (r *reordering) reorderWin2ite(t *bddTree) *bddTree {
	var thisTree *bddTree
	first := t
	if t == nil {
		return nil
	}
	var lastsize int32
	for ok := true; ok; ok = r.reorderNodenum() != lastsize {
		lastsize = r.reorderNodenum()
		thisTree = t
		for thisTree.next != nil {
			best := r.reorderNodenum()
			r.blockdown(thisTree)
			if best < r.reorderNodenum() {
				r.blockdown(thisTree.prev)
				thisTree = thisTree.next
			} else if first == thisTree {
				first = thisTree.prev
			}
		}
	}
	return first
}

func (r *reordering) reorderWin3(t *bddTree) *bddTree {
	thisTree := t
	first := t
	if t == nil {
		return nil
	}

	for thisTree.next != nil {
		_1, _2 := r.reorderSwapwin3(thisTree)
		thisTree = _1
		if _2 != nil {
			first = _2
		}
	}
	return first
}

func (r *reordering) reorderWin3ite(t *bddTree) *bddTree {
	var thisTree *bddTree
	first := t
	var lastsize int32

	if t == nil {
		return nil
	}

	for ok := true; ok; ok = r.reorderNodenum() != lastsize {
		lastsize = r.reorderNodenum()
		thisTree = first
		for thisTree.next != nil && thisTree.next.next != nil {
			_1, _2 := r.reorderSwapwin3(thisTree)
			thisTree = _1
			if _2 != nil {
				first = _2
			}
		}
	}
	return first
}

func (r *reordering) reorderSwapwin3(thisTree *bddTree) (*bddTree, *bddTree) {
	var first *bddTree
	setfirst := thisTree.prev == nil
	next := thisTree
	best := r.reorderNodenum()

	if thisTree.next.next == nil {
		// Only two blocks left -> win2 swap
		r.blockdown(thisTree)

		if best < r.reorderNodenum() {
			r.blockdown(thisTree.prev)
			next = thisTree.next
		} else {
			if setfirst {
				first = thisTree.prev
			}
		}
	} else {
		// Real win3 swap
		pos := 0
		r.blockdown(thisTree) // B A* C (4)
		pos++
		if best > r.reorderNodenum() {
			pos = 0
			best = r.reorderNodenum()
		}

		r.blockdown(thisTree) // B C A* (3)
		pos++
		if best > r.reorderNodenum() {
			pos = 0
			best = r.reorderNodenum()
		}

		thisTree = thisTree.prev.prev
		r.blockdown(thisTree) // C B* A (2)
		pos++
		if best > r.reorderNodenum() {
			pos = 0
			best = r.reorderNodenum()
		}

		r.blockdown(thisTree) // C A B* (1)
		pos++
		if best > r.reorderNodenum() {
			pos = 0
			best = r.reorderNodenum()
		}

		thisTree = thisTree.prev.prev
		r.blockdown(thisTree) // A C* B (0)
		pos++
		if best > r.reorderNodenum() {
			pos = 0
		}

		if pos >= 1 {
			// A C B -> C A* B
			thisTree = thisTree.prev
			r.blockdown(thisTree)
			next = thisTree
			if setfirst {
				first = thisTree.prev
			}
		}

		if pos >= 2 {
			// C A B -> C B A*
			r.blockdown(thisTree)
			next = thisTree.prev
			if setfirst {
				first = thisTree.prev.prev
			}
		}

		if pos >= 3 {
			// C B A -> B C* A
			thisTree = thisTree.prev.prev
			r.blockdown(thisTree)
			next = thisTree
			if setfirst {
				first = thisTree.prev
			}
		}

		if pos >= 4 {
			// B C A -> B A C*
			r.blockdown(thisTree)
			next = thisTree.prev
			if setfirst {
				first = thisTree.prev.prev
			}
		}

		if pos >= 5 {
			// B A C -> A B* C
			thisTree = thisTree.prev.prev
			r.blockdown(thisTree)
			next = thisTree
			if setfirst {
				first = thisTree.prev
			}
		}
	}
	return next, first
}

func (r *reordering) reorderSiftite(t *bddTree) *bddTree {
	first := t
	var lastsize int32
	if t == nil {
		return nil
	}
	for ok := true; ok; ok = r.reorderNodenum() != lastsize {
		lastsize = r.reorderNodenum()
		first = r.reorderSift(first)
	}
	return first
}

func (r *reordering) reorderSift(t *bddTree) *bddTree {
	var thisTree *bddTree
	var seq []*bddTree
	var p []bddSizePair
	num := int32(0)
	for thisTree = t; thisTree != nil; thisTree = thisTree.next {
		thisTree.pos = num
		num++
	}

	p = make([]bddSizePair, num)
	for i := 0; i < len(p); i++ {
		p[i] = bddSizePair{}
	}
	seq = make([]*bddTree, num)

	n := 0
	for thisTree = t; thisTree != nil; thisTree, n = thisTree.next, n+1 {
		// Accumulate number of nodes for each block
		p[n].val = 0
		for v := thisTree.first; v <= thisTree.last; v++ {
			p[n].val = p[n].val - r.levels[v].nodenum
		}
		p[n].block = thisTree
	}

	// Sort according to the number of nodes at each level
	sort.Slice(p, func(i, j int) bool { return p[i].val < p[j].val })

	// Create sequence
	for n := 0; n < int(num); n++ {
		seq[n] = p[n].block
	}

	// Do the sifting on this sequence
	t = r.reorderSiftSeq(t, seq, num)

	return t
}

func (r *reordering) reorderSiftSeq(t *bddTree, seq []*bddTree, num int32) *bddTree {
	var thisTree *bddTree
	if t == nil {
		return nil
	}
	for n := 0; n < int(num); n++ {
		r.reorderSiftBestpos(seq[n], num/2)
	}
	// Find first block
	for thisTree = t; thisTree.prev != nil; thisTree = thisTree.prev {
	}
	return thisTree
}

func (r *reordering) reorderSiftBestpos(blk *bddTree, middlePos int32) {
	best := r.reorderNodenum()
	maxAllowed := best/5 + best
	bestpos := 0
	dirIsUp := blk.pos <= middlePos

	// Move block back and forth
	for n := 0; n < 2; n++ {
		first := true

		if dirIsUp {
			for blk.prev != nil && (r.reorderNodenum() <= maxAllowed || first) {
				first = false
				r.blockdown(blk.prev)
				bestpos--
				if r.reorderNodenum() < best {
					best = r.reorderNodenum()
					bestpos = 0
					maxAllowed = best/5 + best
				}
			}
		} else {
			for blk.next != nil && (r.reorderNodenum() <= maxAllowed || first) {
				first = false
				r.blockdown(blk)
				bestpos++
				if r.reorderNodenum() < best {
					best = r.reorderNodenum()
					bestpos = 0
					maxAllowed = best/5 + best
				}
			}
		}
		dirIsUp = !dirIsUp
	}

	// Move to best pos
	for bestpos < 0 {
		r.blockdown(blk)
		bestpos++
	}
	for bestpos > 0 {
		r.blockdown(blk.prev)
		bestpos--
	}
}

// === Random reordering (mostly for debugging and test ) =============

func (r *reordering) reorderRandom(t *bddTree) *bddTree {
	var thisTree *bddTree
	var seq []*bddTree
	if t == nil {
		return nil
	}

	num := 0
	for thisTree = t; thisTree != nil; thisTree = thisTree.next {
		num++
	}
	seq = make([]*bddTree, num)
	num = 0
	for thisTree = t; thisTree != nil; thisTree = thisTree.next {
		seq[num] = thisTree
		num++
	}

	for n := 0; n < 4*num; n++ {
		blk := rand.Intn(num)
		if seq[blk].next != nil {
			r.blockdown(seq[blk])
		}
	}

	// Find first block
	for thisTree = t; thisTree.prev != nil; thisTree = thisTree.prev {
	}
	return thisTree
}

func (r *reordering) blockdown(left *bddTree) {
	right := left.next
	var n int32
	leftsize := left.last - left.first
	rightsize := right.last - right.first
	leftstart := r.k.var2level[left.seq[0]]
	lseq := left.seq
	rseq := right.seq

	// Move left past right
	for r.k.var2level[lseq[0]] < r.k.var2level[rseq[rightsize]] {
		for n = 0; n < leftsize; n++ {
			if r.k.var2level[lseq[n]]+1 != r.k.var2level[lseq[n+1]] &&
				r.k.var2level[lseq[n]] < r.k.var2level[rseq[rightsize]] {
				r.reorderVardown(lseq[n])
			}
		}
		if r.k.var2level[lseq[leftsize]] < r.k.var2level[rseq[rightsize]] {
			r.reorderVardown(lseq[leftsize])
		}
	}

	// Move right to where left started
	for r.k.var2level[rseq[0]] > leftstart {
		for n = rightsize; n > 0; n-- {
			if r.k.var2level[rseq[n]]-1 != r.k.var2level[rseq[n-1]] && r.k.var2level[rseq[n]] > leftstart {
				r.reorderVarup(rseq[n])
			}
		}
		if r.k.var2level[rseq[0]] > leftstart {
			r.reorderVarup(rseq[0])
		}
	}

	// Swap left and right data in the order
	left.next = right.next
	right.prev = left.prev
	left.prev = right
	right.next = left

	if right.prev != nil {
		right.prev.next = right
	}
	if left.next != nil {
		left.next.prev = left
	}
	n = left.pos
	left.pos = right.pos
	right.pos = n
}

func (r *reordering) reorderVarup(variable int32) {
	if variable < 0 || variable >= r.k.varnum {
		panic(errorx.IllegalState("illegal variable in reordering"))
	}
	if r.k.var2level[variable] != 0 {
		r.reorderVardown(r.k.level2var[r.k.var2level[variable]-1])
	}
}

func (r *reordering) reorderVardown(variable int32) {
	if variable < 0 || variable >= r.k.varnum {
		panic(errorx.IllegalState("illegal variable in reordering"))
	}
	level := r.k.var2level[variable]
	if level >= r.k.varnum-1 {
		return
	}
	r.resizedInMakenode = false

	if r.interactionMatrix.depends(variable, r.k.level2var[level+1]) > 0 {
		toBeProcessed := r.reorderDownSimple(variable)
		r.reorderSwap(toBeProcessed, variable)
		r.reorderLocalGbc(variable)
	}

	// Swap the var<->level tables
	n := r.k.level2var[level]
	r.k.level2var[level] = r.k.level2var[level+1]
	r.k.level2var[level+1] = n
	n = r.k.var2level[variable]
	r.k.var2level[variable] = r.k.var2level[r.k.level2var[level]]
	r.k.var2level[r.k.level2var[level]] = n

	if r.resizedInMakenode {
		r.reorderRehashAll()
	}
}

func (r *reordering) reorderDownSimple(var0 int32) int32 {
	toBeProcessed := int32(0)
	var1 := r.k.level2var[r.k.var2level[var0]+1]
	vl0 := r.levels[var0].start
	size0 := r.levels[var0].size
	r.levels[var0].nodenum = 0
	for n := int32(0); n < size0; n++ {
		q := r.k.hash(n + vl0)
		r.k.setHash(n+vl0, 0)
		for q != 0 {
			next := r.k.next(q)
			if r.vari(r.k.low(q)) != var1 && r.vari(r.k.high(q)) != var1 {
				// Node does not depend on next var, let it stay in the chain
				r.k.setNext(q, r.k.hash(n+vl0))
				r.k.setHash(n+vl0, q)
				r.levels[var0].nodenum++
			} else {
				// Node depends on next var - save it for later processing
				r.k.setNext(q, toBeProcessed)
				toBeProcessed = q
			}
			q = next
		}
	}
	return toBeProcessed
}

func (r *reordering) reorderSwap(toBeProcessed, var0 int32) {
	var1 := r.k.level2var[r.k.var2level[var0]+1]
	for toBeProcessed > 0 {
		next := r.k.next(toBeProcessed)
		f0 := r.k.low(toBeProcessed)
		f1 := r.k.high(toBeProcessed)
		var f00, f01, f10, f11, hash int32

		// Find the cofactors for the newBdd nodes
		if r.vari(f0) == var1 {
			f00 = r.k.low(f0)
			f01 = r.k.high(f0)
		} else {
			f00 = f0
			f01 = f0
		}
		if r.vari(f1) == var1 {
			f10 = r.k.low(f1)
			f11 = r.k.high(f1)
		} else {
			f10 = f1
			f11 = f1
		}

		// Note: makenode does refcou.
		f0 = r.reorderMakenode(var0, f00, f10)
		f1 = r.reorderMakenode(var0, f01, f11)

		// We know that the refcou of the grandchilds of this node is
		// greater than one (these are f00...f11), so there is no need to do
		// a recursive refcou decrease. It is also possible for the
		// LOWp(node)/high nodes to come alive again, so deref. of the
		// childs is delayed until the local GBC.
		r.k.decRef(r.k.low(toBeProcessed))
		r.k.decRef(r.k.high(toBeProcessed))

		// Update in-place
		r.k.setLevel(toBeProcessed, var1)
		r.k.setLow(toBeProcessed, f0)
		r.k.setHigh(toBeProcessed, f1)
		r.levels[var1].nodenum++
		// Rehash the node since it got newBdd childs
		hash = r.nodehashReorder(r.vari(toBeProcessed), r.k.low(toBeProcessed), r.k.high(toBeProcessed))
		r.k.setNext(toBeProcessed, r.k.hash(hash))
		r.k.setHash(hash, toBeProcessed)
		toBeProcessed = next
	}
}

func (r *reordering) reorderMakenode(variable, low, high int32) int32 {
	// Note: We know that low,high has a refcou greater than zero, so there
	// is no need to add reference *recursively*
	// Check whether childs are equal
	if low == high {
		r.k.incRef(low)
		return low
	}

	// Try to find an existing node of this kind
	hash := r.nodehashReorder(variable, low, high)
	res := r.k.hash(hash)

	for res != 0 {
		if r.k.low(res) == low && r.k.high(res) == high {
			r.k.incRef(res)
			return res
		}
		res = r.k.next(res)
	}
	// No existing node -> build one
	// Any free nodes to use ?
	if r.k.freepos == 0 {
		// Try to allocate more nodes - call noderesize without enabling
		// rehashing. Note: if ever rehashing is allowed here, then remember
		// to update local variable "hash"
		r.k.nodeResize(false)
		r.resizedInMakenode = true
	}

	// Build newBdd node
	res = r.k.freepos
	r.k.freepos = r.k.next(r.k.freepos)
	r.levels[variable].nodenum++
	r.k.produced++
	r.k.freenum--

	r.k.setLevel(res, variable)
	r.k.setLow(res, low)
	r.k.setHigh(res, high)

	// Insert node in hash chain
	r.k.setNext(res, r.k.hash(hash))
	r.k.setHash(hash, res)

	// Make sure it is reference counted
	r.k.setRefcou(res, 1)
	r.k.incRef(r.k.low(res))
	r.k.incRef(r.k.high(res))
	return res
}

func (r *reordering) reorderLocalGbc(var0 int32) {
	var1 := r.k.level2var[r.k.var2level[var0]+1]
	vl1 := r.levels[var1].start
	size1 := r.levels[var1].size
	for n := int32(0); n < size1; n++ {
		hash := n + vl1
		q := r.k.hash(hash)
		r.k.setHash(hash, 0)
		for q > 0 {
			next := r.k.next(q)

			if r.k.refcou(q) > 0 {
				r.k.setNext(q, r.k.hash(hash))
				r.k.setHash(hash, q)
			} else {
				r.k.decRef(r.k.low(q))
				r.k.decRef(r.k.high(q))
				r.k.setLow(q, -1)
				r.k.setNext(q, r.k.freepos)
				r.k.freepos = q
				r.levels[var1].nodenum--
				r.k.freenum++
			}
			q = next
		}
	}
}

func (r *reordering) reorderRehashAll() {
	r.reorderSetLevellookup()
	r.k.freepos = 0
	for n := r.k.nodesize - 1; n >= 0; n-- {
		r.k.setHash(n, 0)
	}
	for n := r.k.nodesize - 1; n >= 2; n-- {
		if r.k.refcou(n) > 0 {
			hash := r.nodehashReorder(r.vari(n), r.k.low(n), r.k.high(n))
			r.k.setNext(n, r.k.hash(hash))
			r.k.setHash(hash, n)
		} else {
			r.k.setNext(n, r.k.freepos)
			r.k.freepos = n
		}
	}
}

func (r *reordering) reorderSetLevellookup() {
	for n := int32(0); n < r.k.varnum; n++ {
		r.levels[n].maxsize = r.k.nodesize / r.k.varnum
		r.levels[n].start = n * r.levels[n].maxsize
		r.levels[n].size = r.levels[n].maxsize
		if r.levels[n].size >= 4 {
			r.levels[n].size = int32(primeLte(int(r.levels[n].size)))
		}
	}
}

func (r *reordering) clrVarBlocks() {
	r.varTree = nil
	r.blockId = 0
}

func (r *reordering) disableReorder() {
	r.reorderDisabled = true
}

func (r *reordering) enableReorder() {
	r.reorderDisabled = false
}

func (r *reordering) reorderReady() bool {
	return r.reorderMethod != ReorderNone && r.varTree != nil && r.bddreorderTimes != 0 && !r.reorderDisabled
}

func (r *reordering) reorderAuto() {
	if !r.reorderReady() {
		return
	}
	r.reorder(r.reorderMethod)
	r.bddreorderTimes--
}

func (r *reordering) reorderInit() int {
	r.levels = make([]levelData, r.k.varnum)
	for n := int32(0); n < r.k.varnum; n++ {
		r.levels[n] = levelData{}
		r.levels[n].start = -1
		r.levels[n].size = 0
		r.levels[n].nodenum = 0
	}
	// First mark and recursive refcou. all roots and childs. Also do some
	// setup here for both setLevellookup and reorder_gbc
	if r.markRoots() < 0 {
		return -1
	}
	// Initialize the hash tables
	r.reorderSetLevellookup()
	// Garbage collect and rehash to newBdd scheme
	r.reorderGbc()
	return 0
}

func (r *reordering) markRoots() int32 {
	dep := make([]int32, r.k.varnum)
	r.extRootSize = int32(0)
	for n := int32(2); n < r.k.nodesize; n++ {
		// This is where we go from .level to .var! - Do NOT use the LEVEL macro here.
		r.k.setLevel(n, r.k.level2var[r.k.level(n)])
		if r.k.refcou(n) > 0 {
			r.extRootSize++
			r.k.setMark(n)
		}
	}
	r.extRoots = make([]int32, r.extRootSize)
	r.interactionMatrix = newInteractionMatrix(r.k.varnum)
	r.extRootSize = 0
	for n := int32(2); n < r.k.nodesize; n++ {
		if r.k.marked(n) {
			r.k.unmarkNode(n)
			r.extRoots[r.extRootSize] = n
			r.extRootSize++
			dep[r.vari(n)] = 1
			r.levels[r.vari(n)].nodenum++
			r.addrefRec(r.k.low(n), dep)
			r.addrefRec(r.k.high(n), dep)
			r.addDependencies(dep)
		}
		// Make sure the hash field is empty. This saves a loop in the initial GBC
		r.k.setHash(n, 0)
	}
	r.k.setHash(0, 0)
	r.k.setHash(1, 0)
	return 0
}

func (r *reordering) reorderGbc() {
	r.k.freepos = 0
	r.k.freenum = 0
	// No need to zero all hash fields - this is done in mark_roots
	for n := r.k.nodesize - 1; n >= 2; n-- {
		if r.k.refcou(n) > 0 {
			hash := r.nodehashReorder(r.vari(n), r.k.low(n), r.k.high(n))
			r.k.setNext(n, r.k.hash(hash))
			r.k.setHash(hash, n)
		} else {
			r.k.setLow(n, -1)
			r.k.setNext(n, r.k.freepos)
			r.k.freepos = n
			r.k.freenum++
		}
	}
}

func (r *reordering) checkReorder() {
	r.reorderAuto()
	// Do not reorder before twice as many nodes have been used
	r.usedNodesNextReorder = 2 * (r.k.nodesize - r.k.freenum)
	// And if very little was gained this time (< 20%) then wait until even
	// more nodes (upto twice as many again) have been used
	if r.reorderGain() < 20 {
		r.usedNodesNextReorder += (r.usedNodesNextReorder * (20 - r.reorderGain())) / 20
	}
}

func (r *reordering) addrefRec(root int32, dep []int32) {
	if root < 2 {
		return
	}
	if r.k.refcou(root) == 0 {
		r.k.freenum--
		// Detect variable dependencies for the interaction matrix
		dep[r.vari(root)&markhide] = 1

		// Make sure the nodenum field is updated. Used in the initial GBC
		r.levels[r.vari(root)&markhide].nodenum++

		r.addrefRec(r.k.low(root), dep)
		r.addrefRec(r.k.high(root), dep)
	} else {
		// Update (from previously found) variable dependencies for the interaction matrix
		for n := int32(0); n < r.k.varnum; n++ {
			dep[n] |= r.interactionMatrix.depends(r.vari(root)&markhide, n)
		}
	}
	r.k.incRef(root)
}

func (r *reordering) addDependencies(dep []int32) {
	for n := int32(0); n < r.k.varnum; n++ {
		for m := n; m < r.k.varnum; m++ {
			if dep[n] > 0 && dep[m] > 0 {
				r.interactionMatrix.set(n, m)
				r.interactionMatrix.set(m, n)
			}
		}
	}
}

func (r *reordering) reorderGain() int32 {
	if r.usednumBefore == 0 {
		return 0
	}
	return (100 * (r.usednumBefore - r.usednumAfter)) / r.usednumBefore
}

type levelData struct {
	start   int32 // Start of this sub-table (entry in "bddnodes")
	size    int32 // Size of this sub-table
	maxsize int32 // Max. allowed size of sub-table
	nodenum int32 // Number of nodes in this level
}

type interactionMatrix struct {
	rows [][]int32
}

func newInteractionMatrix(size int32) *interactionMatrix {
	rows := make([][]int32, size)
	for n := int32(0); n < size; n++ {
		rows[n] = make([]int32, size/8+1)
	}
	return &interactionMatrix{rows}
}

func (m *interactionMatrix) set(a, b int32) {
	m.rows[a][b/8] |= 1 << (b % 8)
}

func (m *interactionMatrix) depends(a, b int32) int32 {
	return m.rows[a][b/8] & (1 << (b % 8))
}

type bddSizePair struct {
	val   int32
	block *bddTree
}
