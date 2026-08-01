"""Placeholder pictures can preserve their aspect ratio (Issue #176)."""

import struct
import zipfile
import zlib

import pytest
from gopptx.api_errors import GopptxError
from gopptx.presentation.presentation import Presentation

EMU_PER_POINT = 12700
# A 4-in square box at one inch in from the top-left, in points.
BOX = (72.0, 72.0, 288.0, 288.0)
BOX_EMU = tuple(int(value * EMU_PER_POINT) for value in BOX)


def _wide_png(tmp_path, width=800, height=200):
    """Write a width x height PNG so the fit math has real dimensions to read."""

    def chunk(tag: bytes, data: bytes) -> bytes:
        return (
            struct.pack(">I", len(data))
            + tag
            + data
            + struct.pack(">I", zlib.crc32(tag + data) & 0xFFFFFFFF)
        )

    raw = b"".join(
        b"\x00" + bytes(v for x in range(width) for v in (x % 256, y % 256, 90))
        for y in range(height)
    )
    path = tmp_path / f"image_{width}x{height}.png"
    path.write_bytes(
        b"\x89PNG\r\n\x1a\n"
        + chunk(b"IHDR", struct.pack(">IIBBBBB", width, height, 8, 2, 0, 0, 0))
        + chunk(b"IDAT", zlib.compress(raw))
        + chunk(b"IEND", b"")
    )
    return path


def _slide_xml(path) -> str:
    with zipfile.ZipFile(path) as archive:
        return archive.read("ppt/slides/slide1.xml").decode("utf-8")


def _insert(tmp_path, image_path, fit, name):
    out = tmp_path / name
    with Presentation.new("Placeholder Image Fit") as pres:
        placeholder = next(iter(pres.slides[0].placeholders))
        placeholder.insert_picture(str(image_path), bounds=BOX, fit=fit)
        pres.save(out)
    return _slide_xml(out)


def test_stretch_fills_the_box_and_is_the_default(tmp_path):
    image = _wide_png(tmp_path)
    stretched = _insert(tmp_path, image, "stretch", "stretch.pptx")
    assert f'<a:ext cx="{BOX_EMU[2]}" cy="{BOX_EMU[3]}"/>' in stretched
    assert "<a:srcRect" not in stretched

    out = tmp_path / "default.pptx"
    with Presentation.new("Default Fit") as pres:
        placeholder = next(iter(pres.slides[0].placeholders))
        placeholder.insert_picture(str(image), bounds=BOX)
        pres.save(out)
    assert _slide_xml(out) == stretched


def test_contain_preserves_aspect_ratio_and_centres(tmp_path):
    image = _wide_png(tmp_path)
    xml = _insert(tmp_path, image, "contain", "contain.pptx")

    # A 4:1 image in a square box keeps the full width at a quarter of the height.
    expected_cy = BOX_EMU[2] // 4
    expected_y = BOX_EMU[1] + (BOX_EMU[3] - expected_cy) // 2
    assert f'<a:ext cx="{BOX_EMU[2]}" cy="{expected_cy}"/>' in xml
    assert f'<a:off x="{BOX_EMU[0]}" y="{expected_y}"/>' in xml
    assert "<a:srcRect" not in xml


def test_cover_fills_the_box_and_crops_the_overflow(tmp_path):
    image = _wide_png(tmp_path)
    xml = _insert(tmp_path, image, "cover", "cover.pptx")

    # Only a quarter of a 4:1 image's width fits a square box, so 75% is cropped,
    # split evenly across the two sides.
    assert f'<a:ext cx="{BOX_EMU[2]}" cy="{BOX_EMU[3]}"/>' in xml
    assert '<a:srcRect l="37500" r="37500"/>' in xml


def test_cover_crops_a_tall_image_vertically(tmp_path):
    image = _wide_png(tmp_path, width=200, height=800)
    xml = _insert(tmp_path, image, "cover", "cover_tall.pptx")
    assert '<a:srcRect t="37500" b="37500"/>' in xml


def test_matching_aspect_ratio_needs_no_crop(tmp_path):
    image = _wide_png(tmp_path, width=400, height=400)
    xml = _insert(tmp_path, image, "cover", "cover_square.pptx")
    assert f'<a:ext cx="{BOX_EMU[2]}" cy="{BOX_EMU[3]}"/>' in xml
    assert "<a:srcRect" not in xml


def test_unknown_fit_mode_is_rejected(tmp_path):
    image = _wide_png(tmp_path)
    with Presentation.new("Bad Fit") as pres:
        placeholder = next(iter(pres.slides[0].placeholders))
        with pytest.raises(GopptxError, match="squish"):
            placeholder.insert_picture(str(image), bounds=BOX, fit="squish")
