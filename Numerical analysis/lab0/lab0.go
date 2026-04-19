package main

import (
	"errors"
	"fmt"
	"math"
)

func main() {
	A := [][]float64{
		{4, 1, 0, 0},
		{1, 4, 1, 0},
		{0, 1, 4, 1},
		{0, 0, 1, 4},
	}

	d := []float64{5, 6, 6, 5}

	x, err := SolveTridiagonal(A, d)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("x:")
	for i, v := range x {
		fmt.Printf("x[%d] = %.10f\n", i+1, v)
	}
}

func SolveTridiagonal(A [][]float64, d []float64) ([]float64, error) {
	n := len(A)
	if n == 0 {
		return nil, errors.New("empty matrix")
	}
	if len(d) != n {
		return nil, fmt.Errorf("dimension mismatch: A is %dx%d, d is %d", n, n, len(d))
	}
	for i := 0; i < n; i++ {
		if len(A[i]) != n {
			return nil, fmt.Errorf("matrix is not square: row %d has length %d, expected %d", i, len(A[i]), n)
		}
	}

	a := make([]float64, n)
	b := make([]float64, n)
	c := make([]float64, n)

	for i := 0; i < n; i++ {
		b[i] = A[i][i]
		if i > 0 {
			a[i] = A[i][i-1]
		}
		if i < n-1 {
			c[i] = A[i][i+1]
		}
	}

	if err := validateTridiagonalAndConditions(A, a, b, c); err != nil {
		return nil, err
	}

	alpha := make([]float64, n)
	beta := make([]float64, n)

	if b[0] == 0 {
		return nil, errors.New("condition failed: b1 != 0 (b[0] == 0)")
	}

	alpha[0] = -c[0] / b[0]
	beta[0] = d[0] / b[0]

	eps := 1e-12
	for i := 1; i < n; i++ {
		den := a[i]*alpha[i-1] + b[i]
		if math.Abs(den) < eps {
			return nil, fmt.Errorf("zero/near-zero denominator at i=%d: a[i]*alpha[i-1] + b[i] = %.6g", i+1, den)
		}
		if i < n-1 {
			alpha[i] = -c[i] / den
		} else {
			alpha[i] = 0
		}
		beta[i] = (d[i] - a[i]*beta[i-1]) / den
	}

	x := make([]float64, n)
	x[n-1] = beta[n-1]
	for i := n - 2; i >= 0; i-- {
		x[i] = alpha[i]*x[i+1] + beta[i]
	}

	return x, nil
}

func validateTridiagonalAndConditions(A [][]float64, a, b, c []float64) error {
	n := len(A)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if absInt(i-j) > 1 && A[i][j] != 0 {
				return fmt.Errorf("matrix is not tridiagonal: A[%d][%d]=%v (must be 0)", i+1, j+1, A[i][j])
			}
		}
	}

	if b[0] == 0 {
		return errors.New("condition failed: b1 != 0")
	}

	return nil
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
