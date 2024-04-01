// Code generated from LogicNGPropositional.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // LogicNGPropositional

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type LogicNGPropositionalParser struct {
	*antlr.BaseParser
}

var LogicNGPropositionalParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func logicngpropositionalParserInit() {
	staticData := &LogicNGPropositionalParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "'$true'", "'$false'", "'('", "')'", "'~'", "'&'", "'|'",
		"'=>'", "'<=>'", "'*'", "'+'", "'='", "'<='", "'<'", "'>='", "'>'",
	}
	staticData.SymbolicNames = []string{
		"", "NUMBER", "LITERAL", "TRUE", "FALSE", "LBR", "RBR", "NOT", "AND",
		"OR", "IMPL", "EQUIV", "MUL", "ADD", "EQ", "LE", "LT", "GE", "GT", "WS",
	}
	staticData.RuleNames = []string{
		"formula", "comparison", "simp", "lit", "conj", "disj", "impl", "equiv",
		"mul", "add",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 19, 107, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 1, 0, 1,
		0, 3, 0, 23, 8, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1,
		45, 8, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3, 2, 56,
		8, 2, 1, 3, 1, 3, 1, 3, 3, 3, 61, 8, 3, 1, 4, 1, 4, 1, 4, 5, 4, 66, 8,
		4, 10, 4, 12, 4, 69, 9, 4, 1, 5, 1, 5, 1, 5, 5, 5, 74, 8, 5, 10, 5, 12,
		5, 77, 9, 5, 1, 6, 1, 6, 1, 6, 3, 6, 82, 8, 6, 1, 7, 1, 7, 1, 7, 3, 7,
		87, 8, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 97, 8,
		8, 1, 9, 1, 9, 1, 9, 5, 9, 102, 8, 9, 10, 9, 12, 9, 105, 9, 9, 1, 9, 0,
		0, 10, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 0, 0, 115, 0, 22, 1, 0, 0, 0,
		2, 44, 1, 0, 0, 0, 4, 55, 1, 0, 0, 0, 6, 60, 1, 0, 0, 0, 8, 62, 1, 0, 0,
		0, 10, 70, 1, 0, 0, 0, 12, 78, 1, 0, 0, 0, 14, 83, 1, 0, 0, 0, 16, 96,
		1, 0, 0, 0, 18, 98, 1, 0, 0, 0, 20, 23, 5, 0, 0, 1, 21, 23, 3, 14, 7, 0,
		22, 20, 1, 0, 0, 0, 22, 21, 1, 0, 0, 0, 23, 1, 1, 0, 0, 0, 24, 25, 3, 18,
		9, 0, 25, 26, 5, 14, 0, 0, 26, 27, 5, 1, 0, 0, 27, 45, 1, 0, 0, 0, 28,
		29, 3, 18, 9, 0, 29, 30, 5, 15, 0, 0, 30, 31, 5, 1, 0, 0, 31, 45, 1, 0,
		0, 0, 32, 33, 3, 18, 9, 0, 33, 34, 5, 16, 0, 0, 34, 35, 5, 1, 0, 0, 35,
		45, 1, 0, 0, 0, 36, 37, 3, 18, 9, 0, 37, 38, 5, 17, 0, 0, 38, 39, 5, 1,
		0, 0, 39, 45, 1, 0, 0, 0, 40, 41, 3, 18, 9, 0, 41, 42, 5, 18, 0, 0, 42,
		43, 5, 1, 0, 0, 43, 45, 1, 0, 0, 0, 44, 24, 1, 0, 0, 0, 44, 28, 1, 0, 0,
		0, 44, 32, 1, 0, 0, 0, 44, 36, 1, 0, 0, 0, 44, 40, 1, 0, 0, 0, 45, 3, 1,
		0, 0, 0, 46, 56, 5, 2, 0, 0, 47, 56, 5, 1, 0, 0, 48, 49, 5, 5, 0, 0, 49,
		50, 3, 14, 7, 0, 50, 51, 5, 6, 0, 0, 51, 56, 1, 0, 0, 0, 52, 56, 3, 2,
		1, 0, 53, 56, 5, 3, 0, 0, 54, 56, 5, 4, 0, 0, 55, 46, 1, 0, 0, 0, 55, 47,
		1, 0, 0, 0, 55, 48, 1, 0, 0, 0, 55, 52, 1, 0, 0, 0, 55, 53, 1, 0, 0, 0,
		55, 54, 1, 0, 0, 0, 56, 5, 1, 0, 0, 0, 57, 61, 3, 4, 2, 0, 58, 59, 5, 7,
		0, 0, 59, 61, 3, 6, 3, 0, 60, 57, 1, 0, 0, 0, 60, 58, 1, 0, 0, 0, 61, 7,
		1, 0, 0, 0, 62, 67, 3, 6, 3, 0, 63, 64, 5, 8, 0, 0, 64, 66, 3, 6, 3, 0,
		65, 63, 1, 0, 0, 0, 66, 69, 1, 0, 0, 0, 67, 65, 1, 0, 0, 0, 67, 68, 1,
		0, 0, 0, 68, 9, 1, 0, 0, 0, 69, 67, 1, 0, 0, 0, 70, 75, 3, 8, 4, 0, 71,
		72, 5, 9, 0, 0, 72, 74, 3, 8, 4, 0, 73, 71, 1, 0, 0, 0, 74, 77, 1, 0, 0,
		0, 75, 73, 1, 0, 0, 0, 75, 76, 1, 0, 0, 0, 76, 11, 1, 0, 0, 0, 77, 75,
		1, 0, 0, 0, 78, 81, 3, 10, 5, 0, 79, 80, 5, 10, 0, 0, 80, 82, 3, 12, 6,
		0, 81, 79, 1, 0, 0, 0, 81, 82, 1, 0, 0, 0, 82, 13, 1, 0, 0, 0, 83, 86,
		3, 12, 6, 0, 84, 85, 5, 11, 0, 0, 85, 87, 3, 14, 7, 0, 86, 84, 1, 0, 0,
		0, 86, 87, 1, 0, 0, 0, 87, 15, 1, 0, 0, 0, 88, 97, 5, 2, 0, 0, 89, 97,
		5, 1, 0, 0, 90, 91, 5, 1, 0, 0, 91, 92, 5, 12, 0, 0, 92, 97, 5, 2, 0, 0,
		93, 94, 5, 1, 0, 0, 94, 95, 5, 12, 0, 0, 95, 97, 5, 1, 0, 0, 96, 88, 1,
		0, 0, 0, 96, 89, 1, 0, 0, 0, 96, 90, 1, 0, 0, 0, 96, 93, 1, 0, 0, 0, 97,
		17, 1, 0, 0, 0, 98, 103, 3, 16, 8, 0, 99, 100, 5, 13, 0, 0, 100, 102, 3,
		16, 8, 0, 101, 99, 1, 0, 0, 0, 102, 105, 1, 0, 0, 0, 103, 101, 1, 0, 0,
		0, 103, 104, 1, 0, 0, 0, 104, 19, 1, 0, 0, 0, 105, 103, 1, 0, 0, 0, 10,
		22, 44, 55, 60, 67, 75, 81, 86, 96, 103,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// LogicNGPropositionalParserInit initializes any static state used to implement LogicNGPropositionalParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewLogicNGPropositionalParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func LogicNGPropositionalParserInit() {
	staticData := &LogicNGPropositionalParserStaticData
	staticData.once.Do(logicngpropositionalParserInit)
}

