package encoding

import (
	"slices"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// Normalize returns a normalized <= constraint of the given pseudo-Boolean or
// cardinality constraint.  Panics if the given constraint is not valid on the
// factory.
func Normalize(fac f.Factory, constraint f.Formula) f.Formula {
	comparator, rhs, literals, coefficients, found := fac.PBCOps(constraint)
	if !found {
		panic(errorx.UnknownFormula(constraint))
	}

	normPs := make([]f.Literal, len(literals))
	copy(normPs, literals)
	normCs := make([]int, len(literals))
	var normRhs int
	switch csort := comparator; csort {
	case f.EQ:
		copy(normCs, coefficients)
		normRhs = rhs
		f1 := normalizeLE(fac, normPs, normCs, normRhs)
		normCs = make([]int, len(literals))
		for i := range literals {
			normCs[i] = -coefficients[i]
		}
		normRhs = -rhs
		f2 := normalizeLE(fac, normPs, normCs, normRhs)
		return fac.And(f1, f2)
	case f.LT, f.LE:
		copy(normCs, coefficients)
		if csort == f.LE {
			normRhs = rhs
		} else {
			normRhs = rhs - 1
		}
		return normalizeLE(fac, normPs, normCs, normRhs)
	case f.GT, f.GE:
		for i := range literals {
			normCs[i] = -coefficients[i]
		}
		if csort == f.GE {
			normRhs = -rhs
		} else {
			normRhs = -rhs - 1
		}
		return normalizeLE(fac, normPs, normCs, normRhs)
	default:
		panic(errorx.UnknownEnumValue(csort))
	}
}

// NegatePBC returns the negation of the given pseudo-Boolean constraint.  In
// case of LE, LT, GE, GT the negation is computed by negating the comparator.
// In case of EQ the negation is computed by negating the whole constraint.
// Panics if the given constraint is not valid on the factory.
func NegatePBC(fac f.Factory, formula f.Formula) f.Formula {
	comparator, rhs, literals, coefficients, found := fac.PBCOps(formula)
	if !found {
		panic(errorx.UnknownFormula(formula))
	}
	switch comparator {
	case f.EQ:
		return fac.Or(
			fac.PBC(f.LT, rhs, literals, coefficients),
			fac.PBC(f.GT, rhs, literals, coefficients),
		)
	case f.LE:
		return fac.PBC(f.GT, rhs, literals, coefficients)
	case f.LT:
		return fac.PBC(f.GE, rhs, literals, coefficients)
	case f.GE:
		return fac.PBC(f.LT, rhs, literals, coefficients)
	case f.GT:
		return fac.PBC(f.LE, rhs, literals, coefficients)
	default:
		panic(errorx.UnknownEnumValue(comparator))
	}
}

func normalizeLE(fac f.Factory, ps []f.Literal, cs []int, rhs int) f.Formula {
	c := rhs
	newSize := 0
	for i := 0; i < len(ps); i++ {
		if cs[i] != 0 {
			ps[newSize] = ps[i]
			cs[newSize] = cs[i]
			newSize++
		}
	}
	removeElements(&ps, len(ps)-newSize)
	removeElements(&cs, len(cs)-newSize)
	var2consts := make(map[f.Variable][2]int)
	for i := 0; i < len(ps); i++ {
		x := ps[i].Variable()
		consts, ok := var2consts[x]
		if !ok {
			consts = [2]int{0, 0}
		}
		if ps[i].IsNeg() {
			var2consts[x] = [2]int{consts[0] + cs[i], consts[1]}
		} else {
			var2consts[x] = [2]int{consts[0], consts[1] + cs[i]}
		}
	}
	csps := make([]ilpair, len(var2consts))
	count := 0
	for k, v := range var2consts {
		if v[0] < v[1] {
			c -= v[0]
			csps[count] = ilpair{v[1] - v[0], k.AsLiteral()}
		} else {
			c -= v[1]
			csps[count] = ilpair{v[0] - v[1], k.Negate(fac)}
		}
		count++
	}
	sum := 0
	zeros := 0
	ps = []f.Literal{}
	cs = []int{}

	slices.SortFunc(csps, func(p1, p2 ilpair) int { return int(p1.second) - int(p2.second) })
	for _, pair := range csps {
		if pair.first != 0 {
			cs = append(cs, pair.first)
			ps = append(ps, pair.second)
			sum += cs[len(cs)-1]
		} else {
			zeros++
		}
	}

	var changed bool
	for ok := true; ok; ok = changed {
		changed = false
		if c < 0 {
			return fac.Falsum()
		}
		if sum <= c {
			return fac.Verum()
		}
		div := c
		for i := 0; i < len(cs); i++ {
			div = gcd(div, cs[i])
		}
		if div != 0 && div != 1 {
			for i := 0; i < len(cs); i++ {
				cs[i] = cs[i] / div
			}
			c = c / div
		}
		if div != 1 && div != 0 {
			changed = true
		}
	}
	return fac.PBC(f.LE, c, ps, cs)
}

func gcd(small, big int) int {
	if small == 0 {
		return big
	} else {
		return gcd(big%small, small)
	}
}

type ilpair struct {
	first  int
	second f.Literal
}

func removeElements[T any](slice *[]T, num int) {
	if num > 0 {
		*slice = (*slice)[:len(*slice)-num]
	}
}
