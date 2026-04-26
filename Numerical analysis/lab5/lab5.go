package main

import (
	"fmt"
	"math"
)

func main() {
	x := []float64{1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5}
	y := []float64{0.31, 0.57, 1.57, 1.40, 1.18, 1.54, 1.55, 1.96, 2.42}

	// xa := 3.000
	// xg := 2.236
	// xh := 1.666

	ya := 1.365
	yg := 0.866
	yh := 0.549

	z_xa := 1.296
	z_xg := 1.033
	z_xh := 0.872

	deltas := []float64{
		math.Abs(z_xa - ya), // δ1
		math.Abs(z_xg - yg), // δ2
		math.Abs(z_xa - yg), // δ3
		math.Abs(z_xg - ya), // δ4
		math.Abs(z_xh - ya), // δ5
		math.Abs(z_xa - yh), // δ6
		math.Abs(z_xh - yh), // δ7
		math.Abs(z_xh - yg), // δ8
		math.Abs(z_xg - yh), // δ9
	}

	minDelta := deltas[0]
	minIndex := 1

	for i := 1; i < len(deltas); i++ {
		if deltas[i] < minDelta {
			minDelta = deltas[i]
			minIndex = i + 1
		}
	}

	fmt.Println("Дельты:")
	for i, d := range deltas {
		fmt.Printf("δ%d = %.6f\n", i+1, d)
	}

	fmt.Printf("\nМинимальная δ%d = %.6f\n", minIndex, minDelta)

	// δ8 → z8(x) = a * e^(b/x)

	// линеаризация
	// ln(y) = ln(a) + b/x
	// Y = A + bX, где X = 1/x, A = ln(a)

	n := float64(len(x))

	var sumX, sumX2, sumY, sumXY float64

	for i := 0; i < len(x); i++ {
		X := 1.0 / x[i]
		Y := math.Log(y[i])

		sumX += X
		sumX2 += X * X
		sumY += Y
		sumXY += X * Y
	}

	// СЛАУ
	// A*n + b*ΣX = ΣY
	// A*ΣX + b*ΣX² = ΣXY

	det := n*sumX2 - sumX*sumX

	A := (sumY*sumX2 - sumX*sumXY) / det
	b := (n*sumXY - sumX*sumY) / det

	// делинеаризация
	a := math.Exp(A)

	fmt.Println("\nКоэффициенты:")
	fmt.Printf("a = %.6f\n", a)
	fmt.Printf("b = %.6f\n", b)

	fmt.Println("\nИтоговая функция:")
	fmt.Printf("z(x) = %.6f * e^(%.6f / x)\n", a, b)

	// среднеквадратичное отклонение
	var sumErr float64

	for i := 0; i < len(x); i++ {
		z := a * math.Exp(b/x[i])
		sumErr += math.Pow(z-y[i], 2)
	}

	fmt.Printf("\nСумма квадратов отклонений Δ = %.6f\n", sumErr)
}
