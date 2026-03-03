"""Notes-slide parity checks for nullable notes API semantics."""

import pathlib

import pytest
from gopptx import Presentation

project_root = (pathlib.Path(__file__).parent / "../..").resolve()
input_deck = project_root / "examples" / "assets" / "01" / "01_basic_pptx.pptx"


def test_notes_slide_is_none_until_created() -> None:
    """Notes slide is absent until notes text is written."""
    if not input_deck.exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.slides[0]
        if slide.notes_slide is not None:
            raise AssertionError("expected notes_slide to be None before note creation")

        slide.notes = "hello notes"
        notes_slide = slide.notes_slide
        if notes_slide is None:
            raise AssertionError("expected notes_slide after setting notes")
        if notes_slide.text != "hello notes":
            raise AssertionError("expected notes text round-trip to match")
