import pathlib

import pytest
from gopptx import SHAPE_RECTANGLE, Presentation

project_root = (pathlib.Path(__file__).parent / "../..").resolve()
input_deck = project_root / "examples/assets/01/01_basic_pptx.pptx"


def test_presentation_content_basic() -> None:
    if not input_deck.exists():
        pytest.skip("smoke sample missing")

    with Presentation(str(input_deck)) as prs:
        slide = prs.add_slide("Content Test")

        # Add shape with text and properties
        shape_id = slide.add_shape(
            SHAPE_RECTANGLE,
            bounds=(100, 100, 200, 200),
            text="Hello",
            properties={"fill_color": "FF0000"},
        )
        assert shape_id > 0

        # List shapes
        shapes = slide.list_shapes()
        assert len(shapes) > 0

        # Update shape
        slide.update_shape(shape_id, {"text": "Updated"})

        # Move shape
        slide.move_shape_to_front(shape_id)
        slide.move_shape_to_back(shape_id)

        # Remove shape
        slide.remove_shape(shape_id)

        # Find and replace
        slide.title = "Find Me"
        count = prs.find_and_replace("Find Me", "Found Me")
        assert count >= 1

        # Search shapes
        results = prs.search_shapes("Found Me")
        assert len(results) > 0


def test_comments_and_authors() -> None:
    if not input_deck.exists():
        pytest.skip("smoke sample missing")

    with Presentation(str(input_deck)) as prs:
        # Authors
        author_id = prs.add_author("Test Author", "TA")
        assert author_id >= 0
        authors = prs.get_authors()
        assert any(a["Name"] == "Test Author" for a in authors)

        # Comments
        slide_idx = 0
        comment_index = prs.add_comment(
            slide_idx, author_id, "This is a comment", x=100, y=100
        )
        comments = prs.get_comments(slide_idx)
        assert len(comments) > 0

        # Remove comment
        prs.remove_comment(comment_index)
        comments_after = prs.get_comments(slide_idx)
        assert len(comments_after) < len(comments)


def test_charts_basic() -> None:
    if not input_deck.exists():
        pytest.skip("smoke sample missing")

    with Presentation(str(input_deck)) as prs:
        slide = prs.add_slide("Chart Test")
        # Basic add_chart
        slide.add_chart("bar", ["A", "B"], [10, 20], title="Test Chart")

        charts = slide.list_charts()
        assert len(charts) > 0

        # Update chart data - now succeeds for charts without formula nodes
        prs.update_chart_data(
            slide.index, ["A", "B"], [{"name": "S1", "values": [15, 25]}]
        )


def test_notes() -> None:
    if not input_deck.exists():
        pytest.skip("smoke sample missing")

    with Presentation(str(input_deck)) as prs:
        slide = prs.slides[0]
        slide.notes = "New Notes"
        assert slide.notes == "New Notes"

        prs.set_notes(0, "Updated Notes")
        assert prs.get_notes(0) == "Updated Notes"
