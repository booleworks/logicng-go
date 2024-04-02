package encoding

import (
	"math"

	f "github.com/booleworks/logicng-go/formula"
)

func cnAmk(result Result, vars []f.Variable, rhs int) {
	fac := result.Factory()
	input := make([]f.Literal, 0, len(vars))
	output := make([]f.Literal, 0, 4)

	if rhs > len(vars)/2 {
		geq := len(vars) - rhs
		for _, v := range vars {
			input = append(input, v.Negate(fac))
		}
		sort(geq, &input, &output, result, outputToInput)
		for i := 0; i < geq; i++ {
			result.AddClause(output[i])
		}
	} else {
		input = append(input, f.VariablesAsLiterals(vars)...)
		sort(rhs+1, &input, &output, result, inputToOutput)
		result.AddClause(output[rhs].Negate(fac))
	}
}

func cnAmkForIncremental(result Result, vars []f.Variable, rhs int) *CCIncrementalData {
	fac := result.Factory()
	input := make([]f.Literal, len(vars))
	copy(input, f.VariablesAsLiterals(vars))
	output := make([]f.Literal, 0, 4)
	sort(rhs+1, &input, &output, result, inputToOutput)
	result.AddClause(output[rhs].Negate(fac))
	return &CCIncrementalData{
		Result:     result,
		amkEncoder: AMKCardinalityNetwork,
		alkEncoder: ALKCardinalityNetwork,
		currentRhs: rhs,
		vector1:    output,
	}
}

func cnAlk(result Result, vars []f.Variable, rhs int) {
	fac := result.Factory()
	input := make([]f.Literal, 0, len(vars))
	output := make([]f.Literal, 0, 4)
	newRhs := len(vars) - rhs

	if newRhs > len(vars)/2 {
		geq := len(vars) - newRhs
		input = append(input, f.VariablesAsLiterals(vars)...)
		sort(geq, &input, &output, result, outputToInput)
		for i := 0; i < geq; i++ {
			result.AddClause(output[i])
		}
	} else {
		for _, v := range vars {
			input = append(input, v.Negate(fac))
		}
		sort(newRhs+1, &input, &output, result, inputToOutput)
		result.AddClause(output[newRhs].Negate(fac))
	}
}

func cnAlkForIncremental(result Result, vars []f.Variable, rhs int) *CCIncrementalData {
	fac := result.Factory()
	input := make([]f.Literal, len(vars))
	for i, v := range vars {
		input[i] = v.Negate(fac)
	}
	output := make([]f.Literal, 0, 4)
	newRhs := len(vars) - rhs
	sort(newRhs+1, &input, &output, result, inputToOutput)
	result.AddClause(output[newRhs].Negate(fac))
	return &CCIncrementalData{
		Result:     result,
		amkEncoder: AMKCardinalityNetwork,
		alkEncoder: ALKCardinalityNetwork,
		currentRhs: rhs,
		nVars:      len(vars),
		vector1:    output,
	}
}

func cnExk(result Result, vars []f.Variable, rhs int) {
	fac := result.Factory()
	input := make([]f.Literal, 0, len(vars))
	output := make([]f.Literal, 0, 4)
	input = append(input, f.VariablesAsLiterals(vars)...)
	sort(rhs+1, &input, &output, result, both)
	result.AddClause(output[rhs].Negate(fac))
	result.AddClause(output[rhs-1])
}

type dir byte

const (
	inputToOutput dir = iota
	outputToInput
	both
)

func counterSorterValue(m, n int) int {
	return 2*n + (m-1)*(2*(n-1)-1) - (m - 2) - 2*((m-1)*(m-2)/2)
}

func directSorterValue(n int) int {
	if n > 30 {
		return math.MaxInt
	}
	return int(math.Pow(2, float64(n)) - 1)
}

func sort(m int, input, output *[]f.Literal, result Result, direction dir) {
	if m == 0 {
		clear(*output)
		return
	}
	n := len(*input)
	m2 := m
	if m2 > n {
		m2 = n
	}
	if n == 0 {
		clear(*output)
		return
	}
	if n == 1 {
		clear(*output)
		*output = append(*output, (*input)[0])
		return
	}
	if n == 2 {
		clear(*output)
		o1 := result.NewAuxVar(f.AuxCC).AsLiteral()
		if m2 == 2 {
			o2 := result.NewAuxVar(f.AuxCC).AsLiteral()
			comparator4((*input)[0], (*input)[1], o1, o2, result, direction)
			*output = append(*output, o1, o2)
		} else {
			comparator3((*input)[0], (*input)[1], o1, result, direction)
			*output = append(*output, o1)
		}
		return
	}
	if direction != inputToOutput {
		recursiveSorter(m2, input, output, result, direction)
		return
	}
	counter := counterSorterValue(m2, n)
	direct := directSorterValue(n)

	if counter < direct {
		counterSorter(m2, input, output, result, direction)
	} else {
		directSorter(m2, input, output, result)
	}
}

