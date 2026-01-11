package encoding

import (
	"math"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

func amoPure(result Result, vars []f.Variable) {
	fac := result.Factory()
	for i := range vars {
		for j := i + 1; j < len(vars); j++ {
			result.AddClause(vars[i].Negate(fac), vars[j].Negate(fac))
		}
	}
}

func amoLadder(result Result, vars []f.Variable) {
	fac := result.Factory()
	seqAuxiliary := make([]f.Variable, len(vars)-1)
	for i := 0; i < len(vars)-1; i++ {
		seqAuxiliary[i] = result.NewAuxVar(f.AuxCC)
	}
	for i := range vars {
		if i == 0 {
			result.AddClause(vars[0].Negate(fac), seqAuxiliary[0].AsLiteral())
		} else if i == len(vars)-1 {
			result.AddClause(vars[i].Negate(fac), seqAuxiliary[i-1].Negate(fac))
		} else {
			result.AddClause(vars[i].Negate(fac), seqAuxiliary[i].AsLiteral())
			result.AddClause(seqAuxiliary[i-1].Negate(fac), seqAuxiliary[i].AsLiteral())
			result.AddClause(vars[i].Negate(fac), seqAuxiliary[i-1].Negate(fac))
		}
	}
}

func amoProduct(result Result, recursiveBound int, vars []f.Variable) {
	if recursiveBound == 0 {
		recursiveBound = 20
	}
	fac := result.Factory()
	n := len(vars)
	p := int(math.Ceil(math.Sqrt(float64(n))))
	q := int(math.Ceil(float64(n) / float64(p)))
	us := make([]f.Variable, p)
	for i := range us {
		us[i] = result.NewAuxVar(f.AuxCC)
	}
	vs := make([]f.Variable, q)
	for i := range vs {
		vs[i] = result.NewAuxVar(f.AuxCC)
	}
	if len(us) <= recursiveBound {
		buildPure(result, us)
	} else {
		amoProduct(result, recursiveBound, us)
	}
	if len(vs) <= recursiveBound {
		buildPure(result, vs)
	} else {
		amoProduct(result, recursiveBound, vs)
	}
	for i := range p {
		for j := range q {
			k := i*q + j
			if k >= 0 && k < n {
				result.AddClause(vars[k].Negate(fac), us[i].AsLiteral())
				result.AddClause(vars[k].Negate(fac), vs[j].AsLiteral())
			}
		}
	}
}

func buildPure(result Result, vars []f.Variable) {
	for i := range vars {
		for j := i + 1; j < len(vars); j++ {
			result.AddClause(vars[i].Negate(result.Factory()), vars[j].Negate(result.Factory()))
		}
	}
}

func buildPureLit(result Result, lits []f.Literal) {
	for i := range lits {
		for j := i + 1; j < len(lits); j++ {
			result.AddClause(lits[i].Negate(result.Factory()), lits[j].Negate(result.Factory()))
		}
	}
}

func amoNested(result Result, groupSize int, vars []f.Literal) {
	if groupSize == 0 {
		groupSize = 4
	}
	fac := result.Factory()
	if len(vars) <= groupSize {
		for i := 0; i+1 < len(vars); i++ {
			for j := i + 1; j < len(vars); j++ {
				result.AddClause(vars[i].Negate(fac), vars[j].Negate(fac))
			}
		}
	} else {
		l1 := make([]f.Literal, 0, len(vars)/2)
		l2 := make([]f.Literal, 0, len(vars)/2)
		i := 0
		for ; i < len(vars)/2; i++ {
			l1 = append(l1, vars[i])
		}
		for ; i < len(vars); i++ {
			l2 = append(l2, vars[i])
		}
		newVariable := result.NewAuxVar(f.AuxCC)
		l1 = append(l1, newVariable.AsLiteral())
		l2 = append(l2, newVariable.Negate(fac))
		amoNested(result, groupSize, l1)
		amoNested(result, groupSize, l2)
	}
}

func amoCommander(result Result, groupSize int, vars []f.Literal) {
	if groupSize == 0 {
		groupSize = 3
	}
	fac := result.Factory()
	isExactlyOne := false
	for len(vars) > groupSize {
		literals := make([]f.Literal, 0, len(vars))
		nextLiterals := make([]f.Literal, 0, 4)
		for i := 0; i < len(vars); i++ {
			literals = append(literals, vars[i])
			if i%groupSize == groupSize-1 || i == len(vars)-1 {
				buildPureLit(result, literals)
				literals = append(literals, result.NewAuxVar(f.AuxCC).AsLiteral())
				nextLiterals = append(nextLiterals, literals[len(literals)-1].Negate(fac))
				if isExactlyOne && len(literals) > 0 {
					result.AddClause(literals...)
				}
				for j := 0; j < len(literals)-1; j++ {
					result.AddClause(literals[len(literals)-1].Negate(fac), literals[j].Negate(fac))
				}
				literals = make([]f.Literal, 0, len(vars))
			}
		}
		vars = nextLiterals
		isExactlyOne = true
	}
	buildPureLit(result, vars)
	if isExactlyOne && len(vars) > 0 {
		result.AddClause(vars...)
	}
}

func amoBinary(result Result, vars []f.Variable) {
	fac := result.Factory()
	numberOfBits := int(math.Ceil(math.Log(float64(len(vars))) / math.Log(2)))
	twoPowNBits := int(math.Pow(2, float64(numberOfBits)))
	k := (twoPowNBits - len(vars)) * 2
	bits := make([]f.Variable, numberOfBits)
	for i := range numberOfBits {
		bits[i] = result.NewAuxVar(f.AuxCC)
	}
	var grayCode, nextGray int
	i := 0
	index := -1
	for i < k {
		index++
		grayCode = i ^ (i >> 1)
		i++
		nextGray = i ^ (i >> 1)
		for j := range numberOfBits {
			if (grayCode & (1 << j)) == (nextGray & (1 << j)) {
				if (grayCode & (1 << j)) != 0 {
					result.AddClause(vars[index].Negate(fac), bits[j].AsLiteral())
				} else {
					result.AddClause(vars[index].Negate(fac), bits[j].Negate(fac))
				}
			}
		}
		i++
	}
	for i < twoPowNBits {
		index++
		grayCode = i ^ (i >> 1)
		for j := range numberOfBits {
			if (grayCode & (1 << j)) != 0 {
				result.AddClause(vars[index].Negate(fac), bits[j].AsLiteral())
			} else {
				result.AddClause(vars[index].Negate(fac), bits[j].Negate(fac))
			}
		}
		i++
	}
}

func amoBimander(result Result, groupSize BimanderGroupSize, fixed int, vars []f.Variable) {
	gs := computeGroupSize(groupSize, fixed, len(vars))
	bimanderIntern(result, gs, vars)
}

func bimanderIntern(result Result, groupSize int, vars []f.Variable) {
	groups := initializeGroups(result, groupSize, vars)
	bits := initializeBits(result, groupSize)
	var grayCode, nextGray int
	i := 0
	index := -1
	for ; i < bits.k; i++ {
		index++
		grayCode = i ^ (i >> 1)
		i++
		nextGray = i ^ (i >> 1)
		for j := 0; j < bits.numberOfBits; j++ {
			if (grayCode & (1 << j)) == (nextGray & (1 << j)) {
				handleGrayCode(result, groups, bits, grayCode, index, j)
			}
		}
	}
	for ; i < bits.twoPowNBits; i++ {
		index++
		grayCode = i ^ (i >> 1)
		for j := 0; j < bits.numberOfBits; j++ {
			handleGrayCode(result, groups, bits, grayCode, index, j)
		}
	}
}

func initializeGroups(result Result, groupSize int, vars []f.Variable) [][]f.Literal {
	groups := make([][]f.Literal, 0, 4)
	n := len(vars)
	for range groupSize {
		groups = append(groups, make([]f.Literal, 0, 4))
	}

	g := int(math.Ceil(float64(n) / float64(groupSize)))
	ig := 0
	for i := 0; i < len(vars); {
		for i < g {
			groups[ig] = append(groups[ig], vars[i].AsLiteral())
			i++
		}
		ig++
		g = g + int(math.Ceil(float64(n-i)/float64(groupSize-ig)))
	}
	for i := 0; i < len(groups); i++ {
		buildPureLit(result, groups[i])
	}
	return groups
}

func initializeBits(result Result, groupSize int) *bimanderbits {
	bits := &bimanderbits{}
	bits.numberOfBits = int(math.Ceil(math.Log(float64(groupSize)) / math.Log(2)))
	bits.twoPowNBits = int(math.Pow(2, float64(bits.numberOfBits)))
	bits.k = (bits.twoPowNBits - groupSize) * 2
	for i := 0; i < bits.numberOfBits; i++ {
		bits.bits = append(bits.bits, result.NewAuxVar(f.AuxCC).AsLiteral())
	}
	return bits
}

func handleGrayCode(result Result, groups [][]f.Literal, bits *bimanderbits, grayCode, index, j int) {
	if (grayCode & (1 << j)) != 0 {
		for p := 0; p < len(groups[index]); p++ {
			result.AddClause(groups[index][p].Negate(result.Factory()), bits.bits[j])
		}
	} else {
		for p := 0; p < len(groups[index]); p++ {
			result.AddClause(groups[index][p].Negate(result.Factory()), bits.bits[j].Negate(result.Factory()))
		}
	}
}

func computeGroupSize(groupSize BimanderGroupSize, fixed, numVars int) int {
	switch groupSize {
	case BimanderFixed:
		if fixed == 0 {
			fixed = 3
		}
		return fixed
	case BimanderHalf:
		return numVars / 2
	case BimanderSqrt:
		return int(math.Sqrt(float64(numVars)))
	default:
		panic(errorx.UnknownEnumValue(groupSize))
	}
}

type bimanderbits struct {
	bits         []f.Literal
	numberOfBits int
	twoPowNBits  int
	k            int
}
