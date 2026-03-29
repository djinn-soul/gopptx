#!/usr/bin/env python3
"""Architectural guardrails for LOC ceiling ratcheting and import-cycle budget."""

from __future__ import annotations

import argparse
import ast
import json
import sys
from collections import defaultdict
from dataclasses import dataclass
from pathlib import Path
from typing import TextIO, cast


@dataclass(frozen=True)
class GuardrailConfig:
    """Configuration for architectural guardrails.

    Attributes:
        line_ceiling: Maximum allowed lines for a file before being flagged.
        roots: List of root directories to monitor for LOC checking.
        extensions: Set of file extensions to include in monitoring.
        line_count_baseline: Baseline line counts for over-ceiling files.
        python_import_cycle_budget: Maximum allowed Python import cycles.
        python_cycle_root: Root directory for Python cycle detection.
    """

    line_ceiling: int
    roots: list[str]
    extensions: set[str]
    line_count_baseline: dict[str, int]
    python_import_cycle_budget: int
    python_cycle_root: str


def _load_config(config_path: Path) -> GuardrailConfig:
    raw = json.loads(config_path.read_text(encoding="utf-8"))
    return GuardrailConfig(
        line_ceiling=int(raw["line_ceiling"]) + raw["line_ceiling"] * 0.05,
        roots=list(cast("list[str]", raw["roots"])),
        extensions=set(cast("list[str]", raw["extensions"])),
        line_count_baseline=dict(cast("dict[str, int]", raw["line_count_baseline"])),
        python_import_cycle_budget=int(raw["python_import_cycle_budget"]),
        python_cycle_root=str(raw["python_cycle_root"]),
    )


def _is_test_file(path: Path) -> bool:
    name = path.name
    if path.suffix == ".go" and name.endswith("_test.go"):
        return True
    if path.suffix == ".py" and (name.startswith("test_") or name.endswith("_test.py")):
        return True
    return "tests" in path.parts


def _line_count(path: Path) -> int:
    with path.open("rb") as file:
        return sum(1 for _ in file)


def _iter_monitored_files(repo_root: Path, cfg: GuardrailConfig) -> list[Path]:
    files: list[Path] = []
    for root in cfg.roots:
        root_path = repo_root / root
        if not root_path.exists():
            continue
        for path in root_path.rglob("*"):
            if not path.is_file():
                continue
            if path.suffix not in cfg.extensions:
                continue
            if _is_test_file(path):
                continue
            files.append(path)
    return files


def _check_loc_guardrail(repo_root: Path, cfg: GuardrailConfig) -> list[str]:
    violations: list[str] = []
    for path in _iter_monitored_files(repo_root, cfg):
        rel = path.relative_to(repo_root).as_posix()
        lines = _line_count(path)
        if lines <= cfg.line_ceiling:
            continue
        baseline = cfg.line_count_baseline.get(rel)
        if baseline is None:
            violations.append(
                f"NEW over-ceiling file: {rel} has {lines} lines (ceiling {cfg.line_ceiling})"
            )
            continue
        if lines > baseline:
            violations.append(
                f"Grew over baseline: {rel} has {lines} lines (baseline {baseline}, ceiling {cfg.line_ceiling})"
            )
    return violations


def _module_name(pkg_root: Path, file_path: Path) -> str:
    rel = file_path.relative_to(pkg_root).with_suffix("")
    return "gopptx." + ".".join(rel.parts)


def _resolve_module_candidate(
    available: set[str], candidate: str, import_aliases: list[str] | None = None
) -> set[str]:
    results: set[str] = set()
    parts = candidate.split(".")
    for idx in range(len(parts), 1, -1):
        maybe = ".".join(parts[:idx])
        if maybe in available:
            results.add(maybe)
            break
    if import_aliases:
        for alias in import_aliases:
            maybe = f"{candidate}.{alias}"
            if maybe in available:
                results.add(maybe)
    return results


def _process_import_node(node: ast.Import, available: set[str]) -> set[str]:
    """Process an ast.Import node and return resolved module edges."""
    results: set[str] = set()
    for alias in node.names:
        imported = alias.name
        if imported == "gopptx" or imported.startswith("gopptx."):
            results.update(_resolve_module_candidate(available, imported))
    return results


def _process_import_from_node(
    node: ast.ImportFrom,
    current_pkg_parts: list[str],
    available: set[str],
) -> set[str] | None:
    """Process an ast.ImportFrom node and return resolved module edges, or None if not relevant."""
    level = node.level or 0
    module = node.module or ""

    if level > 0:
        # level 1 is current package, level 2 is parent, etc.
        base_parts = current_pkg_parts[: len(current_pkg_parts) - level + 1]
        target = ".".join(base_parts + ([module] if module else []))
    else:
        target = module

    if not target:
        return None
    if target != "gopptx" and not target.startswith("gopptx."):
        return None

    import_aliases = [alias.name for alias in node.names if alias.name != "*"]
    return _resolve_module_candidate(available, target, import_aliases)


