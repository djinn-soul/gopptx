"""Parity tests for gradient and pattern shape fill controls."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from pathlib import Path

import pytest
from gopptx import Presentation


def test_shape_gradient_fill_round_trip_fixture(tmp_path: Path) -> None:
    """Test shape gradient fill round trip fixture."""
    out_path = tmp_path / "dml_gradient_fill.pptx"

    with Presentation.new("DML Gradient") as prs:
        slide = prs.slides[0]
        slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            properties={
                "fill": {
                    "gradient": {
                        "angle_deg": 90.0,
                        "stops": [
                            {"position_pct": 0.0, "color": "FF0000"},
                            {"position_pct": 100.0, "color": "0000FF"},
                        ],
                    }
                }
            },
        )
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert "<a:gradFill>" in slide_xml
    assert 'ang="5400000"' in slide_xml
    assert 'pos="0"' in slide_xml and 'pos="100000"' in slide_xml


def test_shape_pattern_fill_round_trip_fixture(tmp_path: Path) -> None:
    """Test shape pattern fill round trip fixture."""
    out_path = tmp_path / "dml_pattern_fill.pptx"

    with Presentation.new("DML Pattern") as prs:
        slide = prs.slides[0]
        slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            properties={
                "fill": {
                    "pattern": {
                        "preset": "diagCross",
                        "fg_color": "112233",
                        "bg_color": "AABBCC",
                    }
                }
            },
        )
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert 'a:pattFill prst="diagCross"' in slide_xml
    assert 'val="112233"' in slide_xml and 'val="AABBCC"' in slide_xml


def test_shape_fill_modes_reject_mutual_exclusive() -> None:
    """Test shape fill modes reject mutual exclusive."""
    with Presentation.new("DML Fill Invalid") as prs:
        slide = prs.slides[0]
        with pytest.raises(Exception, match="mutually exclusive"):
            slide.add_shape(
                "rect",
                (1000000, 1000000, 5000000, 1500000),
                properties={
                    "fill": {
                        "solid": "FF0000",
                        "gradient": {"stops": [{"color": "000000"}]},
                    }
                },
            )


def test_shape_fill_modes_reject_background_with_other_modes() -> None:
    """Test shape fill modes reject background with other modes."""
    with Presentation.new("DML Fill Invalid Background") as prs:
        slide = prs.slides[0]
        with pytest.raises(Exception, match="mutually exclusive"):
            slide.add_shape(
                "rect",
                (1000000, 1000000, 5000000, 1500000),
                properties={"fill": {"background": True, "solid": "FFFFFF"}},
            )
