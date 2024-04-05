package main

import (
	"fmt"
	"time"
)

var start time.Time

func init() {
	start = time.Now()
}

func main() {
	numQ1, numIn1, startQ1, QInToQ1, QInToOut1 := scanAutomat()
	newNumQ1, newStartQ1, newQInToQ1, newQInToOut1 := aufenkampHohn(numQ1, numIn1, startQ1, QInToQ1, QInToOut1)
	TimeToQ1, newQ1 := numerate(newNumQ1, numIn1, newStartQ1, newQInToQ1)

	numQ2, numIn2, startQ2, QInToQ2, QInToOut2 := scanAutomat()
	newNumQ2, newStartQ2, newQInToQ2, newQInToOut2 := aufenkampHohn(numQ2, numIn2, startQ2, QInToQ2, QInToOut2)
	TimeToQ2, newQ2 := numerate(newNumQ2, numIn2, newStartQ2, newQInToQ2)

	if compareMealy(newNumQ1, numIn1, newQInToQ1, newQInToOut1, TimeToQ1, newQ1, newNumQ2, numIn2, newQInToQ2, newQInToOut2, TimeToQ2, newQ2) {
		fmt.Println("EQUAL")
	} else {
		fmt.Println("NOT EQUAL")
	}
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

type dsu struct {
	parent []int
	rank   []int
}

func (d *dsu) initDsu(numQ int) {
	d.parent, d.rank = make([]int, numQ), make([]int, numQ)
	for i := range d.parent {
		d.parent[i] = i
	}
}
func (d *dsu) find(q int) int {
	if d.parent[q] != q {
		d.parent[q] = d.find(d.parent[q])
	}
	return d.parent[q]
}
func (d *dsu) union(q1, q2 int) {
	p1, p2 := d.find(q1), d.find(q2)
	if d.rank[p1] > d.rank[p2] {
		p1, p2 = p2, p1
	}
	d.parent[p1] = p2
	if d.rank[p1] == d.rank[p2] {
		d.rank[p2] += 1
	}
}

func split1(numQ, numIn int, QInToQ [][]int, QInToOut [][]string) (m int, pi []int) {
	m = numQ
	dsuQ := dsu{[]int{}, []int{}}
	dsuQ.initDsu(numQ)
	// fmt.Println(dsuQ)
	for q1 := 0; q1 < numQ; q1++ {
		for q2 := q1 + 1; q2 < numQ; q2++ {

			if dsuQ.find(q1) == dsuQ.find(q2) {
				continue
			}

			eq := true
			for x := 0; x < numIn; x++ {
				if QInToOut[q1][x] != QInToOut[q2][x] {
					eq = false
					break
				}
			}

			if eq {
				dsuQ.union(q1, q2)
				m -= 1
			}
		}
	}
	for q := 0; q < numQ; q++ {
		pi = append(pi, dsuQ.find(q))
	}
	// fmt.Println("END SPLIT1")
	return m, pi
}

func split(numQ, numIn int, QInToQ [][]int, QInToOut [][]string, pi []int) (int, []int) {
	m := numQ
	dsuQ := dsu{[]int{}, []int{}}
	dsuQ.initDsu(numQ)
	for q1 := 0; q1 < numQ; q1++ {
		for q2 := q1 + 1; q2 < numQ; q2++ {

			if pi[q1] != pi[q2] || dsuQ.find(q1) == dsuQ.find(q2) {
				continue
			}

			eq := true
			for x := 0; x < numIn; x++ {
				if pi[QInToQ[q1][x]] != pi[QInToQ[q2][x]] {
					eq = false
					break
				}
			}

			if eq {
				dsuQ.union(q1, q2)
				m -= 1
			}
		}
	}
	for q := 0; q < numQ; q++ {
		pi[q] = dsuQ.find(q)
	}
	return m, pi
}

func aufenkampHohn(numQ, numIn, startQ int, QInToQ [][]int, QInToOut [][]string) (newNumQ, newStartQ int, newQInToQ [][]int, newQInToOut [][]string) {
	m, pi := split1(numQ, numIn, QInToQ, QInToOut)
	for {
		newM, newPi := split(numQ, numIn, QInToQ, QInToOut, pi)
		pi = newPi
		if m == newM {
			break
		}
		m = newM

	}
	// classInQ: отображает классовый/корневой Q на тот q, из которого мы дошли до корневого
	// classIndToQ отображает индекс корневого Q, который будет в новой матрице, на индекс, который был в старой
	// QToClassInd по старому индексу находит новый корневой индекс (обратная к classIndToQ)
	classIndQ := map[int]int{}
	classIndToQ := map[int]int{}
	QToClassInd := map[int]int{}
	classCnt := 0
	for q := 0; q < numQ; q++ {
		classQ := pi[q]
		if _, ok := classIndQ[classQ]; ok {
			continue
		}
		classIndToQ[classCnt] = classQ
		classIndQ[classQ] = q
		QToClassInd[classQ] = classCnt
		classCnt += 1
	}

	for classInd := 0; classInd < classCnt; classInd++ {
		newQInToQ = append(newQInToQ, []int{})
		newQInToOut = append(newQInToOut, []string{})
		for x := 0; x < numIn; x++ {
			q := classIndQ[classIndToQ[classInd]]
			newQInToQ[classInd] = append(newQInToQ[classInd], QToClassInd[pi[QInToQ[q][x]]])
			newQInToOut[classInd] = append(newQInToOut[classInd], QInToOut[q][x])
		}
	}
	newStartQ = QToClassInd[pi[startQ]]
	newNumQ = m
	// fmt.Println(m, len(classIndQ), len(classIndToQ), len(QToClassInd))
	return newNumQ, newStartQ, newQInToQ, newQInToOut
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

func compareMealy(newNumQ1 int, numIn1 int, newQInToQ1 [][]int, newQInToOut1 [][]string, TimeToQ1 map[int]int, newQ1 map[int]int, newNumQ2 int, numIn2 int, newQInToQ2 [][]int, newQInToOut2 [][]string, TimeToQ2 map[int]int, newQ2 map[int]int) bool {
	if newNumQ1 != newNumQ2 || numIn1 != numIn2 {
		return false
	}
	for i := 0; i < newNumQ1; i++ {
		for j := 0; j < numIn1; j++ {
			if newQ1[newQInToQ1[TimeToQ1[i]][j]] != newQ2[newQInToQ2[TimeToQ2[i]][j]] || newQInToOut1[TimeToQ1[i]][j] != newQInToOut2[TimeToQ2[i]][j] {
				return false
			}
		}
	}
	return true
}
