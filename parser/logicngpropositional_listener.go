// Code generated from LogicNGPropositional.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // LogicNGPropositional

import "github.com/antlr4-go/antlr/v4"

// LogicNGPropositionalListener is a complete listener for a parse tree produced by LogicNGPropositionalParser.
type LogicNGPropositionalListener interface {
	antlr.ParseTreeListener

	// EnterFormula is called when entering the formula production.
	EnterFormula(c *FormulaContext)

	// EnterComparison is called when entering the comparison production.
	EnterComparison(c *ComparisonContext)

	// EnterSimp is called when entering the simp production.
	EnterSimp(c *SimpContext)

	// EnterLit is called when entering the lit production.
	EnterLit(c *LitContext)

	// EnterConj is called when entering the conj production.
	EnterConj(c *ConjContext)

	// EnterDisj is called when entering the disj production.
	EnterDisj(c *DisjContext)

	// EnterImpl is called when entering the impl production.
	EnterImpl(c *ImplContext)

	// EnterEquiv is called when entering the equiv production.
	EnterEquiv(c *EquivContext)

	// EnterMul is called when entering the mul production.
	EnterMul(c *MulContext)

	// EnterAdd is called when entering the add production.
	EnterAdd(c *AddContext)

	// ExitFormula is called when exiting the formula production.
	ExitFormula(c *FormulaContext)

	// ExitComparison is called when exiting the comparison production.
	ExitComparison(c *ComparisonContext)

	// ExitSimp is called when exiting the simp production.
	ExitSimp(c *SimpContext)

	// ExitLit is called when exiting the lit production.
	ExitLit(c *LitContext)

	// ExitConj is called when exiting the conj production.
	ExitConj(c *ConjContext)

	// ExitDisj is called when exiting the disj production.
	ExitDisj(c *DisjContext)

	// ExitImpl is called when exiting the impl production.
	ExitImpl(c *ImplContext)

	// ExitEquiv is called when exiting the equiv production.
	ExitEquiv(c *EquivContext)

	// ExitMul is called when exiting the mul production.
	ExitMul(c *MulContext)

	// ExitAdd is called when exiting the add production.
	ExitAdd(c *AddContext)
}
