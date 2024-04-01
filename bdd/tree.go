package bdd

// A BDD tree used to represent nested variables blocks or variable reorderings.
type bddTree struct {
	first int32
	last  int32   // First and last variable in this block
	pos   int32   // Sifting position
	seq   []int32 // Sequence of first...last in the current order

	// Are the sub-blocks fixed or may they be reordered
	fixed     bool
	id        int32 // A sequential id number given by addblock
	next      *bddTree
	prev      *bddTree
	nextlevel *bddTree
}

func newBddTree(id int32) *bddTree {
	return &bddTree{
		id:        id,
		first:     -1,
		last:      -1,
		fixed:     true,
		next:      nil,
		prev:      nil,
		nextlevel: nil,
		seq:       nil,
	}
}

func addRange(tree *bddTree, first, last int32, fixed bool, id int32, level2var []int32) *bddTree {
	return addRangeRec(tree, nil, first, last, fixed, id, level2var)
}

func addRangeRec(t, prev *bddTree, first, last int32, fixed bool, id int32, level2var []int32) *bddTree {
	if first < 0 || last < 0 || last < first {
		return nil
	}

	// Empty tree -> build one
	if t == nil {
		t = newBddTree(id)
		t.first = first
		t.fixed = fixed
		t.seq = make([]int32, last-first+1)
		t.last = last
		t.updateSeq(level2var)
		t.prev = prev
		return t
	}

	// Check for identity
	if first == t.first && last == t.last {
		return t
	}

	// Before this section -> insert
	if last < t.first {
		tnew := newBddTree(id)
		tnew.first = first
		tnew.last = last
		tnew.fixed = fixed
		tnew.seq = make([]int32, last-first+1)
		tnew.updateSeq(level2var)
		tnew.next = t
		tnew.prev = t.prev
		t.prev = tnew
		return tnew
	}

	// After this section -> go to next
	if first > t.last {
		next := addRangeRec(t.next, t, first, last, fixed, id, level2var)
		if next != nil {
			t.next = next
		}
		return t
	}

	// Inside this section -> insert in next level
	if first >= t.first && last <= t.last {
		nextlevel := addRangeRec(t.nextlevel, nil, first, last, fixed, id, level2var)
		if nextlevel != nil {
			t.nextlevel = nextlevel
		}
		return t
	}

	// Covering this section -> insert above this level
	if first <= t.first {
		var tnew *bddTree
		thisTree := t

		for {
			// Partial cover ->error
			if last >= thisTree.first && last < thisTree.last {
				return nil
			}
			if thisTree.next == nil || last < thisTree.next.first {
				tnew = newBddTree(id)
				tnew.first = first
				tnew.last = last
				tnew.fixed = fixed
				tnew.seq = make([]int32, last-first+1)
				tnew.updateSeq(level2var)
				tnew.nextlevel = t
				tnew.next = thisTree.next
				tnew.prev = t.prev
				if thisTree.next != nil {
					thisTree.next.prev = tnew
				}
				thisTree.next = nil
				t.prev = nil
				return tnew
			}
			thisTree = thisTree.next
		}
	}
	// partial cover
	return nil
}

func (t *bddTree) updateSeq(bddvar2level []int32) {
	var n int32
	low := t.first
	for n = t.first; n <= t.last; n++ {
		if bddvar2level[n] < bddvar2level[low] {
			low = n
		}
	}
	for n = t.first; n <= t.last; n++ {
		t.seq[bddvar2level[n]-bddvar2level[low]] = n
	}
}
