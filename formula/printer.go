package formula

import (
	"fmt"
	"strings"

	"github.com/booleworks/logicng-go/errorx"
)

// PrintSymbols gathers all symbols for Boolean and pseudo-Boolean formulas for
// printing formulas.  You can use the DefaultSymbols method to use exactly
// these symbols which can be parsed again be the default parser, or you can
// define your own symbols.
type PrintSymbols struct {
	Verum          string // default: $true
	Falsum         string // default: $false
	Not            string // default: ~
	Implication    string // default: =>
	Equivalence    string // default: <=>
	And            string // default: &
	Or             string // default: |
	LeftBracket    string // default: (
	RightBracket   string // default: )
	Plus           string // default: +
	Minus          string // default: -
	Multiplication string // default: *
	Equal          string // default: =
	Less           string // default: <
	LessOrEqual    string // default: <=
	Greater        string // default: >
	GreaterOrEqual string // default: >=
}

// DefaultSymbols returns the standard symbols for printing formulas.
func DefaultSymbols() *PrintSymbols {
	return &PrintSymbols{
		Verum:          "$true",
		Falsum:         "$false",
		Not:            "~",
		Implication:    " => ",
		Equivalence:    " <=> ",
		And:            " & ",
		Or:             " | ",
		LeftBracket:    "(",
		RightBracket:   ")",
		Plus:           " + ",
		Minus:          "-",
		Multiplication: "*",
		Equal:          " = ",
		Less:           " < ",
		LessOrEqual:    " <= ",
		Greater:        " > ",
		GreaterOrEqual: " >= ",
	}
}

// Sprint prints a formula in human-readable form.  In order to this, the
// generating formula factor must be passed as a parameter.
//
// This method panics if the formula cannot be found on the factory.
func (f Formula) Sprint(fac Factory) string {
	return toInnerString(fac, f, fac.Symbols())
}

// Sprint prints a variable in human-readable form.  In order to this, the
// generating formula factor must be passed as a parameter.
//
// This method panics if the variable cannot be found on the factory.
func (v Variable) Sprint(fac Factory) string {
	return toInnerString(fac, Formula(v), fac.Symbols())
}

// Sprint prints a literal in human-readable form.  In order to this, the
// generating formula factor must be passed as a parameter.
//
// This method panics if the literal cannot be found on the factory.
func (l Literal) Sprint(fac Factory) string {
	return toInnerString(fac, Formula(l), fac.Symbols())
}

func toInnerString(fac Factory, formula Formula, s *PrintSymbols) string {
	var printString string
	var err bool
	switch fsort := formula.Sort(); fsort {
	case SortTrue:
		printString = s.Verum
	case SortFalse:
		printString = s.Falsum
	case SortLiteral:
		name, phase, found := fac.LiteralNamePhase(formula)
		err = !found
		if phase {
			printString = name
		} else {
			printString = s.Not + name
		}
	case SortNot:
		not, found := fac.NotOperand(formula)
		err = !found
		printString = s.Not + formatBracket(fac, not, s)
	case SortImpl, SortEquiv:
		left, right, found := fac.BinaryLeftRight(formula)
		err = !found
		if fsort == SortImpl {
			printString = formatBinaryOperator(fac, fsort, left, right, s.Implication, s)
		} else {
			printString = formatBinaryOperator(fac, fsort, left, right, s.Equivalence, s)
		}
	case SortAnd, SortOr:
		ops, found := fac.NaryOperands(formula)
		err = !found
		if fsort == SortAnd {
			printString = formatNaryOperator(fac, fsort, ops, s.And, s)
		} else {
			printString = formatNaryOperator(fac, fsort, ops, s.Or, s)
		}
	case SortCC, SortPBC:
		comparator, rhs, literals, coefficients, found := fac.PBCOps(formula)
		err = !found
		printString = formatPBC(fac, comparator, rhs, literals, coefficients, s)
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
	if err {
		panic(errorx.UnknownFormula(formula))
	}
	return printString
}

func formatBracket(fac Factory, f Formula, s *PrintSymbols) string {
	return s.LeftBracket + toInnerString(fac, f, s) + s.RightBracket
}

func formatBinaryOperator(fac Factory, fsort FSort, left, right Formula, sym string, s *PrintSymbols) string {
	var leftString, rightString string
	if fsort > left.Sort() {
		leftString = toInnerString(fac, left, s)
	} else {
		leftString = formatBracket(fac, left, s)
	}
	if fsort > right.Sort() {
		rightString = toInnerString(fac, right, s)
	} else {
		rightString = formatBracket(fac, right, s)
	}
	return leftString + sym + rightString
}

func formatNaryOperator(fac Factory, fsort FSort, operands []Formula, sym string, s *PrintSymbols) string {
	var sb strings.Builder
	size := len(operands)
	var last Formula
	for i, op := range operands {
		if i == size-1 {
			last = op
		} else {
			if fsort > op.Sort() {
				sb.WriteString(toInnerString(fac, op, s))
			} else {
				sb.WriteString(formatBracket(fac, op, s))
			}
			sb.WriteString(sym)
		}
	}
	if last != 0 {
		if fsort > last.Sort() {
			sb.WriteString(toInnerString(fac, last, s))
		} else {
			sb.WriteString(formatBracket(fac, last, s))
		}
	}
	return sb.String()
}

func formatPBC(fac Factory, comp CSort, rhs int, lits []Literal, coeffs []int, s *PrintSymbols) string {
	var sb strings.Builder
	mul := s.Multiplication
	add := s.Plus
	numOps := len(lits)
	for i := 0; i < numOps-1; i++ {
		if coeffs[i] != 1 {
			sb.WriteString(fmt.Sprintf("%d%s%s%s", coeffs[i], mul, lits[i].Sprint(fac), add))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s", lits[i].Sprint(fac), add))
		}
	}
	if numOps > 0 {
		if coeffs[numOps-1] != 1 {
			sb.WriteString(fmt.Sprintf("%d%s%s", coeffs[numOps-1], mul, lits[numOps-1].Sprint(fac)))
		} else {
			sb.WriteString(lits[numOps-1].Sprint(fac))
		}
	}
	sb.WriteString(fmt.Sprintf("%s%d", formatPBComparator(comp, s), rhs))
	return sb.String()
}

func formatPBComparator(comparator CSort, s *PrintSymbols) string {
	switch comparator {
	case EQ:
		return s.Equal
	case LE:
		return s.LessOrEqual
	case LT:
		return s.Less
	case GE:
		return s.GreaterOrEqual
	case GT:
		return s.Greater
	default:
		panic(errorx.UnknownEnumValue(comparator))
	}
}
