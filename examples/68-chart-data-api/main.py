"""Demonstrate chart data builders with multiple chart types.

This example demonstrates:
- Bar, column, line, and pie charts via add_chart()
- Multiple category/series patterns (quarterly revenue, market share, etc.)
- Using ChartType constants for type-safe chart creation
- Varying bounds placement with Inches()
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.charts import ChartType
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches

QUARTERS = ["Q1", "Q2", "Q3", "Q4"]


def _add_bar_chart_slide(prs: Presentation) -> None:
    """Quarterly revenue bar chart."""
    idx = prs.add_slide(
        "Bar Chart - Quarterly Revenue", layout=SlideLayoutType.TITLE_ONLY
    ).index
    prs.add_chart(
        idx,
        ChartType.BAR,
        QUARTERS,
        [42, 55, 61, 73],
        title="Quarterly Revenue ($k)",
        bounds=(Inches(0.5), Inches(1.2), Inches(9), Inches(5.5)),
    )


def _add_column_chart_slide(prs: Presentation) -> None:
    """Quarterly costs column chart."""
    idx = prs.add_slide(
        "Column Chart - Quarterly Costs", layout=SlideLayoutType.TITLE_ONLY
    ).index
    prs.add_chart(
        idx,
        ChartType.COLUMN,
        QUARTERS,
        [18, 24, 30, 41],
        title="Quarterly Costs ($k)",
        bounds=(Inches(0.5), Inches(1.2), Inches(9), Inches(5.5)),
    )


def _add_line_chart_slide(prs: Presentation) -> None:
    """Units sold line chart."""
    idx = prs.add_slide(
        "Line Chart - Units Sold", layout=SlideLayoutType.TITLE_ONLY
    ).index
    prs.add_chart(
        idx,
        ChartType.LINE,
        QUARTERS,
        [30, 45, 55, 68],
        title="Units Sold per Quarter",
        bounds=(Inches(0.5), Inches(1.2), Inches(9), Inches(5.5)),
    )


def _add_pie_chart_slide(prs: Presentation) -> None:
    """Market share pie chart."""
    idx = prs.add_slide(
        "Pie Chart - Market Share", layout=SlideLayoutType.TITLE_ONLY
    ).index
    prs.add_chart(
        idx,
        ChartType.PIE,
        ["Product A", "Product B", "Product C", "Product D"],
        [35, 25, 20, 20],
        title="Market Share by Product",
        bounds=(Inches(1.0), Inches(1.2), Inches(8), Inches(5.5)),
    )


def _add_chart_summary_slide(prs: Presentation) -> None:
    """Summary slide listing available ChartType constants."""
    prs.add_bullet_slide(
        "ChartType Constants",
        [
            f"ChartType.BAR    = {ChartType.BAR!r}",
            f"ChartType.COLUMN = {ChartType.COLUMN!r}",
            f"ChartType.LINE   = {ChartType.LINE!r}",
            f"ChartType.PIE    = {ChartType.PIE!r}",
            "All use: prs.add_chart(slide_idx, ChartType.X, categories, values, title, bounds)",
        ],
    )


def main() -> None:
    """Create presentation demonstrating chart data API."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "68_chart_data_api.pptx"

    with Presentation.new("Chart Data API Demo") as prs:
        _add_bar_chart_slide(prs)
        _add_column_chart_slide(prs)
        _add_line_chart_slide(prs)
        _add_pie_chart_slide(prs)
        _add_chart_summary_slide(prs)

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Demonstrated: add_chart with ChartType.BAR, COLUMN, LINE, PIE")
    print("  across quarterly revenue, costs, units, and market-share datasets")


if __name__ == "__main__":
    main()
