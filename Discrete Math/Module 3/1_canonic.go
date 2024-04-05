package main

import "fmt"

func main() {

	numQ, numIn, startQ, QInToQ, QInToOut := scanAutomat()
	TimeToQ, newQ := numerate(numQ, numIn, startQ, QInToQ)
	printAutomat(numQ, numIn, startQ, QInToQ, QInToOut, TimeToQ, newQ)
}

func scanAutomat() (numQ, numIn, startQ int, QInToQ [][]int, QInToOut [][]string) {
	fmt.Scan(&numQ, &numIn, &startQ)
	for i := 0; i < numQ; i++ {
		QInToQ = append(QInToQ, []int{})
		for j := 0; j < numIn; j++ {
			v := 0
			fmt.Scan(&v)
			QInToQ[i] = append(QInToQ[i], v)
		}
	}
	for i := 0; i < numQ; i++ {
		QInToOut = append(QInToOut, []string{})
		for j := 0; j < numIn; j++ {
			v := ""
			fmt.Scan(&v)
			QInToOut[i] = append(QInToOut[i], v)
		}
	}
	return numQ, numIn, startQ, QInToQ, QInToOut
}

func numerate(numQ, numIn, startQ int, QInToQ [][]int) (map[int]int, map[int]int) {
	time := 0
	TimeToQ := map[int]int{}
	newQ := map[int]int{}
	used := map[int]struct{}{}
	var dfs func(int)

	dfs = func(q int) {
		used[q] = struct{}{}
		TimeToQ[time] = q
		newQ[q] = time
		time += 1
		for _, to := range QInToQ[q] {
			if _, ok := used[to]; !ok {
				dfs(to)
			}
		}
	}

	dfs(startQ)
	for i := 0; i < numQ; i++ {
		if _, ok := used[i]; !ok {
			dfs(i)
		}
	}
	return TimeToQ, newQ
}

func printAutomat(numQ, numIn, startQ int, QInToQ [][]int, QInToOut [][]string, TimeToQ map[int]int, newQ map[int]int) {
	fmt.Printf("%d\n%d\n%d\n", numQ, numIn, newQ[startQ])
	for i := 0; i < numQ; i++ {
		for _, q := range QInToQ[TimeToQ[i]] {
			fmt.Printf("%d ", newQ[q])
		}
		fmt.Println()
	}

	for i := 0; i < numQ; i++ {
		for _, q := range QInToOut[TimeToQ[i]] {
			fmt.Printf("%s ", q)
		}
		fmt.Println()
	}
}
