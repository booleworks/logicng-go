package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
	"github.com/antlr4-go/antlr/v4"
)

// PropositionalParser is a parser for LogicNG formulas.
type PropositionalParser struct {
	fac f.Factory
}

// New generates a new parser for LogicNG formulas with the
// given factory.  All formulas will be generated on this factory.
func New(fac f.Factory) *PropositionalParser {
	return &PropositionalParser{fac}
}

// ParseUnsafe parses a string and returns the parsed formula.  If there was a
// parser error, it just panics.  If you are not 100% sure that your formula
// can be parsed, you should use the Parse method with proper error management.
func (p *PropositionalParser) ParseUnsafe(data string) f.Formula {
	parsed, err := p.parse(data)
	if err != nil {
		panic(err)
	}
	return parsed
}

// Parse parses a string and returns the parsed formula.  If there was a parser
// error, the error is returned.
func (p *PropositionalParser) Parse(data string) (f.Formula, error) {
	return p.parse(data)
}

func (p *PropositionalParser) parse(data string) (f.Formula, error) {
	if utf8.RuneCountInString(strings.TrimSpace(data)) == 0 {
		return p.fac.Verum(), nil
	}
	is := antlr.NewInputStream(data)
	errorListener := new(errorListener)
	lexer := NewLogicNGPropositionalLexer(is)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	lngParser := NewLogicNGPropositionalParser(stream)
	lngParser.RemoveErrorListeners()
	lngParser.AddErrorListener(errorListener)
	listener := &formulaListener{fac: p.fac}
	antlr.ParseTreeWalkerDefault.Walk(listener, lngParser.Formula())
	if errorListener.err {
		return p.fac.Falsum(), errorx.BadInput(errorListener.message)
	}
	return listener.formula(), nil
}

type errorListener struct {
	*antlr.DefaultErrorListener
	message string
	err     bool
}

func (l *errorListener) SyntaxError(
	_ antlr.Recognizer,
	_ interface{},
	line, column int,
	msg string,
	_ antlr.RecognitionException,
) {
	l.err = true
	l.message = fmt.Sprintf("Syntax error at line %d, column %d: %s\n", line, column, msg)
}

type stackFormula struct {
	formula f.Formula
	divider bool
}

type formulaListener struct {
	*BaseLogicNGPropositionalListener
	fac   f.Factory
	stack []stackFormula
	pbc   *pbc
}

type pbc struct {
	literals     []f.Literal
	coefficients []int
	rhs          int
	comparator   f.CSort
}

func (l *formulaListener) formula() f.Formula {
	return l.stack[0].formula
}

func (l *formulaListener) pushDivider() {
	l.stack = append(l.stack, stackFormula{0, true})
}

func (l *formulaListener) pushFormula(f f.Formula) {
	l.stack = append(l.stack, stackFormula{f, false})
}

func (l *formulaListener) pop() stackFormula {
	if len(l.stack) < 1 {
		panic(errorx.IllegalState("stack is empty unable to pop"))
	}
	result := l.stack[len(l.stack)-1]
	l.stack = l.stack[:len(l.stack)-1]
	return result
}

func (l *formulaListener) ExitSimp(c *SimpContext) {
	if c.TRUE() != nil {
		l.pushFormula(l.fac.Verum())
	} else if c.FALSE() != nil {
		l.pushFormula(l.fac.Falsum())
	} else if c.LITERAL() != nil || c.NUMBER() != nil {
		name, phase := parseLiteral(c.GetText())
		l.pushFormula(l.fac.Literal(name, phase))
	}
}

func parseLiteral(text string) (name string, phase bool) {
	if strings.HasPrefix(text, "~") {
		name = text[1:]
		phase = false
	} else {
		name = text
		phase = true
	}
	return
}

func (l *formulaListener) ExitLit(c *LitContext) {
	if c.NOT() != nil {
		negated := l.fac.Not(l.pop().formula)
		l.pushFormula(negated)
	}
}

