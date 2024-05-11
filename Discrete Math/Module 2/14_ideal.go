package main

import (
	"container/heap"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	n, m, graph := scanGraph()
	it := bfs(n, m, graph)
	if it == nil {
		fmt.Println("ERROR NIL POINTER")
	}
	fmt.Printf("%d\n%s\n", it.depth, it.path)
}

func scanGraph() (n int, m int, graph []map[int]int) {
	fmt.Scan(&n, &m)
	graph = make([]map[int]int, n+1)
	for i := 0; i < n+1; i++ {
		graph[i] = map[int]int{}
	}
	for i := 0; i < m; i++ {
		vFrom, vTo, colorNum := 0, 0, 0
		fmt.Scan(&vFrom, &vTo, &colorNum)
		if vFrom == vTo {
			continue
		}
		if val, ok := graph[vFrom][vTo]; !ok || ok && val > colorNum {
			graph[vFrom][vTo] = colorNum
		}
		if val, ok := graph[vTo][vFrom]; !ok || ok && val > colorNum {
			graph[vTo][vFrom] = colorNum
		}
	}
	return n, m, graph
}

type item struct {
	color  int
	depth  int
	index  int
	vertex int
	path   string
}
type priorityQueue []*item

func (pq priorityQueue) Len() int {
	return len(pq)
}
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].depth < pq[j].depth || pq[i].depth == pq[j].depth && pq[i].color < pq[j].color
}
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}
func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	it := x.(*item)
	it.index = n
	*pq = append(*pq, it)
}
func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	it := old[n-1]
	old[n-1] = nil
	it.index = -1
	*pq = old[0 : n-1]
	return it
}

func bfs(n, m int, graph []map[int]int) *item {
	used := map[int]item{}
	pq := make(priorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &item{color: 0, depth: 0, index: -1, vertex: 1, path: ""})
	used[1] = item{color: 0, depth: 0, index: -1, vertex: 1, path: ""}

	minV := &item{}
	for pq.Len() > 0 {

		v := heap.Pop(&pq).(*item)
		if v.vertex == n {

			if minV.depth == 0 || minV.depth > v.depth || minV.depth == v.depth && compInt(minV.path, v.path) > 0 {
				minV = v

			}
		}
		for to, col := range graph[v.vertex] {
			if it, ok := used[to]; !ok || it.depth > v.depth+1 || it.depth == v.depth+1 && compInt(it.path, v.path+" "+strconv.Itoa(col)) > 0 {

				newIt := item{
					color:  col,
					depth:  v.depth + 1,
					index:  -1,
					vertex: to,
					path:   v.path + " " + strconv.Itoa(col),
				}
				used[to] = newIt
				if ok {
					for i := range pq {
						if pq[i].vertex == to {
							heap.Remove(&pq, i)
							break
						}
					}
				}
				heap.Push(&pq, &newIt)

			}
		}
	}
	return minV
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func compInt(str1, str2 string) (res int) {
	nums1, nums2 := toInt(str1), toInt(str2)
	for i := 0; i < min(len(nums1), len(nums2)); i++ {
		if nums1[i] < nums2[i] {
			return -1
		}
		if nums1[i] > nums2[i] {
			return 1
		}
	}
	if len(nums1) < len(nums2) {
		return -1
	}
	if len(nums1) > len(nums2) {
		return 1
	}
	return 0
}

func toInt(input string) []int {
	numStr := strings.Fields(input)
	numbers := make([]int, len(numStr))

	for i, str := range numStr {
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Ошибка при преобразовании числа: %v\n", err)
			return nil
		}
		numbers[i] = num
	}
	return numbers
}
