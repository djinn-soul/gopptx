"""Parity tests for generic shape-level blur controls."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from pathlib import Path

import pytest
from gopptx import Presentation


def test_shape_blur_round_trip_fixture(tmp_path: Path) -> None:
    """Test shape blur round trip fixture."""
    out_path = tmp_path / "dml_shape_blur.pptx"

    with Presentation.new("DML Blur") as prs:
        slide = prs.slides[0]
        slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            properties={"blur": {"radius_emu": 61000}},
        )
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert '<a:blur rad="61000"/>' in slide_xml  # noqa: S101


def test_shape_blur_rejects_shadow_inherit_true_combo() -> None:
    """Test shape blur rejects shadow inherit true combo."""
    with Presentation.new("DML Blur Invalid") as prs:
        slide = prs.slides[0]
        with pytest.raises(Exception, match=r"shadow.inherit"):
            slide.add_shape(
                "rect",
                (1000000, 1000000, 5000000, 1500000),
                properties={"shadow": {"inherit": True}, "blur": {"radius_emu": 61000}},
            )
