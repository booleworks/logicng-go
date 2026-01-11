package dnnf

import (
	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/sat"
	"github.com/booleworks/logicng-go/simplification"
	"github.com/emirpasic/gods/maps/treemap"
)

var succ = handler.Success()

// Compile returns a compiled DNNF for the given formula.
func Compile(fac f.Factory, formula f.Formula) *DNNF {
	dnnf, _ := CompileWithHandler(fac, formula, handler.NopHandler)
	return dnnf
}

// CompileWithHandler returns a compiled DNNF for the given formula.  The
// handler can be used to cancel the DNNF compilation.
func CompileWithHandler(fac f.Factory, formula f.Formula, hdl handler.Handler) (*DNNF, handler.State) {
	originalVariables := f.NewMutableVarSet()
	originalVariables.AddAll(f.Variables(fac, formula))
	cnf := normalform.CNF(fac, formula)
	originalVariables.AddAll(f.Variables(fac, cnf))
	simplifiedFormula := simplifyFormula(fac, cnf)
	compiler := newCompiler(fac, simplifiedFormula)
	dnnfFormula, state := compiler.compile(hdl)
	if !state.Success {
		return nil, state
	} else {
		return &DNNF{fac, dnnfFormula, originalVariables.AsImmutable()}, succ
	}
}

type compiler struct {
	fac               f.Factory
	originalCNF       f.Formula
	unitClauses       f.Formula
	nonUnitClauses    f.Formula
	solver            sat.DnnfSatSolver
	numberOfVariables int32
	cache             *treemap.Map
	localCacheKeys    [][]*bitset
	localOccurrences  [][][]int32
}

func simplifyFormula(fac f.Factory, formula f.Formula) f.Formula {
	simp := simplification.PropagateBackbone(fac, formula)
	simp, _ = simplification.CNFSubsumption(fac, simp) // formula is in cnf - no error possible
	return simp
}

func newCompiler(fac f.Factory, formula f.Formula) *compiler {
	numVars := f.Variables(fac, formula).Size()
	units, nonUnits := initializeClauses(fac, formula)
	compiler := compiler{
		fac:               fac,
		originalCNF:       formula,
		unitClauses:       units,
		nonUnitClauses:    nonUnits,
		solver:            sat.NewDnnfSolver(fac, numVars),
		numberOfVariables: int32(numVars),
		cache:             treemap.NewWith(bitsetComp),
	}
	compiler.solver.Add(formula)
	return &compiler
}

func bitsetComp(a, b any) int {
	bitset1 := a.(*bitset)
	bitset2 := b.(*bitset)
	if len(bitset1.bits) < len(bitset2.bits) {
		return -1
	}
	if len(bitset2.bits) < len(bitset1.bits) {
		return 1
	}
	// lengths are equal
	for i := 0; i < len(bitset1.bits); i++ {
		if !bitset1.bits[i] && bitset2.bits[i] {
			return -1
		}
		if !bitset2.bits[i] && bitset1.bits[i] {
			return 1
		}
	}
	return 0
}

func (c *compiler) compile(hdl handler.Handler) (f.Formula, handler.State) {
	if !sat.IsSatisfiable(c.fac, c.originalCNF) {
		return c.fac.Falsum(), succ
	}
	dTree, state := c.generateDtree(c.fac, hdl)
	if !state.Success {
		return 0, state
	}
	return c.compileWithTree(dTree, hdl)
}

func initializeClauses(fac f.Factory, cnf f.Formula) (f.Formula, f.Formula) {
	var units []f.Formula
	var nonUnits []f.Formula
	switch cnf.Sort() {
	case f.SortAnd:
		for _, clause := range fac.Operands(cnf) {
			if clause.IsAtomic() {
				units = append(units, clause)
			} else {
				nonUnits = append(nonUnits, clause)
			}
		}
	case f.SortOr:
		nonUnits = append(nonUnits, cnf)
	default:
		units = append(units, cnf)
	}
	return fac.And(units...), fac.And(nonUnits...)
}

