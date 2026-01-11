package encoding

import (
	"slices"

	f "github.com/booleworks/logicng-go/formula"
)

const nullLit = "@NULL_LIT"

func encodePBCAdder(result Result, lits []f.Literal, coeffs []int, rhs int) {
	nb := ldInt(rhs)
	literals := make([]f.Literal, nb)
	buckets := make([][]f.Literal, 0, nb)
	nullLit := result.Factory().Lit(nullLit, true)
	for iBit := range nb {
		buckets = append(buckets, []f.Literal{})
		literals[iBit] = nullLit
		for iVar := range lits {
			if ((1 << iBit) & coeffs[iVar]) != 0 {
				buckets[len(buckets)-1] = append(buckets[len(buckets)-1], lits[iVar])
			}
		}
	}
	adderTree(result, &buckets, &literals, nullLit)
	kBits := numToBits(len(buckets), rhs)
	lessThanOrEqual(result, literals, kBits, nullLit)
}

func ldInt(x int) int {
	ldretutn := 0
	for i := range 31 {
		if (x & (1 << i)) > 0 {
			ldretutn = i + 1
		}
	}
	return ldretutn
}

func adderTree(result Result, buckets *[][]f.Literal, literals *[]f.Literal, nullLit f.Literal) {
	var x, y, z f.Literal
	for i := 0; i < len(*buckets); i++ {
		if len((*buckets)[i]) == 0 {
			continue
		}
		if i == len(*buckets)-1 && len((*buckets)[i]) >= 2 {
			*buckets = append(*buckets, []f.Literal{})
			*literals = append(*literals, nullLit)
		}
		for len((*buckets)[i]) >= 3 {
			x = removeFirst(&(*buckets)[i])
			y = removeFirst(&(*buckets)[i])
			z = removeFirst(&(*buckets)[i])
			xs := faSum(result, x, y, z)
			xc := faCarry(result, x, y, z)
			(*buckets)[i] = append((*buckets)[i], xs)
			(*buckets)[i+1] = append((*buckets)[i+1], xc)
			faExtra(result, xc, xs, x, y, z)
		}
		if len((*buckets)[i]) == 2 {
			x = removeFirst(&(*buckets)[i])
			y = removeFirst(&(*buckets)[i])
			(*buckets)[i] = append((*buckets)[i], haSum(result, x, y))
			(*buckets)[i+1] = append((*buckets)[i+1], haCarry(result, x, y))
		}
		(*literals)[i] = removeFirst(&(*buckets)[i])
	}
}

func faSum(result Result, a, b, c f.Literal) f.Literal {
	fac := result.Factory()
	x := result.NewAuxVar(f.AuxPBC).AsLiteral()
	result.AddClause(a, b, c, x.Negate(fac))
	result.AddClause(a, b.Negate(fac), c.Negate(fac), x.Negate(fac))
	result.AddClause(a.Negate(fac), b, c.Negate(fac), x.Negate(fac))
	result.AddClause(a.Negate(fac), b.Negate(fac), c, x.Negate(fac))
	result.AddClause(a.Negate(fac), b.Negate(fac), c.Negate(fac), x)
	result.AddClause(a.Negate(fac), b, c, x)
	result.AddClause(a, b.Negate(fac), c, x)
	result.AddClause(a, b, c.Negate(fac), x)
	return x
}

func faCarry(result Result, a, b, c f.Literal) f.Literal {
	fac := result.Factory()
	x := result.NewAuxVar(f.AuxPBC).AsLiteral()
	result.AddClause(b, c, x.Negate(fac))
	result.AddClause(a, c, x.Negate(fac))
	result.AddClause(a, b, x.Negate(fac))
	result.AddClause(b.Negate(fac), c.Negate(fac), x)
	result.AddClause(a.Negate(fac), c.Negate(fac), x)
	result.AddClause(a.Negate(fac), b.Negate(fac), x)
	return x
}

func faExtra(result Result, xc, xs, a, b, c f.Literal) {
	fac := result.Factory()
	result.AddClause(xc.Negate(fac), xs.Negate(fac), a)
	result.AddClause(xc.Negate(fac), xs.Negate(fac), b)
	result.AddClause(xc.Negate(fac), xs.Negate(fac), c)
	result.AddClause(xc, xs, a.Negate(fac))
	result.AddClause(xc, xs, b.Negate(fac))
	result.AddClause(xc, xs, c.Negate(fac))
}

func haSum(result Result, a, b f.Literal) f.Literal {
	fac := result.Factory()
	x := result.NewAuxVar(f.AuxPBC).AsLiteral()
	result.AddClause(a.Negate(fac), b.Negate(fac), x.Negate(fac))
	result.AddClause(a, b, x.Negate(fac))
	result.AddClause(a.Negate(fac), b, x)
	result.AddClause(a, b.Negate(fac), x)
	return x
}

func haCarry(result Result, a, b f.Literal) f.Literal {
	fac := result.Factory()
	x := result.NewAuxVar(f.AuxPBC).AsLiteral()
	result.AddClause(a, x.Negate(fac))
	result.AddClause(b, x.Negate(fac))
	result.AddClause(a.Negate(fac), b.Negate(fac), x)
	return x
}

func numToBits(n, num int) []bool {
	number := num
	bits := make([]bool, 0, n)
	for i := n - 1; i >= 0; i-- {
		tmp := 1 << i
		if number < tmp {
			bits = append(bits, false)
		} else {
			bits = append(bits, true)
			number -= tmp
		}
	}
	slices.Reverse(bits)
	return bits
}

func lessThanOrEqual(result Result, xs []f.Literal, ys []bool, nullLit f.Literal) {
	fac := result.Factory()
	var clause []f.Literal
	var skip bool
	for i := range xs {
		if ys[i] || xs[i] == nullLit {
			continue
		}
		clause = make([]f.Literal, 0, len(xs)+1)
		skip = false
		for j := i + 1; j < len(xs); j++ {
			if ys[j] {
				if xs[j] == nullLit {
					skip = true
					break
				}
				clause = append(clause, xs[j].Negate(fac))
			} else {
				if xs[j] == nullLit {
					continue
				}
				clause = append(clause, xs[j])
			}
		}
		if skip {
			continue
		}
		clause = append(clause, xs[i].Negate(fac))
		result.AddClause(clause...)
	}
}

func removeFirst(slice *[]f.Literal) f.Literal {
	first := (*slice)[0]
	*slice = (*slice)[1:]
	return first
}
