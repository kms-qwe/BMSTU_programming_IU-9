package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Error struct {
	Message string // Сообщение об ошибке
}

// Функция создания новой ошибки
func NewError(message string) *Error {
	return &Error{
		Message: message,
	}
}

// Метод для получения сообщения об ошибке
func (err *Error) Error() string {
	return err.Message
}

type vertex struct {
	group     int
	neighbors []int
}

func main() {
	n, graph := scanGraph()
	groupN, err := segmentaion(n, graph)
	if err != nil {
		fmt.Println("No solution")
		os.Exit(0)
	}
	groupN = groupN / 2
	bruteForce(groupN, n, graph)
}
func scanGraph() (int, map[int]vertex) {
	n := 0
	fmt.Scan(&n)
	graph := make(map[int]vertex)
	scanner := bufio.NewScanner(os.Stdin)

	for i := 0; i < n; i++ {
		scanner.Scan()
		line := strings.ReplaceAll(scanner.Text(), " ", "")
		graph[i] = vertex{group: -1, neighbors: []int{}}
		for j, signum := range line {
			if signum == rune('+') {
				graph[i] = vertex{-1, append(graph[i].neighbors, j)}
			}
		}
	}
	return n, graph
}
func segmentaion(n int, graph map[int]vertex) (groupNum int, ok error) {
	used := make(map[int]struct{})
	var dfs func(v int, g map[int]vertex) error
	dfs = func(v int, g map[int]vertex) error {
		for _, childVertex := range graph[v].neighbors {
			if _, ok := used[childVertex]; ok && graph[v].group == graph[childVertex].group {
				// fmt.Println("ERROR!:", childVertex, v)
				return NewError("segmentaion is not possible")
			}
			if _, ok := used[childVertex]; ok {
				continue
			}
			used[childVertex] = struct{}{}
			graph[childVertex] = vertex{groupNum + (graph[v].group+1)%2, graph[childVertex].neighbors}
			if err := dfs(childVertex, g); err != nil {
				return err
			}
		}
		return nil
	}
	for i := range graph {
		if _, ok := used[i]; ok {
			continue
		}
		used[i] = struct{}{}
		graph[i] = vertex{groupNum, graph[i].neighbors}
		if err := dfs(i, graph); err != nil {
			return 0, err
		}
		groupNum += 2
	}
	return groupNum, ok
}
func bruteForce(groupN int, n int, graph map[int]vertex) {
	type variant struct {
		delta   int
		vertexs []int
	}
	abs := func(a int) int {
		if a < 0 {
			a = -a
		}
		return a
	}
	compare := func(slice1 []int, slice2 []int) int {
		if len(slice1) > len(slice2) {
			return 1
		}
		if len(slice1) < len(slice2) {
			return -1
		}
		for i := 0; i < len(slice1) || i < len(slice2); i++ {
			if i >= len(slice1) {
				return -1
			}
			if i >= len(slice2) {
				return 1
			}
			if slice1[i] > slice2[i] {
				return 1
			}
			if slice1[i] < slice2[i] {
				return -1
			}
		}
		return 0
	}
	ans := variant{-1, []int{}}
	for i := 0; i < (1 << groupN); i++ {
		cur := variant{-1, []int{}}
		for vInd := range graph {
			if zero, gP2 := (i&(1<<(graph[vInd].group/2))) == 0, graph[vInd].group%2 == 0; zero && gP2 || !zero && !gP2 {
				cur.vertexs = append(cur.vertexs, vInd)
			}
		}
		cur.delta = abs(len(cur.vertexs) - (n - len(cur.vertexs)))
		sort.Ints(cur.vertexs)
		// fmt.Printf("%03b %v\n", i, cur)
		if ans.delta == -1 || ans.delta > cur.delta || (ans.delta == cur.delta && compare(ans.vertexs, cur.vertexs) > 0) {
			ans = cur
		}
	}
	for _, el := range ans.vertexs {
		fmt.Printf("%d ", el+1)
	}
	fmt.Println()
}
