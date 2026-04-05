"""Demonstrate working with slide placeholders via available Python APIs.

This example demonstrates:
- add_bullet_slide() as the primary way to populate title/body placeholders
- add_title_slide() for centered-title layout
- add_paragraph_slide() for a single body text block
- Enumerating SlideLayoutType constants that correspond to placeholder types
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches


def _add_placeholder_overview(prs: Presentation) -> None:
    """Slides showing placeholder-driven content."""
    # Title + body placeholders via add_bullet_slide
    prs.add_bullet_slide(
        "Title Placeholder + Body Placeholder",
        [
            "add_bullet_slide(title, bullets) fills the title placeholder",
            "and the body/content placeholder with bullet text",
            "This is the most common slide type in PowerPoint",
        ],
    )

    # Title only - body placeholder unused
    prs.add_slide(
        "Title Only Placeholder",
        layout=SlideLayoutType.TITLE_ONLY,
    )

    # Centered title (title + subtitle placeholders)
    prs.add_title_slide("Centered Title Placeholder")

    # Body as single paragraph block
    prs.add_paragraph_slide(
        "Body as Paragraph Text",
        "The body placeholder can hold a free-form paragraph of text "
        "instead of a bulleted list. Use add_paragraph_slide() for this.",
    )

    # Blank - no placeholders active
    prs.add_slide("", layout=SlideLayoutType.BLANK)


def _add_layout_reference_slide(prs: Presentation) -> None:
    """Slide listing SlideLayoutType constants relevant to placeholders."""
    prs.add_bullet_slide(
        "SlideLayoutType Constants",
        [
            f"BLANK          = {SlideLayoutType.BLANK!r}",
            f"TITLE_ONLY     = {SlideLayoutType.TITLE_ONLY!r}",
            f"TITLE_CONTENT  = {SlideLayoutType.TITLE_AND_CONTENT!r}",
        ],
    )


def _add_table_placeholder_slide(prs: Presentation) -> None:
    """Demonstrate placing a table on a slide (placeholder-like content)."""
    idx = prs.add_slide("Table in Slide Body", layout=SlideLayoutType.TITLE_ONLY).index
    rows = [
        ["Column A", "Column B", "Column C"],
        ["Value 1", "Value 2", "Value 3"],
        ["Value 4", "Value 5", "Value 6"],
    ]
    prs.add_table_from_rows(
        idx,
        rows,
        bounds=(Inches(0.5), Inches(1.5), Inches(9), Inches(2.5)),
        first_row=True,
        band_row=True,
    )


def main() -> None:
    """Create presentation demonstrating placeholder-driven slide content."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "66_placeholders_api.pptx"

    with Presentation.new("Placeholders API Demo") as prs:
        _add_placeholder_overview(prs)
        _add_table_placeholder_slide(prs)
        _add_layout_reference_slide(prs)

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Demonstrated: add_bullet_slide, add_title_slide, add_paragraph_slide,")
    print(
        "  BLANK / TITLE_ONLY / TITLE_AND_CONTENT layouts, add_table_from_rows in body"
    )


if __name__ == "__main__":
    main()
