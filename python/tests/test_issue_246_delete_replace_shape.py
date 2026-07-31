"""Issue #246: delete a shape, replace an image in place, or reorder shapes.

The report asked for any one of three things: a way to replace an image, a way
to delete the stale shape, or a way to reorder so a freshly added picture is not
stuck at the end. All three are asserted here against the saved package.
"""

from __future__ import annotations

import re
import zipfile
from typing import TYPE_CHECKING

from gopptx import Presentation
from gopptx.schemas import Inches

if TYPE_CHECKING:
    from pathlib import Path

# 1x1 PNGs, one red and one blue, so a swap is visible in the media part bytes.
_RED_PNG = bytes.fromhex(
    "89504e470d0a1a0a0000000d4948445200000001000000010802000000907753"
    "de0000000c4944415408d763f8cfc00000030101003c3f5e8e0000000049454e"
    "44ae426082"
)
_BLUE_PNG = bytes.fromhex(
    "89504e470d0a1a0a0000000d4948445200000001000000010802000000907753"
    "de0000000c4944415408d76360f8cf000000030101003ea55e8e0000000049454e"
    "44ae426082"
)


def _slide_xml(pptx_path: Path) -> str:
    with zipfile.ZipFile(pptx_path) as zf:
        return zf.read("ppt/slides/slide1.xml").decode("utf-8")


def _media_bytes(pptx_path: Path) -> list[bytes]:
    with zipfile.ZipFile(pptx_path) as zf:
        return [zf.read(n) for n in sorted(zf.namelist()) if n.startswith("ppt/media/")]


def test_delete_removes_the_stale_shape(tmp_path: Path) -> None:
    """Deleting is the simplest answer to 'the old one is still in there too'."""
    output_path = tmp_path / "deleted.pptx"

    with Presentation.new(title="Issue 246 delete") as prs:
        slide = prs.slides[0]
        stale = slide.shapes.add_textbox(
            Inches(1), Inches(1), Inches(3), Inches(1), text="STALE_SHAPE"
        )
        slide.shapes.add_textbox(
            Inches(1), Inches(3), Inches(3), Inches(1), text="FRESH_SHAPE"
        )

        stale.delete()
        prs.save(output_path)

    xml = _slide_xml(output_path)
    assert "STALE_SHAPE" not in xml
    assert "FRESH_SHAPE" in xml


def test_image_is_replaced_without_moving_the_picture(tmp_path: Path) -> None:
    """swap_image_by_index repoints the picture at new bytes, keeping its place."""
    image_path = tmp_path / "red.png"
    image_path.write_bytes(_RED_PNG)
    output_path = tmp_path / "swapped.pptx"

    with Presentation.new(title="Issue 246 replace") as prs:
        slide = prs.slides[0]
        slide.add_picture(str(image_path), Inches(1), Inches(1), Inches(2), Inches(2))
        slide.shapes.add_textbox(
            Inches(1), Inches(4), Inches(3), Inches(1), text="TRAILING_SHAPE"
        )
        before = [shape.id for shape in slide.shapes]

        prs.swap_image_by_index(0, 0, _BLUE_PNG, "png")

        assert [shape.id for shape in prs.slides[0].shapes] == before
        prs.save(output_path)

    xml = _slide_xml(output_path)
    # The picture keeps its place: the textbox added after it is still last.
    assert xml.index("<p:pic") < xml.index("TRAILING_SHAPE")

    embed = re.search(r'r:embed="([^"]+)"', xml)
    assert embed is not None, xml
    with zipfile.ZipFile(output_path) as zf:
        rels = zf.read("ppt/slides/_rels/slide1.xml.rels").decode("utf-8")
    target = re.search(rf'Id="{embed.group(1)}"[^>]*Target="\.\./(media/[^"]+)"', rels)
    assert target is not None, rels

    with zipfile.ZipFile(output_path) as zf:
        assert zf.read(f"ppt/{target.group(1)}") == _BLUE_PNG

    assert _BLUE_PNG in _media_bytes(output_path)


def test_replaced_image_leaves_its_old_media_part_behind(tmp_path: Path) -> None:
    """Documents current behaviour: the superseded part stays in the package.

    The swap adds a new media part and repoints the relationship, so the picture
    renders correctly, but the old bytes are left orphaned. Harmless for one
    swap; repeated swaps grow the file.
    """
    image_path = tmp_path / "red.png"
    image_path.write_bytes(_RED_PNG)
    output_path = tmp_path / "orphan.pptx"

    with Presentation.new(title="Issue 246 orphan") as prs:
        prs.slides[0].add_picture(
            str(image_path), Inches(1), Inches(1), Inches(2), Inches(2)
        )
        prs.swap_image_by_index(0, 0, _BLUE_PNG, "png")
        prs.save(output_path)

    with zipfile.ZipFile(output_path) as zf:
        rels = zf.read("ppt/slides/_rels/slide1.xml.rels").decode("utf-8")
    referenced = set(re.findall(r'Target="\.\./(media/[^"]+)"', rels))

    assert _RED_PNG in _media_bytes(output_path)
    assert len(referenced) == 1, referenced


def test_reordering_moves_a_late_addition_into_place(tmp_path: Path) -> None:
    """A shape added last can be moved, answering the 'stuck at the end' part."""
    output_path = tmp_path / "reordered.pptx"

    with Presentation.new(title="Issue 246 reorder") as prs:
        slide = prs.slides[0]
        slide.shapes.add_textbox(
            Inches(1), Inches(1), Inches(3), Inches(1), text="ORD_FIRST"
        )
        late = slide.shapes.add_textbox(
            Inches(1), Inches(3), Inches(3), Inches(1), text="ORD_LATE"
        )

        late.move_to_back()

        assert late.z_order == 0
        prs.save(output_path)

    order = re.findall(r"ORD_(?:FIRST|LATE)", _slide_xml(output_path))
    assert order == ["ORD_LATE", "ORD_FIRST"]
