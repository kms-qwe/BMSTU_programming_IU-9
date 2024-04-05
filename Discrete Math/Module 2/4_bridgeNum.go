package main

import "fmt"

func main() {
	graph := make(map[int][]int)
	scanGraph(graph)
	fmt.Println(findBridges(graph))
}
func scanGraph(graph map[int][]int) {
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
}
func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
func findBridges(graph map[int][]int) int {
	count := 0
	timer := 0
	used := make(map[int]struct{})
	tin := make(map[int]int)
	fup := make(map[int]int)
	var dfs func(int, int)
	dfs = func(v int, p int) {
		used[v] = struct{}{}
		tin[v], fup[v] = timer, timer
		timer += 1
		for _, to := range graph[v] {
			if to == p {
				continue
			}
			if _, ok := used[to]; ok {
				fup[v] = min(fup[v], tin[to])
			} else {
				dfs(to, v)
				fup[v] = min(fup[v], fup[to])
				if fup[to] > tin[v] {
					count += 1
				}
			}
		}

	}

	for v := range graph {
		if _, ok := used[v]; !ok {
			dfs(v, -1)
		}
	}
	return count
}
