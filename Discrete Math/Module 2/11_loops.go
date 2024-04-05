package main

import "fmt"

type vertex struct {
	dom, sdom, ancestor, label, parent, time int
	inEdge, outEdge                          []int
	bucket                                   map[int]struct{}
}

func printG(n, start int, g map[int]*vertex) {
	fmt.Println("=== PRINT GRAPH ===")
	fmt.Println("n =", n, "| start = ", start)
	for k, v := range g {
		fmt.Println(k, v.time, v.parent, v.dom, v.inEdge, v.outEdge)
	}
	fmt.Println("=== END PRINT === ")
}
func main() {
	n, start, graph := scanGraph()
	// printG(n, start, graph)
	tin := numAndDelete(n, start, graph)
	// fmt.Println("TIN:", tin)
	// printG(n, start, graph)
	buildDominatorsTree(graph, tin)
	// printG(n, start, graph)
	// fmt.Println("COUNT LOOPS:", toCountLoops(graph))
	fmt.Println(toCountLoops(graph))
}
func scanGraph() (int, int, map[int]*vertex) {
	g := make(map[int]*vertex)
	n := 0
	fmt.Scan(&n)

	prevInd := -1
	prevCommand := ""
	start := -1

	for i := 0; i < n; i++ {
		ind := 0
		command := ""
		operand := 0
		fmt.Scan(&ind, &command)
		if _, ok := g[ind]; !ok {
			g[ind] = &vertex{
				dom:      -1,
				sdom:     ind,
				label:    ind,
				ancestor: -1,
				parent:   -1,
				time:     -1,
				inEdge:   []int{},
				outEdge:  []int{},
				bucket:   map[int]struct{}{},
			}
		}
		if command == "BRANCH" || command == "JUMP" {
			fmt.Scan(&operand)
			if _, ok := g[operand]; !ok {
				g[operand] = &vertex{
					dom:      -1,
					sdom:     operand,
					label:    operand,
					ancestor: -1,
					parent:   -1,
					time:     -1,
					inEdge:   []int{},
					outEdge:  []int{},
					bucket:   map[int]struct{}{},
				}
			}
			g[ind].outEdge = append(g[ind].outEdge, operand)
			g[operand].inEdge = append(g[operand].inEdge, ind)
		}
		if prevInd == -1 {
			start = ind
		}
		if prevInd != -1 && prevCommand != "JUMP" {
			g[ind].inEdge = append(g[ind].inEdge, prevInd)
			g[prevInd].outEdge = append(g[prevInd].outEdge, ind)
		}
		prevInd = ind
		prevCommand = command
	}
	return n, start, g
}

func numAndDelete(n, s int, g map[int]*vertex) []int {
	tin := []int{}
	used := map[int]struct{}{}
	time := 0
	var dfs func(int, int)
	dfs = func(parent, v int) {
		used[v] = struct{}{}
		tin = append(tin, v)
		g[v].parent = parent
		g[v].time = time
		time += 1
		for _, to := range g[v].outEdge {
			if _, ok := used[to]; !ok {
				dfs(v, to)
			}
		}
	}
	dfs(-1, s)
	for v := range g {
		if _, ok := used[v]; !ok {
			delete(g, v)
			continue
		}
		newInEdge := []int{}
		for _, w := range g[v].inEdge {
			if _, ok := used[w]; ok {
				newInEdge = append(newInEdge, w)
			}
		}
		g[v].inEdge = newInEdge
		newOutEdge := []int{}
		for _, w := range g[v].outEdge {
			if _, ok := used[w]; ok {
				newOutEdge = append(newOutEdge, w)
			}
		}
		g[v].outEdge = newOutEdge
	}
	return tin
}

func findMin(v int, g map[int]*vertex) int {
	searchAndCut(v, g)
	return g[v].label
}
func searchAndCut(v int, g map[int]*vertex) int {
	if g[v].ancestor == -1 {
		return v
	}
	root := searchAndCut(g[v].ancestor, g)
	if g[g[g[g[v].ancestor].label].sdom].time < g[g[g[v].label].sdom].time {
		g[v].label = g[g[v].ancestor].label
	}
	g[v].ancestor = root
	return root
}
func buildDominatorsTree(g map[int]*vertex, tin []int) {
	for i := 0; i < len(g); i++ {
		w := tin[len(g)-1-i]
		for _, v := range g[w].inEdge {
			u := findMin(v, g)
			if g[g[u].sdom].time < g[g[w].sdom].time {
				g[w].sdom = g[u].sdom
			}
		}
		g[w].ancestor = g[w].parent
		g[g[w].sdom].bucket[w] = struct{}{}

		if w == tin[0] {
			continue
		}
		for v := range g[g[w].parent].bucket {
			u := findMin(v, g)
			if g[g[u].sdom].time == g[g[v].sdom].time {
				g[v].dom = g[v].sdom
			} else {
				g[v].dom = u
			}
		}
		g[g[w].parent].bucket = map[int]struct{}{}

	}

	for i := 1; i < len(g); i++ {
		w := tin[i]
		if g[g[w].dom].time != g[g[w].sdom].time {
			g[w].dom = g[g[w].dom].dom
		}
	}
}
func toCountLoops(g map[int]*vertex) int {
	countLoops := 0
	for v := range g {
		for _, u := range g[v].inEdge {
			for u != -1 && v != u {
				u = g[u].dom
			}
			if v == u {
				countLoops += 1
				break
			}
		}
	}
	return countLoops
}
