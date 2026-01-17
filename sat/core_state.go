package sat

import f "github.com/booleworks/logicng-go/formula"

// CoreState bundles the inner state of a SAT solver and is only used for
// serializing and de-serializing solvers.  There should be no need to use it
// in other contexts, and it should *never* be used to manipulate solver state.
type CoreState struct {
	Config                *Config
	LLConfig              *LowLevelConfig
	Ok                    bool
	QHead                 int
	UnitClauses           *[]int32
	Clauses               *[]*clause
	Learnts               *[]*clause
	Watches               *[][]*watcher
	Vars                  *[]*variable
	OrderHeap             *lngheap
	Trail                 *[]int32
	TrailLim              *[]int
	Model                 *[]bool
	Conflict              *[]int32
	Assumptions           *[]int32
	Seen                  *[]bool
	AnalyzeBtLevel        int
	ClaInc                float64
	ClausesLiterals       int
	LearntsLiterals       int
	Name2idx              *map[string]int32
	IDX2name              *map[int32]string
	CanceledByHandler     bool
	AssumptionProps       *[]f.Proposition
	PgOriginalClauses     *[]proofInformation
	PgProof               *[][]int32
	backboneCandidates    *[]int32
	BackboneAssumptions   *[]int32
	BackboneMap           *map[int32]f.Tristate
	ComputingBackbone     bool
	VarDecay              float64
	VarInc                float64
	LearntsizeAdjustConfl float64
	LearntsizeAdjustCnt   int
	LearntsizeAdjustInc   float64
	MaxLearnts            float64
	WatchesBin            *[][]*watcher
	PermDiff              *[]int
	LastDecisionLevel     *[]int32
	LBDQueue              *boundedQueue
	TrailQueue            *boundedQueue
	Myflag                int
	AnalyzeLBD            int
	NbClausesBeforeReduce int
	Conflicts             int
	ConflictsRestarts     int
	SumLBD                float64
	CurRestart            int
	StateId               int32
	ValidStates           *[]int32
	InSatCall             bool
}
