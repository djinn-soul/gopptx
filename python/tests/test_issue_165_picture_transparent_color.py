"""Regression coverage for upstream python-pptx issue #165."""
# pyright: reportImplicitStringConcatenation=false

from __future__ import annotations

from typing import TYPE_CHECKING
from zipfile import ZipFile

from gopptx import Presentation
from gopptx.schemas import Inches

if TYPE_CHECKING:
    from pathlib import Path

_RED_PNG = bytes.fromhex(
    "89504e470d0a1a0a0000000d4948445200000001000000010802000000907753"
    "de0000000c4944415408d763f8cfc00000030101003c3f5e8e0000000049454e"
    "44ae426082"
)


def test_picture_color_can_be_made_transparent(tmp_path: Path) -> None:
    output = tmp_path / "transparent-color.pptx"
    with Presentation.new(title="Issue 165") as prs:
        slide = prs.slides[0]
        slide.add_picture(
            _RED_PNG,
            Inches(1),
            Inches(1),
            Inches(3),
            Inches(2),
            format="png",
        )
        slide.shapes[-1].make_color_transparent("FF0000")
        prs.save(str(output))

    with ZipFile(output) as package:
        xml = package.read("ppt/slides/slide1.xml").decode()
    assert '<a:srgbClr val="FF0000"/>' in xml
    assert '<a:alpha val="0"/>' in xml
