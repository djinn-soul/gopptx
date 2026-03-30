"""Demonstrate adding and retrieving speaker notes on slides.

This example demonstrates:
- Adding speaker notes to individual slides with set_notes()
- Retrieving notes text with get_notes()
- Updating notes after initial creation
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation

NOTES_PREVIEW_LIMIT = 60


def _add_slides_with_notes(prs: Presentation) -> None:
    """Add slides with speaker notes attached."""
    prs.add_bullet_slide(
        "Slide with Notes",
        ["Bullet 1", "Bullet 2"],
    )
    prs.set_notes(
        prs.slide_count - 1,
        "This is a speaker note for slide 1.\nIt has multiple lines.",
    )

    prs.add_bullet_slide(
        "New Slide with Notes",
        ["New Bullet 1"],
    )
    prs.set_notes(
        prs.slide_count - 1,
        "Secret speaker notes for slide 2.",
    )


def _update_slide_notes(prs: Presentation) -> None:
    """Update the notes on the first slide."""
    prs.set_notes(0, "Updated notes content for slide 1.")


def _verify_notes(prs: Presentation) -> None:
    """Read and print notes to verify round-trip correctness."""
    for i in range(prs.slide_count):
        notes = prs.get_notes(i)
        if notes:
            print(
                f"  Slide {i + 1} notes: {notes[:NOTES_PREVIEW_LIMIT]}{'...' if len(notes) > NOTES_PREVIEW_LIMIT else ''}"
            )


def main() -> None:
    """Create a presentation demonstrating speaker notes."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Notes Template") as prs:
        _add_slides_with_notes(prs)
        _update_slide_notes(prs)

        print("Notes attached to slides:")
        _verify_notes(prs)

        output_path = output_dir / "22-speaker-notes.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created a presentation with speaker notes:")
    print("  - Notes added with set_notes(slide_idx, text)")
    print("  - Notes retrieved with get_notes(slide_idx)")
    print("  - Notes updated after initial creation")


if __name__ == "__main__":
    main()
