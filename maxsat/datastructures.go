package maxsat

type problemType byte

const (
	unweighted problemType = iota
	weighted
)

//go:generate stringer -type=problemType

type softClause struct {
	clause         []int32
	relaxationVars []int32
	weight         int
	assumptionVar  int32
}

func newSoftClause(clause, relaxationVars []int32, weight int, assumptionVar int32) *softClause {
	return &softClause{clause, relaxationVars, weight, assumptionVar}
}

type hardClause struct {
	clause []int32
}

func newHardClause(clause []int32) *hardClause {
	return &hardClause{clause}
}
