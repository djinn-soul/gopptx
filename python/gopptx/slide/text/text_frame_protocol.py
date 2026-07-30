"""Dependency-neutral protocol for live shape text-frame operations."""

from __future__ import annotations

from typing import Protocol


class ShapeTextFrameProtocol(Protocol):
    """Bridge operations required by paragraph and run proxies."""

    def get_paragraphs(self) -> list[dict[str, object]]:
        """Return all paragraph payloads."""
        ...

    def replace_paragraphs(self, paragraphs: list[dict[str, object]]) -> None:
        """Replace all shape-text paragraphs."""
        ...

    def get_paragraph_runs(self, paragraph_index: int) -> list[dict[str, object]]:
        """Return runs for one paragraph."""
        ...

    def replace_paragraph_runs(
        self, paragraph_index: int, runs: list[dict[str, object]]
    ) -> None:
        """Replace runs for one paragraph."""
        ...

    def get_paragraph_payload(self, paragraph_index: int) -> dict[str, object]:
        """Return one paragraph's formatting payload."""
        ...

    def set_paragraph_field(
        self, paragraph_index: int, field: str, value: object
    ) -> None:
        """Update one paragraph field."""
        ...
