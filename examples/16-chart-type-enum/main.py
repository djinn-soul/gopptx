"""Example showing improved chart type API using enum-style constants.

This example demonstrates the ONLY way to specify chart types:
- Using ChartType constants (required): ChartType.COLUMN, ChartType.LINE, etc.

No string values are accepted - type safety through constants only.
"""

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.charts import ChartType
from gopptx.schemas import Inches


def main() -> None:
    """Create a presentation with charts using the new ChartType enum."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Chart Type Examples") as prs:
        # =======================================================================
        # Example 1: Column Chart
        # =======================================================================
        print("\n" + "=" * 70)
        print("EXAMPLE 1: Column Chart (ChartType.COLUMN)")
        print("=" * 70)

        prs.add_slide("Column Chart", layout="title_only")

        print("\nCreating a column/bar chart with ChartType.COLUMN constant")
        prs.add_chart(
            0,
            ChartType.COLUMN,
            ["Q1", "Q2", "Q3", "Q4"],
            [12.0, 18.0, 15.0, 24.0],
            title="Revenue by Quarter",
            bounds=(Inches(0.45), Inches(1.2), Inches(4.5), Inches(3.0)),
        )

        # =======================================================================
        # Example 2: Line Chart
        # =======================================================================
        print("\n" + "=" * 70)
        print("EXAMPLE 2: Line Chart (ChartType.LINE)")
        print("=" * 70)

        prs.add_slide("Line Chart", layout="title_only")

        print("\nCreating a line chart with ChartType.LINE constant")
        prs.add_chart(
            1,
            ChartType.LINE,
            ["Jan", "Feb", "Mar", "Apr", "May"],
            [4.0, 7.0, 6.0, 9.0, 8.0],
            title="Monthly Trend",
            bounds=(Inches(0.45), Inches(1.2), Inches(9.0), Inches(3.0)),
        )

        # =======================================================================
        # Example 3: Pie Chart
        # =======================================================================
        print("\n" + "=" * 70)
        print("EXAMPLE 3: Pie Chart (ChartType.PIE)")
        print("=" * 70)

        prs.add_slide("Pie Chart", layout="title_only")

        print("\nCreating a pie chart with ChartType.PIE constant")
        prs.add_chart(
            2,
            ChartType.PIE,
            ["Direct", "Search", "Referral", "Social"],
            [45.0, 35.0, 12.0, 8.0],
            title="Traffic Sources",
            bounds=(Inches(0.45), Inches(1.2), Inches(4.0), Inches(3.0)),
        )

        print("\nCreating another pie chart with different data")
        prs.add_chart(
            2,
            ChartType.PIE,
            ["Product A", "Product B", "Product C", "Product D"],
            [25.0, 35.0, 25.0, 15.0],
            title="Sales by Product",
            bounds=(Inches(5.0), Inches(1.2), Inches(4.0), Inches(3.0)),
        )

        # =======================================================================
        # Example 4: Comparing Different Chart Types
        # =======================================================================
        print("\n" + "=" * 70)
        print("EXAMPLE 4: Same Data, Different Chart Types")
        print("=" * 70)

        prs.add_slide("Same Data, Different Visualizations", layout="title_only")

        data_categories = ["A", "B", "C", "D"]
        data_values = [30.0, 40.0, 20.0, 10.0]

        print("\nColumn view of data")
        prs.add_chart(
            3,
            ChartType.COLUMN,
            data_categories,
            data_values,
            title="Column View",
            bounds=(Inches(0.45), Inches(1.2), Inches(3.0), Inches(2.5)),
        )

        print("Line view of data")
        prs.add_chart(
            3,
            ChartType.LINE,
            data_categories,
            data_values,
            title="Line View",
            bounds=(Inches(3.55), Inches(1.2), Inches(3.0), Inches(2.5)),
        )

        print("Pie view of data")
        prs.add_chart(
            3,
            ChartType.PIE,
            data_categories,
            data_values,
            title="Pie View",
            bounds=(Inches(6.65), Inches(1.2), Inches(3.0), Inches(2.5)),
        )

        # =======================================================================
        # Example 5: Discovery - List All Available Chart Types
        # =======================================================================
        print("\n" + "=" * 70)
        print("EXAMPLE 5: Available Chart Types")
        print("=" * 70)

        all_types = ChartType.get_all()
        print("\nCurrently supported chart types:")
        for name, value in sorted(all_types.items()):
            print(f"  ChartType.{name:10} = '{value}'")

        print(f"\nTotal: {len(set(all_types.values()))} chart types available")

        # =======================================================================
        # Example 6: Error Handling
        # =======================================================================
        print("\n" + "=" * 70)
        print("EXAMPLE 6: Error Handling")
        print("=" * 70)

        prs.add_slide("Error Handling Demo", layout="title_only")

        # Show what happens with invalid input
        try:
            print("\nTrying invalid type: 'COLUMN' (string constant name not allowed)")
            prs.add_chart(
                5,
                "COLUMN",  # Error - must use ChartType.COLUMN, not the string "COLUMN"
                ["A", "B"],
                [1.0, 2.0],
                bounds=(100, 100, 400, 300),
            )
        except ValueError as e:
            print(f"Caught error (as expected): {str(e)[:80]}...")

        try:
            print("Trying invalid type: 'scatter' (not yet supported)")
            prs.add_chart(
                5,
                "scatter",  # Error - this chart type is not supported yet
                [1.0, 2.0, 3.0],
                [2.0, 4.0, 5.0],
                bounds=(100, 100, 400, 300),
            )
        except ValueError as e:
            print(f"Caught error (as expected): {str(e)[:80]}...")

        # =======================================================================
        # Save presentation
        # =======================================================================
        output_path = output_dir / "16-chart-type-enum.pptx"
        prs.save(str(output_path))
        print(f"\n\nPresentation saved to {output_path}")

        # =======================================================================
        # Summary
        # =======================================================================
        print("\n" + "=" * 70)
        print("SUMMARY")
        print("=" * 70)
        print(
            """
