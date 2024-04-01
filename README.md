<a href="https://www.logicng.org"><img src="https://github.com/booleworks/logicng-go/blob/main/doc/logos/logicng_logo_gopher.png?raw=true" alt="logo" width="400"></a>

[![license](https://img.shields.io/badge/license-MIT-blue?style=flat-square)]()

# A Library for Creating, Manipulating, and Solving Boolean Formulas

__THIS IS AN ALPHA VERSION! THE API MAY STILL CHANGE IN SIGNIFICANT WAYS! DO NOT USE IN PRODUCTION!__

## Introduction

This is [LogicNG](https://logicng.org/) for Go. The [original
version](https://github.com/logic-ng/LogicNG) of LogicNG is a Java Library for
creating, manipulating and solving Boolean and Pseudo-Boolean formulas.

Its main focus lies on memory-efficient data-structures for Boolean formulas
and efficient algorithms for manipulating and solving them. The library is
designed and most notably used in industrial systems which have to manipulate
and solve millions of formulas per day. The Java version of LogicNG is heavily
used by the German automotive industry to validate and optimize their product
documentation, support the configuration process of vehicles, and compute WLTP
values for emission and consumption.

## Implemented Algorithms

The Go version of LogicNG currently provides the following key functionalities:

- Support for Boolean formulas, cardinality constraints, and pseudo-Boolean
  formulas
- Parsing of Boolean formulas from strings or files
- Transformations of formulas, like
  - Normal-Forms NNF, DNF, or CNF with various configurable algorithms
  - Substitution in formulas
  - Subsumption
  - Simplification of formulas
- Encoding cardinality constraints and pseudo-Boolean formulas to purely
  Boolean formulas with a multitude of different algorithms
- Solving formulas with an integrated SAT Solver including
  - Fast backbone computation on the solver
  - Incremental/Decremental solver interface
  - Proof generation
  - Optimization with incremental cardinality constraints
  - Fast model and projected model enumeration
- Optimizing formulas with an integrated MaxSAT solver
- Knowledge compilation with BDDs or DNNFs
- Computation of minimum prime implicants and prime implicant covers
- and many more...

## Philosophy

The most important philosophy of the library is to avoid unnecessary object
creation. Therefore, formulas can only be generated via formula factories. A
formula factory assures that a formula is only created once in memory. If
another instance of the same formula is created by the user, the already
existing one is returned by the factory. This leads to a small memory footprint
and fast execution of algorithms. Formulas can cache the results of algorithms
executed on them and since every formula is hold only once in memory it is
assured that the same algorithm on the same formula is also executed only once.

## Whitepaper

If you want a high-level overview of LogicNG and how it is used in many
applications in the area of product configuration, you can read our
[whitepaper](https://logicng.org/whitepaper/abstract/).

## First Steps

The following code creates the Boolean Formula _A and not (B or not C)_
programmatically.

```go
import "booleworks.com/logicng/formula"

fac := formula.NewFactory()
a := fac.Variable("A")
b := fac.Variable("B")
c := fac.Variable("C")
form := fac.And(a, fac.Or(b, fac.Not(c)))
```

Alternatively you can just parse the formula from a string:

```go
import (
    "booleworks.com/logicng/formula"
    "booleworks.com/logicng/parser" 
)

fac := formula.NewFactory()
parser := parser.New(fac)
form, err := parser.Parse("A & (B | ~C)")
```

Once you created the formula you can for example convert it to NNF or CNF or
solve it with a SAT solver:

```go
import (
    "fmt"
    "booleworks.com/logicng/formula"
    "booleworks.com/logicng/normalform"
    "booleworks.com/logicng/parser"
    "booleworks.com/logicng/sat"
)

fac := formula.NewFactory()
parser := parser.New(fac)
form, err := parser.Parse("A & ~(B | ~C)")

nnf := normalform.NNF(fac, form)
cnf := normalform.CNF(fac, form)

fmt.Println(cnf.Sprint(fac)) // pretty-print the formula

solver := sat.NewSolver(fac)
solver.Add(form)
result := Solver.Sat() // is true
```

