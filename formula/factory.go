package formula

import (
	"fmt"
	"strings"

	"github.com/booleworks/logicng-go/configuration"
	"github.com/booleworks/logicng-go/errorx"
)

// A Factory is the central concept of LogicNG and is always required
// when working with the library. A formula factory is an object consisting of
// two major components:
//
//  1. a factory, which creates formulas
//  2. a container, which stores created formulas.
//
// Formula factories are not thread safe!
type Factory interface {
	Verum() Formula
	Falsum() Formula
	Constant(value bool) Formula
	Var(name string) Variable
	Variable(name string) Formula
	Vars(name ...string) []Variable
	Lit(name string, phase bool) Literal
	Literal(name string, phase bool) Formula
	Not(operand Formula) Formula
	BinaryOperator(sort FSort, left, right Formula) (Formula, error)
	Implication(left, right Formula) Formula
	Equivalence(left Formula, right Formula) Formula
	NaryOperator(sort FSort, operands ...Formula) (Formula, error)
	And(operands ...Formula) Formula
	Minterm(operands ...Literal) Formula
	Or(operands ...Formula) Formula
	Clause(operands ...Literal) Formula
	CC(comparator CSort, rhs uint32, variables ...Variable) Formula
	AMO(variables ...Variable) Formula
	EXO(variables ...Variable) Formula
	PBC(comparator CSort, rhs int, literals []Literal, coefficients []int) Formula

	VarName(variable Variable) (name string, found bool)
	LitNamePhase(literal Literal) (name string, phase, found bool)
	LiteralNamePhase(formula Formula) (name string, phase, found bool)
	NotOperand(formula Formula) (op Formula, found bool)
	BinaryLeftRight(formula Formula) (left, right Formula, found bool)
	NaryOperands(formula Formula) (ops []Formula, found bool)
	PBCOps(formula Formula) (comparator CSort, rhs int, literals []Literal, coefficients []int, found bool)
	Operands(formula Formula) []Formula

	NewAuxVar(sort AuxVarSort) Variable

	transformationCacheEntry(entry TransformationCacheSort) *map[Formula]Formula
	predicateCacheEntry(entry PredicateCacheSort) *map[Formula]bool
	functionCacheEntry(entry FunctionCacheSort) *map[Formula]any
	ConfigurationFor(sort configuration.Sort) (configuration.Config, bool)
	PutConfiguration(configuration configuration.Config) error
	Symbols() *PrintSymbols
	SetPrintSymbols(symbols *PrintSymbols)

	Statistics() string
}

// A CachingFactory is the default (and currently only) implementation of the
// formula factory in LogicNG.
//
// In this implementation, the container function is 'smart': A formula factory
// guarantees that syntactically equivalent formulas are created only once.
// This mechanism also extends to variants of the formula in terms of
// associativity and commutativity. Therefore, if the user creates formulas for
//
//	A & B
//	B & A
//	(B & A)
//
// all of them are represented by only one formula in memory. This approach is
// only possible, because formulas in LogicNG are immutable data structures. So
// once created, a formula can never be altered again.
type CachingFactory struct {
	cFalse Formula
	cTrue  Formula
	id     uint32

	literals     map[Literal]literal
	nots         map[Formula]not
	implications map[Formula]binaryOp
	equivalences map[Formula]binaryOp
	ands         map[Formula]naryOp
	ors          map[Formula]naryOp
	ccs          map[Formula]pbc
	pbcs         map[Formula]pbc

	posLitCache map[string]Variable
	negLitCache map[string]Literal
	notCache    map[Formula]Formula
	implCache   map[fpair]Formula
	equivCache  map[fpair]Formula
	andCache    map[uint64][]Formula
	orCache     map[uint64][]Formula
	ccCache     map[uint64][]Formula
	pbcCache    map[uint64][]Formula

	transformationCache map[TransformationCacheSort]map[Formula]Formula
	predicateCache      map[PredicateCacheSort]map[Formula]bool
	functionCache       map[FunctionCacheSort]map[Formula]any

	auxVarCounters map[AuxVarSort]int

	configurations map[configuration.Sort]configuration.Config
	symbols        *PrintSymbols
	conserveVars   bool
}

