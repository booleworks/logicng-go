package encoding

import f "github.com/booleworks/logicng-go/formula"

type bound byte

const (
	boundLower bound = iota
	boundUpper
	boundBoth
)

func totalizerAMK(result Result, vars []f.Variable, rhs int) *CCIncrementalData {
	tv := initializeConstraint(result, vars)
	toCNF(result, tv, rhs, boundUpper)
	for i := rhs; i < len(*tv.outvars); i++ {
		result.AddClause((*tv.outvars)[i].Negate(result.Factory()))
	}
	outvars := f.VariablesAsLiterals(*tv.outvars)
	return &CCIncrementalData{
		Result:     result,
		amkEncoder: AMKTotalizer,
		alkEncoder: ALKTotalizer,
		currentRhs: rhs,
		vector1:    outvars,
	}
}

func totalizerALK(result Result, vars []f.Variable, rhs int) *CCIncrementalData {
	tv := initializeConstraint(result, vars)
	toCNF(result, tv, rhs, boundLower)
	for i := range rhs {
		result.AddClause((*tv.outvars)[i].AsLiteral())
	}
	outvars := f.VariablesAsLiterals(*tv.outvars)
	return &CCIncrementalData{
		Result:     result,
		amkEncoder: AMKTotalizer,
		alkEncoder: ALKTotalizer,
		currentRhs: rhs,
		nVars:      len(vars),
		vector1:    outvars,
	}
}

func totalizerEXK(result Result, vars []f.Variable, rhs int) {
	tv := initializeConstraint(result, vars)
	toCNF(result, tv, rhs, boundBoth)
	for i := range rhs {
		result.AddClause((*tv.outvars)[i].AsLiteral())
	}
	for i := rhs; i < len(*tv.outvars); i++ {
		result.AddClause((*tv.outvars)[i].Negate(result.Factory()))
	}
}

func initializeConstraint(result Result, vars []f.Variable) *totalizerVars {
	invars := make([]f.Variable, len(vars))
	copy(invars, vars)
	outvars := make([]f.Variable, len(vars))
	for i := range vars {
		outvars[i] = result.NewAuxVar(f.AuxCC)
	}
	return &totalizerVars{&invars, &outvars}
}

func toCNF(result Result, tv *totalizerVars, rhs int, bound bound) {
	split := len(*tv.outvars) / 2
	left := make([]f.Variable, 0, split)
	right := make([]f.Variable, 0, len(*tv.outvars)-split)
	for i := 0; i < len(*tv.outvars); i++ {
		if i < split {
			if split == 1 {
				left = append(left, (*tv.invars)[len(*tv.invars)-1])
				*tv.invars = (*tv.invars)[:len(*tv.invars)-1]
			} else {
				left = append(left, result.NewAuxVar(f.AuxCC))
			}
		} else {
			if len(*tv.outvars)-split == 1 {
				right = append(right, (*tv.invars)[len(*tv.invars)-1])
				*tv.invars = (*tv.invars)[:len(*tv.invars)-1]
			} else {
				right = append(right, result.NewAuxVar(f.AuxCC))
			}
		}
	}
	if bound == boundUpper || bound == boundBoth {
		adderAMK(result, &left, &right, tv.outvars, rhs)
	}
	if bound == boundLower || bound == boundBoth {
		adderALK(result, &left, &right, tv.outvars)
	}
	if len(left) > 1 {
		toCNF(result, &totalizerVars{tv.invars, &left}, rhs, bound)
	}
	if len(right) > 1 {
		toCNF(result, &totalizerVars{tv.invars, &right}, rhs, bound)
	}
}

func adderAMK(result Result, left, right, output *[]f.Variable, rhs int) {
	fac := result.Factory()
	for i := 0; i <= len(*left); i++ {
		for j := 0; j <= len(*right); j++ {
			if i == 0 && j == 0 {
				continue
			}
			if i+j > rhs+1 {
				continue
			}
			if i == 0 {
				result.AddClause((*right)[j-1].Negate(fac), (*output)[j-1].AsLiteral())
			} else if j == 0 {
				result.AddClause((*left)[i-1].Negate(fac), (*output)[i-1].AsLiteral())
			} else {
				result.AddClause((*left)[i-1].Negate(fac), (*right)[j-1].Negate(fac), (*output)[i+j-1].AsLiteral())
			}
		}
	}
}

func adderALK(result Result, left, right, output *[]f.Variable) {
	fac := result.Factory()
	for i := 0; i <= len(*left); i++ {
		for j := 0; j <= len(*right); j++ {
			if i == 0 && j == 0 {
				continue
			}
			if i == 0 {
				result.AddClause((*right)[j-1].AsLiteral(), (*output)[len(*left)+j-1].Negate(fac))
			} else if j == 0 {
				result.AddClause((*left)[i-1].AsLiteral(), (*output)[len(*right)+i-1].Negate(fac))
			} else {
				result.AddClause((*left)[i-1].AsLiteral(), (*right)[j-1].AsLiteral(), (*output)[i+j-2].Negate(fac))
			}
		}
	}
}

type totalizerVars struct {
	invars  *[]f.Variable
	outvars *[]f.Variable
}
