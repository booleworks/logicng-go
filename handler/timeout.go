package handler

import (
	"time"

	"github.com/booleworks/logicng-go/event"
)

// A Timeout handler is used to cancel computations after a given time.
type Timeout struct {
	designatedEnd time.Time
}

// NewTimeoutWithEnd generates a new timeout handler which cancels a computation
// at the given time.
func NewTimeoutWithEnd(time time.Time) *Timeout {
	return &Timeout{time}
}

// NewTimeoutWithDuration generates a new timeout handler which cancels a computation
// after the given duration.
func NewTimeoutWithDuration(duration time.Duration) *Timeout {
	designatedEnd := time.Now().Add(duration)
	return &Timeout{designatedEnd}
}

// ShouldResume returns true if the time limit is not yet reached.
func (t Timeout) ShouldResume(event.Event) bool {
	return !time.Now().After(t.designatedEnd)
}
