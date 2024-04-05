package main

import "fmt"

func main() {
	numQ, numIn, startQ, QInToQ, QInToOut, m := scanAutomat()
	for key := range generateLanguage(numQ, numIn, startQ, QInToQ, QInToOut, m) {
		fmt.Println(key)
	}

}

func scanAutomat() (numQ, numIn, startQ int, QInToQ [][]int, QInToOut [][]string, m int) {
	numIn = 2
	fmt.Scan(&numQ)
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
	fmt.Scan(&startQ, &m)
	return numQ, numIn, startQ, QInToQ, QInToOut, m
}

func generateLanguage(numQ, numIn, startQ int, QInToQ [][]int, QInToOut [][]string, maxLen int) map[string]struct{} {
	language := map[string]struct{}{}
	type state struct {
		word string
		curQ int
	}
	memo := map[state]struct{}{}
	type edge struct {
		from int
		to   int
	}
	var recoursiveGenerate func(curLen int, word string, curQ int, lambdaQ map[edge]struct{})
	recoursiveGenerate = func(curLen int, word string, curQ int, lambdaQ map[edge]struct{}) {
		memo[state{word, curQ}] = struct{}{}
		if curLen > maxLen {
			panic("curLen > maxLen")
		}
		if curLen == maxLen {
			language[word] = struct{}{}
			return
		}
		if word != "" {
			language[word] = struct{}{}
		}
		for i := range QInToQ[curQ] {

			newQ := QInToQ[curQ][i]
			newSymbol := QInToOut[curQ][i]
			newEdge := edge{curQ, newQ}
			if _, ok := lambdaQ[newEdge]; ok || newQ == curQ && newSymbol == "-" {
				continue
			}
			if _, ok := memo[state{word + newSymbol, newQ}]; ok {
				continue
			}
			if newSymbol == "-" {
				lambdaQ[newEdge] = struct{}{}
				recoursiveGenerate(curLen, word, newQ, lambdaQ)
				delete(lambdaQ, newEdge)
			} else {
				recoursiveGenerate(curLen+1, word+newSymbol, newQ, map[edge]struct{}{})
			}
		}
	}

	recoursiveGenerate(0, "", startQ, map[edge]struct{}{})
	return language
}
