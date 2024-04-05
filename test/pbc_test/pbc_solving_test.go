package pbc_test

import (
	"fmt"
	"testing"

	"github.com/booleworks/logicng-go/model/enum"

	"github.com/booleworks/logicng-go/assignment"
	e "github.com/booleworks/logicng-go/encoding"
	f "github.com/booleworks/logicng-go/formula"
	s "github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

var pbcConfigs = []e.Config{
	{PBCEncoder: e.PBCSWC},
	{PBCEncoder: e.PBCAdderNetworks},
	{PBCEncoder: e.PBCBinaryMerge, BinaryMergeUseGAC: true, BinaryMergeNoSupportForSingleBit: true, BinaryMergeUseWatchDog: true},
	{PBCEncoder: e.PBCBinaryMerge, BinaryMergeUseGAC: true, BinaryMergeNoSupportForSingleBit: true, BinaryMergeUseWatchDog: false},
	{PBCEncoder: e.PBCBinaryMerge, BinaryMergeUseGAC: true, BinaryMergeNoSupportForSingleBit: false, BinaryMergeUseWatchDog: true},
	{PBCEncoder: e.PBCBinaryMerge, BinaryMergeUseGAC: true, BinaryMergeNoSupportForSingleBit: false, BinaryMergeUseWatchDog: false},
	{PBCEncoder: e.PBCBinaryMerge, BinaryMergeUseGAC: false, BinaryMergeNoSupportForSingleBit: true, BinaryMergeUseWatchDog: true},
	{PBCEncoder: e.PBCBinaryMerge, BinaryMergeUseGAC: false, BinaryMergeNoSupportForSingleBit: true, BinaryMergeUseWatchDog: false},
	{PBCEncoder: e.PBCBinaryMerge, BinaryMergeUseGAC: false, BinaryMergeNoSupportForSingleBit: false, BinaryMergeUseWatchDog: true},
	{PBCEncoder: e.PBCBinaryMerge, BinaryMergeUseGAC: false, BinaryMergeNoSupportForSingleBit: false, BinaryMergeUseWatchDog: false},
}

func getConfigs() []*s.Config {
	return []*s.Config{
		s.DefaultConfig().UseAtMost(false),
		s.DefaultConfig().UseAtMost(true),
	}
}

func TestPbEncodingLess(t *testing.T) {
	assert := assert.New(t)
	for _, config := range pbcConfigs {
		fac := f.NewFactory()
		for _, satConfig := range getConfigs() {
			coeffs10 := []int{3, 2, 2, 2, 2, 2, 2, 2, 2, 2}
			vars10 := make([]f.Variable, 10)
			for i := 0; i < 10; i++ {
				vars10[i] = fac.Var(fmt.Sprintf("v%d", i))
			}
			literals10 := f.VariablesAsLiterals(vars10)
			pbc := fac.PBC(f.LE, 6, literals10, coeffs10)
			encoding, err := e.EncodePBC(fac, pbc, &config)
			assert.Nil(err)
			solver := s.NewSolver(fac, satConfig)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models := enum.OnSolver(solver, vars10)
			assert.Len(models, 140)
			for _, model := range models {
				assert.True(len(model.PosVars()) <= 3)
			}

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.LT, 7, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models = enum.OnSolver(solver, vars10)
			assert.Len(models, 140)
			for _, model := range models {
				assert.True(len(model.PosVars()) <= 3)
			}

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.LE, 0, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models = enum.OnSolver(solver, vars10)
			assert.Len(models, 1)

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.LE, 1, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models = enum.OnSolver(solver, vars10)
			assert.Len(models, 1)

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.LT, 2, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models = enum.OnSolver(solver, vars10)
			assert.Len(models, 1)

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.LT, 1, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models = enum.OnSolver(solver, vars10)
			assert.Len(models, 1)
		}
	}
}

func TestPbEncodingGreater(t *testing.T) {
	assert := assert.New(t)
	for _, config := range pbcConfigs {
		fac := f.NewFactory()
		for _, satConfig := range getConfigs() {
			coeffs10 := []int{3, 2, 2, 2, 2, 2, 2, 2, 2, 2}
			vars10 := make([]f.Variable, 10)
			for i := 0; i < 10; i++ {
				vars10[i] = fac.Var(fmt.Sprintf("v%d", i))
			}
			literals10 := f.VariablesAsLiterals(vars10)
			pbc := fac.PBC(f.GE, 17, literals10, coeffs10)
			encoding, err := e.EncodePBC(fac, pbc, &config)
			assert.Nil(err)
			solver := s.NewSolver(fac, satConfig)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models := enum.OnSolver(solver, vars10)
			assert.Len(models, 47)
			for _, model := range models {
				assert.True(len(model.PosVars()) >= 8)
			}

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.GT, 16, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models = enum.OnSolver(solver, vars10)
			assert.Len(models, 47)
			for _, model := range models {
				assert.True(len(model.PosVars()) >= 8)
			}

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.GE, 21, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models = enum.OnSolver(solver, vars10)
			assert.Len(models, 1)

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.GE, 22, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.False(solver.Sat())

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.GT, 42, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.False(solver.Sat())
		}
	}
}

func TestPbEncodingEq(t *testing.T) {
	assert := assert.New(t)
	for _, config := range pbcConfigs {
		fac := f.NewFactory()
		for _, satConfig := range getConfigs() {
			coeffs10 := []int{3, 2, 2, 2, 2, 2, 2, 2, 2, 2}
			vars10 := make([]f.Variable, 10)
			for i := 0; i < 10; i++ {
				vars10[i] = fac.Var(fmt.Sprintf("v%d", i))
			}
			literals10 := f.VariablesAsLiterals(vars10)
			pbc := fac.PBC(f.EQ, 5, literals10, coeffs10)
			encoding, err := e.EncodePBC(fac, pbc, &config)
			assert.Nil(err)
			solver := s.NewSolver(fac, satConfig)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models := enum.OnSolver(solver, vars10)
			assert.Len(models, 9)
			for _, model := range models {
				assert.True(len(model.PosVars()) == 2)
				assert.Contains(model.PosVars(), fac.Var("v0"))
			}

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.EQ, 7, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models = enum.OnSolver(solver, vars10)
			assert.Len(models, 36)
			for _, model := range models {
				assert.True(len(model.PosVars()) == 3)
				assert.Contains(model.PosVars(), fac.Var("v0"))
			}

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.EQ, 0, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models = enum.OnSolver(solver, vars10)
			assert.Len(models, 1)

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.EQ, 1, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.False(solver.Sat())

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.EQ, 22, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.False(solver.Sat())
		}
	}
}

func TestPbEncodingNegative(t *testing.T) {
	assert := assert.New(t)
	for _, config := range pbcConfigs {
		fac := f.NewFactory()
		for _, satConfig := range getConfigs() {
			coeffs10 := []int{2, 2, 2, 2, 2, 2, 2, 2, 2, -2}
			vars10 := make([]f.Variable, 10)
			for i := 0; i < 10; i++ {
				vars10[i] = fac.Var(fmt.Sprintf("v%d", i))
			}
			literals10 := f.VariablesAsLiterals(vars10)
			pbc := fac.PBC(f.EQ, 2, literals10, coeffs10)
			encoding, err := e.EncodePBC(fac, pbc, &config)
			assert.Nil(err)
			solver := s.NewSolver(fac, satConfig)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models := enum.OnSolver(solver, vars10)
			assert.Len(models, 45)

			solver = s.NewSolver(fac, satConfig)
			encoding, err = e.EncodePBC(fac, fac.PBC(f.EQ, 4, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models = enum.OnSolver(solver, vars10)
			assert.Len(models, 120)

			solver = s.NewSolver(fac, satConfig)
			coeffs10 = []int{2, 2, -3, 2, -7, 2, 2, 2, 2, -2}
			encoding, err = e.EncodePBC(fac, fac.PBC(f.EQ, 4, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models = enum.OnSolver(solver, vars10)
			assert.Len(models, 57)

			solver = s.NewSolver(fac, satConfig)
			coeffs10 = []int{2, 2, -3, 2, -7, 2, 2, 2, 2, -2}
			encoding, err = e.EncodePBC(fac, fac.PBC(f.EQ, -10, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models = enum.OnSolver(solver, vars10)
			assert.Len(models, 8)

			solver = s.NewSolver(fac, satConfig)
			coeffs10 = []int{2, 2, -4, 2, -6, 2, 2, 2, 2, -2}
			encoding, err = e.EncodePBC(fac, fac.PBC(f.EQ, -12, literals10, coeffs10), &config)
			assert.Nil(err)
			solver.Add(encoding...)
			assert.True(solver.Sat())
			models = enum.OnSolver(solver, vars10)
			assert.Len(models, 1)
		}
	}
}

func TestPbEncodingLarge(t *testing.T) {
	assert := assert.New(t)
	for _, config := range pbcConfigs {
		fac := f.NewFactory()
		solver := s.NewSolver(fac, getConfigs()[0])
		numLits := 100
		coeffs := make([]int, numLits)
		vars := make([]f.Variable, numLits)
		for i := 0; i < numLits; i++ {
			vars[i] = fac.Var(fmt.Sprintf("v%d", i))
			coeffs[i] = i + 1
		}
		lits := f.VariablesAsLiterals(vars)
		pbc := fac.PBC(f.GE, 5000, lits, coeffs)
		encoding, err := e.EncodePBC(fac, pbc, &config)
		assert.Nil(err)
		solver.Add(encoding...)
		sResult := solver.Call(s.Params().ModelIfSat(vars))
		assert.True(sResult.Sat())
		ass, _ := sResult.Model().Assignment(fac)
		assert.True(assignment.Evaluate(fac, pbc, ass))
	}
}
