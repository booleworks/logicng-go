package sat

import (
	"strconv"

	f "booleworks.com/logicng/formula"
)

// GeneratePigeonHole generates a pigeon hole problem of size n and returns it
// as a formula.
func GeneratePigeonHole(fac f.Factory, n int) f.Formula {
	prefix := "v"
	return fac.And(placeInSomeHole(n, fac, prefix), onlyOnePigeonInHole(n, fac, prefix))
}

func placeInSomeHole(n int, fac f.Factory, prefix string) f.Formula {
	if n == 1 {
		return fac.And(fac.Variable(prefix+"1"), fac.Variable(prefix+"2"))
	}
	ors := make([]f.Formula, 0, n)
	for i := 1; i <= n+1; i++ {
		orOps := make([]f.Formula, 0, n)
		for j := 1; j <= n; j++ {
			orOps = append(orOps, fac.Variable(prefix+strconv.FormatInt(int64(n*(i-1)+j), 10)))
		}
		ors = append(ors, fac.Or(orOps...))
	}
	return fac.And(ors...)
}

func onlyOnePigeonInHole(n int, fac f.Factory, prefix string) f.Formula {
	if n == 1 {
		return fac.Or(fac.Literal(prefix+"1", false), fac.Literal(prefix+"2", false))
	}
	ors := make([]f.Formula, 0, n*n*(n+1)/2)
	for j := 1; j <= n; j++ {
		for i := 1; i <= n; i++ {
			for k := i + 1; k <= n+1; k++ {
				ors = append(ors, fac.Or(fac.Literal(prefix+strconv.FormatInt(int64(n*(i-1)+j), 10), false),
					fac.Literal(prefix+strconv.FormatInt(int64(n*(k-1)+j), 10), false)))
			}
		}
	}
	return fac.And(ors...)
}
