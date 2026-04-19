package main

import (
	"fmt"
	"math"
)

var (
	a       = -1.0
	b       = 2.0
	epsilon = 0.001
	Itrue   = 2 - 5/math.Exp(3)
)

func f(x float64) float64 {
	return x * x * math.Exp(x-2)
}

func midpoint(n int) float64 {
	h := (b - a) / float64(n)

	sum := 0.0
	for i := 0; i < n; i++ {
		x := a + (float64(i)+0.5)*h
		sum += f(x)
	}

	return h * sum
}

func trap(n int) float64 {
	h := (b - a) / float64(n)

	sum := (f(a) + f(b)) / 2

	for i := 1; i < n; i++ {
		x := a + float64(i)*h
		sum += f(x)
	}

	return h * sum
}

func simpson(n int) float64 {
	h := (b - a) / float64(n)

	sum := f(a) + f(b)

	for i := 1; i < n; i++ {
		x := a + float64(i)*h

		if i%2 == 0 {
			sum += 2 * f(x)
		} else {
			sum += 4 * f(x)
		}
	}

	return h * sum / 3
}

func richardson(Ih, I2h float64, k int) float64 {
	return (Ih - I2h) / (math.Pow(2, float64(k)) - 1)
}

type Result struct {
	n     int
	Istar float64
	R     float64
	Icorr float64
	Err   float64
}

func runMethod(method func(int) float64, order int) Result {

	n := 2
	prev := method(n)

	for {
		n *= 2

		Istar := method(n)
		R := richardson(Istar, prev, order)
		Icorr := Istar + R
		err := math.Abs(Icorr - Itrue)

		if math.Abs(R) < epsilon {
			return Result{
				n:     n,
				Istar: Istar,
				R:     R,
				Icorr: Icorr,
				Err:   err,
			}
		}

		prev = Istar
	}
}

func main() {

	fmt.Printf("Integral on [%.1f, %.1f]\n", a, b)
	fmt.Printf("epsilon = %.6f\n", epsilon)
	fmt.Printf("Analytical I = %.6f\n\n", Itrue)

	fmt.Printf("%-30s %-6s %-12s %-12s %-12s %-12s\n",
		"Метод", "n", "I*", "R", "I*+R", "|ΔI|")

	mid := runMethod(midpoint, 2)
	trap := runMethod(trap, 2)
	simp := runMethod(simpson, 4)

	fmt.Printf("%-30s %-6d %-12.6f %-12.6f %-12.6f %-12.6f\n",
		"Метод средних прямоугольников",
		mid.n, mid.Istar, mid.R, mid.Icorr, mid.Err)

	fmt.Printf("%-30s %-6d %-12.6f %-12.6f %-12.6f %-12.6f\n",
		"Метод трапеций",
		trap.n, trap.Istar, trap.R, trap.Icorr, trap.Err)

	fmt.Printf("%-30s %-6d %-12.6f %-12.6f %-12.6f %-12.6f\n",
		"Метод Симпсона",
		simp.n, simp.Istar, simp.R, simp.Icorr, simp.Err)
}