func (l *formulaListener) EnterConj(c *ConjContext) {
	if c.AND(0) != nil {
		l.pushDivider()
	}
}

func (l *formulaListener) ExitConj(c *ConjContext) {
	if c.AND(0) != nil {
		operands := make([]f.Formula, 0)
		current := l.pop()
		for !current.divider {
			operands = append(operands, current.formula)
			current = l.pop()
		}
		for i, j := 0, len(operands)-1; i < j; i, j = i+1, j-1 {
			operands[i], operands[j] = operands[j], operands[i]
		}
		l.pushFormula(l.fac.And(operands...))
	}
}

func (l *formulaListener) EnterDisj(c *DisjContext) {
	if c.OR(0) != nil {
		l.pushDivider()
	}
}

func (l *formulaListener) ExitDisj(c *DisjContext) {
	if c.OR(0) != nil {
		operands := make([]f.Formula, 0)
		current := l.pop()
		for !current.divider {
			operands = append(operands, current.formula)
			current = l.pop()
		}
		for i, j := 0, len(operands)-1; i < j; i, j = i+1, j-1 {
			operands[i], operands[j] = operands[j], operands[i]
		}
		l.pushFormula(l.fac.Or(operands...))
	}
}

func (l *formulaListener) ExitImpl(c *ImplContext) {
	if c.IMPL() != nil {
		right := l.pop().formula
		left := l.pop().formula
		l.pushFormula(l.fac.Implication(left, right))
	}
}

func (l *formulaListener) ExitEquiv(c *EquivContext) {
	if c.EQUIV() != nil {
		right := l.pop().formula
		left := l.pop().formula
		l.pushFormula(l.fac.Equivalence(left, right))
	}
}

func (l *formulaListener) EnterAdd(_ *AddContext) {
	l.pbc.literals = make([]f.Literal, 0)
	l.pbc.coefficients = make([]int, 0)
}

func (l *formulaListener) ExitMul(c *MulContext) {
	coeff := 1
	var lit string
	switch {
	case c.LITERAL() != nil && len(c.AllNUMBER()) == 0:
		lit = c.LITERAL().GetText()
	case c.LITERAL() == nil && len(c.AllNUMBER()) == 1:
		lit = c.LITERAL().GetText()
	case c.LITERAL() != nil && len(c.AllNUMBER()) == 1:
		lit = c.LITERAL().GetText()
		coeff, _ = strconv.Atoi(c.NUMBER(0).GetText())
	case c.LITERAL() == nil && len(c.AllNUMBER()) == 2:
		lit = c.NUMBER(0).GetText()
		coeff, _ = strconv.Atoi(c.NUMBER(1).GetText())
	}
	name, phase := parseLiteral(lit)
	literal := l.fac.Lit(name, phase)
	l.pbc.literals = append(l.pbc.literals, literal)
	l.pbc.coefficients = append(l.pbc.coefficients, coeff)
}

func (l *formulaListener) EnterComparison(_ *ComparisonContext) {
	l.pbc = &pbc{}
}

func (l *formulaListener) ExitComparison(c *ComparisonContext) {
	l.pbc.rhs, _ = strconv.Atoi(c.NUMBER().GetText())
	l.pbc.comparator = l.parseComparator(c)
	l.pushFormula(l.fac.PBC(l.pbc.comparator, l.pbc.rhs, l.pbc.literals, l.pbc.coefficients))
}

func (l *formulaListener) parseComparator(c *ComparisonContext) f.CSort {
	switch {
	case c.EQ() != nil:
		return f.EQ
	case c.LE() != nil:
		return f.LE
	case c.LT() != nil:
		return f.LT
	case c.GE() != nil:
		return f.GE
	case c.GT() != nil:
		return f.GT
	default:
		panic(errorx.IllegalState("unknown comparator in parser"))
	}
}
