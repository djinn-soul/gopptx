"""Type-only bridge contract for slide proxy objects."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import Protocol

if TYPE_CHECKING:
    from collections.abc import Mapping, Sequence

    from ...schemas import TextRun


class TextOperationsProtocol(Protocol):
    """Text and run management."""

    def get_slide_text_states(self, slide_index: int) -> list[dict[str, object]]:
        """Protocol member."""
        ...

    def get_shape_text_state(
        self, slide_index: int, shape_id: int
    ) -> dict[str, object]:
        """Protocol member."""
        ...

    def get_shape_runs(self, slide_index: int, shape_id: int) -> list[TextRun]:
        """Protocol member."""
        ...

    def set_shape_runs(
        self, slide_index: int, shape_id: int, runs: list[TextRun]
    ) -> None:
        """Protocol member."""
        ...

    def update_shape_run_text(
        self, slide_index: int, shape_id: int, run_index: int, text: str
    ) -> None:
        """Protocol member."""
        ...

    def update_slide_run_texts(
        self, slide_index: int, updates: list[dict[str, object]]
    ) -> None:
        """Protocol member."""
        ...

    def append_shape_run(self, slide_index: int, shape_id: int, run: TextRun) -> None:
        """Protocol member."""
        ...

    def add_textbox(
        self,
        slide_index: int,
        left: float,
        top: float,
        width: float,
        height: float,
        *,
        text: str = "",
        **kwargs: object,
    ) -> int:
        """Protocol member."""
        ...

    def add_textboxes(
        self, slide_index: int, textboxes: Sequence[Mapping[str, object]]
    ) -> list[int]:
        """Protocol member."""
        ...

    def flush_pending_textbox_adds(self, slide_index: int) -> list[int]:
        """Protocol member."""
        ...

    def has_pending_textbox_adds(self, slide_index: int) -> bool:
        """Protocol member."""
        ...

    def queue_textbox_add(self, slide_index: int, payload: dict[str, object]) -> int:
        """Protocol member."""
        ...

    def flush_pending_slide_run_text_updates(self, slide_index: int) -> None:
        """Protocol member."""
        ...

    def has_pending_slide_run_text_updates(self, slide_index: int) -> bool:
        """Protocol member."""
        ...

    def queue_shape_runs_replace(
        self, slide_index: int, shape_id: int, runs: list[dict[str, object]]
    ) -> None:
        """Protocol member."""
        ...

    def queue_shape_run_text_update(
        self, slide_index: int, shape_id: int, run_index: int, text: str
    ) -> None:
        """Protocol member."""
        ...

    def flush_pending_shape_runs_replacements(
        self, slide_index: int, shape_id: int
    ) -> None:
        """Protocol member."""
        ...

    def has_pending_shape_runs_replace(self, slide_index: int, shape_id: int) -> bool:
        """Protocol member."""
        ...
