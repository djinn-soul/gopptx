"""EMF/WMF picture support regression tests (upstream python-pptx issues #1042, #25)."""

import pathlib
import zipfile

from gopptx.presentation.presentation import Presentation
from gopptx.schemas import Inches


def _emf_bytes() -> bytes:
    """Return a minimal EMR_HEADER carrying the ' EMF' signature at offset 40."""
    data = bytearray(88)
    data[0:4] = b"\x01\x00\x00\x00"
    data[40:44] = b" EMF"
    return bytes(data)


def _wmf_bytes() -> bytes:
    """Return a placeable WMF header."""
    return b"\xd7\xcd\xc6\x9a" + bytes(20)


def test_add_picture_emf_declares_content_type(tmp_path: pathlib.Path) -> None:
    """An embedded EMF gets a Default content-type entry, so the package opens."""
    img_path = tmp_path / "diagram.emf"
    img_path.write_bytes(_emf_bytes())
    output_path = tmp_path / "emf_output.pptx"

    with Presentation.new("EMF Test") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_picture(
            img_path,
            left=Inches(1),
            top=Inches(1),
            width=Inches(3),
            height=Inches(2),
        )
        assert shape_id > 0
        pres.save(output_path)

    with zipfile.ZipFile(output_path) as archive:
        content_types = archive.read("[Content_Types].xml").decode("utf-8")
        media = [n for n in archive.namelist() if n.startswith("ppt/media/")]

    assert 'Extension="emf" ContentType="image/x-emf"' in content_types
    assert any(name.endswith(".emf") for name in media)


def test_add_picture_wmf_declares_content_type(tmp_path: pathlib.Path) -> None:
    """An embedded WMF gets a Default content-type entry, so the package opens."""
    img_path = tmp_path / "legacy.wmf"
    img_path.write_bytes(_wmf_bytes())
    output_path = tmp_path / "wmf_output.pptx"

    with Presentation.new("WMF Test") as pres:
        slide = pres.slides[0]
        slide.add_picture(
            img_path,
            left=Inches(1),
            top=Inches(1),
            width=Inches(3),
            height=Inches(2),
        )
        pres.save(output_path)

    with zipfile.ZipFile(output_path) as archive:
        content_types = archive.read("[Content_Types].xml").decode("utf-8")

    assert 'Extension="wmf" ContentType="image/x-wmf"' in content_types


def test_image_part_reports_emf_content_type_and_ext(tmp_path: pathlib.Path) -> None:
    """Picture.image reports image/x-emf rather than application/octet-stream."""
    img_path = tmp_path / "diagram.emf"
    img_path.write_bytes(_emf_bytes())
    output_path = tmp_path / "emf_metadata.pptx"

    with Presentation.new("EMF Metadata") as pres:
        slide = pres.slides[0]
        slide.add_picture(
            img_path,
            left=Inches(1),
            top=Inches(1),
            width=Inches(3),
            height=Inches(2),
        )
        pres.save(output_path)

    reopened = Presentation()
    reopened.open(output_path)
    try:
        picture = next(
            shape for shape in reopened.slides[0].shapes if shape.shape_type == "pic"
        )
        image = picture.image
        assert image.content_type == "image/x-emf"
        assert image.ext == "emf"
        assert image.blob == _emf_bytes()
    finally:
        reopened.close()
