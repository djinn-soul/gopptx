from __future__ import annotations

import json
import os
import pathlib
import time

from gopptx import Presentation, ops

RUNS = 25
OPS_PER_RUN = 40


def bench_single(pres: Presentation, ops_per_run: int, runs: int) -> float:
    start = time.perf_counter()
    for run in range(runs):
        for i in range(ops_per_run):
            pres.execute(
                ops.OP_SET_SLIDE_TITLE,
                {"slide_index": 0, "title": f"single-{run}-{i}"},
            )
    return time.perf_counter() - start


def bench_batch(pres: Presentation, ops_per_run: int, runs: int) -> float:
    start = time.perf_counter()
    for run in range(runs):
        commands = [
            {
                "op": ops.OP_SET_SLIDE_TITLE,
                "payload": {"slide_index": 0, "title": f"batch-{run}-{i}"},
            }
            for i in range(ops_per_run)
        ]
        pres.execute_batch(commands)
    return time.perf_counter() - start


def main() -> None:
    with Presentation.new("Python Batch Benchmark") as pres:
        pres.add_slide("Slide 1")

        single_sec = bench_single(pres, OPS_PER_RUN, RUNS)
        batch_sec = bench_batch(pres, OPS_PER_RUN, RUNS)

        total_ops = OPS_PER_RUN * RUNS
        single_ms_per_op = (single_sec / total_ops) * 1000
        batch_ms_per_op = (batch_sec / total_ops) * 1000
        speedup = single_sec / batch_sec if batch_sec > 0 else 0.0

        report = {
            "runs": RUNS,
            "ops_per_run": OPS_PER_RUN,
            "total_ops": total_ops,
            "single_total_sec": round(single_sec, 6),
            "batch_total_sec": round(batch_sec, 6),
            "single_ms_per_op": round(single_ms_per_op, 6),
            "batch_ms_per_op": round(batch_ms_per_op, 6),
            "speedup_x": round(speedup, 3),
        }

        print(json.dumps(report, indent=2))

        output_dir = pathlib.Path("tmp")
        output_dir.mkdir(parents=True, exist_ok=True)
        out_path = output_dir / "python_batch_latency_benchmark.json"
        out_path.write_text(json.dumps(report, indent=2), encoding="utf-8")
        print(f"Wrote benchmark report: {out_path}")


if __name__ == "__main__":
    os.environ.setdefault("PYTHONUTF8", "1")
    main()