func (c *compiler) generateDtree(fac f.Factory, hdl handler.Handler) (dtree, handler.State) {
	if c.nonUnitClauses.IsAtomic() {
		return nil, succ
	}
	tree, state := generateMinFillDtree(fac, c.nonUnitClauses, hdl)
	if !state.Success {
		return nil, state
	}
	tree.initialize(c.solver)
	return tree, succ
}

func (c *compiler) compileWithTree(dtree dtree, hdl handler.Handler) (f.Formula, handler.State) {
	if c.nonUnitClauses.IsAtomic() {
		return c.originalCNF, succ
	}
	if !c.solver.Start() {
		return c.fac.Falsum(), succ
	}
	c.initializeCaches(dtree)
	if e := event.DnnfComputationStarted; !hdl.ShouldResume(e) {
		return 0, handler.Cancelation(e)
	}

	result, state := c.cnf2ddnnf(dtree, hdl)
	if !state.Success {
		return 0, state
	} else {
		return c.fac.And(c.unitClauses, result), state
	}
}

func (c *compiler) initializeCaches(dtree dtree) {
	depth := dtree.depth() + 1
	sep := dtree.widestSeparator() + 1
	variables := f.Variables(c.fac, c.originalCNF).Size()

	c.localCacheKeys = make([][]*bitset, depth)
	for i := range depth {
		c.localCacheKeys[i] = make([]*bitset, sep)
		for j := range sep {
			c.localCacheKeys[i][j] = newBitset(dtree.size() + int32(variables))
		}
	}
	c.localOccurrences = make([][][]int32, depth)
	for i := range depth {
		c.localOccurrences[i] = make([][]int32, sep)
		for j := range sep {
			c.localOccurrences[i][j] = make([]int32, variables)
			for k := range variables {
				c.localOccurrences[i][j][k] = -1
			}
		}
	}
}

func (c *compiler) cnf2ddnnf(tree dtree, hdl handler.Handler) (f.Formula, handler.State) {
	return c.cnf2ddnnfInner(tree, 0, hdl)
}

func (c *compiler) cnf2ddnnfInner(tree dtree, currentShannons int, hdl handler.Handler) (f.Formula, handler.State) {
	separator := tree.dynamicSeparator()
	implied := c.newlyImpliedLiterals(tree.staticVarSet())

	if separator.cardinality() == 0 {
		switch node := tree.(type) {
		case *dtreeLeaf:
			return c.fac.And(implied, c.leaf2ddnnf(node)), succ
		default:
			return c.conjoin(implied, node.(*dtreeNode), currentShannons, hdl)
		}
	} else {
		variable := c.chooseShannonVariable(tree, separator, currentShannons)
		if e := event.DnnfShannonExpansion; !hdl.ShouldResume(e) {
			return 0, handler.Cancelation(e)
		}

		positiveDnnf := c.fac.Falsum()
		if c.solver.Decide(variable, true) {
			res, state := c.cnf2ddnnfInner(tree, currentShannons+1, hdl)
			if !state.Success {
				return 0, state
			} else {
				positiveDnnf = res
			}
		}
		c.solver.UndoDecide(variable)
		if positiveDnnf.Sort() == f.SortFalse {
			if c.solver.AtAssertionLevel() && c.solver.AssertCdLiteral() {
				return c.cnf2ddnnf(tree, hdl)
			} else {
				return c.fac.Falsum(), succ
			}
		}

		negativeDnnf := c.fac.Falsum()
		if c.solver.Decide(variable, false) {
			res, state := c.cnf2ddnnfInner(tree, currentShannons+1, hdl)
			if !state.Success {
				return 0, state
			} else {
				negativeDnnf = res
			}
		}
		c.solver.UndoDecide(variable)
		if negativeDnnf == c.fac.Falsum() {
			if c.solver.AtAssertionLevel() && c.solver.AssertCdLiteral() {
				return c.cnf2ddnnf(tree, hdl)
			} else {
				return c.fac.Falsum(), succ
			}
		}

		lit := c.solver.LitForIdx(variable)
		positiveBranch := c.fac.And(lit, positiveDnnf)
		negativeBranch := c.fac.And(lit.Negate(c.fac), negativeDnnf)
		return c.fac.And(implied, c.fac.Or(positiveBranch, negativeBranch)), succ
	}
}

