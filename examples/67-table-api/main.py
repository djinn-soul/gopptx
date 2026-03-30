"""Demonstrate table creation and formatting.

This example demonstrates:
- add_table_from_rows() with plain text rows
- first_row banding (header row styling)
- band_row alternating row colors
- Tables on TITLE_ONLY slides with custom bounds
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches


def _add_plain_table(prs: Presentation) -> None:
    """Plain text table: name / role / location."""
    idx = prs.add_slide("Plain Text Table", layout=SlideLayoutType.TITLE_ONLY)
    rows = [
        ["Name", "Role", "Location"],
        ["Alice", "Engineer", "New York"],
        ["Bob", "Designer", "London"],
        ["Carol", "Manager", "Tokyo"],
        ["David", "Analyst", "Sydney"],
    ]
    prs.add_table_from_rows(
        idx,
        rows,
        bounds=(Inches(0.5), Inches(1.5), Inches(9), Inches(3.0)),
        first_row=True,
        band_row=True,
    )


def _add_product_table(prs: Presentation) -> None:
    """Product / units / price table with header row."""
    idx = prs.add_slide(
        "Product Table - Header + Banded Rows",
        layout=SlideLayoutType.TITLE_ONLY,
    )
    rows = [
        ["Product", "Units", "Price"],
        ["Widget A", "150", "$12.50"],
        ["Widget B", "200", "$9.99"],
        ["Widget C", "75", "$24.00"],
        ["Widget D", "310", "$5.49"],
    ]
    prs.add_table_from_rows(
        idx,
        rows,
        bounds=(Inches(0.5), Inches(1.5), Inches(9), Inches(3.2)),
        first_row=True,
        band_row=True,
    )


def _add_status_table(prs: Presentation) -> None:
    """Status table without banding to show difference."""
    idx = prs.add_slide(
        "Status Table - No Row Banding",
        layout=SlideLayoutType.TITLE_ONLY,
    )
    rows = [
        ["Name", "Notes", "Status"],
        ["Alpha", "First entry in the system", "Active"],
        ["Beta", "Second entry, deprecated soon", "Pending"],
        ["Gamma", "Archived", "Inactive"],
    ]
    prs.add_table_from_rows(
        idx,
        rows,
        bounds=(Inches(0.5), Inches(1.5), Inches(9), Inches(2.8)),
        first_row=True,
        band_row=False,
    )


def _add_wide_table(prs: Presentation) -> None:
    """Wide table spanning full slide width."""
    idx = prs.add_slide(
        "Wide Table - Full Width",
        layout=SlideLayoutType.TITLE_ONLY,
    )
    rows = [
        ["Quarter", "Revenue ($k)", "Costs ($k)", "Profit ($k)", "Growth (%)"],
        ["Q1", "42", "18", "24", "+5%"],
        ["Q2", "55", "24", "31", "+12%"],
        ["Q3", "61", "30", "31", "+8%"],
        ["Q4", "73", "41", "32", "+15%"],
    ]
    prs.add_table_from_rows(
        idx,
        rows,
        bounds=(Inches(0.3), Inches(1.5), Inches(9.4), Inches(3.2)),
        first_row=True,
        band_row=True,
    )


def main() -> None:
    """Create presentation demonstrating table creation via add_table_from_rows."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "67_table_api.pptx"

    with Presentation.new("Table API Demo") as prs:
        _add_plain_table(prs)
        _add_product_table(prs)
        _add_status_table(prs)
        _add_wide_table(prs)

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print(
        "Demonstrated: add_table_from_rows with first_row=True/False, band_row=True/False"
    )
    print("  and varying bounds using Inches()")


if __name__ == "__main__":
    main()
