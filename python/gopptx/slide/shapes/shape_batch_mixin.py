"""Batch shape helpers for slide proxies."""

from __future__ import annotations

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from collections.abc import Mapping

    from ...presentation.presentation import Presentation


class SlideShapeBatchMixin:
    """Mixin exposing high-throughput batch shape inserts on Slide objects."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int:
            """Current slide index."""
            ...

    def add_textboxes(
        self,
        textboxes: list[Mapping[str, float | str]],
    ) -> list[int]:
        """Add multiple textboxes to this slide in a single bridge call."""
        shape_ids = self._presentation.add_textboxes(self.index, textboxes)
        invalidate_shapes = getattr(self, "_invalidate_shape_cache", None)
        if callable(invalidate_shapes):
            invalidate_shapes()
        invalidate_text = getattr(self, "_invalidate_text_state_cache", None)
        if callable(invalidate_text):
            invalidate_text()
        return shape_ids

    def add_connectors(
        self,
        connectors: list[Mapping[str, object]],
    ) -> list[int]:
        """Add multiple connectors to this slide in a single bridge call."""
        shape_ids = self._presentation.add_connectors(self.index, connectors)
        invalidate_shapes = getattr(self, "_invalidate_shape_cache", None)
        if callable(invalidate_shapes):
            invalidate_shapes()
        invalidate_text = getattr(self, "_invalidate_text_state_cache", None)
        if callable(invalidate_text):
            invalidate_text()
        return shape_ids