// NewLogicNGPropositionalParser produces a new parser instance for the optional input antlr.TokenStream.
func NewLogicNGPropositionalParser(input antlr.TokenStream) *LogicNGPropositionalParser {
	LogicNGPropositionalParserInit()
	this := new(LogicNGPropositionalParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &LogicNGPropositionalParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "LogicNGPropositional.g4"

	return this
}

// LogicNGPropositionalParser tokens.
const (
	LogicNGPropositionalParserEOF     = antlr.TokenEOF
	LogicNGPropositionalParserNUMBER  = 1
	LogicNGPropositionalParserLITERAL = 2
	LogicNGPropositionalParserTRUE    = 3
	LogicNGPropositionalParserFALSE   = 4
	LogicNGPropositionalParserLBR     = 5
	LogicNGPropositionalParserRBR     = 6
	LogicNGPropositionalParserNOT     = 7
	LogicNGPropositionalParserAND     = 8
	LogicNGPropositionalParserOR      = 9
	LogicNGPropositionalParserIMPL    = 10
	LogicNGPropositionalParserEQUIV   = 11
	LogicNGPropositionalParserMUL     = 12
	LogicNGPropositionalParserADD     = 13
	LogicNGPropositionalParserEQ      = 14
	LogicNGPropositionalParserLE      = 15
	LogicNGPropositionalParserLT      = 16
	LogicNGPropositionalParserGE      = 17
	LogicNGPropositionalParserGT      = 18
	LogicNGPropositionalParserWS      = 19
)

// LogicNGPropositionalParser rules.
const (
	LogicNGPropositionalParserRULE_formula    = 0
	LogicNGPropositionalParserRULE_comparison = 1
	LogicNGPropositionalParserRULE_simp       = 2
	LogicNGPropositionalParserRULE_lit        = 3
	LogicNGPropositionalParserRULE_conj       = 4
	LogicNGPropositionalParserRULE_disj       = 5
	LogicNGPropositionalParserRULE_impl       = 6
	LogicNGPropositionalParserRULE_equiv      = 7
	LogicNGPropositionalParserRULE_mul        = 8
	LogicNGPropositionalParserRULE_add        = 9
)

// IFormulaContext is an interface to support dynamic dispatch.
type IFormulaContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EOF() antlr.TerminalNode
	Equiv() IEquivContext

	// IsFormulaContext differentiates from other interfaces.
	IsFormulaContext()
}

type FormulaContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFormulaContext() *FormulaContext {
	var p = new(FormulaContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_formula
	return p
}

func InitEmptyFormulaContext(p *FormulaContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_formula
}

func (*FormulaContext) IsFormulaContext() {}

func NewFormulaContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FormulaContext {
	var p = new(FormulaContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LogicNGPropositionalParserRULE_formula

	return p
}

func (s *FormulaContext) GetParser() antlr.Parser { return s.parser }

func (s *FormulaContext) EOF() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserEOF, 0)
}

func (s *FormulaContext) Equiv() IEquivContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEquivContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEquivContext)
}

func (s *FormulaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FormulaContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FormulaContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.EnterFormula(s)
	}
}

func (s *FormulaContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.ExitFormula(s)
	}
}

func (p *LogicNGPropositionalParser) Formula() (localctx IFormulaContext) {
	localctx = NewFormulaContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, LogicNGPropositionalParserRULE_formula)
	p.SetState(22)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case LogicNGPropositionalParserEOF:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(20)
			p.Match(LogicNGPropositionalParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case LogicNGPropositionalParserNUMBER, LogicNGPropositionalParserLITERAL, LogicNGPropositionalParserTRUE, LogicNGPropositionalParserFALSE, LogicNGPropositionalParserLBR, LogicNGPropositionalParserNOT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(21)
			p.Equiv()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IComparisonContext is an interface to support dynamic dispatch.
type IComparisonContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Add() IAddContext
	EQ() antlr.TerminalNode
	NUMBER() antlr.TerminalNode
	LE() antlr.TerminalNode
	LT() antlr.TerminalNode
	GE() antlr.TerminalNode
	GT() antlr.TerminalNode

	// IsComparisonContext differentiates from other interfaces.
	IsComparisonContext()
}

type ComparisonContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComparisonContext() *ComparisonContext {
	var p = new(ComparisonContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_comparison
	return p
}

func InitEmptyComparisonContext(p *ComparisonContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_comparison
}

func (*ComparisonContext) IsComparisonContext() {}

func NewComparisonContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparisonContext {
	var p = new(ComparisonContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LogicNGPropositionalParserRULE_comparison

	return p
}

func (s *ComparisonContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparisonContext) Add() IAddContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAddContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAddContext)
}

