// Package errorx contains LogicNG specific extensions to error types.
//
// When it comes to error handling,  ee try to adhere to the following
// guidelines:
//   - If a problem is clearly a programming error on the caller's side and is
//     deterministically reproducible, we usually panic.  Examples for this are
//     mixing formulas which were not produced by the same formula factory or
//     conjoining BDDs with a different kernel.  Handling these situations as
//     errors would make the API very unpleasant to work with.
//   - When failures can happen in a typical workflow using the library, we
//     usually use errors.  Examples for this are adding a literal to an
//     assignment which is already present with the different polarity or
//     constructing an incremental cardinality constraint from a tautological or
//     contradictory constraint.
package errorx
