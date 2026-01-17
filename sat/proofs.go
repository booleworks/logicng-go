package sat

import (
	"math"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

type proofInformation struct {
	clause      []int32
	proposition f.Proposition
}

const (
	proofBigInit = 1000000
	proofUnsat   = 0
	proofSat     = 1
	proofExtra   = 2
	proofMark    = 3
)

type drupResult struct {
	trivialUnsat bool
	unsatCore    [][]int32
}

func drupCompute(originalProblem, proof *[][]int32) drupResult {
	result := drupResult{}
	solver := newDrupSolver(originalProblem, proof)
	parseReturnValue := solver.parse()
	if !parseReturnValue {
		result.trivialUnsat = true
		result.unsatCore = [][]int32{}
	} else {
		result.trivialUnsat = false
		result.unsatCore = solver.verify()
	}
	return result
}

type drupSolver struct {
	originalProblem [][]int32
	proof           [][]int32
	core            [][]int32
	delete          bool
	db              []int
	nVars           int
	nClauses        int
	falseStack      []int
	reason          []int
	internalFalse   []int
	forcedPtr       int
	processedPtr    int
	assignedPtr     int
	adlist          []int
	wlist           [][]int
	count           int
	adlemmas        int
	lemmas          int
	time            int
}

func newDrupSolver(originalProblem, proof *[][]int32) *drupSolver {
	solver := &drupSolver{}

	solver.originalProblem = *originalProblem
	solver.proof = *proof
	solver.core = [][]int32{}
	solver.delete = true
	return solver
}

func (s *drupSolver) parse() bool {
	s.nVars = 0
	for _, slice := range s.originalProblem {
		for i := range slice {
			abs := int(math.Abs(float64(slice[i])))
			if abs > s.nVars {
				s.nVars = abs
			}
		}
	}
	nClauses := len(s.originalProblem)

	del := false
	nZeros := nClauses
	var buffer []int
	s.db = []int{}

	s.count = 1
	s.falseStack = make([]int, s.nVars+1)
	s.reason = make([]int, s.nVars+1)
	s.internalFalse = make([]int, 2*s.nVars+3)
	s.wlist = make([][]int, 2*s.nVars+3)
	for i := 1; i <= s.nVars; i++ {
		s.wlist[drupIndex(i)] = []int{}
		s.wlist[drupIndex(-i)] = []int{}
	}
	s.adlist = []int{}
	marks := make([]int, 2*s.nVars+3)
	mark := 0

	hashTable := make(map[int32]*[]int)
	currentFile := s.originalProblem
	var fileSwitchFlag bool
	clauseNr := 0
	for {
		fileSwitchFlag = nZeros <= 0
		var clause []int32
		if clauseNr < len(currentFile) {
			clause = currentFile[clauseNr]
		}
		clauseNr++
		if clause == nil {
			s.lemmas = len(s.db) + 1
			break
		}
		toks := make([]int, 0, len(clause)-1)
		if fileSwitchFlag && clause[0] == -1 {
			del = true
		}
		var i int
		if fileSwitchFlag {
			i = 1
		} else {
			i = 0
		}
		for ; i < len(clause); i++ {
			toks = append(toks, int(clause[i]))
		}
		buffer = append(buffer, toks...)
		if clauseNr >= len(currentFile) && !fileSwitchFlag {
			fileSwitchFlag = true
			clauseNr = 0
			currentFile = s.proof
		}
		if clauseNr > len(currentFile) && fileSwitchFlag && len(currentFile) > 0 {
			break
		}
		mark++
		hash := drupGetHash(&marks, mark, buffer)
		if del {
			if s.delete {
				match := s.matchClause(hashTable[hash], marks, mark, len(buffer))
				pop(hashTable[hash])
				s.adlist = append(s.adlist, (match<<1)+1)
			}
			del = false
			buffer = []int{}
			continue
		}
		clausePtr := len(s.db) + 1
		s.db = append(s.db, 2*s.count)
		s.count++
		for i := 0; i < len(buffer); i++ {
			s.db = append(s.db, buffer[i])
		}
		s.db = append(s.db, 0)

		slice, ok := hashTable[hash]
		if !ok {
			slice = &[]int{}
			hashTable[hash] = slice
		}
		*slice = append(*slice, clausePtr)
		s.adlist = append(s.adlist, clausePtr<<1)

		if nZeros == 0 {
			s.lemmas = clausePtr
			s.adlemmas = len(s.adlist) - 1
		}
		if nZeros > 0 {
			if len(buffer) == 0 || ((len(buffer) == 1) && s.internalFalse[drupIndex(s.db[clausePtr])] != 0) {
				return false
			} else if len(buffer) == 1 {
				if s.internalFalse[drupIndex(-s.db[clausePtr])] == 0 {
					pos := int(math.Abs(float64(s.db[clausePtr])))
					s.reason[pos] = clausePtr + 1
					s.assign(s.db[clausePtr])
				}
			} else {
				s.addWatch(clausePtr, 0)
				s.addWatch(clausePtr, 1)
			}
		} else if len(buffer) == 0 {
			break
		}
		buffer = []int{}
		nZeros--
	}
	return true
}

func (s *drupSolver) matchClause(clauselist *[]int, marks []int, mark int, inputSize int) int {
	var i, matchsize int
	for i = 0; i < len(*clauselist); i++ {
		matchsize = 0
		aborted := false
		for l := (*clauselist)[i]; s.db[l] != 0; l++ {
			if marks[drupIndex(s.db[l])] != mark {
				aborted = true
				break
			}
			matchsize++
		}
		if !aborted && inputSize == matchsize {
			result := (*clauselist)[i]
			(*clauselist)[i] = (*clauselist)[len(*clauselist)-1]
			return result
		}
	}
	panic(errorx.IllegalState("could not match deleted clause"))
}

func (s *drupSolver) assign(a int) {
	s.internalFalse[drupIndex(-a)] = 1
	s.falseStack[s.assignedPtr] = -a
	s.assignedPtr++
}

func (s *drupSolver) addWatch(cPtr, index int) {
	lit := s.db[cPtr+index]
	s.wlist[drupIndex(lit)] = append(s.wlist[drupIndex(lit)], cPtr<<1)
}

func (s *drupSolver) verify() [][]int32 {
	var ad, d int
	flag := false
	clausePtr := 0
	lemmasPtr := s.lemmas
	lastPtr := s.lemmas
	endPtr := s.lemmas
	checked := s.adlemmas
	var buffer []int

	gotoPostProcess := false
	if s.processedPtr < s.assignedPtr {
		if s.propagate() == proofUnsat {
			gotoPostProcess = true
		}
	}
	s.forcedPtr = s.processedPtr

	if !gotoPostProcess {
		gotoVerification := false
		for !gotoVerification {
			flag = false
			buffer = []int{}
			clausePtr = lemmasPtr
			for ok := true; ok; ok = d != 0 {
				ad = s.adlist[checked]
				checked++
				d = ad & 1
				cPtr := ad >> 1
				if d != 0 && s.db[cPtr+1] != 0 {
					pos := int(math.Abs(float64(s.db[cPtr])))
					if s.reason[pos]-1 == ad>>1 {
						continue
					}
					s.removeWatch(cPtr, 0)
					s.removeWatch(cPtr, 1)
				}
			}

			for s.db[lemmasPtr] != 0 {
				lit := s.db[lemmasPtr]
				lemmasPtr++
				if s.internalFalse[drupIndex(-lit)] != 0 {
					flag = true
				}
				if s.internalFalse[drupIndex(lit)] == 0 {
					if len(buffer) <= 1 {
						s.db[lemmasPtr-1] = s.db[clausePtr+len(buffer)]
						s.db[clausePtr+len(buffer)] = lit
					}
					buffer = append(buffer, lit)
				}
			}

			if s.db[clausePtr+1] != 0 {
				s.addWatch(clausePtr, 0)
				s.addWatch(clausePtr, 1)
			}
			lemmasPtr += proofExtra

			if flag {
				s.adlist[checked-1] = 0
			}
			if flag {
				continue // Clause is already satisfied
			}
			if len(buffer) == 0 {
				panic(errorx.IllegalState("conflict claimed, but not detected"))
			}

			if len(buffer) == 1 {
				s.assign(buffer[0])
				pos := int(math.Abs(float64(buffer[0])))
				s.reason[pos] = clausePtr + 1
				s.forcedPtr = s.processedPtr
				if s.propagate() == proofUnsat {
					gotoVerification = true
				}
			}

			if lemmasPtr >= len(s.db) {
				break
			}
		}
		if !gotoVerification {
			panic(errorx.IllegalState("no conflit"))
		}

		s.forcedPtr = s.processedPtr
		lemmasPtr = clausePtr - proofExtra

		for {
			buffer = []int{}
			clausePtr = lemmasPtr + proofExtra
			for ok := true; ok; ok = d != 0 {
				checked--
				ad = s.adlist[checked]
				d = ad & 1
				cPtr := ad >> 1
				if d != 0 && s.db[cPtr+1] != 0 {
					pos := int(math.Abs(float64(s.db[cPtr])))
					if s.reason[pos]-1 == ad>>1 {
						continue
					}
					s.addWatch(cPtr, 0)
					s.addWatch(cPtr, 1)
				}
			}

			time := s.db[clausePtr-1]
			if s.db[clausePtr+1] != 0 {
				s.removeWatch(clausePtr, 0)
				s.removeWatch(clausePtr, 1)
			}

			gotoNextLemma := ad == 0
			if !gotoNextLemma {
				for s.db[clausePtr] != 0 {
					lit := s.db[clausePtr]
					clausePtr++
					if s.internalFalse[drupIndex(-lit)] != 0 {
						flag = true
					}
					if s.internalFalse[drupIndex(lit)] == 0 {
						buffer = append(buffer, lit)
					}
				}

				if flag && len(buffer) == 1 {
					for ok := true; ok; ok = s.falseStack[s.forcedPtr] != -buffer[0] {
						s.forcedPtr--
						s.internalFalse[drupIndex(s.falseStack[s.forcedPtr])] = 0
					}
					s.processedPtr = s.forcedPtr
					s.assignedPtr = s.forcedPtr
				}

				if (time & 1) != 0 {
					for i := 0; i < len(buffer); i++ {
						s.assign(-buffer[i])
						pos := int(math.Abs(float64(buffer[i])))
						s.reason[pos] = 0
					}
					if s.propagate() == proofSat {
						panic(errorx.IllegalState("formula is SAT"))
					}
				}
			}

			if lemmasPtr+proofExtra == lastPtr {
				break
			}
			lemmasPtr--
			for s.db[lemmasPtr] != 0 {
				lemmasPtr--
			}
		}
	}

	var marked int
	lemmasPtr = 0
	for lemmasPtr+proofExtra <= lastPtr {
		if (s.db[lemmasPtr] & 1) != 0 {
			s.count++
		}
		lemmasPtr++
		for s.db[lemmasPtr] != 0 {
			lemmasPtr++
		}
		lemmasPtr++
	}
	lemmasPtr = 0

	for lemmasPtr+proofExtra <= lastPtr {
		var coreSlice []int32
		marked = s.db[lemmasPtr] & 1
		lemmasPtr++
		for s.db[lemmasPtr] != 0 {
			if marked != 0 {
				coreSlice = append(coreSlice, int32(s.db[lemmasPtr]))
			}
			lemmasPtr++
		}
		if marked != 0 {
			s.core = append(s.core, coreSlice)
		}
		lemmasPtr++
	}

	s.count = 0
	for lemmasPtr+proofExtra <= endPtr {
		marked = s.db[lemmasPtr] & 1
		lemmasPtr++
		if marked != 0 {
			s.count++
		}
		for s.db[lemmasPtr] != 0 {
			lemmasPtr++
		}
		lemmasPtr++
	}
	return s.core
}

func (s *drupSolver) propagate() int32 {
	start := make([]int, 2)
	check := 0
	var i, lit int
	_lit := 0
	var watch *[]int
	_watchPtr := 0
	start[0] = s.processedPtr
	start[1] = s.processedPtr
	gotoFlipCheck := true
	for gotoFlipCheck {
		gotoFlipCheck = false
		check ^= 1
		for !gotoFlipCheck && start[check] < s.assignedPtr {
			lit = s.falseStack[start[check]]
			start[check]++
			watch = &s.wlist[drupIndex(lit)]
			var watchPtr int
			if lit == _lit {
				watchPtr = _watchPtr
			} else {
				watchPtr = 0
			}

			for watchPtr < len(*watch) {
				if ((*watch)[watchPtr] & 1) != check {
					watchPtr++
					continue
				}
				clausePtr := (*watch)[watchPtr] / 2
				if s.internalFalse[drupIndex(-s.db[clausePtr])] != 0 ||
					s.internalFalse[drupIndex(-s.db[clausePtr+1])] != 0 {
					watchPtr++
					continue
				}
				if s.db[clausePtr] == lit {
					s.db[clausePtr] = s.db[clausePtr+1]
				}
				gotoNextClause := false
				for i = 2; s.db[clausePtr+i] != 0; i++ {
					if s.internalFalse[drupIndex(s.db[clausePtr+i])] == 0 {
						s.db[clausePtr+1] = s.db[clausePtr+i]
						s.db[clausePtr+i] = lit
						s.addWatchLit(s.db[clausePtr+1], (*watch)[watchPtr])
						(*watch)[watchPtr] = s.wlist[drupIndex(lit)][len(s.wlist[drupIndex(lit)])-1]
						pop(&s.wlist[drupIndex(lit)])
						gotoNextClause = true
						break
					}
				}
				if !gotoNextClause {
					s.db[clausePtr+1] = lit
					watchPtr++
					if s.internalFalse[drupIndex(s.db[clausePtr])] == 0 {
						s.assign(s.db[clausePtr])
						pos := int(math.Abs(float64(s.db[clausePtr])))
						s.reason[pos] = clausePtr + 1
						if check == 0 {
							start[0]--
							_lit = lit
							_watchPtr = watchPtr
							gotoFlipCheck = true
							break
						}
					} else {
						s.analyze(clausePtr)
						return proofUnsat
					}
				}
			}
		}
		if check != 0 {
			gotoFlipCheck = true
		}
	}
	s.processedPtr = s.assignedPtr
	return proofSat
}

func (s *drupSolver) removeWatch(cPtr, index int) {
	lit := s.db[cPtr+index]
	watch := &s.wlist[drupIndex(lit)]
	watchPtr := int32(0)
	for {
		_cPtr := (*watch)[watchPtr] >> 1
		watchPtr++
		if _cPtr == cPtr {
			(*watch)[watchPtr-1] = s.wlist[drupIndex(lit)][len(s.wlist[drupIndex(lit)])-1]
			pop(&s.wlist[drupIndex(lit)])
			return
		}
	}
}

func (s *drupSolver) addWatchLit(l, m int) {
	s.wlist[drupIndex(l)] = append(s.wlist[drupIndex(l)], m)
}

func (s *drupSolver) analyze(clausePtr int) {
	s.markClause(clausePtr, 0)
	for s.assignedPtr > 0 {
		s.assignedPtr--
		lit := s.falseStack[s.assignedPtr]
		if (s.internalFalse[drupIndex(lit)] == proofMark) && s.reason[int(math.Abs(float64(lit)))] != 0 {
			s.markClause(s.reason[int(math.Abs(float64(lit)))], -1)
		}
		if s.assignedPtr < s.forcedPtr {
			s.internalFalse[drupIndex(lit)] = 1
		} else {
			s.internalFalse[drupIndex(lit)] = 0
		}
	}
	s.processedPtr = s.forcedPtr
	s.assignedPtr = s.forcedPtr
}

func (s *drupSolver) markClause(clausePtr, index int) {
	if (s.db[clausePtr+index-1] & 1) == 0 {
		s.db[clausePtr+index-1] = s.db[clausePtr+index-1] | 1
		if s.db[clausePtr+1+index] == 0 {
			return
		}
		s.markWatch(clausePtr, index, -index)
		s.markWatch(clausePtr, 1+index, -index)
	}
	for s.db[clausePtr] != 0 {
		s.internalFalse[drupIndex(s.db[clausePtr])] = proofMark
		clausePtr++
	}
}

func (s *drupSolver) markWatch(clausePtr, index, offset int) {
	watch := &s.wlist[drupIndex(s.db[clausePtr+index])]
	clause := s.db[clausePtr-offset-1]
	watchPtr := int32(0)
	for {
		_clause := s.db[((*watch)[watchPtr]>>1)-1]
		watchPtr++
		if _clause == clause {
			(*watch)[watchPtr-1] = (*watch)[watchPtr-1] | 1
			return
		}
	}
}

func drupGetHash(marks *[]int, mark int, input []int) int32 {
	sum := 0
	xor := 0
	prod := 1
	for i := range input {
		prod *= input[i]
		sum += input[i]
		xor ^= input[i]
		(*marks)[drupIndex(input[i])] = mark
	}
	return int32(math.Abs(float64((1023*sum + prod ^ (31 * xor)) % proofBigInit)))
}

func drupIndex(lit int) int {
	if lit > 0 {
		return lit * 2
	}
	return (-lit * 2) ^ 1
}