def _build_python_import_graph(pkg_root: Path) -> dict[str, set[str]]:
    modules: dict[str, Path] = {}
    for path in pkg_root.rglob("*.py"):
        if "tests" in path.parts:
            continue
        modules[_module_name(pkg_root, path)] = path

    available = set(modules.keys())
    edges: dict[str, set[str]] = defaultdict(set)

    for mod, path in modules.items():
        tree = ast.parse(path.read_text(encoding="utf-8"))
        current_pkg_parts = mod.split(".")[:-1]

        for node in ast.walk(tree):
            if isinstance(node, ast.Import):
                edges[mod].update(_process_import_node(node, available))
                continue

            if not isinstance(node, ast.ImportFrom):
                continue

            resolved = _process_import_from_node(node, current_pkg_parts, available)
            if resolved:
                edges[mod].update(resolved)

    return edges


def _tarjan_scc(graph: dict[str, set[str]]) -> list[list[str]]:
    index = 0
    stack: list[str] = []
    in_stack: set[str] = set()
    indices: dict[str, int] = {}
    low: dict[str, int] = {}
    sccs: list[list[str]] = []

    def dfs(node: str) -> None:
        nonlocal index
        indices[node] = index
        low[node] = index
        index += 1
        stack.append(node)
        in_stack.add(node)

        for nxt in graph.get(node, set()):
            if nxt not in indices:
                dfs(nxt)
                low[node] = min(low[node], low[nxt])
            elif nxt in in_stack:
                low[node] = min(low[node], indices[nxt])

        if low[node] != indices[node]:
            return
        component: list[str] = []
        while True:
            popped = stack.pop()
            in_stack.remove(popped)
            component.append(popped)
            if popped == node:
                break
        sccs.append(component)

    for node in graph:
        if node not in indices:
            dfs(node)
    return sccs


def _check_python_cycle_budget(repo_root: Path, cfg: GuardrailConfig) -> list[str]:
    pkg_root = repo_root / cfg.python_cycle_root
    if not pkg_root.exists():
        return [f"python cycle root does not exist: {cfg.python_cycle_root}"]

    graph = _build_python_import_graph(pkg_root)
    sccs = _tarjan_scc(graph)
    cycles = [sorted(component) for component in sccs if len(component) > 1]
    if len(cycles) <= cfg.python_import_cycle_budget:
        return []
    lines = [
        f"Python import-cycle budget exceeded: {len(cycles)} > {cfg.python_import_cycle_budget}"
    ]
    lines.extend(
        f"  - cycle({len(component)}): {', '.join(component)}"
        for component in sorted(cycles, key=lambda c: (-len(c), c[0]))
    )
    return lines


def _write_current_baseline(
    repo_root: Path, cfg: GuardrailConfig, config_path: Path
) -> None:
    baseline: dict[str, int] = {}
    for path in _iter_monitored_files(repo_root, cfg):
        rel = path.relative_to(repo_root).as_posix()
        lines = _line_count(path)
        if lines > cfg.line_ceiling:
            baseline[rel] = lines

    raw = json.loads(config_path.read_text(encoding="utf-8"))
    raw["line_count_baseline"] = dict(sorted(baseline.items()))
    config_path.write_text(json.dumps(raw, indent=2) + "\n", encoding="utf-8")


def _find_repo_root(start: Path) -> Path:
    for candidate in (start, *start.parents):
        if (candidate / ".git").exists():
            return candidate
    return start


def _write_line(message: str, *, stream: TextIO = sys.stdout) -> None:
    stream.write(f"{message}\n")


def main() -> int:
    """Run architectural guardrail checks.

    Returns:
        Exit code: 0 if all checks pass, 1 if there are violations.
    """
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--config",
        default="scripts/ci/architectural_guardrails.json",
        help="Path to guardrail config JSON.",
    )
    parser.add_argument(
        "--write-current-baseline",
        action="store_true",
        help="Rewrite config baseline with current over-ceiling file line counts.",
    )
    args = parser.parse_args()

    repo_root = _find_repo_root(Path.cwd())
    config_path = (repo_root / args.config).resolve()
    cfg = _load_config(config_path)

    if args.write_current_baseline:
        _write_current_baseline(repo_root, cfg, config_path)
        _write_line(
            f"Updated baseline in {config_path.relative_to(repo_root).as_posix()}"
        )
        return 0

    violations: list[str] = []
    violations.extend(_check_loc_guardrail(repo_root, cfg))
    violations.extend(_check_python_cycle_budget(repo_root, cfg))

    if not violations:
        _write_line("Architectural guardrails passed.")
        return 0

    _write_line("Architectural guardrails failed:", stream=sys.stderr)
    sys.stderr.writelines(f"- {violation}\n" for violation in violations)
    return 1


if __name__ == "__main__":
    sys.exit(main())
