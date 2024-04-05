package main

import (
	"fmt"
	"math"
	"sort"
)

type vertex struct {
	ind int
	x   int
	y   int
}

func (v vertex) getDist(to vertex) float64 {
	return math.Sqrt(math.Pow(float64(v.x-to.x), 2) + math.Pow(float64(v.y-to.y), 2))
}

type edge struct {
	v1, v2 int
	dist   float64
}
type edgeSlice []edge

func (e edgeSlice) Len() int {
	return len(e)
}
func (e edgeSlice) Less(i, j int) bool {
	if e[i].dist < e[j].dist {
		return true
	}
	return false
}
func (e edgeSlice) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
func main() {
	verts := scanVertex()
	edges := getEdge(verts)
	sort.Sort(edges)
	fmt.Printf("%0.2f\n", kruskal(edges, len(verts)))

}

func scanVertex() []vertex {
	n := 0
	fmt.Scan(&n)
	verts := make([]vertex, n)
	for i := 0; i < n; i++ {
		x, y := 0, 0
		fmt.Scan(&x, &y)
		verts[i] = vertex{i, x, y}
	}
	return verts
}
func getEdge(verts []vertex) edgeSlice {
	edges := make([]edge, 0)
	for i := 0; i < len(verts); i++ {
		for j := 0; j < i; j++ {
			edges = append(edges, edge{j, i, verts[i].getDist(verts[j])})
		}
	}
	return edges
}

type subset struct {
	parent int
	rank   int
}

func find(subsets []subset, i int) int {
	if subsets[i].parent != i {
		subsets[i].parent = find(subsets, subsets[i].parent)
	}
	return subsets[i].parent
}

func union(subsets []subset, x, y int) {
	rootX := find(subsets, x)
	rootY := find(subsets, y)

	if subsets[rootX].rank < subsets[rootY].rank {
		subsets[rootX].parent = rootY
	} else if subsets[rootX].rank > subsets[rootY].rank {
		subsets[rootY].parent = rootX
	} else {
		subsets[rootY].parent = rootX
		subsets[rootX].rank++
	}
}

func kruskal(edges []edge, v int) float64 {
	result := 0.0
	subsets := make([]subset, v)

	for i := 0; i < v; i++ {
		subsets[i].parent = i
		subsets[i].rank = 0
	}
	e := 0
	i := 0
	for e < v-1 && i < len(edges) {
		nextEdge := edges[i]
		i++

		x := find(subsets, nextEdge.v1)
		y := find(subsets, nextEdge.v2)

		if x != y {
			result += nextEdge.dist
			union(subsets, x, y)
			e++
		}
	}

	return result
}
