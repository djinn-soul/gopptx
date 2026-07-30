"""Test picture slide background support (Issue #1095)."""

import pathlib
import zipfile

from gopptx import Presentation


def test_picture_slide_background(tmp_path: pathlib.Path) -> None:
    """slide.background.set_picture(path_or_bytes) sets full-slide picture background (Issue #1095)."""
    # Create 1x1 dummy PNG bytes
    png_bytes = (
        b"\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01"
        b"\x08\x06\x00\x00\x00\x1f\x15c4\x00\x00\x00\rIDATx\x9cc\xf8\xff\xff?"
        b"\x03\x00\x05\xfe\x02\xfe\xa7\x96a\x1d\x00\x00\x00\x00IEND\xaeB`\x82"
    )

    img_file = tmp_path / "bg_image.png"
    img_file.write_bytes(png_bytes)

    output_path = tmp_path / "picture_bg_test.pptx"

    with Presentation.new(title="Picture BG Deck") as pres:
        slide = pres.slides[0]

        # Apply picture background via file path
        slide.background.set_picture(img_file)

        pres.save(output_path)

    # Inspect slide XML for p:bg and blipFill
    with zipfile.ZipFile(output_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")
        assert "p:bg" in slide_xml or "a:blipFill" in slide_xml or "blip" in slide_xml
