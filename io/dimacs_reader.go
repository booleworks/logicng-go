package io

import (
	"bufio"
	"os"
	"strings"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// ReadDimacs reads a CNF from the given filename in Dimacs format.  The
// optional prefix parameter is used to generate the variable names.  The
// default value is `v` therefore variable v1, v2, ... will be generated from
// the input problem.  Returns the CNF as a list of clauses and an optional
// error if there was a problem reading the file.
func ReadDimacs(fac f.Factory, filename string, prefix ...string) (*[]f.Formula, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ReadDimacsFile(fac, file, prefix...)
}

// ReadDimacsFile reads a CNF from the given file in Dimacs format.  The
// optional prefix parameter is used to generate the variable names.  The
// default value is `v` therefore variable v1, v2, ... will be generated from
// the input problem.  Returns the CNF as a list of clauses and an optional
// error if there was a problem reading the file.
func ReadDimacsFile(fac f.Factory, file *os.File, prefix ...string) (*[]f.Formula, error) {
	result := make([]f.Formula, 0, 64)
	var pfx string
	if len(prefix) == 0 {
		pfx = "v"
	} else {
		pfx = prefix[0]
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "c") && !strings.HasPrefix(line, "p") && strings.TrimSpace(line) != "" {
			split := strings.Fields(line)
			if split[len(split)-1] != "0" {
				return nil, errorx.BadInput("line %s did not end with 0", line)
			}
			vars := make([]f.Formula, 0, len(split)-1)
			for _, lit := range split[:len(split)-1] {
				if lit != "" {
					if strings.HasPrefix(lit, "-") {
						vars = append(vars, fac.Literal(pfx+lit[1:], false))
					} else {
						vars = append(vars, fac.Variable(pfx+lit))
					}
				}
			}
			result = append(result, fac.Or(vars...))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &result, nil
}
