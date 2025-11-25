package np

import (
	"fmt"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/maxsat"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/sets/linkedhashset"
)

// MinimumSetCover computes a minimum set cover for a given collection of sets.
// This is a simple MaxSAT based implementation of an algorithm and is really
// only meant for small set cover problems with perhaps some tens or hundreds
// of set and hundreds of variables.
func MinimumSetCover[T any](sets [][]T) [][]T {
	fac := f.NewFactory()
	if len(sets) == 0 {
		return [][]T{}
	}
	setMap := make(map[f.Variable][]T)
	elementOccurrences := hashmap.New()
	for _, set := range sets {
		setVar := fac.Var(fmt.Sprintf("@SET_SEL_%d", len(setMap)))
		setMap[setVar] = set
		for _, element := range set {
			occs, ok := elementOccurrences.Get(element)
			if !ok {
				occs = linkedhashset.New()
				elementOccurrences.Put(element, occs)
			}
			occs.(*linkedhashset.Set).Add(setVar)
		}
	}
	solver := maxsat.OLL(fac)
	for _, _occs := range elementOccurrences.Values() {
		occs := _occs.(*linkedhashset.Set)
		ops := make([]f.Variable, 0, occs.Size())
		occs.Each(func(_ int, val any) {
			ops = append(ops, val.(f.Variable))
		})
		_ = solver.AddHardFormula(fac.Or(f.VariablesAsFormulas(ops)...))
	}
	for setVar := range setMap {
		_ = solver.AddSoftFormula(setVar.Negate(fac).AsFormula(), 1)
	}
	solverResult := solver.Solve()
	if !solverResult.Satisfiable {
		panic(errorx.IllegalState("optimization problem was not satisfiable"))
	}
	model := solverResult.Model
	pos := f.NewVarSet(model.PosVars()...)
	var minimumCover []f.Variable
	for key := range setMap {
		if pos.Contains(key) {
			minimumCover = append(minimumCover, key)
		}
	}
	result := make([][]T, len(minimumCover))
	for i, setVar := range minimumCover {
		result[i] = setMap[setVar]
	}
	return result
}
