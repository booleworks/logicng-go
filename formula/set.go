package formula

import (
	"booleworks.com/logicng/errorx"
	"slices"
)

type (
	VarSet     = fset[Variable]
	LitSet     = fset[Literal]
	FormulaSet = fset[Formula]
)

// FormulaSet represents a set of formulas/variables/literals.  In a formula
// set the contained formulas are unique.  The order within the set is not
// deterministic.  When using the Content method, the formulas sorted by their
// unique ID and are therefore deterministic.
type fset[T Formula | Variable | Literal] struct {
	elements map[T]present
}

// NewFormulaSet generates a new formula set with the given formulas as
// content.
func NewFormulaSet(formula ...Formula) *FormulaSet {
	fs := &FormulaSet{make(map[Formula]present, len(formula))}
	for _, f := range formula {
		fs.elements[f] = present{}
	}
	return fs
}

// NewVarSet generates a new variable set with the given variables as
// content.
func NewVarSet(variable ...Variable) *VarSet {
	fs := &VarSet{make(map[Variable]present, len(variable))}
	for _, f := range variable {
		fs.elements[f] = present{}
	}
	return fs
}

// NewLitSet generates a new literal set with the given literal as content.
func NewLitSet(literal ...Literal) *LitSet {
	fs := &LitSet{make(map[Literal]present, len(literal))}
	for _, f := range literal {
		fs.elements[f] = present{}
	}
	return fs
}

// NewVariableSetCopy returns a new variable set with the content of all the
// given variable sets as content.
func NewVariableSetCopy(variableSet ...*VarSet) *VarSet {
	fs := &VarSet{make(map[Variable]present)}
	for _, set := range variableSet {
		for f := range set.elements {
			fs.elements[f] = present{}
		}
	}
	return fs
}

// Add adds a new element to the set.
func (f *fset[T]) Add(element T) {
	f.elements[element] = present{}
}

// AddAll adds the content of the given other set to the set.
func (f *fset[T]) AddAll(other *fset[T]) {
	for formula := range other.elements {
		f.elements[formula] = present{}
	}
}

// AddAllElements adds the other elements to the set.
func (f *fset[T]) AddAllElements(elements *[]T) {
	for _, e := range *elements {
		f.elements[e] = present{}
	}
}

// Remove removes the given element from the set.
func (f *fset[T]) Remove(element T) {
	delete(f.elements, element)
}

// RemoveAll removes all the content from the other given set from the other
// set from the set.
func (f *fset[T]) RemoveAll(elements *fset[T]) {
	for formula := range elements.elements {
		delete(f.elements, formula)
	}
}

// RemoveAllElements removes the content of the other set from the set.
func (f *fset[T]) RemoveAllElements(elements *[]T) {
	for _, formula := range *elements {
		delete(f.elements, formula)
	}
}

// Each takes a function which is executed on each element of the set.
func (f *fset[T]) Each(function func(index int, element T)) {
	count := 0
	for element := range f.elements {
		function(count, element)
		count++
	}
}

// Any return any element from the set or returns an error if the set is empty.
func (f *fset[T]) Any() (T, error) {
	if len(f.elements) == 0 {
		return 0, errorx.BadInput("empty formula set")
	}
	for lit := range f.elements {
		return lit, nil
	}
	return 0, nil
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
	slice := make([]T, len(f.elements))
	i := 0
	for k := range f.elements {
		slice[i] = k
		i++
	}
	slices.Sort(slice)
	return slice
}

// Variables returns all variables of the given formula as a variable set.
func Variables(fac Factory, formula ...Formula) *VarSet {
	if len(formula) == 1 {
		return variables(fac, formula[0])
	}
	result := NewVarSet()
	for _, f := range formula {
		result.AddAll(variables(fac, f))
	}
	return result
}

func variables(fac Factory, formula Formula) *VarSet {
	cached, ok := LookupFunctionCache(fac, FuncVariables, formula)
	if ok {
		return cached.(*VarSet)
	}

	result := NewVarSet()
	switch fsort := formula.Sort(); fsort {
	case SortFalse, SortTrue:
		break
	case SortLiteral:
		variable := Literal(formula).Variable()
		result.Add(variable)
	case SortNot:
		op, _ := fac.NotOperand(formula)
		result.AddAll(Variables(fac, op))
	case SortImpl, SortEquiv:
		left, right, _ := fac.BinaryLeftRight(formula)
		result.AddAll(Variables(fac, left))
		result.AddAll(Variables(fac, right))
	case SortOr, SortAnd:
		ops, _ := fac.NaryOperands(formula)
		for _, op := range ops {
			result.AddAll(Variables(fac, op))
		}
	case SortCC, SortPBC:
		_, _, lits, _, _ := fac.PBCOps(formula)
		for _, lit := range lits {
			variable := lit.Variable()
			result.Add(variable)
		}
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
	SetFunctionCache(fac, FuncVariables, formula, result)
	return result
}

// Literals returns all literals of the given formula as a literal set.
func Literals(fac Factory, formula Formula) *LitSet {
	cached, ok := LookupFunctionCache(fac, FuncLiterals, formula)
	if ok {
		return cached.(*LitSet)
	}

	result := NewLitSet()
	switch fsort := formula.Sort(); fsort {
	case SortFalse, SortTrue:
		break
	case SortLiteral:
		result.Add(Literal(formula))
	case SortNot:
		op, _ := fac.NotOperand(formula)
		result.AddAll(Literals(fac, op))
	case SortImpl, SortEquiv:
		left, right, _ := fac.BinaryLeftRight(formula)
		result.AddAll(Literals(fac, left))
		result.AddAll(Literals(fac, right))
	case SortOr, SortAnd:
		ops, _ := fac.NaryOperands(formula)
		for _, op := range ops {
			result.AddAll(Literals(fac, op))
		}
	case SortCC, SortPBC:
		_, _, lits, _, _ := fac.PBCOps(formula)
		for _, lit := range lits {
			result.Add(lit)
		}
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
	SetFunctionCache(fac, FuncLiterals, formula, result)
	return result
}
