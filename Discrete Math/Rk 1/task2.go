package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	n, _ := strconv.Atoi(strings.TrimSpace(line))

	graph := make([][]int, n)
	inDegree := make([]int, n)

	for i := 0; i < n; i++ {
		line, _ = reader.ReadString('\n')
		parts := strings.Fields(line)
		k, _ := strconv.Atoi(parts[0])
		if k > 0 {
			for _, part := range parts[1:] {
				pre, _ := strconv.Atoi(part)
				graph[pre-1] = append(graph[pre-1], i)
				inDegree[i]++
			}
		}
	}

	queue := []int{}
	semesters := make([]int, n)
	for i := 0; i < n; i++ {
		if inDegree[i] == 0 {
			queue = append(queue, i)
			semesters[i] = 1
		}
	}

	index := 0
	for index < len(queue) {
		course := queue[index]
		index++
		for _, nextCourse := range graph[course] {
			inDegree[nextCourse]--
			if inDegree[nextCourse] == 0 {
				queue = append(queue, nextCourse)
				semesters[nextCourse] = semesters[course] + 1
			}
		}
	}

	if len(queue) != n {
		fmt.Println("cycle")
	} else {
		maxSemesters := 0
		for _, s := range semesters {
			if s > maxSemesters {
				maxSemesters = s
			}
		}
		fmt.Println(maxSemesters)
	}
}
