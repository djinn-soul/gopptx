"""Dynamic Report Generation - Build PPTX reports from JSON, CSV, or database rows.

This example demonstrates:
- Loading metadata from JSON
- Parsing data from CSV (can easily be adapted to database queries)
- Calculating aggregates and metrics
- Creating professional report layouts with multiple slide types
- Using themes for consistent styling
- Building charts and tables from structured data
"""

from __future__ import annotations

import csv
import json
import operator
from io import StringIO

from gopptx import Presentation, ShapeType
from gopptx.constants import SIZE_16X9_HEIGHT, SIZE_16X9_WIDTH
from gopptx.presentation.charts import ChartType
from gopptx.presentation.slides import SlideLayoutType
from gopptx.presentation.theme import get_theme
from gopptx.schemas import Inches

REPORT_JSON = """
{
  "title": "Dynamic Report Generation",
  "subtitle": "Build full PPTX reports from JSON, CSV, or database rows",
  "period": "Q1 2026",
  "owner": "Finance Operations",
  "generated_by": "gopptx"
}
"""

REGION_CSV = """
region,revenue,orders,target
North,460000,128,420000
West,370000,102,350000
South,290000,96,310000
East,264000,88,250000
"""


def load_rows() -> tuple[dict[str, str], list[dict[str, object]]]:
    """Load JSON metadata and CSV data rows."""
    meta = json.loads(REPORT_JSON)

    rows: list[dict[str, object]] = []
    for raw in csv.DictReader(StringIO(REGION_CSV.strip())):
        revenue = float(raw["revenue"])
        orders = int(raw["orders"])
        target = float(raw["target"])
        rows.append({
            "region": raw["region"],
            "revenue": revenue,
            "orders": orders,
            "target": target,
            "attainment": revenue / target,
        })
    return meta, rows


def format_currency(value: float) -> str:
    """Format value as currency in millions."""
    return f"${value / 1_000_000:.2f}M"


def format_percent(value: float) -> str:
    """Format value as percentage."""
    return f"{value * 100:.0f}%"


def add_card(
    slide: object,
    bounds: tuple[float, float, float, float],
    text: str,
    fill: str,
    line: str,
) -> None:
    """Add a metric card to the slide."""
    x, y, w, h = bounds
    slide.add_shape(
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(x), Inches(y), Inches(w), Inches(h)),
        text=text,
        properties={
            "fill": {"solid": fill},
            "line": {"color": line, "width_emu": 12700},
        },
    )


# Load data
meta, rows = load_rows()
total_revenue = sum(row["revenue"] for row in rows)
total_orders = sum(row["orders"] for row in rows)
average_attainment = sum(row["attainment"] for row in rows) / len(rows)
max_region = max(rows, key=operator.itemgetter("revenue"))

