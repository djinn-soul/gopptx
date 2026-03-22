"""Type-only bridge contract for slide proxy objects."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import Protocol

if TYPE_CHECKING:
    from ...schemas import ShapeUpdate


class NotesOperationsProtocol(Protocol):
    """Speaker notes management."""

    def get_notes(self, slide_index: int) -> str:
        """Protocol member."""
        ...

    def get_notes_payload(self, slide_index: int) -> dict[str, object]:
        """Protocol member."""
        ...

    def set_notes(self, slide_index: int, text: str) -> None:
        """Protocol member."""
        ...

    def set_notes_shape_text(self, slide_index: int, shape_id: int, text: str) -> None:
        """Protocol member."""
        ...

    def set_notes_shape_props(
        self, slide_index: int, shape_id: int, updates: ShapeUpdate
    ) -> None:
        """Protocol member."""
        ...