func comparator3(x1, x2, y f.Literal, result Result, direction dir) {
	if direction == inputToOutput || direction == both {
		result.AddClause(x1.Negate(result.Factory()), y)
		result.AddClause(x2.Negate(result.Factory()), y)
	}
	if direction == outputToInput || direction == both {
		result.AddClause(y.Negate(result.Factory()), x1, x2)
	}
}

func comparator4(x1, x2, y1, y2 f.Literal, result Result, direction dir) {
	if direction == inputToOutput || direction == both {
		result.AddClause(x1.Negate(result.Factory()), y1)
		result.AddClause(x2.Negate(result.Factory()), y1)
		result.AddClause(x1.Negate(result.Factory()), x2.Negate(result.Factory()), y2)
	}
	if direction == outputToInput || direction == both {
		result.AddClause(y1.Negate(result.Factory()), x1, x2)
		result.AddClause(y2.Negate(result.Factory()), x1)
		result.AddClause(y2.Negate(result.Factory()), x2)
	}
}

func recursiveSorter(m int, input, output *[]f.Literal, result Result, direction dir) {
	clear(*output)
	n := len(*input)
	l := n / 2
	tmpLitsA := make([]f.Literal, l)
	tmpLitsB := make([]f.Literal, 0, n-l)
	tmpLitsO1 := make([]f.Literal, 0)
	tmpLitsO2 := make([]f.Literal, 0)
	for i := 0; i < l; i++ {
		tmpLitsA[i] = (*input)[i]
	}
	for i := l; i < n; i++ {
		tmpLitsB = append(tmpLitsB, (*input)[i])
	}
	sort(m, &tmpLitsA, &tmpLitsO1, result, direction)
	sort(m, &tmpLitsB, &tmpLitsO2, result, direction)
	merge(m, &tmpLitsO1, &tmpLitsO2, output, result, direction)
}

func counterSorter(k int, x, output *[]f.Literal, result Result, direction dir) {
	fac := result.Factory()
	n := len(*x)
	auxVars := make([][]f.Literal, n)
	for i := 0; i < n; i++ {
		auxVars[i] = make([]f.Literal, k)
	}

	for j := 0; j < k; j++ {
		for i := j; i < n; i++ {
			auxVars[i][j] = result.NewAuxVar(f.AuxCC).AsLiteral()
		}
	}
	if direction == inputToOutput || direction == both {
		for i := 0; i < n; i++ {
			result.AddClause((*x)[i].Negate(fac), auxVars[i][0])
			if i > 0 {
				result.AddClause(auxVars[i-1][0].Negate(fac), auxVars[i][0])
			}
		}
		for j := 1; j < k; j++ {
			for i := j; i < n; i++ {
				result.AddClause((*x)[i].Negate(fac), auxVars[i-1][j-1].Negate(fac), auxVars[i][j])
				if i > j {
					result.AddClause(auxVars[i-1][j].Negate(fac), auxVars[i][j])
				}
			}
		}
	}
	clear(*output)
	for i := 0; i < k; i++ {
		*output = append(*output, auxVars[n-1][i])
	}
}

func directSorter(m int, input, output *[]f.Literal, result Result) {
	n := len(*input)
	bitmask := 1
	clause := make([]f.Literal, 0)
	clear(*output)
	for i := 0; i < m; i++ {
		*output = append(*output, result.NewAuxVar(f.AuxCC).AsLiteral())
	}
	for bitmask < int(math.Pow(2, float64(n))) {
		count := 0
		clear(clause)
		for i := 0; i < n; i++ {
			if ((1 << i) & bitmask) != 0 {
				count++
				if count > m {
					break
				}
				clause = append(clause, (*input)[i].Negate(result.Factory()))
			}
		}
		if count <= m {
			clause = append(clause, (*output)[count-1])
			result.AddClause(clause...)
		}
		bitmask++
	}
}

