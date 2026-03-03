"""Parity tests for generic shape-level reflection controls."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from pathlib import Path

import pytest
from gopptx import Presentation


def test_shape_reflection_round_trip_fixture(tmp_path: Path) -> None:
    """Test shape reflection round trip fixture."""
    out_path = tmp_path / "dml_shape_reflection.pptx"

    with Presentation.new("DML Reflection") as prs:
        slide = prs.slides[0]
        slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            properties={"reflection": {"blur_emu": 14000, "distance_emu": 9000}},
        )
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert '<a:reflection blurRad="14000" dist="9000"/>' in slide_xml  # noqa: S101


def test_shape_reflection_rejects_shadow_inherit_true_combo() -> None:
    """Test shape reflection rejects shadow inherit true combo."""
    with Presentation.new("DML Reflection Invalid") as prs:
        slide = prs.slides[0]
        with pytest.raises(Exception, match=r"shadow.inherit"):
            slide.add_shape(
                "rect",
                (1000000, 1000000, 5000000, 1500000),
                properties={
                    "shadow": {"inherit": True},
                    "reflection": {"distance_emu": 9000},
                },
            )
