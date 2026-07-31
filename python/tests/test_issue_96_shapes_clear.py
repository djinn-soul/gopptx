"""Issue #96: remove every shape from a slide in one call.

Asserts on the saved slide part, not just the collection length, so a facade
that forgets to flush the removal cannot pass.
"""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

import pytest
from gopptx import Presentation
from gopptx.schemas import Inches

if TYPE_CHECKING:
    from pathlib import Path


def _slide_xml(pptx_path: Path) -> str:
    with zipfile.ZipFile(pptx_path) as zf:
        return zf.read("ppt/slides/slide1.xml").decode("utf-8")


def test_clear_removes_every_shape(tmp_path: Path) -> None:
    """clear() empties the slide and the shapes are gone from the saved spTree."""
    output_path = tmp_path / "cleared.pptx"

    with Presentation.new(title="Issue 96") as prs:
        slide = prs.slides[0]
        slide.shapes.add_textbox(
            Inches(1), Inches(1), Inches(4), Inches(1), text="ALPHA"
        )
        slide.shapes.add_textbox(
            Inches(1), Inches(3), Inches(4), Inches(1), text="BETA"
        )
        initial = len(slide.shapes)
        assert initial >= 2

        removed = slide.shapes.clear()

        assert removed == initial
        assert len(slide.shapes) == 0
        prs.save(output_path)

    xml = _slide_xml(output_path)
    assert "ALPHA" not in xml
    assert "BETA" not in xml
    assert "<p:sp>" not in xml


def test_clear_on_empty_collection_is_a_no_op(tmp_path: Path) -> None:
    """Calling clear() twice is safe and reports nothing removed the second time."""
    output_path = tmp_path / "cleared_twice.pptx"

    with Presentation.new(title="Issue 96 empty") as prs:
        slide = prs.slides[0]
        slide.shapes.clear()

        assert slide.shapes.clear() == 0
        assert len(slide.shapes) == 0
        prs.save(output_path)

    assert "<p:sp>" not in _slide_xml(output_path)


def test_clear_leaves_other_slides_untouched(tmp_path: Path) -> None:
    """clear() is scoped to its own slide."""
    output_path = tmp_path / "cleared_scoped.pptx"

    with Presentation.new(title="Issue 96 scope") as prs:
        prs.add_slide("KEEP_SLIDE")
        prs.slides[0].shapes.clear()

        assert len(prs.slides[0].shapes) == 0
        assert len(prs.slides[1].shapes) > 0
        prs.save(output_path)

    with zipfile.ZipFile(output_path) as zf:
        second = zf.read("ppt/slides/slide2.xml").decode("utf-8")
    assert "KEEP_SLIDE" in second


def test_remove_rejects_a_non_shape() -> None:
    """remove() still refuses anything that is not a proxy or an integer id."""
    with (
        Presentation.new(title="Issue 96 guard") as prs,
        pytest.raises(TypeError),
    ):
        prs.slides[0].shapes.remove("not-a-shape")
