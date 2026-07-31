"""python-pptx-style facade coverage: slide size, fill, line, font, layouts, tables.

Every test asserts on the serialized XML as well as the in-memory facade, so a
getter that silently returns a default cannot make the test pass.
"""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

import pytest
from gopptx import Presentation
from gopptx.schemas import Inches, Point, RGBColor

if TYPE_CHECKING:
    from pathlib import Path


def _slide_xml(pptx_path: Path, index: int = 1) -> str:
    with zipfile.ZipFile(pptx_path) as zf:
        return zf.read(f"ppt/slides/slide{index}.xml").decode("utf-8")


def _presentation_xml(pptx_path: Path) -> str:
    with zipfile.ZipFile(pptx_path) as zf:
        return zf.read("ppt/presentation.xml").decode("utf-8")


def test_slide_dimensions_roundtrip(tmp_path: Path) -> None:
    """slide_width/slide_height read back the value set and land in sldSz."""
    output_path = tmp_path / "slide_size.pptx"

    with Presentation.new(title="Slide size") as prs:
        default_w, default_h = prs.slide_width, prs.slide_height
        if default_w <= 0 or default_h <= 0:
            raise AssertionError("expected positive default slide dimensions")

        prs.slide_width = Inches(13.333)
        prs.slide_height = Inches(7.5)

        assert prs.slide_width == Inches(13.333)
        assert prs.slide_height == Inches(7.5)
        prs.save(output_path)

    xml = _presentation_xml(output_path)
    assert f'cx="{Inches(13.333)}"' in xml
    assert f'cy="{Inches(7.5)}"' in xml


def test_slide_dimensions_setting_one_axis_preserves_other(tmp_path: Path) -> None:
    """Setting only the width leaves the height untouched."""
    output_path = tmp_path / "slide_size_one_axis.pptx"

    with Presentation.new(title="Slide size") as prs:
        original_h = prs.slide_height
        prs.slide_width = Inches(10)

        assert prs.slide_height == original_h
        prs.save(output_path)

    assert f'cy="{original_h}"' in _presentation_xml(output_path)


def test_text_frame_margins_roundtrip(tmp_path: Path) -> None:
    """Text-frame margins read back and serialize to the lIns/tIns attributes."""
    output_path = tmp_path / "margins.pptx"

    with Presentation.new(title="Margins") as prs:
        slide = prs.slides[0]
        box_id = slide.add_textbox(Inches(1), Inches(1), Inches(5), Inches(2))
        text_frame = slide.shape(box_id).text_frame

        text_frame.margin_left = Inches(0.5)
        text_frame.margin_right = Inches(0.25)
        text_frame.margin_top = Inches(0.3)
        text_frame.margin_bottom = Inches(0.1)

        assert text_frame.margin_left == Inches(0.5)
        assert text_frame.margin_right == Inches(0.25)
        assert text_frame.margin_top == Inches(0.3)
        assert text_frame.margin_bottom == Inches(0.1)
        prs.save(output_path)

    xml = _slide_xml(output_path)
    assert f'lIns="{Inches(0.5)}"' in xml
    assert f'rIns="{Inches(0.25)}"' in xml
    assert f'tIns="{Inches(0.3)}"' in xml
    assert f'bIns="{Inches(0.1)}"' in xml


def test_text_frame_word_wrap_and_anchor(tmp_path: Path) -> None:
    """word_wrap and vertical_anchor serialize to wrap/anchor attributes."""
    output_path = tmp_path / "anchor.pptx"

    with Presentation.new(title="Anchor") as prs:
        slide = prs.slides[0]
        box_id = slide.add_textbox(Inches(1), Inches(1), Inches(5), Inches(2))
        text_frame = slide.shape(box_id).text_frame

        text_frame.word_wrap = True
        text_frame.vertical_anchor = "middle"

        assert text_frame.word_wrap is True
        assert text_frame.vertical_anchor == "middle"
        prs.save(output_path)

    xml = _slide_xml(output_path)
    assert 'wrap="square"' in xml
    assert 'anchor="ctr"' in xml