func (s *ComparisonContext) EQ() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserEQ, 0)
}

func (s *ComparisonContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserNUMBER, 0)
}

func (s *ComparisonContext) LE() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserLE, 0)
}

func (s *ComparisonContext) LT() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserLT, 0)
}

func (s *ComparisonContext) GE() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserGE, 0)
}

func (s *ComparisonContext) GT() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserGT, 0)
}

func (s *ComparisonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComparisonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.EnterComparison(s)
	}
}

func (s *ComparisonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.ExitComparison(s)
	}
}

func (p *LogicNGPropositionalParser) Comparison() (localctx IComparisonContext) {
	localctx = NewComparisonContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, LogicNGPropositionalParserRULE_comparison)
	p.SetState(44)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(24)
			p.Add()
		}
		{
			p.SetState(25)
			p.Match(LogicNGPropositionalParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(26)
			p.Match(LogicNGPropositionalParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(28)
			p.Add()
		}
		{
			p.SetState(29)
			p.Match(LogicNGPropositionalParserLE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(30)
			p.Match(LogicNGPropositionalParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(32)
			p.Add()
		}
		{
			p.SetState(33)
			p.Match(LogicNGPropositionalParserLT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(34)
			p.Match(LogicNGPropositionalParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(36)
			p.Add()
		}
		{
			p.SetState(37)
			p.Match(LogicNGPropositionalParserGE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(38)
			p.Match(LogicNGPropositionalParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(40)
			p.Add()
		}
		{
			p.SetState(41)
			p.Match(LogicNGPropositionalParserGT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(42)
			p.Match(LogicNGPropositionalParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISimpContext is an interface to support dynamic dispatch.
type ISimpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LITERAL() antlr.TerminalNode
	NUMBER() antlr.TerminalNode
	LBR() antlr.TerminalNode
	Equiv() IEquivContext
	RBR() antlr.TerminalNode
	Comparison() IComparisonContext
	TRUE() antlr.TerminalNode
	FALSE() antlr.TerminalNode

	// IsSimpContext differentiates from other interfaces.
	IsSimpContext()
}

type SimpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySimpContext() *SimpContext {
	var p = new(SimpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_simp
	return p
}

func InitEmptySimpContext(p *SimpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_simp
}

func (*SimpContext) IsSimpContext() {}

func NewSimpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SimpContext {
	var p = new(SimpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LogicNGPropositionalParserRULE_simp

	return p
}

func (s *SimpContext) GetParser() antlr.Parser { return s.parser }

func (s *SimpContext) LITERAL() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserLITERAL, 0)
}

func (s *SimpContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserNUMBER, 0)
}

func (s *SimpContext) LBR() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserLBR, 0)
}

func (s *SimpContext) Equiv() IEquivContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEquivContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEquivContext)
}

func (s *SimpContext) RBR() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserRBR, 0)
}

func (s *SimpContext) Comparison() IComparisonContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComparisonContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComparisonContext)
}

func (s *SimpContext) TRUE() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserTRUE, 0)
}

func (s *SimpContext) FALSE() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserFALSE, 0)
}

func (s *SimpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SimpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.EnterSimp(s)
	}
}

func (s *SimpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.ExitSimp(s)
	}
}

