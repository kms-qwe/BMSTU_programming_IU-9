/*
реализация метода Гаусса для неоднородных слау n * n
*/

package main

import (
	"fmt"
	"math/big"
)

var (
	zero = big.NewRat(0, 1)
	one  = big.NewRat(1, 1)
)

func main() {
	matrix := initMatrix()
	gauss(matrix)

}

// Возвращает первую строку с индексом >=j, у которой jый элемент != 0
func nonZeroRow(j int, matrix [][]big.Rat) int {
	for i := j; i < len(matrix); i++ {
		if matrix[i][j].Cmp(zero) != 0 {
			return i
		}
	}
	return -1
}

func swapRow(i, j int, matrix [][]big.Rat) {
	matrix[i], matrix[j] = matrix[j], matrix[i]
}
func normalizeColomn(j int, matrix [][]big.Rat) string {
	rowsIndex := nonZeroRow(j, matrix)
	if rowsIndex == -1 {
		return "can't normalize colomn"
	}
	if rowsIndex != j {
		swapRow(j, rowsIndex, matrix)
	}
	multiplyRow(big.NewRat(0, 1).Quo(one, &matrix[j][j]), j, matrix)
	for i := 0; i < len(matrix); i++ {
		if i == j || matrix[i][j].Cmp(zero) == 0 {
			continue
		}
		multiplyAndAddRow(i, j, *big.NewRat(0, 1).Neg(&matrix[i][j]), matrix)
	}
	return ""
}
func multiplyRow(num *big.Rat, row int, matrix [][]big.Rat) {
	for j := 0; j < len(matrix)+1; j++ {
		matrix[row][j].Mul(&matrix[row][j], num)
	}
}
func multiplyAndAddRow(dest int, source int, factor big.Rat, matrix [][]big.Rat) {
	for j := 0; j < len(matrix)+1; j++ {
		matrix[dest][j].Add(&matrix[dest][j], big.NewRat(0, 1).Mul(&matrix[source][j], &factor))
	}
}
func gauss(matrix [][]big.Rat) {
	for j := 0; j < len(matrix); j++ {
		// printMatrix(matrix, "gauss")
		err := normalizeColomn(j, matrix)
		if err != "" {
			fmt.Println("No solution")
			return
		}
	}
	for i := 0; i < len(matrix); i++ {
		fmt.Println(&matrix[i][len(matrix[i])-1])
	}
}
func printMatrix(matrix [][]big.Rat, s string) {
	fmt.Printf("------------------- %s\n", s)
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix)+1; j++ {
			if j == len(matrix) {
				fmt.Printf("| %s\n", &matrix[i][j])
			} else {
				fmt.Printf("%s ", &matrix[i][j])
			}

		}
	}
}
func initMatrix() [][]big.Rat {
	n := 0
	fmt.Scan(&n)
	matrix := make([][]big.Rat, n)
	for i := range matrix {
		matrix[i] = make([]big.Rat, n+1)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n+1; j++ {
			var val int
			fmt.Scan(&val)
			matrix[i][j] = *big.NewRat(int64(val), 1)
		}
	}
	return matrix
}
