"""Coverage for issues #41, #49, #67, #68, #144, #194, #319, #339, #452 and #547.

Each test drives the facade and then asserts on the saved package, because an
in-memory getter that echoes the value it was handed passes a facade-only test
while nothing reaches the file.
"""

from __future__ import annotations

import re
import zipfile
from typing import TYPE_CHECKING

from gopptx import Presentation
from gopptx.schemas import Inches

if TYPE_CHECKING:
    from pathlib import Path

_SLIDE_PART_RE = re.compile(r"^ppt/slides/slide\d+\.xml$")


def _slide_xml(pptx_path: Path, index: int = 1) -> str:
    with zipfile.ZipFile(pptx_path) as zf:
        return zf.read(f"ppt/slides/slide{index}.xml").decode("utf-8")


def _slide_parts(pptx_path: Path) -> list[str]:
    with zipfile.ZipFile(pptx_path) as zf:
        return sorted(name for name in zf.namelist() if _SLIDE_PART_RE.match(name))


def test_issue_41_shape_and_picture_deletion(tmp_path: Path) -> None:
    """delete() and shapes.remove() drop the shape from the saved spTree."""
    output_path = tmp_path / "issue_41.pptx"

    with Presentation.new(title="Test 41 Deck") as prs:
        slide = prs.slides[0]
        keep = slide.shapes.add_textbox(
            Inches(1.0), Inches(1.0), Inches(4.0), Inches(2.0), text="SHAPE_KEEP"
        )
        doomed = slide.shapes.add_textbox(
            Inches(1.0), Inches(3.0), Inches(4.0), Inches(2.0), text="SHAPE_DELETE"
        )
        removed = slide.shapes.add_textbox(
            Inches(1.0), Inches(5.0), Inches(4.0), Inches(1.0), text="SHAPE_REMOVE"
        )
        initial_count = len(slide.shapes)

        doomed.delete()
        assert len(slide.shapes) == initial_count - 1

        slide.shapes.remove(removed)
        assert len(slide.shapes) == initial_count - 2

        assert keep.text_frame.text == "SHAPE_KEEP"
        prs.save(output_path)

    xml = _slide_xml(output_path)
    assert "SHAPE_KEEP" in xml
    assert "SHAPE_DELETE" not in xml
    assert "SHAPE_REMOVE" not in xml


def test_issue_49_shape_z_order(tmp_path: Path) -> None:
    """move_to_front/back reorder the shape within the saved spTree."""
    output_path = tmp_path / "issue_49.pptx"

    with Presentation.new(title="Test 49 Deck") as prs:
        slide = prs.slides[0]
        first = slide.shapes.add_textbox(
            Inches(1.0), Inches(1.0), Inches(3.0), Inches(1.0), text="ZBOX_A"
        )
        second = slide.shapes.add_textbox(
            Inches(1.0), Inches(2.5), Inches(3.0), Inches(1.0), text="ZBOX_B"
        )

        assert first.z_order < second.z_order

        first.move_to_front()
        assert first.z_order == len(slide.shapes) - 1
        assert first.z_order > second.z_order

        prs.save(output_path)

    order = re.findall(r"ZBOX_[AB]", _slide_xml(output_path))
    assert order == ["ZBOX_B", "ZBOX_A"]


def test_issue_49_move_to_back(tmp_path: Path) -> None:
    """move_to_back puts the shape first in document order."""
    output_path = tmp_path / "issue_49_back.pptx"

    with Presentation.new(title="Test 49 Back") as prs:
        slide = prs.slides[0]
        slide.shapes.add_textbox(
            Inches(1.0), Inches(1.0), Inches(3.0), Inches(1.0), text="ZBACK_A"
        )
        last = slide.shapes.add_textbox(
            Inches(1.0), Inches(2.5), Inches(3.0), Inches(1.0), text="ZBACK_B"
        )

        last.move_to_back()
        assert last.z_order == 0
        prs.save(output_path)

    order = re.findall(r"ZBACK_[AB]", _slide_xml(output_path))
    assert order == ["ZBACK_B", "ZBACK_A"]


