package bdd

import (
	"math/big"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/normalform"
)

// A BDD is a canonical representation of a Boolean formula. It contains a
// pointer to the kernel which was used to generate the BDD and the node
// index of the BDD within this kernel.
type BDD struct {
	Kernel *Kernel
	Index  int32
}

func newBdd(index int32, kernel *Kernel) *BDD {
	return &BDD{kernel, index}
}

// Compile creates a BDD for a given formula.  The variable ordering in this
// case is the order in which the variables occur in the formula.
func Compile(fac f.Factory, formula f.Formula) *BDD {
	bdd, _ := CompileWithHandler(fac, formula, nil)
	return bdd
}

// CompileWithHandler creates a BDD for a given formula with the given
// bddHandler.  The handler can abort the BDD creation based on the number of
// nodes created during the BDD compilation process.  If the BDD compilation
// was aborted, the ok flag is false.
func CompileWithHandler(fac f.Factory, formula f.Formula, bddHandler Handler) (bdd *BDD, ok bool) {
	handler.Start(bddHandler)
	varNum := int32(f.Variables(fac, formula).Size())
	kernel := NewKernel(fac, varNum, varNum*30, varNum*20)
	bddIndex, ok := compile(fac, formula, kernel, bddHandler)
	if !ok {
		return nil, ok
	} else {
		return newBdd(bddIndex, kernel), ok
	}
}

// CompileWithVarOrder creates a BDD for a given formula and a variable
// ordering.
func CompileWithVarOrder(fac f.Factory, formula f.Formula, order []f.Variable) *BDD {
	bdd, _ := CompileWithVarOrderAndHandler(fac, formula, order, nil)
	return bdd
}

// CompileWithVarOrderAndHandler creates a BDD for a given formula, variable
// ordering, and bddHandler.  The handler can abort the BDD creation based on
// the number of nodes created during the BDD compilation process.  If the BDD
// compilation was aborted, the ok flag is false.
func CompileWithVarOrderAndHandler(
	fac f.Factory,
	formula f.Formula,
	order []f.Variable,
	bddHandler Handler,
) (bdd *BDD, ok bool) {
	handler.Start(bddHandler)
	varNum := len(order)
	kernel := NewKernelWithOrdering(fac, order, int32(varNum)*30, int32(varNum)*20)
	bddIndex, ok := compile(fac, formula, kernel, bddHandler)
	if !ok {
		return nil, ok
	} else {
		return newBdd(bddIndex, kernel), ok
	}
}

// CompileWithKernel creates a BDD for a given formula with a given kernel.
func CompileWithKernel(fac f.Factory, formula f.Formula, kernel *Kernel) *BDD {
	bdd, _ := CompileWithKernelAndHandler(fac, formula, kernel, nil)
	return bdd
}

// CompileWithKernelAndHandler creates a BDD for a given formula with a given
// kernel and bddHandler.  The handler can abort the BDD creation based on the
// number of nodes created during the BDD compilation process.  If the BDD
// compilation was aborted, the ok flag is false.
func CompileWithKernelAndHandler(
	fac f.Factory,
	formula f.Formula,
	kernel *Kernel,
	bddHandler Handler,
) (bdd *BDD, ok bool) {
	handler.Start(bddHandler)
	bddIndex, ok := compile(fac, formula, kernel, bddHandler)
	if !ok {
		return nil, ok
	} else {
		return newBdd(bddIndex, kernel), ok
	}
}

