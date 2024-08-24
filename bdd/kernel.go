package bdd

import (
	"math"

	"github.com/booleworks/logicng-go/errorx"
	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
)

const (
	bddTrue  = int32(1)
	bddFalse = int32(0)
)

const (
	maxvar   = int32(0x1FFFFF)
	maxref   = int32(0x3FF)
	markon   = int32(0x200000)
	markoff  = int32(0x1FFFFF)
	markhide = int32(0x1FFFFF)
)

const (
	cacheidRestrict    = 0x1
	cacheidSatcou      = 0x2
	cacheidPathcouOne  = 0x4
	cacheidPathcouZero = 0x8
	cacheidForall      = 0x1
)

type operand struct {
	v  int32
	tt [4]int
}

var (
	bddAnd   = operand{0, [4]int{0, 0, 0, 1}}
	bddOr    = operand{2, [4]int{0, 1, 1, 1}}
	bddImp   = operand{5, [4]int{1, 1, 0, 1}}
	bddEquiv = operand{6, [4]int{1, 0, 0, 1}}
	bddNot   = operand{10, [4]int{1, 1, 0, 0}}
)

var succ = handler.Success()

// A Kernel holds all internal data structures used during the compilation,
// especially the node cache. Algorithms on BDDs always require the kernel
// which was used to compile the BDD.  All fields of the kernel are private
// since there should be no reason ever to manually access or change them.
type Kernel struct {
	fac     f.Factory
	var2idx map[f.Variable]int32
	idx2var map[int32]f.Variable

	reordering *reordering

	nodes           []node  // All the bdd nodes
	vars            []int32 // Set of defined BDD variables
	minfreenodes    int32   // Minimal % of nodes that has to be left after a garbage collection
	gbcollectnum    int32   // Number of garbage collections
	cachesize       int32   // Size of the operator caches
	nodesize        int32   // Number of allocated nodes
	maxnodeincrease int32   // Max. # of nodes used to inc. table
	freepos         int32   // First free node
	freenum         int32   // Number of free nodes
	produced        int     // Number of newBdd nodes ever produced
	varnum          int32   // Number of defined BDD variables
	refstack        []int32 // Internal node reference stack
	refstacktop     int32   // Internal node reference stack top
	level2var       []int32 // Level -> variable table
	var2level       []int32 // Variable -> level table

	quantvarset   []int32 // Current variable set for quant.
	quantvarsetId int32   // Current id used in quantvarset
	quantlast     int32   // Current last variable to be quant.

	applycache   *cache // Cache for apply results
	itecache     *cache // Cache for ITE results
	quantcache   *cache // Cache for exist/forall results
	appexcache   *cache // Cache for appex/appall results
	replacecache *cache // Cache for replace results
	misccache    *cache // Cache for other results
}

// NewKernel constructs a newBdd BDD kernel with numVars the number of variables
// on the kernel, nodeSize the initial number of nodes in the internal node
// table, and cacheSize the fixed size of the internal node cache.
func NewKernel(fac f.Factory, numVars, nodeSize, cacheSize int32) *Kernel {
	kernel := &Kernel{}
	kernel.fac = fac
	kernel.var2idx = make(map[f.Variable]int32, 16)
	kernel.idx2var = make(map[int32]f.Variable, 16)
	kernel.nodesize = int32(primeGte(max(int(nodeSize), 3)))
	kernel.nodes = make([]node, kernel.nodesize)
	kernel.minfreenodes = 20
	for n := int32(0); n < kernel.nodesize; n++ {
		kernel.setRefcou(n, 0)
		kernel.setLow(n, -1)
		kernel.setHash(n, 0)
		kernel.setLevel(n, 0)
		kernel.setNext(n, n+1)
	}
	kernel.setNext(kernel.nodesize-1, 0)
	kernel.setRefcou(0, maxref)
	kernel.setRefcou(1, maxref)
	kernel.setLow(0, 0)
	kernel.setHigh(0, 0)
	kernel.setLow(1, 1)
	kernel.setHigh(1, 1)
	kernel.initOperators(max(cacheSize, 3))
	kernel.freepos = 2
	kernel.freenum = kernel.nodesize - 2
	kernel.varnum = 0
	kernel.gbcollectnum = 0
	kernel.cachesize = cacheSize
	kernel.maxnodeincrease = 50000
	kernel.reordering = newReordering(kernel)
	kernel.reordering.usedNodesNextReorder = kernel.nodesize
	kernel.setNumberOfVars(numVars)
	return kernel
}