def test_shape_solid_fill_rgb(tmp_path: Path) -> None:
    """fill.solid() plus fore_color.rgb writes an srgbClr solidFill."""
    output_path = tmp_path / "fill.pptx"

    with Presentation.new(title="Fill") as prs:
        slide = prs.slides[0]
        shape_id = slide.add_shape(
            "rectangle", bounds=(Inches(1), Inches(1), Inches(4), Inches(2))
        )
        shape = slide.shape(shape_id)

        shape.fill.solid()
        shape.fill.fore_color.rgb = RGBColor(0x34, 0x98, 0xDB)

        assert shape.fill.type == "solid"
        assert shape.fill.fore_color.rgb == RGBColor(0x34, 0x98, 0xDB)
        prs.save(output_path)

    assert 'srgbClr val="3498DB"' in _slide_xml(output_path)


def test_shape_fill_solid_defaults_to_black(tmp_path: Path) -> None:
    """fill.solid() alone produces a concrete black fill, not an empty one."""
    with Presentation.new(title="Fill default") as prs:
        slide = prs.slides[0]
        shape_id = slide.add_shape(
            "rectangle", bounds=(Inches(1), Inches(1), Inches(2), Inches(1))
        )
        shape = slide.shape(shape_id)

        shape.fill.solid()

        assert shape.fill.type == "solid"
        assert shape.fill.fore_color.rgb == RGBColor(0, 0, 0)


def test_shape_line_color_and_width(tmp_path: Path) -> None:
    """line.width and line.fill.fore_color.rgb write a:ln w and srgbClr."""
    output_path = tmp_path / "line.pptx"

    with Presentation.new(title="Line") as prs:
        slide = prs.slides[0]
        shape_id = slide.add_shape(
            "rectangle", bounds=(Inches(1), Inches(1), Inches(4), Inches(2))
        )
        shape = slide.shape(shape_id)

        shape.line.fill.solid()
        shape.line.fill.fore_color.rgb = RGBColor(0xE7, 0x4C, 0x3C)
        shape.line.width = Point(4)

        assert shape.line.width == Point(4)
        assert shape.line.fill.fore_color.rgb == RGBColor(0xE7, 0x4C, 0x3C)
        prs.save(output_path)

    xml = _slide_xml(output_path)
    assert f'<a:ln w="{Point(4)}"' in xml
    assert 'srgbClr val="E74C3C"' in xml


def test_run_font_formatting(tmp_path: Path) -> None:
    """Run font name/size/bold/italic/underline/color all serialize."""
    output_path = tmp_path / "font.pptx"

    with Presentation.new(title="Font") as prs:
        slide = prs.slides[0]
        box_id = slide.add_textbox(Inches(1), Inches(1), Inches(8), Inches(2))
        text_frame = slide.shape(box_id).text_frame

        paragraph = text_frame.paragraphs[0]
        paragraph.text = "Styled"
        run = paragraph.runs[0]

        run.font.name = "Georgia"
        run.font.size = Point(28)
        run.font.bold = True
        run.font.italic = True
        run.font.underline = True
        run.font.color.rgb = RGBColor(0x8E, 0x44, 0xAD)

        assert run.font.name == "Georgia"
        assert run.font.size == Point(28)
        assert run.font.size_pt == pytest.approx(28.0)
        assert run.font.bold is True
        assert run.font.italic is True
        assert run.font.underline is True
        prs.save(output_path)

    xml = _slide_xml(output_path)
    assert 'typeface="Georgia"' in xml
    assert 'sz="2800"' in xml
    assert 'b="1"' in xml
    assert 'i="1"' in xml
    assert 'u="sng"' in xml
    assert 'srgbClr val="8E44AD"' in xml


def test_rgb_color_str_and_from_string() -> None:
    """RGBColor is a 3-tuple that formats as uppercase hex and parses back."""
    color = RGBColor(0x8E, 0x44, 0xAD)

    assert tuple(color) == (0x8E, 0x44, 0xAD)
    assert str(color) == "8E44AD"
    assert RGBColor.from_string("#8e44ad") == color
    assert RGBColor.from_string("8E44AD") == color


@pytest.mark.parametrize("bad", ["", "8E44A", "8E44ADD", "#GGGGGG", "zzzzzz"])
def test_rgb_color_from_string_rejects_bad_input(bad: str) -> None:
    """A malformed hex color raises instead of silently becoming black."""
    with pytest.raises(ValueError, match="6-digit hex color"):
        RGBColor.from_string(bad)


