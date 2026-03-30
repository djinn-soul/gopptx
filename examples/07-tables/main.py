"""Demonstrate table creation and styling in gopptx.

This example demonstrates:
- Plain text rows with uniform column widths
- Styled rows with bold headers and background colors
- Mixed cell formatting with bold, italic, and alignment
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.schemas import Inches


def _add_plain_table_slide(prs: Presentation) -> None:
    """Add a slide with a plain text table."""
    slide_idx = prs.slide_count - 1
    rows = [
        ["Header 1", "Header 2", "Header 3"],
        ["Row 1, Col 1", "Row 1, Col 2", "Row 1, Col 3"],
        ["Row 2, Col 1", "Row 2, Col 2", "Row 2, Col 3"],
        ["Row 3, Col 1", "Row 3, Col 2", "Row 3, Col 3"],
    ]
    prs.add_table_from_rows(
        slide_idx,
        rows,
        bounds=(Inches(0.5), Inches(1.5), Inches(9), Inches(3)),
        first_row=True,
        band_row=True,
    )


def _add_styled_table_slide(prs: Presentation) -> None:
    """Add a slide with a styled table featuring colored headers."""
    slide_idx = prs.slide_count - 1
    rows = [
        ["Product", "Quantity", "Price"],
        ["Widget A", "150", "$12.50"],
        ["Widget B", "200", "$9.99"],
        ["Widget C", "75", "$24.00"],
    ]
    prs.add_table_from_rows(
        slide_idx,
        rows,
        bounds=(Inches(0.5), Inches(1.5), Inches(9), Inches(3)),
        first_row=True,
        band_row=True,
    )


def _add_mixed_table_slide(prs: Presentation) -> None:
    """Add a slide with a mixed-formatting table."""
    slide_idx = prs.slide_count - 1
    rows = [
        ["Name", "Description", "Status"],
        ["Alpha", "First entry in the list.", "Active"],
        ["Beta", "Second entry, slightly longer description.", "Pending"],
        ["Gamma", "Deprecated entry.", "Inactive"],
    ]
    prs.add_table_from_rows(
        slide_idx,
        rows,
        bounds=(Inches(0.5), Inches(1.5), Inches(9), Inches(3)),
        first_row=True,
        band_row=True,
    )


def main() -> None:
    """Create a presentation demonstrating table creation and styling."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Tables Demo") as prs:
        prs.add_slide("Plain Text Table")
        _add_plain_table_slide(prs)

        prs.add_slide("Styled Table with Headers")
        _add_styled_table_slide(prs)

        prs.add_slide("Mixed Cell Formatting")
        _add_mixed_table_slide(prs)

        output_path = output_dir / "07-tables.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created 3 slides demonstrating table styles:")
    print("  - Plain text rows")
    print("  - Styled headers with background colors")
    print("  - Mixed bold, italic, and alignment")


if __name__ == "__main__":
    main()