// NewKernelWithOrdering constructs a newBdd BDD kernel with ordering the list of
// ordered variables, nodeSize the initial number of nodes in the internal node
// table, and cacheSize the fixed size of the internal node cache.
func NewKernelWithOrdering(fac f.Factory, ordering []f.Variable, nodeSize, cacheSize int32) *Kernel {
	kernel := NewKernel(fac, int32(len(ordering)), nodeSize, cacheSize)
	for _, variable := range ordering {
		kernel.getOrAddVarIndex(variable)
	}
	return kernel
}

// SwapVariables swaps the two variables first and second on the kernel. Beware
// that if the kernel was used for multiple BDDs, the variables are swapped in
// all of these BDDs.  Returns an error if one of the variables cannot be found
// on the kernel.
func (k *Kernel) SwapVariables(first, second f.Variable) error {
	firstVar, err := k.IndexForVariable(first)
	if err != nil {
		return err
	}
	secondVar, err := k.IndexForVariable(second)
	if err != nil {
		return err
	}
	k.reordering.swapVariables(firstVar, secondVar)
	return nil
}

// Reorder reorders the variables on the kernel with the given reordering
// method. Beware that if the kernel was used for multiple BDDs, the reordering
// is performed on all of these BDDs.
//
// Only blocks of variables will be reordered. See the documentation of
// AddVariableBlock to learn more about such variable blocks. Without the
// definition of any block, nothing will be reordered.
//
// If the reordering should be performed without any restrictions,
// AddVariableBlockAll can be called before this method.
func (k *Kernel) Reorder(method ReorderingMethod) {
	k.reordering.reorder(method)
}

// ActivateReorderDuringBuild activates automatic reordering during the BDD
// compilation process with the given reordering method and an upper bound for
// the number of reorderings performed.
func (k *Kernel) ActivateReorderDuringBuild(method ReorderingMethod, bound int32) {
	k.reordering.setReorderDuringConstruction(method, bound)
}

// AddVariableBlock adds a variable block starting at variable first and ending
// in variable last (both inclusive).
//
// Variable blocks are used in the BDD reordering or in the automatic
// reordering during the construction of the BDD (configured by
// ActivateReorderDuringBuild). Variable blocks can be nested, i.e. one block
// can contain an arbitrary number of ("child") blocks. Furthermore, a variable
// block can also be a single variable.
//
// During reordering, the child blocks of a parent block can be reordered, but
// they are kept together. So no other block can be moved in between the child
// blocks. Furthermore, variables in a block which are not in a child block
// will be left untouched.
//
// Example: Lets assume we have a BDD with the variable ordering v1, v2, v3,
// v4, v5, v6, v7 We create the following blocks:
//
//	A  reaching from v1 to v5
//	B  reaching from v6 to v7
//	A1 reaching from v1 to v2
//	A2 reaching from v3 to v3
//	A3 reaching from v4 to v5
//
// This means that the variables of A and B can never be mixed up in the order.
// So during reordering the variables v6 and v7 can either be moved to the
// front (before A) or remain at their position. Furthermore, for example v1
// and v2 will always stay together and neither v3 nor any other variable can
// be moved in between them. On the other hand, the blocks A1, A2, and A3 can
// be swapped arbitrarily.
//
// These are valid result of a reordering based on the above blocks:
//
//	v3, v1, v2, v4, v5, v6, v7
//	v6, v7, v4, v5, v3, v1, v2
//	v6, v7, v1, v2, v3, v4, v5
//
// These however would be illegal:
//
//	v2, v1, v3, v4, v5, v6, v7 (variables in a block which are not in a child block will not be reordered)
//	v1, v3, v2, v4, v5, v6, v7 (variables of different block will not be mixed up)
//
// If a block is fixed (the example above assumed always blocks which are not
// fixed), its immediate child blocks will remain in their order. E.g. if block
// A was fixed, the blocks A1, A2, and A3 would not be allowed to be swapped.
// Let's assume block A to be fixed and that we have two other unfixed blocks:
//
//	A11 reaching from v1 to v1
//	A12 reaching from v2 to v2
//
// These are examples for legal reorderings:
//
//	v2, v1, v3, v4, v5, v6, v7 (block A is fixed, but "grandchildren" are still allowed to be reordered)
//	v6, v7, v2, v1, v3, v4, v5
//
// These are examples for illegal reorderings:
//
//	v3, v2, v1, v4, v5, v6, v7 (block A is fixed, so it's child blocks must be reordered)
//	v1, v2, v4, v5, v3, v6, v7
//
// Each block (including all nested blocks) must be defined by a separate call
// to this method. The blocks may be added in an arbitrary order, so it is not
// required to add them top-down or bottom-up. However, the blocks must not
// intersect, except of one block containing the other. Furthermore, both
// the first and the last variable must be known by the kernel
// and the level first must be lower than the level of last.
func (k *Kernel) AddVariableBlock(first, last int32, fixed bool) {
	k.reordering.addVariableBlock(first, last, fixed)
}

