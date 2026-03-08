"""Batch shape-insert helpers for the Presentation API."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .. import ops
from .helpers import PresentationMixinBase

if TYPE_CHECKING:
    from collections.abc import Mapping


class PresentationShapeBatchMixin(PresentationMixinBase):
    """Mixin providing high-throughput shape insertion helpers."""

    def add_textboxes(
        self,
        slide_index: int,
        textboxes: list[Mapping[str, float | str]],
    ) -> list[int]:
        """Add multiple textboxes to one slide in a single bridge call."""
        payload: dict[str, object] = {"slide_index": slide_index}
        payload["textboxes"] = cast(
            "object",
            [dict(textbox) for textbox in textboxes],
        )
        result = self.execute(ops.OP_ADD_TEXTBOXES, payload)
        return cast("list[int]", result.get("shape_ids", []))
