# OpenMP-like (Numba) CG in Python

Install:
  python -m pip install numba numpy matplotlib

Run:
  OMP_NUM_THREADS=8 python cg_omp.py --N 40000000 --eps 1e-8 --matrix model --problem mix

Benchmark:
  python bench_omp.py --solver ./cg_omp.py --N 40000000 --eps 1e-8 --threads 1 2 4 8 16 --repeats 3 --matrix model --problem mix --outdir bench_omp_out
