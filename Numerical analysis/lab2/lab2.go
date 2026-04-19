package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
)

func exactY(x float64) float64 {
	// Вариант 12:
	// y(x) = (sin(2x) - cos(2x) + e^(-2x)) / 8
	return (math.Sin(2*x) - math.Cos(2*x) + math.Exp(-2*x)) / 8.0
}

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

	const eps = 1e-14
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

func solveByProgonka(n int, a, b, A, B float64) ([]float64, error) {
	h := (b - a) / float64(n)
	m := n + 1

	lower := make([]float64, m)
	diag := make([]float64, m)
	upper := make([]float64, m)
	rhs := make([]float64, m)

	diag[0] = 1.0
	rhs[0] = A

	diag[n] = 1.0
	rhs[n] = B

	for i := 1; i <= n-1; i++ {
		xi := a + float64(i)*h

		p := 0.0
		q := 4.0
		f := math.Exp(-2 * xi)

		lower[i] = 1.0 - (h/2.0)*p
		diag[i] = h*h*q - 2.0
		upper[i] = 1.0 + (h/2.0)*p
		rhs[i] = h * h * f
	}

	return SolveTridiagonal(lower, diag, upper, rhs)
}

func solveByShooting(n int, a, b, A, B float64) ([]float64, error) {
	h := (b - a) / float64(n)

	y0 := make([]float64, n+1)
	y1 := make([]float64, n+1)
	y := make([]float64, n+1)

	// Базовые начальные значения
	D0 := A
	D1 := h

	y0[0] = A
	y0[1] = D0

	y1[0] = 0.0
	y1[1] = D1

	for i := 1; i <= n-1; i++ {
		xi := a + float64(i)*h

		p := 0.0
		q := 4.0
		f := math.Exp(-2 * xi)

		den := 1.0 + p*h/2.0
		if math.Abs(den) < 1e-14 {
			return nil, fmt.Errorf("zero denominator at i=%d", i)
		}

		y0[i+1] = (f*h*h + (2.0-q*h*h)*y0[i] - (1.0-p*h/2.0)*y0[i-1]) / den
		y1[i+1] = ((2.0-q*h*h)*y1[i] - (1.0-p*h/2.0)*y1[i-1]) / den
	}

	if math.Abs(y1[n]) < 1e-14 {
		return nil, errors.New("cannot compute C1: y1[n] is zero")
	}

	C1 := (B - y0[n]) / y1[n]

	for i := 0; i <= n; i++ {
		y[i] = y0[i] + C1*y1[i]
	}

	return y, nil
}

func main() {
	nFlag := flag.Int("n", 10, "число разбиений")
	flag.Parse()

	n := *nFlag
	if n < 2 {
		fmt.Println("n must be >= 2")
		return
	}

	a := 0.0
	b := 1.0

	// Для сравнения двух методов приводим задачу к краевой:
	// y(0) = exactY(0), y(1) = exactY(1)
	A := exactY(a)
	B := exactY(b)

	yProg, err := solveByProgonka(n, a, b, A, B)
	if err != nil {
		fmt.Println("Ошибка метода прогонки:", err)
		return
	}

	yShot, err := solveByShooting(n, a, b, A, B)
	if err != nil {
		fmt.Println("Ошибка метода стрельбы:", err)
		return
	}

	h := (b - a) / float64(n)

	fmt.Println("ЛР2. Вариант 12")
	fmt.Println("y'' + 4y = e^(-2x)")
	fmt.Println("Сравнение методов: прогонка и стрельба")
	fmt.Printf("Граничные условия для сравнения: y(0)=%.8f, y(1)=%.8f\n\n", A, B)

	fmt.Printf("%12s%18s%18s%18s%16s%16s\n",
		"x", "Точное", "Прогонка", "Стрельба", "|y-y*|1", "|y-y*|2")

	maxErrProg := 0.0
	maxErrShot := 0.0

	for i := 0; i <= n; i++ {
		xi := a + float64(i)*h
		yEx := exactY(xi)

		errProg := math.Abs(yEx - yProg[i])
		errShot := math.Abs(yEx - yShot[i])

		if errProg > maxErrProg {
			maxErrProg = errProg
		}
		if errShot > maxErrShot {
			maxErrShot = errShot
		}

		fmt.Printf("%12.8f%18.8f%18.8f%18.8f%16.8f%16.8f\n",
			xi, yEx, yProg[i], yShot[i], errProg, errShot)
	}

	fmt.Printf("\nМакс ошибка (прогонка) = %.8f\n", maxErrProg)
	fmt.Printf("Макс ошибка (стрельба) = %.8f\n", maxErrShot)
}
