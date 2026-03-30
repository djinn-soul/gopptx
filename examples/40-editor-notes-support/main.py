"""Demonstrate reading and writing speaker notes on slides.

This example demonstrates:
- Setting speaker notes at slide creation time
- Overwriting existing notes with set_notes()
- Adding new slides with notes via the editor
- Reading notes back with get_notes()
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def main() -> None:
    """Create presentation demonstrating editor notes support."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Notes Demo") as prs:
        # Slide 1: initial notes (will be overwritten)
        prs.add_bullet_slide(
            "Slide 1",
            ["Point A", "Point B"],
        )
        prs.set_notes(
            0, "Speaker notes for slide 1.\nThese will be overwritten by the editor."
        )

        # Slide 2: notes that remain unchanged
        prs.add_bullet_slide(
            "Slide 2",
            ["Point C", "Point D"],
        )
        prs.set_notes(1, "Speaker notes for slide 2.\nThese remain unchanged.")

        # Overwrite slide 1 notes
        prs.set_notes(0, "Updated notes for slide 1 — written via set_notes().")

        # Slide 3: added programmatically with notes
        prs.add_bullet_slide(
            "Slide 3 (Added via Editor)",
            [
                "Added programmatically",
                "Includes notes written at add time",
            ],
        )
        prs.set_notes(2, "Notes for the new slide — set during add_slide.")

        # Verify notes are readable
        notes_0 = prs.get_notes(0)
        notes_1 = prs.get_notes(1)
        notes_2 = prs.get_notes(2)

        output_path = output_dir / "40_editor_notes_support.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print(f"Slide 1 notes: {notes_0!r}")
    print(f"Slide 2 notes: {notes_1!r}")
    print(f"Slide 3 notes: {notes_2!r}")


if __name__ == "__main__":
    main()
