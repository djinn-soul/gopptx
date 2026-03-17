"""Slide text-state mixin for shape text operations."""

from __future__ import annotations

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from ...presentation.presentation import Presentation
    from ...schemas import TextRun


class SlideTextMixin:
    """Mixin providing text-state operations for slide shapes."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int:
            """Slide index."""
            ...

        def _apply_cached_run_text_update(
            self,
            shape_id: int,
            run_index: int,
            text: str,
        ) -> None:
            """Apply text update to local run caches."""
            ...

        def _replace_cached_runs(
            self,
            shape_id: int,
            runs: list[dict[str, object]],
        ) -> None:
            """Replace cached run list for a shape."""
            ...

        def _append_cached_run(
            self,
            shape_id: int,
            run: dict[str, object],
        ) -> None:
            """Append a run entry to local caches."""
            ...

    def _invalidate_text_state_cache_if_present(self) -> None:
        """Invalidate text-state cache if the concrete slide exposes it."""
        invalidate = getattr(self, "_invalidate_text_state_cache", None)
        if callable(invalidate):
            invalidate()

    def _flush_pending_text_updates_if_present(self) -> None:
        """Flush deferred run-text updates queued at presentation layer."""
        flush = getattr(
            self._presentation, "flush_pending_slide_run_text_updates", None
        )
        has_pending = getattr(
            self._presentation, "has_pending_slide_run_text_updates", None
        )
        if not callable(flush) or not callable(has_pending):
            return
        if has_pending(self.index):
            flush(self.index)

    def _flush_pending_shape_runs_replace_if_present(self, shape_id: int) -> None:
        """Flush pending whole-run replacement queued for a specific shape."""
        flush = getattr(
            self._presentation, "flush_pending_shape_runs_replacements", None
        )
        has_pending = getattr(
            self._presentation, "has_pending_shape_runs_replace", None
        )
        if not callable(flush) or not callable(has_pending):
            return
        if has_pending(self.index, shape_id):
            flush(slide_index=self.index, shape_id=shape_id)

    def get_shape_text_state(self, shape_id: int) -> dict[str, object]:
        """Return text state payload for a shape."""
        self._flush_pending_text_updates_if_present()
        self._flush_pending_shape_runs_replace_if_present(shape_id)
        return self._presentation.get_shape_text_state(self.index, shape_id)

    def get_shape_runs(self, shape_id: int) -> list[TextRun]:
        """Return rich run payloads for a shape."""
        self._flush_pending_text_updates_if_present()
        self._flush_pending_shape_runs_replace_if_present(shape_id)
        return self._presentation.get_shape_runs(self.index, shape_id)

    def set_shape_runs(self, shape_id: int, runs: list[TextRun]) -> None:
        """Replace all runs for a shape, using queueing when available."""
        self._flush_pending_text_updates_if_present()
        queue = getattr(self._presentation, "queue_shape_runs_replace", None)
        normalized = [dict(run) for run in runs]
        if not callable(queue):
            self._presentation.set_shape_runs(self.index, shape_id, runs)
        else:
            queue(self.index, shape_id, normalized)
        self._replace_cached_runs(shape_id, normalized)

    def update_shape_run_text(self, shape_id: int, run_index: int, text: str) -> None:
        """Update one run's text by index, preserving cache coherence."""
        self._flush_pending_shape_runs_replace_if_present(shape_id)
        queue = getattr(self._presentation, "queue_shape_run_text_update", None)
        if not callable(queue):
            self._presentation.update_shape_run_text(
                self.index, shape_id, run_index, text
            )
            self._invalidate_shape_cache_if_present()
            self._invalidate_text_state_cache_if_present()
        else:
            queue(self.index, shape_id, run_index, text)
            self._apply_cached_run_text_update(shape_id, run_index, text)

    def update_shape_run_texts(
        self,
        updates: list[tuple[int, int, str]] | list[dict[str, object]],
    ) -> None:
        """Apply multiple run-text updates on the current slide."""
        self._flush_pending_text_updates_if_present()
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
        for update in normalized:
            shape_id = update.get("shape_id")
            run_index = update.get("run_index")
            text = update.get("text")
            if (
                isinstance(shape_id, int)
                and isinstance(run_index, int)
                and isinstance(text, str)
            ):
                self._apply_cached_run_text_update(shape_id, run_index, text)

    def append_shape_run(self, shape_id: int, run: TextRun) -> None:
        """Append a run to a shape and update local caches."""
        self._flush_pending_shape_runs_replace_if_present(shape_id)
        self._flush_pending_text_updates_if_present()
        self._presentation.append_shape_run(self.index, shape_id, run)
        self._append_cached_run(shape_id, dict(run))

    def _invalidate_shape_cache_if_present(self) -> None:
        """Invalidate shape cache if the concrete slide exposes it."""
        invalidate = getattr(self, "_invalidate_shape_cache", None)
        if callable(invalidate):
            invalidate()
