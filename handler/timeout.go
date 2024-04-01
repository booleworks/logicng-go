package handler

import "time"

// A Timeout handler is used to abort computations after a given time.
type Timeout struct {
	Computation
	designatedEnd time.Time
}

// NewTimeoutWithEnd generates a new timeout handler which aborts a computation
// at the given time.
func NewTimeoutWithEnd(time time.Time) *Timeout {
	return &Timeout{Computation{}, time}
}

// NewTimeoutWithDuration generates a new timeout handler which aborts a computation
// after the given duration.
func NewTimeoutWithDuration(duration time.Duration) *Timeout {
	designatedEnd := time.Now().Add(duration)
	return &Timeout{Computation{}, designatedEnd}
}

// TimeLimitExceeded reports whether the internal time limit of the handler was
// exceeded.
func (t *Timeout) TimeLimitExceeded() bool {
	t.aborted = time.Now().After(t.designatedEnd)
	return t.aborted
}

// Aborted reports whether the computation was aborted by the handler.
func (t *Timeout) Aborted() bool {
	t.TimeLimitExceeded()
	return t.aborted
}
