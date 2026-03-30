"""Demonstrate enum constants available in the gopptx Python API.

This example demonstrates:
- ShapeType constants (RECTANGLE, ELLIPSE, ROUNDED_RECTANGLE, TRIANGLE)
- ConnectorType constants (STRAIGHT, ELBOW, CURVED)
- SlideLayoutType constants (BLANK, TITLE_ONLY, TITLE_CONTENT)
- ChartType constants (BAR, COLUMN, LINE, PIE)
- Shapes gallery using each ShapeType constant
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.constants import ConnectorType, ShapeType
from gopptx.presentation.charts import ChartType
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches

_SHAPE_COLORS = ["4472C4", "C0504D", "9BBB59", "F79646", "8064A2", "4BACC6"]


def _add_core_shape_types_slide(prs: Presentation) -> None:
    """Gallery of core shape types."""
    idx = prs.add_slide("Core Shape Type Constants", layout=SlideLayoutType.TITLE_ONLY)
    types = [
        (ShapeType.RECTANGLE, "Rect"),
        (ShapeType.ROUNDED_RECTANGLE, "RoundRect"),
        (ShapeType.ELLIPSE, "Ellipse"),
        (ShapeType.TRIANGLE, "Triangle"),
    ]
    for i, (shape_type, label) in enumerate(types):
        col = i % 2
        row = i // 2
        prs.add_shape(
            idx,
            shape_type,
            (
                Inches(col * 4.5 + 0.5),
                Inches(row * 2.0 + 1.5),
                Inches(3.5),
                Inches(1.5),
            ),
            text=label,
            properties={"fill_color": _SHAPE_COLORS[i]},
        )


def _add_connector_types_slide(prs: Presentation) -> None:
    """Enumerate connector type constants."""
    idx = prs.add_slide("Connector Type Constants", layout=SlideLayoutType.TITLE_ONLY)
    connectors = [
        (ConnectorType.STRAIGHT, Inches(0.5), Inches(2.0), Inches(3.5), Inches(2.0)),
        (ConnectorType.ELBOW, Inches(0.5), Inches(3.0), Inches(4.5), Inches(4.5)),
        (ConnectorType.CURVED, Inches(5.0), Inches(2.0), Inches(9.0), Inches(4.5)),
    ]
    labels = ["STRAIGHT", "ELBOW", "CURVED"]
    for _label, (ctype, x1, y1, x2, y2) in zip(labels, connectors):
        prs.add_connector(idx, ctype, x1, y1, x2, y2)
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(1.2),
        Inches(9),
        Inches(0.5),
        text="ConnectorType.STRAIGHT / ELBOW / CURVED",
    )


def _add_layout_constants_slide(prs: Presentation) -> None:
    """Enumerate slide layout type constants."""
    prs.add_bullet_slide(
        "SlideLayoutType Constants",
        [
            f"BLANK          = {SlideLayoutType.BLANK!r}",
            f"TITLE_ONLY     = {SlideLayoutType.TITLE_ONLY!r}",
            f"TITLE_CONTENT  = {SlideLayoutType.TITLE_CONTENT!r}",
            "Use layout= parameter in prs.add_slide(title, layout=...)",
        ],
    )


def _add_chart_type_constants_slide(prs: Presentation) -> None:
    """Enumerate ChartType constants."""
    prs.add_bullet_slide(
        "ChartType Constants",
        [
            f"ChartType.BAR    = {ChartType.BAR!r}",
            f"ChartType.COLUMN = {ChartType.COLUMN!r}",
            f"ChartType.LINE   = {ChartType.LINE!r}",
            f"ChartType.PIE    = {ChartType.PIE!r}",
            "Pass to prs.add_chart(slide_idx, ChartType.X, ...)",
        ],
    )


def _add_shape_type_constants_slide(prs: Presentation) -> None:
    """List ShapeType and ConnectorType string values."""
    prs.add_bullet_slide(
        "ShapeType & ConnectorType Values",
        [
            f"ShapeType.RECTANGLE         = {ShapeType.RECTANGLE!r}",
            f"ShapeType.ELLIPSE           = {ShapeType.ELLIPSE!r}",
            f"ShapeType.ROUNDED_RECTANGLE = {ShapeType.ROUNDED_RECTANGLE!r}",
            f"ShapeType.TRIANGLE          = {ShapeType.TRIANGLE!r}",
            f"ConnectorType.STRAIGHT      = {ConnectorType.STRAIGHT!r}",
            f"ConnectorType.ELBOW         = {ConnectorType.ELBOW!r}",
            f"ConnectorType.CURVED        = {ConnectorType.CURVED!r}",
        ],
    )


def main() -> None:
    """Create presentation demonstrating enum constants."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "78_enum_api.pptx"

    with Presentation.new("Enum API Demo") as prs:
        _add_core_shape_types_slide(prs)
        _add_connector_types_slide(prs)
        _add_layout_constants_slide(prs)
        _add_chart_type_constants_slide(prs)
        _add_shape_type_constants_slide(prs)

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print(
        "Demonstrated: ShapeType, ConnectorType, SlideLayoutType, ChartType constants"
    )
    print("  with a shapes gallery and enumeration slides")


if __name__ == "__main__":
    main()
