#!/usr/bin/env python3
"""Verify that Go and Python development versions stay synchronized."""

from __future__ import annotations

import re
import sys
from pathlib import Path

VERSION_PATTERNS = {
    "pyproject.toml": re.compile(r'(?m)^version = "([^"]+)"$'),
    "pkg/pptx/version.go": re.compile(r'(?m)^const Version = "([^"]+)"$'),
    "python/gopptx/_version.py": re.compile(r'(?m)^\s+__version__ = "([^"]+)"$'),
}


def read_version(repo_root: Path, relative_path: str, pattern: re.Pattern[str]) -> str:
    """Extract one version string from a repository file."""
    content = (repo_root / relative_path).read_text(encoding="utf-8")
    match = pattern.search(content)
    if match is None:
        raise ValueError(f"version not found in {relative_path}")
    return match.group(1)


def main() -> int:
    """Return nonzero when the development version sources disagree."""
    repo_root = Path(__file__).resolve().parents[2]
    versions = {
        path: read_version(repo_root, path, pattern)
        for path, pattern in VERSION_PATTERNS.items()
    }
    unique_versions = set(versions.values())
    if len(unique_versions) == 1:
        version = unique_versions.pop()
        sys.stdout.write(f"Version consistency passed: {version}\n")
        return 0

    sys.stderr.write("Version consistency failed:\n")
    for path, version in versions.items():
        sys.stderr.write(f"- {path}: {version}\n")
    return 1


if __name__ == "__main__":
    sys.exit(main())
