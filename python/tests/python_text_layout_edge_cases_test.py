"""Advanced text-layout edge case fixtures for parity coverage."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from pathlib import Path

from gopptx import ParagraphProps, Presentation, Run, TextFrameProps


def test_complex_text_layout_fixture_emits_expected_tokens(tmp_path: Path) -> None:
    """Test complex text layout fixture emits expected tokens."""
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
            paragraph=ParagraphProps(
                indent=228600,
                hanging=114300,
                tab_stops=[457200, 914400],
            ),
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

    assert 'marL="228600"' in slide_xml
    assert 'indent="-114300"' in slide_xml
    if (
        '<a:tabLst><a:tab pos="457200"/><a:tab pos="914400"/></a:tabLst>'
        not in slide_xml
    ):
        raise AssertionError("expected serialized tab list in slide XML")
    assert 'wrap="none"' in slide_xml
    assert 'vert="vert270"' in slide_xml
    assert 'numCol="2"' in slide_xml
    assert "<a:noAutofit/>" in slide_xml
    assert "hlinkClick" in slide_xml
    assert 'cap="small"' in slide_xml
    assert 'baseline="-25000"' in slide_xml


def test_update_shape_paragraph_fixture_rewrites_text_body(tmp_path: Path) -> None:
    """Test update shape paragraph fixture rewrites text body."""
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
            {
                "paragraph": {
                    "left_margin": 114300,
                    "hanging_indent": 57150,
                    "tabs": [228600],
                }
            },
        )
        prs.save(out_path)

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")

    assert 'marL="114300"' in slide_xml
    assert 'indent="-57150"' in slide_xml
    assert '<a:tabLst><a:tab pos="228600"/></a:tabLst>' in slide_xml


def test_paragraph_props_rejects_negative_tab_stop() -> None:
    """Negative paragraph tab-stop positions should fail fast."""
    props = ParagraphProps(tab_stops=[-1])
    try:
        props.to_payload()
    except ValueError as err:
        assert "tab_stops" in str(err)
        return
    raise AssertionError("expected ValueError for negative tab stop")
