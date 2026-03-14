"""Buffered run-text updates for presentation save/read flows."""
# ruff: noqa: D102

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ..helpers import PresentationMixinBase
from ..runtime import PresentationRuntimeMixin

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
        self._pending_shape_runs_replacements: dict[
            int,
            dict[int, list[dict[str, object]]],
        ] = {}

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

    def queue_shape_runs_replace(
        self,
        slide_index: int,
        shape_id: int,
        runs: list[dict[str, object]],
    ) -> None:
        """Buffer full run replacement for one shape until a flush boundary."""
        slide_replacements = self._pending_shape_runs_replacements.setdefault(
            slide_index, {}
        )
        slide_replacements[shape_id] = [dict(run) for run in runs]

    def has_pending_shape_runs_replace(self, slide_index: int, shape_id: int) -> bool:
        """Return whether one shape has a pending full run replacement."""
        slide_replacements = self._pending_shape_runs_replacements.get(slide_index)
        if not slide_replacements:
            return False
        return shape_id in slide_replacements

    def has_pending_slide_run_text_updates(self, slide_index: int) -> bool:
        """Return whether a slide has buffered run-text updates."""
        return bool(self._pending_slide_run_text_updates.get(slide_index))

    def flush_pending_slide_run_text_updates(self, slide_index: int) -> None:
        """Flush buffered run-text updates for one slide."""
        slide_updates = self._pending_slide_run_text_updates.get(slide_index)
        if not slide_updates:
            return
        self.flush_pending_shape_runs_replacements(slide_index=slide_index)
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
        self.flush_all_pending_shape_runs_replacements()
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

    def flush_pending_shape_runs_replacements(
        self,
        *,
        slide_index: int,
        shape_id: int | None = None,
    ) -> None:
        """Flush buffered full run replacements for one slide or shape."""
        slide_replacements = self._pending_shape_runs_replacements.get(slide_index)
        if not slide_replacements:
            return
        runtime_self = cast("PresentationRuntimeMixin", self)
        target_shape_ids: list[int]
        if shape_id is None:
            target_shape_ids = sorted(slide_replacements)
        elif shape_id in slide_replacements:
            target_shape_ids = [shape_id]
        else:
            target_shape_ids = []
        if target_shape_ids:
            payload_updates = [
                {
                    "shape_id": pending_shape_id,
                    "runs": [
                        dict(run) for run in slide_replacements[pending_shape_id]
                    ],
                }
                for pending_shape_id in target_shape_ids
            ]
            PresentationRuntimeMixin.execute(
                runtime_self,
                ops.OP_SET_SLIDE_SHAPE_RUNS,
                {"slide_index": slide_index, "updates": payload_updates},
            )
            for pending_shape_id in target_shape_ids:
                slide_replacements.pop(pending_shape_id, None)
        if not slide_replacements:
            self._pending_shape_runs_replacements.pop(slide_index, None)

    def flush_all_pending_shape_runs_replacements(self) -> None:
        """Flush buffered full run replacements for all slides."""
        for slide_index in sorted(self._pending_shape_runs_replacements):
            self.flush_pending_shape_runs_replacements(slide_index=slide_index)

    def _pending_slide_indexes(self) -> Iterable[int]:
        return sorted(self._pending_slide_run_text_updates)

    def save(self, path: str) -> None:
        """Flush buffered text writes before saving the deck."""
        self.flush_all_pending_slide_run_text_updates()
        runtime_self = cast("PresentationRuntimeMixin", self)
        PresentationRuntimeMixin.save(runtime_self, path)
