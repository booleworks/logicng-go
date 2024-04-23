// Package handler provides datastructures for handlers which can be used to
// abort potentially long-running computations in LogicNG.  There are some
// standard handlers like timeout handler already implemented in LogicNG.  If
// you need other criteria when to abort computations, you can implement the
// [Handler] interface yourself.
package handler
