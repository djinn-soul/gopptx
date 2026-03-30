"""Run pip-audit with an ignore list for pre-commit/CI parity."""

from __future__ import annotations

import argparse
import runpy
import sys
from pathlib import Path


def _read_ignore_ids(path: Path) -> list[str]:
    if not path.exists():
        return []
    try:
        lines = path.read_text(encoding="utf-8").splitlines()
    except OSError as e:
        print(f"Error reading ignore file '{path}': {e}", file=sys.stderr)
        sys.exit(1)
    return [
        line
        for raw in lines
        if (line := raw.strip()) and not line.startswith("#")
    ]


def _parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Run pip-audit with optional ignore list."
    )
    parser.add_argument(
        "--requirements",
        default="requirements-dev.txt",
        help="Path to requirements file to audit.",
    )
    parser.add_argument(
        "--ignore-file",
        default="pip-audit-ignore.txt",
        help="Path to newline-delimited vulnerability IDs to ignore.",
    )
    args, _unknown = parser.parse_known_args()
    return args


def main() -> int:
    """Execute pip-audit and return its process-style exit code."""
    args = _parse_args()
    ignore_ids = _read_ignore_ids(Path(args.ignore_file))

    cmd = ["pip-audit", "-r", args.requirements, "--skip-editable"]
    for vuln_id in ignore_ids:
        cmd.extend(["--ignore-vuln", vuln_id])

    old_argv = sys.argv[:]
    try:
        sys.argv = cmd
        try:
            runpy.run_module("pip_audit", run_name="__main__")
            return 0
        except SystemExit as exc:
            code = exc.code
            return int(code) if isinstance(code, int) else 1
    finally:
        sys.argv = old_argv


if __name__ == "__main__":
    sys.exit(main())
