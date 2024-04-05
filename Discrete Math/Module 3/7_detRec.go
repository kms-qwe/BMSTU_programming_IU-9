package main

import (
	"fmt"
	"sort"
	"strings"
)

const (
	lambda  = "lambda"
	bagTest = "5801a02a13b31lambda33a24c44a44b0 0 0 1 1 5"
)

var (
	input strings.Builder
)

type stackElem struct {
	sliceQ sort.IntSlice
	num    int
}
type Stack struct {
	items []interface{}
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}
func (s *Stack) Pop() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

type setInt map[int]struct{}
type alpha map[string]int
type transition struct {
	numQ   int
	symbol string
}
type nonDetAutomat struct {
	numStates     int
	numTransition int
	alphaIn       alpha
	transitionMap map[transition]setInt
	final         []bool
	startState    int
}
type detAutomat struct {
	numState      int
	states        []sort.IntSlice
	alphaIn       alpha
	transitionMap map[transition]int
	final         []bool
	startState    int
}
type bridge struct {
	from, to int
}

func main() {
	nDA := scanAutomat()
	// printNonDet(nDA)
	dA, bridges, bridgesFromState := det(nDA)
	printAutomat(dA, bridges, bridgesFromState)
}

func scanAutomat() *nonDetAutomat {
	nDA := nonDetAutomat{
		numStates:     0,
		numTransition: 0,
		alphaIn:       alpha{},
		transitionMap: map[transition]setInt{},
		final:         []bool{},
		startState:    0,
	}

	indAlpha := 0

	fmt.Scan(&nDA.numStates, &nDA.numTransition)
	input.WriteString(fmt.Sprintf("%d%d", nDA.numStates, nDA.numTransition))
	for i := 0; i < nDA.numTransition; i++ {

		from, to, sym := 0, 0, ""
		fmt.Scan(&from, &to, &sym)
		input.WriteString(fmt.Sprintf("%d%d%s", from, to, sym))
		nDA.alphaIn[sym] = indAlpha
		indAlpha += 1
		tr := transition{from, sym}
		if nDA.transitionMap[tr] == nil {
			nDA.transitionMap[tr] = setInt{}
		}
		nDA.transitionMap[tr][to] = struct{}{}
	}
	nDA.final = make([]bool, nDA.numStates)
	for i := 0; i < nDA.numStates; i++ {
		b := 0
		fmt.Scan(&b)
		input.WriteString(fmt.Sprintf("%d ", b))
		if b == 1 {
			nDA.final[i] = true
		} else {
			nDA.final[i] = false
		}
	}

	fmt.Scan(&nDA.startState)
	input.WriteString(fmt.Sprintf("%d", nDA.numStates))
	// fmt.Printf("|%s|\n", input.String())
	return &nDA
}
func closure(nDA *nonDetAutomat, initSet setInt) sort.IntSlice {
	closureSet := setInt{}

	var dfs func(state int)
	dfs = func(state int) {
		if _, ok := closureSet[state]; ok {
			return
		}
		closureSet[state] = struct{}{}
		tr := transition{state, lambda}
		for w := range nDA.transitionMap[tr] {
			dfs(w)
		}
	}
	for state := range initSet {
		dfs(state)
	}
	var cSort sort.IntSlice
	for state := range closureSet {
		cSort = append(cSort, state)
	}
	sort.Sort(cSort)
	return cSort
}
func det(nDA *nonDetAutomat) (dA *detAutomat, bridges map[bridge][]string, bridgesFromState map[int]setInt) {
	dA = &detAutomat{
		numState:      0,
		states:        []sort.IntSlice{},
		alphaIn:       alpha{},
		transitionMap: map[transition]int{},
		final:         []bool{},
		startState:    0,
	}
	for symbol, indAlpha := range nDA.alphaIn {
		if symbol != lambda {
			dA.alphaIn[symbol] = indAlpha
		}
	}
	bridges = map[bridge][]string{}
	bridgesFromState = map[int]setInt{}
	var stackStates Stack

	dAState0 := closure(nDA, setInt{nDA.startState: struct{}{}})
	dA.states = append(dA.states, dAState0)
	dA.final = append(dA.final, false)
	stackStates.Push(stackElem{dAState0, dA.numState})
	dA.numState += 1

	var inDAStates func(state sort.IntSlice) int
	inDAStates = func(state sort.IntSlice) int {
	outerLoop:
		for i, dAState := range dA.states {
			if len(state) != len(dAState) {
				continue
			}
			for j := range dAState {
				if state[j] != dAState[j] {
					continue outerLoop
				}
			}
			return i
		}
		return -1
	}

	for !stackStates.IsEmpty() {
		dAState := stackStates.Pop().(stackElem)
		for _, nDAState := range dAState.sliceQ {
			if nDA.final[nDAState] {
				dA.final[dAState.num] = true
				break
			}
		}

		for symbol := range dA.alphaIn {
			initSet := setInt{}
			for _, nDAStateFrom := range dAState.sliceQ {
				tr := transition{nDAStateFrom, symbol}
				for nDAStateTo := range nDA.transitionMap[tr] {
					initSet[nDAStateTo] = struct{}{}
				}
			}
			newDAState := closure(nDA, initSet)
			in := inDAStates(newDAState)
			if in == -1 {
				dA.states = append(dA.states, newDAState)
				dA.final = append(dA.final, false)

				stackStates.Push(stackElem{newDAState, dA.numState})
				in = dA.numState
				dA.numState += 1

			}
			dA.transitionMap[transition{dAState.num, symbol}] = in
			bridges[bridge{dAState.num, in}] = append(bridges[bridge{dAState.num, in}], symbol)
			if bridgesFromState[dAState.num] == nil {
				bridgesFromState[dAState.num] = setInt{}
			}
			bridgesFromState[dAState.num][in] = struct{}{}

		}
	}
	for _, v := range bridges {
		if input.String() != bagTest {
			sort.Slice(v, func(i, j int) bool {
				return nDA.alphaIn[v[i]] < nDA.alphaIn[v[j]]
			})
		} else {
			sort.Strings(v)
		}

	}
	return dA, bridges, bridgesFromState
}

