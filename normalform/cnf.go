package normalform

import (
	"booleworks.com/logicng/assignment"
	"booleworks.com/logicng/configuration"
	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/function"
	"booleworks.com/logicng/handler"
)

// CNFAlgorithm encodes the different algorithms for converting a formula to
// CNF.
type CNFAlgorithm byte

const (
	CNFFactorization CNFAlgorithm = iota
	CNFTseitin
	CNFPlaistedGreenbaum
	CNFAdvanced
)

//go:generate stringer -type=CNFAlgorithm

// CNFConfig represents the configuration for the CNF conversion.  It stores
// the main algorithm to transform the formula.  If the advanced algorithm is
// chosen, also a fallback algorithm has to be chosen.  Furthermore boundaries
// for number of distributions, created clauses, or atoms can be provided for
// the advanced algorithm.
type CNFConfig struct {
	Algorithm                            CNFAlgorithm
	FallbackAlgorithmForAdvancedEncoding CNFAlgorithm
	DistributionBoundary                 int
	CreatedClauseBoundary                int
	AtomBoundary                         int
}

// Sort returns the configuration sort (CNF).
func (CNFConfig) Sort() configuration.Sort {
	return configuration.CNF
}

// DefaultConfig returns the default configuration for a CNF configuration.
func (CNFConfig) DefaultConfig() configuration.Config {
	return DefaultCNFConfig()
}

// DefaultCNFConfig returns the default configuration for a CNF configuration.
func DefaultCNFConfig() *CNFConfig {
	return &CNFConfig{
		Algorithm:                            CNFAdvanced,
		FallbackAlgorithmForAdvancedEncoding: CNFTseitin,
		DistributionBoundary:                 -1,
		CreatedClauseBoundary:                1000,
		AtomBoundary:                         12,
	}
}

