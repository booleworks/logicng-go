// Code generated from LogicNGPropositional.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type LogicNGPropositionalLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var LogicNGPropositionalLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func logicngpropositionallexerLexerInit() {
	staticData := &LogicNGPropositionalLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "", "", "'$true'", "'$false'", "'('", "')'", "'~'", "'&'", "'|'",
		"'=>'", "'<=>'", "'*'", "'+'", "'='", "'<='", "'<'", "'>='", "'>'",
	}
	staticData.SymbolicNames = []string{
		"", "NUMBER", "LITERAL", "TRUE", "FALSE", "LBR", "RBR", "NOT", "AND",
		"OR", "IMPL", "EQUIV", "MUL", "ADD", "EQ", "LE", "LT", "GE", "GT", "WS",
	}
	staticData.RuleNames = []string{
		"NUMBER", "LITERAL", "TRUE", "FALSE", "LBR", "RBR", "NOT", "AND", "OR",
		"IMPL", "EQUIV", "MUL", "ADD", "EQ", "LE", "LT", "GE", "GT", "WS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 19, 110, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 1, 0, 3, 0, 41, 8, 0,
		1, 0, 4, 0, 44, 8, 0, 11, 0, 12, 0, 45, 1, 1, 3, 1, 49, 8, 1, 1, 1, 1,
		1, 5, 1, 53, 8, 1, 10, 1, 12, 1, 56, 9, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2,
		1, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 5, 1, 5,
		1, 6, 1, 6, 1, 7, 1, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1,
		10, 1, 10, 1, 11, 1, 11, 1, 12, 1, 12, 1, 13, 1, 13, 1, 14, 1, 14, 1, 14,
		1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 17, 1, 17, 1, 18, 4, 18, 105, 8,
		18, 11, 18, 12, 18, 106, 1, 18, 1, 18, 0, 0, 19, 1, 1, 3, 2, 5, 3, 7, 4,
		9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14,
		29, 15, 31, 16, 33, 17, 35, 18, 37, 19, 1, 0, 6, 1, 0, 45, 45, 1, 0, 48,
		57, 1, 0, 126, 126, 5, 0, 35, 35, 48, 57, 64, 90, 95, 95, 97, 122, 5, 0,
		35, 35, 48, 57, 65, 90, 95, 95, 97, 122, 3, 0, 9, 10, 13, 13, 32, 32, 114,
		0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0,
		0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0,
		0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0,
		0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1,
		0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 35, 1, 0, 0, 0, 0, 37, 1, 0, 0, 0, 1, 40,
		1, 0, 0, 0, 3, 48, 1, 0, 0, 0, 5, 57, 1, 0, 0, 0, 7, 63, 1, 0, 0, 0, 9,
		70, 1, 0, 0, 0, 11, 72, 1, 0, 0, 0, 13, 74, 1, 0, 0, 0, 15, 76, 1, 0, 0,
		0, 17, 78, 1, 0, 0, 0, 19, 80, 1, 0, 0, 0, 21, 83, 1, 0, 0, 0, 23, 87,
		1, 0, 0, 0, 25, 89, 1, 0, 0, 0, 27, 91, 1, 0, 0, 0, 29, 93, 1, 0, 0, 0,
		31, 96, 1, 0, 0, 0, 33, 98, 1, 0, 0, 0, 35, 101, 1, 0, 0, 0, 37, 104, 1,
		0, 0, 0, 39, 41, 7, 0, 0, 0, 40, 39, 1, 0, 0, 0, 40, 41, 1, 0, 0, 0, 41,
		43, 1, 0, 0, 0, 42, 44, 7, 1, 0, 0, 43, 42, 1, 0, 0, 0, 44, 45, 1, 0, 0,
		0, 45, 43, 1, 0, 0, 0, 45, 46, 1, 0, 0, 0, 46, 2, 1, 0, 0, 0, 47, 49, 7,
		2, 0, 0, 48, 47, 1, 0, 0, 0, 48, 49, 1, 0, 0, 0, 49, 50, 1, 0, 0, 0, 50,
		54, 7, 3, 0, 0, 51, 53, 7, 4, 0, 0, 52, 51, 1, 0, 0, 0, 53, 56, 1, 0, 0,
		0, 54, 52, 1, 0, 0, 0, 54, 55, 1, 0, 0, 0, 55, 4, 1, 0, 0, 0, 56, 54, 1,
		0, 0, 0, 57, 58, 5, 36, 0, 0, 58, 59, 5, 116, 0, 0, 59, 60, 5, 114, 0,
		0, 60, 61, 5, 117, 0, 0, 61, 62, 5, 101, 0, 0, 62, 6, 1, 0, 0, 0, 63, 64,
		5, 36, 0, 0, 64, 65, 5, 102, 0, 0, 65, 66, 5, 97, 0, 0, 66, 67, 5, 108,
		0, 0, 67, 68, 5, 115, 0, 0, 68, 69, 5, 101, 0, 0, 69, 8, 1, 0, 0, 0, 70,
		71, 5, 40, 0, 0, 71, 10, 1, 0, 0, 0, 72, 73, 5, 41, 0, 0, 73, 12, 1, 0,
		0, 0, 74, 75, 5, 126, 0, 0, 75, 14, 1, 0, 0, 0, 76, 77, 5, 38, 0, 0, 77,
		16, 1, 0, 0, 0, 78, 79, 5, 124, 0, 0, 79, 18, 1, 0, 0, 0, 80, 81, 5, 61,
		0, 0, 81, 82, 5, 62, 0, 0, 82, 20, 1, 0, 0, 0, 83, 84, 5, 60, 0, 0, 84,
		85, 5, 61, 0, 0, 85, 86, 5, 62, 0, 0, 86, 22, 1, 0, 0, 0, 87, 88, 5, 42,
		0, 0, 88, 24, 1, 0, 0, 0, 89, 90, 5, 43, 0, 0, 90, 26, 1, 0, 0, 0, 91,
		92, 5, 61, 0, 0, 92, 28, 1, 0, 0, 0, 93, 94, 5, 60, 0, 0, 94, 95, 5, 61,
		0, 0, 95, 30, 1, 0, 0, 0, 96, 97, 5, 60, 0, 0, 97, 32, 1, 0, 0, 0, 98,
		99, 5, 62, 0, 0, 99, 100, 5, 61, 0, 0, 100, 34, 1, 0, 0, 0, 101, 102, 5,
		62, 0, 0, 102, 36, 1, 0, 0, 0, 103, 105, 7, 5, 0, 0, 104, 103, 1, 0, 0,
		0, 105, 106, 1, 0, 0, 0, 106, 104, 1, 0, 0, 0, 106, 107, 1, 0, 0, 0, 107,
		108, 1, 0, 0, 0, 108, 109, 6, 18, 0, 0, 109, 38, 1, 0, 0, 0, 6, 0, 40,
		45, 48, 54, 106, 1, 6, 0, 0,
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

// LogicNGPropositionalLexerInit initializes any static state used to implement LogicNGPropositionalLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewLogicNGPropositionalLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func LogicNGPropositionalLexerInit() {
	staticData := &LogicNGPropositionalLexerLexerStaticData
	staticData.once.Do(logicngpropositionallexerLexerInit)
}

// NewLogicNGPropositionalLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewLogicNGPropositionalLexer(input antlr.CharStream) *LogicNGPropositionalLexer {
	LogicNGPropositionalLexerInit()
	l := new(LogicNGPropositionalLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &LogicNGPropositionalLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "LogicNGPropositional.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// LogicNGPropositionalLexer tokens.
const (
	LogicNGPropositionalLexerNUMBER  = 1
	LogicNGPropositionalLexerLITERAL = 2
	LogicNGPropositionalLexerTRUE    = 3
	LogicNGPropositionalLexerFALSE   = 4
	LogicNGPropositionalLexerLBR     = 5
	LogicNGPropositionalLexerRBR     = 6
	LogicNGPropositionalLexerNOT     = 7
	LogicNGPropositionalLexerAND     = 8
	LogicNGPropositionalLexerOR      = 9
	LogicNGPropositionalLexerIMPL    = 10
	LogicNGPropositionalLexerEQUIV   = 11
	LogicNGPropositionalLexerMUL     = 12
	LogicNGPropositionalLexerADD     = 13
	LogicNGPropositionalLexerEQ      = 14
	LogicNGPropositionalLexerLE      = 15
	LogicNGPropositionalLexerLT      = 16
	LogicNGPropositionalLexerGE      = 17
	LogicNGPropositionalLexerGT      = 18
	LogicNGPropositionalLexerWS      = 19
)
