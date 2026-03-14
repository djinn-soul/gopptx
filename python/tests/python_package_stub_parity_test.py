"""Parity checks for top-level package runtime exports vs stub exports."""

from __future__ import annotations

import ast
from pathlib import Path


def _read_all_tuple(path: Path) -> tuple[str, ...]:
    tree = ast.parse(path.read_text(encoding="utf-8"))
    for node in tree.body:
        if not isinstance(node, ast.Assign):
            continue
        has_all_target = any(
            isinstance(target, ast.Name) and target.id == "__all__"
            for target in node.targets
        )
        if not has_all_target:
            continue
        value = ast.literal_eval(node.value)
        if isinstance(value, tuple) and all(isinstance(item, str) for item in value):
            return value
        raise AssertionError(f"__all__ in {path} must be a tuple[str, ...]")
    raise AssertionError(f"missing __all__ in {path}")


def test_top_level_stub_all_matches_runtime() -> None:
    root = Path(__file__).resolve().parents[2]
    runtime_path = root / "python" / "gopptx" / "__init__.py"
    stub_path = root / "python" / "gopptx" / "__init__.pyi"
    assert stub_path.exists(), "missing root package stub: python/gopptx/__init__.pyi"

    runtime_all = _read_all_tuple(runtime_path)
    stub_all = _read_all_tuple(stub_path)
    assert stub_all == runtime_all, (
        "root package stub __all__ drift detected; keep __init__.pyi aligned with "
        "__init__.py exports for IntelliSense parity"
    )
