package main

import (
	"fmt"
)

type comp struct {
	startVertex int
	countVertex int

	countEdge int
}

func main() {
	graph := make(map[int][]int)
	countVertexInComponent := make(map[int]comp)
	scanGraph(graph)
	// printDOTWithColor(graph, 0)
	// fmt.Println(graph)
	maxComp(graph, countVertexInComponent)
	maxNumberVertex := maxComponent(countVertexInComponent)
	// fmt.Println(countVertexInComponent)
	// fmt.Printf("%#v\n", maxNumberVertex)
	printDOTWithColor(graph, maxNumberVertex.startVertex)

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
func printDOTWithColor(graph map[int][]int, maxNumberVertex int) {
	fmt.Println("graph {")

	type edge struct {
		v1, v2 int
	}
	vertexes := make(map[int]struct{})
	vertexes[maxNumberVertex] = struct{}{}
	alreadyWorked := make(map[int]struct{})
	queueVertex := queue{}
	alreadyWorked[maxNumberVertex] = struct{}{}
	queueVertex.enqueue(maxNumberVertex)
	for !queueVertex.empty() {
		// fmt.Println(queueVertex.items, alreadyWorked)
		vertex := (queueVertex.dequeue()).(int)
		for _, key2 := range graph[vertex] {
			if _, ok := alreadyWorked[key2]; ok {
				continue
			}
			queueVertex.enqueue(key2)
			vertexes[key2] = struct{}{}
			alreadyWorked[key2] = struct{}{}
		}
	}
	// fmt.Println("Hl")
	alreadyPrint := make(map[edge]struct{})
	for key := range graph {
		if _, ok1 := vertexes[key]; ok1 {
			fmt.Printf("\t %d [color = red]\n", key)
			continue
		}
		fmt.Println("\t", key)
	}
	for key, val := range graph {
		for _, vertex := range val {
			if _, ok := alreadyPrint[edge{key, vertex}]; !ok {
				if _, ok := alreadyPrint[edge{vertex, key}]; !ok {
					_, ok1 := vertexes[key]
					_, ok2 := vertexes[vertex]
					if ok1 || ok2 {
						fmt.Printf("\t %d--%d [color = red]\n", key, vertex)
					} else {
						fmt.Printf("\t %d--%d\n", key, vertex)
					}

				}

			}
			alreadyPrint[edge{key, vertex}] = struct{}{}
			alreadyPrint[edge{vertex, key}] = struct{}{}
		}
	}
	fmt.Println("}")
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
func maxComp(graph map[int][]int, countVertexInComponent map[int]comp) {
	alreadyWorked := make(map[int]struct{})
	queueVertex := queue{}
	numberOfComponent := 0
	for key := range graph {
		if _, ok := alreadyWorked[key]; ok {
			continue
		}
		alreadyWorked[key] = struct{}{}
		queueVertex.enqueue(key)
		countVertexInComponent[numberOfComponent] = comp{key, 1, 0}
		for !queueVertex.empty() {
			// fmt.Println(key, queueVertex.items, alreadyWorked, countVertexInComponent)
			vertex := (queueVertex.dequeue()).(int)
			for _, key2 := range graph[vertex] {
				countVertexInComponent[numberOfComponent] = comp{countVertexInComponent[numberOfComponent].startVertex, countVertexInComponent[numberOfComponent].countVertex, countVertexInComponent[numberOfComponent].countEdge + 1}
				if countVertexInComponent[numberOfComponent].startVertex > key2 {
					countVertexInComponent[numberOfComponent] = comp{key2, countVertexInComponent[numberOfComponent].countVertex, countVertexInComponent[numberOfComponent].countEdge}
				}
				if _, ok := alreadyWorked[key2]; ok {
					continue
				}
				countVertexInComponent[numberOfComponent] = comp{countVertexInComponent[numberOfComponent].startVertex, countVertexInComponent[numberOfComponent].countVertex + 1, 0}
				alreadyWorked[key2] = struct{}{}
				queueVertex.enqueue(key2)
				// fmt.Println("VERTEX KEY2", vertex, key2)
			}
		}
		numberOfComponent += 1
	}
}
func maxComponent(countVertexInComponent map[int]comp) comp {
	maxComp := comp{-1, -1, -1}
	for _, val := range countVertexInComponent {
		if val.countVertex < maxComp.countVertex {
			continue
		}
		if val.countVertex > maxComp.countVertex {
			maxComp = val
			continue
		}
		if val.countEdge < maxComp.countEdge {
			continue
		}
		if val.countEdge > maxComp.countEdge {
			maxComp = val
			continue
		}
		if val.startVertex > maxComp.startVertex {
			continue
		}
		maxComp = val
	}
	return maxComp
}
