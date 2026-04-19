package main

import (
	"fmt"
	"math"
)

const eps = 0.01

func f1(x, y float64) float64 {
	return math.Sin(y-1) + x - 1.3
}

func f2(x, y float64) float64 {
	return y - math.Sin(x+1) - 0.8
}

func df1dx(x, y float64) float64 {
	return 1
}

func df1dy(x, y float64) float64 {
	return math.Cos(y - 1)
}

func df2dx(x, y float64) float64 {
	return -math.Cos(x + 1)
}

func df2dy(x, y float64) float64 {
	return 1
}

func main() {
	x := 0.57
	y := 1.8

	xExact := 0.589844
	yExact := 1.799383

	fmt.Println("Начальное приближение:")
	fmt.Printf("x0 = %.5f, y0 = %.5f\n\n", x, y)

	iter := 0

	for {
		F1 := f1(x, y)
		F2 := f2(x, y)

		a := df1dx(x, y)
		b := df1dy(x, y)
		c := df2dx(x, y)
		d := df2dy(x, y)

		det := a*d - b*c

		if math.Abs(det) < 1e-10 {
			fmt.Println("Якобиан вырожден")
			return
		}

		dx := (-F1*d + b*F2) / det
		dy := (-a*F2 + F1*c) / det

		xNew := x + dx
		yNew := y + dy

		err := math.Max(math.Abs(xNew-x), math.Abs(yNew-y))

		iter++
		fmt.Printf("Итерация %d:\n", iter)
		fmt.Printf("  x = %.5f, y = %.5f\n", xNew, yNew)
		fmt.Printf("  dx = %.5f, dy = %.5f\n", dx, dy)
		fmt.Printf("  |Δ| = %.5f\n", err)
		fmt.Printf("  f1 = %.5f, f2 = %.5f\n\n", F1, F2)

		if err < eps {
			fmt.Println("Критерий остановки выполнен:")
			fmt.Printf("|Δ| = %.5f < eps = %.5f\n\n", err, eps)

			fmt.Println("Ответ:")
			fmt.Printf("x = %.5f, y = %.5f\n", xNew, yNew)

			errX := math.Abs(xNew - xExact)
			errY := math.Abs(yNew - yExact)

			fmt.Println("\nАбсолютная погрешность относительно точного решения:")
			fmt.Printf("|x - xExact| = %.6f\n", errX)
			fmt.Printf("|y - yExact| = %.6f\n", errY)

			break
		}

		x = xNew
		y = yNew
	}
}
