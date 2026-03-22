"""Type-only bridge contract for slide proxy objects."""

from __future__ import annotations

from typing_extensions import Protocol


class BaseEngineProtocol(Protocol):
    """Core engine and cache operations."""

    _batch_active: bool

    @property
    def batch_active(self) -> bool:
        """Protocol member."""
        ...

    def execute(
        self, op: str, payload: dict[str, object] | None = None
    ) -> dict[str, object]:
        """Protocol member."""
        ...

    def invalidate_cache(self) -> None:
        """Protocol member."""
        ...
