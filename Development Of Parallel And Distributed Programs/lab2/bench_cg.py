#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Benchmark runner for the MPI CG solver.
Runs both variants (dup, dist) over a set of process counts,
parses the solver output, computes speedup/efficiency, and saves plots.

Usage example:
  python bench_cg.py --solver ./cg_mpi.py --N 2000000 --eps 1e-8 --procs 1 2 4 8 16 --repeats 3

Notes:
- This script launches "mpiexec -n P python SOLVER ...".
- Ensure you run it on a machine/cluster with MPI and mpi4py available.
"""
import argparse, subprocess, sys, json, time, csv, os, datetime
import numpy as np
import matplotlib.pyplot as plt

def run_once(mpi_launcher, solver, procs, N, eps, variant, problem, matrix):
    cmd = [mpi_launcher, "-n", str(procs), sys.executable, solver,
           "--N", str(N), "--eps", str(eps), "--variant", variant, "--problem", problem, "--matrix", matrix]
    p = subprocess.run(cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True, check=False)
    out = p.stdout
    # Find JSON_RESULT_BEGIN/END
    jbeg = out.find("JSON_RESULT_BEGIN")
    jend = out.find("JSON_RESULT_END")
    if jbeg == -1 or jend == -1:
        print(out)
        raise RuntimeError("Failed to parse solver output (no JSON_RESULT markers).")
    jtxt = out[jbeg + len("JSON_RESULT_BEGIN"):jend].strip()
    data = json.loads(jtxt)
    return data, out

def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--solver", default="./cg_mpi.py")
    ap.add_argument("--N", type=int, required=True)
    ap.add_argument("--eps", type=float, default=1e-6)
    ap.add_argument("--procs", type=int, nargs="+", default=[1,2,4,8,16])
    ap.add_argument("--repeats", type=int, default=1)
    ap.add_argument("--variant", choices=["both","dup","dist"], default="both")
    ap.add_argument("--problem", choices=["ones","sine"], default="ones")
    ap.add_argument("--matrix", choices=["model","dense"], default="model")
    ap.add_argument("--mpi", default="mpiexec", help="MPI launcher: mpiexec or mpirun")
    ap.add_argument("--outdir", default="bench_out")
    args = ap.parse_args()

    os.makedirs(args.outdir, exist_ok=True)
    stamp = datetime.datetime.now().strftime("%Y%m%d_%H%M%S")
    csv_path = os.path.join(args.outdir, f"results_{stamp}.csv")

    variants = ["dup","dist"] if args.variant == "both" else [args.variant]
    rows = []
    for var in variants:
        for p in args.procs:
            for r in range(args.repeats):
                data, out = run_once(args.mpi, args.solver, p, args.N, args.eps, var, args.problem, args.matrix)
                print(f"[bench] {var} procs={p} iters={data['iterations']} time={data['time_sec']:.4f}s relres={data['rel_residual']:.3e}")
                rows.append({
                    "variant": var,
                    "procs": p,
                    "N": args.N,
                    "eps": args.eps,
                    "iterations": data["iterations"],
                    "time_sec": data["time_sec"],
                    "rel_residual": data["rel_residual"],
                    "matrix": args.matrix,
                    "problem": args.problem,
                    "timestamp": stamp,
                })

    # Aggregate by (variant, procs): average time
    import collections
    agg = collections.defaultdict(list)
    for row in rows:
        agg[(row["variant"], row["procs"])].append(row["time_sec"])
    avg = {k: float(np.mean(v)) for k,v in agg.items()}

    # Compute speedup and efficiency using the p=1 baseline per variant
    speedup = {}
    efficiency = {}
    for var in variants:
        t1 = avg.get((var, 1), None)
        if t1 is None:
            continue
        for p in args.procs:
            tp = avg.get((var, p), None)
            if tp is None: continue
            sp = t1 / tp if tp > 0 else np.nan
            ef = sp / p
            speedup[(var, p)] = sp
            efficiency[(var, p)] = ef

    # Save CSV
    with open(csv_path, "w", newline="") as f:
        w = csv.writer(f)
        w.writerow(["variant","procs","N","eps","iterations","time_sec","rel_residual","matrix","problem","timestamp"])
        for row in rows:
            w.writerow([row[c] for c in ["variant","procs","N","eps","iterations","time_sec","rel_residual","matrix","problem","timestamp"]])
    print(f"[bench] Saved raw results to {csv_path}")

    # Make plots
    def plot_metric(metric_dict, ylabel, fname):
        plt.figure()
        for var in variants:
            xs = []
            ys = []
            for p in sorted(args.procs):
                val = metric_dict.get((var, p), None)
                if val is not None:
                    xs.append(p)
                    ys.append(val)
            plt.plot(xs, ys, marker="o", label=var)
        plt.xlabel("Cores / MPI processes")
        plt.ylabel(ylabel)
        plt.title(ylabel + " vs cores")
        plt.legend()
        out = os.path.join(args.outdir, f"{fname}_{stamp}.png")
        plt.savefig(out, dpi=150, bbox_inches="tight")
        print(f"[bench] Saved plot {out}")

    # Time vs cores (lower is better)
    # Construct the time metric
    time_metric = {}
    for (var,p), _ in avg.items():
        time_metric[(var,p)] = avg[(var,p)]
    plot_metric(time_metric, "Time, s (lower is better)", "time_vs_cores")
    plot_metric(speedup, "Speedup S(p) = T(1)/T(p)", "speedup_vs_cores")
    plot_metric(efficiency, "Efficiency E(p) = S(p)/p", "efficiency_vs_cores")

if __name__ == "__main__":
    main()
