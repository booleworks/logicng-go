package formula

import (
	"fmt"

	"booleworks.com/logicng/errorx"
)

// Formula represents a Boolean (or pseudo-Boolean) formula in LogicNG.  The
// datatype itself is just an alias to an uint32.  The first four bits encode
// the formula sort whereas the remaining 28 bit are a unique identifier for
// the formula on the factory.  Therefore, a formula is only valid and useful
// in the context of the factory which was used to generate it.
type Formula uint32

// Variable represents a Boolean variable in LogicNG.  This is just a type
// alias for the uint32 wrapped by the Formula type.  But for type-safety
// reasons it is often desirable to explicitly know that something is a
// variable.  You can convert between variable and formula with the AsFormula
// and AsVariable methods.
type Variable uint32

// Literal represents a Boolean literal in LogicNG.  This is just a type
// alias for the uint32 wrapped by the Formula type.  But for type-safety
// reasons it is often desirable to explicitly know that something is a
// literal.  You can convert between literal and formula with the AsFormula
// and AsLiteral methods.
type Literal uint32

// EncodeVariable takes a unique ID an returns its encoding as a variable.
func EncodeVariable(id uint32) Variable {
	return Variable((uint32(SortLiteral) << 28) | (id & idMask))
}

// AsFormula returns the variable as a formula type.
func (v Variable) AsFormula() Formula {
	return Formula(v)
}

// AsLiteral returns the variable as a literal type.
func (v Variable) AsLiteral() Literal {
	return Literal(v)
}

// ID returns the extracted ID of the variable.
func (v Variable) ID() uint32 {
	return uint32(v) & idMask
}

// Negate returns the negation of the variable.
func (v Variable) Negate(fac Factory) Literal {
	name, _ := fac.VarName(v)
	return fac.Lit(name, false)
}

// EncodeLiteral takes a unique ID an returns its encoding as a literal.
func EncodeLiteral(id uint32) Literal {
	return Literal((uint32(SortLiteral) << 28) | (id & idMask))
}

// AsFormula returns the literal as a formula type.
func (l Literal) AsFormula() Formula {
	return Formula(l)
}

// ID returns the extracted ID of the literal.
func (l Literal) ID() uint32 {
	return uint32(l) & idMask
}

// IsPos reports whether the literal is positive.
func (l Literal) IsPos() bool {
	return l.ID()%2 == 1
}

// IsNeg reports whether the literal is negative.
func (l Literal) IsNeg() bool {
	return l.ID()%2 == 0
}

// Variable extracts the variable from a literal.
func (l Literal) Variable() Variable {
	if l.IsPos() {
		return Variable(l)
	} else {
		return EncodeVariable(negId(l.ID()))
	}
}

// Negate returns the negation of the literal.
func (l Literal) Negate(fac Factory) Literal {
	name, phase, _ := fac.LitNamePhase(l)
	return fac.Lit(name, !phase)
}

// VariablesAsLiterals returns a list of variables as a list of literals.
func VariablesAsLiterals(variables []Variable) []Literal {
	literals := make([]Literal, len(variables))
	for i, v := range variables {
		literals[i] = Literal(v)
	}
	return literals
}

// VariablesAsFormulas returns a list of variables as a list of formulas.
func VariablesAsFormulas(variables []Variable) []Formula {
	literals := make([]Formula, len(variables))
	for i, v := range variables {
		literals[i] = Formula(v)
	}
	return literals
}

// LiteralsAsVariables returns a list of literals as a list of variables or
// returns an error if there is a negative literal in the list.
func LiteralsAsVariables(literals []Literal) ([]Variable, error) {
	variables := make([]Variable, len(literals))
	for i, l := range literals {
		if l.IsNeg() {
			return nil, errorx.BadInput("negative literal")
		}
		variables[i] = Variable(l)
	}
	return variables, nil
}

// LiteralsAsFormulas returns a list of literals as a list of formulas.
func LiteralsAsFormulas(literals []Literal) []Formula {
	formuals := make([]Formula, len(literals))
	for i, l := range literals {
		formuals[i] = Formula(l)
	}
	return formuals
}

