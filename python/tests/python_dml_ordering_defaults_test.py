"""Parity fixtures for DML XML ordering and defaults."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from pathlib import Path

from gopptx import Presentation


def test_shape_dml_style_xml_ordering_fixture(tmp_path: Path) -> None:
    """Test shape dml style xml ordering fixture."""
    out_path = tmp_path / "dml_style_ordering.pptx"

    with Presentation.new("DML Ordering") as prs:
        slide = prs.slides[0]
        shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            properties={
                "fill": {"solid": "112233"},
                "line": {"color": "445566", "width_emu": 25400, "dash_style": "dash"},
                "shadow": {"color": "778899", "distance_emu": 20000},
                "glow": {"color": "AABBCC", "radius_emu": 40000},
            },
        )
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    marker = f'id="{shape_id}"'
    marker_idx = slide_xml.index(marker)
    shape_start = slide_xml.rfind("<p:sp ", 0, marker_idx)
    shape_end = slide_xml.index("</p:sp>", marker_idx) + len("</p:sp>")
    shape_xml = slide_xml[shape_start:shape_end]

    idx_fill = shape_xml.index('<a:solidFill><a:srgbClr val="112233"/></a:solidFill>')
    idx_line = shape_xml.index('<a:ln w="25400"><a:prstDash val="dash"/>')
    idx_effect = shape_xml.index("<a:effectLst><a:outerShdw")
    idx_geom = shape_xml.index("<a:prstGeom")
    assert idx_fill < idx_line < idx_effect < idx_geom  # noqa: S101


def test_shape_line_defaults_emit_without_solid_fill_fixture(tmp_path: Path) -> None:
    """Test shape line defaults emit without solid fill fixture."""
    out_path = tmp_path / "dml_line_defaults.pptx"

    with Presentation.new("DML Defaults") as prs:
        slide = prs.slides[0]
        slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            properties={"line": {"width_emu": 38100, "dash_style": "round_dot"}},
        )
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert '<a:ln w="38100"><a:prstDash val="sysDot"/></a:ln>' in slide_xml  # noqa: S101
