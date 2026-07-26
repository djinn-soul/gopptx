"""Test SVG vector graphic embedding support (Issue #1112)."""

import pathlib
import zipfile

from gopptx import Presentation

TINY_SVG_BYTES = b'<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100"><rect width="100" height="100" fill="red"/></svg>'


def test_add_picture_with_svg_bytes(tmp_path: pathlib.Path) -> None:
    """Adding an SVG image via bytes creates an image/svg+xml part and asvg:svgBlip drawingML extension (Issue #1112)."""
    output_path = tmp_path / "svg_test.pptx"

    with Presentation.new(title="SVG Test Deck") as pres:
        slide = pres.slides[0]

        shape_id = slide.add_picture(
            TINY_SVG_BYTES,
            left=1000000,
            top=1000000,
            width=2000000,
            height=2000000,
            image_format="svg",
            description="Vector SVG Logo",
        )
        assert shape_id > 0

        pres.save(output_path)

    # Inspect presentation ZIP contents
    with zipfile.ZipFile(output_path) as zf:
        namelist = zf.namelist()
        svg_parts = [name for name in namelist if name.endswith(".svg")]
        assert len(svg_parts) > 0, "expected embedded .svg media part"

        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")
        assert "asvg:svgBlip" in slide_xml, (
            "expected asvg:svgBlip drawingML extension element"
        )
        assert 'descr="Vector SVG Logo"' in slide_xml


def test_add_picture_with_svg_file(tmp_path: pathlib.Path) -> None:
    """Adding an SVG image file path creates an image/svg+xml part and slide picture shape."""
    svg_file = tmp_path / "logo.svg"
    svg_file.write_bytes(TINY_SVG_BYTES)

    output_path = tmp_path / "svg_file_test.pptx"

    with Presentation.new(title="SVG File Test Deck") as pres:
        slide = pres.slides[0]

        shape_id = slide.add_picture(
            str(svg_file),
            left=1000000,
            top=1000000,
            width=2000000,
            height=2000000,
        )
        assert shape_id > 0

        pres.save(output_path)

    with Presentation(output_path) as pres:
        shapes = pres.slides[0].list_shapes()
        assert len(shapes) > 0
