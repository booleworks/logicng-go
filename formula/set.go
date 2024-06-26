package formula

import (
	"slices"

	"github.com/booleworks/logicng-go/errorx"
)

type (
	VarSet            = fset[Variable]        // Immutable set of variables
	LitSet            = fset[Literal]         // Immutable set of literals
	FormulaSet        = fset[Formula]         // Immutable set of formulas
	MutableVarSet     = mutableFset[Variable] // Mutable set of variables
	MutableLitSet     = mutableFset[Literal]  // Mutable set of literals
	MutableFormulaSet = mutableFset[Formula]  // Mutable set of formulas
)

type fset[T Formula | Variable | Literal] struct {
	elements map[T]present
	content  []T
}

type mutableFset[T Formula | Variable | Literal] struct {
	*fset[T]
}

// NewFormulaSet generates a new formula set with the given formulas as
// content.
func NewFormulaSet(formula ...Formula) *FormulaSet {
	fs := &FormulaSet{make(map[Formula]present, len(formula)), nil}
	for _, f := range formula {
		fs.elements[f] = present{}
	}
	fs.setContent()
	return fs
}

// NewVarSet generates a new variable set with the given variables as
// content.
func NewVarSet(variable ...Variable) *VarSet {
	fs := &VarSet{make(map[Variable]present, len(variable)), nil}
	for _, f := range variable {
		fs.elements[f] = present{}
	}
	fs.setContent()
	return fs
}

// NewLitSet generates a new literal set with the given literal as content.
func NewLitSet(literal ...Literal) *LitSet {
	fs := &LitSet{make(map[Literal]present, len(literal)), nil}
	for _, f := range literal {
		fs.elements[f] = present{}
	}
	fs.setContent()
	return fs
}

// NewVarSetCopy returns a new variable set with the content of all the given
// variable sets as content.
func NewVarSetCopy(variableSet ...*VarSet) *VarSet {
	fs := &VarSet{make(map[Variable]present), nil}
	for _, set := range variableSet {
		for f := range set.elements {
			fs.elements[f] = present{}
		}
	}
	fs.setContent()
	return fs
}

// NewMutableFormulaSet generates a new mutable formula set with the given
// formulas as content.
func NewMutableFormulaSet(formula ...Formula) *MutableFormulaSet {
	return &MutableFormulaSet{NewFormulaSet(formula...)}
}

// NewMutableVarSet generates a new mutable variable set with the given
// variables as content.
func NewMutableVarSet(variable ...Variable) *MutableVarSet {
	return &MutableVarSet{NewVarSet(variable...)}
}

// NewMutableLitSet generates a new mutable literal set with the given literal
// as content.
func NewMutableLitSet(literal ...Literal) *MutableLitSet {
	return &MutableLitSet{NewLitSet(literal...)}
}

// NewMutableVarSetCopy returns a new mutable variable set with the
// content of all the given variable sets as content.
func NewMutableVarSetCopy(variableSet ...*VarSet) *MutableVarSet {
	return &MutableVarSet{NewVarSetCopy(variableSet...)}
}

// Empty reports whether the set is empty.
func (f *fset[T]) Empty() bool {
	return len(f.elements) == 0
}

// Size returns the size of the set.
func (f *fset[T]) Size() int {
	return len(f.elements)
}

// Contains reports whether the given element is in the set.
func (f *fset[T]) Contains(element T) bool {
	_, ok := f.elements[element]
	return ok
}

// ContainsAll reports whether all elements from the given set are in the set.
func (f *fset[T]) ContainsAll(elements *fset[T]) bool {
	for formula := range elements.elements {
		_, ok := f.elements[formula]
		if !ok {
			return false
		}
	}
	return true
}

// Content returns the content of the set as a sorted (by ID) slice.
func (f *fset[T]) Content() []T {
	if f.content == nil {
		f.setContent()
	}
	return f.content
}

// Add adds a new element to the set.
func (f *mutableFset[T]) Add(element T) {
	f.elements[element] = present{}
	f.content = nil
}

// AddAll adds the content of the given other set to the set.
func (f *mutableFset[T]) AddAll(other *fset[T]) {
	for formula := range other.elements {
		f.elements[formula] = present{}
	}
	f.content = nil
}

