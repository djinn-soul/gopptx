"""Regressions for the review findings on PR #67."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

from gopptx import MSO_LINE, Presentation
from gopptx.constants import LineDashStyle

if TYPE_CHECKING:
    from pathlib import Path


def test_str_enum_members_carry_their_value() -> None:
    """The 3.10 StrEnum shim must stringify to the value, as 3.11+ does."""
    assert str(LineDashStyle.ROUND_DOT) == "sysDot"
    assert f"{MSO_LINE.ROUND_DOT}" == "sysDot"


def test_vertical_anchor_round_trips(tmp_path: Path) -> None:
    """The setter must send the wire key the Go TextFrame is tagged with."""
    output = tmp_path / "anchor.pptx"
    with Presentation.new("Anchor") as presentation:
        slide = presentation.slides[0]
        shape = slide.shapes.add_textbox(500000, 500000, 3000000, 1000000, text="hi")
        shape.text_frame.vertical_anchor = "middle"
        assert shape.text_frame.vertical_anchor == "middle"
        presentation.save(output)

    with zipfile.ZipFile(output) as package:
        slide_xml = package.read("ppt/slides/slide1.xml").decode()
    assert 'anchor="ctr"' in slide_xml


def test_placeholder_text_can_be_cleared() -> None:
    """``placeholder.text = ""`` is how python-pptx empties a placeholder."""
    with Presentation.new("Clear") as presentation:
        placeholder = presentation.slides[0].placeholders[0]
        placeholder.text = "populated"
        assert placeholder.text == "populated"
        placeholder.text = ""
        assert not placeholder.text


def test_add_slide_accepts_a_layout_object() -> None:
    """python-pptx style ``add_slide(prs.slide_layouts[i])`` binds that layout."""
    with Presentation.new("Layouts") as presentation:
        layout = presentation.slide_layouts[1]
        slide = presentation.slides.add_slide(layout, title="From layout object")
        assert slide.slide_layout.part_name == layout.part_name


def test_deleting_a_run_reindexes_retained_proxies() -> None:
    """Proxies for later runs must follow their runs after an earlier delete."""
    with Presentation.new("Runs") as presentation:
        shape = presentation.slides[0].shapes.add_textbox(
            500000, 500000, 4000000, 1000000, text="A"
        )
        paragraph = shape.text_frame.paragraphs[0]
        for value in ("B", "C"):
            paragraph.runs.add_run(value)

        first, second, third = (paragraph.runs[i] for i in range(3))
        assert [first.text, second.text, third.text] == ["A", "B", "C"]

        second.remove()
        assert third.text == "C"
        third.text = "C2"
        assert [run.text for run in paragraph.runs] == ["A", "C2"]