CHART TYPE API - ENUM CONSTANTS REQUIRED:

[+] ONLY CORRECT WAY
    from gopptx.presentation.charts import ChartType
    slide.add_chart(ChartType.COLUMN, ...)

    Benefits:
    - IDE autocomplete and type hints
    - Type-safe: no string typos possible
    - Discoverable: ChartType.get_all() lists all options
    - Self-documenting code

[!] NOT SUPPORTED
    - String constant names: "COLUMN", "LINE", "PIE" - NOT allowed
    - Raw values: "bar", "line", "pie" - NOT allowed
    - Must use the actual ChartType constants

CURRENTLY SUPPORTED CHART TYPES:

    ChartType.COLUMN   = "bar"   - Vertical bar/column chart
    ChartType.BAR      = "bar"   - Alias for COLUMN
    ChartType.LINE     = "line"  - Line chart
    ChartType.PIE      = "pie"   - Pie/circular chart

PLANNED FOR FUTURE (Not Yet Supported):

    Scatter/Bubble (ChartType.SCATTER, BUBBLE)
    Area (ChartType.AREA, AREA_STACKED)
    Radar (ChartType.RADAR, RADAR_FILLED)
    Stock (ChartType.STOCK_HLC, STOCK_OHLC)
    Combo (ChartType.COMBO)
    And more...

USAGE PATTERNS:

  # Create any supported chart type
  chart_id = slide.add_chart(
      ChartType.LINE,
      ["Q1", "Q2", "Q3"],
      [100, 150, 120],
      title="Sales Trend",
      bounds=(100, 100, 400, 300),
  )

  # List all available types
  for name, value in ChartType.get_all().items():
      print(f"ChartType.{name}")

  # Find value by constant name
  value = ChartType.get_by_name("COLUMN")  # -> "bar"

  # Validate a type value
  ChartType.validate("bar")  # OK - returns "bar"
  ChartType.validate("COLUMN")  # ERROR - use ChartType.COLUMN instead
        """
        )


if __name__ == "__main__":
    main()