// NewFactory returns a new caching formula factory.  If the optional
// conserveVars flag is set to true, trivial contradictions and tautologies are
// not simplified in formulas. E.g. a formula like A & ~A, A | ~A, or A => A
// can be generated on the formula factory.  Therefore, it is guaranteed that
// all variables of the original formula are still present on the formula
// factory. If set to false, the formulas will be simplified to false or true
// and therefore variables of the original formula can not be present on the
// formula factory.  The default behaviour is that the flag is set to false.
func NewFactory(conserveVars ...bool) Factory {
	f := &CachingFactory{
		cFalse:              EncodeFormula(SortFalse, 0),
		cTrue:               EncodeFormula(SortTrue, 1),
		id:                  2,
		literals:            make(map[Literal]literal),
		nots:                make(map[Formula]not),
		implications:        make(map[Formula]binaryOp),
		equivalences:        make(map[Formula]binaryOp),
		ands:                make(map[Formula]naryOp),
		ors:                 make(map[Formula]naryOp),
		ccs:                 make(map[Formula]pbc),
		pbcs:                make(map[Formula]pbc),
		posLitCache:         make(map[string]Variable),
		negLitCache:         make(map[string]Literal),
		notCache:            make(map[Formula]Formula),
		implCache:           make(map[fpair]Formula),
		equivCache:          make(map[fpair]Formula),
		andCache:            make(map[uint64][]Formula),
		orCache:             make(map[uint64][]Formula),
		ccCache:             make(map[uint64][]Formula),
		pbcCache:            make(map[uint64][]Formula),
		transformationCache: make(map[TransformationCacheSort]map[Formula]Formula),
		predicateCache:      make(map[PredicateCacheSort]map[Formula]bool),
		functionCache:       make(map[FunctionCacheSort]map[Formula]any),
		auxVarCounters:      make(map[AuxVarSort]int),
		configurations:      make(map[configuration.Sort]configuration.Config),
		symbols:             DefaultSymbols(),
		conserveVars:        conserveVars != nil && conserveVars[0],
	}
	return f
}

// Verum returns the Boolean true constant.
func (fac *CachingFactory) Verum() Formula {
	return fac.cTrue
}

// Falsum returns the Boolean false constant.
func (fac *CachingFactory) Falsum() Formula {
	return fac.cFalse
}

// Constant returns the Boolean constant represented by the given value.
func (fac *CachingFactory) Constant(value bool) Formula {
	if value {
		return fac.cTrue
	}
	return fac.cFalse
}

// Var returns a Boolean variable with the given name.  In contrast to the
// Variable method, a variable type is returned.
func (fac *CachingFactory) Var(name string) Variable {
	variable, ok := fac.posLitCache[name]
	if !ok {
		id := fac.nextPosId()
		variable, _ = EncodeFormula(SortLiteral, id).AsVariable()
		fac.posLitCache[name] = variable
		fac.literals[Literal(variable)] = literal{name, true}
	}
	return variable
}

// Variable returns a Boolean variable with the given name as a Formula.
func (fac *CachingFactory) Variable(name string) Formula {
	return fac.Var(name).AsFormula()
}

// Vars returns a list of Boolean variable with the given names.
func (fac *CachingFactory) Vars(name ...string) []Variable {
	variables := make([]Variable, len(name))
	for i := range name {
		variables[i] = fac.Var(name[i])
	}
	return variables
}

// Lit returns a Boolean literal with the given name and phase. For a name "A"
// the positive phase returns the variable "A" whereas the negative phase
// returns the negated variable "~A".  In contrast to the Literal function, a
// literal type is returned.
func (fac *CachingFactory) Lit(name string, phase bool) Literal {
	variable := fac.Variable(name)
	if phase {
		return Literal(variable)
	}
	lit, ok := fac.negLitCache[name]
	if !ok {
		id := negId(variable.ID())
		if id > fac.id {
			fac.id = id
		}
		lit = EncodeLiteral(id)
		fac.negLitCache[name] = lit
		fac.literals[lit] = literal{name, false}
	}
	return lit
}

// Literal returns a Boolean literal with the given name and phase. For a name "A"
// the positive phase returns the variable "A" whereas the negative phase
// returns the negated variable "~A".
func (fac *CachingFactory) Literal(name string, phase bool) Formula {
	return fac.Lit(name, phase).AsFormula()
}

