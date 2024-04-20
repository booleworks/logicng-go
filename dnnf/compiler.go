package dnnf

import (
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/sat"
	"github.com/booleworks/logicng-go/simplification"
	"github.com/emirpasic/gods/maps/treemap"
)

// A Handler for a DNNF can abort the compilation of a DNNF.  The method
// ShannonExpansion is called after each performed Shannon expansion within the
// DNNF compiler.
type Handler interface {
	handler.Handler
	ShannonExpansion() bool
}

// Compile returns a compiled DNNF for the given formula.
func Compile(fac f.Factory, formula f.Formula) *DNNF {
	dnnf, _ := CompileWithHandler(fac, formula, nil)
	return dnnf
}

// CompileWithHandler returns a compiled DNNF for the given formula.  The
// handler can be used to abort the DNNF compilation.  If the DNNF compilation
// was aborted, the ok flag is false.
func CompileWithHandler(fac f.Factory, formula f.Formula, dnnfHandler Handler) (dnnf *DNNF, ok bool) {
	originalVariables := f.NewMutableVarSet()
	originalVariables.AddAll(f.Variables(fac, formula))
	cnf := normalform.CNF(fac, formula)
	originalVariables.AddAll(f.Variables(fac, cnf))
	simplifiedFormula := simplifyFormula(fac, cnf)
	compiler := newCompiler(fac, simplifiedFormula)
	dnnfFormula, ok := compiler.compile(dnnfHandler)
	if !ok {
		return nil, ok
	} else {
		return &DNNF{fac, dnnfFormula, originalVariables.AsImmutable()}, ok
	}
}

type compiler struct {
	fac               f.Factory
	cnf               f.Formula
	unitClauses       f.Formula
	nonUnitClauses    f.Formula
	solver            sat.DnnfSatSolver
	numberOfVariables int32
	cache             *treemap.Map
	handler           Handler
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
		cnf:               formula,
		unitClauses:       units,
		nonUnitClauses:    nonUnits,
		solver:            sat.NewDnnfSolver(fac, numVars),
		numberOfVariables: int32(numVars),
		cache:             treemap.NewWith(bitsetComp),
	}
	compiler.solver.Add(formula)
	return &compiler
}