func (p *LogicNGPropositionalParser) Simp() (localctx ISimpContext) {
	localctx = NewSimpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, LogicNGPropositionalParserRULE_simp)
	p.SetState(55)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(46)
			p.Match(LogicNGPropositionalParserLITERAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(47)
			p.Match(LogicNGPropositionalParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(48)
			p.Match(LogicNGPropositionalParserLBR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(49)
			p.Equiv()
		}
		{
			p.SetState(50)
			p.Match(LogicNGPropositionalParserRBR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(52)
			p.Comparison()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(53)
			p.Match(LogicNGPropositionalParserTRUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(54)
			p.Match(LogicNGPropositionalParserFALSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILitContext is an interface to support dynamic dispatch.
type ILitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Simp() ISimpContext
	NOT() antlr.TerminalNode
	Lit() ILitContext

	// IsLitContext differentiates from other interfaces.
	IsLitContext()
}

type LitContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLitContext() *LitContext {
	var p = new(LitContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_lit
	return p
}

func InitEmptyLitContext(p *LitContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_lit
}

func (*LitContext) IsLitContext() {}

func NewLitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LitContext {
	var p = new(LitContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LogicNGPropositionalParserRULE_lit

	return p
}

func (s *LitContext) GetParser() antlr.Parser { return s.parser }

func (s *LitContext) Simp() ISimpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISimpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISimpContext)
}

func (s *LitContext) NOT() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserNOT, 0)
}

func (s *LitContext) Lit() ILitContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILitContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILitContext)
}

func (s *LitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.EnterLit(s)
	}
}

func (s *LitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.ExitLit(s)
	}
}

func (p *LogicNGPropositionalParser) Lit() (localctx ILitContext) {
	localctx = NewLitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, LogicNGPropositionalParserRULE_lit)
	p.SetState(60)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case LogicNGPropositionalParserNUMBER, LogicNGPropositionalParserLITERAL, LogicNGPropositionalParserTRUE, LogicNGPropositionalParserFALSE, LogicNGPropositionalParserLBR:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(57)
			p.Simp()
		}

	case LogicNGPropositionalParserNOT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(58)
			p.Match(LogicNGPropositionalParserNOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(59)
			p.Lit()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConjContext is an interface to support dynamic dispatch.
type IConjContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllLit() []ILitContext
	Lit(i int) ILitContext
	AllAND() []antlr.TerminalNode
	AND(i int) antlr.TerminalNode

	// IsConjContext differentiates from other interfaces.
	IsConjContext()
}

type ConjContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConjContext() *ConjContext {
	var p = new(ConjContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_conj
	return p
}

func InitEmptyConjContext(p *ConjContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_conj
}

func (*ConjContext) IsConjContext() {}

func NewConjContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConjContext {
	var p = new(ConjContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LogicNGPropositionalParserRULE_conj

	return p
}

func (s *ConjContext) GetParser() antlr.Parser { return s.parser }

func (s *ConjContext) AllLit() []ILitContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILitContext); ok {
			len++
		}
	}

	tst := make([]ILitContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILitContext); ok {
			tst[i] = t.(ILitContext)
			i++
		}
	}

	return tst
}

func (s *ConjContext) Lit(i int) ILitContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILitContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILitContext)
}

func (s *ConjContext) AllAND() []antlr.TerminalNode {
	return s.GetTokens(LogicNGPropositionalParserAND)
}

func (s *ConjContext) AND(i int) antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserAND, i)
}

func (s *ConjContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConjContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConjContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.EnterConj(s)
	}
}

func (s *ConjContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.ExitConj(s)
	}
}

func (p *LogicNGPropositionalParser) Conj() (localctx IConjContext) {
	localctx = NewConjContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, LogicNGPropositionalParserRULE_conj)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(62)
		p.Lit()
	}
	p.SetState(67)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == LogicNGPropositionalParserAND {
		{
			p.SetState(63)
			p.Match(LogicNGPropositionalParserAND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(64)
			p.Lit()
		}

		p.SetState(69)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDisjContext is an interface to support dynamic dispatch.
type IDisjContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllConj() []IConjContext
	Conj(i int) IConjContext
	AllOR() []antlr.TerminalNode
	OR(i int) antlr.TerminalNode

	// IsDisjContext differentiates from other interfaces.
	IsDisjContext()
}

type DisjContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDisjContext() *DisjContext {
	var p = new(DisjContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_disj
	return p
}

func InitEmptyDisjContext(p *DisjContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_disj
}

func (*DisjContext) IsDisjContext() {}

func NewDisjContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DisjContext {
	var p = new(DisjContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LogicNGPropositionalParserRULE_disj

	return p
}

func (s *DisjContext) GetParser() antlr.Parser { return s.parser }

func (s *DisjContext) AllConj() []IConjContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IConjContext); ok {
			len++
		}
	}

	tst := make([]IConjContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IConjContext); ok {
			tst[i] = t.(IConjContext)
			i++
		}
	}

	return tst
}

