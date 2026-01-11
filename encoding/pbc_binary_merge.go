package encoding

import (
	"math"
	"slices"

	f "github.com/booleworks/logicng-go/formula"
)

func encodePBCBinaryMerge(result Result, lts []f.Literal, cffs []int, rhs int, config *Config) {
	fac := result.Factory()
	lits := make([]f.Literal, len(lts))
	copy(lits, lts)
	coeffs := make([]int, len(cffs))
	copy(coeffs, cffs)
	maxWeight := slices.Max(cffs)
	nullLit := fac.Lit(nullLit, true)

	if !config.BinaryMergeUseGAC {
		binaryMerge(result, lts, cffs, rhs, maxWeight, len(lits), nullLit, config, nullLit)
	} else {
		var x f.Literal
		encodeCompleteConstraint := false
		for i := 0; i < len(lits); i++ {
			f1 := math.Floor(math.Log(float64(coeffs[i]) / math.Log(2)))
			f2 := math.Log(float64(coeffs[i])) / math.Log(2)
			condition := math.Abs(f1-f2) <= 1e-9
			if config.BinaryMergeNoSupportForSingleBit && condition {
				encodeCompleteConstraint = true
				continue
			}
			tmpLit := lits[i]
			tmpCoeff := coeffs[i]
			lits[i] = lits[len(lits)-1]
			coeffs[i] = coeffs[len(coeffs)-1]
			lits = lits[:len(lits)-1]
			coeffs = coeffs[:len(coeffs)-1]
			x = tmpLit
			if maxWeight == tmpCoeff {
				mw := slices.Max(coeffs)
				if rhs-tmpCoeff <= 0 {
					for j := 0; j < len(lits); j++ {
						result.AddClause(x.Negate(fac), lits[j].Negate(fac))
					}
				} else {
					binaryMerge(result, lits, coeffs, rhs-tmpCoeff, mw, len(lits), x.Negate(fac), config, nullLit)
				}
			} else {
				if rhs-tmpCoeff <= 0 {
					for j := 0; j < len(lits); j++ {
						result.AddClause(x.Negate(fac), lits[j].Negate(fac))
					}
				}
				binaryMerge(result, lits, coeffs, rhs-tmpCoeff, maxWeight, len(lits), x.Negate(fac), config, nullLit)
			}
			if i < len(lits) {
				lits = append(lits, lits[i])
				lits[i] = tmpLit
				coeffs = append(coeffs, coeffs[i])
				coeffs[i] = tmpCoeff
			}
		}
		if config.BinaryMergeNoSupportForSingleBit && encodeCompleteConstraint {
			binaryMerge(result, lts, cffs, rhs, maxWeight, len(lits), nullLit, config, nullLit)
		}
	}
}

