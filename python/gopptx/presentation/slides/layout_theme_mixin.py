"""Presentation layout and theme mixins."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ...schemas import SlideLayoutInfo, SlideMasterCloneResult


class PresentationLayoutMixin(PresentationMixinBase):
    """Mixin providing slide layout management methods."""

    def list_slide_layouts(self) -> list[SlideLayoutInfo]:
        """List all available slide layouts."""
        result = self.execute(ops.OP_LIST_SLIDE_LAYOUTS, {})
        layouts = cast("list[dict[str, object]]", result.get("layouts", []))
        for item in layouts:
            if "Name" in item and "name" not in item:
                item["name"] = item["Name"]
            if "Part" in item and "part" not in item:
                item["part"] = item["Part"]
            if "MasterPart" in item and "master_part" not in item:
                item["master_part"] = item["MasterPart"]
        return cast("list[SlideLayoutInfo]", layouts)

    def rebind_slide_layout(self, slide_index: int, layout_part: str) -> None:
        """Rebind a slide to a different layout."""
        target_layout = layout_part
        if "/" not in target_layout:
            for layout in self.list_slide_layouts():
                if layout.get("name") == target_layout:
                    target_layout = cast("str", layout.get("part", target_layout))
                    break
        self.execute(
            ops.OP_REBIND_SLIDE_LAYOUT,
            {"slide_index": slide_index, "layout_part": target_layout},
        )
        self.invalidate_cache()

    def clone_layout_master_family(self, layout_part: str) -> SlideMasterCloneResult:
        """Clone a layout and its master family."""
        result = self.execute(
            ops.OP_CLONE_LAYOUT_MASTER_FAMILY, {"layout_part": layout_part}
        )
        self.invalidate_cache()
        return cast("SlideMasterCloneResult", result)


class PresentationThemeMixin(PresentationMixinBase):
    """Mixin providing theme and slide size configuration methods."""

    def apply_theme(self, theme_name: str) -> None:
        """Apply a theme to the presentation."""
        theme = theme_name
        if theme_name.lower() == "office":
            theme = "Corporate"
        self.execute(ops.OP_APPLY_THEME, {"theme_name": theme})
        self.invalidate_cache()

    def set_slide_size(self, width: int, height: int) -> None:
        """Set the slide size."""
        self.execute(ops.OP_SET_SLIDE_SIZE, {"width": width, "height": height})
        self.invalidate_cache()

    def set_global_theme_preset(self, name: str) -> None:
        """Apply a named built-in theme preset (e.g. 'facet', 'ion', 'office')."""
        self.execute(ops.OP_SET_GLOBAL_THEME_PRESET, {"name": name})
        self.invalidate_cache()

    def set_theme_font_scheme(self, major: str, minor: str) -> None:
        """Update major/minor latin typefaces across all theme parts."""
        self.execute(ops.OP_SET_THEME_FONT_SCHEME, {"major": major, "minor": minor})

    def set_theme_color_scheme(self, **colors: str) -> None:
        """Update one or more standard theme color slots.

        Args:
            **colors: Keyword arguments mapping color slot names to hex values.
                Valid keys: dk1, lt1, dk2, lt2, accent1..accent6, hlink, fol_hlink.
        """
        valid_keys = {
            "dk1",
            "lt1",
            "dk2",
            "lt2",
            "accent1",
            "accent2",
            "accent3",
            "accent4",
            "accent5",
            "accent6",
            "hlink",
            "fol_hlink",
        }
        payload: dict[str, object] = {
            k: v for k, v in colors.items() if k in valid_keys
        }
        self.execute(ops.OP_SET_THEME_COLOR_SCHEME, payload)

    def get_theme_inventory(self) -> dict[str, object]:
        """Return all theme parts and master/theme bindings in the package."""
        return self.execute(ops.OP_GET_THEME_INVENTORY, {})

    def get_layout_shapes(self, layout_part: str) -> list[str]:
        """Return the shape names defined in a slide layout."""
        result = self.execute(ops.OP_GET_LAYOUT_SHAPES, {"layout_part": layout_part})
        return cast("list[str]", result.get("shapes", []))

    def get_master_shapes(self, master_part: str) -> list[str]:
        """Return the shape names defined in a slide master."""
        result = self.execute(ops.OP_GET_MASTER_SHAPES, {"master_part": master_part})
        return cast("list[str]", result.get("shapes", []))

    def get_layout_placeholders(self, layout_part: str) -> list[dict[str, object]]:
        """Return placeholder metadata for a slide layout."""
        result = self.execute(
            ops.OP_GET_LAYOUT_PLACEHOLDERS, {"layout_part": layout_part}
        )
        return cast("list[dict[str, object]]", result.get("placeholders", []))

    def get_master_placeholders(self, master_part: str) -> list[dict[str, object]]:
        """Return placeholder metadata for a slide master."""
        result = self.execute(
            ops.OP_GET_MASTER_PLACEHOLDERS, {"master_part": master_part}
        )
        return cast("list[dict[str, object]]", result.get("placeholders", []))
