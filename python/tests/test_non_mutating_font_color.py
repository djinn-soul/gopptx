"""Test non-mutating text run font color read access (Issue #1111)."""

import pathlib
import zipfile

from gopptx import Presentation
from gopptx.schemas import Inches


def test_non_mutating_font_color_read(tmp_path: pathlib.Path) -> None:
    """Reading run.font.color.rgb or run.font.color does NOT mutate XML or inject empty solidFill (Issue #1111)."""
    output_path = tmp_path / "non_mutating_font_color_test.pptx"

    with Presentation.new(title="Theme Inherited Text") as pres:
        slide = pres.slides[0]

        # Add textbox without explicit font color (inherited from theme)
        tb_id = slide.add_textbox(Inches(1), Inches(1), Inches(8), Inches(2))
        tf = slide.shape(tb_id).text_frame
        p = tf.paragraphs[0]
        run = p.runs[0]
        run.text = "Inherited Theme Color Text"

        # Read-only access to font color - MUST NOT MUTATE XML
        color_rgb = run.font.color.rgb
        color_type = run.font.color.type
        assert color_rgb is None
        assert color_type is None

        pres.save(output_path)

    # Verify that ppt/slides/slide1.xml contains NO empty solidFill injected into rPr
    with zipfile.ZipFile(output_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")
        assert "<a:solidFill/>" not in slide_xml
        assert "<a:solidFill></a:solidFill>" not in slide_xml
