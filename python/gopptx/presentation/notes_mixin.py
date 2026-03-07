"""Presentation notes mixin."""
# ruff: noqa: D102

from __future__ import annotations

from typing import cast

from .. import ops
from .helpers import PresentationProtocol


class PresentationNotesMixin(PresentationProtocol):
    """Mixin providing speaker notes methods."""

    def get_notes_payload(self, slide_index: int) -> dict[str, object]:
        return self.execute(ops.OP_GET_NOTES, {"slide_index": slide_index})

    def get_notes(self, slide_index: int) -> str:
        result = self.get_notes_payload(slide_index)
        return str(cast("str", result.get("text", "")))

    def set_notes(self, slide_index: int, text: str) -> None:
        self.execute(ops.OP_SET_NOTES, {"slide_index": slide_index, "text": text})
