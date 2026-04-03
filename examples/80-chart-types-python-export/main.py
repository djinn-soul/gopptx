"""Create a Python chart sample deck and export PPTX + PDF with all chart types."""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.charts import ChartType
from gopptx.presentation.export.export_mixin import PDFOptions
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches

OUTPUT_DIR = Path("examples/output")
PPTX_PATH = OUTPUT_DIR / "80_chart_types_python_export.pptx"
PDF_PATH = OUTPUT_DIR / "80_chart_types_python_export.pdf"
CHART_BOUNDS = (Inches(0.8), Inches(1.4), Inches(8.5), Inches(4.6))

ALL_CHART_TYPES = ChartType.get_all()
UNIQUE_CHART_VALUES = sorted(set(ALL_CHART_TYPES.values()))
CATEGORY_CHART_VALUES = [ct for ct in UNIQUE_CHART_VALUES if ct != ChartType.COMBO]


def add_chart_slide(prs: Presentation, chart_type: str) -> None:
    """Add one chart slide for a chart type."""
    prs.add_slide(f"{chart_type} Chart", layout=SlideLayoutType.TITLE_ONLY)
    slide_index = prs.slide_count - 1

    prs.add_chart(
        slide_index,
        chart_type,
        ["Q1", "Q2", "Q3", "Q4"],
        [14, 21, 18, 27],
        title=f"{chart_type} Demo",
        bounds=CHART_BOUNDS,
    )


def add_combo_slide(prs: Presentation) -> None:
    """Add one combo chart slide."""
    prs.add_slide("COMBO Chart", layout=SlideLayoutType.TITLE_ONLY)
    slide_index = prs.slide_count - 1
    prs.add_combo_chart(
        slide_index,
        ["Q1", "Q2", "Q3", "Q4"],
        bar_series=[{"name": "Revenue", "values": [180, 220, 210, 260]}],
        line_series=[{"name": "Growth %", "values": [8, 11, 10, 14]}],
        title="COMBO Demo",
        bounds=CHART_BOUNDS,
    )


def add_chart_surface_reference(prs: Presentation) -> None:
    """Add a reference slide listing all chart constants exposed in Python."""
    constant_lines = [
        f"ChartType.{name} = {value}" for name, value in ALL_CHART_TYPES.items()
    ]
    prs.add_bullet_slide(
        "Python ChartType Surface",
        [
            f"Named constants: {len(ALL_CHART_TYPES)}",
            f"Unique chart kinds: {len(UNIQUE_CHART_VALUES)}",
            *constant_lines,
        ],
    )


def main() -> None:
    """Build sample deck and export both PPTX and PDF."""
    OUTPUT_DIR.mkdir(exist_ok=True)

    with Presentation.new("Python Chart Types Export Demo") as prs:
        prs.add_bullet_slide(
            "Chart Types Demo",
            [
                "This deck demonstrates chart creation from Python API.",
                "It outputs both PPTX and PDF artifacts.",
            ],
        )

        for chart_type in CATEGORY_CHART_VALUES:
            add_chart_slide(prs, chart_type)

        add_combo_slide(prs)
        add_chart_surface_reference(prs)

        prs.save(str(PPTX_PATH))
        pdf_output = prs.save_as_pdf(str(PDF_PATH), PDFOptions(driver="native"))

    print(f"Saved PPTX: {PPTX_PATH}")
    print(f"Saved PDF: {pdf_output}")


if __name__ == "__main__":
    main()
