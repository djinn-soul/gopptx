"""Tests for live text object-model traversal APIs."""
# ruff: noqa: D103,PLR2004,TC003

from __future__ import annotations

import zipfile
from pathlib import Path

from gopptx import Presentation, Run


def test_shape_text_traversal_and_run_updates(tmp_path: Path) -> None:
    out_path = tmp_path / "text_object_model.pptx"

    with Presentation.new("Text Model") as prs:
        slide = prs.add_slide("Text")
        shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1000000),
            runs=[Run("Hello").to_payload(), Run(" World").to_payload()],
        )

        shape = slide.shape(shape_id)
        para = shape.text_frame.paragraphs[0]
        assert len(para.runs) == 2  # noqa: S101
        assert para.runs[0].text == "Hello"  # noqa: S101
        assert para.runs[1].text == " World"  # noqa: S101

        para.runs[1].text = " gopptx"
        para.runs.add_run("!")
        assert para.runs[2].text == "!"  # noqa: S101

        prs.save(str(out_path))


def test_shape_paragraph_advanced_controls(tmp_path: Path) -> None:
    out_path = tmp_path / "text_object_model_paragraph_advanced.pptx"

    with Presentation.new("Text Model Paragraph") as prs:
        slide = prs.add_slide("Text")
        slide_number = slide.index + 1
        shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1000000),
            runs=[Run("Para").to_payload()],
        )
        shape = slide.shape(shape_id)
        para = shape.text_frame.paragraphs[0]
        para.alignment = "center"
        para.level = 1
        para.line_spacing = 1.2
        para.space_before = 200
        para.space_after = 100
        shape.text_frame.fit_text()
        prs.save(str(out_path))

    with zipfile.ZipFile(out_path) as zf:
        slide_path = f"ppt/slides/slide{slide_number}.xml"
        xml = zf.read(slide_path).decode("utf-8")
    assert 'algn="ctr"' in xml  # noqa: S101
    assert 'lvl="1"' in xml  # noqa: S101
    assert '<a:spcPct val="120000"/>' in xml  # noqa: S101
    assert '<a:spcBef><a:spcPts val="200"/></a:spcBef>' in xml  # noqa: S101
    assert '<a:spcAft><a:spcPts val="100"/></a:spcAft>' in xml  # noqa: S101
