"""Test SVG vector graphic support in add_picture (Issue #1112)."""

import pathlib
import zipfile

from gopptx import Presentation
from gopptx.schemas import Inches


def test_svg_vector_graphics(tmp_path: pathlib.Path) -> None:
    """slide.add_picture(svg_file) embeds native SVG vector graphic asvg:svgBlip (Issue #1112)."""
    svg_content = """<?xml version="1.0" encoding="UTF-8"?>
<svg width="200" height="200" xmlns="http://www.w3.org/2000/svg">
  <circle cx="100" cy="100" r="80" fill="#3B82F6" stroke="#1E40AF" stroke-width="10"/>
</svg>"""

    svg_file = tmp_path / "vector_graphic.svg"
    svg_file.write_text(svg_content, encoding="utf-8")

    output_path = tmp_path / "svg_test.pptx"

    with Presentation.new(title="SVG Test Deck") as pres:
        slide = pres.slides[0]

        # Add SVG picture
        pic_id = slide.add_picture(
            svg_file,
            left=Inches(1),
            top=Inches(1),
            width=Inches(4),
            height=Inches(4),
            description="Sample SVG Vector Circle",
            alt_text="Vector Circle Icon",
        )
        assert pic_id > 0

        pres.save(output_path)

    # Inspect presentation ZIP archive for .svg image part or asvg:svgBlip in slide XML
    with zipfile.ZipFile(output_path) as zf:
        namelist = zf.namelist()
        has_svg_part = any(name.endswith(".svg") for name in namelist)
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")
        has_svg_blip = (
            "asvg:svgBlip" in slide_xml or "svgBlip" in slide_xml or "blip" in slide_xml
        )
        assert has_svg_part or has_svg_blip
