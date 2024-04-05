package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type vertex struct {
	duration int
	color    string
	to       []string
}
type edge struct {
	from, to string
	color    string
}

func main() {
	graph, err := scanGraph()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	paintCycle(graph)
	red := criticalPath(graph)
	printGraph(graph, red)
}
func scanGraph() (map[string]vertex, error) {
	g := make(map[string]vertex)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}
		if line[len(line)-1] != ';' {
			for scanner.Scan() {
				add := strings.TrimSpace(scanner.Text())
				line += add
				if add == "" || add[len(add)-1] == ';' {
					break
				}
			}
		}

		line = strings.ReplaceAll(line, "<", "")
		line = strings.ReplaceAll(line, ";", "")
		line = strings.ReplaceAll(line, "(", " ")
		line = strings.ReplaceAll(line, ")", "")
		dependencies := strings.Fields(line)
		i := 1
		incDif := false
		if _, ok := g[dependencies[0]]; !ok {
			num, err := strconv.Atoi(dependencies[1])
			if err != nil {
				return nil, err
			}
			g[dependencies[0]] = vertex{num, "black", []string{}}
			i += 1
			incDif = true
		}
		dif := 1
		if incDif {
			dif = 2
		} else {
			dif = 1
		}
		for ; i < len(dependencies); i++ {
			incDif = false
			incI := false
			if _, ok := g[dependencies[i]]; !ok {
				num, err := strconv.Atoi(dependencies[i+1])
				if err != nil {
					return nil, err
				}
				g[dependencies[i]] = vertex{num, "black", []string{}}
				incI = true
				incDif = true
			}

			v := dependencies[i-dif]
			color := "black"
			if v == dependencies[i] || g[v].color == "blue" {
				color = "blue"
			}
			g[v] = vertex{g[v].duration, color, append(g[v].to, dependencies[i])}
			if incDif {
				dif = 2
			} else {
				dif = 1
			}
			if incI {
				i += 1
			}
		}

	}
	return g, nil
}
func printGraph(g map[string]vertex, red map[edge]struct{}) {
	fmt.Printf("digraph{\n")
	for name, v := range g {
		if v.color == "black" {
			fmt.Printf("\t%s [label = \"%s\"]\n", name, name+"("+strconv.Itoa(v.duration)+")")
		} else {
			fmt.Printf("\t%s [label = \"%s\", color = %s]\n", name, name+"("+strconv.Itoa(v.duration)+")", v.color)
		}

	}
	for from, v := range g {
		for _, name := range v.to {
			color := v.color
			if _, ok := red[edge{from, name, "red"}]; ok {
				color = "red"
			} else if v.color == "red" {
				color = "black"
			}
			if color == "black" {
				fmt.Printf("\t%s -> %s\n", from, name)
			} else {
				fmt.Printf("\t%s -> %s [color = %s]\n", from, name, color)
			}

		}
	}
	fmt.Printf("}\n")
}
func reverseGraph(g map[string]vertex) map[string]vertex {
	gR := make(map[string]vertex)
	for name := range g {
		gR[name] = vertex{g[name].duration, g[name].color, []string{}}
	}
	for name, v := range g {
		for _, nameTo := range v.to {
			gR[nameTo] = vertex{gR[nameTo].duration, gR[nameTo].color, append(gR[nameTo].to, name)}
		}
	}
	return gR
}
func buildCondensation(n int, graph, graphT map[string]vertex) []string {
	used := make(map[string]struct{})
	usedT := make(map[string]struct{})
	order := []string{}
	component := []string{}
	vInCycle := []string{}

	var dfs1 func(string)
	dfs1 = func(v string) {
		used[v] = struct{}{}
		for _, to := range graph[v].to {
			if _, ok := used[to]; !ok {
				dfs1(to)
			}
		}
		order = append(order, v)
	}

	var dfs2 func(string)
	dfs2 = func(v string) {
		usedT[v] = struct{}{}
		component = append(component, v)
		for _, to := range graphT[v].to {
			if _, ok := usedT[to]; !ok {
				dfs2(to)
			}
		}
	}

	for name := range graph {
		if _, ok := used[name]; !ok {
			dfs1(name)
		}
	}
	for i := 0; i < n; i++ {
		v := order[n-1-i]
		if _, ok := usedT[v]; !ok {
			dfs2(v)
			if len(component) > 1 {
				vInCycle = append(vInCycle, component...)
			}
			component = []string{}
		}
	}
	return vInCycle
}
func paintCycle(g map[string]vertex) {
	vToPaint := buildCondensation(len(g), g, reverseGraph(g))
	used := make(map[string]struct{})

	var dfs func(string)
	dfs = func(name string) {
		used[name] = struct{}{}
		g[name] = vertex{g[name].duration, "blue", g[name].to}
		for _, to := range g[name].to {
			if _, ok := used[to]; !ok {
				dfs(to)
			}
		}
	}

	for _, name := range vToPaint {
		if _, ok := used[name]; !ok {
			dfs(name)
		}
	}
}
func criticalPath(g map[string]vertex) map[edge]struct{} {
	getRoot := func(g map[string]vertex) []string {
		roots := []string{}
		for nameV, v := range reverseGraph(g) {
			if len(v.to) == 0 {
				roots = append(roots, nameV)
			}
		}
		return roots
	}
	roots := getRoot(g)
	if roots == nil {
		return nil
	}
	ans := [][]string{}
	bestPath := []string{}
	bestScore := -1
	var dfs func(string, int, []string)
	dfs = func(name string, score int, path []string) {
		path = append(path, name)
		score += g[name].duration
		anyChild := false
		for _, to := range g[name].to {
			if g[to].color != "blue" {
				anyChild = true
				dfs(to, score, path)
			}
		}
		if !anyChild {
			if bestScore < score {
				bestScore = score
				bestPath = path
				ans = [][]string{}
				cp := make([]string, len(bestPath))
				copy(cp, bestPath)
				ans = append(ans, cp)
			} else if bestScore == score {
				cp := make([]string, len(path))
				copy(cp, path)
				ans = append(ans, cp)
			}
		}
	}
	for _, root := range roots {
		dfs(root, 0, []string{})
	}
	for _, best := range ans {
		for _, name := range best {
			g[name] = vertex{g[name].duration, "red", g[name].to}
		}
	}
	res := make(map[edge]struct{})
	for _, best := range ans {
		for i := 0; i < len(best)-1; i++ {
			res[edge{best[i], best[i+1], "red"}] = struct{}{}
		}
	}
	return res
}