func merge(m int, inputA, inputB, output *[]f.Literal, result Result, direction dir) {
	if m == 0 {
		clear(*output)
		return
	}
	a := len(*inputA)
	b := len(*inputB)
	n := a + b
	m2 := m
	if m2 > n {
		m2 = n
	}
	if a == 0 || b == 0 {
		if a == 0 {
			*output = *inputB
		} else {
			*output = *inputA
		}
		return
	}
	if direction != inputToOutput {
		recursiveMerger(m2, inputA, len(*inputA), inputB, len(*inputB), result, output, direction)
		return
	}
	directMerger(m2, inputA, inputB, output, result)
}

func recursiveMerger(
	c int, inputA *[]f.Literal, a int, inputB *[]f.Literal, b int, result Result, output *[]f.Literal, direction dir,
) {
	clear(*output)
	a2 := a
	b2 := b
	if a2 > c {
		a2 = c
	}
	if b2 > c {
		b2 = c
	}
	if c == 1 {
		y := result.NewAuxVar(f.AuxCC).AsLiteral()
		comparator3((*inputA)[0], (*inputB)[0], y, result, direction)
		*output = append(*output, y)
		return
	}
	if a2 == 1 && b2 == 1 {
		y1 := result.NewAuxVar(f.AuxCC).AsLiteral()
		y2 := result.NewAuxVar(f.AuxCC).AsLiteral()
		comparator4((*inputA)[0], (*inputB)[0], y1, y2, result, direction)
		*output = append(*output, y1, y2)
		return
	}
	oddMerge := make([]f.Literal, 0)
	evenMerge := make([]f.Literal, 0)
	tmpLitsOddA := make([]f.Literal, 0, (a2/2)+1)
	tmpLitsOddB := make([]f.Literal, 0, (b2/2)+1)
	tmpLitsEvenA := make([]f.Literal, 0, (a2/2)+1)
	tmpLitsEvenB := make([]f.Literal, 0, (b2/2)+1)

	for i := 0; i < a2; i = i + 2 {
		tmpLitsOddA = append(tmpLitsOddA, (*inputA)[i])
	}
	for i := 0; i < b2; i = i + 2 {
		tmpLitsOddB = append(tmpLitsOddB, (*inputB)[i])
	}
	for i := 1; i < a2; i = i + 2 {
		tmpLitsEvenA = append(tmpLitsEvenA, (*inputA)[i])
	}
	for i := 1; i < b2; i = i + 2 {
		tmpLitsEvenB = append(tmpLitsEvenB, (*inputB)[i])
	}
	merge(c/2+1, &tmpLitsOddA, &tmpLitsOddB, &oddMerge, result, direction)
	merge(c/2, &tmpLitsEvenA, &tmpLitsEvenB, &evenMerge, result, direction)

	*output = append(*output, oddMerge[0])

	i := 1
	j := 0
	for {
		if i < len(oddMerge) && j < len(evenMerge) {
			if len(*output)+2 <= c {
				z0 := result.NewAuxVar(f.AuxCC).AsLiteral()
				z1 := result.NewAuxVar(f.AuxCC).AsLiteral()
				comparator4(oddMerge[i], evenMerge[j], z0, z1, result, direction)
				*output = append(*output, z0, z1)
				if len(*output) == c {
					break
				}
			} else if len(*output)+1 == c {
				z0 := result.NewAuxVar(f.AuxCC).AsLiteral()
				comparator3(oddMerge[i], evenMerge[j], z0, result, direction)
				*output = append(*output, z0)
				break
			}
		} else if i >= len(oddMerge) && j >= len(evenMerge) {
			break
		} else if i >= len(oddMerge) {
			*output = append(*output, evenMerge[len(evenMerge)-1])
			break
		} else {
			*output = append(*output, oddMerge[len(oddMerge)-1])
			break
		}
		i++
		j++
	}
}

func directMerger(m int, inputA, inputB, output *[]f.Literal, result Result) {
	fac := result.Factory()
	a := len(*inputA)
	b := len(*inputB)
	for i := 0; i < m; i++ {
		*output = append(*output, result.NewAuxVar(f.AuxCC).AsLiteral())
	}
	j := min(m, a)
	for i := 0; i < j; i++ {
		result.AddClause((*inputA)[i].Negate(fac), (*output)[i])
	}
	j = min(m, b)
	for i := 0; i < j; i++ {
		result.AddClause((*inputB)[i].Negate(fac), (*output)[i])
	}
	for i := 0; i < a; i++ {
		for k := 0; k < b; k++ {
			if i+k+1 < m {
				result.AddClause((*inputA)[i].Negate(fac), (*inputB)[k].Negate(fac), (*output)[i+k+1])
			}
		}
	}
}
