package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Matrix struct {
	n    int
	data []float64
}

func NewMatrix(n int) *Matrix               { return &Matrix{n: n, data: make([]float64, n*n)} }
func (m *Matrix) At(i, j int) float64       { return m.data[i*m.n+j] }
func (m *Matrix) AddAt(i, j int, v float64) { m.data[i*m.n+j] += v }
func (m *Matrix) FillRand(r *rand.Rand) {
	for i := range m.data {
		m.data[i] = r.Float64()*2 - 1
	}
}

// Базовый последовательный (row-major) для проверки корректности/замера
func MulRowMajor(A, B *Matrix) *Matrix {
	n := A.n
	C := NewMatrix(n)
	for i := 0; i < n; i++ {
		ci := C.data[i*n : i*n+n]
		aiBase := i * n
		for k := 0; k < n; k++ {
			aik := A.data[aiBase+k]
			bk := B.data[k*n : k*n+n]
			for j := 0; j < n; j++ {
				ci[j] += aik * bk[j]
			}
		}
	}
	return C
}

// Полезная параллельность: шардим строки C между воркерами + 3D-блокировка.
// Внутри тайлов используем порядок циклов i (мелкий) -> k (средний) -> j (внутренний, по непрерывным сегментам).
// Это даёт непрерывный доступ к A по k и к B по j, C пишем кусками по j.
func MulParallelBlocked(A, B *Matrix, workers, Ti, Tj, Tk int) *Matrix {
	n := A.n
	if workers < 1 {
		workers = 1
	}
	// Не плодим лишних горутин, ограничиваем реально доступными потоками
	if maxp := runtime.GOMAXPROCS(0); workers > maxp {
		workers = maxp
	}
	// Разумные дефолты блоков
	if Ti <= 0 {
		Ti = 32
	}
	if Tj <= 0 {
		Tj = 256
	}
	if Tk <= 0 {
		Tk = 128
	}

	C := NewMatrix(n)

	rowsPer := (n + workers - 1) / workers
	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		i0 := w * rowsPer
		i1 := i0 + rowsPer
		if i1 > n {
			i1 = n
		}
		if i0 >= i1 {
			wg.Done()
			continue
		}

		go func(i0, i1 int) {
			defer wg.Done()

			n := A.n
			Ad := A.data
			Bd := B.data
			Cd := C.data

			// Блокируем по i,k,j
			for ii := i0; ii < i1; ii += Ti {
				iEnd := ii + Ti
				if iEnd > i1 {
					iEnd = i1
				}
				for kk := 0; kk < n; kk += Tk {
					kEnd := kk + Tk
					if kEnd > n {
						kEnd = n
					}
					for i := ii; i < iEnd; i++ {
						ci := Cd[i*n : i*n+n] // вся строка C[i,:]
						aiBase := i * n
						// Пробегаем блок по k: A[i,k] читается по порядку, B[k,:] читаем строками
						for k := kk; k < kEnd; k++ {
							aik := Ad[aiBase+k]
							bk := Bd[k*n : k*n+n] // строка B[k,:] — непрерывная по j
							// Разбиваем по j-сегментам, чтобы держать C и B в L1/L2
							for jj := 0; jj < n; jj += Tj {
								jEnd := jj + Tj
								if jEnd > n {
									jEnd = n
								}
								cseg := ci[jj:jEnd]
								bseg := bk[jj:jEnd]
								// SAXPY: cseg += aik * bseg
								// Небольшая ручная "размотка" даёт плюс к perf в Go
								ln := jEnd - jj
								j := 0
								for ; j+3 < ln; j += 4 {
									cseg[j+0] += aik * bseg[j+0]
									cseg[j+1] += aik * bseg[j+1]
									cseg[j+2] += aik * bseg[j+2]
									cseg[j+3] += aik * bseg[j+3]
								}
								for ; j < ln; j++ {
									cseg[j] += aik * bseg[j]
								}
							}
						}
					}
				}
			}
		}(i0, i1)
	}
	wg.Wait()
	return C
}

func maxAbsDiff(X, Y *Matrix) float64 {
	var mx float64
	for i := range X.data {
		if d := math.Abs(X.data[i] - Y.data[i]); d > mx {
			mx = d
		}
	}
	return mx
}

func fmtSpeedup(base, cur time.Duration) string {
	if cur <= 0 {
		return "∞x"
	}
	return fmt.Sprintf("%.2fx", float64(base)/float64(cur))
}

func main() {
	n := flag.Int("n", 512, "размер квадратных матриц (n×n)")
	flag.Parse()

	fmt.Printf("GOMAXPROCS=%d\n", runtime.GOMAXPROCS(0))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	A := NewMatrix(*n)
	B := NewMatrix(*n)
	A.FillRand(r)
	B.FillRand(r)

	// База — последовательная row-major версия
	t0 := time.Now()
	Cseq := MulRowMajor(A, B)
	tSeq := time.Since(t0)

	fmt.Printf("n=%d\n", *n)
	fmt.Printf("[seq  row] time=%v\n", tSeq)

	// Параллельная блокированная версия (без транспонирования B)
	for _, w := range []int{2, 4, 8, 16, 32, 64} {
		t0 = time.Now()
		Cpar := MulParallelBlocked(A, B, w, 32, 256, 128)
		tPar := time.Since(t0)
		diff := maxAbsDiff(Cseq, Cpar)
		fmt.Printf("[par %3d] time=%v  speedup=%s  max|Δ|=%.3e\n",
			w, tPar, fmtSpeedup(tSeq, tPar), diff)
	}
}