// AddAllVariablesAsBlock adds a single variable block for all variables known
// by the kernel.
func (k *Kernel) AddAllVariablesAsBlock() {
	k.reordering.addVariableBlockAll()
}

// IndexForVariable returns the kernel's internal index for the given variable.
// Returns an error if the given variable is not found on the kernel.
func (k *Kernel) IndexForVariable(variable f.Variable) (int32, error) {
	index, ok := k.var2idx[variable]
	if !ok {
		return -1, errorx.BadInput("variable %s unknown on the kernel", variable.Sprint(k.fac))
	} else {
		return index, nil
	}
}

func (k *Kernel) setNumberOfVars(num int32) {
	if num < 0 || num > maxvar {
		panic(errorx.IllegalState("illegal variable number: %d", num))
	}
	k.reordering.disableReorder()
	k.vars = make([]int32, num*2)
	k.level2var = make([]int32, num+1)
	k.var2level = make([]int32, num+1)
	k.refstack = make([]int32, num*2+4)
	k.refstacktop = 0
	for k.varnum < num {
		node, _ := k.makeNode(k.varnum, 0, 1)
		k.vars[k.varnum*2] = k.pushRef(node)
		node, _ = k.makeNode(k.varnum, 1, 0)
		k.vars[k.varnum*2+1] = node
		k.popref(1)
		k.setRefcou(k.vars[k.varnum*2], maxref)
		k.setRefcou(k.vars[k.varnum*2+1], maxref)
		k.level2var[k.varnum] = k.varnum
		k.var2level[k.varnum] = k.varnum
		k.varnum++
	}
	k.setLevel(0, num)
	k.setLevel(1, num)
	k.level2var[num] = num
	k.var2level[num] = num
	k.varResize()
	k.reordering.enableReorder()
}

func (k *Kernel) getOrAddVarIndex(variable f.Variable) int32 {
	index, ok := k.var2idx[variable]
	if !ok {
		if len(k.var2idx) >= int(k.varnum) {
			panic(errorx.IllegalState("no free variables left"))
		} else {
			index = int32(len(k.var2idx))
			k.var2idx[variable] = index
			k.idx2var[index] = variable
		}
	}
	return index
}

func (k *Kernel) getVariableForIndex(idx int32) (variable f.Variable, found bool) {
	variable, found = k.idx2var[idx]
	return
}

func (k *Kernel) getLevel(variable f.Variable) int32 {
	idx, ok := k.var2idx[variable]
	if ok && idx >= 0 && int(idx) < len(k.var2level) {
		return k.var2level[idx]
	} else {
		return -1
	}
}

func (k *Kernel) doWithPotentialReordering(operation func() (int32, bool)) int32 {
	k.initRef()
	res, reorder := operation()
	if !reorder {
		return res
	} else {
		k.reordering.checkReorder()
		k.initRef()
		k.reordering.disableReorder()
		res, reorder = operation()
		if reorder {
			panic(errorx.IllegalState("must never happen"))
		}
		k.reordering.enableReorder()
		return res
	}
}

