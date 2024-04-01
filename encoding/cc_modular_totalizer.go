package encoding

import (
	"math"

	f "booleworks.com/logicng/formula"
)

func modtotalizerAMK(result Result, vars []f.Literal, rhs int) *CCIncrementalData {
	state := newState(result.Factory())
	mod := initialize(result, rhs, len(vars), state)
	copy(state.inlits, vars)
	mtToCNF(result, mod, &state.cardinalityUpOutvars, &state.cardinalityLwOutvars, len(vars), state)
	encodeOutput(result, rhs, mod, state)
	state.currentCardinalityRhs = rhs + 1
	return &CCIncrementalData{
		Result:     result,
		amkEncoder: AMKModularTotalizer,
		alkEncoder: ALKModularTotalizer,
		mod:        mod,
		currentRhs: rhs,
		vector1:    state.cardinalityUpOutvars,
		vector2:    state.cardinalityLwOutvars,
	}
}

func modtotalizerALK(result Result, vars []f.Literal, rhs int) *CCIncrementalData {
	state := newState(result.Factory())
	newRhs := len(vars) - rhs
	mod := initialize(result, newRhs, len(vars), state)
	for i, v := range vars {
		state.inlits[i] = v.Negate(result.Factory())
	}
	mtToCNF(result, mod, &state.cardinalityUpOutvars, &state.cardinalityLwOutvars, len(vars), state)
	encodeOutput(result, newRhs, mod, state)
	state.currentCardinalityRhs = newRhs + 1
	return &CCIncrementalData{
		Result:     result,
		amkEncoder: AMKModularTotalizer,
		alkEncoder: ALKModularTotalizer,
		mod:        mod,
		currentRhs: rhs,
		nVars:      len(vars),
		vector1:    state.cardinalityUpOutvars,
		vector2:    state.cardinalityLwOutvars,
	}
}

func initialize(result Result, rhs, n int, state *mtstate) int {
	mod := int(math.Ceil(math.Sqrt(float64(rhs) + 1.0)))
	state.cardinalityUpOutvars = make([]f.Literal, n/mod)
	for i := 0; i < n/mod; i++ {
		state.cardinalityUpOutvars[i] = result.NewCcVariable().AsLiteral()
	}
	state.cardinalityLwOutvars = make([]f.Literal, mod-1)
	for i := 0; i < mod-1; i++ {
		state.cardinalityLwOutvars[i] = result.NewCcVariable().AsLiteral()
	}
	state.inlits = make([]f.Literal, n)
	state.currentCardinalityRhs = rhs + 1
	if len(state.cardinalityUpOutvars) == 0 {
		state.cardinalityUpOutvars = append(state.cardinalityUpOutvars, state.h0.AsLiteral())
	}
	return mod
}

func mtToCNF(result Result, mod int, ubvars, lwvars *[]f.Literal, rhs int, state *mtstate) {
	lupper := make([]f.Literal, 0, 4)
	llower := make([]f.Literal, 0, 4)
	rupper := make([]f.Literal, 0, 4)
	rlower := make([]f.Literal, 0, 4)
	split := rhs / 2
	left := 1
	right := 1
	if split == 1 {
		lupper = append(lupper, state.h0.AsLiteral())
		lupper = append(lupper, state.h0.AsLiteral())
		llower = append(llower, state.inlits[len(state.inlits)-1])
		state.inlits = state.inlits[:len(state.inlits)-1]
	} else {
		left = split / mod
		for i := 0; i < left; i++ {
			lupper = append(lupper, result.NewCcVariable().AsLiteral())
		}
		limit := mod - 1
		if left%mod == 0 && split < mod-1 {
			limit = split
		}
		for i := 0; i < limit; i++ {
			llower = append(llower, result.NewCcVariable().AsLiteral())
		}
	}
	if rhs-split == 1 {
		rupper = append(rupper, state.h0.AsLiteral())
		rlower = append(rlower, state.inlits[len(state.inlits)-1])
		state.inlits = state.inlits[:len(state.inlits)-1]
	} else {
		right := (rhs - split) / mod
		for i := 0; i < right; i++ {
			rupper = append(rupper, result.NewCcVariable().AsLiteral())
		}
		limit := mod - 1
		if right%mod == 0 && rhs-split < mod-1 {
			limit = rhs - split
		}
		for i := 0; i < limit; i++ {
			rlower = append(rlower, result.NewCcVariable().AsLiteral())
		}
	}
	if len(lupper) == 0 {
		lupper = append(lupper, state.h0.AsLiteral())
	}
	if len(rupper) == 0 {
		rupper = append(rupper, state.h0.AsLiteral())
	}
	adder(result, mod, ubvars, lwvars, &rupper, &rlower, &lupper, &llower, state)
	val := left*mod + split - left*mod
	if val > 1 {
		mtToCNF(result, mod, &lupper, &llower, val, state)
	}
	val = right*mod + (rhs - split) - right*mod
	if val > 1 {
		mtToCNF(result, mod, &rupper, &rlower, val, state)
	}
}

