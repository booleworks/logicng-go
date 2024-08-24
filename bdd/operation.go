package bdd

import (
	"math/big"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
)

func (k *Kernel) satOne(r int32) int32 {
	if r < 2 {
		return r
	}
	k.reordering.disableReorder()
	k.initRef()
	res := k.satOneRec(r)
	k.reordering.enableReorder()
	return res
}

func (k *Kernel) satOneRec(r int32) int32 {
	if isConst(r) {
		return r
	}
	if isZero(k.low(r)) {
		res := k.satOneRec(k.high(r))
		return k.pushRef(k.makeNodeUnsafe(k.level(r), bddFalse, res))
	} else {
		res := k.satOneRec(k.low(r))
		return k.pushRef(k.makeNodeUnsafe(k.level(r), res, bddFalse))
	}
}

func (k *Kernel) satOneSet(r, variable, pol int32) int32 {
	if isZero(r) {
		return r
	}
	if !isConst(pol) {
		panic(errorx.IllegalState("polarity must be a constant"))
	}
	k.reordering.disableReorder()
	k.initRef()
	res := k.satOneSetRec(r, variable, pol)
	k.reordering.enableReorder()
	return res
}

func (k *Kernel) satOneSetRec(r, variable, satPolarity int32) int32 {
	if isConst(r) && isConst(variable) {
		return r
	}
	if k.level(r) < k.level(variable) {
		if isZero(k.low(r)) {
			res := k.satOneSetRec(k.high(r), variable, satPolarity)
			return k.pushRef(k.makeNodeUnsafe(k.level(r), bddFalse, res))
		} else {
			res := k.satOneSetRec(k.low(r), variable, satPolarity)
			return k.pushRef(k.makeNodeUnsafe(k.level(r), res, bddFalse))
		}
	} else if k.level(variable) < k.level(r) {
		res := k.satOneSetRec(r, k.high(variable), satPolarity)
		if satPolarity == bddTrue {
			return k.pushRef(k.makeNodeUnsafe(k.level(variable), bddFalse, res))
		} else {
			return k.pushRef(k.makeNodeUnsafe(k.level(variable), res, bddFalse))
		}
	} else {
		if isZero(k.low(r)) {
			res := k.satOneSetRec(k.high(r), k.high(variable), satPolarity)
			return k.pushRef(k.makeNodeUnsafe(k.level(r), bddFalse, res))
		} else {
			res := k.satOneSetRec(k.low(r), k.high(variable), satPolarity)
			return k.pushRef(k.makeNodeUnsafe(k.level(r), res, bddFalse))
		}
	}
}

func (k *Kernel) fullSatOne(r int32) int32 {
	if r == 0 {
		return 0
	}
	k.reordering.disableReorder()
	k.initRef()
	res := k.fullSatOneRec(r)
	for v := k.level(r) - 1; v >= 0; v-- {
		node, _ := k.makeNode(v, res, 0)
		res = k.pushRef(node)
	}
	k.reordering.enableReorder()
	return res
}

func (k *Kernel) fullSatOneRec(r int32) int32 {
	if r < 2 {
		return r
	}
	if k.low(r) != 0 {
		res := k.fullSatOneRec(k.low(r))
		for v := k.level(k.low(r)) - 1; v > k.level(r); v-- {
			res = k.pushRef(k.makeNodeUnsafe(v, res, 0))
		}
		return k.pushRef(k.makeNodeUnsafe(k.level(r), res, 0))
	} else {
		res := k.fullSatOneRec(k.high(r))
		for v := k.level(k.high(r)) - 1; v > k.level(r); v-- {
			res = k.pushRef(k.makeNodeUnsafe(v, res, 0))
		}
		return k.pushRef(k.makeNodeUnsafe(k.level(r), 0, res))
	}
}

func (k *Kernel) allSat(r int32) [][]byte {
	allsatProfile := make([]byte, k.varnum)
	for v := k.level(r) - 1; v >= 0; v-- {
		allsatProfile[k.level2var[v]] = 2
	}
	k.initRef()
	allSat := make([][]byte, 0, 16)
	k.allSatRec(r, &allSat, &allsatProfile)
	return allSat
}

