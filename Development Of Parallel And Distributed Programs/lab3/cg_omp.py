#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
OpenMP-like (Numba) Conjugate Gradient solver for Ax = b.
Matrix kinds: model (A=J+I), dense (for small N).
Problems: ones, sine, mix.
"""
from __future__ import annotations
import argparse, math, numpy as np

try:
    from numba import njit, prange, get_num_threads
    NUMBA_AVAILABLE = True
except Exception:
    NUMBA_AVAILABLE = False

def require_numba():
    if not NUMBA_AVAILABLE:
        raise RuntimeError("Requires 'numba'. Install via: python -m pip install numba")

if NUMBA_AVAILABLE:
    @njit(parallel=True, fastmath=True)
    def dot_parallel(x, y):
        s = 0.0
        n = x.shape[0]
        for i in prange(n):
            s += x[i] * y[i]
        return s

    @njit(parallel=True, fastmath=True)
    def norm2_parallel(x):
        return math.sqrt(dot_parallel(x, x))

    @njit(parallel=True, fastmath=True)
    def matvec_model_parallel(z, y_out):
        s = 0.0
        n = z.shape[0]
        for i in prange(n):
            s += z[i]
        for i in prange(n):
            y_out[i] = s + z[i]

    @njit(parallel=True, fastmath=True)
    def axpy_parallel(a, x, y):
        n = x.shape[0]
        for i in prange(n):
            y[i] += a * x[i]

    @njit(parallel=True, fastmath=True)
    def add_scaled_parallel(a, x, y, out):
        n = x.shape[0]
        for i in prange(n):
            out[i] = x[i] + a * y[i]

def build_dense(N:int)->np.ndarray:
    A = np.ones((N,N), dtype=np.float64)
    np.fill_diagonal(A, 2.0)
    return A

def _matvec_dense_numba(A, z, y_out):
    require_numba()
    _matvec_dense_jit(A, z, y_out)

if NUMBA_AVAILABLE:
    @njit(parallel=True, fastmath=True)
    def _matvec_dense_jit(A, z, y_out):
        n = z.shape[0]
        for i in prange(n):
            s = 0.0
            Ai = A[i]
            for j in range(n):
                s += Ai[j] * z[j]
            y_out[i] = s

def build_problem(N:int, matrix_kind:str, problem:str):
    b = np.zeros(N, dtype=np.float64)
    u_true = None
    if problem == "ones":
        b.fill(float(N+1))
        u_true = np.ones(N, dtype=np.float64)
    elif problem in ("sine","mix"):
        u = np.empty(N, dtype=np.float64)
        two_pi_over_N = 2.0 * math.pi / float(N)
        c = 0.12345 if problem == "mix" else 0.0
        for i in range(N):
            u[i] = math.sin(two_pi_over_N * float(i)) + c
        u_true = u.copy()
        if matrix_kind == "model":
            s = float(np.sum(u))
            b = s + u
        else:
            A = build_dense(N)
            b = A @ u
    else:
        raise ValueError("Unknown --problem")
    return b, u_true

def cg_solve(N:int, eps:float, max_it:int, matrix_kind:str, problem:str, verbose:bool):
    require_numba()
    b, _ = build_problem(N, matrix_kind, problem)
    x = np.zeros(N, dtype=np.float64)
    r = b.copy()
    z = r.copy()
    Az = np.zeros_like(r)

    bnorm = norm2_parallel(b)
    if bnorm == 0.0:
        bnorm = 1.0

    import time
    t0 = time.perf_counter()
    it = 0
    converged = False

    A_dense = None
    if matrix_kind == "dense":
        A_dense = build_dense(N)

    while it < max_it:
        if matrix_kind == "model":
            matvec_model_parallel(z, Az)
        else:
            _matvec_dense_numba(A_dense, z, Az)

        rr = dot_parallel(r, r)
        zAz = dot_parallel(z, Az)
        if zAz == 0.0:
            break
        alpha = rr / zAz

        axpy_parallel(alpha, z, x)
        axpy_parallel(-alpha, Az, r)

        rnorm = norm2_parallel(r)
        relres = rnorm / bnorm
        it += 1
        if verbose and (it % 5 == 0 or relres < eps):
            print(f"[iter {it}] relres={relres:.3e}", flush=True)
        if relres < eps:
            converged = True
            break

        rr_new = dot_parallel(r, r)
        beta = (rr_new / rr) if rr != 0.0 else 0.0
        add_scaled_parallel(beta, r, z, z)

    t1 = time.perf_counter()
    elapsed = t1 - t0
    final_relres = norm2_parallel(r) / bnorm

    threads = None
    if NUMBA_AVAILABLE:
        try:
            threads = get_num_threads()
        except Exception:
            threads = None

    print(f"[CG-OMP-Py] matrix={matrix_kind} problem={problem} N={N} threads={threads} iters={it} relres={final_relres:.3e} time={elapsed:.4f}s")
    print("x[:5] ≈", " ".join(f"{v:.6f}" for v in x[:min(5, N)]))

    print("JSON_RESULT_BEGIN")
    import json as _json
    print(_json.dumps({
        "N": N, "eps": eps, "iterations": it, "converged": bool(converged),
        "rel_residual": float(final_relres), "time_sec": float(elapsed),
        "matrix": matrix_kind, "problem": problem, "threads": threads
    }))
    print("JSON_RESULT_END")

def main():
    import argparse
    parser = argparse.ArgumentParser(description="Numba-parallel (OpenMP-like) CG solver (Python)")
    parser.add_argument("--N", type=int, required=True)
    parser.add_argument("--eps", type=float, default=1e-6)
    parser.add_argument("--max-it", type=int, default=-1)
    parser.add_argument("--matrix", choices=["model","dense"], default="model")
    parser.add_argument("--problem", choices=["ones","sine","mix"], default="ones")
    parser.add_argument("--verbose", action="store_true")
    args = parser.parse_args()
    if args.max_it < 0:
        args.max_it = max(10, min(args.N, 5*args.N))
    cg_solve(args.N, args.eps, args.max_it, args.matrix, args.problem, args.verbose)

if __name__ == "__main__":
    main()
