package bdd

import (
	"math"

	"github.com/booleworks/logicng-go/errorx"
)

func (k *Kernel) ithVar(i int32) (int32, bool) {
	if i < 0 || i >= k.varnum {
		return -1, false
	}
	return k.vars[i*2], true
}

func (k *Kernel) nithVar(i int32) (int32, bool) {
	if i < 0 || i >= k.varnum {
		return -1, false
	}
	return k.vars[i*2+1], true
}

func (k *Kernel) bddVar(root int32) int32 {
	if root < 2 {
		panic(errorx.IllegalState("illegal node number: %d", root))
	}
	return k.level2var[k.level(root)]
}

func (k *Kernel) bddLow(root int32) int32 {
	if root < 2 {
		panic(errorx.IllegalState("illegal node number: %d", root))
	}
	return k.low(root)
}

func (k *Kernel) bddHigh(root int32) int32 {
	if root < 2 {
		panic(errorx.IllegalState("illegal node number: %d", root))
	}
	return k.high(root)
}

func (k *Kernel) and(l, r int32) int32 {
	return k.apply(l, r, bddAnd)
}

func (k *Kernel) or(l, r int32) int32 {
	return k.apply(l, r, bddOr)
}

func (k *Kernel) implication(l, r int32) int32 {
	return k.apply(l, r, bddImp)
}

func (k *Kernel) equivalence(l, r int32) int32 {
	return k.apply(l, r, bddEquiv)
}

func (k *Kernel) not(r int32) int32 {
	return k.doWithPotentialReordering(func() (int32, bool) {
		return k.notRec(r)
	})
}

func (k *Kernel) notRec(r int32) (int32, bool) {
	if isZero(r) {
		return bddTrue, false
	}
	if isOne(r) {
		return bddFalse, false
	}
	entry := k.applycache.lookup(r)
	if entry.a == r && entry.c == bddNot.v {
		return entry.res, false
	}
	node, _ := k.notRec(k.low(r))
	k.pushRef(node)
	node, _ = k.notRec(k.high(r))
	k.pushRef(node)
	res, reorder := k.makeNode(k.level(r), k.readRef(2), k.readRef(1))
	if reorder {
		return -1, true
	}
	k.popref(2)
	entry.a = r
	entry.c = bddNot.v
	entry.res = res
	return res, false
}

func (k *Kernel) restrict(r, variable int32) int32 {
	if variable < 2 {
		return r
	}
	k.varset2svartable(variable)
	return k.doWithPotentialReordering(func() (int32, bool) {
		return k.restrictRec(r, (variable<<3)|cacheidRestrict)
	})
}

func (k *Kernel) restrictRec(r, miscid int32) (int32, bool) {
	if isConst(r) || k.level(r) > k.quantlast {
		return r, false
	}
	entry := k.misccache.lookup(pair(r, miscid))
	if entry.a == r && entry.c == miscid {
		return entry.res, false
	}
	var res int32
	var reorder bool
	if k.insvarset(k.level(r)) {
		if k.quantvarset[k.level(r)] > 0 {
			res, reorder = k.restrictRec(k.high(r), miscid)
		} else {
			res, reorder = k.restrictRec(k.low(r), miscid)
		}
	} else {
		node, reorder := k.restrictRec(k.low(r), miscid)
		if reorder {
			return -1, true
		}
		k.pushRef(node)
		node, reorder = k.restrictRec(k.high(r), miscid)
		if reorder {
			return -1, true
		}
		k.pushRef(node)
		res, reorder = k.makeNode(k.level(r), k.readRef(2), k.readRef(1))
		if reorder {
			return -1, true
		}
		k.popref(2)
	}
	if reorder {
		return -1, true
	}
	entry.a = r
	entry.c = miscid
	entry.res = res
	return res, false
}

func (k *Kernel) exists(r, variable int32) int32 {
	if variable < 2 {
		return r
	}
	k.varset2vartable(variable)
	return k.doWithPotentialReordering(func() (int32, bool) {
		return k.quantRec(r, bddOr, variable<<3)
	})
}

func (k *Kernel) forAll(r, variable int32) int32 {
	if variable < 2 {
		return r
	}
	k.varset2vartable(variable)
	return k.doWithPotentialReordering(func() (int32, bool) {
		return k.quantRec(r, bddAnd, (variable<<3)|cacheidForall)
	})
}

func (k *Kernel) quantRec(r int32, op operand, quantid int32) (int32, bool) {
	if r < 2 || k.level(r) > k.quantlast {
		return r, false
	}
	entry := k.quantcache.lookup(r)
	if entry.a == r && entry.c == quantid {
		return entry.res, false
	}
	var res int32
	var reorder bool
	node, reorder := k.quantRec(k.low(r), op, quantid)
	if reorder {
		return -1, true
	}
	k.pushRef(node)
	node, reorder = k.quantRec(k.high(r), op, quantid)
	if reorder {
		return -1, true
	}
	k.pushRef(node)
	if k.invarset(k.level(r)) {
		res, reorder = k.applyRec(k.readRef(2), k.readRef(1), op)
	} else {
		res, reorder = k.makeNode(k.level(r), k.readRef(2), k.readRef(1))
	}
	if reorder {
		return -1, true
	}
	k.popref(2)
	entry.a = r
	entry.c = quantid
	entry.res = res
	return res, false
}

func (k *Kernel) varset2svartable(root int32) {
	if root < 2 {
		panic(errorx.IllegalState("illegal variable: %d", root))
	}
	k.quantvarsetId++
	if k.quantvarsetId == math.MaxInt32 {
		k.quantvarset = make([]int32, k.varnum)
		k.quantvarsetId = 1
	}
	for n := root; !isConst(n); {
		if isZero(k.low(n)) {
			k.quantvarset[k.level(n)] = k.quantvarsetId
			n = k.high(n)
		} else {
			k.quantvarset[k.level(n)] = -k.quantvarsetId
			n = k.low(n)
		}
		k.quantlast = k.level(n)
	}
}

func (k *Kernel) varset2vartable(root int32) {
	if root < 2 {
		panic(errorx.IllegalState("illegal variable: %d", root))
	}
	k.quantvarsetId++
	if k.quantvarsetId == math.MaxInt32 {
		k.quantvarset = make([]int32, k.varnum)
		k.quantvarsetId = 1
	}
	for n := root; n > 1; n = k.high(n) {
		k.quantvarset[k.level(n)] = k.quantvarsetId
		k.quantlast = k.level(n)
	}
}

func (k *Kernel) insvarset(a int32) bool {
	return int32(math.Abs(float64(k.quantvarset[a]))) == k.quantvarsetId
}

func (k *Kernel) invarset(a int32) bool {
	return k.quantvarset[a] == k.quantvarsetId
}
