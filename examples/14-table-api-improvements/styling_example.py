"""Example showing improved table styling system."""

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.tables.table_styles import TableStyle
from gopptx.schemas import Inches


def main():
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Table Styling Examples") as prs:
        # =======================================================================
        # OLD WAY - Raw GUIDs (no longer necessary)
        # =======================================================================
        prs.add_slide("OLD WAY - Raw GUIDs (Don't Do This)", layout="title_only")

        table_id = prs.add_table(
            slide=0,
            rows=2,
            cols=2,
            bounds=(Inches(1), Inches(1.5), Inches(8), Inches(2)),
            data=[["A", "B"], ["C", "D"]],
        )

        # OLD: Using raw GUID (confusing!)
        prs.set_table_style(0, table_id, "{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}")

        # =======================================================================
        # NEW WAY #1 - Using TableStyle Constants (Recommended)
        # =======================================================================
        prs.add_slide("NEW WAY #1 - TableStyle Constants", layout="title_only")

        table_id = prs.add_table(
            slide=1,
            rows=3,
            cols=3,
            bounds=(Inches(1), Inches(1.5), Inches(8), Inches(2.5)),
            data=[
                ["Product", "Q1", "Q2"],
                ["Widgets", "50", "60"],
                ["Gadgets", "30", "40"],
            ],
            first_row=True,
            band_row=True,
        )

        table = prs.slides[1].table(table_id)

        # NEW: Using named constants (clear and readable!)
        table.apply_style(TableStyle.MEDIUM_STYLE_2)

        # Other style options:
        # table.apply_style(TableStyle.LIGHT_STYLE_1)
        # table.apply_style(TableStyle.DARK_STYLE_1)
        # table.apply_style(TableStyle.MEDIUM_STYLE_1_ACCENT_1)

        # =======================================================================
        # NEW WAY #2 - Using Style Name String
        # =======================================================================
        prs.add_slide("NEW WAY #2 - Style Name Strings", layout="title_only")

        table_id = prs.add_table(
            slide=2,
            rows=2,
            cols=2,
            bounds=(Inches(1), Inches(1.5), Inches(8), Inches(2)),
            data=[["Name", "Value"], ["Item", "100"]],
        )

        # Apply by style name string
        prs.set_table_style(2, table_id, "MEDIUM_STYLE_1")

        # =======================================================================
        # NEW WAY #3 - Backward Compat: Raw GUID Still Works
        # =======================================================================
        prs.add_slide("NEW WAY #3 - Backward Compat (Raw GUID)", layout="title_only")

        table_id = prs.add_table(
            slide=3,
            rows=2,
            cols=2,
            bounds=(Inches(1), Inches(1.5), Inches(8), Inches(2)),
            data=[["X", "Y"], ["1", "2"]],
        )

        # Old way still works for backward compatibility
        table = prs.slides[3].table(table_id)
        table.apply_style("{B9AC3A68-259E-4EED-9050-4AE35E7F2B2D}")

        # =======================================================================
        # Available Styles
        # =======================================================================
        prs.add_slide("Available Built-in Styles", layout="title_only")

        # List all available style constants
        print("\nAvailable TableStyle Constants:")
        print("=" * 60)
        for name, guid in TableStyle.get_all().items():
            print(f"  TableStyle.{name} = {guid}")
        print()

        # =======================================================================
        # Discovery: Find Styles in Presentation
        # =======================================================================
        print("\nTable Styles in Presentation:")
        print("=" * 60)

        # Get all style names available in this presentation
        style_names = prs.get_all_table_style_names()
        print(f"Available styles: {style_names}")

        # Find a specific style by name
        style_guid = prs.get_table_style_by_name("Medium Style 2 - Accent 1")
        print(f"Found 'Medium Style 2 - Accent 1': {style_guid}")

        # List all available styles with details
        all_styles = prs.list_table_styles()
        for style in all_styles:
            print(f"  {style.get('name', 'Unknown')}: {style.get('guid', 'N/A')}")

        print()

        # Save
        output_path = output_dir / "table_styling_examples.pptx"
        prs.save(str(output_path))
        print(f"Presentation saved to {output_path}")


if __name__ == "__main__":
    main()
