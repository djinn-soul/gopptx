"""Text-frame and paragraph typed schema definitions."""

from __future__ import annotations

try:
    from typing import TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import TypedDict


class TextFrame(TypedDict, total=False):
    """Text frame settings."""

    margin_top: int
    margin_bottom: int
    margin_left: int
    margin_right: int
    word_wrap: bool
    auto_fit: bool
    auto_fit_type: str
    vertical_align: str
    orientation: str
    columns: int
    rotation: float


class Paragraph(TypedDict, total=False):
    """Paragraph settings."""

    indent: int
    hanging: int
    tab_stops: list[int]
    alignment: str
    level: int
    bullet_style: str
    bullet_char: str
    bullet_color: str
    bullet_size_pct: int
    line_spacing_pct: int
    line_spacing_pts: int
    space_before_pts: int
    space_after_pts: int
    rtl: bool
