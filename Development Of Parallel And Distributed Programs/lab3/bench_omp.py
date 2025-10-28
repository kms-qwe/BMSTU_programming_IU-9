#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import argparse, subprocess, os, sys, json, csv, datetime
import numpy as np
import matplotlib.pyplot as plt

def run_once(solver, threads, N, eps, matrix, problem):
    env = os.environ.copy()
    env['OMP_NUM_THREADS'] = str(threads)
    env['NUMBA_NUM_THREADS'] = str(threads)
    cmd = [sys.executable, solver, '--N', str(N), '--eps', str(eps), '--matrix', matrix, '--problem', problem]
    p = subprocess.run(cmd, env=env, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
    out = p.stdout
    jbeg = out.find('JSON_RESULT_BEGIN')
    jend = out.find('JSON_RESULT_END')
    if jbeg == -1 or jend == -1:
        print(out)
        raise RuntimeError('Failed to parse solver output (no JSON_RESULT markers).')
    jtxt = out[jbeg+len('JSON_RESULT_BEGIN'):jend].strip()
    data = json.loads(jtxt)
    return data, out

def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--solver', default='./cg_omp.py')
    ap.add_argument('--N', type=int, required=True)
    ap.add_argument('--eps', type=float, default=1e-6)
    ap.add_argument('--threads', type=int, nargs='+', default=[1,2,4,8,16])
    ap.add_argument('--repeats', type=int, default=1)
    ap.add_argument('--matrix', choices=['model','dense'], default='model')
    ap.add_argument('--problem', choices=['ones','sine','mix'], default='mix')
    ap.add_argument('--outdir', default='bench_omp_out')
    args = ap.parse_args()

    os.makedirs(args.outdir, exist_ok=True)
    stamp = datetime.datetime.now().strftime('%Y%m%d_%H%M%S')
    csv_path = os.path.join(args.outdir, f'results_omp_{stamp}.csv')

    rows = []
    for t in args.threads:
        for r in range(args.repeats):
            data, out = run_once(args.solver, t, args.N, args.eps, args.matrix, args.problem)
            print(f"[bench-omp] threads={t} iters={data['iterations']} time={data['time_sec']:.4f}s relres={data['rel_residual']:.3e}")
            rows.append({
                'threads': t, 'N': args.N, 'eps': args.eps,
                'iterations': data['iterations'], 'time_sec': data['time_sec'],
                'rel_residual': data['rel_residual'],
                'matrix': args.matrix, 'problem': args.problem, 'timestamp': stamp
            })

    # aggregate
    from collections import defaultdict
    agg = defaultdict(list)
    for row in rows:
        agg[row['threads']].append(row['time_sec'])
    avg = {k: float(np.mean(v)) for k,v in agg.items()}
    t1 = avg.get(1, None)
    speedup = {k: (t1/avg[k] if avg[k]>0 else np.nan) for k in avg.keys()} if t1 is not None else {}
    efficiency = {k: (speedup[k]/k) for k in speedup.keys()} if t1 is not None else {}

    with open(csv_path, 'w', newline='') as f:
        w = csv.writer(f); w.writerow(['threads','N','eps','iterations','time_sec','rel_residual','matrix','problem','timestamp'])
        for row in rows:
            w.writerow([row[c] for c in ['threads','N','eps','iterations','time_sec','rel_residual','matrix','problem','timestamp']])
    print(f"[bench-omp] Saved raw results to {csv_path}")

    def plot_metric(x_to_val, ylabel, fname):
        xs = sorted(x_to_val.keys())
        ys = [x_to_val[x] for x in xs]
        plt.figure()
        plt.plot(xs, ys, marker='o')
        plt.xlabel('Threads'); plt.ylabel(ylabel); plt.title(ylabel+' vs threads')
        out = os.path.join(args.outdir, f'{fname}_{stamp}.png')
        plt.savefig(out, dpi=150, bbox_inches='tight')
        print(f"[bench-omp] Saved plot {out}")

    plot_metric(avg, 'Time, s (lower is better)', 'time_vs_threads')
    if t1 is not None:
        plot_metric(speedup, 'Speedup S(p)=T(1)/T(p)', 'speedup_vs_threads')
        plot_metric(efficiency, 'Efficiency E(p)=S(p)/p', 'efficiency_vs_threads')

if __name__ == '__main__':
    main()