# Create presentation
with Presentation.new(str(meta["title"])) as p:
    # Set slide size and apply theme
    p.set_slide_size(SIZE_16X9_WIDTH, SIZE_16X9_HEIGHT)
    p.apply_theme(get_theme("aurora"))

    # ===== Slide 1: Overview with metrics and chart =====
    p.update_slide(0, layout=SlideLayoutType.BLANK)
    slide = p.slides[0]

    # Title and subtitle
    slide.add_textbox(
        Inches(0.8), Inches(0.35), Inches(6.8), Inches(0.5), text=str(meta["title"])
    )
    slide.add_textbox(
        Inches(0.8), Inches(0.82), Inches(6.8), Inches(0.32), text=str(meta["subtitle"])
    )

    # Metric cards
    add_card(
        slide,
        (0.8, 1.35, 1.8, 1.0),
        f"Revenue\n{format_currency(total_revenue)}",
        "EEF4FB",
        "A9C4E2",
    )
    add_card(
        slide, (2.75, 1.35, 1.8, 1.0), f"Orders\n{total_orders}", "E8F5E9", "B8D5B8"
    )
    add_card(
        slide,
        (0.8, 2.55, 1.8, 1.0),
        f"Avg attainment\n{format_percent(average_attainment)}",
        "FCE4D6",
        "E8B89C",
    )
    add_card(
        slide,
        (2.75, 2.55, 1.8, 1.0),
        f"Top region\n{max_region['region']}",
        "FFF2CC",
        "E0C75C",
    )

    # Chart container background
    slide.add_shape(
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(4.8), Inches(1.25), Inches(4.0), Inches(3.05)),
        properties={
            "fill": {"solid": "FFFFFF"},
            "line": {"color": "C9D3E0", "width_emu": 12700},
        },
    )

    # Bar chart with revenue by region (using ChartType enum)
    slide.add_chart(
        ChartType.BAR,  # Using enum instead of string
        [str(row["region"]) for row in rows],
        [float(row["revenue"]) / 1000 for row in rows],
        title="Revenue by Region (USD K)",
        bounds=(Inches(4.95), Inches(1.45), Inches(3.7), Inches(2.65)),
    )

    # Description text
    slide.add_textbox(
        Inches(0.8),
        Inches(4.0),
        Inches(8.0),
        Inches(0.5),
        text="JSON metadata drives the title block. CSV rows feed the chart. "
        "Database query rows can use the same dict schema.",
    )

    # ===== Slide 2: Detail table =====
    p.add_slide("Regional Detail Table", layout=SlideLayoutType.BLANK)
    detail = p.slides[1]

    # Title and description
    detail.add_textbox(
        Inches(0.8),
        Inches(0.35),
        Inches(5.8),
        Inches(0.45),
        text="Regional Detail Table",
    )
    detail.add_textbox(
        Inches(0.8),
        Inches(0.82),
        Inches(8.4),
        Inches(0.3),
        text="The same row structure can come from JSON records, CSV files, "
        "or SQL query results.",
    )

    # Build table data
    detail_rows = [["Region", "Revenue", "Orders", "Target", "Attainment"]]
    detail_rows.extend(
        [
            str(row["region"]),
            format_currency(row["revenue"]),
            str(row["orders"]),
            format_currency(row["target"]),
            format_percent(row["attainment"]),
        ]
        for row in rows
    )

    # Add table
    p.add_table_from_rows(
        1,
        detail_rows,
        (Inches(0.8), Inches(1.35), Inches(8.6), Inches(2.95)),
        first_row=True,
        band_row=True,
    )

    # Footer with metadata
    detail.add_textbox(
        Inches(0.8),
        Inches(4.45),
        Inches(8.6),
        Inches(0.3),
        text=f"Generated by {meta['generated_by']} for {meta['owner']} "
        f"({meta['period']}).",
    )

    # Save
    p.save("examples/output/08-dynamic-report-generation.pptx")
    print("Presentation created: examples/output/08-dynamic-report-generation.pptx")

    # ===== Summary =====
    print("\n" + "=" * 70)
    print("DYNAMIC REPORT GENERATION")
    print("=" * 70)
    print("\nPattern: Data-Driven Presentations")
    print("\nBenefits:")
    print("  [+] JSON metadata for configuration")
    print("  [+] CSV/database rows for content")
    print("  [+] Calculated metrics and aggregates")
    print("  [+] Professional layout with cards, charts, and tables")
    print("  [+] Theme consistency across all slides")
    print("\nData Sources:")
    print("  - JSON: Meta information (title, period, owner)")
    print("  - CSV: Regional data (revenue, orders, target)")
    print("  - Calculations: Totals, averages, rankings")
    print("\nSlide Composition:")
    print("  Slide 1 (Overview):")
    print("    - Title from JSON")
    print("    - Four metric cards (Revenue, Orders, Attainment, Top Region)")
    print("    - Bar chart showing revenue by region")
    print("  Slide 2 (Details):")
    print("    - Table with all row data")
    print("    - Footer with generation metadata")
    print("\nKey APIs Used:")
    print("  - prs.set_slide_size() - Set 16:9 format")
    print("  - prs.apply_theme() - Aurora theme")
    print("  - SlideLayoutType.BLANK - Enum for blank layout")
    print("  - ChartType.BAR - Enum for bar chart")
    print("  - add_table_from_rows() - Table from list of lists")
    print("  - add_card() - Reusable metric card builder")
    print("=" * 70 + "\n")
