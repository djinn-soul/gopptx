"""Font color proxies return RGBColor, matching fill and line color proxies."""

from gopptx.presentation.presentation import Presentation
from gopptx.schemas import Inches, RGBColor

PURPLE = RGBColor(0x8E, 0x44, 0xAD)


def _paragraph(slide, shape_id):
    return slide.shape(shape_id).text_frame.paragraphs[0]


def test_run_font_color_rgb_is_rgbcolor(tmp_path):
    path = str(tmp_path / "run_font_color.pptx")
    with Presentation.new("Run Font Color") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_textbox(Inches(1), Inches(2), Inches(6), Inches(1))
        paragraph = _paragraph(slide, shape_id)
        paragraph.text = "colored"
        run = paragraph.runs[0]
        run.font.color.rgb = PURPLE

        assert run.font.color.rgb == PURPLE
        assert isinstance(run.font.color.rgb, RGBColor)
        assert run.font.color.type == "RGB"
        pres.save(path)

    with Presentation(path) as reopened:
        reloaded = list(reopened.slides[0].shapes)[-1].text_frame.paragraphs[0].runs[0]
        assert reloaded.font.color.rgb == PURPLE


def test_paragraph_font_color_rgb_is_rgbcolor(tmp_path):
    path = str(tmp_path / "paragraph_font_color.pptx")
    with Presentation.new("Paragraph Font Color") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_textbox(Inches(1), Inches(2), Inches(6), Inches(1))
        paragraph = _paragraph(slide, shape_id)
        paragraph.text = "colored"
        paragraph.font.color.rgb = PURPLE

        assert paragraph.font.color.rgb == PURPLE
        assert isinstance(paragraph.font.color.rgb, RGBColor)
        for run in paragraph.runs:
            assert run.font.color.rgb == PURPLE
        pres.save(path)

    with Presentation(path) as reopened:
        reloaded = list(reopened.slides[0].shapes)[-1].text_frame.paragraphs[0]
        assert reloaded.font.color.rgb == PURPLE


def test_font_color_matches_fill_color_proxy_type():
    with Presentation.new("Color Proxy Parity") as pres:
        slide = pres.slides[0]
        rect_id = slide.add_shape(
            "rectangle", bounds=(Inches(1), Inches(2), Inches(3), Inches(2))
        )
        shape = slide.shape(rect_id)
        shape.fill.solid()
        shape.fill.fore_color.rgb = PURPLE

        text_id = slide.add_textbox(Inches(1), Inches(5), Inches(6), Inches(1))
        paragraph = _paragraph(slide, text_id)
        paragraph.text = "colored"
        paragraph.runs[0].font.color.rgb = PURPLE

        assert type(paragraph.runs[0].font.color.rgb) is type(shape.fill.fore_color.rgb)


def test_unset_font_color_reads_none():
    with Presentation.new("Unset Font Color") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_textbox(Inches(1), Inches(2), Inches(6), Inches(1))
        paragraph = _paragraph(slide, shape_id)
        paragraph.text = "plain"
        assert paragraph.runs[0].font.color.rgb is None
        assert paragraph.runs[0].font.color.type is None
