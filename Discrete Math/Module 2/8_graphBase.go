package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	n, m, graph, graphT := scanGraph()
	countScc, scc, vComp, order := buildCondensation(n, m, graph, graphT)
	for _, el := range baseNum(n, graph, countScc, scc, vComp, order) {
		fmt.Printf("%d ", el)
	}
	fmt.Println()

}
func scanGraph() (int, int, map[int][]int, map[int][]int) {
	n, m := 0, 0
	fmt.Scan(&n)
	fmt.Scan(&m)
	graph := make(map[int][]int)
	graphT := make(map[int][]int)
	for i := 0; i < n; i++ {
		graph[i] = []int{}
		graphT[i] = []int{}
	}
	for i := 0; i < m; i++ {
		from, to := 0, 0
		fmt.Scan(&from, &to)
		graph[from] = append(graph[from], to)
		graphT[to] = append(graphT[to], from)
	}
	return n, m, graph, graphT
}
func buildCondensation(n, m int, graph, graphT map[int][]int) (int, map[int]map[int]struct{}, map[int]int, []int) {
	scc := make(map[int]map[int]struct{})
	countScc := 0
	used := make(map[int]struct{})
	usedT := make(map[int]struct{})
	order := []int{}
	component := []int{}
	vComp := make(map[int]int)

	var dfs1 func(int)
	dfs1 = func(v int) {
		used[v] = struct{}{}
		for _, to := range graph[v] {
			if _, ok := used[to]; !ok {
				dfs1(to)
			}
		}
		order = append(order, v)
	}

	var dfs2 func(int)
	dfs2 = func(v int) {
		usedT[v] = struct{}{}
		component = append(component, v)
		for _, to := range graphT[v] {
			if _, ok := usedT[to]; !ok {
				dfs2(to)
			}
		}
	}

	for i := 0; i < n; i++ {
		if _, ok := used[i]; !ok {
			dfs1(i)
		}
	}
	for i := 0; i < n; i++ {
		v := order[n-1-i]
		if _, ok := usedT[v]; !ok {
			dfs2(v)

			countScc += 1
			scc[countScc] = make(map[int]struct{})
			for _, w := range component {
				scc[countScc][w] = struct{}{}
				vComp[w] = countScc
			}
			component = []int{}
		}
	}
	return countScc, scc, vComp, order
}
func baseNum(n int, graph map[int][]int, countScc int, scc map[int]map[int]struct{}, vComp map[int]int, order []int) []int {
	ans := []int{}
	withParents := make(map[int]struct{})
	used := make(map[int]struct{})

	var dfs func(parentV int, v int)
	dfs = func(parentV int, v int) {
		used[v] = struct{}{}
		// if vComp[parentV] == vComp[v] {
		// 	delete(used, v)
		// }
		if vComp[parentV] != vComp[v] {
			withParents[vComp[v]] = struct{}{}
		}
		for _, to := range graph[v] {
			if _, ok := used[to]; !ok {
				dfs(parentV, to)
			}
		}
	}
	for i := 0; i < n; i++ {
		v := order[n-1-i]
		if _, ok := used[v]; !ok {
			dfs(v, v)
		}
	}
	setGetMin := func(set map[int]struct{}) int {
		m := math.MaxInt64
		for k := range set {
			if m > k {
				m = k
			}
		}
		return m
	}

	for sccNum, comp := range scc {
		if _, ok := withParents[sccNum]; !ok {
			ans = append(ans, setGetMin(comp))
		}
	}
	sort.Sort(sort.IntSlice(ans))
	return ans

}
