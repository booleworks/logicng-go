package simplification

import f "github.com/booleworks/logicng-go/formula"

// A RatingFunction is used during simplification to compute the simplicity
// rating of a formula.  This rating is then used to compare it during
// simplification and choosing the formula with the lowest rating.
type RatingFunction func(fac f.Factory, formula f.Formula) float64

// The DefaultRatingFunction rates a formula by its string length.
func DefaultRatingFunction(fac f.Factory, formula f.Formula) float64 {
	return float64(len(formula.Sprint(fac)))
}