// Not returns the negation of the given formula.  Constants, literals and
// negations are negated directly and returned. For all other formulas a new
// Not formula is returned.
func (fac *CachingFactory) Not(operand Formula) Formula {
	switch {
	case operand == fac.cFalse:
		return fac.cTrue
	case operand == fac.cTrue:
		return fac.cFalse
	case operand.Sort() == SortLiteral:
		lit := fac.literals[Literal(operand)]
		return fac.Literal(lit.name, !lit.phase)
	case operand.Sort() == SortNot:
		not := fac.nots[operand]
		return not.operand
	default:
		neg, ok := fac.notCache[operand]
		if !ok {
			id := negId(operand.ID())
			if id > fac.id {
				fac.id = id
			}
			neg = EncodeFormula(SortNot, id)
			fac.notCache[operand] = neg
			fac.nots[neg] = not{operand}
		}
		return neg
	}
}

// BinaryOperator returns a new binary operator with the given sort and the two
// operands left and right.  Returns an error if the given sort is not a binary
// operator (implication or equivalence).
func (fac *CachingFactory) BinaryOperator(sort FSort, left, right Formula) (Formula, error) {
	switch sort {
	case SortImpl:
		return fac.Implication(left, right), nil
	case SortEquiv:
		return fac.Equivalence(left, right), nil
	default:
		return 0, errorx.BadFormulaSort(sort)
	}
}

// Implication returns an implication left => right.
func (fac *CachingFactory) Implication(left, right Formula) Formula {
	switch {
	case left == fac.cFalse, right == fac.cTrue:
		return fac.cTrue
	case left == fac.cTrue:
		return right
	case right == fac.cFalse, left.ID() == negId(right.ID()):
		return fac.Not(left)
	case !fac.conserveVars && left == right:
		return fac.cTrue
	default:
		key := fpair{left, right}
		impl, ok := fac.implCache[key]
		if !ok {
			id := fac.nextPosId()
			impl = EncodeFormula(SortImpl, id)
			fac.implCache[key] = impl
			fac.implications[impl] = binaryOp{left, right}
		}
		return impl
	}
}

// Equivalence returns an equivalence left <=> right.
func (fac *CachingFactory) Equivalence(left, right Formula) Formula {
	switch {
	case left == fac.cTrue:
		return right
	case right == fac.cTrue:
		return left
	case left == fac.cFalse:
		return fac.Not(right)
	case right == fac.cFalse:
		return fac.Not(left)
	case !fac.conserveVars && left == right:
		return fac.cTrue
	case !fac.conserveVars && left.ID() == negId(right.ID()):
		return fac.cFalse
	default:
		key := fpair{left, right}
		equiv, ok := fac.equivCache[key]
		if !ok {
			id := fac.nextPosId()
			equiv = EncodeFormula(SortEquiv, id)
			fac.equivCache[key] = equiv
			fac.equivalences[equiv] = binaryOp{left, right}
		}
		return equiv
	}
}

// NaryOperator returns a new n-ary operator with the given sort and the list
// of operands.  Returns an error if the given sort is not an n-ary operator
// (conjunction or disjunction).
func (fac *CachingFactory) NaryOperator(sort FSort, operands ...Formula) (Formula, error) {
	switch sort {
	case SortAnd:
		return fac.And(operands...), nil
	case SortOr:
		return fac.Or(operands...), nil
	default:
		return 0, errorx.BadFormulaSort(sort)
	}
}

// And returns a conjunction of the given operands.  If the result is an And
// formula, it is guaranteed to have at least two operands.  An empty
// conjunction is treated as true, a conjunction with one operand, is treated
// as this operand.
func (fac *CachingFactory) And(operands ...Formula) Formula {
	hash := hashOperands(operands)
	and, ok := fac.findAnd(hash, operands)
	if !ok {
		condensed, isCnf, isContradiction := fac.condenseOperandsAnd(operands...)
		switch {
		case isContradiction:
			return fac.cFalse
		case len(condensed) == 0:
			return fac.cTrue
		case len(condensed) == 1:
			return condensed[0]
		default:
			hash = hashOperands(condensed)
			and, ok = fac.findAnd(hash, condensed)
			if !ok {
				id := fac.nextPosId()
				and = EncodeFormula(SortAnd, id)
				fac.setCNFCaches(and, isCnf)
				fac.andCache[hash] = append(fac.andCache[hash], and)
				fac.ands[and] = naryOp{condensed}
			}
		}
	}
	return and
}

