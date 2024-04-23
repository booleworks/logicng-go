// Package primeimplicant provides algorithms for computing minimum prime
// implicants and minimum prime implicant coverages in LogicNG.
//
//   - An implicant of a formula is any min-term – a conjunction of literals –
//     such that the implicant logically implies the formula.
//   - A prime implicant is an implicant which cannot be further reduced (i.e.
//     literals being removed) such that the reduced term yields an implicant.
//   - A minimum-size prime implicant is a prime implicant with minimum size,
//     in terms of the number of literals, among all prime implicants of a
//     formula.
//
// LogicNG provides a function to compute a minimum prime implicant:
//
//	fac := formula.NewFactory()
//	p := parser.New(fac)
//	f1 := p.ParseUnsafe("(A | B) & (A | C ) & (C | D) & (B | ~D)")
//	implicant, err := primimplicant.Minimum(fac, f1) // implicant B, C
//
// A prime implicant cover of a formula is a number of prime implicants
// which cover all min-terms of the formula.  To compute such a cover of
// minimal size, you can use the following code.  Computing minimal prime
// covers is an important step in simplifying algorithms like QuineMcCluskey.
//
//	fac := formula.NewFactory()
//	p := parser.New(fac)
//	f1 := p.ParseUnsafe("(A | B) & (A | C ) & (C | D) & (B | ~D)")
//	primes := primimplicant.CoverMin(fac, f1, primimplicant.CoverImplicants)
//	implicants := primes.Implicants // [B, C], [A, C, ~D], [A, B, D]
package primeimplicant
