"""Parity tests for paragraph indent/hanging facade support."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from pathlib import Path

import pytest
from gopptx import ParagraphProps, Presentation


def test_paragraph_props_aliases_emit_ooxml_attrs(tmp_path: Path) -> None:
    """Test paragraph props aliases emit ooxml attrs."""
    out_path = tmp_path / "paragraph_props_facade.pptx"

    with Presentation.new("Paragraph Facade") as prs:
        slide = prs.slides[0]
        shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            text="paragraph",
            paragraph=ParagraphProps(left_margin=228600, hanging_indent=228600),
        )
        assert shape_id > 0
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert 'marL="228600"' in slide_xml
    assert 'indent="-228600"' in slide_xml


def test_paragraph_props_rejects_negative_hanging() -> None:
    """Test paragraph props rejects negative hanging."""
    with pytest.raises(ValueError, match=r"paragraph\.hanging must be >= 0"):
        ParagraphProps(hanging=-1).to_payload()


def test_paragraph_props_tab_stops_payload() -> None:
    """Tab-stop positions should serialize into paragraph payload."""
    payload = ParagraphProps(tab_stops=[228600, 457200]).to_payload()
    assert payload["tab_stops"] == [228600, 457200]


def test_paragraph_props_rejects_negative_tab_stop() -> None:
    """Negative tab-stop positions should fail fast."""
    with pytest.raises(ValueError, match=r"paragraph\.tab_stops values must be >= 0"):
        ParagraphProps(tab_stops=[-1]).to_payload()


def test_paragraph_props_bullet_payload() -> None:
    """Bullet properties should serialize in normalized payload."""
    payload = ParagraphProps(
        bullet_style="roman_upper",
        bullet_char="*",
        bullet_color="00AAFF",
        bullet_size_pct=90,
    ).to_payload()
    assert payload["bullet_style"] == "roman_upper"
    assert payload["bullet_char"] == "*"
    assert payload["bullet_color"] == "00AAFF"
    assert payload["bullet_size_pct"] == 90


def test_paragraph_props_custom_bullet_requires_char() -> None:
    """Custom bullet style must provide bullet char."""
    with pytest.raises(
        ValueError,
        match=r"paragraph\.bullet_char is required when bullet_style='custom'",
    ):
        ParagraphProps(bullet_style="custom").to_payload()