// Minterm returns a conjunction between the given literals.  In contrast to
// the [And] Method, the operands are not condensed.  Meaning, if the operands
// contain duplicate literals, the result will also contain these duplicates
// and therefore lead to potential problems down the road.
func (fac *CachingFactory) Minterm(operands ...Literal) Formula {
	hash := hashOperands(operands)
	formulas := LiteralsAsFormulas(operands)
	and, ok := fac.findAnd(hash, formulas)
	if !ok {
		switch {
		case len(operands) == 0:
			return fac.cTrue
		case len(operands) == 1:
			return operands[0].AsFormula()
		default:
			id := fac.nextPosId()
			and = EncodeFormula(SortAnd, id)
			fac.setCNFCaches(and, true)
			fac.andCache[hash] = append(fac.andCache[hash], and)
			fac.ands[and] = naryOp{formulas}
		}
	}
	return and
}

// Or returns a disjunction of the given operands.  If the result is an Or
// formula, it is guaranteed to have at least two operands.  An empty
// disjunction is treated as false, a disjunction with one operand, is treated
// as this operand.
func (fac *CachingFactory) Or(operands ...Formula) Formula {
	hash := hashOperands(operands)
	or, ok := fac.findOr(hash, operands)
	if !ok {
		condensed, isCnf, isTautology := fac.condenseOperandsOr(operands...)
		switch {
		case isTautology:
			return fac.cTrue
		case len(condensed) == 0:
			return fac.cFalse
		case len(condensed) == 1:
			return condensed[0]
		default:
			hash = hashOperands(condensed)
			or, ok = fac.findOr(hash, condensed)
			if !ok {
				id := fac.nextPosId()
				or = EncodeFormula(SortOr, id)
				fac.setCNFCaches(or, isCnf)
				fac.orCache[hash] = append(fac.orCache[hash], or)
				fac.ors[or] = naryOp{condensed}
			}
		}
	}
	return or
}

// Clause returns a disjunction between the given literals.  In contrast to
// the [Or] Method, the operands are not condensed.  Meaning, if the operands
// contain duplicate literals, the result will also contain these duplicates
// and therefore lead to potential problems down the road.
func (fac *CachingFactory) Clause(operands ...Literal) Formula {
	hash := hashOperands(operands)
	formulas := LiteralsAsFormulas(operands)
	or, ok := fac.findOr(hash, formulas)
	if !ok {
		switch {
		case len(operands) == 0:
			return fac.cFalse
		case len(operands) == 1:
			return operands[0].AsFormula()
		default:
			id := fac.nextPosId()
			or = EncodeFormula(SortOr, id)
			fac.setCNFCaches(or, true)
			fac.orCache[hash] = append(fac.orCache[hash], or)
			fac.ors[or] = naryOp{formulas}
		}
	}
	return or
}

// CC returns a new cardinality constraint with the given comparator and
// right-hand-side representing a constraint v_1 + ... + v_n C rhs with C in [<,
// >, <=, >=, =].
func (fac *CachingFactory) CC(comparator CSort, rhs uint32, variables ...Variable) Formula {
	return fac.constructCC(comparator, int(rhs), variables)
}

// AMO returns an at-most-one (<= 1) constraint over the given variables.
func (fac *CachingFactory) AMO(variables ...Variable) Formula {
	return fac.constructCC(LE, 1, variables)
}

// EXO returns an exactly-one (= 1) constraint over the given variables.
func (fac *CachingFactory) EXO(variables ...Variable) Formula {
	return fac.constructCC(EQ, 1, variables)
}

// PBC returns a pseudo-Boolean constraint with the given comparator and
// right-hand-side.  The single operands are represented by multiplications x_i
// * l_i with x_i from the coefficients and l_i from the literals.  The
// constraint represented is therefore x_1 * l_1 + ... + x_n * x_n C rhs with C
// in [<, >, <=, >=, =].
func (fac *CachingFactory) PBC(comparator CSort, rhs int, literals []Literal, coefficients []int) Formula {
	if len(literals) == 0 {
		return fac.Constant(evaluateTrivialPBConstraint(comparator, rhs))
	}
	if fac.isCC(comparator, rhs, literals, &coefficients) {
		variables, _ := LiteralsAsVariables(literals)
		return fac.constructCC(comparator, rhs, variables)
	}
	return fac.constructPBC(comparator, rhs, literals, coefficients)
}