// CompileLiterals creates a BDD for a conjunction of literals with a given
// kernel.
func CompileLiterals(literals []f.Literal, kernel *Kernel) *BDD {
	var bdd int32
	if len(literals) == 0 {
		bdd = bddFalse
	} else if len(literals) == 1 {
		lit := literals[0]
		variable := lit.Variable()
		idx := kernel.getOrAddVarIndex(variable)
		if lit.IsPos() {
			bdd, _ = kernel.ithVar(idx)
		} else {
			bdd, _ = kernel.nithVar(idx)
		}
	} else {
		lit := literals[0]
		variable := lit.Variable()
		idx := kernel.getOrAddVarIndex(variable)
		if lit.IsPos() {
			bdd, _ = kernel.ithVar(idx)
		} else {
			bdd, _ = kernel.nithVar(idx)
		}
		for i := 1; i < len(literals); i++ {
			lit = literals[i]
			variable := lit.Variable()
			idx = kernel.getOrAddVarIndex(variable)
			var operand int32
			if lit.IsPos() {
				operand, _ = kernel.ithVar(idx)
			} else {
				operand, _ = kernel.nithVar(idx)
			}
			previous := bdd
			var ok bool
			bdd, ok = kernel.addRef(kernel.and(bdd, operand), nil)
			if !ok {
				panic(errorx.IllegalState("bdd generation was aborted by handler"))
			}
			kernel.delRef(previous)
			kernel.delRef(operand)
		}
	}
	return newBdd(bdd, kernel)
}

func compile(fac f.Factory, formula f.Formula, kernel *Kernel, handler Handler) (int32, bool) {
	switch formula.Sort() {
	case f.SortFalse:
		return bddFalse, true
	case f.SortTrue:
		return bddTrue, true
	case f.SortLiteral:
		variable := f.Literal(formula).Variable()
		idx := kernel.getOrAddVarIndex(variable)
		if formula.IsPos() {
			return kernel.ithVar(idx)
		} else {
			return kernel.nithVar(idx)
		}
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		operand, ok := compile(fac, op, kernel, handler)
		if !ok {
			return 0, false
		}
		res, ok := kernel.addRef(kernel.not(operand), handler)
		kernel.delRef(operand)
		return res, ok
	case f.SortImpl, f.SortEquiv:
		l, r, _ := fac.BinaryLeftRight(formula)
		left, ok := compile(fac, l, kernel, handler)
		if !ok {
			return 0, false
		}
		right, ok := compile(fac, r, kernel, handler)
		if !ok {
			return 0, false
		}
		var res int32
		if formula.Sort() == f.SortImpl {
			res, ok = kernel.addRef(kernel.implication(left, right), handler)
		} else {
			res, ok = kernel.addRef(kernel.equivalence(left, right), handler)
		}
		kernel.delRef(left)
		kernel.delRef(right)
		return res, ok
	case f.SortAnd, f.SortOr:
		ops, _ := fac.NaryOperands(formula)
		res, ok := compile(fac, ops[0], kernel, handler)
		if !ok {
			return 0, false
		}
		for i := 1; i < len(ops); i++ {
			operand, ok := compile(fac, ops[i], kernel, handler)
			if !ok {
				return 0, false
			}
			previous := res
			if formula.Sort() == f.SortAnd {
				res, ok = kernel.addRef(kernel.and(res, operand), handler)
			} else {
				res, ok = kernel.addRef(kernel.or(res, operand), handler)
			}
			kernel.delRef(previous)
			kernel.delRef(operand)
			if !ok {
				return 0, false
			}
		}
		return res, ok
	case f.SortCC, f.SortPBC:
		return compile(fac, normalform.NNF(fac, formula), kernel, handler)
	default:
		panic(errorx.UnknownEnumValue(formula.Sort()))
	}
}

// ToFormula returns a formula representation of the BDD. This is done by using
// the Shannon expansion. If followPathsToTrue is activated, the paths leading
// to the true terminal are followed to generate the formula. If
// followPathsToTrue is deactivated, the paths leading to the false terminal
// are followed to generate the formula and the resulting formula is negated.
// Depending on the formula and the number of satisfying assignments, the
// generated formula can be more compact using the true paths or false paths,
// respectively.
func (b *BDD) ToFormula(fac f.Factory, followPathsToTrue ...bool) f.Formula {
	var fptt bool
	if len(followPathsToTrue) > 1 {
		fptt = followPathsToTrue[0]
	} else {
		fptt = true
	}
	return b.Kernel.toFormula(fac, b.Index, fptt)
}

// Negate returns a newBdd BDD which is the negation of the BDD.
func (b *BDD) Negate() *BDD {
	bdd, _ := b.Kernel.addRef(b.Kernel.not(b.Index), nil)
	return newBdd(bdd, b.Kernel)
}

