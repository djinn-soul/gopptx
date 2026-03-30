"""Demonstrate chart formatting options across multiple chart types.

This example demonstrates:
- Bar, column, line, and pie charts with varied data sets
- Using add_chart() with title and bounds on TITLE_ONLY slides
- Multi-slide chart gallery with monthly, quarterly, and categorical data
- Legend position awareness via ChartType constants
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.charts import ChartType
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches

MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun"]
QUARTERS = ["Q1", "Q2", "Q3", "Q4"]
CHART_BOUNDS = (Inches(0.5), Inches(1.2), Inches(9), Inches(5.5))


def _add_monthly_bar(prs: Presentation) -> None:
    """Monthly sales bar chart with full formatting options."""
    idx = prs.add_slide("Bar Chart - Full Options", layout=SlideLayoutType.TITLE_ONLY)
    prs.add_chart(
        idx,
        ChartType.BAR,
        MONTHS,
        [12, 19, 15, 22, 28, 24],
        title="Monthly Sales (Bar Chart - Full Options)",
        bounds=CHART_BOUNDS,
    )


def _add_monthly_column(prs: Presentation) -> None:
    """Monthly revenue column chart."""
    idx = prs.add_slide(
        "Column Chart - Monthly Revenue", layout=SlideLayoutType.TITLE_ONLY
    )
    prs.add_chart(
        idx,
        ChartType.COLUMN,
        MONTHS,
        [5, 7, 6, 10, 12, 9],
        title="Revenue Trend ($k)",
        bounds=CHART_BOUNDS,
    )


def _add_monthly_line(prs: Presentation) -> None:
    """Monthly costs line chart."""
    idx = prs.add_slide("Line Chart - Monthly Costs", layout=SlideLayoutType.TITLE_ONLY)
    prs.add_chart(
        idx,
        ChartType.LINE,
        MONTHS,
        [8, 11, 9, 14, 16, 13],
        title="Monthly Costs ($k)",
        bounds=CHART_BOUNDS,
    )


def _add_market_share_pie(prs: Presentation) -> None:
    """Market share pie chart."""
    idx = prs.add_slide("Pie Chart - Market Share", layout=SlideLayoutType.TITLE_ONLY)
    prs.add_chart(
        idx,
        ChartType.PIE,
        ["Product A", "Product B", "Product C", "Product D", "Other"],
        [30, 25, 20, 15, 10],
        title="Market Share by Product",
        bounds=(Inches(1.0), Inches(1.2), Inches(8), Inches(5.5)),
    )


def _add_quarterly_bar(prs: Presentation) -> None:
    """Quarterly units column chart."""
    idx = prs.add_slide(
        "Column Chart - Quarterly Units", layout=SlideLayoutType.TITLE_ONLY
    )
    prs.add_chart(
        idx,
        ChartType.COLUMN,
        QUARTERS,
        [42, 55, 61, 73],
        title="Quarterly Units Sold",
        bounds=CHART_BOUNDS,
    )


def _add_radar_reference_slide(prs: Presentation) -> None:
    """Reference slide for chart axis/legend constants."""
    prs.add_bullet_slide(
        "Chart API Notes",
        [
            "add_chart(slide_idx, ChartType, categories, values, title, bounds)",
            "ChartType.BAR    - horizontal bar chart",
            "ChartType.COLUMN - vertical bar (column) chart",
            "ChartType.LINE   - line chart",
            "ChartType.PIE    - pie chart",
            "bounds = (x, y, width, height) as Inches()",
        ],
    )


def main() -> None:
    """Create presentation demonstrating chart formatting options."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "69_chart_api.pptx"

    with Presentation.new("Chart API Demo") as prs:
        _add_monthly_bar(prs)
        _add_monthly_column(prs)
        _add_monthly_line(prs)
        _add_market_share_pie(prs)
        _add_quarterly_bar(prs)
        _add_radar_reference_slide(prs)

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print(
        "Demonstrated: add_chart with BAR, COLUMN, LINE, PIE across monthly/quarterly data"
    )


if __name__ == "__main__":
    main()
