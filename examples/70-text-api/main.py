"""Demonstrate text formatting with bullet styles and slide notes.

This example demonstrates:
- add_bullet_slide() for standard bullet lists
- add_paragraph_slide() for free-form text blocks
- set_notes() for plain-text speaker notes
- get_notes() to verify notes survive a save/open round-trip
- Slide title and content via add_slide() with TITLE_CONTENT layout
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def _add_title_styling_slide(prs: Presentation) -> None:
    """Slide demonstrating title and content text."""
    prs.add_bullet_slide(
        "Slide Title Styling",
        [
            "Titles are set as the first argument to add_bullet_slide()",
            "Bullet text appears in the content placeholder",
            "Use add_paragraph_slide() for free-form body text",
        ],
    )


def _add_bullet_variants_slide(prs: Presentation) -> None:
    """Slide showing different bullet content patterns."""
    prs.add_bullet_slide(
        "Bullet Content Variants",
        [
            "Standard bullet point",
            "Another bullet at the same level",
            "Bullets can reference any content",
            "Keep bullets concise and parallel",
        ],
    )


def _add_paragraph_text_slide(prs: Presentation) -> None:
    """Slide using paragraph text instead of bullets."""
    prs.add_paragraph_slide(
        "Free-Form Paragraph Text",
        "This slide uses add_paragraph_slide() to render a free-form block "
        "of text in the body placeholder. This is useful for introductory "
        "context, longer explanations, or any content that does not benefit "
        "from bullet formatting.",
    )


def _add_notes_slides(prs: Presentation) -> int:
    """Add a slide with speaker notes and return its index."""
    idx = prs.add_bullet_slide(
        "Speaker Notes Demo",
        [
            "This slide has speaker notes attached.",
            "Open the notes panel in PowerPoint to read them.",
            "Notes are set with prs.set_notes(slide_idx, text).",
        ],
    ).index
    prs.set_notes(
        idx,
        "These are plain-text speaker notes for the demo slide.\n"
        "Use them to remind yourself of talking points.\n"
        "Notes survive a save/open round-trip.",
    )
    return idx


def _add_rich_notes_slide(prs: Presentation) -> int:
    """Add a slide with multi-line speaker notes."""
    idx = prs.add_bullet_slide(
        "Multi-Line Notes",
        [
            "This slide has multi-paragraph notes.",
            "First talking point: introduce the topic.",
            "Second talking point: key concept overview.",
        ],
    ).index
    prs.set_notes(
        idx,
        "Opening paragraph - introduce the topic.\n"
        "Key concept #1: explain the first idea.\n"
        "Key concept #2: explain the second idea.\n"
        "Closing: summarise and transition to the next slide.",
    )
    return idx


def main() -> None:
    """Create presentation demonstrating text and notes APIs."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "70_text_api.pptx"

    notes_slide_idx: int

    with Presentation.new("Text API Demo") as prs:
        _add_title_styling_slide(prs)
        _add_bullet_variants_slide(prs)
        _add_paragraph_text_slide(prs)
        notes_slide_idx = _add_notes_slides(prs)
        _add_rich_notes_slide(prs)

        # Verify notes before save
        notes = prs.get_notes(notes_slide_idx)
        print(f"Notes on slide {notes_slide_idx}: {notes[:60]!r}...")

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    # Verify notes survive round-trip
    with Presentation(str(output_path)) as prs2:
        notes = prs2.get_notes(notes_slide_idx)
        print(f"Round-trip notes verify: {notes[:40]!r}...")

    print("\n=== SUMMARY ===")
    print("Demonstrated: add_bullet_slide, add_paragraph_slide,")
    print("  set_notes, get_notes, round-trip verification")


if __name__ == "__main__":
    main()