func bitsetComp(a, b interface{}) int {
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

func (c *compiler) compile(dnnfHandler Handler) (f.Formula, bool) {
	if !sat.IsSatisfiable(c.fac, c.cnf) {
		return c.fac.Falsum(), true
	}
	dTree := c.generateDtree(c.fac)
	return c.compileWithTree(dTree, dnnfHandler)
}

func initializeClauses(fac f.Factory, cnf f.Formula) (f.Formula, f.Formula) {
	units := []f.Formula{}
	nonUnits := []f.Formula{}
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

func (c *compiler) generateDtree(fac f.Factory) dtree {
	if c.nonUnitClauses.IsAtomic() {
		return nil
	}
	tree := generateMinFillDtree(fac, c.nonUnitClauses)
	tree.initialize(c.solver)
	return tree
}

func (c *compiler) compileWithTree(dtree dtree, compilationHandler Handler) (f.Formula, bool) {
	if c.nonUnitClauses.IsAtomic() {
		return c.cnf, true
	}
	if !c.solver.Start() {
		return c.fac.Falsum(), true
	}
	c.initializeCaches(dtree)
	c.handler = compilationHandler
	handler.Start(compilationHandler)

	result, ok := c.cnf2ddnnf(dtree)
	c.handler = nil
	if !ok {
		return 0, false
	} else {
		return c.fac.And(c.unitClauses, result), true
	}
}

func (c *compiler) initializeCaches(dtree dtree) {
	depth := dtree.depth() + 1
	sep := dtree.widestSeparator() + 1
	variables := f.Variables(c.fac, c.cnf).Size()

	c.localCacheKeys = make([][]*bitset, depth)
	for i := 0; i < depth; i++ {
		c.localCacheKeys[i] = make([]*bitset, sep)
		for j := 0; j < sep; j++ {
			c.localCacheKeys[i][j] = newBitset(dtree.size() + int32(variables))
		}
	}
	c.localOccurrences = make([][][]int32, depth)
	for i := 0; i < depth; i++ {
		c.localOccurrences[i] = make([][]int32, sep)
		for j := 0; j < sep; j++ {
			c.localOccurrences[i][j] = make([]int32, variables)
			for k := 0; k < variables; k++ {
				c.localOccurrences[i][j][k] = -1
			}
		}
	}
}

func (c *compiler) cnf2ddnnf(tree dtree) (f.Formula, bool) {
	return c.cnf2ddnnfInner(tree, 0)
}

func (c *compiler) cnf2ddnnfInner(tree dtree, currentShannons int) (f.Formula, bool) {
	separator := tree.dynamicSeparator()
	implied := c.newlyImpliedLiterals(tree.staticVarSet())

	if separator.cardinality() == 0 {
		switch node := tree.(type) {
		case *dtreeLeaf:
			return c.fac.And(implied, c.leaf2ddnnf(node)), true
		default:
			return c.conjoin(implied, node.(*dtreeNode), currentShannons)
		}
	} else {
		variable := c.chooseShannonVariable(tree, separator, currentShannons)
		if c.handler != nil && !c.handler.ShannonExpansion() {
			return 0, false
		}

		positiveDnnf := c.fac.Falsum()
		if c.solver.Decide(variable, true) {
			res, ok := c.cnf2ddnnfInner(tree, currentShannons+1)
			if !ok {
				return 0, false
			} else {
				positiveDnnf = res
			}
		}
		c.solver.UndoDecide(variable)
		if positiveDnnf.Sort() == f.SortFalse {
			if c.solver.AtAssertionLevel() && c.solver.AssertCdLiteral() {
				return c.cnf2ddnnf(tree)
			} else {
				return c.fac.Falsum(), true
			}
		}

		negativeDnnf := c.fac.Falsum()
		if c.solver.Decide(variable, false) {
			res, ok := c.cnf2ddnnfInner(tree, currentShannons+1)
			if !ok {
				return 0, false
			} else {
				negativeDnnf = res
			}
		}
		c.solver.UndoDecide(variable)
		if negativeDnnf == c.fac.Falsum() {
			if c.solver.AtAssertionLevel() && c.solver.AssertCdLiteral() {
				return c.cnf2ddnnf(tree)
			} else {
				return c.fac.Falsum(), true
			}
		}

		lit := c.solver.LitForIdx(variable)
		positiveBranch := c.fac.And(lit, positiveDnnf)
		negativeBranch := c.fac.And(lit.Negate(c.fac), negativeDnnf)
		return c.fac.And(implied, c.fac.Or(positiveBranch, negativeBranch)), true
	}
}

func (c *compiler) chooseShannonVariable(tree dtree, separator *bitset, currentShannons int) int32 {
	occurrences := c.localOccurrences[tree.depth()][currentShannons]
	for i := 0; i < len(occurrences); i++ {
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

func (c *compiler) conjoin(implied f.Formula, tree *dtreeNode, currentShannons int) (f.Formula, bool) {
	if implied.Sort() == f.SortFalse {
		return c.fac.Falsum(), true
	}
	left, ok := c.cnfAux(tree.left, currentShannons)
	if !ok {
		return 0, false
	}
	if left.Sort() == f.SortFalse {
		return c.fac.Falsum(), true
	}
	right, ok := c.cnfAux(tree.right, currentShannons)
	if !ok {
		return 0, false
	}
	if right.Sort() == f.SortFalse {
		return c.fac.Falsum(), true
	}
	return c.fac.And(implied, left, right), true
}

func (c *compiler) cnfAux(tree dtree, currentShannons int) (f.Formula, bool) {
	switch node := tree.(type) {
	case *dtreeLeaf:
		return c.leaf2ddnnf(node), true
	default:
		key := c.computeCacheKey(tree.(*dtreeNode), currentShannons)
		if val, ok := c.cache.Get(key); ok {
			return val.(f.Formula), true
		} else {
			dnnf, ok := c.cnf2ddnnf(tree)
			if !ok {
				return 0, false
			}
			if dnnf.Sort() != f.SortFalse {
				c.cache.Put(key.clone(), dnnf)
			}
			return dnnf, true
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
	leafResultOperands := []f.Formula{}
	leafCurrentLiterals := []f.Literal{}
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
