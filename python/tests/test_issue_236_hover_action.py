"""Issue #236: BaseShape.hover_action.

The hover element inside p:cNvPr is a:hlinkHover (CT_NonVisualDrawingProps).
a:hlinkMouseOver is the run-level spelling; PowerPoint round-trips it there but
never reads it, so these tests pin the element name, not just the round trip.
"""

from __future__ import annotations

import re
import zipfile
from typing import TYPE_CHECKING

from gopptx import Presentation
from gopptx.schemas import Inches

if TYPE_CHECKING:
    from pathlib import Path

_CNVPR_RE = re.compile(r"<p:cNvPr[^>]*>.*?</p:cNvPr>", re.DOTALL)


def _slide_xml(pptx_path: Path) -> str:
    with zipfile.ZipFile(pptx_path) as zf:
        return zf.read("ppt/slides/slide1.xml").decode("utf-8")


def _rels_xml(pptx_path: Path) -> str:
    with zipfile.ZipFile(pptx_path) as zf:
        return zf.read("ppt/slides/_rels/slide1.xml.rels").decode("utf-8")


def test_hover_action_writes_hlink_hover(tmp_path: Path) -> None:
    """A shape hover action lands as a:hlinkHover, not a:hlinkMouseOver."""
    output_path = tmp_path / "hover.pptx"

    with Presentation.new(title="Issue 236") as prs:
        slide = prs.slides[0]
        slide.add_shape(
            "rectangle",
            bounds=(Inches(1), Inches(1), Inches(3), Inches(1)),
            hover_action={"address": "https://example.com/hover"},
        )
        prs.save(output_path)

    xml = _slide_xml(output_path)
    assert "<a:hlinkHover" in xml
    assert "hlinkMouseOver" not in xml


def test_hover_and_click_coexist(tmp_path: Path) -> None:
    """Both actions are written, each against its own relationship."""
    output_path = tmp_path / "hover_click.pptx"

    with Presentation.new(title="Issue 236 both") as prs:
        slide = prs.slides[0]
        slide.add_shape(
            "rectangle",
            bounds=(Inches(1), Inches(1), Inches(3), Inches(1)),
            click_action={"address": "https://example.com/click"},
            hover_action={"address": "https://example.com/hover"},
        )
        prs.save(output_path)

    xml = _slide_xml(output_path)
    click = re.search(r'<a:hlinkClick r:id="([^"]+)"', xml)
    hover = re.search(r'<a:hlinkHover r:id="([^"]+)"', xml)

    assert click is not None
    assert hover is not None
    assert click.group(1) != hover.group(1)

    rels = _rels_xml(output_path)
    assert f'Id="{click.group(1)}"' in rels
    assert f'Id="{hover.group(1)}"' in rels
    assert "https://example.com/hover" in rels


def test_hover_action_reads_back(tmp_path: Path) -> None:
    """The facade reports the hover address it wrote."""
    output_path = tmp_path / "hover_readback.pptx"

    with Presentation.new(title="Issue 236 readback") as prs:
        slide = prs.slides[0]
        shape_id = slide.add_shape(
            "rectangle",
            bounds=(Inches(1), Inches(1), Inches(3), Inches(1)),
            hover_action={"address": "https://example.com/hover"},
        )

        assert slide.shape(shape_id).hover_action.address == "https://example.com/hover"
        prs.save(output_path)

    assert "<a:hlinkHover" in _slide_xml(output_path)


def test_hover_action_is_inside_cnvpr(tmp_path: Path) -> None:
    """The element sits in p:cNvPr, where PowerPoint looks for it."""
    output_path = tmp_path / "hover_placement.pptx"

    with Presentation.new(title="Issue 236 placement") as prs:
        slide = prs.slides[0]
        slide.add_shape(
            "rectangle",
            bounds=(Inches(1), Inches(1), Inches(3), Inches(1)),
            hover_action={"address": "https://example.com/hover"},
        )
        prs.save(output_path)

    blocks = [
        block
        for block in _CNVPR_RE.findall(_slide_xml(output_path))
        if "hlinkHover" in block
    ]

    assert len(blocks) == 1, blocks
