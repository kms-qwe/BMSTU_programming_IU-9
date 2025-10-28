#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
MPI Conjugate Gradient (CG) solver for Ax = b

Two variants:
  --variant dup  : vectors x and b are duplicated on every process
  --variant dist : vectors x and b are distributed across processes (like A)

Matrix options:
  --matrix model : implicit SPD matrix from the lab handout:
                   A has 2.0 on the diagonal and 1.0 elsewhere (A = J + I).
                   Matvec is computed without storing A: (A v) = sum(v) * 1 + v
  --matrix dense : (for small N) local row-block of a dense matrix is stored.
                   Use only to validate correctness; memory O(N^2).

Problems:
  --problem ones : "model problem with known solution": b = (N+1)*1, x* = 1
  --problem sine : "model problem with arbitrary solution":
                   u_i = sin(2π i / N), b = A u, so solution is u

Stopping criterion:
    ||r||_2 / ||b||_2 < eps
"""
from __future__ import annotations
import numpy as np
from mpi4py import MPI
import argparse
import math

def distribute_rows(N:int, size:int, rank:int):
    """Return (start, end, counts, displs) for contiguous row-block decomposition."""
    q, r = divmod(N, size)
    counts = np.array([q + (1 if p < r else 0) for p in range(size)], dtype=np.int64)
    displs = np.zeros(size, dtype=np.int64)
    displs[1:] = np.cumsum(counts[:-1])
    start = int(displs[rank])
    end = int(start + counts[rank])
    return start, end, counts, displs

def build_vectors(N:int, problem:str, start:int, end:int, comm:MPI.Comm, variant:str):
    """Return (x_local, b_local, maybe also x_full, b_full depending on variant)."""
    # Initial guess x0 = 0
    x_local = np.zeros(end - start, dtype=np.float64)

    if problem == "ones":
        # b_i = N + 1 for all i
        b_local = np.full(end - start, float(N + 1), dtype=np.float64)
    elif problem == "sine":
        # u_i = sin(2π i / N), b = A u (for model A): (A u) = (sum u) * 1 + u
        idx = np.arange(start, end, dtype=np.float64)
        u_local = np.sin(2.0 * np.pi * idx / float(N))
        sum_u_local = np.array([u_local.sum()], dtype=np.float64)
        sum_u_global = np.array([0.0], dtype=np.float64)
        comm.Allreduce(sum_u_local, sum_u_global, op=MPI.SUM)
        b_local = sum_u_global[0] + u_local
    else:
        raise ValueError("Unknown problem: choose 'ones' or 'sine'")

    if variant == "dup":
        # Gather b to all processes (duplicate)
        _, _, counts, displs = distribute_rows(N, comm.Get_size(), comm.Get_rank())
        b_full = np.zeros(N, dtype=np.float64)
        comm.Allgatherv(b_local, [b_full,
                                  counts.tolist(), displs.tolist(), MPI.DOUBLE])
        x_full = np.zeros_like(b_full)  # x0 = 0, duplicated
        return x_local, b_local, x_full, b_full
    else:
        return x_local, b_local, None, None

def build_dense_local_A(N:int, start:int, end:int):
    """Build local row-block of dense 'model' matrix with 2.0 on diag and 1.0 elsewhere.
       WARNING: O(N * local_rows) memory; only for small N!"""
    rows = end - start
    A_local = np.ones((rows, N), dtype=np.float64)
    for i in range(rows):
        A_local[i, start + i] = 2.0
    return A_local

def matvec_model(z_local:np.ndarray, comm:MPI.Comm) -> (np.ndarray, float):
    """Compute y_local = A z for 'model' matrix: A = J + I -> A z = sum(z) * 1 + z."""
    sum_local = np.array([z_local.sum()], dtype=np.float64)
    sum_global = np.array([0.0], dtype=np.float64)
    comm.Allreduce(sum_local, sum_global, op=MPI.SUM)
    y_local = sum_global[0] + z_local
    return y_local, sum_global[0]

def matvec_dense(z_local:np.ndarray, A_local:np.ndarray, counts, displs, comm:MPI.Comm,
                 need_full_z:bool) -> np.ndarray:
    """Compute y_local = A_local @ z, where A_local are local rows of the global dense A.
       If need_full_z==True, gather z to a full vector on every rank (for distributed variant)."""
    if need_full_z:
        # SAFER: zeros + python lists for counts/displs
        N = int(np.sum(counts))
        z_full = np.zeros(N, dtype=np.float64)
        comm.Allgatherv(z_local, [z_full,
                                  counts.tolist() if hasattr(counts, "tolist") else counts,
                                  displs.tolist() if hasattr(displs, "tolist") else displs,
                                  MPI.DOUBLE])
        with np.errstate(divide='ignore', over='ignore', invalid='ignore'):
            return A_local @ z_full
    else:
        # here v is already "full" if variant==dup (we pass full vector)
        with np.errstate(divide='ignore', over='ignore', invalid='ignore'):
            return A_local @ z_local

def dot_global(x_local:np.ndarray, y_local:np.ndarray, comm:MPI.Comm) -> float:
    """Global dot product of distributed vectors."""
    local = np.array([float(np.dot(x_local, y_local))], dtype=np.float64)
    glob = np.array([0.0], dtype=np.float64)
    comm.Allreduce(local, glob, op=MPI.SUM)
    return float(glob[0])

def norm_global(x_local:np.ndarray, comm:MPI.Comm) -> float:
    """Global 2-norm of a distributed vector."""
    return math.sqrt(dot_global(x_local, x_local, comm))

def cg_solve(N:int, eps:float, max_it:int, variant:str, matrix_kind:str, problem:str,
             comm:MPI.Comm):
    """Run CG and return a result dict on rank 0."""
    rank = comm.Get_rank()
    size = comm.Get_size()

    start, end, counts, displs = distribute_rows(N, size, rank)

    # Build vectors (and duplicates if needed)
    x_local, b_local, x_full, b_full = build_vectors(N, problem, start, end, comm, variant)

    # Prepare matvec operator
    A_local = None
    if matrix_kind == "dense":
        A_local = build_dense_local_A(N, start, end)

    if matrix_kind == "model":
        def A_times(v_local):
            y_loc, _ = matvec_model(v_local, comm)
            return y_loc
        need_full_z = False
    elif matrix_kind == "dense":
        if variant == "dup":
            need_full_z = False
        else:
            need_full_z = True
        def A_times(v_local):
            return matvec_dense(v_local, A_local, counts, displs, comm, need_full_z)
    else:
        raise ValueError("Unknown --matrix option")

    # Initialize residuals: r0 = b - A x0 ; x0 = 0 -> r0 = b ; z0 = r0
    r_local = b_local.copy()
    z_local = r_local.copy()

    if variant == "dup":
        if b_full is None:
            b_full = np.zeros(N, dtype=np.float64)
            comm.Allgatherv(b_local, [b_full, counts.tolist(), displs.tolist(), MPI.DOUBLE])
        x_local = None  # use x_full as the source of truth

    # Precompute ||b||
    bnorm = norm_global(b_local, comm)
    if bnorm == 0.0:
        bnorm = 1.0

    # Prepare z_full for dense/dup
    if matrix_kind == "dense" and variant == "dup":
        z_full = np.zeros(N, dtype=np.float64)
        comm.Allgatherv(z_local, [z_full, counts.tolist(), displs.tolist(), MPI.DOUBLE])
    else:
        z_full = None

    comm.Barrier()
    t0 = MPI.Wtime()

    it = 0
    converged = False
    while it < max_it:
        # Compute A z
        if matrix_kind == "dense" and variant == "dup":
            Az_local = A_times(z_full)  # expects full vector
        else:
            Az_local = A_times(z_local)

        rr = dot_global(r_local, r_local, comm)
        zAz = dot_global(z_local, Az_local, comm)
        if zAz == 0.0:
            break
        alpha = rr / zAz

        # x_{k+1} update
        if variant == "dup":
            s, e = start, end
            x_full[s:e] += alpha * z_local
        else:
            x_local += alpha * z_local

        # r_{k+1}
        r_local = r_local - alpha * Az_local

        # Check convergence
        rnorm = norm_global(r_local, comm)
        relres = rnorm / bnorm
        it += 1

        rr_new = dot_global(r_local, r_local, comm)
        beta = (rr_new / rr) if rr != 0.0 else 0.0

        # z_{k+1}
        z_local = r_local + beta * z_local

        # Maintain z_full for dense/dup
        if matrix_kind == "dense" and variant == "dup":
            comm.Allgatherv(z_local, [z_full, counts.tolist(), displs.tolist(), MPI.DOUBLE])

        if relres < eps:
            converged = True
            break

        # Keep duplicated x in sync (dup)
        if variant == "dup":
            chunk = x_full[start:end].copy()
            comm.Allgatherv(chunk, [x_full, counts.tolist(), displs.tolist(), MPI.DOUBLE])

    comm.Barrier()
    t1 = MPI.Wtime()

    if variant == "dup":
        chunk = x_full[start:end].copy()
        comm.Allgatherv(chunk, [x_full, counts.tolist(), displs.tolist(), MPI.DOUBLE])
        final_relres = norm_global(r_local, comm) / bnorm
        result = {
            "N": N, "eps": eps, "iterations": it, "converged": bool(converged),
            "rel_residual": float(final_relres), "time_sec": float(t1 - t0),
            "variant": variant, "matrix": matrix_kind, "problem": problem, "proc_size": size,
        }
        x_preview = x_full[:min(5, N)].copy()
    else:
        final_relres = norm_global(r_local, comm) / bnorm
        result = {
            "N": N, "eps": eps, "iterations": it, "converged": bool(converged),
            "rel_residual": float(final_relres), "time_sec": float(t1 - t0),
            "variant": variant, "matrix": matrix_kind, "problem": problem, "proc_size": size,
        }
        x_preview = None

    # Pretty print on rank 0
    if rank == 0:
        if x_preview is None:
            # gather first 5 elements just for preview
            k = min(5, N)
            preview = np.zeros(k, dtype=np.float64)
            _, _, counts_all, displs_all = distribute_rows(N, size, rank)
            for g in range(k):
                owner = int(np.searchsorted((displs_all + counts_all), g, side='right'))
                owner = owner if (displs_all[owner] <= g < displs_all[owner] + counts_all[owner]) else owner-1
                offset = g - displs_all[owner]
                if rank == owner:
                    val = (x_local[offset]).copy()
                else:
                    val = None
                val = comm.bcast(val, root=owner)
                preview[g] = val
            x_preview = preview

        result_str = (f"[CG] variant={variant} matrix={matrix_kind} problem={problem} "
                      f"N={N} procs={size} iters={it} relres={result['rel_residual']:.3e} "
                      f"time={result['time_sec']:.4f}s")
        print(result_str, flush=True)
        print("x[:5] ≈", " ".join(f"{v:.6f}" for v in x_preview), flush=True)
        print("JSON_RESULT_BEGIN")
        import json as _json
        print(_json.dumps(result), flush=True)
        print("JSON_RESULT_END", flush=True)

    return result

def main():
    parser = argparse.ArgumentParser(description="MPI Conjugate Gradient solver (lab #2)")
    parser.add_argument("--N", type=int, required=True, help="Problem size N")
    parser.add_argument("--eps", type=float, default=1e-6, help="Stopping tolerance (relative residual)")
    parser.add_argument("--max-it", type=int, default=10**9, help="Maximum iterations (default huge)")
    parser.add_argument("--variant", choices=["dup","dist"], default="dup",
                        help="dup: duplicate x,b on each rank; dist: distribute x,b across ranks")
    parser.add_argument("--matrix", choices=["model","dense"], default="model",
                        help="Matrix representation")
    parser.add_argument("--problem", choices=["ones","sine"], default="ones",
                        help="Right-hand side / true solution model")
    parser.add_argument("--verbose", action="store_true", help="Extra prints per iteration")
    args = parser.parse_args()

    comm = MPI.COMM_WORLD
    rank = comm.Get_rank()

    # Soft hint about divisibility
    if rank == 0:
        for p in [1,2,4,8,16]:
            if args.N % p != 0:
                print(f"[warn] N={args.N} is not divisible by {p}. It's allowed to choose an N divisible by all.", flush=True)
                break

    if args.max_it == 10**9:
        args.max_it = max(10, min(args.N, 5 * args.N))

    cg_solve(args.N, args.eps, args.max_it, args.variant, args.matrix, args.problem, comm)

if __name__ == "__main__":
    main()