def test_issue_67_slide_deletion(tmp_path: Path) -> None:
    """Deleting slides removes their parts from the saved package."""
    output_path = tmp_path / "issue_67.pptx"

    with Presentation.new(title="Test 67 Deck") as prs:
        prs.add_slide("Slide 2")
        prs.add_slide("Slide 3")
        initial_count = len(prs.slides)
        assert initial_count == 3

        prs.slides[2].delete()
        assert len(prs.slides) == initial_count - 1

        prs.slides.remove(0)
        assert len(prs.slides) == initial_count - 2
        prs.save(output_path)

    assert len(_slide_parts(output_path)) == 1


def test_issue_68_slide_reordering(tmp_path: Path) -> None:
    """Reordering survives a save/reload round trip."""
    output_path = tmp_path / "issue_68.pptx"

    with Presentation.new(title="DECK_FIRST") as prs:
        second = prs.add_slide("DECK_SECOND")
        third = prs.add_slide("DECK_THIRD")
        assert second.title == "DECK_SECOND"
        assert third.title == "DECK_THIRD"

        moved = prs.slides[0]
        moved_title = moved.title
        moved.move_to(2)
        assert [slide.title for slide in prs.slides] == [
            "DECK_SECOND",
            "DECK_THIRD",
            moved_title,
        ]

        prs.slides.move(2, 0)
        assert prs.slides[0].title == moved_title
        prs.save(output_path)

    with Presentation(str(output_path)) as reloaded:
        assert [slide.title for slide in reloaded.slides] == [
            "DECK_FIRST",
            "DECK_SECOND",
            "DECK_THIRD",
        ]


def test_issue_144_run_and_paragraph_removal(tmp_path: Path) -> None:
    """Removed runs and paragraphs are absent from the saved text body."""
    output_path = tmp_path / "issue_144.pptx"

    with Presentation.new(title="Test 144 Deck") as prs:
        slide = prs.slides[0]
        shape = slide.shapes.add_textbox(
            Inches(1.0), Inches(1.0), Inches(5.0), Inches(2.0), text=""
        )

        text_frame = shape.text_frame
        paragraph = text_frame.paragraphs[0]
        keep_run = paragraph.runs.add_run("RUN_KEEP")
        doomed_run = paragraph.runs.add_run("RUN_DELETE")
        paragraph.runs.add_run("RUN_REMOVE")
        initial_runs = len(paragraph.runs)

        doomed_run.delete()
        assert len(paragraph.runs) == initial_runs - 1

        # RUN_REMOVE is now the last run, since RUN_DELETE sat between them.
        paragraph.remove_run(len(paragraph.runs) - 1)
        assert len(paragraph.runs) == initial_runs - 2
        assert keep_run.text == "RUN_KEEP"

        doomed_paragraph = text_frame.add_paragraph("PARA_DELETE")
        paragraph_count = len(text_frame.paragraphs)
        doomed_paragraph.delete()
        assert len(text_frame.paragraphs) == paragraph_count - 1
        prs.save(output_path)

    xml = _slide_xml(output_path)
    assert "RUN_KEEP" in xml
    assert "RUN_DELETE" not in xml
    assert "RUN_REMOVE" not in xml
    assert "PARA_DELETE" not in xml


def test_issue_194_insert_slide_at_index(tmp_path: Path) -> None:
    """insert_slide and add_slide(index=) place the slide and it persists."""
    output_path = tmp_path / "issue_194.pptx"

    with Presentation.new(title="DECK_ORIGINAL") as prs:
        initial_count = len(prs.slides)

        inserted = prs.slides.insert_slide(0, title="INSERTED_FIRST")
        assert inserted.title == "INSERTED_FIRST"
        assert prs.slides[0].title == "INSERTED_FIRST"
        assert len(prs.slides) == initial_count + 1

        middle = prs.add_slide("INSERTED_MIDDLE", index=1)
        assert middle.title == "INSERTED_MIDDLE"
        assert prs.slides[1].title == "INSERTED_MIDDLE"
        assert len(prs.slides) == initial_count + 2
        prs.save(output_path)

    with Presentation(str(output_path)) as reloaded:
        assert [slide.title for slide in reloaded.slides] == [
            "INSERTED_FIRST",
            "INSERTED_MIDDLE",
            "DECK_ORIGINAL",
        ]


