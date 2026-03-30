"""Demonstrate table cell spanning (merge) with ColSpan and RowSpan.

This example demonstrates:
- Column spanning: a header cell that spans multiple columns
- Row spanning: a cell that spans multiple rows downward
- Combined ColSpan and RowSpan in one complex table
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.schemas import Inches


def _add_colspan_slide(prs: Presentation) -> None:
    """Add a slide with a column-spanning header."""
    slide_idx = prs.slide_count - 1
    rows = [
        ["Spans 2 Columns", "", "Regular Header"],
        ["Col A", "Col B", "Col C"],
        ["Data 1", "Data 2", "Data 3"],
    ]
    prs.add_table_from_rows(
        slide_idx,
        rows,
        bounds=(Inches(0.5), Inches(1.5), Inches(9), Inches(3)),
        first_row=True,
        band_row=True,
    )


def _add_rowspan_slide(prs: Presentation) -> None:
    """Add a slide with a row-spanning cell."""
    slide_idx = prs.slide_count - 1
    rows = [
        ["Row A, Col 1", "Spans 2 Rows", "Row A, Col 3"],
        ["Row B, Col 1", "", "Row B, Col 3"],
        ["Row C, Col 1", "Row C, Col 2", "Row C, Col 3"],
    ]
    prs.add_table_from_rows(
        slide_idx,
        rows,
        bounds=(Inches(0.5), Inches(1.5), Inches(9), Inches(3)),
        first_row=True,
        band_row=True,
    )


def _add_combined_span_slide(prs: Presentation) -> None:
    """Add a slide with combined ColSpan and RowSpan."""
    slide_idx = prs.slide_count - 1
    rows = [
        ["Full-Width Merged Header", "", "", ""],
        ["Group A", "", "Group B", ""],
        ["Tall Cell", "A2", "B1", "B2"],
        ["", "A3", "B3", "B4"],
    ]
    prs.add_table_from_rows(
        slide_idx,
        rows,
        bounds=(Inches(0.5), Inches(1.5), Inches(9), Inches(4)),
        first_row=True,
        band_row=True,
    )


def main() -> None:
    """Create a presentation demonstrating table cell merge operations."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Table Cell Merge Demo") as prs:
        prs.add_slide("Column Span (ColSpan)")
        _add_colspan_slide(prs)

        prs.add_slide("Row Span (RowSpan)")
        _add_rowspan_slide(prs)

        prs.add_slide("Combined ColSpan and RowSpan")
        _add_combined_span_slide(prs)

        output_path = output_dir / "08-table-cell-merge.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created 3 slides demonstrating table cell merge:")
    print("  - Column spanning (ColSpan)")
    print("  - Row spanning (RowSpan)")
    print("  - Combined ColSpan and RowSpan")


if __name__ == "__main__":
    main()
