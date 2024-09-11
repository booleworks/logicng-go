package simplification

import (
	"github.com/booleworks/logicng-go/assignment"
	"github.com/booleworks/logicng-go/configuration"
	"github.com/booleworks/logicng-go/event"
	"github.com/booleworks/logicng-go/explanation/smus"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/primeimplicant"
	"github.com/booleworks/logicng-go/sat"
)

// Config describes the configuration of an advanced simplifier.  It
// holds flags for activating different simplification steps as well as an
// optional rating function.
type Config struct {
	FactorOut         bool
	RestrictBackbone  bool
	SimplifyNegations bool
	UseRatingFunction bool
	RatingFunction    RatingFunction
}

// Sort returns the configuration sort (Advanced Simplifier).
func (Config) Sort() configuration.Sort {
	return configuration.AdvancedSimplifier
}

// DefaultConfig returns the default configuration for an advanced simplifier
// configuration.
func (Config) DefaultConfig() configuration.Config {
	return DefaultConfig()
}

// DefaultConfig returns the default configuration for an advanced simplifier
// configuration.
func DefaultConfig() *Config {
	return &Config{
		RestrictBackbone:  true,
		FactorOut:         true,
		SimplifyNegations: true,
		UseRatingFunction: true,
		RatingFunction:    DefaultRatingFunction,
	}
}

// QMC computes a simplification on the formula based on the
// algorithm by Quine and McCluskey.  This implementation is however not based
// on the traditional term tables but uses a SAT solver based implementation
// based on the advanced simplifier.  The resulting formula is in DNF.
func QMC(fac f.Factory, formula f.Formula) f.Formula {
	result, _ := QMCWithHandler(fac, formula, handler.NopHandler)
	return result
}

// QMCWithHandler computes a simplification on the formula based on
// the algorithm by Quine and McCluskey.  This implementation is however not
// based on the traditional term tables but uses a SAT solver based
// implementation based on the advanced simplifier.  The resulting formula is
// in DNF.  The given handler can be used to cancel the optimization function
// during the prime implicant computation.
func QMCWithHandler(fac f.Factory, formula f.Formula, hdl handler.Handler) (f.Formula, handler.State) {
	config := Config{
		RestrictBackbone:  false,
		FactorOut:         false,
		SimplifyNegations: false,
		UseRatingFunction: false,
		RatingFunction:    DefaultRatingFunction,
	}
	return AdvancedWithHandler(fac, formula, hdl, &config)
}

// Advanced simplifies the given formula by performing the following steps
//
//  1. Computation of all prime implicants
//  2. Finding the minimal coverage over the found prime implicants (by finding one smallest MUS)
//  3. Building a DNF from the minimal prime implicant coverage
//  4. Factoring out common factors of the DNF using the FactorOut function
//  5. Minimizing negations of the factored-out DNF using the SimplifyNegations function
//
// It can be configured with an optional advanced simplifier configuration.
func Advanced(fac f.Factory, formula f.Formula, config ...*Config) f.Formula {
	result, _ := AdvancedWithHandler(fac, formula, handler.NopHandler, config...)
	return result
}

