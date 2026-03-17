"""Tests for live text object-model traversal APIs."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

from gopptx import Presentation, Run, RunHyperlink

if TYPE_CHECKING:
    from pathlib import Path


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
        assert len(para.runs) == 2
        assert para.runs[0].text == "Hello"
        assert para.runs[1].text == " World"

        para.runs[1].text = " gopptx"
        para.runs.add_run("!")
        assert para.runs[2].text == "!"

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
    assert 'algn="ctr"' in xml
    assert 'lvl="1"' in xml
    assert '<a:spcPct val="120000"/>' in xml
    assert '<a:spcBef><a:spcPts val="200"/></a:spcBef>' in xml
    assert '<a:spcAft><a:spcPts val="100"/></a:spcAft>' in xml


def test_slide_bulk_run_text_updates() -> None:
    with Presentation.new("Bulk Text") as prs:
        slide = prs.add_slide("Bulk")
        first_shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1000000),
            runs=[Run("One").to_payload()],
        )
        second_shape_id = slide.add_shape(
            "rect",
            (1000000, 2200000, 5000000, 1000000),
            runs=[Run("Two").to_payload()],
        )

        slide.update_shape_run_texts([
            {"shape_id": first_shape_id, "run_index": 0, "text": "Alpha"},
            (second_shape_id, 0, "Beta"),
        ])

        assert (
            slide.shape(first_shape_id).text_frame.paragraphs[0].runs[0].text == "Alpha"
        )
        assert (
            slide.shape(second_shape_id).text_frame.paragraphs[0].runs[0].text == "Beta"
        )


def test_presentation_bulk_run_text_updates() -> None:
    with Presentation.new("Bulk Text Deck") as prs:
        first_slide = prs.add_slide("Bulk One")
        second_slide = prs.add_slide("Bulk Two")
        first_shape_id = first_slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1000000),
            runs=[Run("One").to_payload()],
        )
        second_shape_id = second_slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1000000),
            runs=[Run("Two").to_payload()],
        )

        prs.update_deck_run_texts([
            {
                "slide_index": first_slide.index,
                "updates": [
                    {"shape_id": first_shape_id, "run_index": 0, "text": "Alpha"}
                ],
            },
            {
                "slide_index": second_slide.index,
                "updates": [
                    {"shape_id": second_shape_id, "run_index": 0, "text": "Beta"}
                ],
            },
        ])

        assert (
            first_slide.shape(first_shape_id).text_frame.paragraphs[0].runs[0].text
            == "Alpha"
        )
        assert (
            second_slide.shape(second_shape_id).text_frame.paragraphs[0].runs[0].text
            == "Beta"
        )


def test_shape_run_text_buffer_flushes_on_fresh_read_and_save(tmp_path: Path) -> None:
    out_path = tmp_path / "text_object_model_buffered_save.pptx"

    with Presentation.new("Buffered Text") as prs:
        slide = prs.add_slide("Text")
        slide_number = slide.index + 1
        shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1000000),
            runs=[Run("Hello").to_payload(), Run(" World").to_payload()],
        )
        run = slide.shape(shape_id).text_frame.paragraphs[0].runs[0]
        run.text = "Hi"

        fresh_run = slide.shape(shape_id).text_frame.paragraphs[0].runs[0]
        assert fresh_run.text == "Hi"

        prs.save(str(out_path))

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read(f"ppt/slides/slide{slide_number}.xml").decode("utf-8")
    assert "<a:t>Hi</a:t>" in slide_xml


def test_shape_run_text_buffer_keeps_slide_local_reads_coherent() -> None:
    with Presentation.new("Buffered Multi Shape") as prs:
        slide = prs.add_slide("Text")
        first_shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1000000),
            runs=[Run("One").to_payload()],
        )
        second_shape_id = slide.add_shape(
            "rect",
            (1000000, 2200000, 5000000, 1000000),
            runs=[Run("Two").to_payload()],
        )

        first_run = slide.shape(first_shape_id).text_frame.paragraphs[0].runs[0]
        second_run = slide.shape(second_shape_id).text_frame.paragraphs[0].runs[0]

        first_run.text = "Alpha"
        assert slide.shape(first_shape_id).text == "Alpha"
        assert (
            slide.shape(second_shape_id).text_frame.paragraphs[0].runs[0].text == "Two"
        )

        second_run.text = "Beta"
        assert (
            slide.shape(first_shape_id).text_frame.paragraphs[0].runs[0].text == "Alpha"
        )
        assert slide.shape(second_shape_id).text == "Beta"


def test_shape_run_hyperlink_updates_buffer_and_persist_final_state(
    tmp_path: Path,
) -> None:
    out_path = tmp_path / "text_object_model_hyperlink_buffered_save.pptx"

    with Presentation.new("Buffered Hyperlink") as prs:
        slide = prs.add_slide("Text")
        slide_number = slide.index + 1
        shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1000000),
            runs=[Run("Link").to_payload()],
        )
        run = slide.shape(shape_id).text_frame.paragraphs[0].runs[0]

        run.hyperlink = RunHyperlink(address="https://example.com/one")
        assert run.hyperlink.address == "https://example.com/one"
        run.hyperlink = RunHyperlink(address="https://example.com/two")
        assert run.hyperlink.address == "https://example.com/two"
        run.hyperlink = None
        assert run.hyperlink.address is None

        prs.save(str(out_path))

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read(f"ppt/slides/slide{slide_number}.xml").decode("utf-8")
    assert "hlinkClick" not in slide_xml
