package main

import (
	"fmt"
	"math"
)

const (
	eps     = 1e-3
	maxIter = 1000
)

const (
	x1Exact = 5.0 / 13.0
	x2Exact = -5.0 / 26.0
)

func f(x1, x2 float64) float64 {
	return x1*x1 + 4*x1*x2 + 17*x2*x2 + 5*x2
}

func grad(x1, x2 float64) (float64, float64) {
	gx := 2*x1 + 4*x2
	gy := 4*x1 + 34*x2 + 5
	return gx, gy
}

func norm(gx, gy float64) float64 {
	if math.Abs(gx) > math.Abs(gy) {
		return math.Abs(gx)
	}
	return math.Abs(gy)
}

func step(x1, x2 float64) float64 {
	gx, gy := grad(x1, x2)

	phi1 := -(gx*gx + gy*gy)
	phi2 := 2*gx*gx + 8*gx*gy + 34*gy*gy

	return -phi1 / phi2
}

func main() {

	fmt.Println("АНАЛИТИЧЕСКОЕ РЕШЕНИЕ:")
	fmt.Printf("x* = (%.6f, %.6f)\n", x1Exact, x2Exact)
	fmt.Printf("f(x*) = %.6f\n", f(x1Exact, x2Exact))
	fmt.Println()

	x1 := 0.0
	x2 := 0.0

	fmt.Println("ИТЕРАЦИИ:")

	for k := 0; k < maxIter; k++ {
		gx, gy := grad(x1, x2)
		norm := norm(gx, gy)

		fmt.Printf("k=%d  x=(%.6f, %.6f)  f=%.6f  ||grad||=%.6f\n",
			k, x1, x2, f(x1, x2), norm)

		if norm < eps {
			break
		}

		t := step(x1, x2)

		x1 = x1 - t*gx
		x2 = x2 - t*gy
	}

	fmt.Println()
	fmt.Println("РЕЗУЛЬТАТ:")
	fmt.Printf("x ≈ (%.6f, %.6f)\n", x1, x2)
	fmt.Printf("f(x) ≈ %.6f\n", f(x1, x2))

	errX1 := math.Abs(x1 - x1Exact)
	errX2 := math.Abs(x2 - x2Exact)
	errNorm := math.Max(errX1, errX2)
	errF := math.Abs(f(x1, x2) - f(x1Exact, x2Exact))

	fmt.Println()
	fmt.Println("ПОГРЕШНОСТИ:")
	fmt.Printf("|x1 - x1*| = %.6e\n", errX1)
	fmt.Printf("|x2 - x2*| = %.6e\n", errX2)
	fmt.Printf("||x - x*|| = %.6e\n", errNorm)
	fmt.Printf("|f(x) - f(x*)| = %.6e\n", errF)

	fmt.Println()
	fmt.Println("СРАВНЕНИЕ:")
	fmt.Printf("Аналитически: (%.6f, %.6f)\n", x1Exact, x2Exact)
	fmt.Printf("Численно:     (%.6f, %.6f)\n", x1, x2)
}
