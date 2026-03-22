"""Placeholder operations mixin for presentations."""

from __future__ import annotations

from typing import cast

from ... import ops
from ...utils import is_four_number_bounds
from ..helpers import PresentationMixinBase
from ..placeholders import (
    build_placeholder_chart_payload,
    build_placeholder_table_payload,
)


class PresentationPlaceholderMixin(PresentationMixinBase):
    """Mixin providing placeholder operations for Presentation."""

    def list_placeholders(self, slide_index: int) -> list[dict[str, object]]:
        """Bridge op: list all placeholders on a slide."""
        result = self.execute(
            ops.OP_LIST_PLACEHOLDERS,
            {"slide_index": slide_index},
        )
        return cast("list[dict[str, object]]", result.get("placeholders", []))

    def get_slide_layout_ref(self, slide_index: int) -> tuple[str, str]:
        """Return (layout_part, master_part) for a slide."""
        result = self.execute(
            ops.OP_GET_SLIDE_LAYOUT_REF,
            {"slide_index": slide_index},
        )
        layout_part = str(result.get("layout_part", ""))
        master_part = str(result.get("master_part", ""))
        return layout_part, master_part

    def set_placeholder_content(
        self,
        slide_index: int,
        ph_index: int,
        ph_type: str = "",
        **kwargs: object,
    ) -> None:
        """Bridge op: insert rich content into a placeholder."""
        bounds = kwargs.get("bounds")

        payload: dict[str, object] = {
            "slide_index": slide_index,
            "index": ph_index,
            "ph_type": ph_type,
        }
        text = kwargs.get("text")
        if isinstance(text, str):
            payload["text"] = text
        image_path = kwargs.get("image_path")
        if isinstance(image_path, str):
            payload["image_path"] = image_path
        if is_four_number_bounds(bounds):
            payload["bounds"] = list(bounds)
        text_style = kwargs.get("text_style")
        if isinstance(text_style, dict):
            payload["text_style"] = text_style
        force_rect = kwargs.get("force_rect_geometry")
        if isinstance(force_rect, bool):
            payload["force_rect_geometry"] = force_rect
        table_payload = build_placeholder_table_payload(
            kwargs.get("table"),
            kwargs.get("table_rows"),
            kwargs.get("table_cols"),
        )
        if table_payload is not None:
            payload["table"] = table_payload
        chart_payload = build_placeholder_chart_payload(
            kwargs,
            bounds,
        )
        if chart_payload is not None:
            payload["chart"] = chart_payload

        self.execute(ops.OP_SET_PLACEHOLDER_CONTENT, payload)
        self.invalidate_cache()
