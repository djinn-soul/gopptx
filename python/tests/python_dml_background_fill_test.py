"""Parity tests for generic shape-level background/no-fill controls."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from pathlib import Path

import pytest
from gopptx import Presentation


def test_shape_background_fill_round_trip_fixture(tmp_path: Path) -> None:
    """Test shape background fill round trip fixture."""
    out_path = tmp_path / "dml_background_fill.pptx"

    with Presentation.new("DML Background Fill") as prs:
        slide = prs.slides[0]
        slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            properties={"fill": {"background": True}},
        )
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert "<a:noFill/>" in slide_xml  # noqa: S101


def test_shape_background_fill_rejects_false_flag() -> None:
    """Test shape background fill rejects false flag."""
    with Presentation.new("DML Background Fill Invalid") as prs:
        slide = prs.slides[0]
        with pytest.raises(Exception, match=r"fill.background"):
            slide.add_shape(
                "rect",
                (1000000, 1000000, 5000000, 1500000),
                properties={"fill": {"background": False}},
            )
