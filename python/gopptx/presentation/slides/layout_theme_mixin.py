"""Presentation layout and theme mixins."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from typing_extensions import Protocol

    from ...presentation.theme.theme import Theme
    from ...schemas import SlideLayoutInfo, SlideMasterCloneResult

    class _PresentationThemeOps(Protocol):
        def execute(
            self, op: str, payload: dict[str, object] | None = None
        ) -> dict[str, object]: ...

        def invalidate_cache(self) -> None: ...

        def set_theme_color_scheme(self, **colors: str) -> None: ...

        def set_theme_font_scheme(self, major: str, minor: str) -> None: ...


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

    def reorder_slide_layouts(
        self, layout_parts: list[str], master_part: str | None = None
    ) -> None:
        """Reorder slide layouts within a slide master (Issue #1080).

        Args:
            layout_parts: List of layout part paths or layout names in desired order.
            master_part: Optional master part path (e.g. 'ppt/slideMasters/slideMaster1.xml').
        """
        layouts = self.list_slide_layouts()
        name_to_part = {
            cast("str", layout.get("name", "")): cast("str", layout.get("part", ""))
            for layout in layouts
        }
        resolved = [name_to_part.get(p, p) for p in layout_parts]
        payload: dict[str, object] = {"layout_parts": resolved}
        if master_part:
            payload["master_part"] = master_part
        self.execute(ops.OP_REORDER_SLIDE_LAYOUTS, payload)
        self.invalidate_cache()

    def clone_layout_master_family(self, layout_part: str) -> SlideMasterCloneResult:
        """Clone a layout and its master family."""
        result = self.execute(
            ops.OP_CLONE_LAYOUT_MASTER_FAMILY, {"layout_part": layout_part}
        )
        self.invalidate_cache()
        return cast("SlideMasterCloneResult", result)

    def add_slide_master(self) -> str:
        """Add a new slide master and return its part path."""
        result = self.execute(ops.OP_ADD_SLIDE_MASTER, {})
        self.invalidate_cache()
        master_part = result.get("master_part")
        if not isinstance(master_part, str):
            raise TypeError("bridge response master_part must be a string")
        return master_part

    def remove_slide_master(self, master_part: str) -> None:
        """Remove a slide master by part path."""
        self.execute(ops.OP_REMOVE_SLIDE_MASTER, {"master_part": master_part})
        self.invalidate_cache()

    def add_slide_layout(
        self, master_part: str, layout_name: str = "Custom Layout"
    ) -> str:
        """Add a slide layout under a slide master and return the layout part path."""
        result = self.execute(
            ops.OP_ADD_SLIDE_LAYOUT,
            {"master_part": master_part, "layout_name": layout_name},
        )
        self.invalidate_cache()
        layout_part = result.get("layout_part")
        if not isinstance(layout_part, str):
            raise TypeError("bridge response layout_part must be a string")
        return layout_part

    def remove_slide_layout(self, layout_part: str) -> None:
        """Remove a slide layout by part path."""
        self.execute(ops.OP_REMOVE_SLIDE_LAYOUT, {"layout_part": layout_part})
        self.invalidate_cache()


class PresentationThemeMixin(PresentationMixinBase):
    """Mixin providing theme and slide size configuration methods."""

    def apply_theme(self: _PresentationThemeOps, theme: str | Theme) -> None:
        """Apply a theme to the presentation."""
        if isinstance(theme, str):
            theme_name = "Corporate" if theme.lower() == "office" else theme
            self.execute(ops.OP_APPLY_THEME, {"theme_name": theme_name})
            self.invalidate_cache()
            return

        color_dict = theme.colors.to_dict()
        self.set_theme_color_scheme(**color_dict)
        font_dict = theme.fonts.to_dict()
        self.set_theme_font_scheme(font_dict["major_font"], font_dict["minor_font"])

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

    def add_layout_shape(
        self,
        layout_part: str,
        shape_type: str,
        bounds: tuple[float, float, float, float],
    ) -> int:
        """Add an autoshape to a slide layout and return its shape id.

        The shape appears on every slide that uses the layout.

        Args:
            layout_part: Layout part path, e.g. "ppt/slideLayouts/slideLayout1.xml".
            shape_type: Preset geometry name, e.g. "rect".
            bounds: (left, top, width, height) in EMU.
        """
        left, top, width, height = bounds
        result = self.execute(
            ops.OP_ADD_LAYOUT_SHAPE,
            {
                "layout_part": layout_part,
                "shape_type": shape_type,
                "left": left,
                "top": top,
                "width": width,
                "height": height,
            },
        )
        return int(cast("int", result.get("shape_id", 0)))

    def add_layout_textbox(
        self,
        layout_part: str,
        text: str,
        bounds: tuple[float, float, float, float],
    ) -> int:
        """Add a text box to a slide layout and return its shape id.

        This is the operation upstream python-pptx has no API for (issue #1044).
        """
        left, top, width, height = bounds
        result = self.execute(
            ops.OP_ADD_LAYOUT_TEXTBOX,
            {
                "layout_part": layout_part,
                "text": text,
                "left": left,
                "top": top,
                "width": width,
                "height": height,
            },
        )
        return int(cast("int", result.get("shape_id", 0)))

    def add_master_shape(
        self,
        master_part: str,
        shape_type: str,
        bounds: tuple[float, float, float, float],
    ) -> int:
        """Add an autoshape to a slide master and return its shape id."""
        left, top, width, height = bounds
        result = self.execute(
            ops.OP_ADD_MASTER_SHAPE,
            {
                "master_part": master_part,
                "shape_type": shape_type,
                "left": left,
                "top": top,
                "width": width,
                "height": height,
            },
        )
        return int(cast("int", result.get("shape_id", 0)))

    def add_master_textbox(
        self,
        master_part: str,
        text: str,
        bounds: tuple[float, float, float, float],
    ) -> int:
        """Add a text box to a slide master and return its shape id."""
        left, top, width, height = bounds
        result = self.execute(
            ops.OP_ADD_MASTER_TEXTBOX,
            {
                "master_part": master_part,
                "text": text,
                "left": left,
                "top": top,
                "width": width,
                "height": height,
            },
        )
        return int(cast("int", result.get("shape_id", 0)))

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