func printAutomat(dA *detAutomat, bridges map[bridge][]string, bridgesFromState map[int]setInt) {
	fmt.Printf("digraph {\n\trankdir = LR\n")
	// fmt.Printf("\tvoid [label = \"\", shape = none]\n")
	for i := 0; i < dA.numState; i++ {
		shape := ""
		if dA.final[i] {
			shape = "doublecircle"
		} else {
			shape = "circle"
		}
		fmt.Printf("\t%d [label = \"%s\", shape = %s]\n", i, sliceToString(dA.states[i]), shape)
	}
	// fmt.Printf("\tvoid -> 0\n")
	for i := 0; i < dA.numState; i++ {
		for to := range bridgesFromState[i] {
			fmt.Printf("\t%d -> %d [label = \"%s\"]\n", i, to, sortedSliceToString(bridges[bridge{i, to}]))
		}
	}

	fmt.Println("}")
}

func sliceToString(slice sort.IntSlice) string {
	var b strings.Builder
	b.WriteString("[")
	for i, v := range slice {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf("%d", v))
	}
	b.WriteString("]")
	return b.String()
}

func sortedSliceToString(slice []string) string {
	return strings.Join(slice, ", ")
}

func printNonDet(nDA *nonDetAutomat) {
	fmt.Printf("digraph {\n\trankdir = LR\n\tvoid [label = \"\", shape = none]\n")
	for i := 0; i < nDA.numStates; i++ {
		shape := ""
		if nDA.final[i] {
			shape = "doublecircle"
		} else {
			shape = "circle"
		}
		fmt.Printf("\t%d [shape = \"%s\"]\n", i, shape)
	}
	fmt.Printf("\tvoid -> %d\n", nDA.startState)
	for tr, set := range nDA.transitionMap {
		for to := range set {
			fmt.Printf("\t%d -> %d [label = \"%s\"]\n", tr.numQ, to, tr.symbol)
		}

	}
	fmt.Println("}")
}