func (s *DisjContext) Conj(i int) IConjContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConjContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConjContext)
}

func (s *DisjContext) AllOR() []antlr.TerminalNode {
	return s.GetTokens(LogicNGPropositionalParserOR)
}

func (s *DisjContext) OR(i int) antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserOR, i)
}

func (s *DisjContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DisjContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DisjContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.EnterDisj(s)
	}
}

func (s *DisjContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.ExitDisj(s)
	}
}

func (p *LogicNGPropositionalParser) Disj() (localctx IDisjContext) {
	localctx = NewDisjContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, LogicNGPropositionalParserRULE_disj)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(70)
		p.Conj()
	}
	p.SetState(75)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == LogicNGPropositionalParserOR {
		{
			p.SetState(71)
			p.Match(LogicNGPropositionalParserOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(72)
			p.Conj()
		}

		p.SetState(77)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IImplContext is an interface to support dynamic dispatch.
type IImplContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Disj() IDisjContext
	IMPL() antlr.TerminalNode
	Impl() IImplContext

	// IsImplContext differentiates from other interfaces.
	IsImplContext()
}

type ImplContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImplContext() *ImplContext {
	var p = new(ImplContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_impl
	return p
}

func InitEmptyImplContext(p *ImplContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_impl
}

func (*ImplContext) IsImplContext() {}

func NewImplContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImplContext {
	var p = new(ImplContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LogicNGPropositionalParserRULE_impl

	return p
}

func (s *ImplContext) GetParser() antlr.Parser { return s.parser }

func (s *ImplContext) Disj() IDisjContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDisjContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDisjContext)
}

func (s *ImplContext) IMPL() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserIMPL, 0)
}

func (s *ImplContext) Impl() IImplContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImplContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImplContext)
}

func (s *ImplContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImplContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImplContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.EnterImpl(s)
	}
}

func (s *ImplContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.ExitImpl(s)
	}
}

func (p *LogicNGPropositionalParser) Impl() (localctx IImplContext) {
	localctx = NewImplContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, LogicNGPropositionalParserRULE_impl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(78)
		p.Disj()
	}
	p.SetState(81)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == LogicNGPropositionalParserIMPL {
		{
			p.SetState(79)
			p.Match(LogicNGPropositionalParserIMPL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(80)
			p.Impl()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEquivContext is an interface to support dynamic dispatch.
type IEquivContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Impl() IImplContext
	EQUIV() antlr.TerminalNode
	Equiv() IEquivContext

	// IsEquivContext differentiates from other interfaces.
	IsEquivContext()
}

type EquivContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEquivContext() *EquivContext {
	var p = new(EquivContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_equiv
	return p
}

func InitEmptyEquivContext(p *EquivContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_equiv
}

func (*EquivContext) IsEquivContext() {}

func NewEquivContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EquivContext {
	var p = new(EquivContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LogicNGPropositionalParserRULE_equiv

	return p
}

func (s *EquivContext) GetParser() antlr.Parser { return s.parser }

func (s *EquivContext) Impl() IImplContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImplContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImplContext)
}

func (s *EquivContext) EQUIV() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserEQUIV, 0)
}

func (s *EquivContext) Equiv() IEquivContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEquivContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEquivContext)
}

func (s *EquivContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EquivContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EquivContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.EnterEquiv(s)
	}
}

func (s *EquivContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.ExitEquiv(s)
	}
}

