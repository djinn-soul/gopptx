"""Demonstrate adding bar and line charts to an existing presentation slide.

This example demonstrates:
- Creating a base presentation with a content slide
- Adding a bar chart with quarterly sales data
- Adding a line chart with monthly growth data to the same slide
- Chart positioning via bounds parameter
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.charts import ChartType
from gopptx.schemas import Inches


def main() -> None:
    """Create presentation demonstrating editor chart support."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Editor Chart Support") as prs:
        prs.add_slide("Chart Playground")
        slide_idx = prs.slide_count - 1

        # Bar chart: quarterly sales
        prs.add_chart(
            slide_idx,
            ChartType.BAR,
            categories=["Q1", "Q2", "Q3", "Q4"],
            values_or_series=[100, 200, 150, 300],
            title="Quarterly Sales",
            bounds=(Inches(0.5), Inches(1.5), Inches(4.5), Inches(3.5)),
        )

        # Line chart: monthly growth (offset to right)
        prs.add_chart(
            slide_idx,
            ChartType.LINE,
            categories=["Jan", "Feb", "Mar"],
            values_or_series=[5, 10, 8],
            title="Monthly Growth",
            bounds=(Inches(5), Inches(1.5), Inches(4.5), Inches(3.5)),
        )

        # Summary slide
        prs.add_bullet_slide(
            "Chart Support Features",
            [
                "add_chart() embeds charts directly onto existing slides",
                "BAR chart: categorical comparisons",
                "LINE chart: trend data over time",
                "COLUMN, PIE charts also available",
                "bounds=(x, y, w, h) controls chart position and size",
            ],
        )

        output_path = output_dir / "39_editor_chart_support_output.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("2 slides: Chart Playground (bar + line chart), Features summary")


if __name__ == "__main__":
    main()