func negId(id uint32) uint32 {
	return id ^ 1
}

func (fac *CachingFactory) nextPosId() uint32 {
	if fac.id%2 == 0 {
		fac.id = fac.id + 1
	} else {
		fac.id = fac.id + 2
	}
	return fac.id
}

func hashOperands[T Formula | Literal](operands []T) uint64 {
	hash := uint64(4243)
	for _, op := range operands {
		hash = ((hash << 5) + hash) + uint64(op)
	}
	return hash
}

func hashPbc(constraint *pbc) uint64 {
	hash := uint64(4243)
	hash *= uint64(constraint.rhs)
	hash *= uint64(constraint.comparator)
	for _, op := range constraint.literals {
		hash = ((hash << 5) + hash) + uint64(op)
	}
	for _, coeff := range constraint.coefficients {
		hash = ((hash << 5) + hash) + uint64(coeff)
	}
	return hash
}

func (fac *CachingFactory) condenseOperandsAnd(operands ...Formula) ([]Formula, bool, bool) {
	length := 0
	for _, op := range operands {
		if op.Sort() == SortAnd {
			length += len(fac.ands[op].operands)
		} else {
			length++
		}
	}

	opCache := make(map[uint32]present, length)
	ops := make([]Formula, 0, length)
	cnfCheck := true
	for _, form := range operands {
		if form.Sort() == SortAnd {
			for _, op := range fac.ands[form].operands {
				ret := fac.addFormulaAnd(&opCache, &ops, op)
				if ret == 0x00 {
					return nil, cnfCheck, true
				}
				if ret == 0x02 {
					cnfCheck = false
				}
			}
		} else {
			ret := fac.addFormulaAnd(&opCache, &ops, form)
			if ret == 0x00 {
				return nil, cnfCheck, true
			}
			if ret == 0x02 {
				cnfCheck = false
			}
		}
	}
	return ops, cnfCheck, false
}

func (fac *CachingFactory) addFormulaAnd(opCache *map[uint32]present, ops *[]Formula, formula Formula) byte {
	if formula == fac.cTrue {
		return 0x01
	} else if formula == fac.cFalse || fac.containsComplement(opCache, formula) {
		return 0x00
	}
	id := formula.ID()
	if _, ok := (*opCache)[id]; !ok {
		(*opCache)[id] = present{}
		*ops = append(*ops, formula)
	}
	if fac.isCNFClause(formula) {
		return 0x01
	}
	return 0x02
}

func (fac *CachingFactory) condenseOperandsOr(operands ...Formula) ([]Formula, bool, bool) {
	length := 0
	for _, op := range operands {
		if op.Sort() == SortOr {
			length += len(fac.ors[op].operands)
		} else {
			length++
		}
	}

	opCache := make(map[uint32]present, length)
	ops := make([]Formula, 0, length)
	cnfCheck := true
	for _, form := range operands {
		if form.Sort() == SortOr {
			for _, op := range fac.ors[form].operands {
				ret := fac.addFormulaOr(&opCache, &ops, op)
				if ret == 0x00 {
					return nil, cnfCheck, true
				}
				if ret == 0x02 {
					cnfCheck = false
				}
			}
		} else {
			ret := fac.addFormulaOr(&opCache, &ops, form)
			if ret == 0x00 {
				return nil, cnfCheck, true
			}
			if ret == 0x02 {
				cnfCheck = false
			}
		}
	}
	return ops, cnfCheck, false
}

func (fac *CachingFactory) addFormulaOr(opCache *map[uint32]present, ops *[]Formula, formula Formula) byte {
	if formula == fac.cFalse {
		return 0x01
	} else if formula == fac.cTrue || fac.containsComplement(opCache, formula) {
		return 0x00
	}
	id := formula.ID()
	if _, ok := (*opCache)[id]; !ok {
		(*opCache)[id] = present{}
		*ops = append(*ops, formula)
	}
	if formula.Sort() == SortLiteral {
		return 0x01
	}
	return 0x02
}

func (fac *CachingFactory) findAnd(hash uint64, ops []Formula) (Formula, bool) {
	ands, ok := fac.andCache[hash]
	if !ok {
		return 0, false
	}
	for _, and := range ands {
		if opsEquals(fac.ands[and].operands, ops) {
			return and, true
		}
	}
	return 0, false
}

