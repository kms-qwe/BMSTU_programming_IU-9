package main

import (
	"flag"
	"fmt"
	"math"
)

// Тестовый вариант:
// y” + 5y' - 3y = 3e^x,  x in [0, 1]
// y(0) = 1, y(1) = e
func exactY(x float64) float64 {
	return math.Exp(x)
}

func main() {
	nFlag := flag.Int("n", 10, "number of subintervals")
	flag.Parse()

	n := *nFlag
	if n < 2 {
		fmt.Println("n must be >= 2")
		return
	}

	a := 0.0
	b := 1.0
	A := 1.0
	B := math.E

	h := (b - a) / float64(n)

	D0 := A
	D1 := h

	y0 := make([]float64, n+1)
	y1 := make([]float64, n+1)
	y := make([]float64, n+1)

	y0[0] = A
	y0[1] = D0

	y1[0] = 0.0
	y1[1] = D1

	for i := 1; i <= n-1; i++ {
		xi := a + float64(i)*h

		pi := 5.0
		qi := -3.0
		fi := 3.0 * math.Exp(xi)

		den := 1.0 + pi*h/2.0
		if math.Abs(den) < 1e-14 {
			fmt.Printf("zero denominator at i=%d\n", i)
			return
		}

		y0[i+1] = (fi*h*h + (2.0-qi*h*h)*y0[i] - (1.0-pi*h/2.0)*y0[i-1]) / den
		y1[i+1] = ((2.0-qi*h*h)*y1[i] - (1.0-pi*h/2.0)*y1[i-1]) / den
	}

	if math.Abs(y1[n]) < 1e-14 {
		fmt.Println("cannot compute C1: y1[n] is zero")
		return
	}

	C1 := (B - y0[n]) / y1[n]

	for i := 0; i <= n; i++ {
		y[i] = y0[i] + C1*y1[i]
	}

	fmt.Println("Method: shooting method")
	fmt.Println("Problem: y'' + 5y' - 3y = 3e^x, y(0)=1, y(1)=e")
	fmt.Printf("n = %d, h = %.10g\n", n, h)
	fmt.Printf("A = %.15g\n", A)
	fmt.Printf("B = %.15g\n", B)
	fmt.Printf("D0 = %.15g\n", D0)
	fmt.Printf("D1 = %.15g\n", D1)
	fmt.Printf("C1 = %.15g\n\n", C1)

	fmt.Printf("%-4s %-12s %-20s %-20s %-20s %-20s %-20s\n",
		"i", "x_i", "y0[i]", "y1[i]", "y[i]", "y_exact(x_i)", "|err|")

	maxErr := 0.0
	maxI := 0

	for i := 0; i <= n; i++ {
		xi := a + float64(i)*h
		yExact := exactY(xi)
		err := math.Abs(yExact - y[i])

		if err > maxErr {
			maxErr = err
			maxI = i
		}

		fmt.Printf("%-4d %-12.6f %-20.12f %-20.12f %-20.12f %-20.12f %-20.12e\n",
			i, xi, y0[i], y1[i], y[i], yExact, err)
	}

	fmt.Printf("\n||y - y_num||_inf = %.15g (at i=%d)\n", maxErr, maxI)
}
