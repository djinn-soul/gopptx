from __future__ import annotations

import json
import time
from statistics import mean

from gopptx import Presentation, ops

ITERATIONS = 5
OPS_PER_ITERATION = 200


def benchmark_single_calls() -> list[float]:
    timings: list[float] = []
    with Presentation.new("Bridge Perf Single") as pres:
        for i in range(OPS_PER_ITERATION):
            pres.add_slide(f"Seed {i}")

        for it in range(ITERATIONS):
            start = time.perf_counter()
            for i in range(OPS_PER_ITERATION):
                pres.set_slide_title(0, f"Single {it}-{i}")
            timings.append(time.perf_counter() - start)
    return timings


def benchmark_batched_calls() -> list[float]:
    timings: list[float] = []
    with Presentation.new("Bridge Perf Batch") as pres:
        for i in range(OPS_PER_ITERATION):
            pres.add_slide(f"Seed {i}")

        for it in range(ITERATIONS):
            commands = [
                {
                    "op": ops.OP_SET_SLIDE_TITLE,
                    "payload": {"slide_index": 0, "title": f"Batch {it}-{i}"},
                }
                for i in range(OPS_PER_ITERATION)
            ]
            start = time.perf_counter()
            pres.execute_batch(commands)
            timings.append(time.perf_counter() - start)
    return timings


def benchmark_json_codec() -> tuple[float, float]:
    payload = {
        "api_version": 1,
        "request_id": "bench",
        "op": "set_slide_title",
        "payload": {"slide_index": 0, "title": "bench"},
    }
    encoded = json.dumps(payload).encode("utf-8")

    encode_start = time.perf_counter()
    for _ in range(OPS_PER_ITERATION * 100):
        json.dumps(payload).encode("utf-8")
    encode_elapsed = time.perf_counter() - encode_start

    decode_start = time.perf_counter()
    for _ in range(OPS_PER_ITERATION * 100):
        json.loads(encoded.decode("utf-8"))
    decode_elapsed = time.perf_counter() - decode_start

    return encode_elapsed, decode_elapsed


def main() -> None:
    single = benchmark_single_calls()
    batched = benchmark_batched_calls()
    json_encode, json_decode = benchmark_json_codec()

    single_avg = mean(single)
    batch_avg = mean(batched)
    speedup = single_avg / batch_avg if batch_avg > 0 else 0.0

    print("Bridge Performance Benchmark")
    print(f"ITERATIONS={ITERATIONS} OPS_PER_ITERATION={OPS_PER_ITERATION}")
    print(f"single_avg_seconds={single_avg:.6f}")
    print(f"batch_avg_seconds={batch_avg:.6f}")
    print(f"batch_speedup_x={speedup:.2f}")
    print(f"json_encode_total_seconds={json_encode:.6f}")
    print(f"json_decode_total_seconds={json_decode:.6f}")


if __name__ == "__main__":
    main()