func (c *compiler) chooseShannonVariable(tree dtree, separator *bitset, currentShannons int) int32 {
	occurrences := c.localOccurrences[tree.depth()][currentShannons]
	for i := range occurrences {
		if separator.get(i) {
			occurrences[i] = 0
		} else {
			occurrences[i] = -1
		}
	}
	tree.countUnsubsumedOccurrences(occurrences)

	max := int32(-1)
	maxVal := int32(-1)
	for i := separator.nextSetBit(0); i != -1; i = separator.nextSetBit(i + 1) {
		val := occurrences[i]
		if val > maxVal {
			max = i
			maxVal = val
		}
	}
	return max
}

func (c *compiler) conjoin(
	implied f.Formula, tree *dtreeNode, currentShannons int, hdl handler.Handler,
) (f.Formula, handler.State) {
	if implied.Sort() == f.SortFalse {
		return c.fac.Falsum(), succ
	}
	left, state := c.cnfAux(tree.left, currentShannons, hdl)
	if !state.Success {
		return 0, state
	}
	if left.Sort() == f.SortFalse {
		return c.fac.Falsum(), succ
	}
	right, state := c.cnfAux(tree.right, currentShannons, hdl)
	if !state.Success {
		return 0, state
	}
	if right.Sort() == f.SortFalse {
		return c.fac.Falsum(), succ
	}
	return c.fac.And(implied, left, right), succ
}

func (c *compiler) cnfAux(tree dtree, currentShannons int, hdl handler.Handler) (f.Formula, handler.State) {
	switch node := tree.(type) {
	case *dtreeLeaf:
		return c.leaf2ddnnf(node), succ
	default:
		key := c.computeCacheKey(tree.(*dtreeNode), currentShannons)
		if val, ok := c.cache.Get(key); ok {
			return val.(f.Formula), succ
		} else {
			dnnf, state := c.cnf2ddnnf(tree, hdl)
			if !state.Success {
				return 0, state
			}
			if dnnf.Sort() != f.SortFalse {
				c.cache.Put(key.clone(), dnnf)
			}
			return dnnf, succ
		}
	}
}

func (c *compiler) computeCacheKey(tree *dtreeNode, currentShannons int) *bitset {
	key := c.localCacheKeys[tree.depth()][currentShannons]
	key.clear()
	tree.cacheKey(key, c.numberOfVariables)
	return key
}

func (c *compiler) leaf2ddnnf(leaf *dtreeLeaf) f.Formula {
	var leafResultOperands []f.Formula
	var leafCurrentLiterals []f.Literal
	index := 0
	for _, lit := range f.Literals(c.fac, leaf.clause).Content() {
		switch c.solver.ValueOf(sat.MkLit(c.solver.VariableIndex(lit), lit.IsNeg())) {
		case f.TristateTrue:
			return c.fac.Verum()
		case f.TristateUndef:
			leafCurrentLiterals = append(leafCurrentLiterals, lit)
			leafResultOperands = append(leafResultOperands, c.fac.And(f.LiteralsAsFormulas(leafCurrentLiterals)...))
			leafCurrentLiterals[index] = lit.Negate(c.fac)
			index++
		}
	}
	return c.fac.Or(leafResultOperands...)
}

func (c *compiler) newlyImpliedLiterals(knownVariables *bitset) f.Formula {
	return c.solver.NewlyImplied(knownVariables.bits)
}