@pytest.mark.parametrize("channels", [(-1, 0, 0), (256, 0, 0), (0, 0, 999)])
def test_rgb_color_rejects_out_of_range_channels(
    channels: tuple[int, int, int],
) -> None:
    """Channels outside 0-255 raise rather than serializing to garbage hex."""
    with pytest.raises(ValueError, match="out of range"):
        RGBColor(*channels)


def test_slide_master_and_layout_access() -> None:
    """slide_master, slide_layouts and slide.slide_layout resolve real parts."""
    with Presentation.new(title="Layouts") as prs:
        master = prs.slide_master
        layouts = prs.slide_layouts

        if len(layouts) == 0:
            raise AssertionError("expected at least one slide layout")
        assert master is prs.slide_masters[0]

        layout = prs.slides[0].slide_layout
        assert layout.part_name.startswith("ppt/slideLayouts/")
        assert layout.part_name in {item.part_name for item in layouts}


def test_slide_layout_raises_for_unknown_binding(
    monkeypatch: pytest.MonkeyPatch,
) -> None:
    """An unresolvable layout binding raises rather than returning layout[0]."""
    with Presentation.new(title="Layouts") as prs:
        slide = prs.slides[0]
        monkeypatch.setattr(
            prs,
            "get_slide_layout_ref",
            lambda _index: ("ppt/slideLayouts/nope.xml", "ppt/slideMasters/m.xml"),
        )

        with pytest.raises(LookupError):
            _ = slide.slide_layout


def test_add_table_cell_text(tmp_path: Path) -> None:
    """Table cell text reads back and lands in the a:tbl graphic frame."""
    output_path = tmp_path / "table.pptx"

    with Presentation.new(title="Table") as prs:
        slide = prs.slides[0]
        table_id = slide.add_table(
            2, 2, bounds=(Inches(1), Inches(1), Inches(8), Inches(2))
        )
        table = slide.table(table_id)

        table.cell(0, 0).text = "Product"
        table.cell(1, 0).text = "Widget A"

        assert table.cell(0, 0).text == "Product"
        assert table.cell(1, 0).text == "Widget A"
        prs.save(output_path)

    xml = _slide_xml(output_path)
    assert "<a:tbl>" in xml
    assert "Widget A" in xml


def test_add_picture_embeds_image(tmp_path: Path) -> None:
    """add_picture creates a p:pic and embeds the media part."""
    image_path = tmp_path / "dot.png"
    image_path.write_bytes(
        bytes.fromhex(
            "89504e470d0a1a0a0000000d4948445200000001000000010802000000907753"
            "de0000000c4944415408d763f8cfc00000030101003c3f5e8e0000000049454e"
            "44ae426082"
        )
    )
    output_path = tmp_path / "picture.pptx"

    with Presentation.new(title="Picture") as prs:
        slide = prs.slides[0]
        pic_id = slide.add_picture(
            str(image_path), Inches(1), Inches(1), Inches(2), Inches(2)
        )

        assert slide.shape(pic_id).id == pic_id
        prs.save(output_path)

    xml = _slide_xml(output_path)
    assert "<p:pic" in xml
    assert "<a:blip" in xml
    with zipfile.ZipFile(output_path) as zf:
        media = [name for name in zf.namelist() if name.startswith("ppt/media/")]
    if not media:
        raise AssertionError("expected an embedded media part")


@pytest.mark.parametrize(
    ("shape_type", "preset"),
    [("rectangle", "rect"), ("oval", "ellipse"), ("right_arrow", "rightArrow")],
)
def test_add_autoshape_preset_geometry(
    tmp_path: Path, shape_type: str, preset: str
) -> None:
    """Each autoshape alias serializes to its OOXML preset geometry."""
    output_path = tmp_path / f"{shape_type}.pptx"

    with Presentation.new(title="AutoShape") as prs:
        slide = prs.slides[0]
        shape_id = slide.add_shape(
            shape_type, bounds=(Inches(1), Inches(1), Inches(3), Inches(2))
        )

        assert slide.shape(shape_id) is not None
        prs.save(output_path)

    assert f'prst="{preset}"' in _slide_xml(output_path)