func (p *LogicNGPropositionalParser) Equiv() (localctx IEquivContext) {
	localctx = NewEquivContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, LogicNGPropositionalParserRULE_equiv)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(83)
		p.Impl()
	}
	p.SetState(86)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == LogicNGPropositionalParserEQUIV {
		{
			p.SetState(84)
			p.Match(LogicNGPropositionalParserEQUIV)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(85)
			p.Equiv()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMulContext is an interface to support dynamic dispatch.
type IMulContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LITERAL() antlr.TerminalNode
	AllNUMBER() []antlr.TerminalNode
	NUMBER(i int) antlr.TerminalNode
	MUL() antlr.TerminalNode

	// IsMulContext differentiates from other interfaces.
	IsMulContext()
}

type MulContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMulContext() *MulContext {
	var p = new(MulContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_mul
	return p
}

func InitEmptyMulContext(p *MulContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_mul
}

func (*MulContext) IsMulContext() {}

func NewMulContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MulContext {
	var p = new(MulContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LogicNGPropositionalParserRULE_mul

	return p
}

func (s *MulContext) GetParser() antlr.Parser { return s.parser }

func (s *MulContext) LITERAL() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserLITERAL, 0)
}

func (s *MulContext) AllNUMBER() []antlr.TerminalNode {
	return s.GetTokens(LogicNGPropositionalParserNUMBER)
}

func (s *MulContext) NUMBER(i int) antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserNUMBER, i)
}

func (s *MulContext) MUL() antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserMUL, 0)
}

func (s *MulContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MulContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MulContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.EnterMul(s)
	}
}

func (s *MulContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.ExitMul(s)
	}
}

func (p *LogicNGPropositionalParser) Mul() (localctx IMulContext) {
	localctx = NewMulContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, LogicNGPropositionalParserRULE_mul)
	p.SetState(96)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 8, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(88)
			p.Match(LogicNGPropositionalParserLITERAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(89)
			p.Match(LogicNGPropositionalParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(90)
			p.Match(LogicNGPropositionalParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(91)
			p.Match(LogicNGPropositionalParserMUL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(92)
			p.Match(LogicNGPropositionalParserLITERAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(93)
			p.Match(LogicNGPropositionalParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(94)
			p.Match(LogicNGPropositionalParserMUL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(95)
			p.Match(LogicNGPropositionalParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAddContext is an interface to support dynamic dispatch.
type IAddContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllMul() []IMulContext
	Mul(i int) IMulContext
	AllADD() []antlr.TerminalNode
	ADD(i int) antlr.TerminalNode

	// IsAddContext differentiates from other interfaces.
	IsAddContext()
}

type AddContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAddContext() *AddContext {
	var p = new(AddContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_add
	return p
}

func InitEmptyAddContext(p *AddContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LogicNGPropositionalParserRULE_add
}

func (*AddContext) IsAddContext() {}

func NewAddContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AddContext {
	var p = new(AddContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LogicNGPropositionalParserRULE_add

	return p
}

func (s *AddContext) GetParser() antlr.Parser { return s.parser }

func (s *AddContext) AllMul() []IMulContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IMulContext); ok {
			len++
		}
	}

	tst := make([]IMulContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IMulContext); ok {
			tst[i] = t.(IMulContext)
			i++
		}
	}

	return tst
}

func (s *AddContext) Mul(i int) IMulContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMulContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMulContext)
}

func (s *AddContext) AllADD() []antlr.TerminalNode {
	return s.GetTokens(LogicNGPropositionalParserADD)
}

func (s *AddContext) ADD(i int) antlr.TerminalNode {
	return s.GetToken(LogicNGPropositionalParserADD, i)
}

func (s *AddContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AddContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AddContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.EnterAdd(s)
	}
}

func (s *AddContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LogicNGPropositionalListener); ok {
		listenerT.ExitAdd(s)
	}
}

func (p *LogicNGPropositionalParser) Add() (localctx IAddContext) {
	localctx = NewAddContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, LogicNGPropositionalParserRULE_add)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(98)
		p.Mul()
	}
	p.SetState(103)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == LogicNGPropositionalParserADD {
		{
			p.SetState(99)
			p.Match(LogicNGPropositionalParserADD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(100)
			p.Mul()
		}

		p.SetState(105)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