func (k *Kernel) apply(l, r int32, op operand) int32 {
	return k.doWithPotentialReordering(func() (int32, bool) {
		return k.applyRec(l, r, op)
	})
}

func (k *Kernel) applyRec(l, r int32, op operand) (int32, bool) {
	var res int32
	switch op {
	case bddAnd:
		if l == r {
			return l, false
		}
		if isZero(l) || isZero(r) {
			return 0, false
		}
		if isOne(l) {
			return r, false
		}
		if isOne(r) {
			return l, false
		}
	case bddOr:
		if l == r {
			return l, false
		}
		if isOne(l) || isOne(r) {
			return 1, false
		}
		if isZero(l) {
			return r, false
		}
		if isZero(r) {
			return l, false
		}
	case bddImp:
		if isZero(l) {
			return 1, false
		}
		if isOne(l) {
			return r, false
		}
		if isOne(r) {
			return 1, false
		}
	}
	if isConst(l) && isConst(r) {
		res = int32(op.tt[l<<1|r])
	} else {
		entry := k.applycache.lookup(triple(l, r, op.v))
		if entry.a == l && entry.b == r && entry.c == op.v {
			return entry.res, false
		}
		if k.level(l) == k.level(r) {
			node, reorder := k.applyRec(k.low(l), k.low(r), op)
			if reorder {
				return -1, true
			}
			k.pushRef(node)
			node, reorder = k.applyRec(k.high(l), k.high(r), op)
			if reorder {
				return -1, true
			}
			k.pushRef(node)
			res, reorder = k.makeNode(k.level(l), k.readRef(2), k.readRef(1))
			if reorder {
				return -1, true
			}
		} else if k.level(l) < k.level(r) {
			node, reorder := k.applyRec(k.low(l), r, op)
			if reorder {
				return -1, true
			}
			k.pushRef(node)
			node, reorder = k.applyRec(k.high(l), r, op)
			if reorder {
				return -1, true
			}
			k.pushRef(node)
			res, reorder = k.makeNode(k.level(l), k.readRef(2), k.readRef(1))
			if reorder {
				return -1, true
			}
		} else {
			node, reorder := k.applyRec(l, k.low(r), op)
			if reorder {
				return -1, true
			}
			k.pushRef(node)
			node, reorder = k.applyRec(l, k.high(r), op)
			if reorder {
				return -1, true
			}
			k.pushRef(node)
			res, reorder = k.makeNode(k.level(r), k.readRef(2), k.readRef(1))
			if reorder {
				return -1, true
			}
		}
		k.popref(2)
		entry.a = l
		entry.b = r
		entry.c = op.v
		entry.res = res
	}
	return res, false
}

func (k *Kernel) addRef(root int32, hdl handler.Handler) (int32, handler.State) {
	if !hdl.ShouldResume(event.BddNewRefAdded) {
		return -1, handler.Cancellation(event.BddNewRefAdded)
	}
	if root < 2 {
		return root, succ
	}
	if root >= k.nodesize {
		panic(errorx.IllegalState("not a valid BDD root node: %d", root))
	}
	if k.low(root) == -1 {
		panic(errorx.IllegalState("not a valid BDD root node: %d", root))
	}
	k.incRef(root)
	return root, succ
}

func (k *Kernel) delRef(root int32) {
	if root < 2 {
		return
	}
	if root >= k.nodesize {
		panic(errorx.IllegalState("cannot dereference a variable > varnum"))
	}
	if k.low(root) == -1 {
		panic(errorx.IllegalState("cannot dereference variable -1"))
	}
	if !k.hasref(root) {
		panic(errorx.IllegalState("cannot dereference a variable which has no reference"))
	}
	k.decRef(root)
}

func (k *Kernel) decRef(n int32) {
	if k.refcou(n) != maxref && k.refcou(n) > 0 {
		k.setRefcou(n, k.refcou(n)-1)
	}
}

func (k *Kernel) incRef(n int32) {
	if k.refcou(n) < maxref {
		k.setRefcou(n, k.refcou(n)+1)
	}
}

func (k *Kernel) makeNodeUnsafe(level, low, high int32) int32 {
	node, _ := k.makeNode(level, low, high)
	return node
}

