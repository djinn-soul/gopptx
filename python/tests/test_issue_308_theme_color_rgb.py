"""Issue #308: resolve the RGB value behind a theme colour.

ThemeColor is a str-Enum, so a bare slot name reads as a valid argument at the
call site. Passing one used to raise AttributeError deep in the colour model.
"""

from __future__ import annotations

import re
import zipfile
from typing import TYPE_CHECKING

import pytest
from gopptx import Presentation
from gopptx.color.model import ThemeColor

if TYPE_CHECKING:
    from pathlib import Path

_HEX_RE = re.compile(r"^[0-9A-F]{6}$")


def test_slot_name_string_resolves() -> None:
    """A bare slot name resolves instead of raising AttributeError."""
    with Presentation.new(title="Issue 308") as prs:
        resolved = prs.theme_color_rgb("accent1")

        assert resolved is not None
        assert _HEX_RE.match(resolved), resolved


def test_string_and_enum_agree() -> None:
    """The str form and the ThemeColor member resolve identically."""
    with Presentation.new(title="Issue 308 parity") as prs:
        assert prs.theme_color_rgb("accent1") == prs.theme_color_rgb(
            ThemeColor.ACCENT_1
        )


def test_alias_slots_resolve_to_their_scheme_entry() -> None:
    """tx1/bg1 are UI aliases for dk1/lt1 and resolve to the same values."""
    with Presentation.new(title="Issue 308 alias") as prs:
        assert prs.theme_color_rgb("tx1") == prs.theme_color_rgb("dk1")
        assert prs.theme_color_rgb("bg1") == prs.theme_color_rgb("lt1")


def test_brightness_shifts_the_resolved_value() -> None:
    """A brightness tweak changes the resolved colour."""
    with Presentation.new(title="Issue 308 brightness") as prs:
        base = prs.theme_color_rgb("accent1")
        lighter = prs.theme_color_rgb("accent1", brightness=0.4)

        assert base is not None
        assert lighter is not None
        assert lighter != base
        assert _HEX_RE.match(lighter), lighter


def test_unknown_slot_raises_value_error() -> None:
    """An unknown slot names the valid ones rather than failing obscurely."""
    with (
        Presentation.new(title="Issue 308 bad") as prs,
        pytest.raises(ValueError, match="unknown theme colour"),
    ):
        prs.theme_color_rgb("chartreuse")


def test_resolved_value_matches_the_theme_part(tmp_path: Path) -> None:
    """The resolved accent1 is the value stored in the theme part."""
    output_path = tmp_path / "theme.pptx"

    with Presentation.new(title="Issue 308 theme") as prs:
        resolved = prs.theme_color_rgb("accent1")
        prs.save(output_path)

    with zipfile.ZipFile(output_path) as zf:
        theme_xml = zf.read("ppt/theme/theme1.xml").decode("utf-8")
    match = re.search(r'<a:accent1>\s*<a:srgbClr val="([0-9A-Fa-f]{6})"', theme_xml)

    assert match is not None, theme_xml[:400]
    assert resolved == match.group(1).upper()
