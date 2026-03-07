"""Parity tests for generic shape-level shadow controls."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from pathlib import Path

import pytest
from gopptx import Presentation


def test_shape_shadow_round_trip_fixture(tmp_path: Path) -> None:
    """Test shape shadow round trip fixture."""
    out_path = tmp_path / "dml_shape_shadow.pptx"

    with Presentation.new("DML Shadow") as prs:
        slide = prs.slides[0]
        shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            text="shadow",
            properties={
                "shadow": {
                    "color": "123456",
                    "blur_emu": 60000,
                    "distance_emu": 40000,
                    "angle_deg": 45.0,
                }
            },
        )
        slide.update_shape(shape_id, {"shadow": {"inherit": True}})
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert "effectLst" not in slide_xml


def test_shape_shadow_rejects_inherit_with_explicit_fields() -> None:
    """Test shape shadow rejects inherit with explicit fields."""
    with Presentation.new("DML Shadow Invalid") as prs:
        slide = prs.slides[0]
        with pytest.raises(Exception, match=r"shadow.inherit"):
            slide.add_shape(
                "rect",
                (1000000, 1000000, 5000000, 1500000),
                properties={"shadow": {"inherit": True, "color": "123456"}},
            )


def test_shape_shadow_inherit_false_emits_empty_effect_list(tmp_path: Path) -> None:
    """Test shape shadow inherit false emits empty effect list."""
    out_path = tmp_path / "dml_shape_shadow_inherit_false.pptx"
    with Presentation.new("DML Shadow inherit false") as prs:
        slide = prs.slides[0]
        slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            properties={"shadow": {"inherit": False}},
        )
        prs.save(out_path)
    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")
    assert "<a:effectLst/>" in slide_xml
