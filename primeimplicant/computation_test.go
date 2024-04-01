package primeimplicant

import (
	"fmt"
	"testing"
	"time"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/handler"
	"booleworks.com/logicng/io"
	"booleworks.com/logicng/parser"
	"booleworks.com/logicng/randomizer"
	"booleworks.com/logicng/sat"
	"github.com/stretchr/testify/assert"
)

func TestPrimeComputationDoc(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	f1 := p.ParseUnsafe("(A | B) & (A | C ) & (C | D) & (B | ~D)")
	primes := CoverMin(fac, f1, CoverImplicants)
	implicants := primes.Implicants
	assert.Equal(t, 3, len(implicants))
}

func TestPrimeComputationSimple(t *testing.T) {
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	computeAndVerify(t, fac, d.True)
	computeAndVerify(t, fac, d.False)
	computeAndVerify(t, fac, d.A)
	computeAndVerify(t, fac, d.NA)
	computeAndVerify(t, fac, d.AND1)
	computeAndVerify(t, fac, d.AND2)
	computeAndVerify(t, fac, d.AND3)
	computeAndVerify(t, fac, d.OR1)
	computeAndVerify(t, fac, d.OR2)
	computeAndVerify(t, fac, d.OR3)
	computeAndVerify(t, fac, d.NOT1)
	computeAndVerify(t, fac, d.NOT2)
	computeAndVerify(t, fac, d.IMP1)
	computeAndVerify(t, fac, d.IMP2)
	computeAndVerify(t, fac, d.IMP3)
	computeAndVerify(t, fac, d.IMP4)
	computeAndVerify(t, fac, d.EQ1)
	computeAndVerify(t, fac, d.EQ2)
	computeAndVerify(t, fac, d.EQ3)
	computeAndVerify(t, fac, d.EQ4)
	computeAndVerify(t, fac, d.PBC1)
	computeAndVerify(t, fac, d.PBC2)
	computeAndVerify(t, fac, d.PBC3)
	computeAndVerify(t, fac, d.PBC4)
	computeAndVerify(t, fac, d.PBC5)
}

func TestPrimeComputationCornerCases(t *testing.T) {
	fac := f.NewFactory()
	for _, formula := range f.NewCornerCases(fac) {
		computeAndVerify(t, fac, formula)
	}
}

func TestPrimeComputationRandom(t *testing.T) {
	fac := f.NewFactory()
	numTests := 100
	if testing.Short() {
		numTests = 10
	}
	for i := 0; i < numTests; i++ {
		config := randomizer.DefaultConfig()
		config.NumVars = 10
		config.WeightPBC = 0
		config.Seed = int64(i * 42)
		formula := randomizer.New(fac, config).Formula(4)
		computeAndVerify(t, fac, formula)
	}
}

func TestPrimeComputationOriginalFormulas(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	fac := f.NewFactory()
	formulas, _ := io.ReadFormulas(fac, "../test/data/formulas/simplify_formulas.txt")
	for _, formula := range formulas {
		resultImplicantsMin := CoverMin(fac, formula, CoverImplicants)
		verify(t, fac, resultImplicantsMin, formula)
		resultImplicatesMin := CoverMin(fac, formula, CoverImplicates)
		verify(t, fac, resultImplicatesMin, formula)
	}
}

func TestPrimeComputationHandlerSmall(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("a & b | ~c & a")
	duration, _ := time.ParseDuration("1s")
	handler := sat.OptimizationHandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))
	testHandlerImplicant(t, fac, handler, formula, false)
	testHandlerImplicate(t, fac, handler, formula, false)
}

func TestPrimeComputationHandlerLarge(t *testing.T) {
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/large.txt")
	duration, _ := time.ParseDuration("500ms")
	handler := sat.OptimizationHandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))
	testHandlerImplicant(t, fac, handler, formula, true)
	testHandlerImplicate(t, fac, handler, formula, true)
}

func computeAndVerify(t *testing.T, fac f.Factory, formula f.Formula) {
	assert := assert.New(t)
	resultImplicantsMax := CoverMax(fac, formula, CoverImplicants)
	verify(t, fac, resultImplicantsMax, formula)
	resultImplicantsMin := CoverMin(fac, formula, CoverImplicants)
	verify(t, fac, resultImplicantsMin, formula)
	assert.Equal(CoverImplicants, resultImplicantsMax.CoverSort)
	assert.Equal(CoverImplicants, resultImplicantsMin.CoverSort)

	resultImplicatesMax := CoverMax(fac, formula, CoverImplicates)
	verify(t, fac, resultImplicatesMax, formula)
	resultImplicatesMin := CoverMin(fac, formula, CoverImplicates)
	verify(t, fac, resultImplicatesMin, formula)
	assert.Equal(CoverImplicates, resultImplicatesMax.CoverSort)
	assert.Equal(CoverImplicates, resultImplicatesMin.CoverSort)
}

func verify(t *testing.T, fac f.Factory, result *PrimeResult, formula f.Formula) {
	verifyImplicants(t, fac, result.Implicants, formula)
	verifyImplicates(t, fac, result.Implicates, formula)
}

func verifyImplicants(t *testing.T, fac f.Factory, implicantSets [][]f.Literal, formula f.Formula) {
	implicants := make([]f.Formula, len(implicantSets))
	for i, implicant := range implicantSets {
		implicants[i] = fac.And(f.LiteralsAsFormulas(implicant)...)
		testPrimeImplicantProperty(t, fac, formula, implicant)
	}
	assert.True(t, sat.IsEquivalent(fac, fac.Or(implicants...), formula))
}

func verifyImplicates(t *testing.T, fac f.Factory, implicateSets [][]f.Literal, formula f.Formula) {
	implicates := make([]f.Formula, len(implicateSets))
	for i, implicate := range implicateSets {
		implicates[i] = fac.Or(f.LiteralsAsFormulas(implicate)...)
		testPrimeImplicateProperty(t, fac, formula, implicate)
	}
	assert.True(t, sat.IsEquivalent(fac, fac.And(implicates...), formula))
}

func testHandlerImplicant(t *testing.T, fac f.Factory, handler sat.OptimizationHandler, formula f.Formula, exp bool) {
	assert := assert.New(t)
	result, ok := CoverMaxWithHandler(fac, formula, CoverImplicants, handler)
	if exp {
		assert.True(handler.Aborted())
		assert.False(ok)
		assert.Nil(result)
	} else {
		fmt.Println(result)
		assert.False(handler.Aborted())
		assert.True(ok)
		assert.NotNil(result)
	}
	result, ok = CoverMinWithHandler(fac, formula, CoverImplicants, handler)
	if exp {
		assert.True(handler.Aborted())
		assert.False(ok)
		assert.Nil(result)
	} else {
		assert.False(handler.Aborted())
		assert.True(ok)
		assert.NotNil(result)
	}
}

func testHandlerImplicate(t *testing.T, fac f.Factory, handler sat.OptimizationHandler, formula f.Formula, exp bool) {
	assert := assert.New(t)
	result, ok := CoverMaxWithHandler(fac, formula, CoverImplicates, handler)
	if exp {
		assert.True(handler.Aborted())
		assert.False(ok)
		assert.Nil(result)
	} else {
		assert.False(handler.Aborted())
		assert.True(ok)
		assert.NotNil(result)
	}
	result, ok = CoverMinWithHandler(fac, formula, CoverImplicates, handler)
	if exp {
		assert.True(handler.Aborted())
		assert.False(ok)
		assert.Nil(result)
	} else {
		assert.False(handler.Aborted())
		assert.True(ok)
		assert.NotNil(result)
	}
}