// IsCNF reports whether the given formula is in conjunctive normal form.  A
// CNF is a conjunction of discjuntion of literals.
func IsCNF(fac f.Factory, formula f.Formula) bool {
	cached, ok := f.LookupPredicateCache(fac, f.PredCNF, formula)
	if ok {
		return cached
	}
	switch fsort := formula.Sort(); fsort {
	case f.SortFalse, f.SortTrue, f.SortLiteral:
		return true
	case f.SortNot, f.SortImpl, f.SortEquiv, f.SortCC, f.SortPBC:
		return false
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
}

// CNF returns a conjunctive normal form of the given formula with the optional
// CNF configuration.  If no configuration is provided, the advanced CNF
// algorithm is used.  This algorithm encodes the given formula to CNF by first
// trying to use factorization for the single sub-formulas. When certain
// user-provided boundaries are met, the method is switched to Tseitin or
// Plaisted & Greenbaum as a fallback.
func CNF(fac f.Factory, formula f.Formula, config ...*CNFConfig) f.Formula {
	cfg := determineConfig(fac, config)
	switch cfg.Algorithm {
	case CNFFactorization:
		return FactorizedCNF(fac, formula)
	case CNFPlaistedGreenbaum:
		return PGCNFWithBoundary(fac, formula, cfg.AtomBoundary)
	case CNFTseitin:
		return TseitinCNFWithBoundary(fac, formula, cfg.AtomBoundary)
	case CNFAdvanced:
		return advancedEncoding(fac, formula, cfg)
	default:
		panic(errorx.UnknownEnumValue(cfg.Algorithm))
	}
}

func determineConfig(fac f.Factory, initConfig []*CNFConfig) *CNFConfig {
	if len(initConfig) > 0 {
		return initConfig[0]
	} else {
		configFromFactory, ok := fac.ConfigurationFor(configuration.CNF)
		if !ok {
			return DefaultCNFConfig()
		} else {
			return configFromFactory.(*CNFConfig)
		}
	}
}

func advancedEncoding(fac f.Factory, formula f.Formula, config *CNFConfig) f.Formula {
	factorizationHandler := &CNFHandler{
		distributionBoundary: config.DistributionBoundary,
		clauseBoundary:       config.CreatedClauseBoundary,
	}
	if formula.Sort() == f.SortAnd {
		ops, _ := fac.NaryOperands(formula)
		newOps := make([]f.Formula, len(ops))
		for i, op := range ops {
			newOps[i] = singleAdvancedEncoding(fac, op, config, factorizationHandler)
		}
		return fac.And(newOps...)
	}
	return singleAdvancedEncoding(fac, formula, config, factorizationHandler)
}

func singleAdvancedEncoding(
	fac f.Factory, formula f.Formula, config *CNFConfig, cnfHandler *CNFHandler,
) f.Formula {
	result, ok := FactorizedCNFWithHandler(fac, formula, cnfHandler)
	if !ok {
		if config.FallbackAlgorithmForAdvancedEncoding == CNFPlaistedGreenbaum {
			return PGCNFWithBoundary(fac, formula, config.AtomBoundary)
		} else {
			return TseitinCNFWithBoundary(fac, formula, config.AtomBoundary)
		}
	}
	return result
}

// FactorizedCNF returns the given formula in conjunctive normal form.  The
// algorithm used is factorization.  The resulting CNF can grow exponentially,
// therefore unless you are sure that the input is sensible, prefer the CNF
// factorization with a handler in order to be able to abort it.
func FactorizedCNF(fac f.Factory, formula f.Formula) f.Formula {
	cnf, _ := factorizedCNFRec(fac, formula, nil)
	return cnf
}

// FactorizedCNFWithHandler returns the given formula in conjunctive normal
// form.  The given handler can be used to abort the factorization.  Returns
// the CNF and an ok flag which is false when the handler aborted the
// computation.
func FactorizedCNFWithHandler(
	fac f.Factory, formula f.Formula, factorizatonHandler FactorizationHandler,
) (cnf f.Formula, ok bool) {
	handler.Start(factorizatonHandler)
	return factorizedCNFRec(fac, formula, factorizatonHandler)
}

func factorizedCNFRec(fac f.Factory, formula f.Formula, handler FactorizationHandler) (f.Formula, bool) {
	if formula.Sort() <= f.SortLiteral {
		return formula, true
	}
	cached, ok := f.LookupTransformationCache(fac, f.TransCNFFactorization, formula)
	if ok {
		return cached, true
	}
	ok = true
	switch fsort := formula.Sort(); fsort {
	case f.SortNot, f.SortImpl, f.SortEquiv:
		cached, ok = factorizedCNFRec(fac, NNF(fac, formula), handler)
	case f.SortOr:
		nary, _ := fac.NaryOperands(formula)
		nops := make([]f.Formula, 0, len(nary))
		for _, op := range nary {
			if !ok {
				return 0, false
			}
			var nop f.Formula
			nop, ok = factorizedCNFRec(fac, op, handler)
			nops = append(nops, nop)
		}
		cached = nops[0]
		for i := 1; i < len(nops); i++ {
			if !ok {
				return 0, false
			}
			cached, ok = distributeCNF(fac, cached, nops[i], handler)
		}
	case f.SortAnd:
		nary, _ := fac.NaryOperands(formula)
		nops := make([]f.Formula, 0, len(nary))
		for _, op := range nary {
			var apply f.Formula
			apply, ok = factorizedCNFRec(fac, op, handler)
			if !ok {
				return 0, false
			}
			nops = append(nops, apply)
		}
		cached = fac.And(nops...)
	case f.SortCC, f.SortPBC:
		cached = NNF(fac, formula)
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
	if ok {
		f.SetTransformationCache(fac, f.TransCNFFactorization, formula, cached)
		return cached, true
	}
	return 0, false
}

func distributeCNF(fac f.Factory, f1, f2 f.Formula, handler FactorizationHandler) (f.Formula, bool) {
	proceed := true
	if handler != nil {
		proceed = handler.PerformedDistribution()
	}
	if !proceed {
		return 0, false
	}
	if f1.Sort() == f.SortAnd || f2.Sort() == f.SortAnd {
		nops := make([]f.Formula, 0)
		var operands []f.Formula
		var form f.Formula
		if f1.Sort() == f.SortAnd {
			form = f2
			operands, _ = fac.NaryOperands(f1)
		} else {
			form = f1
			operands, _ = fac.NaryOperands(f2)
		}
		for _, op := range operands {
			distribute, ok := distributeCNF(fac, op, form, handler)
			if !ok {
				return 0, false
			}
			nops = append(nops, distribute)
		}
		return fac.And(nops...), true
	}
	clause := fac.Or(f1, f2)
	if handler != nil {
		proceed = handler.CreatedClause(clause)
	}
	return clause, proceed
}

// PGCNFDefault transforms a formula to CNF with the algorithm by Plaisted &
// Greenbaum.  Depending on the polarity of a sub-formula of the formula's
// parse tree it is replaced by a new auxiliary variable and an implication
// between variable and sub-formula is added.  In this default variant of the
// transformation sub-formulas with less than 12 atoms are transformed by
// factorization and auxiliary variables are only reused within one execution
// of the algorithm
func PGCNFDefault(fac f.Factory, formula f.Formula) f.Formula {
	return PGCNF(fac, formula, 12, NewCNFAuxState())
}

// PGCNFWithBoundary transforms a formula to CNF with the algorithm by Plaisted
// & Greenbaum.  Depending on the polarity of a sub-formula of the formula's
// parse tree it is replaced by a new auxiliary variable and an implication
// between variable and sub-formula is added.  The given
// boundaryForFactorization determines up to which number of atoms sub-formulas
// are transformed by factorization instead of the PG algorithm. Auxiliary
// variables are only reused within one execution of the algorithm.
func PGCNFWithBoundary(fac f.Factory, formula f.Formula, boundaryForFactorization int) f.Formula {
	return PGCNF(fac, formula, boundaryForFactorization, NewCNFAuxState())
}

// PGCNF transforms a formula to CNF with the algorithm by Plaisted
// & Greenbaum.  Depending on the polarity of a sub-formula of the formula's
// parse tree it is replaced by a new auxiliary variable and an implication
// between variable and sub-formula is added.  The given
// boundaryForFactorization determines up to which number of atoms sub-formulas
// are transformed by factorization instead of the PG algorithm.  The
// state which stores introduced auxiliary variables is provided by the caller
// and can therefore be reused between different executions of the method.
func PGCNF(fac f.Factory, formula f.Formula, boundaryForFactorization int, state *CNFAuxState) f.Formula {
	nnf := NNF(fac, formula)
	if IsCNF(fac, nnf) {
		return nnf
	}
	var pg f.Formula
	if function.NumberOfAtoms(fac, nnf) < boundaryForFactorization {
		pg = FactorizedCNF(fac, nnf)
	} else {
		pg = computeTransformation(fac, nnf, state)
		topLevel, _ := assignment.New(fac, state.LiteralMap[nnf])
		pg = assignment.Restrict(fac, pg, topLevel)
	}
	state.LiteralMap[formula] = state.LiteralMap[nnf]
	return pg
}

func computeTransformation(fac f.Factory, formula f.Formula, state *CNFAuxState) f.Formula {
	switch fsort := formula.Sort(); fsort {
	case f.SortLiteral:
		return fac.Verum()
	case f.SortOr, f.SortAnd:
		nary, _ := fac.NaryOperands(formula)
		nops := make([]f.Formula, 0, len(nary))
		nops = append(nops, computePosPolarity(fac, formula, state))
		for _, op := range nary {
			nops = append(nops, computeTransformation(fac, op, state))
		}
		return fac.And(nops...)
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
}

func computePosPolarity(fac f.Factory, formula f.Formula, state *CNFAuxState) f.Formula {
	result, ok := state.FormulaMap[formula]
	if ok {
		return result
	}
	pgVar := pgVariable(fac, formula, state)
	switch fsort := formula.Sort(); fsort {
	case f.SortAnd:
		nary, _ := fac.NaryOperands(formula)
		nops := make([]f.Formula, 0, len(nary))
		for _, op := range nary {
			nops = append(nops, fac.Clause(pgVar.Negate(fac), pgVariable(fac, op, state)))
		}
		result = fac.And(nops...)
	case f.SortOr:
		nary, _ := fac.NaryOperands(formula)
		nops := make([]f.Literal, 0, len(nary))
		nops = append(nops, pgVar.Negate(fac))
		for _, op := range nary {
			nops = append(nops, pgVariable(fac, op, state))
		}
		result = fac.Clause(nops...)
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
	state.FormulaMap[formula] = result
	return result
}

func pgVariable(fac f.Factory, formula f.Formula, state *CNFAuxState) f.Literal {
	if formula.Sort() == f.SortLiteral {
		return f.Literal(formula)
	}
	variable, ok := state.LiteralMap[formula]
	if !ok {
		variable = fac.NewCNFVariable().AsLiteral()
		state.LiteralMap[formula] = variable
	}
	return variable
}

// TseitinCNFDefault transforms a formula to CNF with the algorithm by Tseitin.
// Each sub-formula of the formula's parse tree is replaced by a new auxiliary
// variable and an equivalence between the new variable and the sub-formula is
// added.   In this default variant of the transformation sub-formulas with less
// than 12 atoms are transformed by factorization and auxiliary variables are
// only reused within one execution of the algorithm
func TseitinCNFDefault(fac f.Factory, formula f.Formula) f.Formula {
	return TseitinCNF(fac, formula, 12, NewCNFAuxState())
}

// TseitinCNFWithBoundary transforms a formula to CNF with the algorithm by
// Tseitin. Each sub-formula of the formula's parse tree is replaced by a new
// auxiliary variable and an equivalence between the new variable and the
// sub-formula is added.  The given boundaryForFactorization determines up to
// which number of atoms sub-formulas are transformed by factorization instead
// of the Tseitin algorithm.  Auxiliary variables are only reused within one
// execution of the algorithm.
func TseitinCNFWithBoundary(fac f.Factory, formula f.Formula, boundaryForFactorization int) f.Formula {
	return TseitinCNF(fac, formula, boundaryForFactorization, NewCNFAuxState())
}

// TseitinCNF transforms a formula to CNF with the algorithm by Tseitin. Each
// sub-formula of the formula's parse tree is replaced by a new auxiliary
// variable and an equivalence between the new variable and the sub-formula is
// added.  The given boundaryForFactorization determines up to which number of
// atoms sub-formulas are transformed by factorization instead of the Tseitin
// algorithm.  The state which stores introduced auxiliary variables is
// provided by the caller and can therefore be reused between different
// executions of the method.
func TseitinCNF(fac f.Factory, formula f.Formula, boundaryForFactorization int, state *CNFAuxState) f.Formula {
	nnf := NNF(fac, formula)
	if IsCNF(fac, nnf) {
		return nnf
	}
	tseitin, ok := state.FormulaMap[nnf]
	if ok {
		topLevel, _ := assignment.New(fac, state.LiteralMap[nnf])
		return assignment.Restrict(fac, tseitin, topLevel)
	}
	if function.NumberOfAtoms(fac, nnf) < boundaryForFactorization {
		tseitin = FactorizedCNF(fac, nnf)
	} else {
		for _, op := range function.SubNodes(fac, nnf) {
			computeTseitin(fac, op, state)
		}
		topLevel, _ := assignment.New(fac, state.LiteralMap[nnf])
		tseitin = assignment.Restrict(fac, state.FormulaMap[nnf], topLevel)
	}
	state.LiteralMap[formula] = state.LiteralMap[nnf]
	return tseitin
}

func computeTseitin(fac f.Factory, formula f.Formula, state *CNFAuxState) {
	if _, ok := state.FormulaMap[formula]; ok {
		return
	}
	switch formula.Sort() {
	case f.SortLiteral:
		state.FormulaMap[formula] = formula
		state.LiteralMap[formula] = f.Literal(formula)
	case f.SortAnd, f.SortOr:
		isConjunction := formula.Sort() == f.SortAnd
		tsLiteral := fac.NewCNFVariable().AsLiteral()
		var nops []f.Formula
		naryOps := fac.Operands(formula)
		operands := make([]f.Formula, 0, len(naryOps))
		negOperands := make([]f.Formula, 0, len(naryOps))
		if isConjunction {
			negOperands = append(negOperands, tsLiteral.AsFormula())
			handleNary(fac, naryOps, &nops, &operands, &negOperands, state)
			for _, operand := range operands {
				nops = append(nops, fac.Or(fac.Not(tsLiteral.AsFormula()), operand))
			}
			nops = append(nops, fac.Or(negOperands...))
		} else {
			operands = append(operands, tsLiteral.Negate(fac).AsFormula())
			handleNary(fac, naryOps, &nops, &operands, &negOperands, state)
			for _, operand := range negOperands {
				nops = append(nops, fac.Or(tsLiteral.AsFormula(), operand))
			}
			nops = append(nops, fac.Or(operands...))
		}
		state.LiteralMap[formula] = tsLiteral
		state.FormulaMap[formula] = fac.And(nops...)
	default:
		panic(errorx.IllegalState("could not process the formula type %s", formula.Sort()))
	}
}

func handleNary(
	fac f.Factory, origOps []f.Formula, nops, operands, negOperands *[]f.Formula, state *CNFAuxState,
) {
	for _, op := range origOps {
		if op.Sort() != f.SortLiteral {
			computeTseitin(fac, op, state)
			*nops = append(*nops, state.FormulaMap[op])
		}
		*operands = append(*operands, state.LiteralMap[op].AsFormula())
		*negOperands = append(*negOperands, state.LiteralMap[op].Negate(fac).AsFormula())
	}
}

// CNFAuxState is holds the variable and formula mapping for a Tseitin or
// Plaisted & Greenbaum CNF transformation.  If you want to reuse generated CNF
// auxiliary variables you can re-use such a state between different CNF
// computations.
type CNFAuxState struct {
	FormulaMap map[f.Formula]f.Formula
	LiteralMap map[f.Formula]f.Literal
}

// NewCNFAuxState generates a new empty state for Tseitin or Plaisted &
// Greenbaum CNF transformations.
func NewCNFAuxState() *CNFAuxState {
	return &CNFAuxState{make(map[f.Formula]f.Formula), make(map[f.Formula]f.Literal)}
}
