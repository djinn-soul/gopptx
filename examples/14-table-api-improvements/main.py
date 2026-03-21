"""Example showcasing improved table API with data loading and convenience builders."""

import os
from pathlib import Path

from gopptx import Presentation
from gopptx.schemas import Inches


def main():
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    # Create presentation
    with Presentation.new("Table API Improvements") as prs:
        # =======================================================================
        # Slide 1: New named parameter style
        # =======================================================================
        prs.add_slide("New Named Parameter Style", layout="title_only")

        prs.add_table(
            slide=0,
            rows=3,
            cols=2,
            bounds=(Inches(0.8), Inches(1.5), Inches(8), Inches(2.5)),
        )

        # =======================================================================
        # Slide 2: Table with data loaded in one call
        # =======================================================================
        prs.add_slide("Table with Data", layout="title_only")

        table_id = prs.add_table(
            slide=1,
            rows=4,
            cols=3,
            bounds=(Inches(0.8), Inches(1.5), Inches(8), Inches(3)),
            data=[
                ["Product", "Q1", "Q2"],
                ["Widgets", "120", "140"],
                ["Gadgets", "90", "110"],
                ["Gizmos", "60", "70"],
            ],
            first_row=True,
            band_row=True,
        )

        # Access and modify via proxy
        table = prs.slides[1].table(table_id)
        print(f"Created table with {table.row_count} rows, {table.col_count} cols")
        print(f"Cell [0,0] = '{table[0, 0].text}'")

        # =======================================================================
        # Slide 3: Using add_table_from_rows() convenience method
        # =======================================================================
        prs.add_slide("Using add_table_from_rows()", layout="title_only")

        inventory_data = [
            ["Item", "Stock", "Status"],
            ["Widgets", "150", "Available"],
            ["Gadgets", "80", "Low Stock"],
            ["Gizmos", "200", "Available"],
            ["Doohickeys", "15", "Critical"],
        ]

        prs.add_table_from_rows(
            slide=2,
            rows=inventory_data,
            bounds=(Inches(0.8), Inches(1.5), Inches(8), Inches(3)),
            first_row=True,
            band_row=True,
        )

        # =======================================================================
        # Slide 4: Using add_table_from_dicts() convenience method
        # =======================================================================
        prs.add_slide("Using add_table_from_dicts()", layout="title_only")

        sales_data = [
            {"region": "North", "sales": "45000", "growth": "12%"},
            {"region": "South", "sales": "38000", "growth": "8%"},
            {"region": "East", "sales": "52000", "growth": "15%"},
            {"region": "West", "sales": "41000", "growth": "10%"},
        ]

        prs.add_table_from_dicts(
            slide=3,
            rows=sales_data,
            column_names=["region", "sales", "growth"],  # Explicit column order
            bounds=(Inches(0.8), Inches(1.5), Inches(8), Inches(3)),
            first_row=True,
            band_row=True,
        )

        # =======================================================================
        # Slide 5: Table with all formatting flags
        # =======================================================================
        prs.add_slide("All Formatting Options", layout="title_only")

        prs.add_table(
            slide=4,
            rows=3,
            cols=3,
            bounds=(Inches(0.8), Inches(1.5), Inches(8), Inches(2.5)),
            data=[
                ["Name", "Type", "Status"],
                ["Item A", "Category 1", "Active"],
                ["Item B", "Category 2", "Inactive"],
            ],
            first_row=True,  # Header row formatting
            first_col=True,  # First column emphasis
            last_row=True,   # Last row emphasis
            last_col=True,   # Last column emphasis
            band_row=True,   # Alternating row colors
            band_col=True,   # Alternating column colors
        )

        # =======================================================================
        # Slide 6: Bulk updating table data with set_data()
        # =======================================================================
        prs.add_slide("Bulk Update with set_data()", layout="title_only")

        table_id = prs.add_table(
            slide=5,
            rows=2,
            cols=2,
            bounds=(Inches(0.8), Inches(1.5), Inches(8), Inches(2)),
            data=[["Old A", "Old B"], ["Old C", "Old D"]],
        )

        # Later, bulk replace all data
        table = prs.slides[5].table(table_id)
        table.set_data([
            ["New A", "New B"],
            ["New C", "New D"],
        ])

        # =======================================================================
        # Slide 7: Comparison - Old vs New API
        # =======================================================================
        prs.add_slide("API Comparison", layout="title_only")

        # Old way (still works - backward compatible)
        x, y, cx, cy = Inches(0.8), Inches(1.5), Inches(4), Inches(2.5)
        old_table_id = prs.add_table(0, 2, 2, (x, y, cx, cy))
        prs.set_table_cell_text(0, old_table_id, 0, 0, "Old API")
        prs.set_table_cell_text(0, old_table_id, 0, 1, "Still Works")
        prs.set_table_cell_text(0, old_table_id, 1, 0, "A")
        prs.set_table_cell_text(0, old_table_id, 1, 1, "B")
        prs.set_table_flags(0, old_table_id, {"first_row": True})

        # New way (recommended)
        new_table_id = prs.add_table(
            slide=6,
            rows=2,
            cols=2,
            bounds=(Inches(5), Inches(1.5), Inches(4), Inches(2.5)),
            data=[["New API", "Cleaner"], ["A", "B"]],
            first_row=True,
        )

        # Save presentation
        output_path = output_dir / "table_api_improvements.pptx"
        prs.save(str(output_path))
        print(f"\nPresentation saved to {output_path}")


if __name__ == "__main__":
    main()
