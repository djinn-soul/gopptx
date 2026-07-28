"""Colour model, brightness, and theme-colour resolution tests."""

from __future__ import annotations

import pytest
from gopptx import ColorFormat, ColorType, Presentation, ThemeColor
from gopptx.color import apply_brightness, normalize_hex_color


def test_rgb_color_round_trips() -> None:
    color = ColorFormat.from_rgb("#ff0000")
    assert color.type is ColorType.RGB
    assert color.rgb == "FF0000"
    assert color.theme_color is None
    assert color.brightness == pytest.approx(0.0)
    assert color.to_rgb() == "FF0000"
    assert repr(color) == "ColorFormat(rgb='FF0000', brightness=0.0)"


def test_theme_color_needs_a_scheme() -> None:
    color = ColorFormat.from_theme(ThemeColor.ACCENT_1)
    assert color.type is ColorType.SCHEME
    assert color.theme_color is ThemeColor.ACCENT_1
    # Without a scheme there is nothing to resolve against.
    assert color.to_rgb() is None
    assert color.to_rgb({"accent1": "2563EB"}) == "2563EB"
    # An undefined slot stays unresolved rather than guessing.
    assert color.to_rgb({"accent2": "2563EB"}) is None


def test_theme_color_aliases_resolve_to_scheme_slots() -> None:
    scheme = {"dk1": "000000", "lt1": "FFFFFF", "fol_hlink": "7C3AED"}
    assert ColorFormat.from_theme(ThemeColor.TEXT_1).to_rgb(scheme) == "000000"
    assert ColorFormat.from_theme(ThemeColor.BACKGROUND_1).to_rgb(scheme) == "FFFFFF"
    assert (
        ColorFormat.from_theme(ThemeColor.FOLLOWED_HYPERLINK).to_rgb(scheme) == "7C3AED"
    )


def test_brightness_lightens_and_darkens() -> None:
    assert apply_brightness("808080", 0.0) == "808080"
    assert apply_brightness("000000", 1.0) == "FFFFFF"
    assert apply_brightness("FFFFFF", -1.0) == "000000"
    # Halfway to white, then halfway to black.
    assert apply_brightness("000000", 0.5) == "808080"
    assert apply_brightness("FFFFFF", -0.5) == "808080"

    lighter = ColorFormat.from_rgb("2563EB", brightness=0.5)
    assert lighter.to_rgb() == "92B1F5"
    assert lighter.with_brightness(0.0).to_rgb() == "2563EB"


def test_none_color_resolves_to_nothing() -> None:
    color = ColorFormat()
    assert color.type is ColorType.NONE
    assert color.to_rgb({"accent1": "2563EB"}) is None


def test_color_equality_and_hashing() -> None:
    assert ColorFormat.from_rgb("FF0000") == ColorFormat.from_rgb("#ff0000")
    assert ColorFormat.from_rgb("FF0000") != ColorFormat.from_rgb("FF0001")
    assert len({ColorFormat.from_rgb("FF0000"), ColorFormat.from_rgb("FF0000")}) == 1
    assert ColorFormat.from_rgb("FF0000") != "FF0000"


def test_color_validation_errors() -> None:
    with pytest.raises(ValueError, match="6-digit hex"):
        normalize_hex_color("FFF")
    with pytest.raises(ValueError, match="6-digit hex"):
        normalize_hex_color("GGGGGG")
    with pytest.raises(ValueError, match="brightness must be between"):
        ColorFormat.from_rgb("FF0000", brightness=1.5)
    with pytest.raises(ValueError, match="requires rgb"):
        ColorFormat(ColorType.RGB)
    with pytest.raises(ValueError, match="requires theme_color"):
        ColorFormat(ColorType.SCHEME)


def test_presentation_resolves_theme_colors() -> None:
    with Presentation.new("Colour resolution") as prs:
        scheme = prs.theme_color_scheme
        assert "accent1" in scheme
        assert len(scheme["accent1"]) == len("RRGGBB")

        accent1 = prs.theme_color_rgb(ThemeColor.ACCENT_1)
        assert accent1 == scheme["accent1"]
        assert prs.resolve_color(ColorFormat.from_theme(ThemeColor.ACCENT_1)) == accent1

        # tx1 is an alias of dk1 and resolves to the same value.
        assert prs.theme_color_rgb(ThemeColor.TEXT_1) == scheme["dk1"]

        lighter = prs.theme_color_rgb(ThemeColor.ACCENT_1, brightness=0.5)
        assert lighter == apply_brightness(scheme["accent1"], 0.5)


def test_presentation_theme_color_scheme_follows_edits() -> None:
    with Presentation.new("Colour edits") as prs:
        prs.set_theme_color_scheme(accent1="123456")
        assert prs.theme_color_scheme["accent1"] == "123456".upper()
        assert prs.theme_color_rgb(ThemeColor.ACCENT_1) == "123456"
