"""Colour model for gopptx: RGB values, theme colours, and RGB resolution."""

from .model import (
    ColorFormat,
    ColorType,
    ThemeColor,
    apply_brightness,
    normalize_hex_color,
)

__all__ = [
    "ColorFormat",
    "ColorType",
    "ThemeColor",
    "apply_brightness",
    "normalize_hex_color",
]