func (k *Kernel) allSatRec(r int32, models *[][]byte, allsatProfile *[]byte) {
	if isOne(r) {
		cpy := make([]byte, len(*allsatProfile))
		copy(cpy, *allsatProfile)
		*models = append(*models, cpy)
		return
	}
	if isZero(r) {
		return
	}
	if !isZero(k.low(r)) {
		(*allsatProfile)[k.level2var[k.level(r)]] = 0
		for v := k.level(k.low(r)) - 1; v > k.level(r); v-- {
			(*allsatProfile)[k.level2var[v]] = 2
		}
		k.allSatRec(k.low(r), models, allsatProfile)
	}
	if !isZero(k.high(r)) {
		(*allsatProfile)[k.level2var[k.level(r)]] = 1
		for v := k.level(k.high(r)) - 1; v > k.level(r); v-- {
			(*allsatProfile)[k.level2var[v]] = 2
		}
		k.allSatRec(k.high(r), models, allsatProfile)
	}
}

func (k *Kernel) allUnsat(r int32) [][]byte {
	allunsatProfile := make([]byte, k.varnum)
	for v := k.level(r) - 1; v >= 0; v-- {
		allunsatProfile[k.level2var[v]] = 2
	}
	k.initRef()
	allUnsat := make([][]byte, 0, 16)
	k.allUnsatRec(r, &allUnsat, &allunsatProfile)
	return allUnsat
}

func (k *Kernel) allUnsatRec(r int32, models *[][]byte, allunsatProfile *[]byte) {
	if isZero(r) {
		cpy := make([]byte, len(*allunsatProfile))
		copy(cpy, *allunsatProfile)
		*models = append(*models, cpy)
		return
	}
	if isOne(r) {
		return
	}
	if !isOne(k.low(r)) {
		(*allunsatProfile)[k.level2var[k.level(r)]] = 0
		for v := k.level(k.low(r)) - 1; v > k.level(r); v-- {
			(*allunsatProfile)[k.level2var[v]] = 2
		}
		k.allUnsatRec(k.low(r), models, allunsatProfile)
	}
	if !isOne(k.high(r)) {
		(*allunsatProfile)[k.level2var[k.level(r)]] = 1
		for v := k.level(k.high(r)) - 1; v > k.level(r); v-- {
			(*allunsatProfile)[k.level2var[v]] = 2
		}
		k.allUnsatRec(k.high(r), models, allunsatProfile)
	}
}

func (k *Kernel) satCount(r int32) *big.Int {
	size, i, e := big.NewInt(0), big.NewInt(2), big.NewInt(int64(k.level(r)))
	size.Exp(i, e, nil)
	satcount := k.satCountRec(r, cacheidSatcou)
	satcount.Mul(satcount, size)
	return satcount
}

func (k *Kernel) satCountRec(root, miscid int32) *big.Int {
	if root < 2 {
		return big.NewInt(int64(root))
	}
	entry := k.misccache.lookup(root)
	if entry.a == root && entry.c == miscid {
		return entry.bdres
	}
	size := big.NewInt(0)
	s := big.NewInt(1)
	val, i, e := big.NewInt(0), big.NewInt(2), big.NewInt(int64(k.level(k.low(root))-k.level(root)-1))
	s.Mul(s, val.Exp(i, e, nil))
	size.Add(size, s.Mul(s, k.satCountRec(k.low(root), miscid)))
	s = big.NewInt(1)
	val, i, e = big.NewInt(0), big.NewInt(2), big.NewInt(int64(k.level(k.high(root))-k.level(root)-1))
	s.Mul(s, val.Exp(i, e, nil))
	size.Add(size, s.Mul(s, k.satCountRec(k.high(root), miscid)))
	entry.a = root
	entry.c = miscid
	entry.bdres = size
	return size
}

func (k *Kernel) pathCountOne(r int32) *big.Int {
	return k.pathCountRecOne(r, cacheidPathcouOne)
}

func (k *Kernel) pathCountRecOne(r, miscid int32) *big.Int {
	var size *big.Int
	if isZero(r) {
		return big.NewInt(0)
	}
	if isOne(r) {
		return big.NewInt(1)
	}
	entry := k.misccache.lookup(r)
	if entry.a == r && entry.c == miscid {
		return entry.bdres
	}
	count := k.pathCountRecOne(k.low(r), miscid)
	size = count.Add(count, k.pathCountRecOne(k.high(r), miscid))
	entry.a = r
	entry.c = miscid
	entry.bdres = size
	return size
}

func (k *Kernel) pathCountZero(r int32) *big.Int {
	return k.pathCountRecZero(r, cacheidPathcouZero)
}

