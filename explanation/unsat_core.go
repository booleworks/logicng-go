package explanation

import f "github.com/booleworks/logicng-go/formula"

// UnsatCore represents an unsatisfiable core of a formula.  If the core is
// guaranteed to be a MUS, the flag IsGuaranteedMUS is set to true.  If it is
// set to false, the core could be a MUS, but for efficiency reasons this is
// not computed.
type UnsatCore struct {
	Propositions    []f.Proposition // propositions of the unsat core
	IsGuaranteedMUS bool            // flag whether the core is guaranteed to be a MUS
}

// NewUnsatCore creates a new UnsatCore with the given propositions and the
// given flag whether the core is guaranteed to be a MUS.
func NewUnsatCore(propositions []f.Proposition, isMus bool) *UnsatCore {
	return &UnsatCore{propositions, isMus}
}