func (fac *CachingFactory) findOr(hash uint64, ops []Formula) (Formula, bool) {
	ors, ok := fac.orCache[hash]
	if !ok {
		return 0, false
	}
	for _, or := range ors {
		if opsEquals(fac.ors[or].operands, ops) {
			return or, true
		}
	}
	return 0, false
}

func (fac *CachingFactory) findCC(hash uint64, constraint *pbc) (Formula, bool) {
	ccs, ok := fac.ccCache[hash]
	if !ok {
		return 0, false
	}
	for _, ccFormula := range ccs {
		cc := fac.ccs[ccFormula]
		if fac.pbcsEquals(&cc, constraint) {
			return ccFormula, true
		}
	}
	return 0, false
}

func (fac *CachingFactory) findPBC(hash uint64, constraint *pbc) (Formula, bool) {
	pbcs, ok := fac.pbcCache[hash]
	if !ok {
		return 0, false
	}
	for _, pbcFormula := range pbcs {
		pbc := fac.pbcs[pbcFormula]
		if fac.pbcsEquals(&pbc, constraint) {
			return pbcFormula, true
		}
	}
	return 0, false
}

func (fac *CachingFactory) isCC(comparator CSort, rhs int, literals []Literal, coefficients *[]int) bool {
	for _, lit := range literals {
		if !lit.IsPos() {
			return false
		}
	}
	if coefficients != nil {
		for _, c := range *coefficients {
			if c != 1 {
				return false
			}
		}
	}
	return comparator == LE && rhs >= 0 ||
		comparator == LT && rhs >= 1 ||
		comparator == GE && rhs >= 0 ||
		comparator == GT && rhs >= 0 ||
		comparator == EQ && rhs >= 0
}

func (fac *CachingFactory) constructCC(comparator CSort, rhs int, variables []Variable) Formula {
	if len(variables) == 0 {
		return fac.Constant(evaluateTrivialPBConstraint(comparator, rhs))
	}
	pbc := pbc{
		literals:     VariablesAsLiterals(variables),
		coefficients: make([]int, len(variables)),
		rhs:          rhs,
		comparator:   comparator,
	}
	for i := range variables {
		pbc.coefficients[i] = 1
	}
	hash := hashPbc(&pbc)
	constraint, ok := fac.findCC(hash, &pbc)
	if !ok {
		id := fac.nextPosId()
		constraint = EncodeFormula(SortCC, id)
		fac.ccCache[hash] = append(fac.ccCache[hash], constraint)
		fac.ccs[constraint] = pbc
	}
	return constraint
}

func (fac *CachingFactory) constructPBC(comparator CSort, rhs int, literals []Literal, coefficients []int) Formula {
	if len(literals) == 0 {
		return fac.Constant(evaluateTrivialPBConstraint(comparator, rhs))
	}
	pbc := pbc{
		literals:     literals,
		coefficients: coefficients,
		rhs:          rhs,
		comparator:   comparator,
	}
	hash := hashPbc(&pbc)
	constraint, ok := fac.findPBC(hash, &pbc)
	if !ok {
		id := fac.nextPosId()
		constraint = EncodeFormula(SortPBC, id)
		fac.pbcCache[hash] = append(fac.pbcCache[hash], constraint)
		fac.pbcs[constraint] = pbc
	}
	return constraint
}

func opsEquals[t comparable](o1, o2 []t) bool {
	if len(o1) != len(o2) {
		return false
	}
	for i, v := range o1 {
		if v != o2[i] {
			return false
		}
	}
	return true
}

func (fac *CachingFactory) pbcsEquals(p1, p2 *pbc) bool {
	return p1.comparator == p2.comparator &&
		p1.rhs == p2.rhs &&
		opsEquals(p1.literals, p2.literals) &&
		opsEquals(p1.coefficients, p2.coefficients)
}

func (fac *CachingFactory) containsComplement(opCache *map[uint32]present, formula Formula) bool {
	if fac.conserveVars {
		return false
	}
	negatedId := negId(formula.ID())
	_, exists := (*opCache)[negatedId]
	return exists
}

func evaluateTrivialPBConstraint(comparator CSort, rhs int) bool {
	switch comparator {
	case EQ:
		return rhs == 0
	case LE:
		return rhs >= 0
	case LT:
		return rhs > 0
	case GE:
		return rhs <= 0
	case GT:
		return rhs < 0
	default:
		panic(errorx.UnknownEnumValue(comparator))
	}
}