// AdvancedWithHandler simplifies the given formula by performing the following steps
//
//  1. Computation of all prime implicants
//  2. Finding the minimal coverage over the found prime implicants (by finding one smallest MUS)
//  3. Building a DNF from the minimal prime implicant coverage
//  4. Factoring out common factors of the DNF using the FactorOut function
//  5. Minimizing negations of the factored-out DNF using the SimplifyNegations function
//
// It can be configured with an optional advanced simplifier configuration.
// The given  handler can be used to cancel the optimization function during the
// prime implicant computation.
func AdvancedWithHandler(
	fac f.Factory,
	formula f.Formula,
	hdl handler.Handler,
	config ...*Config,
) (f.Formula, handler.State) {
	cfg := determineConfig(fac, config)
	if e := event.AdvancedSimplificationStarted; !hdl.ShouldResume(e) {
		return 0, handler.Cancelation(e)
	}
	simplified := formula
	var backboneLiterals []f.Literal
	if cfg.RestrictBackbone {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		backbone, state := solver.ComputeBackboneWithHandler(fac, f.Variables(fac, formula).Content(), hdl)
		if !state.Success {
			return 0, state
		}
		if !backbone.Sat {
			return fac.Falsum(), handler.Success()
		}
		backboneLiterals = append(backboneLiterals, backbone.CompleteBackbone(fac)...)
		ass, _ := assignment.New(fac, backboneLiterals...)
		simplified = assignment.Restrict(fac, formula, ass)
	}
	simplifyMinDnf, state := computeMinDNF(fac, simplified, hdl)
	if !state.Success {
		return 0, state
	}
	simplified = simplifyWithRating(fac, simplified, simplifyMinDnf, cfg)
	if cfg.FactorOut {
		factoredOut := FactorOut(fac, simplified, cfg.RatingFunction)
		simplified = simplifyWithRating(fac, simplified, factoredOut, cfg)
	}
	if cfg.RestrictBackbone {
		simplified = fac.And(fac.Minterm(backboneLiterals...), simplified)
	}
	if cfg.SimplifyNegations {
		negationSimplified := MinimizeNegations(fac, simplified)
		simplified = simplifyWithRating(fac, simplified, negationSimplified, cfg)
	}
	return simplified, handler.Success()
}

func determineConfig(fac f.Factory, initConfig []*Config) *Config {
	if len(initConfig) > 0 {
		return initConfig[0]
	} else {
		configFromFactory, ok := fac.ConfigurationFor(configuration.AdvancedSimplifier)
		if !ok {
			return DefaultConfig()
		} else {
			return configFromFactory.(*Config)
		}
	}
}

func computeMinDNF(fac f.Factory, simplified f.Formula, hdl handler.Handler) (f.Formula, handler.State) {
	primeResult, state := primeimplicant.CoverMinWithHandler(
		fac, simplified, primeimplicant.CoverImplicants, hdl,
	)
	if !state.Success {
		return 0, state
	}
	primeImplicants := primeResult.Implicants
	minimizedPIs, state := smus.ComputeForFormulasWithHandler(
		fac,
		negateAllLiterals(fac, primeImplicants),
		hdl,
		simplified,
	)
	if !state.Success {
		return 0, state
	}
	simplified = fac.Or(negateAllLiteralsInFormulas(fac, minimizedPIs)...)
	return simplified, handler.Success()
}

func negateAllLiterals(fac f.Factory, literalSets [][]f.Literal) []f.Formula {
	result := make([]f.Formula, len(literalSets))
	for i, literals := range literalSets {
		negated := make([]f.Formula, len(literals))
		for j, literal := range literals {
			negated[j] = literal.Negate(fac).AsFormula()
		}
		result[i] = fac.Or(negated...)
	}
	return result
}

func negateAllLiteralsInFormulas(fac f.Factory, formulas []f.Formula) []f.Formula {
	result := make([]f.Formula, len(formulas))
	for i, formula := range formulas {
		literals := f.Literals(fac, formula).Content()
		negated := make([]f.Literal, len(literals))
		for j, literal := range literals {
			negated[j] = literal.Negate(fac)
		}
		result[i] = fac.Minterm(negated...)
	}
	return result
}

func simplifyWithRating(fac f.Factory, formula, simplifiedOneStep f.Formula, config *Config) f.Formula {
	if !config.UseRatingFunction {
		return simplifiedOneStep
	}
	ratingSimplified := config.RatingFunction(fac, simplifiedOneStep)
	ratingFormula := config.RatingFunction(fac, formula)
	if ratingSimplified < ratingFormula {
		return simplifiedOneStep
	} else {
		return formula
	}
}
