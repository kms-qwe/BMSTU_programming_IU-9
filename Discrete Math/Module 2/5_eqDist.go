package main

import "fmt"

func main() {
	n, m, graph := scanGraph()
	vSupport := scanSlice()
	printIntSlice(eqDist(n, m, graph, vSupport))
}
func printIntSlice(slice []int) {
	if len(slice) == 0 {
		fmt.Printf("-")
	}
	for i := 0; i < len(slice); i++ {
		fmt.Printf("%d ", slice[i])
	}
	fmt.Println()
}
func eqDist(n, m int, graph map[int][]int, vSup map[int]int) []int {
	distMatrix := make([][]int, n)
	for i := 0; i < n; i++ {
		distMatrix[i] = make([]int, len(vSup))
	}
	for v := range vSup {
		// fmt.Println("NEW", v)
		bfs(graph, v, distMatrix, vSup)
	}
	// fmt.Println(distMatrix)
	return checkMatrix(distMatrix)
}
func checkMatrix(distMatrix [][]int) []int {
	res := []int{}
	for i := 0; i < len(distMatrix); i++ {
		dist := distMatrix[i][0]
		for j := 1; j < len(distMatrix[i]); j++ {
			if distMatrix[i][j] != dist {
				dist = -1
				break
			}
		}
		if dist > 0 {
			res = append(res, i)
		}
	}
	return res
}
func bfs(graph map[int][]int, v int, distMatrix [][]int, vSup map[int]int) {
	used := make(map[int]struct{})
	vQueue := queue{}
	used[v] = struct{}{}
	vQueue.enqueue(elem{v, 0})
	for !vQueue.empty() {
		to := (vQueue.dequeue()).(elem)
		distMatrix[to.vertex][vSup[v]] = to.dist
		if i, ok := vSup[to.vertex]; ok {
			if v == to.vertex {
				distMatrix[v][i] = 0
			} else {
				distMatrix[v][i] = to.dist
			}
		}
		// fmt.Println(v, to, to.dist, vQueue.items, distMatrix, used, graph)
		for _, toto := range graph[to.vertex] {
			if _, ok := used[toto]; ok {
				continue
			}
			// fmt.Println(to, toto, graph[to.vertex])
			vQueue.enqueue(elem{toto, to.dist + 1})
			used[toto] = struct{}{}
		}
	}
}
func scanSlice() map[int]int {
	k := 0
	fmt.Scan(&k)
	vSup := make(map[int]int)
	for i := 0; i < k; i++ {
		v := 0
		fmt.Scan(&v)
		vSup[v] = i
	}
	return vSup
}
func scanGraph() (int, int, map[int][]int) {
	graph := make(map[int][]int)
	n, m := 0, 0
	fmt.Scan(&n, &m)
	for i := 0; i < n; i++ {
		graph[i] = make([]int, 0)
	}
	for i := 0; i < m; i++ {
		vertex1, vertex2 := 0, 0
		fmt.Scan(&vertex1, &vertex2)
		if _, ok := graph[vertex1]; !ok {
			graph[vertex1] = make([]int, 0)
		}
		if _, ok := graph[vertex2]; !ok {
			graph[vertex2] = make([]int, 0)
		}
		graph[vertex1] = append(graph[vertex1], vertex2)
		if vertex1 != vertex2 {
			graph[vertex2] = append(graph[vertex2], vertex1)
		}
	}
	return n, m, graph
}

type elem struct {
	vertex int
	dist   int
}
type queue struct {
	items []interface{}
}

func (q *queue) enqueue(item interface{}) {
	q.items = append(q.items, item)
}
func (q *queue) dequeue() interface{} {
	if len(q.items) == 0 {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}
func (q *queue) empty() bool {
	return len(q.items) == 0
}
