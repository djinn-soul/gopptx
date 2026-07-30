"""Line-format helpers shared by gridlines, series lines and label borders."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

if TYPE_CHECKING:
    from ...schemas import ChartLineFormatSpec

_HEX_LENGTHS = (6, 8)


def normalize_hex_color(value: str, name: str) -> str:
    """Return an uppercase RRGGBB string, rejecting anything else.

    The engine normalizes colours too; this keeps the error next to the call
    that made it rather than surfacing it from the bridge.
    """
    normalized = value.strip().lstrip("#").upper()
    if len(normalized) not in _HEX_LENGTHS or any(
        char not in "0123456789ABCDEF" for char in normalized
    ):
        raise ValueError(f"{name} must be a hex colour such as 'FF0000'")
    return normalized


def build_line_format(
    *,
    color: str | None = None,
    width_emu: int | None = None,
    dash: str | None = None,
    none: bool | None = None,
    name: str = "line",
) -> ChartLineFormatSpec:
    """Build a line-format payload from keyword arguments.

    ``none`` renders the element with no line at all and wins over the other
    three arguments, matching how PowerPoint hides a gridline or a border.
    """
    spec: dict[str, object] = {}
    if none:
        spec["none"] = True
        return cast("ChartLineFormatSpec", spec)
    if color is not None:
        spec["color"] = normalize_hex_color(color, f"{name} color")
    if width_emu is not None:
        if width_emu < 0:
            raise ValueError(f"{name} width_emu must not be negative")
        spec["width_emu"] = int(width_emu)
    if dash is not None:
        if not dash.strip():
            raise ValueError(f"{name} dash must not be empty")
        spec["dash"] = dash.strip()
    if not spec:
        raise ValueError(f"{name} needs one of color, width_emu, dash, none")
    return cast("ChartLineFormatSpec", spec)