func binaryMerge(
	result Result,
	literals []f.Literal,
	coefficients []int,
	leq, maxWeight, n int,
	gacLit f.Literal,
	config *Config,
	nullLit f.Literal,
) {
	fac := result.Factory()
	lessThen := leq + 1
	p := int(math.Floor(math.Log(float64(maxWeight)) / math.Log(2)))
	m := int(math.Ceil(float64(lessThen) / math.Pow(2, float64(p))))
	newLessThen := int(float64(m) * math.Pow(2, float64(p)))
	t := int(float64(m)*math.Pow(2, float64(p)) - float64(lessThen))

	trueLit := result.NewAuxVar(f.AuxPBC).AsLiteral()
	result.AddClause(trueLit)
	buckets := make([][]f.Literal, p+1)
	bit := 1
	for i := 0; i <= p; i++ {
		buckets[i] = make([]f.Literal, 0, 4)
		if (t & bit) != 0 {
			buckets[i] = append(buckets[i], trueLit)
		}
		for j := range n {
			if (coefficients[j] & bit) != 0 {
				if gacLit != nullLit && coefficients[j] >= lessThen {
					result.AddClause(gacLit, literals[j].Negate(fac))
				} else {
					buckets[i] = append(buckets[i], literals[j])
				}
			}
		}
		bit = bit << 1
	}
	bucketCard := make([][]f.Literal, p+1)
	bucketMerge := make([][]f.Literal, p+1)
	for i := 0; i < p+1; i++ {
		bucketCard[i] = make([]f.Literal, 0, 4)
		bucketMerge[i] = make([]f.Literal, 0, 4)
	}
	carries := make([]f.Literal, 0, 4)
	for i := range buckets {
		k := int(math.Ceil(float64(newLessThen) / math.Pow(2, float64(i))))
		if config.BinaryMergeUseWatchDog {
			totalizer(result, &buckets[i], &bucketCard[i])
		} else {
			sort(k, &buckets[i], &bucketCard[i], result, inputToOutput)
		}
		if k <= len(buckets[i]) {
			if gacLit != nullLit {
				result.AddClause(gacLit, bucketCard[i][k-1].Negate(fac))
			} else {
				result.AddClause(bucketCard[i][k-1].Negate(fac))
			}
		}
		if i > 0 {
			if len(carries) > 0 {
				if len(bucketCard[i]) == 0 {
					bucketMerge[i] = carries
				} else {
					if config.BinaryMergeUseWatchDog {
						unaryAdder(result, &bucketCard[i], &carries, &bucketMerge[i])
					} else {
						merge(k, &bucketCard[i], &carries, &bucketMerge[i], result, inputToOutput)
					}
					if k == len(bucketMerge[i]) || config.BinaryMergeUseWatchDog && k <= len(bucketMerge[i]) {
						if gacLit != nullLit {
							result.AddClause(gacLit, bucketMerge[i][k-1].Negate(fac))
						} else {
							result.AddClause(bucketMerge[i][k-1].Negate(fac))
						}
					}
				}
			} else {
				bucketMerge[i] = bucketCard[i]
			}
		}
		carries = make([]f.Literal, 0, 4)
		if i == 0 {
			for j := 1; j < len(bucketCard[0]); j = j + 2 {
				carries = append(carries, bucketCard[0][j])
			}
		} else {
			for j := 1; j < len(bucketMerge[i]); j = j + 2 {
				carries = append(carries, bucketMerge[i][j])
			}
		}
	}
}

func totalizer(result Result, x, ux *[]f.Literal) {
	*ux = make([]f.Literal, 0, 4)
	if len(*x) == 0 {
		return
	}
	if len(*x) == 1 {
		*ux = append(*ux, (*x)[0])
	} else {
		for i := 0; i < len(*x); i++ {
			*ux = append(*ux, result.NewAuxVar(f.AuxPBC).AsLiteral())
		}
		x1 := make([]f.Literal, len(*x)/2)
		x2 := make([]f.Literal, len(*x)-(len(*x)/2))

		i := 0
		for ; i < len(*x)/2; i++ {
			x1[i] = (*x)[i]
		}
		for ; i < len(*x); i++ {
			x2[i-(len(*x)/2)] = (*x)[i]
		}
		ux1 := make([]f.Literal, 0, 4)
		ux2 := make([]f.Literal, 0, 4)
		totalizer(result, &x1, &ux1)
		totalizer(result, &x2, &ux2)
		unaryAdder(result, &ux1, &ux2, ux)
	}
}

func unaryAdder(result Result, u, v, w *[]f.Literal) {
	fac := result.Factory()
	*w = make([]f.Literal, 0, 4)
	if len(*u) == 0 {
		for i := 0; i < len(*v); i++ {
			*w = append(*w, (*v)[i])
		}
	} else if len(*v) == 0 {
		for i := 0; i < len(*u); i++ {
			*w = append(*w, (*u)[i])
		}
	} else {
		for i := 0; i < len(*u)+len(*v); i++ {
			*w = append(*w, result.NewAuxVar(f.AuxPBC).AsLiteral())
		}
		for a := 0; a < len(*u); a++ {
			for b := 0; b < len(*v); b++ {
				result.AddClause((*u)[a].Negate(fac), (*v)[b].Negate(fac), (*w)[a+b+1])
			}
		}
		for i := 0; i < len(*v); i++ {
			result.AddClause((*v)[i].Negate(fac), (*w)[i])
		}
		for i := 0; i < len(*u); i++ {
			result.AddClause((*u)[i].Negate(fac), (*w)[i])
		}
	}
}
