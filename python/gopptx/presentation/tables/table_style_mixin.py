"""Table style operations mixin."""

from __future__ import annotations

from typing import cast

from ... import ops
from ..helpers import PresentationMixinBase


class PresentationTableStyleMixin(PresentationMixinBase):
    """Mixin providing table style query and application methods."""

    def set_table_style(self, slide_index: int, shape_id: int, style: str) -> None:
        """Apply a table style by name or GUID."""
        from .table_styles import TableStyle

        style_guid = style
        if isinstance(style, str) and not style.startswith("{"):
            styles = TableStyle.get_all()
            if style not in styles:
                available = ", ".join(sorted(styles.keys()))
                raise ValueError(
                    f"Unknown style name '{style}'. Available: {available}"
                )
            style_guid = styles[style]

        self.execute(
            ops.OP_SET_TABLE_STYLE,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "style_guid": style_guid,
            },
        )

    def define_table_style(self, name: str, style_id: str | None = None) -> str:
        """Define a custom table style and return its resolved style ID."""
        payload: dict[str, object] = {"name": name}
        if style_id is not None:
            payload["style_id"] = style_id
        result = self.execute(ops.OP_DEFINE_TABLE_STYLE, payload)
        return str(result.get("style_id", ""))

    def list_table_styles(self) -> list[dict[str, str]]:
        """List available table styles visible to the presentation."""
        result = self.execute(ops.OP_LIST_TABLE_STYLES, {})
        return cast("list[dict[str, str]]", result.get("styles", []))

    def get_table_style_by_name(self, name: str) -> str | None:
        """Find a presentation table style GUID by name."""
        styles = self.list_table_styles()
        for style in styles:
            if style.get("name", "").lower() == name.lower():
                return style.get("guid")
        return None

    def get_all_table_style_names(self) -> list[str]:
        """Get all available table style names in the presentation."""
        styles = self.list_table_styles()
        return [style.get("name", "") for style in styles if "name" in style]
