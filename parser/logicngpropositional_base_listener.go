// Code generated from LogicNGPropositional.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // LogicNGPropositional

import "github.com/antlr4-go/antlr/v4"

// BaseLogicNGPropositionalListener is a complete listener for a parse tree produced by LogicNGPropositionalParser.
type BaseLogicNGPropositionalListener struct{}

var _ LogicNGPropositionalListener = &BaseLogicNGPropositionalListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseLogicNGPropositionalListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseLogicNGPropositionalListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseLogicNGPropositionalListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseLogicNGPropositionalListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterFormula is called when production formula is entered.
func (s *BaseLogicNGPropositionalListener) EnterFormula(ctx *FormulaContext) {}

// ExitFormula is called when production formula is exited.
func (s *BaseLogicNGPropositionalListener) ExitFormula(ctx *FormulaContext) {}

// EnterComparison is called when production comparison is entered.
func (s *BaseLogicNGPropositionalListener) EnterComparison(ctx *ComparisonContext) {}

// ExitComparison is called when production comparison is exited.
func (s *BaseLogicNGPropositionalListener) ExitComparison(ctx *ComparisonContext) {}

// EnterSimp is called when production simp is entered.
func (s *BaseLogicNGPropositionalListener) EnterSimp(ctx *SimpContext) {}

// ExitSimp is called when production simp is exited.
func (s *BaseLogicNGPropositionalListener) ExitSimp(ctx *SimpContext) {}

// EnterLit is called when production lit is entered.
func (s *BaseLogicNGPropositionalListener) EnterLit(ctx *LitContext) {}

// ExitLit is called when production lit is exited.
func (s *BaseLogicNGPropositionalListener) ExitLit(ctx *LitContext) {}

// EnterConj is called when production conj is entered.
func (s *BaseLogicNGPropositionalListener) EnterConj(ctx *ConjContext) {}

// ExitConj is called when production conj is exited.
func (s *BaseLogicNGPropositionalListener) ExitConj(ctx *ConjContext) {}

// EnterDisj is called when production disj is entered.
func (s *BaseLogicNGPropositionalListener) EnterDisj(ctx *DisjContext) {}

// ExitDisj is called when production disj is exited.
func (s *BaseLogicNGPropositionalListener) ExitDisj(ctx *DisjContext) {}

// EnterImpl is called when production impl is entered.
func (s *BaseLogicNGPropositionalListener) EnterImpl(ctx *ImplContext) {}

// ExitImpl is called when production impl is exited.
func (s *BaseLogicNGPropositionalListener) ExitImpl(ctx *ImplContext) {}

// EnterEquiv is called when production equiv is entered.
func (s *BaseLogicNGPropositionalListener) EnterEquiv(ctx *EquivContext) {}

// ExitEquiv is called when production equiv is exited.
func (s *BaseLogicNGPropositionalListener) ExitEquiv(ctx *EquivContext) {}

// EnterMul is called when production mul is entered.
func (s *BaseLogicNGPropositionalListener) EnterMul(ctx *MulContext) {}

// ExitMul is called when production mul is exited.
func (s *BaseLogicNGPropositionalListener) ExitMul(ctx *MulContext) {}

// EnterAdd is called when production add is entered.
func (s *BaseLogicNGPropositionalListener) EnterAdd(ctx *AddContext) {}

// ExitAdd is called when production add is exited.
func (s *BaseLogicNGPropositionalListener) ExitAdd(ctx *AddContext) {}