// Implies returns a newBdd BDD which is the implication of the BDD to the given
// other BDD.  This method panics if the BDDs were constructed by different
// kernels.
func (b *BDD) Implies(other *BDD) *BDD {
	if other.Kernel != b.Kernel {
		panic(errorx.BadInput("other BDD and receiver BDD have different kernels"))
	}
	bdd, _ := b.Kernel.addRef(b.Kernel.implication(b.Index, other.Index), nil)
	return newBdd(bdd, b.Kernel)
}

// ImpliedBy returns a newBdd BDD which is the implication of the other given BDD
// to the BDD.  This method panics if the BDDs were constructed by different
// kernels.
func (b *BDD) ImpliedBy(other *BDD) *BDD {
	if other.Kernel != b.Kernel {
		panic(errorx.BadInput("other BDD and receiver BDD have different kernels"))
	}
	bdd, _ := b.Kernel.addRef(b.Kernel.implication(other.Index, b.Index), nil)
	return newBdd(bdd, b.Kernel)
}

// Equivalence returns a newBdd BDD which is the equivalence of the BDD and the
// other given BDD.  This method panics if the BDDs were constructed by
// different kernels.
func (b *BDD) Equivalence(other *BDD) *BDD {
	if other.Kernel != b.Kernel {
		panic(errorx.BadInput("other BDD and receiver BDD have different kernels"))
	}
	bdd, _ := b.Kernel.addRef(b.Kernel.equivalence(b.Index, other.Index), nil)
	return newBdd(bdd, b.Kernel)
}

// And returns a newBdd BDD which is the conjunction of the BDD and the given
// other BDD.  This method panics if the BDDs were constructed by different
// kernels.
func (b *BDD) And(other *BDD) *BDD {
	if other.Kernel != b.Kernel {
		panic(errorx.BadInput("other BDD and receiver BDD have different kernels"))
	}
	bdd, _ := b.Kernel.addRef(b.Kernel.and(b.Index, other.Index), nil)
	return newBdd(bdd, b.Kernel)
}

// Or returns a newBdd BDD which is the disjunction of the BDD and the given other
// BDD.  This method panics if the BDDs were constructed by different kernels.
func (b *BDD) Or(other *BDD) *BDD {
	if other.Kernel != b.Kernel {
		panic(errorx.BadInput("other BDD and receiver BDD have different kernels"))
	}
	bdd, _ := b.Kernel.addRef(b.Kernel.or(b.Index, other.Index), nil)
	return newBdd(bdd, b.Kernel)
}

// IsTautology reports whether the BDD is a tautology.
func (b *BDD) IsTautology() bool {
	return b.Index == bddTrue
}

// IsContradiction reports whether the BDD is a tautology.
func (b *BDD) IsContradiction() bool {
	return b.Index == bddFalse
}

// ModelCount returns the number of satisfying models of the BDD.
func (b *BDD) ModelCount() *big.Int {
	return b.Kernel.satCount(b.Index)
}

// ModelEnumeration enumerates all models of the BDD wrt. a given set of variables.
func (b *BDD) ModelEnumeration(variables ...f.Variable) []*model.Model {
	return bddModelEnum(b, variables)
}

// CNF returns a CNF formula for the BDD.
func (b *BDD) CNF() f.Formula {
	return cnf(b)
}

// DNF returns a DNF formula for the BDD.
func (b *BDD) DNF() f.Formula {
	return dnf(b)
}

// NumberOfCNFClauses returns the number of clauses for the CNF formula of the
// BDD.
func (b *BDD) NumberOfCNFClauses() *big.Int {
	return b.Kernel.pathCountZero(b.Index)
}

// Restrict returns a newBdd BDD where the literals of the restriction are
// assigned to their respective polarity and therefore the BDD does not contain
// the respective variables anymore.
func (b *BDD) Restrict(restriction ...f.Literal) *BDD {
	resBdd := CompileLiterals(restriction, b.Kernel)
	return newBdd(b.Kernel.restrict(b.Index, resBdd.Index), b.Kernel)
}

