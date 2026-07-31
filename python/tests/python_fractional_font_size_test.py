"""Fractional point font sizes survive the bridge and a save/reload round-trip."""

import math
import zipfile

from gopptx.presentation.presentation import Presentation
from gopptx.schemas import Inches


def _first_run(slide):
    shape = list(slide.shapes)[-1]
    return shape.text_frame.paragraphs[0].runs[0]


def _assert_size_pt(actual, expected):
    assert actual is not None
    assert math.isclose(actual, expected, rel_tol=0, abs_tol=1e-9)


def test_half_point_font_size_round_trips(tmp_path):
    path = str(tmp_path / "fractional_size.pptx")
    with Presentation.new("Fractional Size") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_textbox(Inches(1), Inches(2), Inches(6), Inches(1))
        paragraph = slide.shape(shape_id).text_frame.paragraphs[0]
        paragraph.text = "half point"
        paragraph.runs[0].font.size_pt = 11.5
        _assert_size_pt(paragraph.runs[0].font.size_pt, 11.5)
        pres.save(path)

    with zipfile.ZipFile(path) as archive:
        slide_xml = archive.read("ppt/slides/slide1.xml").decode("utf-8")
    assert 'sz="1150"' in slide_xml

    with Presentation(path) as reopened:
        _assert_size_pt(_first_run(reopened.slides[0]).font.size_pt, 11.5)


def test_paragraph_font_size_pt_applies_to_every_run(tmp_path):
    path = str(tmp_path / "paragraph_size.pptx")
    with Presentation.new("Paragraph Size") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_textbox(Inches(1), Inches(2), Inches(6), Inches(1))
        paragraph = slide.shape(shape_id).text_frame.paragraphs[0]
        paragraph.text = "styled"
        paragraph.font.size_pt = 18.5
        _assert_size_pt(paragraph.font.size_pt, 18.5)
        pres.save(path)

    with Presentation(path) as reopened:
        _assert_size_pt(_first_run(reopened.slides[0]).font.size_pt, 18.5)
