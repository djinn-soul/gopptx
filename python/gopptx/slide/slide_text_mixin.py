"""Slide text-state mixin for shape text operations."""
# ruff: noqa: D102

from __future__ import annotations

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from ..presentation.presentation import Presentation
    from ..schemas import TextRun


class SlideTextMixin:
    """Mixin providing text-state operations for slide shapes."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int: ...

    def _invalidate_text_state_cache_if_present(self) -> None:
        invalidate = getattr(self, "_invalidate_text_state_cache", None)
        if callable(invalidate):
            invalidate()

    def get_shape_text_state(self, shape_id: int) -> dict[str, object]:
        return self._presentation.get_shape_text_state(self.index, shape_id)

    def get_shape_runs(self, shape_id: int) -> list[TextRun]:
        return self._presentation.get_shape_runs(self.index, shape_id)

    def set_shape_runs(self, shape_id: int, runs: list[TextRun]) -> None:
        self._presentation.set_shape_runs(self.index, shape_id, runs)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()

    def update_shape_run_text(self, shape_id: int, run_index: int, text: str) -> None:
        self._presentation.update_shape_run_text(self.index, shape_id, run_index, text)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()

    def update_shape_run_texts(
        self,
        updates: list[tuple[int, int, str]] | list[dict[str, object]],
    ) -> None:
        normalized: list[dict[str, object]] = []
        for update in updates:
            if isinstance(update, tuple):
                shape_id, run_index, text = update
                normalized.append({
                    "shape_id": shape_id,
                    "run_index": run_index,
                    "text": text,
                })
                continue
            normalized.append(dict(update))
        self._presentation.update_slide_run_texts(self.index, normalized)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()

    def append_shape_run(self, shape_id: int, run: TextRun) -> None:
        self._presentation.append_shape_run(self.index, shape_id, run)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()

    def _invalidate_shape_cache_if_present(self) -> None:
        invalidate = getattr(self, "_invalidate_shape_cache", None)
        if callable(invalidate):
            invalidate()
