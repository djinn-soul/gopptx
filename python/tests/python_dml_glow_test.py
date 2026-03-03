"""Parity tests for generic shape-level glow controls."""

from __future__ import annotations

import zipfile
from pathlib import Path

import pytest
from gopptx import Presentation


def test_shape_glow_round_trip_fixture(tmp_path: Path) -> None:
    out_path = tmp_path / "dml_shape_glow.pptx"

    with Presentation.new("DML Glow") as prs:
        slide = prs.slides[0]
        slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            properties={"glow": {"color": "ABCDEF", "radius_emu": 50000}},
        )
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert '<a:glow rad="50000"><a:srgbClr val="ABCDEF"/></a:glow>' in slide_xml  # noqa: S101


def test_shape_glow_rejects_shadow_inherit_true_combo() -> None:
    with Presentation.new("DML Glow Invalid") as prs:
        slide = prs.slides[0]
        with pytest.raises(Exception, match="shadow.inherit"):
            slide.add_shape(
                "rect",
                (1000000, 1000000, 5000000, 1500000),
                properties={"shadow": {"inherit": True}, "glow": {"radius_emu": 50000}},
            )