func (fac *CachingFactory) isCNFClause(formula Formula) bool {
	switch {
	case formula.Sort() == SortLiteral:
		return true
	case formula.Sort() == SortOr:
		for _, op := range fac.ors[formula].operands {
			if op.Sort() != SortLiteral {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (fac *CachingFactory) setCNFCaches(formula Formula, isCNF bool) {
	if isCNF {
		SetPredicateCache(fac, PredCNF, formula, true)
		SetTransformationCache(fac, TransCNFFactorization, formula, formula)
	} else {
		SetPredicateCache(fac, PredCNF, formula, false)
	}
}

// VarName returns the name of the given variable.  The ok flag indicates
// whether the variable was found on the factory or not.
func (fac *CachingFactory) VarName(variable Variable) (name string, ok bool) {
	lit, ok := fac.literals[variable.AsLiteral()]
	if ok {
		return lit.name, ok
	}
	return "", ok
}

// LitNamePhase returns the name and phase of the given literal.  The ok flag
// indicates whether the literal was found on the factory or not.
func (fac *CachingFactory) LitNamePhase(literal Literal) (name string, phase, ok bool) {
	lit, ok := fac.literals[literal]
	if ok {
		return lit.name, lit.phase, ok
	}
	return "", false, ok
}

// LiteralNamePhase returns the name and phase of a given formula interpreted as
// literal.  The ok flag is false when the given formula was not a literal, or
// it was not found on the factory.
func (fac *CachingFactory) LiteralNamePhase(formula Formula) (name string, phase, ok bool) {
	if formula.Sort() != SortLiteral {
		return "", false, false
	}
	lit, ok := fac.literals[Literal(formula)]
	if ok {
		return lit.name, lit.phase, ok
	}
	return "", false, ok
}

// NotOperand returns the operand of a given formula interpreted as negation.
// The ok flag indicates whether the negation was found on the factory or not.
func (fac *CachingFactory) NotOperand(formula Formula) (op Formula, ok bool) {
	not, ok := fac.nots[formula]
	if ok {
		return not.operand, ok
	}
	return 0, ok
}

// BinaryLeftRight returns the left and right operand of a given formula
// interpreted as a binary operator. The ok flag indicates whether the binary
// operator was found on the factory or not.
func (fac *CachingFactory) BinaryLeftRight(formula Formula) (left, right Formula, ok bool) {
	var binary binaryOp
	switch fsort := formula.Sort(); fsort {
	case SortImpl:
		binary, ok = fac.implications[formula]
	case SortEquiv:
		binary, ok = fac.equivalences[formula]
	default:
		ok = false
	}
	if ok {
		return binary.left, binary.right, true
	}
	return 0, 0, false
}

// NaryOperands returns the operands of a given formula interpreted as an n-ary
// operator. The ok flag indicates whether the n-ary operator was found on the
// factory or not.
func (fac *CachingFactory) NaryOperands(formula Formula) (ops []Formula, ok bool) {
	var nary naryOp
	switch fsort := formula.Sort(); fsort {
	case SortAnd:
		nary, ok = fac.ands[formula]
	case SortOr:
		nary, ok = fac.ors[formula]
	default:
		ok = false
	}
	if ok {
		return nary.operands, true
	}
	return nil, false
}

// PBCOps returns the comparator, right-hand-side, literals, and coefficients of
// a given formula interpreted as a pseudo-Boolean constraint. The ok flag
// indicates whether the constraint was found on the factory or not.
func (fac *CachingFactory) PBCOps(
	formula Formula,
) (comparator CSort, rhs int, literals []Literal, coefficients []int, found bool) {
	switch formula.Sort() {
	case SortCC:
		c, found := fac.ccs[formula]
		if found {
			return c.comparator, c.rhs, c.literals, c.coefficients, true
		}
	case SortPBC:
		c, found := fac.pbcs[formula]
		if found {
			return c.comparator, c.rhs, c.literals, c.coefficients, true
		}
	}
	return 0, 0, nil, nil, false
}

// Operands returns the operands of a given formula.  For a negation this is
// its operand, for an implication and equivalence its left and right operand
// (in this order), for n-ary operators their operands.  All other formulas
// have no operands.
func (fac *CachingFactory) Operands(formula Formula) []Formula {
	switch fsort := formula.Sort(); fsort {
	case SortFalse, SortTrue, SortLiteral, SortCC, SortPBC:
		return []Formula{}
	case SortNot:
		return []Formula{fac.nots[formula].operand}
	case SortImpl:
		impl := fac.implications[formula]
		return []Formula{impl.left, impl.right}
	case SortEquiv:
		equiv := fac.equivalences[formula]
		return []Formula{equiv.left, equiv.right}
	case SortAnd:
		and := fac.ands[formula]
		return and.operands
	case SortOr:
		or := fac.ors[formula]
		return or.operands
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
}

func (fac *CachingFactory) getPBCUnsafe(formula Formula) pbc {
	switch formula.Sort() {
	case SortCC:
		return fac.ccs[formula]
	default:
		return fac.pbcs[formula]
	}
}

// NewAuxVar generates and returns a new auxiliary variable of the given sort.
func (fac *CachingFactory) NewAuxVar(sort AuxVarSort) Variable {
	variable := fac.Var(fmt.Sprintf("%s%d", sort, fac.auxVarCounters[sort]))
	fac.auxVarCounters[sort]++
	return variable
}

func (fac *CachingFactory) transformationCacheEntry(entry TransformationCacheSort) *map[Formula]Formula {
	cache, ok := fac.transformationCache[entry]
	if !ok {
		cache = make(map[Formula]Formula)
		fac.transformationCache[entry] = cache
	}
	return &cache
}

func (fac *CachingFactory) predicateCacheEntry(entry PredicateCacheSort) *map[Formula]bool {
	cache, ok := fac.predicateCache[entry]
	if !ok {
		cache = make(map[Formula]bool)
		fac.predicateCache[entry] = cache
	}
	return &cache
}

func (fac *CachingFactory) functionCacheEntry(entry FunctionCacheSort) *map[Formula]any {
	cache, ok := fac.functionCache[entry]
	if !ok {
		cache = make(map[Formula]any)
		fac.functionCache[entry] = cache
	}
	return &cache
}

// ConfigurationFor returns the configuration for a given configuration sort.
// The ok flag indicates whether a config for the given sort was found in the
// factory.
func (fac *CachingFactory) ConfigurationFor(
	sort configuration.Sort,
) (config configuration.Config, ok bool) {
	config, ok = fac.configurations[sort]
	return
}

// PutConfiguration adds a configuration to the factory.
func (fac *CachingFactory) PutConfiguration(config configuration.Config) error {
	if config.Sort() == configuration.FormulaFactory {
		return errorx.BadInput("configurations for the formula factory itself can only be passed in the constructor")
	}
	fac.configurations[config.Sort()] = config
	return nil
}

// Symbols returns the print symbols for the factory.  These symbols are used
// when calling Sprint on a formula.
func (fac *CachingFactory) Symbols() *PrintSymbols {
	return fac.symbols
}

// SetPrintSymbols sets the symbols for printing formulas with Sprint.
func (fac *CachingFactory) SetPrintSymbols(symbols *PrintSymbols) {
	fac.symbols = symbols
}

// Statistics returns a statistic of the factory as a multi-line string.
func (fac *CachingFactory) Statistics() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Current ID:                %d\n", fac.id)
	fmt.Fprintf(&sb, "# Literals:                %d\n", len(fac.literals))
	fmt.Fprintf(&sb, "# Negations:               %d\n", len(fac.nots))
	fmt.Fprintf(&sb, "# Implications:            %d\n", len(fac.implications))
	fmt.Fprintf(&sb, "# Equivalences:            %d\n", len(fac.equivalences))
	fmt.Fprintf(&sb, "# Conjunctions:            %d\n", len(fac.ands))
	fmt.Fprintf(&sb, "# Disjunctions:            %d\n", len(fac.ors))
	fmt.Fprintf(&sb, "# Cardinality Constraints: %d\n", len(fac.ccs))
	fmt.Fprintf(&sb, "# PB Constraints:          %d\n", len(fac.pbcs))
	return sb.String()
}

type literal struct {
	name  string
	phase bool
}

type not struct {
	operand Formula
}

type binaryOp struct {
	left  Formula
	right Formula
}

type naryOp struct {
	operands []Formula
}

type pbc struct {
	literals     []Literal
	coefficients []int
	rhs          int
	comparator   CSort
}

type fpair struct {
	f1 Formula
	f2 Formula
}

type present struct{}
