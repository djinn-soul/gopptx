"""Presentation notes mixin."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ...schemas import ShapeUpdate


class PresentationNotesMixin(PresentationMixinBase):
    """Mixin providing speaker notes methods."""

    def notes_slide_exists(self, slide_index: int) -> bool:
        """Return True if a notes slide exists for the given slide index."""
        result = self.execute(ops.OP_NOTES_SLIDE_EXISTS, {"slide_index": slide_index})
        return bool(result.get("notes_slide_exists", False))

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

    def list_notes_shapes(self, slide_index: int) -> list[dict[str, object]]:
        """List all shapes in the notes pane of a slide."""
        result = self.execute(ops.OP_LIST_NOTES_SHAPES, {"slide_index": slide_index})
        return cast("list[dict[str, object]]", result.get("shapes", []))

    def list_notes_placeholders(self, slide_index: int) -> list[dict[str, object]]:
        """List all placeholders in the notes pane of a slide."""
        result = self.execute(
            ops.OP_LIST_NOTES_PLACEHOLDERS, {"slide_index": slide_index}
        )
        return cast("list[dict[str, object]]", result.get("placeholders", []))

    def get_handout_master(self) -> dict[str, object]:
        """Return handout master information.

        Returns a dict with keys ``present`` (bool), ``orientation``
        (``"landscape"`` or ``"portrait"``), and ``slides_per_page`` (int).
        """
        return self.execute(ops.OP_GET_HANDOUT_MASTER, {})

    def update_handout_master(
        self,
        *,
        orientation: str = "",
        slides_per_page: int = 0,
    ) -> None:
        """Configure the handout master.

        Args:
            orientation: ``"landscape"`` or ``"portrait"``.
            slides_per_page: Number of slide thumbnails per page
                (1, 2, 3, 4, 6, or 9).
        """
        payload: dict[str, object] = {}
        if orientation:
            payload["orientation"] = orientation
        if slides_per_page:
            payload["slides_per_page"] = slides_per_page
        self.execute(ops.OP_UPDATE_HANDOUT_MASTER, payload)

    def has_digital_signature(self) -> bool:
        """Return True if the presentation has a digital signature."""
        result = self.execute(ops.OP_HAS_DIGITAL_SIGNATURE, {})
        return bool(result.get("has_digital_signature", False))

    def update_notes_master(
        self,
        *,
        header: str = "",
        footer: str = "",
        show_date_time: bool = True,
        show_slide_num: bool = True,
    ) -> None:
        """Configure the global notes master.

        Args:
            header: Header text to display on notes pages.
            footer: Footer text to display on notes pages.
            show_date_time: Whether to show the date/time placeholder.
            show_slide_num: Whether to show the slide number placeholder.
        """
        self.execute(
            ops.OP_UPDATE_NOTES_MASTER,
            {
                "header": header,
                "footer": footer,
                "show_date_time": show_date_time,
                "show_slide_num": show_slide_num,
            },
        )
