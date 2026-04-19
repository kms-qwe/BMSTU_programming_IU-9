package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
)

func SolveTridiagonal(lower, diag, upper, rhs []float64) ([]float64, error) {
	n := len(diag)
	if len(lower) != n || len(upper) != n || len(rhs) != n {
		return nil, errors.New("dimension mismatch in tridiagonal arrays")
	}
	if n == 0 {
		return nil, errors.New("empty system")
	}

	alpha := make([]float64, n)
	beta := make([]float64, n)

	eps := 1e-14
	if math.Abs(diag[0]) < eps {
		return nil, errors.New("zero/near-zero diagonal at i=0")
	}

	alpha[0] = -upper[0] / diag[0]
	beta[0] = rhs[0] / diag[0]

	for i := 1; i < n; i++ {
		den := lower[i]*alpha[i-1] + diag[i]
		if math.Abs(den) < eps {
			return nil, fmt.Errorf("zero/near-zero denominator at i=%d", i)
		}
		if i < n-1 {
			alpha[i] = -upper[i] / den
		} else {
			alpha[i] = 0
		}
		beta[i] = (rhs[i] - lower[i]*beta[i-1]) / den
	}

	x := make([]float64, n)
	x[n-1] = beta[n-1]
	for i := n - 2; i >= 0; i-- {
		x[i] = alpha[i]*x[i+1] + beta[i]
	}
	return x, nil
}

// Вариант 12:
// p = 0, q = 4, f(x) = e^(-2x), y(0)=0, y'(0)=0
//
// Точное решение задачи Коши:
// y(x) = (sin(2x) - cos(2x) + e^(-2x)) / 8
func exactY(x float64) float64 {
	return (math.Sin(2*x) - math.Cos(2*x) + math.Exp(-2*x)) / 8.0
}

func main() {
	nFlag := flag.Int("n", 10, "number of subintervals (n), grid has n+1 nodes")
	flag.Parse()

	n := *nFlag
	if n < 2 {
		fmt.Println("n must be >= 2")
		return
	}

	alphaBC := 0.0
	bBC := exactY(1.0)
	h := 1.0 / float64(n)
	m := n + 1

	lower := make([]float64, m)
	diag := make([]float64, m)
	upper := make([]float64, m)
	rhs := make([]float64, m)

	diag[0] = 1
	rhs[0] = alphaBC

	diag[n] = 1
	rhs[n] = bBC

	for i := 1; i <= n-1; i++ {
		xi := float64(i) * h

		p := 0.0
		q := 4.0
		f := math.Exp(-2 * xi)

		lower[i] = 1.0 - (h/2.0)*p
		diag[i] = h*h*q - 2.0
		upper[i] = 1.0 + (h/2.0)*p
		rhs[i] = h * h * f
	}

	yNum, err := SolveTridiagonal(lower, diag, upper, rhs)
	if err != nil {
		fmt.Println("solve error:", err)
		return
	}

	fmt.Printf("BVP: y'' + 4y = e^(-2x), y(0)=0, y(1)=b\n")
	fmt.Printf("n=%d, h=%.10g\n", n, h)
	fmt.Printf("b = y(1) = %.15g\n\n", bBC)

	fmt.Println("i\t x_i\t\t y_i (num)\t\t y(x_i) (exact)\t\t |err|")
	maxErr := 0.0
	maxI := 0
	for i := 0; i <= n; i++ {
		xi := float64(i) * h
		yEx := exactY(xi)
		errAbs := math.Abs(yEx - yNum[i])
		if errAbs > maxErr {
			maxErr = errAbs
			maxI = i
		}
		fmt.Printf("%d\t %.10g\t %.15g\t %.15g\t %.15g\n", i, xi, yNum[i], yEx, errAbs)
	}

	fmt.Printf("\n||y - y_num||_inf = max_i |y(x_i)-y_i| = %.15g (at i=%d)\n", maxErr, maxI)
}
