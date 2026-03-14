"""Presentation mixin for slide-scoped text batch updates."""
# ruff: noqa: D102

from __future__ import annotations

from typing import TYPE_CHECKING

from ... import ops


class PresentationTextBatchMixin:
    """Bridge-backed helpers for bulk slide text mutations."""

    if TYPE_CHECKING:

        def execute(
            self,
            op: str,
            payload: dict[str, object] | None = None,
        ) -> dict[str, object]: ...

    def update_slide_run_texts(
        self,
        slide_index: int,
        updates: list[dict[str, object]],
    ) -> None:
        self.execute(
            ops.OP_UPDATE_SLIDE_RUN_TEXTS,
            {
                "slide_index": slide_index,
                "updates": updates,
            },
        )

    def update_deck_run_texts(
        self,
        slide_updates: list[dict[str, object]],
    ) -> None:
        self.execute(
            ops.OP_UPDATE_DECK_RUN_TEXTS,
            {
                "slides": slide_updates,
            },
        )
