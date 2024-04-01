package assignment

import (
	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
)

// EvaluatesToTrue reports whether the formula evaluates to true under the
// assignment (mapping from variables to truth values).
//
// If a partial assignment is given, the check only recognizes simple
// unsatisfiable/tautology cases
//   - operand of an AND/OR is false/true
//   - all operators of an OR/AND are false/true
//   - AND/OR has two operands with complementary negations
//
// This evaluation differs from the standard [Evaluate] in two ways:
//  1. It accepts partial assignments
//  2. It tries to avoid the generation of intermediate formulas to speed up the performance
func EvaluatesToTrue(fac f.Factory, formula f.Formula, assignment map[f.Variable]bool) bool {
	context := evaluationContext{fac, true, assignment}
	return context.test(formula, true).Sort() == f.SortTrue
}

// EvaluatesToFalse reports whether the formula evaluates to false under the
// assignment (mapping from variables to truth values).
//
// If a partial assignment is given, the check only recognizes simple
// unsatisfiable/tautology cases
//   - operand of an AND/OR is false/true
//   - all operators of an OR/AND are false/true
//   - AND/OR has two operands with complementary negations
//
// This evaluation differs from the standard [Evaluate] in two ways:
//  1. It accepts partial assignments
//  2. It tries to avoid the generation of intermediate formulas to speed up the performance
func EvaluatesToFalse(fac f.Factory, formula f.Formula, assignment map[f.Variable]bool) bool {
	context := evaluationContext{fac, false, assignment}
	return context.test(formula, true).Sort() == f.SortFalse
}

func (c *evaluationContext) test(formula f.Formula, topLevel bool) f.Formula {
	switch formula.Sort() {
	case f.SortFalse, f.SortTrue:
		return formula
	case f.SortLiteral:
		variable := f.Literal(formula).Variable()
		found, ok := c.mapping[variable]
		if !ok {
			return formula
		} else {
			return c.fac.Constant(formula.IsPos() == found)
		}
	case f.SortNot:
		return c.handleNot(formula, topLevel)
	case f.SortImpl:
		return c.handleImplication(formula, topLevel)
	case f.SortEquiv:
		return c.handleEquivalence(formula, topLevel)
	case f.SortOr:
		return c.handleOr(formula, topLevel)
	case f.SortAnd:
		return c.handleAnd(formula, topLevel)
	case f.SortCC, f.SortPBC:
		return c.handlePbc(formula)
	default:
		panic(errorx.UnknownEnumValue(formula.Sort()))
	}
}

func (c *evaluationContext) handleNot(formula f.Formula, topLevel bool) f.Formula {
	op, _ := c.fac.NotOperand(formula)
	opResult := c.test(op, false)
	if topLevel && !opResult.IsConstant() {
		return c.fac.Constant(!c.evaluatesToTrue)
	}
	if opResult.IsConstant() {
		return c.fac.Constant(isFalsum(opResult))
	} else {
		return c.fac.Not(opResult)
	}
}

func (c *evaluationContext) handleImplication(formula f.Formula, topLevel bool) f.Formula {
	left, right, _ := c.fac.BinaryLeftRight(formula)
	leftResult := c.test(left, false)
	if leftResult.IsConstant() {
		if c.evaluatesToTrue {
			if isFalsum(leftResult) {
				return c.fac.Verum()
			} else {
				return c.test(right, topLevel)
			}
		} else {
			if isVerum(leftResult) {
				return c.test(right, topLevel)
			} else {
				return c.fac.Verum()
			}
		}
	}
	if !c.evaluatesToTrue && topLevel {
		return c.fac.Verum()
	}
	rightResult := c.test(right, false)
	if rightResult.IsConstant() {
		if isVerum(rightResult) {
			return c.fac.Verum()
		} else {
			return c.fac.Not(leftResult)
		}
	}
	return c.fac.Implication(leftResult, rightResult)
}

func (c *evaluationContext) handleEquivalence(formula f.Formula, topLevel bool) f.Formula {
	left, right, _ := c.fac.BinaryLeftRight(formula)
	leftResult := c.test(left, false)
	if leftResult.IsConstant() {
		if isVerum(leftResult) {
			return c.test(right, topLevel)
		} else {
			return c.test(c.fac.Not(right), topLevel)
		}
	}

	rightResult := c.test(right, false)
	if rightResult.IsConstant() {
		if topLevel {
			return c.fac.Constant(!c.evaluatesToTrue)
		}
		if isVerum(rightResult) {
			return leftResult
		} else {
			return c.fac.Not(leftResult)
		}
	}

	return c.fac.Equivalence(leftResult, rightResult)
}

func (c *evaluationContext) handleOr(formula f.Formula, topLevel bool) f.Formula {
	ops, _ := c.fac.NaryOperands(formula)
	nops := make([]f.Formula, 0, len(ops))
	for _, op := range ops {
		opResult := c.test(op, !c.evaluatesToTrue && topLevel)
		if isVerum(opResult) {
			return c.fac.Verum()
		}
		if !opResult.IsConstant() {
			if !c.evaluatesToTrue && topLevel {
				return c.fac.Verum()
			}
			nops = append(nops, opResult)
		}
	}
	return c.fac.Or(nops...)
}

func (c *evaluationContext) handleAnd(formula f.Formula, topLevel bool) f.Formula {
	ops, _ := c.fac.NaryOperands(formula)
	nops := make([]f.Formula, 0, len(ops))
	for _, op := range ops {
		opResult := c.test(op, c.evaluatesToTrue && topLevel)
		if isFalsum(opResult) {
			return c.fac.Falsum()
		}
		if !opResult.IsConstant() {
			if c.evaluatesToTrue && topLevel {
				return c.fac.Falsum()
			}
			nops = append(nops, opResult)
		}
	}
	return c.fac.And(nops...)
}

func (c *evaluationContext) handlePbc(formula f.Formula) f.Formula {
	assignment, _ := New(c.fac)
	for variable, phase := range c.mapping {
		name, _ := c.fac.VarName(variable)
		_ = assignment.AddLit(c.fac, c.fac.Lit(name, phase))
	}
	return Restrict(c.fac, formula, assignment)
}

func isFalsum(formula f.Formula) bool {
	return formula.Sort() == f.SortFalse
}

func isVerum(formula f.Formula) bool {
	return formula.Sort() == f.SortTrue
}

type evaluationContext struct {
	fac             f.Factory
	evaluatesToTrue bool
	mapping         map[f.Variable]bool
}
