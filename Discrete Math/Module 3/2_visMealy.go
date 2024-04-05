package main

import "fmt"

func main() {
	printAutomat(scanAutomat())
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
func printAutomat(numQ, numIn, startQ int, QInToQ [][]int, QInToOut [][]string) {
	asciiCodeOfA := 97
	fmt.Println("digraph {")
	fmt.Printf("\trankdir = LR\n")
	for i := 0; i < numQ; i++ {
		for j := 0; j < numIn; j++ {
			fmt.Printf("\t%d -> %d [label = \"%c(%s)\"]\n", i, QInToQ[i][j], rune(asciiCodeOfA+j), QInToOut[i][j])
		}
	}
	fmt.Println("}")
}
