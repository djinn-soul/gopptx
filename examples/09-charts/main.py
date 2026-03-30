"""Demonstrate chart creation with multiple chart types.

This example demonstrates:
- Bar chart
- Column chart
- Line chart
- Pie chart
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.charts import ChartType
from gopptx.schemas import Inches

_CATEGORIES = ["Q1", "Q2", "Q3", "Q4"]
_VALUES = [20.0, 35.0, 28.0, 42.0]


def _add_bar_chart(prs: Presentation) -> None:
    """Add a horizontal bar chart slide."""
    prs.add_slide("Bar Chart")
    prs.add_chart(
        prs.slide_count - 1,
        ChartType.BAR,
        _CATEGORIES,
        _VALUES,
        title="Quarterly Revenue (Bar)",
        bounds=(Inches(1), Inches(1.5), Inches(8), Inches(5)),
    )


def _add_column_chart(prs: Presentation) -> None:
    """Add a vertical column chart slide."""
    prs.add_slide("Column Chart")
    prs.add_chart(
        prs.slide_count - 1,
        ChartType.COLUMN,
        _CATEGORIES,
        _VALUES,
        title="Quarterly Revenue (Column)",
        bounds=(Inches(1), Inches(1.5), Inches(8), Inches(5)),
    )


def _add_line_chart(prs: Presentation) -> None:
    """Add a line chart slide."""
    prs.add_slide("Line Chart")
    prs.add_chart(
        prs.slide_count - 1,
        ChartType.LINE,
        _CATEGORIES,
        _VALUES,
        title="Quarterly Trend (Line)",
        bounds=(Inches(1), Inches(1.5), Inches(8), Inches(5)),
    )


def _add_pie_chart(prs: Presentation) -> None:
    """Add a pie chart slide."""
    prs.add_slide("Pie Chart")
    prs.add_chart(
        prs.slide_count - 1,
        ChartType.PIE,
        _CATEGORIES,
        _VALUES,
        title="Revenue Distribution (Pie)",
        bounds=(Inches(1), Inches(1.5), Inches(8), Inches(5)),
    )


def main() -> None:
    """Create a presentation demonstrating multiple chart types."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Charts Demo") as prs:
        _add_bar_chart(prs)
        _add_column_chart(prs)
        _add_line_chart(prs)
        _add_pie_chart(prs)

        output_path = output_dir / "09-charts.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created 4 slides demonstrating chart types:")
    print("  - Bar (horizontal)")
    print("  - Column (vertical)")
    print("  - Line")
    print("  - Pie")


if __name__ == "__main__":
    main()
