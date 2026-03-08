"""Buffered run-text updates for presentation save/read flows."""
# ruff: noqa: D102

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .helpers import PresentationMixinBase
from .runtime import PresentationRuntimeMixin

if TYPE_CHECKING:
    from collections.abc import Iterable


class PresentationTextWriteBufferMixin(PresentationMixinBase):
    """Coalesce repeated run-text writes into slide-level bridge updates."""

    if TYPE_CHECKING:

        def update_slide_run_texts(
            self,
            slide_index: int,
            updates: list[dict[str, object]],
        ) -> None: ...

        def update_deck_run_texts(
            self,
            slide_updates: list[dict[str, object]],
        ) -> None: ...

    def __init__(self) -> None:
        """Initialize pending slide text-update state."""
        super().__init__()
        self._pending_slide_run_text_updates: dict[int, dict[tuple[int, int], str]] = {}

    def queue_shape_run_text_update(
        self,
        slide_index: int,
        shape_id: int,
        run_index: int,
        text: str,
    ) -> None:
        """Buffer a run-text update until a flush boundary is reached."""
        slide_updates = self._pending_slide_run_text_updates.setdefault(slide_index, {})
        slide_updates[shape_id, run_index] = text

    def has_pending_slide_run_text_updates(self, slide_index: int) -> bool:
        """Return whether a slide has buffered run-text updates."""
        return bool(self._pending_slide_run_text_updates.get(slide_index))

    def flush_pending_slide_run_text_updates(self, slide_index: int) -> None:
        """Flush buffered run-text updates for one slide."""
        slide_updates = self._pending_slide_run_text_updates.get(slide_index)
        if not slide_updates:
            return
        updates = cast(
            "list[dict[str, object]]",
            [
                {"shape_id": shape_id, "run_index": run_index, "text": text}
                for (shape_id, run_index), text in sorted(slide_updates.items())
            ],
        )
        self.update_slide_run_texts(slide_index, updates)
        self._pending_slide_run_text_updates.pop(slide_index, None)

    def flush_all_pending_slide_run_text_updates(self) -> None:
        """Flush buffered run-text updates for all slides."""
        slide_updates = cast(
            "list[dict[str, object]]",
            [
                {
                    "slide_index": slide_index,
                    "updates": [
                        {"shape_id": shape_id, "run_index": run_index, "text": text}
                        for (shape_id, run_index), text in sorted(
                            self._pending_slide_run_text_updates[slide_index].items()
                        )
                    ],
                }
                for slide_index in self._pending_slide_indexes()
            ],
        )
        if not slide_updates:
            return
        self.update_deck_run_texts(slide_updates)
        self._pending_slide_run_text_updates = {}

    def _pending_slide_indexes(self) -> Iterable[int]:
        return sorted(self._pending_slide_run_text_updates)

    def save(self, path: str) -> None:
        """Flush buffered text writes before saving the deck."""
        self.flush_all_pending_slide_run_text_updates()
        runtime_self = cast("PresentationRuntimeMixin", self)
        PresentationRuntimeMixin.save(runtime_self, path)
