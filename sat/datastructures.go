package sat

import (
	"fmt"
	"sort"
	"strings"

	f "github.com/booleworks/logicng-go/formula"
)

type clause struct {
	data           []int32
	learntOnState  int32
	isAtMost       bool
	activity       float64
	seen           bool
	lbd            int
	canBeDel       bool
	oneWatched     bool
	atMostWatchers int
}

func newClause(ps []int32, learntOnState int32) *clause {
	return &clause{
		data:           ps,
		learntOnState:  learntOnState,
		canBeDel:       true,
		atMostWatchers: -1,
	}
}

func newAtMostClause(ps []int32, learntOnState int32) *clause {
	return &clause{
		data:           ps,
		learntOnState:  learntOnState,
		isAtMost:       true,
		canBeDel:       true,
		atMostWatchers: -1,
	}
}

func (c *clause) get(i int) int32 {
	return c.data[i]
}

func (c *clause) set(i int, lit int32) {
	c.data[i] = lit
}

func (c *clause) learnt() bool {
	return c.learntOnState >= 0
}

func (c *clause) incrementActivity(inc float64) {
	c.activity += inc
}

func (c *clause) rescaleActivity() {
	c.activity *= 1e-20
}

func (c *clause) pop() {
	c.data = c.data[:len(c.data)-1]
}

func (c *clause) cardinality() int {
	return len(c.data) - c.atMostWatchers + 1
}

func (c *clause) size() int {
	return len(c.data)
}

func (c *clause) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	sb.WriteString(fmt.Sprintf("activity=%f,", c.activity))
	sb.WriteString(fmt.Sprintf("learntOnState=%d,", c.learntOnState))
	sb.WriteString(fmt.Sprintf("seen=%t,", c.seen))
	sb.WriteString(fmt.Sprintf("lbd=%d,", c.lbd))
	sb.WriteString(fmt.Sprintf("canBeDel=%t,", c.canBeDel))
	sb.WriteString(fmt.Sprintf("oneWatched=%t,", c.oneWatched))
	sb.WriteString(fmt.Sprintf("isAtMost=%t,", c.isAtMost))
	sb.WriteString(fmt.Sprintf("atMostWatchers=%d,", c.atMostWatchers))
	sb.WriteString("lits=[")
	for i := 0; i < c.size(); i++ {
		lit := c.data[i]
		if (lit & 1) == 1 {
			sb.WriteString("-")
		}
		sb.WriteString(fmt.Sprintf("%d", lit>>1))
		if i != c.size()-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]}")
	return sb.String()
}

func sortClauses(cs *[]*clause) {
	sort.Slice(*cs, func(i, j int) bool {
		x, y := (*cs)[i], (*cs)[j]
		if len(x.data) > 2 && len(y.data) == 2 {
			return true
		}
		if len(y.data) > 2 && len(x.data) == 2 {
			return false
		}
		if len(x.data) == 2 && len(y.data) == 2 {
			return false
		}
		if x.lbd > y.lbd {
			return true
		}
		if x.lbd < y.lbd {
			return false
		}
		return x.activity < y.activity
	})
}

// Variable
type variable struct {
	assignment f.Tristate
	level      int
	reason     *clause
	activity   float64
	polarity   bool
	decision   bool
}

func newVariable(polarity bool) *variable {
	return &variable{
		assignment: f.TristateUndef,
		level:      -1,
		reason:     nil,
		activity:   0,
		polarity:   polarity,
		decision:   false,
	}
}

func (v *variable) rescaleActivity() {
	v.activity *= 1e-100
}

func (v *variable) incrementActivity(inc float64) {
	v.activity += inc
}

func (v *variable) String() string {
	return fmt.Sprintf("{assignment=%s, level=%d, reason=%s, activity=%f, polarity=%t, decision=%t}",
		v.assignment, v.level, v.reason, v.activity, v.polarity, v.decision)
}

// Watcher
type watcher struct {
	clause  *clause
	blocker int32
}

func newWatcher(clause *clause, blocker int32) *watcher {
	return &watcher{clause, blocker}
}

func (w *watcher) String() string {
	return fmt.Sprintf("{clause=%s, blocker=%d}", w.clause, w.blocker)
}
