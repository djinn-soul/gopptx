"""Parity tests for generic shape-level soft-edge controls."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from pathlib import Path

import pytest
from gopptx import Presentation


def test_shape_soft_edge_round_trip_fixture(tmp_path: Path) -> None:
    """Test shape soft edge round trip fixture."""
    out_path = tmp_path / "dml_shape_soft_edge.pptx"

    with Presentation.new("DML SoftEdge") as prs:
        slide = prs.slides[0]
        slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            properties={"soft_edge": {"radius_emu": 62000}},
        )
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert '<a:softEdge rad="62000"/>' in slide_xml


def test_shape_soft_edge_rejects_shadow_inherit_true_combo() -> None:
    """Test shape soft edge rejects shadow inherit true combo."""
    with Presentation.new("DML SoftEdge Invalid") as prs:
        slide = prs.slides[0]
        with pytest.raises(Exception, match=r"shadow.inherit"):
            slide.add_shape(
                "rect",
                (1000000, 1000000, 5000000, 1500000),
                properties={
                    "shadow": {"inherit": True},
                    "soft_edge": {"radius_emu": 62000},
                },
            )
