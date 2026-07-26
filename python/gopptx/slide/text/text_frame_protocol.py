"""Dependency-neutral protocol for live shape text-frame operations."""

from __future__ import annotations

from typing import Protocol


class ShapeTextFrameProtocol(Protocol):
    """Bridge operations required by paragraph and run proxies."""

    def get_runs(self) -> list[dict[str, object]]:
        """Return the current shape text runs."""
        ...

    def replace_runs(self, runs: list[dict[str, object]]) -> None:
        """Replace all shape text runs."""
        ...

    def append_run(self, run: dict[str, object]) -> None:
        """Append one shape text run."""
        ...

    def update_run_text(self, run_index: int, text: str) -> None:
        """Update one run's text."""
        ...

    def get_paragraph_payload(self) -> dict[str, object]:
        """Return the normalized paragraph payload."""
        ...

    def set_paragraph_field(self, field: str, value: object) -> None:
        """Update one normalized paragraph field."""
        ...