// Exists performs existential quantifier elimination for a given set of
// variables and return the resulting BDD.
func (b *BDD) Exists(variable ...f.Variable) *BDD {
	resBdd := CompileLiterals(f.VariablesAsLiterals(variable), b.Kernel)
	return newBdd(b.Kernel.exists(b.Index, resBdd.Index), b.Kernel)
}

// ForAll performs universal quantifier elimination for a given set of
// variables and returns the resulting BDD.
func (b *BDD) ForAll(variable ...f.Variable) *BDD {
	resBdd := CompileLiterals(f.VariablesAsLiterals(variable), b.Kernel)
	return newBdd(b.Kernel.forAll(b.Index, resBdd.Index), b.Kernel)
}

// Model returns an arbitrary model of the BDD.  An error is returned if the
// BDD is a contradiction and therefore has no model.
func (b *BDD) Model() (*model.Model, error) {
	return b.createModel(b.Kernel.satOne(b.Index))
}

// ModelWithVariables returns an arbitrary model of the BDD which contains at
// least the given variables. If a variable is a don't care variable, it will
// be assigned with the given defaultValue.  An error is returned if the BDD is
// a contradiction and therefore has no model.
func (b *BDD) ModelWithVariables(defaultValue bool, variable ...f.Variable) (*model.Model, error) {
	bdd := CompileLiterals(f.VariablesAsLiterals(variable), b.Kernel)
	var pol int32
	if defaultValue {
		pol = bddTrue
	} else {
		pol = bddFalse
	}
	modelBdd := b.Kernel.satOneSet(b.Index, bdd.Index, pol)
	return b.createModel(modelBdd)
}

// FullModel returns a model over all variables of the BDD.  An error is
// returned if the BDD is a contradiction and therefore has no model.
func (b *BDD) FullModel() (*model.Model, error) {
	return b.createModel(b.Kernel.fullSatOne(b.Index))
}

// PathCountOne returns the number of paths leading to the terminal 1 node.
func (b *BDD) PathCountOne() *big.Int {
	return b.Kernel.pathCountOne(b.Index)
}

// PathCountZero returns the number of paths leading to the terminal 0 node.
func (b *BDD) PathCountZero() *big.Int {
	return b.Kernel.pathCountZero(b.Index)
}

// Support returns all the variables the BDD depends on.
func (b *BDD) Support() []f.Variable {
	supportBdd := b.Kernel.support(b.Index)
	model, err := b.createModel(supportBdd)
	if err != nil {
		return []f.Variable{}
	} else {
		return model.PosVars()
	}
}

// NodeCount returns the number of distinct nodes for the BDD.
func (b *BDD) NodeCount() int {
	return b.Kernel.nodeCount(b.Index)
}

// VariableProfile returns how often each variable occurs in the BDD.
func (b *BDD) VariableProfile() map[f.Variable]int {
	varProfile := b.Kernel.varProfile(b.Index)
	profile := make(map[f.Variable]int, len(varProfile))
	for i := 0; i < len(varProfile); i++ {
		variable, _ := b.Kernel.getVariableForIndex(int32(i))
		profile[variable] = varProfile[i]
	}
	return profile
}

// VariableOrder returns the variable order of the BDD.
func (b *BDD) VariableOrder() []f.Variable {
	order := make([]f.Variable, len(b.Kernel.level2var)-1)
	for i := 0; i < len(order); i++ {
		variable, _ := b.Kernel.getVariableForIndex(b.Kernel.level2var[i])
		order[i] = variable
	}
	return order
}

func (b *BDD) createModel(modelBdd int32) (*model.Model, error) {
	if modelBdd == bddFalse {
		return nil, errorx.BadInput("the BDD has no model")
	}
	if modelBdd == bddTrue {
		return model.New(), nil
	}
	nodes := b.Kernel.allNodes(modelBdd)
	lits := make([]f.Literal, len(nodes))
	for i, node := range nodes {
		variable, _ := b.Kernel.getVariableForIndex(node[1])
		if node[2] == bddFalse {
			lits[i] = variable.AsLiteral()
		} else if node[3] == bddFalse {
			lits[i] = variable.Negate(b.Kernel.fac)
		} else {
			panic(errorx.IllegalState("model must have a unique path through the BDD"))
		}
	}
	return model.New(lits...), nil
}
