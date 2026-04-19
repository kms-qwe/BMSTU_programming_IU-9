package main

import (
	"errors"
	"fmt"
	"math"
)

func main() {
	a := 0.0
	b := 1.0
	n := 10

	x, y, h := TabulateExp(a, b, n)

	ai, bi, ci, di, hi, err := BuildNaturalCubicSplineUniformAll(x, y)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("Cubic natural spline for f(x)=exp(x) on [%.2f, %.2f], n=%d, h=%.6f\n\n", a, b, n, h)

	fmt.Println("Nodes table (x_i, y_i=exp(x_i)):")
	fmt.Printf("%-3s %-12s %-16s\n", "i", "x_i", "y_i")
	for i := 0; i <= n; i++ {
		fmt.Printf("%-3d %-12.6f %-16.10f\n", i, x[i], y[i])
	}

	fmt.Println("\nSpline coefficients per interval i=0..n-1:")
	fmt.Printf("%-3s %-14s %-14s %-14s %-14s %-10s\n", "i", "a_i", "b_i", "c_i", "d_i", "h_i")
	for i := 0; i < n; i++ {
		fmt.Printf("%-3d %-14.10f %-14.10f %-14.10f %-14.10f %-10.6f\n",
			i, ai[i], bi[i], ci[i], di[i], hi[i])
	}

	fmt.Println("\nDelta at nodes: delta_i = S(x_i) - y_i (should be ~0)")
	fmt.Printf("%-3s %-12s %-16s %-16s %-16s\n", "i", "x_i", "y_i", "S(x_i)", "delta")
	maxAbsDelta := 0.0
	for i := 0; i <= n; i++ {
		s := SplineValueUniform(x, ai, bi, ci, di, x[i])
		delta := s - y[i]
		if math.Abs(delta) > maxAbsDelta {
			maxAbsDelta = math.Abs(delta)
		}
		fmt.Printf("%-3d %-12.6f %-16.10f %-16.10f %-16.3e\n", i, x[i], y[i], s, delta)
	}
	fmt.Printf("max |delta_i| = %.3e\n", maxAbsDelta)

	fmt.Println("\nDelta at midpoints: delta2_i = |S(x_{i-0.5}) - f(x_{i-0.5})|, i=1..n")
	fmt.Printf("%-3s %-12s %-16s %-16s %-16s %-10s\n", "i", "x_mid", "f(x_mid)", "S(x_mid)", "delta2", "rel%")
	maxRel := 0.0
	for i := 1; i <= n; i++ {
		xMid := 0.5 * (x[i-1] + x[i])
		fMid := math.Exp(xMid)
		sMid := SplineValueUniform(x, ai, bi, ci, di, xMid)
		delta2 := math.Abs(sMid - fMid)
		rel := delta2 / math.Abs(fMid) * 100.0
		if rel > maxRel {
			maxRel = rel
		}
		fmt.Printf("%-3d %-12.6f %-16.10f %-16.10f %-16.3e %-10.4f\n", i, xMid, fMid, sMid, delta2, rel)
	}
	fmt.Printf("max rel%% at midpoints = %.4f%%\n", maxRel)
}

func TabulateExp(a, b float64, n int) (x, y []float64, h float64) {
	h = (b - a) / float64(n)
	x = make([]float64, n+1)
	y = make([]float64, n+1)
	for i := 0; i <= n; i++ {
		x[i] = a + float64(i)*h
		y[i] = math.Exp(x[i])
	}
	return x, y, h
}

func BuildNaturalCubicSplineUniformAll(x, y []float64) (a, b, c, d, h []float64, err error) {
	if len(x) != len(y) {
		return nil, nil, nil, nil, nil, errors.New("x and y size mismatch")
	}
	if len(x) < 2 {
		return nil, nil, nil, nil, nil, errors.New("need at least 2 points")
	}

	n := len(x) - 1
	h0 := x[1] - x[0]
	if h0 <= 0 {
		return nil, nil, nil, nil, nil, errors.New("non-positive step")
	}

	for i := 1; i < n; i++ {
		hi := x[i+1] - x[i]
		if math.Abs(hi-h0) > 1e-9 {
			return nil, nil, nil, nil, nil, fmt.Errorf("grid is not uniform at i=%d: h=%.12g, hi=%.12g", i, h0, hi)
		}
	}

	if n == 1 {
		a = []float64{y[0]}
		b = []float64{(y[1] - y[0]) / h0}
		c = []float64{0, 0}
		d = []float64{0}
		h = []float64{h0}
		return a, b, c, d, h, nil
	}

	m := n - 1
	lower := make([]float64, m)
	diag := make([]float64, m)
	upper := make([]float64, m)
	rhs := make([]float64, m)

	for k := 0; k < m; k++ {
		i := k + 1
		diag[k] = 4.0
		if k > 0 {
			lower[k] = 1.0
		}
		if k < m-1 {
			upper[k] = 1.0
		}
		rhs[k] = (3.0 / (h0 * h0)) * (y[i+1] - 2.0*y[i] + y[i-1])
	}

	cInner, err := SolveTridiagonal(lower, diag, upper, rhs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	c = make([]float64, n+1)
	c[0] = 0
	c[n] = 0
	for k := 0; k < m; k++ {
		c[k+1] = cInner[k]
	}

	a = make([]float64, n)
	b = make([]float64, n)
	d = make([]float64, n)
	h = make([]float64, n)

	for i := 0; i < n; i++ {
		h[i] = x[i+1] - x[i]
		a[i] = y[i]
		b[i] = (y[i+1]-y[i])/h[i] - (h[i]/3.0)*(2.0*c[i]+c[i+1])
		d[i] = (c[i+1] - c[i]) / (3.0 * h[i])
	}

	return a, b, c, d, h, nil
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

func SplineValueUniform(x []float64, a, b []float64, cFull []float64, d []float64, xQuery float64) float64 {
	n := len(x) - 1
	if n <= 0 {
		return math.NaN()
	}

	h0 := x[1] - x[0]
	i := int(math.Floor((xQuery - x[0]) / h0))
	if i < 0 {
		i = 0
	}
	if i >= n {
		i = n - 1
	}

	t := xQuery - x[i]
	return a[i] + b[i]*t + cFull[i]*t*t + d[i]*t*t*t
}