func (k *Kernel) makeNode(level, low, high int32) (int32, bool) {
	if low == high {
		return low, false
	}
	hash := nodehash(level, low, high, k.nodesize)
	res := k.hash(hash)
	for res != 0 {
		if k.level(res) == level && k.low(res) == low && k.high(res) == high {
			return res, false
		}
		res = k.next(res)
	}
	if k.freepos == 0 {
		k.gbc()
		if (k.nodesize-k.freenum) >= k.reordering.usedNodesNextReorder && k.reordering.reorderReady() {
			return -1, true
		}
		if (k.freenum*100)/k.nodesize <= k.minfreenodes {
			k.nodeResize(true)
			hash = nodehash(level, low, high, k.nodesize)
		}
		if k.freepos == 0 {
			panic(errorx.IllegalState("cannot allocate more space for more nodes"))
		}
	}
	res = k.freepos
	k.freepos = k.next(k.freepos)
	k.freenum--
	k.produced++
	k.setLevel(res, level)
	k.setLow(res, low)
	k.setHigh(res, high)
	k.setNext(res, k.hash(hash))
	k.setHash(hash, res)
	return res, false
}

func (k *Kernel) markCount(i int32) int {
	if i < 2 {
		return 0
	}
	if k.marked(i) || k.low(i) == -1 {
		return 0
	}
	k.setMark(i)
	count := 1
	count += k.markCount(k.low(i))
	count += k.markCount(k.high(i))
	return count
}

func (k *Kernel) gbc() {
	for r := 0; r < int(k.refstacktop); r++ {
		k.mark(k.refstack[r])
	}
	for n := int32(0); n < k.nodesize; n++ {
		if k.refcou(n) > 0 {
			k.mark(n)
		}
		k.setHash(n, 0)
	}
	k.freepos = 0
	k.freenum = 0
	for n := k.nodesize - 1; n >= 2; n-- {
		if (k.level(n)&markon) != 0 && k.low(n) != -1 {
			k.setLevel(n, k.level(n)&markoff)
			hash := nodehash(k.level(n), k.low(n), k.high(n), k.nodesize)
			k.setNext(n, k.hash(hash))
			k.setHash(hash, n)
		} else {
			k.setLow(n, -1)
			k.setNext(n, k.freepos)
			k.freepos = n
			k.freenum++
		}
	}
	k.resetCaches()
	k.gbcollectnum++
}

func (k *Kernel) gbcRehash() {
	k.freepos = 0
	k.freenum = 0
	for n := k.nodesize - 1; n >= 2; n-- {
		if k.low(n) != -1 {
			hash := nodehash(k.level(n), k.low(n), k.high(n), k.nodesize)
			k.setNext(n, k.hash(hash))
			k.setHash(hash, n)
		} else {
			k.setNext(n, k.freepos)
			k.freepos = n
			k.freenum++
		}
	}
}

func (k *Kernel) mark(i int32) {
	if i < 2 {
		return
	}
	if (k.level(i)&markon) != 0 || k.low(i) == -1 {
		return
	}
	k.setLevel(i, k.level(i)|markon)
	k.mark(k.low(i))
	k.mark(k.high(i))
}

func (k *Kernel) nodeResize(doRehash bool) {
	oldsize := k.nodesize
	k.nodesize = k.nodesize << 1
	if k.nodesize > oldsize+k.maxnodeincrease {
		k.nodesize = oldsize + k.maxnodeincrease
	}
	k.nodesize = int32(primeLte(int(k.nodesize)))
	newnodes := make([]node, k.nodesize)
	copy(newnodes, k.nodes)
	k.nodes = newnodes
	if doRehash {
		for n := int32(0); n < oldsize; n++ {
			k.setHash(n, 0)
		}
	}
	for n := oldsize; n < k.nodesize; n++ {
		k.setRefcou(n, 0)
		k.setHash(n, 0)
		k.setLevel(n, 0)
		k.setLow(n, -1)
		k.setNext(n, n+1)
	}
	k.setNext(k.nodesize-1, k.freepos)
	k.freepos = oldsize
	k.freenum += k.nodesize - oldsize
	if doRehash {
		k.gbcRehash()
	}
}