def test_issue_319_slide_is_hidden(tmp_path: Path) -> None:
    """is_hidden writes show="0" on p:sld and clears it again."""
    hidden_path = tmp_path / "issue_319_hidden.pptx"
    shown_path = tmp_path / "issue_319_shown.pptx"

    with Presentation.new(title="Test 319 Deck") as prs:
        slide = prs.slides[0]
        assert not slide.is_hidden

        slide.is_hidden = True
        assert slide.is_hidden
        prs.save(hidden_path)

        slide.is_hidden = False
        assert not slide.is_hidden
        prs.save(shown_path)

    assert 'show="0"' in _slide_xml(hidden_path)
    assert 'show="0"' not in _slide_xml(shown_path)


def test_issue_339_strikethrough_support(tmp_path: Path) -> None:
    """strikethrough maps to the a:rPr strike attribute both ways."""
    struck_path = tmp_path / "issue_339_struck.pptx"
    plain_path = tmp_path / "issue_339_plain.pptx"

    with Presentation.new(title="Test 339 Deck") as prs:
        slide = prs.slides[0]
        shape = slide.shapes.add_textbox(
            Inches(1.0),
            Inches(1.0),
            Inches(4.0),
            Inches(2.0),
            text="Strikethrough Test",
        )
        run = shape.text_frame.paragraphs[0].runs[0]

        run.strikethrough = True
        assert run.strikethrough is True
        assert run.font.strikethrough is True
        prs.save(struck_path)

        run.font.strikethrough = False
        assert run.strikethrough is False
        assert run.font.strikethrough is False
        prs.save(plain_path)

    assert 'strike="sngStrike"' in _slide_xml(struck_path)
    assert 'strike="sngStrike"' not in _slide_xml(plain_path)


def test_issue_452_shape_hidden(tmp_path: Path) -> None:
    """hidden writes the cNvPr hidden attribute and clears it again."""
    hidden_path = tmp_path / "issue_452_hidden.pptx"
    shown_path = tmp_path / "issue_452_shown.pptx"

    with Presentation.new(title="Test 452 Deck") as prs:
        slide = prs.slides[0]
        shape = slide.shapes.add_textbox(
            Inches(1.0), Inches(1.0), Inches(4.0), Inches(2.0), text="Hidden Shape Test"
        )
        assert not shape.hidden

        shape.hidden = True
        assert shape.hidden
        prs.save(hidden_path)

        shape.hidden = False
        assert not shape.hidden
        prs.save(shown_path)

    assert 'hidden="1"' in _slide_xml(hidden_path)
    assert 'hidden="1"' not in _slide_xml(shown_path)


def test_issue_547_shape_flip_x_y(tmp_path: Path) -> None:
    """flip_x/flip_y alias flip_horizontal/flip_vertical and write flipH/flipV."""
    flipped_path = tmp_path / "issue_547_flipped.pptx"
    reset_path = tmp_path / "issue_547_reset.pptx"

    with Presentation.new(title="Test 547 Deck") as prs:
        slide = prs.slides[0]
        shape = slide.shapes.add_textbox(
            Inches(1.0), Inches(1.0), Inches(4.0), Inches(2.0), text="Flip Test"
        )

        shape.flip_x = True
        shape.flip_y = True
        assert shape.flip_x is True
        assert shape.flip_horizontal is True
        assert shape.flip_y is True
        assert shape.flip_vertical is True
        prs.save(flipped_path)

        shape.flip_x = False
        shape.flip_y = False
        assert not shape.flip_x
        assert not shape.flip_y
        prs.save(reset_path)

    flipped_xml = _slide_xml(flipped_path)
    assert 'flipH="1"' in flipped_xml
    assert 'flipV="1"' in flipped_xml

    reset_xml = _slide_xml(reset_path)
    assert 'flipH="1"' not in reset_xml
    assert 'flipV="1"' not in reset_xml
