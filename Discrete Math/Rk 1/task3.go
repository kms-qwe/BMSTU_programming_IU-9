package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type State struct {
	isFinal     bool
	transitions []int
}

type DFA struct {
	states []State
	start  int
}

type Pair struct {
	first, second int
}

type QueueItem struct {
	pair Pair
	path string
}

func BFS(dfa1, dfa2 DFA, M int) string {
	queue := []QueueItem{{Pair{dfa1.start, dfa2.start}, ""}}
	visited := make(map[Pair]bool)
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.pair] {
			continue
		}
		visited[current.pair] = true

		state1 := dfa1.states[current.pair.first]
		state2 := dfa2.states[current.pair.second]

		if state1.isFinal != state2.isFinal {
			return current.path
		}

		for i := 0; i < M; i++ {
			next := Pair{
				state1.transitions[i],
				state2.transitions[i],
			}

			if !visited[next] {
				queue = append(queue, QueueItem{next, current.path + string(alphabet[i])})
			}
		}
	}
	return "="
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	line, _ := reader.ReadString('\n')
	M, _ := strconv.Atoi(strings.TrimSpace(line))

	line, _ = reader.ReadString('\n')
	parts := strings.Fields(line)
	N1, _ := strconv.Atoi(parts[0])
	start1, _ := strconv.Atoi(parts[1])
	dfa1 := DFA{
		states: make([]State, N1),
		start:  start1,
	}
	for i := 0; i < N1; i++ {
		line, _ = reader.ReadString('\n')
		parts = strings.Fields(line)
		isFinal := parts[0] == "+"
		transitions := make([]int, M)
		for j := 0; j < M; j++ {
			transitions[j], _ = strconv.Atoi(parts[j+1])
		}
		dfa1.states[i] = State{isFinal, transitions}
	}

	line, _ = reader.ReadString('\n')
	parts = strings.Fields(line)
	N2, _ := strconv.Atoi(parts[0])
	start2, _ := strconv.Atoi(parts[1])
	dfa2 := DFA{
		states: make([]State, N2),
		start:  start2,
	}
	for i := 0; i < N2; i++ {
		line, _ = reader.ReadString('\n')
		parts = strings.Fields(line)
		isFinal := parts[0] == "+"
		transitions := make([]int, M)
		for j := 0; j < M; j++ {
			transitions[j], _ = strconv.Atoi(parts[j+1])
		}
		dfa2.states[i] = State{isFinal, transitions}
	}

	result := BFS(dfa1, dfa2, M)
	if result == "=" {
		fmt.Println(result)
	} else {
		fmt.Println(len(result))
	}
}
