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
