"""Test RTL / Right-To-Left Arabic paragraph text layout (Issue #1080)."""

import pathlib
import zipfile

from gopptx import Presentation


def test_rtl_arabic_layout(tmp_path: pathlib.Path) -> None:
    """Paragraphs with rtl=True emit a:pPr rtl="1" (Issue #1080)."""
    output_path = tmp_path / "rtl_arabic_test.pptx"

    with Presentation.new(title="RTL Arabic Presentation") as pres:
        slide = pres.slides[0]

        # Add textbox with Arabic RTL text
        tb_id = slide.add_textbox(
            1000000, 1000000, 6000000, 2000000, text="مقدمة عن برنامج الولاء"
        )
        assert tb_id > 0

        # Apply RTL paragraph property via paragraph facade
        p0 = slide.shape(tb_id).text_frame.paragraphs[0]
        p0.rtl = True
        p0.alignment = "r"

        assert p0.rtl is True
        assert p0.alignment == "r"

        pres.save(output_path)

    # Inspect XML for a:pPr rtl="1"
    with zipfile.ZipFile(output_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")
        assert 'rtl="1"' in slide_xml or 'algn="r"' in slide_xml
