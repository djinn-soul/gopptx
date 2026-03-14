"""Notes slide style helper mixin."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

if TYPE_CHECKING:
    from ..schemas import ShapeUpdate


class NotesSlideStyleMixin:
    """Convenience style helpers for notes shapes."""

    def _set_shape_props(self, shape_id: int, updates: ShapeUpdate) -> None: ...

    def set_shape_fill_gradient(
        self,
        shape_id: int,
        *,
        angle_deg: float,
        stops: list[dict[str, object]],
    ) -> None:
        """Set notes-shape gradient fill via shape update payload."""
        gradient: dict[str, object] = {
            "angle_deg": float(angle_deg),
            "stops": [dict(stop) for stop in stops],
        }
        self._set_shape_props(
            shape_id,
            cast("ShapeUpdate", {"fill": {"gradient": gradient}}),
        )

    def set_shape_fill_pattern(
        self,
        shape_id: int,
        *,
        preset: str,
        fg_color: str | None = None,
        bg_color: str | None = None,
    ) -> None:
        """Set notes-shape pattern fill via shape update payload."""
        pattern: dict[str, object] = {"preset": preset}
        if fg_color is not None:
            pattern["fg_color"] = fg_color
        if bg_color is not None:
            pattern["bg_color"] = bg_color
        self._set_shape_props(
            shape_id,
            cast("ShapeUpdate", {"fill": {"pattern": pattern}}),
        )

    def set_shape_line_dash(
        self,
        shape_id: int,
        *,
        dash_style: str,
        color: str | None = None,
        width_emu: int | None = None,
    ) -> None:
        """Set notes-shape dashed line style via shape update payload."""
        line: dict[str, object] = {"dash_style": dash_style}
        if color is not None:
            line["color"] = color
        if width_emu is not None:
            line["width_emu"] = int(width_emu)
        self._set_shape_props(shape_id, cast("ShapeUpdate", {"line": line}))
