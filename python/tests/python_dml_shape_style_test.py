"""Parity tests for generic shape-level DML fill/line controls."""

from __future__ import annotations

import zipfile
from pathlib import Path

from gopptx import Presentation


def test_shape_fill_and_line_round_trip_fixture(tmp_path: Path) -> None:
    out_path = tmp_path / "dml_shape_style.pptx"

    with Presentation.new("DML Style") as prs:
        slide = prs.slides[0]
        shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            text="style",
            properties={
                "fill": {"solid": "FF0000"},
                "line": {"color": "00FF00", "width_emu": 25400, "dash_style": "dash"},
            },
        )
        slide.update_shape(
            shape_id,
            {
                "fill": {"solid": "#112233"},
                "line": {
                    "color": "334455",
                    "width_emu": 38100,
                    "dash_style": "long_dash_dot",
                },
            },
        )
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert 'val="112233"' in slide_xml  # noqa: S101
    assert 'a:ln w="38100"' in slide_xml  # noqa: S101
    assert '<a:prstDash val="lgDashDot"/>' in slide_xml  # noqa: S101
    assert 'val="334455"' in slide_xml  # noqa: S101
