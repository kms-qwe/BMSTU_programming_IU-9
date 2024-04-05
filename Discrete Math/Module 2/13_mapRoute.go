package main

import (
	"container/heap"
	"fmt"
	"math"
)

const (
	inf = math.MaxInt64
)

var (
	cost [][]int
	d    [][]int
	n    int
)

type pair struct {
	x, y int
}

type item struct {
	value    pair
	priority int
	index    int
}

type PriorityQueue []*item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	it := x.(*item)
	it.index = n
	*pq = append(*pq, it)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	old[n-1] = nil
	it.index = -1
	*pq = old[0 : n-1]
	return it
}

func dij(x, y int) {
	for i := range d {
		for j := range d[i] {
			d[i][j] = -1
		}
	}
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	it := &item{
		value:    pair{x, y},
		priority: cost[x][y],
	}
	heap.Push(&pq, it)

	for pq.Len() > 0 {
		minItem := heap.Pop(&pq).(*item)
		x, y := minItem.value.x, minItem.value.y
		d[x][y] = minItem.priority

		if x < n-1 && (d[x+1][y] == -1 || d[x+1][y] > d[x][y]+cost[x+1][y]) {
			d[x+1][y] = d[x][y] + cost[x+1][y]
			it := &item{
				value:    pair{x + 1, y},
				priority: d[x+1][y],
			}
			heap.Push(&pq, it)
		}

		if y < n-1 && (d[x][y+1] == -1 || d[x][y+1] > d[x][y]+cost[x][y+1]) {
			d[x][y+1] = d[x][y] + cost[x][y+1]
			it := &item{
				value:    pair{x, y + 1},
				priority: d[x][y+1],
			}
			heap.Push(&pq, it)
		}

		if x > 0 && (d[x-1][y] == -1 || d[x-1][y] > d[x][y]+cost[x-1][y]) {
			d[x-1][y] = d[x][y] + cost[x-1][y]
			it := &item{
				value:    pair{x - 1, y},
				priority: d[x-1][y],
			}
			heap.Push(&pq, it)
		}

		if y > 0 && (d[x][y-1] == -1 || d[x][y-1] > d[x][y]+cost[x][y-1]) {
			d[x][y-1] = d[x][y] + cost[x][y-1]
			it := &item{
				value:    pair{x, y - 1},
				priority: d[x][y-1],
			}
			heap.Push(&pq, it)
		}
	}
}

func main() {
	fmt.Scan(&n)
	cost = make([][]int, n)
	d = make([][]int, n)
	for i := range cost {
		cost[i] = make([]int, n)
		d[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fmt.Scan(&cost[i][j])
		}
	}
	dij(0, 0)
	fmt.Println(d[n-1][n-1])
}
