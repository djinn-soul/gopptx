"""Demonstrate speaker notes on slides.

This example demonstrates:
- set_notes(slide_idx, text) to attach plain-text speaker notes
- get_notes(slide_idx) to read back notes
- Multi-paragraph notes via newline-separated strings
- Updating notes after initial save (overwrite notes via set_notes)
- Notes surviving a save/open round-trip
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def _add_plain_notes_slide(prs: Presentation) -> int:
    """Slide with plain single-paragraph notes."""
    idx = prs.add_bullet_slide(
        "Plain Text Notes",
        [
            "This slide has plain text speaker notes.",
            "Open the notes panel in PowerPoint to read them.",
            "Notes are set with prs.set_notes(slide_idx, text).",
        ],
    ).index
    prs.set_notes(
        idx,
        "These are plain-text speaker notes for slide 1.\n"
        "Use them to remind yourself of key talking points.",
    )
    return idx


def _add_multi_paragraph_notes_slide(prs: Presentation) -> int:
    """Slide with multi-paragraph notes."""
    idx = prs.add_bullet_slide(
        "Multi-Paragraph Notes",
        [
            "This slide uses multi-line notes.",
            "Each line separated by \\n becomes a new paragraph.",
        ],
    ).index
    prs.set_notes(
        idx,
        "Opening paragraph - introduce the topic.\n"
        "Key concept #1: explain the first idea clearly.\n"
        "Key concept #2: cover the supporting evidence.\n"
        "Closing: transition to the next slide.",
    )
    return idx


def _add_bullet_style_notes_slide(prs: Presentation) -> int:
    """Slide demonstrating notes with bullet-like content."""
    idx = prs.add_bullet_slide(
        "Structured Notes Content",
        [
            "Notes can contain bullet-style content via newlines.",
            "Each line is a separate paragraph in the notes pane.",
        ],
    ).index
    prs.set_notes(
        idx,
        "First bullet note item.\n"
        "Second bullet note item.\n"
        "First numbered note item.\n"
        "Second numbered note item.",
    )
    return idx


def _add_empty_notes_slide(prs: Presentation) -> int:
    """Slide with no notes - notes are empty by default."""
    return prs.add_bullet_slide(
        "No Notes (Default)",
        [
            "This slide has no speaker notes attached.",
            "get_notes() returns an empty string for this slide.",
        ],
    ).index


def _add_updated_notes_slide(prs: Presentation) -> int:
    """Slide whose notes will be overwritten after creation."""
    idx = prs.add_bullet_slide(
        "Notes Overwrite Demo",
        ["This slide's notes are set, then overwritten."],
    ).index
    prs.set_notes(idx, "Original notes - will be replaced.")
    return idx


def main() -> None:
    """Create presentation demonstrating the speaker notes API."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "76_notes_api.pptx"

    plain_idx: int
    overwrite_idx: int
    empty_idx: int

    with Presentation.new("Notes API Demo") as prs:
        plain_idx = _add_plain_notes_slide(prs)
        _add_multi_paragraph_notes_slide(prs)
        _add_bullet_style_notes_slide(prs)
        empty_idx = _add_empty_notes_slide(prs)
        overwrite_idx = _add_updated_notes_slide(prs)

        # Overwrite notes before saving
        prs.set_notes(
            overwrite_idx,
            "Overwritten notes - updated programmatically before save.",
        )
        notes_check = prs.get_notes(plain_idx)
        print(f"Notes on slide {plain_idx}: {notes_check[:50]!r}...")

        empty_check = prs.get_notes(empty_idx)
        print(f"Notes on empty slide: {empty_check!r}")

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    # Verify notes survive round-trip
    with Presentation(str(output_path)) as prs2:
        rt_notes = prs2.get_notes(plain_idx)
        print(f"Round-trip verify: {rt_notes[:45]!r}...")
        overwrite_rt = prs2.get_notes(overwrite_idx)
        print(f"Overwrite verify : {overwrite_rt[:45]!r}...")

    print("\n=== SUMMARY ===")
    print("Demonstrated: set_notes, get_notes, multi-paragraph notes,")
    print("  notes overwrite, empty-slide notes, and round-trip verification")


if __name__ == "__main__":
    main()
