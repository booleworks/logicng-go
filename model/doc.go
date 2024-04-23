// Package model provides a data structure for Boolean models in LogicNG.  A
// model is just a list of literals.
//
// In contrast to the assignment found in the [assignment] package, a model
// stores the literals internally just as a list without checking for
// duplicates or contrary assignments.  This makes their generation more
// performant.  Models usually are used by algorithms (like model enumeration,
// or getting a model from a solver) which guarantee that the stored literals
// are unique and contradiction-free.
package model
