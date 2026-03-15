"""Presentation notes mixin."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ...schemas import ShapeUpdate


class PresentationNotesMixin(PresentationMixinBase):
    """Mixin providing speaker notes methods."""

    def get_notes_payload(self, slide_index: int) -> dict[str, object]:
        """Return raw notes payload for a slide index."""
        return self.execute(ops.OP_GET_NOTES, {"slide_index": slide_index})

    def get_notes(self, slide_index: int) -> str:
        """Return speaker notes plain text for a slide index."""
        result = self.get_notes_payload(slide_index)
        return str(cast("str", result.get("text", "")))

    def set_notes(self, slide_index: int, text: str) -> None:
        """Set speaker notes plain text for a slide index."""
        self.execute(ops.OP_SET_NOTES, {"slide_index": slide_index, "text": text})

    def set_notes_shape_text(self, slide_index: int, shape_id: int, text: str) -> None:
        """Set text for one notes shape by shape ID."""
        self.execute(
            ops.OP_SET_NOTES_SHAPE_TEXT,
            {"slide_index": slide_index, "shape_id": shape_id, "text": text},
        )

    def set_notes_shape_props(
        self, slide_index: int, shape_id: int, updates: ShapeUpdate
    ) -> None:
        """Patch style/geometry properties for one notes shape by shape ID."""
        self.execute(
            ops.OP_SET_NOTES_SHAPE_PROPS,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "updates": cast("dict[str, object]", updates),
            },
        )
