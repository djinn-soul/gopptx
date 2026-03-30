"""Demonstrate the slides API with multiple layout types.

This example demonstrates:
- Adding slides with different SlideLayoutType values (BLANK, TITLE_ONLY, TITLE_CONTENT)
- add_title_slide(), add_bullet_slide(), add_paragraph_slide()
- Slide manipulation: duplicate_slide(), move_slide(), remove_slide()
- Reading slide count and accessing slides by index
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType


def _add_layout_slides(prs: Presentation) -> None:
    """Add slides demonstrating each supported layout type."""
    # Title slide (centered title layout)
    prs.add_title_slide("Slides API Demo")

    # Title + content (default bullet layout)
    prs.add_bullet_slide(
        "Slide Layout: Title and Content",
        [
            "This is the default content layout",
            "Use add_bullet_slide() for bullet lists",
            "SlideLayoutType.TITLE_CONTENT",
        ],
    )

    # Title only layout
    prs.add_slide("Title Only Layout", layout=SlideLayoutType.TITLE_ONLY)

    # Blank layout
    prs.add_slide("", layout=SlideLayoutType.BLANK)

    # Paragraph slide
    prs.add_paragraph_slide(
        "Paragraph Slide",
        "This slide uses add_paragraph_slide() to add a block of "
        "free-form text instead of a bullet list.",
    )


def _demonstrate_slide_manipulation(prs: Presentation) -> None:
    """Show duplicate, move, and remove operations."""
    count_before = prs.slide_count
    print(f"Slide count before manipulation : {count_before}")

    # Duplicate the first content slide (index 1)
    prs.duplicate_slide(1)
    print(f"After duplicate_slide(1)        : {prs.slide_count} slides")

    # Move the duplicated slide to the end
    prs.move_slide(prs.slide_count - 1, prs.slide_count - 1)
    print(f"After move_slide                : {prs.slide_count} slides")

    # Remove the last slide
    prs.remove_slide(prs.slide_count - 1)
    print(f"After remove_slide              : {prs.slide_count} slides")


def main() -> None:
    """Create presentation demonstrating the slides API."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "64_slides_api.pptx"

    with Presentation.new("Slides API Demo") as prs:
        _add_layout_slides(prs)
        _demonstrate_slide_manipulation(prs)

        print(f"Final slide count: {prs.slide_count}")
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Demonstrated: add_title_slide, add_bullet_slide, add_paragraph_slide,")
    print("  add_slide with layout types, duplicate_slide, move_slide, remove_slide")


if __name__ == "__main__":
    main()