func (k *Kernel) refcou(node int32) int32 {
	return k.nodes[node].refcou
}

func (k *Kernel) level(node int32) int32 {
	return k.nodes[node].level
}

func (k *Kernel) low(node int32) int32 {
	return k.nodes[node].low
}

func (k *Kernel) high(node int32) int32 {
	return k.nodes[node].high
}

func (k *Kernel) hash(node int32) int32 {
	return k.nodes[node].hash
}

func (k *Kernel) next(node int32) int32 {
	return k.nodes[node].next
}

func (k *Kernel) setRefcou(node, refcou int32) {
	k.nodes[node].refcou = refcou
}

func (k *Kernel) setLevel(node, level int32) {
	k.nodes[node].level = level
}

func (k *Kernel) setLow(node, low int32) {
	k.nodes[node].low = low
}

func (k *Kernel) setHigh(node, high int32) {
	k.nodes[node].high = high
}

func (k *Kernel) setHash(node, hash int32) {
	k.nodes[node].hash = hash
}

func (k *Kernel) setNext(node, next int32) {
	k.nodes[node].next = next
}

func (k *Kernel) initRef() {
	k.refstacktop = 0
}

func (k *Kernel) pushRef(n int32) int32 {
	k.refstack[k.refstacktop] = n
	k.refstacktop++
	return n
}

func (k *Kernel) readRef(n int32) int32 {
	return k.refstack[k.refstacktop-n]
}

func (k *Kernel) popref(n int32) {
	k.refstacktop -= n
}

func (k *Kernel) hasref(n int32) bool {
	return k.refcou(n) > 0
}

func isConst(n int32) bool {
	return n < 2
}

func isOne(n int32) bool {
	return n == 1
}

func isZero(n int32) bool {
	return n == 0
}

func (k *Kernel) marked(n int32) bool {
	return (k.level(n) & markon) != 0
}

func (k *Kernel) setMark(n int32) {
	k.setLevel(n, k.level(n)|markon)
}

func (k *Kernel) unmarkNode(n int32) {
	k.setLevel(n, k.level(n)&markoff)
}

func nodehash(lvl, l, h, nodesize int32) int32 {
	return int32(math.Abs(float64(triple(lvl, l, h) % nodesize)))
}

func pair(a, b int32) int32 {
	return (a+b)*(a+b+1)/2 + a
}

func triple(a, b, c int32) int32 {
	return pair(c, pair(a, b))
}

func (k *Kernel) initOperators(cachesize int32) {
	k.applycache = newCache(cachesize)
	k.itecache = newCache(cachesize)
	k.quantcache = newCache(cachesize)
	k.appexcache = newCache(cachesize)
	k.replacecache = newCache(cachesize)
	k.misccache = newCache(cachesize)
	k.quantvarsetId = 0
	k.quantvarset = nil
}

func (k *Kernel) resetCaches() {
	k.applycache.reset()
	k.itecache.reset()
	k.quantcache.reset()
	k.appexcache.reset()
	k.replacecache.reset()
	k.misccache.reset()
}

func (k *Kernel) varResize() {
	k.quantvarset = make([]int32, k.varnum)
	k.quantvarsetId = 0
}

type node struct {
	refcou int32
	level  int32
	low    int32
	high   int32
	hash   int32
	next   int32
}

// Statistics holds fields with internal kernel statistics.
type Statistics struct {
	Produced int   // number of produced nodes
	Nodes    int32 // number of allocated nodes in the node table
	Free     int32 // number of free nodes in the node table
	Vars     int32 // number of variables
	Cache    int32 // cache size
	GC       int32 // number of performed garbage collections
	Used     int32 // number of used nodes
}

// Statistics returns the statistics for the kernel.
func (k *Kernel) Statistics() Statistics {
	return Statistics{
		Produced: k.produced,
		Nodes:    k.nodesize,
		Free:     k.freenum,
		Vars:     k.varnum,
		Cache:    k.cachesize,
		GC:       k.gbcollectnum,
		Used:     k.nodesize - k.freenum,
	}
}

// Factory returns the formula factory of the kernel.
func (k *Kernel) Factory() f.Factory {
	return k.fac
}