// AddAllElements adds the other elements to the set.
func (f *mutableFset[T]) AddAllElements(elements *[]T) {
	for _, e := range *elements {
		f.elements[e] = present{}
	}
	f.content = nil
}

// Remove removes the given element from the set.
func (f *mutableFset[T]) Remove(element T) {
	delete(f.elements, element)
	f.content = nil
}

// RemoveAll removes all the content from the other given set from the other
// set from the set.
func (f *mutableFset[T]) RemoveAll(elements *fset[T]) {
	for formula := range elements.elements {
		delete(f.elements, formula)
	}
	f.content = nil
}

// RemoveAllElements removes the content of the other set from the set.
func (f *mutableFset[T]) RemoveAllElements(elements *[]T) {
	for _, formula := range *elements {
		delete(f.elements, formula)
	}
	f.content = nil
}

// AsImmutable returns the mutable set as an immutable set.
func (f *mutableFset[T]) AsImmutable() *fset[T] {
	return f.fset
}

func (f *fset[T]) setContent() {
	f.content = make([]T, len(f.elements))
	i := 0
	for k := range f.elements {
		f.content[i] = k
		i++
	}
	slices.Sort(f.content)
}

// Variables returns all variables of the given formula as a variable set.
func Variables(fac Factory, formula ...Formula) *VarSet {
	if len(formula) == 1 {
		return variables(fac, formula[0])
	}
	result := NewMutableVarSet()
	for _, f := range formula {
		result.AddAll(variables(fac, f))
	}
	return result.AsImmutable()
}

func variables(fac Factory, formula Formula) *VarSet {
	cached, ok := LookupFunctionCache(fac, FuncVariables, formula)
	if ok {
		return cached.(*VarSet)
	}

	vars := NewMutableVarSet()
	switch fsort := formula.Sort(); fsort {
	case SortFalse, SortTrue:
		break
	case SortLiteral:
		variable := Literal(formula).Variable()
		vars.Add(variable)
	case SortNot:
		op, _ := fac.NotOperand(formula)
		vars.AddAll(Variables(fac, op))
	case SortImpl, SortEquiv:
		left, right, _ := fac.BinaryLeftRight(formula)
		vars.AddAll(Variables(fac, left))
		vars.AddAll(Variables(fac, right))
	case SortOr, SortAnd:
		ops, _ := fac.NaryOperands(formula)
		for _, op := range ops {
			vars.AddAll(Variables(fac, op))
		}
	case SortCC, SortPBC:
		_, _, lits, _, _ := fac.PBCOps(formula)
		for _, lit := range lits {
			variable := lit.Variable()
			vars.Add(variable)
		}
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
	result := vars.AsImmutable()
	SetFunctionCache(fac, FuncVariables, formula, result)
	return result
}

// Literals returns all literals of the given formula as a literal set.
func Literals(fac Factory, formula ...Formula) *LitSet {
	if len(formula) == 1 {
		return literals(fac, formula[0])
	}
	result := NewMutableLitSet()
	for _, f := range formula {
		result.AddAll(literals(fac, f))
	}
	return result.AsImmutable()
}

func literals(fac Factory, formula Formula) *LitSet {
	cached, ok := LookupFunctionCache(fac, FuncLiterals, formula)
	if ok {
		return cached.(*LitSet)
	}

	lits := NewMutableLitSet()
	switch fsort := formula.Sort(); fsort {
	case SortFalse, SortTrue:
		break
	case SortLiteral:
		lits.Add(Literal(formula))
	case SortNot:
		op, _ := fac.NotOperand(formula)
		lits.AddAll(Literals(fac, op))
	case SortImpl, SortEquiv:
		left, right, _ := fac.BinaryLeftRight(formula)
		lits.AddAll(Literals(fac, left))
		lits.AddAll(Literals(fac, right))
	case SortOr, SortAnd:
		ops, _ := fac.NaryOperands(formula)
		for _, op := range ops {
			lits.AddAll(Literals(fac, op))
		}
	case SortCC, SortPBC:
		_, _, lts, _, _ := fac.PBCOps(formula)
		for _, lit := range lts {
			lits.Add(lit)
		}
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
	result := lits.AsImmutable()
	SetFunctionCache(fac, FuncLiterals, formula, result)
	return result
}