// AsVariable returns the current formula as a variable type or returns an
// error if the formula is not a variable.
func (f Formula) AsVariable() (Variable, error) {
	if f.Sort() != SortLiteral || !f.IsPos() {
		return 0, errorx.BadFormulaSort(f.Sort())
	}
	return Variable(f), nil
}

// AsLiteral returns the current formula as a literal type or returns an
// error if the formula is not a literal.
func (f Formula) AsLiteral() (Literal, error) {
	if f.Sort() != SortLiteral {
		return 0, errorx.BadFormulaSort(f.Sort())
	}
	return Literal(f), nil
}

// FSort encodes the different formula sorts.
type FSort uint32

const (
	SortFalse   FSort = iota // constant false
	SortTrue                 // constant true
	SortLiteral              // literal (variable + phase)
	SortNot                  // negation
	SortAnd                  // conjunction
	SortOr                   // disjunction
	SortImpl                 // implication
	SortEquiv                // equivalence
	SortCC                   // cardinality constraint
	SortPBC                  // pseudo-Boolean constraint
)

// CSort encodes the sort of an integer comparator
type CSort byte

const (
	EQ CSort = iota // equal
	LE              // less or equal
	LT              // less than
	GE              // greater or equal
	GT              // greater than
)

// Evaluate evaluates an integer comparison given its left- and right-hand side.
func (c CSort) Evaluate(lhs, rhs int) bool {
	switch c {
	case EQ:
		return lhs == rhs
	case LE:
		return lhs <= rhs
	case LT:
		return lhs < rhs
	case GE:
		return lhs >= rhs
	case GT:
		return lhs > rhs
	default:
		panic(errorx.UnknownEnumValue(c))
	}
}

//go:generate stringer -type=FSort
//go:generate stringer -type=CSort

const (
	typeMask uint32 = 0xF0000000
	idMask   uint32 = 0x0FFFFFFF
)

// DualSort returns the dual sort for AND and OR.  For all other formula types
// it returns an error.
func DualSort(fsort FSort) (FSort, error) {
	switch fsort {
	case SortAnd:
		return SortOr, nil
	case SortOr:
		return SortAnd, nil
	default:
		return 0, errorx.BadFormulaSort(fsort)
	}
}

// EncodeFormula takes a formula sort and a unique ID an returns its encoding.
func EncodeFormula(fsort FSort, id uint32) Formula {
	return Formula((uint32(fsort) << 28) | (id & idMask))
}

// Sort returns the extracted sort of a formula.
func (f Formula) Sort() FSort {
	return FSort((uint32(f) & typeMask) >> 28)
}

// ID returns the extracted ID of a formula.
func (f Formula) ID() uint32 {
	return uint32(f) & idMask
}

// String returns a simple string representation of a formula.  Since the
// default string method does not know the constructing factory, only the
// formula sort and ID can be printed.  If you want a human-readable string
// output, use the Sprint function which takes the factory as parameter.
func (f Formula) String() string {
	return fmt.Sprintf("{sort: %s, id: %d}", f.Sort(), f.ID())
}

// Negate returns the negation of the formula.
func (f Formula) Negate(fac Factory) Formula {
	return fac.Not(f)
}

// IsPos reports whether the formula is positive, i.e. the true constant, a
// positive literal, or any formula which is not a negation.
func (f Formula) IsPos() bool {
	return f.ID()%2 == 1
}

// IsNeg reports whether the formula is negative, i.e. the false constant, a
// negative literal, or a negation.
func (f Formula) IsNeg() bool {
	return f.ID()%2 == 0
}

// IsConstant reports whether the formula is a constant.
func (f Formula) IsConstant() bool {
	return f.Sort() <= SortTrue
}

// IsAtomic reports whether the formula is an atomic formula, i.e. a constant,
// a literal, or a pseudo-Boolean constraint including cardinality constraints.
func (f Formula) IsAtomic() bool {
	fsort := f.Sort()
	return fsort <= SortLiteral || fsort == SortCC || fsort == SortPBC
}

// Comparator compares two formulas f1 and f2 with their integer
// value.
func Comparator(a, b interface{}) int {
	aFormula := a.(Literal)
	bFormula := b.(Literal)
	switch {
	case aFormula > bFormula:
		return 1
	case aFormula < bFormula:
		return -1
	default:
		return 0
	}
}
