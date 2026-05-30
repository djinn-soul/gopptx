"""Normalize release tags to Python package versions."""

from __future__ import annotations

import argparse
import re
import sys

VERSION_PATTERN = re.compile(
    r"^(?P<base>\d+\.\d+\.\d+)"
    r"(?:(?P<label>alpha|a|beta|b|rc|dev|test)[.-]?(?P<number>\d+)|"
    r"[-.](?P<dash_label>alpha|a|beta|b|rc|dev|test)[.-]?(?P<dash_number>\d+))?$"
)

LABEL_MAP = {
    "alpha": "a",
    "a": "a",
    "beta": "b",
    "b": "b",
    "rc": "rc",
    "dev": ".dev",
    "test": ".dev",
}


def normalize_version(source_version: str, *, require_alpha: bool = False) -> str:
    """Return a PEP 440-compatible version for a git tag or version string."""
    raw_version = source_version.removeprefix("v")
    match = VERSION_PATTERN.fullmatch(raw_version)
    if match is None:
        raise ValueError(
            "version must look like v0.2.0-alpha.1, v0.2.0-alpha1, "
            "v0.2.0-beta.1, v0.2.0-rc.1, or v0.2.0-test1"
        )

    label = match.group("label") or match.group("dash_label")
    number = match.group("number") or match.group("dash_number")
    if require_alpha and label not in {"alpha", "a"}:
        raise ValueError(
            "production PyPI publishing is currently limited to alpha tags"
        )

    base = match.group("base")
    if label is None:
        if require_alpha:
            raise ValueError(
                "production PyPI publishing is currently limited to alpha tags"
            )
        return base

    return f"{base}{LABEL_MAP[label]}{number}"


def parse_args() -> argparse.Namespace:
    """Parse command-line arguments."""
    parser = argparse.ArgumentParser(
        description="Normalize a release tag or version to a Python package version."
    )
    parser.add_argument("source_version", help="Git tag or raw version to normalize.")
    parser.add_argument(
        "--require-alpha",
        action="store_true",
        help="Reject non-alpha versions.",
    )
    parser.add_argument(
        "--github-output",
        action="store_true",
        help="Emit GitHub Actions output lines instead of only the version.",
    )
    return parser.parse_args()


def main() -> int:
    """Run the version normalizer."""
    args = parse_args()
    try:
        version = normalize_version(
            args.source_version,
            require_alpha=args.require_alpha,
        )
    except ValueError as exc:
        sys.stderr.write(f"error: {exc}\n")
        return 1

    if args.github_output:
        sys.stdout.write(f"version={version}\n")
        sys.stdout.write(f"source_version={args.source_version}\n")
    else:
        sys.stdout.write(f"{version}\n")
    return 0


if __name__ == "__main__":
    sys.exit(main())
