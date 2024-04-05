package main

import "fmt"

func isSmple(x int) bool {
	for i := 2; i*i <= x; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}
func simpleDiv(x int) []int {
	divisors := []int{}
	if isSmple(x) {
		divisors = append(divisors, x)
		return divisors
	}
	for i := 2; x > 1; i++ {
		if x%i == 0 {
			divisors = append(divisors, i)
		}
		for x%i == 0 {
			x /= i
		}
	}
	return divisors
}
func buildGraphDFS(graph map[int][]int, vertex int, divs []int, done map[int]bool) {
	done[vertex] = true
	if vertex == 1 {
		return
	}
	if len(divs) == 0 {
		graph[vertex] = append(graph[vertex], 1)
		return
	}
	for _, div := range divs {
		graph[vertex] = append(graph[vertex], vertex/div)
		if ok := done[vertex/div]; !ok {

			buildGraphDFS(graph, vertex/div, simpleDiv(vertex/div), done)
		}
	}
}
func PrintDOT(graph map[int][]int) {
	fmt.Println("graph {")
	for key := range graph {
		fmt.Println("\t", key)
	}
	for key, val := range graph {
		for _, vertex := range val {
			fmt.Println("\t", key, "--", vertex)
		}
	}
	fmt.Println("}")
}
func main() {
	var x int
	fmt.Scan(&x)
	graph := map[int][]int{
		x: {},
	}
	done := map[int]bool{}
	buildGraphDFS(graph, x, simpleDiv(x), done)
	// fmt.Println(graph)
	PrintDOT(graph)

}