func (k *Kernel) pathCountRecZero(r, miscid int32) *big.Int {
	var size *big.Int
	if isZero(r) {
		return big.NewInt(1)
	}
	if isOne(r) {
		return big.NewInt(0)
	}
	entry := k.misccache.lookup(r)
	if entry.a == r && entry.c == miscid {
		return entry.bdres
	}
	count := k.pathCountRecZero(k.low(r), miscid)
	size = count.Add(count, k.pathCountRecZero(k.high(r), miscid))
	entry.a = r
	entry.c = miscid
	entry.bdres = size
	return size
}

func (k *Kernel) support(r int32) int32 {
	supportId := int32(0)
	supportSet := make([]int32, k.varnum)
	res := int32(1)
	if r < 2 {
		return bddFalse
	}
	supportId++
	supportMin := k.level(r)
	supportMax := supportMin
	k.supportRec(r, supportId, supportSet, &supportMax)
	k.unmark(r)

	k.reordering.disableReorder()
	for n := supportMax; n >= supportMin; n-- {
		if supportSet[n] == supportId {
			k.addRef(res, handler.NopHandler)
			tmp := k.makeNodeUnsafe(n, 0, res)
			k.delRef(res)
			res = tmp
		}
	}
	k.reordering.enableReorder()
	return res
}

func (k *Kernel) supportRec(r, supportId int32, support []int32, supportMax *int32) {
	if r < 2 {
		return
	}
	if (k.level(r)&markon) != 0 || k.low(r) == -1 {
		return
	}
	support[k.level(r)] = supportId
	if k.level(r) > *supportMax {
		*supportMax = k.level(r)
	}
	k.setLevel(r, k.level(r)|markon)
	k.supportRec(k.low(r), supportId, support, supportMax)
	k.supportRec(k.high(r), supportId, support, supportMax)
}

func (k *Kernel) unmark(i int32) {
	if i < 2 {
		return
	}
	if !k.marked(i) || k.low(i) == -1 {
		return
	}
	k.unmarkNode(i)
	k.unmark(k.low(i))
	k.unmark(k.high(i))
}

func (k *Kernel) nodeCount(r int32) int {
	count := k.markCount(r)
	k.unmark(r)
	return count
}

func (k *Kernel) varProfile(r int32) []int {
	varprofile := make([]int, k.varnum)
	k.varProfileRec(r, varprofile)
	k.unmark(r)
	return varprofile
}

func (k *Kernel) varProfileRec(r int32, varprofile []int) {
	if r < 2 {
		return
	}
	if (k.level(r) & markon) != 0 {
		return
	}
	varprofile[k.level2var[k.level(r)]]++
	k.setLevel(r, k.level(r)|markon)
	k.varProfileRec(k.low(r), varprofile)
	k.varProfileRec(k.high(r), varprofile)
}

func (k *Kernel) allNodes(r int32) [][]int32 {
	result := make([][]int32, 0)
	if r < 2 {
		return result
	}
	k.mark(r)
	for n := int32(0); n < k.nodesize; n++ {
		if (k.level(n) & markon) != 0 {
			k.setLevel(n, k.level(n)&markoff)
			result = append(result, []int32{n, k.level2var[k.level(n)], k.low(n), k.high(n)})
		}
	}
	return result
}

func (k *Kernel) toFormula(fac f.Factory, r int32, followPathsToTrue bool) f.Formula {
	k.initRef()
	formula := k.toFormulaRec(fac, r, followPathsToTrue)
	if followPathsToTrue {
		return formula
	} else {
		return formula.Negate(fac)
	}
}

func (k *Kernel) toFormulaRec(fac f.Factory, r int32, followPathsToTrue bool) f.Formula {
	if isOne(r) {
		return fac.Constant(followPathsToTrue)
	}
	if isZero(r) {
		return fac.Constant(!followPathsToTrue)
	}
	variable, _ := k.getVariableForIndex(k.level2var[k.level(r)])
	low := k.low(r)
	high := k.high(r)
	var lowFormula, highFormula f.Formula
	if isRelevant(low, followPathsToTrue) {
		lowFormula = fac.And(variable.Negate(fac).AsFormula(), k.toFormulaRec(fac, low, followPathsToTrue))
	} else {
		lowFormula = fac.Falsum()
	}
	if isRelevant(high, followPathsToTrue) {
		highFormula = fac.And(variable.AsFormula(), k.toFormulaRec(fac, high, followPathsToTrue))
	} else {
		highFormula = fac.Falsum()
	}
	return fac.Or(lowFormula, highFormula)
}

func isRelevant(r int32, followPathsToTrue bool) bool {
	return followPathsToTrue && !isZero(r) || !followPathsToTrue && !isOne(r)
}
