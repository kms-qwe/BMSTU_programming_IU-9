package main

import (
	"container/heap"
	"fmt"
)

type edge struct {
	v1, v2 int
	cost   int
}

type edges []edge

func (e edges) Len() int {
	return len(e)
}
func (e edges) Less(i, j int) bool {
	return e[i].cost < e[j].cost
}
func (e edges) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e *edges) Push(x interface{}) {
	*e = append(*e, x.(edge))
}

func (e *edges) Pop() interface{} {
	old := *e
	n := len(old)
	x := old[n-1]
	*e = old[0 : n-1]
	return x
}

func (e *edges) Remove(index int) {
	if index >= 0 && index < len(*e) {
		(*e)[index] = (*e)[len(*e)-1]
		*e = (*e)[:len(*e)-1]
	}
}

func scanGraph() (int, int, map[int][]edge) {
	graph := make(map[int][]edge)
	n, m := 0, 0
	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		graph[i] = make([]edge, 0)
	}
	fmt.Scan(&m)
	for i := 0; i < m; i++ {
		v1, v2, cost := 0, 0, 0
		fmt.Scan(&v1, &v2, &cost)
		graph[v1] = append(graph[v1], edge{v1, v2, cost})
		graph[v2] = append(graph[v2], edge{v2, v1, cost})
	}
	return n, m, graph
}

func prim(n int, graph map[int][]edge) int {
	resultCost := 0
	alreadyWorked := make(map[int]struct{})
	edgesHeap := &edges{}

	for _, e := range graph[0] {
		heap.Push(edgesHeap, e)
	}
	alreadyWorked[0] = struct{}{}

	for edgesHeap.Len() > 0 {
		minEdge := heap.Pop(edgesHeap).(edge)
		if _, ok := alreadyWorked[minEdge.v2]; ok {
			continue
		}
		resultCost += minEdge.cost
		alreadyWorked[minEdge.v2] = struct{}{}
		for _, w := range graph[minEdge.v2] {
			if _, ok := alreadyWorked[w.v2]; !ok {
				heap.Push(edgesHeap, w)
			}
		}
	}

	return resultCost
}

func main() {
	n, _, graph := scanGraph()
	fmt.Println(prim(n, graph))
}
