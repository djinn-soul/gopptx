"""Advanced text-layout edge case fixtures for parity coverage."""

from __future__ import annotations

import zipfile
from pathlib import Path

from gopptx import ParagraphProps, Presentation, Run, TextFrameProps


def test_complex_text_layout_fixture_emits_expected_tokens(tmp_path: Path) -> None:
    out_path = tmp_path / "text_layout_edge_case.pptx"

    with Presentation.new("Text Layout Fixtures") as prs:
        slide = prs.slides[0]
        run = Run("Hyper")
        run.hyperlink.address = "https://example.com/"
        run.small_caps = True
        run.subscript = True

        slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            runs=[run, Run(" layout")],
            paragraph=ParagraphProps(indent=228600, hanging=114300),
            text_frame=TextFrameProps(
                word_wrap=False,
                auto_fit_type="none",
                orientation="vertical_270",
                columns=2,
            ),
        )
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert 'marL="228600"' in slide_xml  # noqa: S101
    assert 'indent="-114300"' in slide_xml  # noqa: S101
    assert 'wrap="none"' in slide_xml  # noqa: S101
    assert 'vert="vert270"' in slide_xml  # noqa: S101
    assert 'numCol="2"' in slide_xml  # noqa: S101
    assert "<a:noAutofit/>" in slide_xml  # noqa: S101
    assert "hlinkClick" in slide_xml  # noqa: S101
    assert 'smCaps="1"' in slide_xml  # noqa: S101
    assert 'baseline="-25000"' in slide_xml  # noqa: S101


def test_update_shape_paragraph_fixture_rewrites_text_body(tmp_path: Path) -> None:
    out_path = tmp_path / "text_layout_update_edge_case.pptx"

    with Presentation.new("Text Layout Update") as prs:
        slide = prs.slides[0]
        shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            text="initial",
        )
        slide.update_shape(
            shape_id,
            {"paragraph": {"left_margin": 114300, "hanging_indent": 57150}},
        )
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert 'marL="114300"' in slide_xml  # noqa: S101
    assert 'indent="-57150"' in slide_xml  # noqa: S101