func adder(result Result, mod int, upper, lower, lupper, llower, rupper, rlower *[]f.Literal, state *mtstate) {
	fac := result.Factory()
	carry := state.varUndef
	if (*upper)[0] != state.h0.AsLiteral() {
		carry = result.NewCcVariable()
	}
	for i := 0; i <= len(*llower); i++ {
		for j := 0; j <= len(*rlower); j++ {
			if i+j > state.currentCardinalityRhs+1 && state.currentCardinalityRhs+1 < mod {
				continue
			}
			if i+j < mod {
				if i == 0 && j != 0 {
					if (*upper)[0] != state.h0.AsLiteral() {
						result.AddClause((*rlower)[j-1].Negate(fac), (*lower)[i+j-1], carry.AsLiteral())
					} else {
						result.AddClause((*rlower)[j-1].Negate(fac), (*lower)[i+j-1])
					}
				} else if j == 0 && i != 0 {
					if (*upper)[0] != state.h0.AsLiteral() {
						result.AddClause((*llower)[i-1].Negate(fac), (*lower)[i+j-1], carry.AsLiteral())
					} else {
						result.AddClause((*llower)[i-1].Negate(fac), (*lower)[i+j-1])
					}
				} else if i != 0 {
					if (*upper)[0] != state.h0.AsLiteral() {
						result.AddClause((*llower)[i-1].Negate(fac), (*rlower)[j-1].Negate(fac), (*lower)[i+j-1], carry.AsLiteral())
					} else {
						result.AddClause((*llower)[i-1].Negate(fac), (*rlower)[j-1].Negate(fac), (*lower)[i+j-1])
					}
				}
			} else if i+j > mod {
				result.AddClause((*llower)[i-1].Negate(fac), (*rlower)[j-1].Negate(fac), (*lower)[(i+j)%mod-1])
			} else {
				result.AddClause((*llower)[i-1].Negate(fac), (*rlower)[j-1].Negate(fac), carry.AsLiteral())
			}
		}
	}
	if (*upper)[0] != state.h0.AsLiteral() {
		finalAdder(result, mod, upper, lupper, rupper, carry, state)
	}
}

func finalAdder(result Result, mod int, upper, lupper, rupper *[]f.Literal, carry f.Variable, state *mtstate) {
	fac := result.Factory()
	for i := 0; i <= len(*lupper); i++ {
		for j := 0; j <= len(*rupper); j++ {
			a := state.varError.AsLiteral()
			b := state.varError.AsLiteral()
			c := state.varError.AsLiteral()
			d := state.varError.AsLiteral()
			closeMod := state.currentCardinalityRhs / mod
			if state.currentCardinalityRhs%mod != 0 {
				closeMod++
			}
			if mod*(i+j) > closeMod*mod {
				continue
			}
			if i != 0 {
				a = (*lupper)[i-1]
			}
			if j != 0 {
				b = (*rupper)[j-1]
			}
			if i+j != 0 && i+j-1 < len(*upper) {
				c = (*upper)[i+j-1]
			}
			if i+j < len(*upper) {
				d = (*upper)[i+j]
			}
			if c != state.varUndef.AsLiteral() && c != state.varError.AsLiteral() {
				clause := make([]f.Literal, 0, 2)
				if a != state.varUndef.AsLiteral() && a != state.varError.AsLiteral() {
					clause = append(clause, a.Negate(fac))
				}
				if b != state.varUndef.AsLiteral() && b != state.varError.AsLiteral() {
					clause = append(clause, b.Negate(fac))
				}
				clause = append(clause, c)
				if len(clause) > 1 {
					result.AddClause(clause...)
				}
			}
			clause := make([]f.Literal, 0, 2)
			clause = append(clause, carry.Negate(fac))
			if a != state.varUndef.AsLiteral() && a != state.varError.AsLiteral() {
				clause = append(clause, a.Negate(fac))
			}
			if b != state.varUndef.AsLiteral() && b != state.varError.AsLiteral() {
				clause = append(clause, b.Negate(fac))
			}
			if d != state.varError.AsLiteral() && d != state.varUndef.AsLiteral() {
				clause = append(clause, d)
			}
			if len(clause) > 1 {
				result.AddClause(clause...)
			}
		}
	}
}

func encodeOutput(result Result, rhs, mod int, state *mtstate) {
	fac := result.Factory()
	ulimit := (rhs + 1) / mod
	llimit := (rhs + 1) - ulimit*mod
	for i := ulimit; i < len(state.cardinalityUpOutvars); i++ {
		result.AddClause(state.cardinalityUpOutvars[i].Negate(fac))
	}
	if ulimit != 0 && llimit != 0 {
		for i := llimit - 1; i < len(state.cardinalityLwOutvars); i++ {
			result.AddClause(state.cardinalityUpOutvars[ulimit-1].Negate(fac), state.cardinalityLwOutvars[i].Negate(fac))
		}
	} else {
		if ulimit == 0 {
			for i := llimit - 1; i < len(state.cardinalityLwOutvars); i++ {
				result.AddClause(state.cardinalityLwOutvars[i].Negate(fac))
			}
		} else {
			result.AddClause(state.cardinalityUpOutvars[ulimit-1].Negate(fac))
		}
	}
}

type mtstate struct {
	varUndef              f.Variable
	varError              f.Variable
	h0                    f.Variable
	inlits                []f.Literal
	cardinalityUpOutvars  []f.Literal
	cardinalityLwOutvars  []f.Literal
	currentCardinalityRhs int
}

func newState(fac f.Factory) *mtstate {
	varUndef := fac.Var("RESERVED@VAR_UNDEF")
	return &mtstate{
		varUndef:              varUndef,
		varError:              fac.Var("RESERVED@VAR_ERROR"),
		h0:                    varUndef,
		inlits:                make([]f.Literal, 0),
		currentCardinalityRhs: -1,
	}
}
