"""Test notes extraction across presentations with >20 slides (Issue #1099)."""

import pathlib

from gopptx import Presentation


def test_notes_extraction_beyond_20_slides(tmp_path: pathlib.Path) -> None:
    """Slide notes on slide 21+ are preserved and correctly extracted (Issue #1099)."""
    output_path = tmp_path / "large_deck_notes.pptx"

    with Presentation.new(title="Large Deck Notes Test") as pres:
        # Add 25 slides
        for i in range(1, 26):
            slide = pres.add_slide(title=f"Slide {i}") if i > 1 else pres.slides[0]
            slide.notes = f"Notes for slide {i}"

        pres.save(output_path)

    # Reopen deck and extract notes from all 25 slides
    with Presentation(output_path) as pres:
        assert pres.slide_count == 25
        for i in range(1, 26):
            notes_text = pres.slides[i - 1].notes
            assert notes_text == f"Notes for slide {i}", (
                f"Slide {i} notes mismatch: got {notes_text!r}"
            )
