"""Tests for live text object-model traversal APIs."""
# ruff: noqa: D103,PLR2004,TC003

from __future__ import annotations

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
