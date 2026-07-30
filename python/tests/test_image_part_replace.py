"""ImagePartProxy.replace must actually swap the bytes behind the shape."""

from __future__ import annotations

from io import BytesIO
from typing import TYPE_CHECKING

import pytest
from gopptx import Presentation
from PIL import Image

if TYPE_CHECKING:
    from pathlib import Path


def _png_bytes(color: str, size: tuple[int, int]) -> bytes:
    buffer = BytesIO()
    Image.new("RGB", size, color).save(buffer, format="PNG")
    return buffer.getvalue()


def _deck_with_picture(path: Path, data: bytes) -> tuple[Presentation, int]:
    path.write_bytes(data)
    pres = Presentation.new(title="Image Replace")
    slide = pres.slides[0]
    shape_id = slide.add_image(str(path), (1000000, 1000000, 2000000, 1000000))
    return pres, shape_id


def test_replace_swaps_image_bytes(tmp_path: Path) -> None:
    original = _png_bytes("blue", (144, 72))
    replacement = _png_bytes("red", (96, 96))

    pres, shape_id = _deck_with_picture(tmp_path / "original.png", original)
    with pres:
        image = pres.slides[0].shape(shape_id).image
        assert image.blob == original

        image.replace(replacement)

        refreshed = pres.slides[0].shape(shape_id).image
        assert refreshed.blob == replacement
        assert refreshed.blob != original


def test_replace_accepts_a_path(tmp_path: Path) -> None:
    original = _png_bytes("blue", (144, 72))
    replacement_path = tmp_path / "replacement.png"
    replacement_path.write_bytes(_png_bytes("green", (32, 32)))

    pres, shape_id = _deck_with_picture(tmp_path / "original.png", original)
    with pres:
        pres.slides[0].shape(shape_id).image.replace(replacement_path)
        assert (
            pres.slides[0].shape(shape_id).image.blob == replacement_path.read_bytes()
        )


def test_replace_reports_a_shape_without_an_image(tmp_path: Path) -> None:
    original = _png_bytes("blue", (144, 72))
    pres, _ = _deck_with_picture(tmp_path / "original.png", original)
    with pres:
        textbox_id = pres.slides[0].add_textbox(1000000, 1000000, 2000000, 1000000)
        with pytest.raises(Exception, match=r"(?i)image|relationship"):
            pres.slides[0].shape(textbox_id).image.replace(original)
